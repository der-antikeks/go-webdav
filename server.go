package webdav

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func Handler(root FileSystem) http.Handler {
	return &Server{Fs: root}
}

type Server struct {
	// trimmed path prefix
	TrimPrefix string

	// files are readonly?
	ReadOnly bool

	// generate directory listings?
	Listings bool

	// access to a collection of named files
	Fs FileSystem
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("DAV:", r.RemoteAddr, r.Method, r.URL)

	switch r.Method {
	case "OPTIONS":
		s.doOptions(w, r)

	case "GET":
		s.doGet(w, r)
	case "HEAD":
		s.doHead(w, r)
	case "POST":
		s.doPost(w, r)
	case "DELETE":
		s.doDelete(w, r)
	case "PUT":
		s.doPut(w, r)

	case "PROPFIND":
		s.doPropfind(w, r)
	case "PROPPATCH":
		s.doProppatch(w, r)
	case "MKCOL":
		s.doMkcol(w, r)
	case "COPY":
		s.doCopy(w, r)
	case "MOVE":
		s.doMove(w, r)

	case "LOCK":
		s.doLock(w, r)
	case "UNLOCK":
		s.doUnlock(w, r)

	default:
		log.Println("DAV:", "unknown method", r.Method)
		w.WriteHeader(StatusBadRequest)
	}
}

func (s *Server) methodsAllowed(path string) string {
	if !s.pathExists(path) {
		return "OPTIONS, MKCOL, PUT, LOCK"
	}

	allowed := "OPTIONS, GET, HEAD, POST, DELETE, TRACE, PROPPATCH, COPY, MOVE, LOCK, UNLOCK"

	if s.Listings {
		allowed += ", PROPFIND"
	}

	if s.pathIsDirectory(path) {
		allowed += ", PUT"
	}

	return allowed
}

// convert request url to path
func (s *Server) url2path(u *url.URL) string {
	if u.Path == "" {
		return "/"
	}

	if p := strings.TrimPrefix(u.Path, s.TrimPrefix); len(p) < len(u.Path) {
		return strings.Trim(p, "/")
	}

	return "/"
}

// convert path to url
func (s *Server) path2url(p string) *url.URL {
	return &url.URL{Path: "/" + s.TrimPrefix + "/" + p}
}

// does path exists?
func (s *Server) pathExists(path string) bool {
	f, err := s.Fs.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	return true
}

// is path a directory?
func (s *Server) pathIsDirectory(path string) bool {
	f, err := s.Fs.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return false
	}

	return fi.IsDir()
}

func (s *Server) directoryContents(path string) []string {
	f, err := s.Fs.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	fi, err := f.Readdir(0)
	if err != nil {
		return nil
	}

	ret := make([]string, len(fi))
	for k, i := range fi {
		name := i.Name()
		if i.IsDir() {
			name += "/"
		}
		ret[k] = name
	}

	return ret
}

// is path in request locked?
func (s *Server) isLockedRequest(r *http.Request) bool {
	return s.isLocked(
		s.url2path(r.URL),
		r.Header.Get("If")+r.Header.Get("Lock-Token"))
}

// is path locked?
func (s *Server) isLocked(path, ifHeader string) bool {
	// TODO
	return false
}

func (s *Server) lockResource(path string) {
	// TODO
}

func (s *Server) unlockResource(path string) {
	// TODO
}

// The PROPFIND method retrieves properties defined on the resource identified by the Request-URI
// http://www.webdav.org/specs/rfc4918.html#METHOD_PROPFIND
func (s *Server) doPropfind(w http.ResponseWriter, r *http.Request) {
	if !s.Listings {
		w.Header().Set("Allow", s.methodsAllowed(s.url2path(r.URL)))
		w.WriteHeader(StatusMethodNotAllowed)
		return
	}

	/*
		TODO:
			return only directory and ics-file
			of current user
	*/

	depth := r.Header.Get("Depth")
	switch depth {
	case "0", "1":
	case "", "infinity":
		// treat as infinity if no depth header was included
		// disable infinity for performance and security concerns
		// http://www.webdav.org/specs/rfc4918.html#rfc.section.9.1.1
		w.WriteHeader(StatusForbidden)
		return
	default:
		log.Println("DAV:", "invalid depth header", depth)
		w.WriteHeader(StatusBadRequest)
		return
	}

	var propnames bool
	var properties []string
	var includes []string

	if r.ContentLength > 0 {
		propfind, err := NodeFromXml(r.Body)
		if err != nil {
			log.Println("DAV:", "invalid propfind xml", err)
			w.WriteHeader(StatusBadRequest)
			return
		}

		// find by property
		// http://www.webdav.org/specs/rfc4918.html#dav.properties
		if propfind.HasChildren("prop") {
			for _, p := range propfind.GetChildrens("prop") {
				properties = append(properties, p.Name.Local)
			}
		}

		// find property names
		if propfind.HasChildren("propname") {
			propnames = true
		}

		// find all properties
		if propfind.HasChildren("allprop") {
			if propfind.HasChildren("include") {
				for _, i := range propfind.GetChildrens("include") {
					includes = append(includes, i.Name.Local)
				}
			}
		}
	}

	path := s.url2path(r.URL)
	if !s.pathExists(path) {
		w.WriteHeader(StatusNotFound)
		return
	}

	paths := []string{path}
	if depth == "1" {
		// fetch all files if directory
		// respect []includes
	}

	log.Println("propnames", propnames)

	w.WriteHeader(StatusMulti)
	w.Header().Set("Content-Type", "application/xml; charset=UTF-8")
	for _, p := range paths {
		// test locks/ authorization
		// if properties, show only given properties, else all
		// if propnames, return names of properties, else names and values
		log.Println(p)
	}

	// TODO: propfind
	w.WriteHeader(StatusNotImplemented)
}

// http://www.webdav.org/specs/rfc4918.html#METHOD_PROPPATCH
func (s *Server) doProppatch(w http.ResponseWriter, r *http.Request) {
	if s.ReadOnly {
		w.WriteHeader(StatusForbidden)
		return
	}

	if s.isLockedRequest(r) {
		w.WriteHeader(StatusLocked)
		return
	}

	// TODO: proppatch
	w.WriteHeader(StatusNotImplemented)
}

// http://www.webdav.org/specs/rfc4918.html#METHOD_MKCOL
func (s *Server) doMkcol(w http.ResponseWriter, r *http.Request) {
	if s.ReadOnly {
		w.WriteHeader(StatusForbidden)
		return
	}

	if s.isLockedRequest(r) {
		w.WriteHeader(StatusLocked)
		return
	}

	path := s.url2path(r.URL)
	if s.pathExists(path) {
		w.Header().Set("Allow", s.methodsAllowed(s.url2path(r.URL)))
		w.WriteHeader(StatusMethodNotAllowed)
		return
	}

	// MKCOL may contain messagebody, precise behavior is undefined
	if r.ContentLength > 0 {
		_, err := NodeFromXml(r.Body)
		if err != nil {
			w.WriteHeader(StatusBadRequest)
			return
		}

		w.WriteHeader(StatusUnsupportedMediaType)
		return
	}

	if err := s.Fs.Mkdir(path); err != nil {
		w.WriteHeader(StatusConflict)
		return
	}

	w.WriteHeader(StatusCreated)
	s.unlockResource(path)
}

// http://www.webdav.org/specs/rfc4918.html#rfc.section.9.4
func (s *Server) doGet(w http.ResponseWriter, r *http.Request) {
	s.serveResource(w, r, true)
}

// http://www.webdav.org/specs/rfc4918.html#rfc.section.9.4
func (s *Server) doHead(w http.ResponseWriter, r *http.Request) {
	s.serveResource(w, r, false)
}

// http://www.webdav.org/specs/rfc4918.html#METHOD_POST
func (s *Server) doPost(w http.ResponseWriter, r *http.Request) {
	s.doGet(w, r)
}

func (s *Server) serveResource(w http.ResponseWriter, r *http.Request, serveContent bool) {
	path := s.url2path(r.URL)

	f, err := s.Fs.Open(path)
	if err != nil {
		http.Error(w, r.RequestURI, StatusNotFound)
		return
	}
	defer f.Close()

	// TODO: what if path is collection?

	fi, err := f.Stat()
	if err != nil {
		http.Error(w, r.RequestURI, StatusNotFound)
	}
	modTime := fi.ModTime()

	if serveContent {
		http.ServeContent(w, r, path, modTime, f)
	} else {
		// TODO: better way to send only head
		http.ServeContent(w, r, path, modTime, emptyFile{})
	}
}

// http://www.webdav.org/specs/rfc4918.html#METHOD_DELETE
func (s *Server) doDelete(w http.ResponseWriter, r *http.Request) {
	if s.ReadOnly {
		w.WriteHeader(StatusForbidden)
		return
	}

	if s.isLockedRequest(r) {
		w.WriteHeader(StatusLocked)
		return
	}

	s.deleteResource(s.url2path(r.URL), w, r, true)
}

func (s *Server) deleteResource(path string, w http.ResponseWriter, r *http.Request, setStatus bool) bool {
	ifHeader := r.Header.Get("If")
	lockToken := r.Header.Get("Lock-Token")

	if s.isLocked(path, ifHeader+lockToken) {
		w.WriteHeader(StatusLocked)
		return false
	}

	if !s.pathExists(path) {
		w.WriteHeader(StatusNotFound)
		return false
	}

	if !s.pathIsDirectory(path) {
		if err := s.Fs.Remove(path); err != nil {
			w.WriteHeader(StatusInternalServerError)
			return false
		}
	} else {
		// http://www.webdav.org/specs/rfc4918.html#delete-collections
		errors := map[string]int{}
		s.deleteCollection(path, w, r, errors)

		if err := s.Fs.Remove(path); err != nil {
			errors[path] = StatusInternalServerError
		}

		if len(errors) != 0 {
			// send multistatus
			abs := r.RequestURI

			buf := new(bytes.Buffer)
			buf.WriteString(`<?xml version="1.0" encoding="utf-8"?>`)
			buf.WriteString(`<multistatus xmlns='DAV:'>`)

			for p, e := range errors {
				buf.WriteString(`<response>`)
				buf.WriteString(`<href>` + abs + p + `</href>`)
				buf.WriteString(`<status>HTTP/1.1 ` + string(e) + ` ` + StatusText(e) + `</status>`)
				buf.WriteString(`</response>`)
			}

			buf.WriteString(`</multistatus>`)

			w.WriteHeader(StatusMulti)
			w.Header().Set("Content-Length", string(buf.Len()))
			w.Header().Set("Content-Type", "application/xml; charset=utf-8")
			buf.WriteTo(w)

			return false
		}
	}

	if setStatus {
		w.WriteHeader(StatusNoContent)
	}
	return true
}

func (s *Server) deleteCollection(path string, w http.ResponseWriter, r *http.Request, errors map[string]int) {
	ifHeader := r.Header.Get("If")
	lockToken := r.Header.Get("Lock-Token")

	for _, p := range s.directoryContents(path) {
		p = path + "/" + p

		if s.isLocked(p, ifHeader+lockToken) {
			errors[p] = StatusLocked
		} else {
			if s.pathIsDirectory(p) {
				s.deleteCollection(p, w, r, errors)
			}

			if err := s.Fs.Remove(p); err != nil {
				errors[p] = StatusInternalServerError
			}
		}
	}

}

// http://www.webdav.org/specs/rfc4918.html#METHOD_PUT
func (s *Server) doPut(w http.ResponseWriter, r *http.Request) {
	if s.ReadOnly {
		w.WriteHeader(StatusForbidden)
		return
	}

	if s.isLockedRequest(r) {
		w.WriteHeader(StatusLocked)
		return
	}

	path := s.url2path(r.URL)

	if s.pathIsDirectory(path) {
		// use MKCOL instead
		w.WriteHeader(StatusMethodNotAllowed)
		return
	}

	exists := s.pathExists(path)

	// TODO: content range / partial put

	// truncate file if exists
	file, err := s.Fs.Create(path)
	if err != nil {
		w.WriteHeader(StatusConflict)
		return
	}
	defer file.Close()

	if _, err := io.Copy(file, r.Body); err != nil {
		w.WriteHeader(StatusConflict)
	} else {
		if exists {
			w.WriteHeader(StatusNoContent)
		} else {
			w.WriteHeader(StatusCreated)
		}
	}

	s.unlockResource(path)
}

// http://www.webdav.org/specs/rfc4918.html#METHOD_COPY
func (s *Server) doCopy(w http.ResponseWriter, r *http.Request) {
	if s.ReadOnly {
		w.WriteHeader(StatusForbidden)
		return
	}

	s.copyResource(w, r)
}

// http://www.webdav.org/specs/rfc4918.html#METHOD_MOVE
func (s *Server) doMove(w http.ResponseWriter, r *http.Request) {
	if s.ReadOnly {
		w.WriteHeader(StatusForbidden)
		return
	}

	if s.isLockedRequest(r) {
		w.WriteHeader(StatusLocked)
		return
	}

	if s.copyResource(w, r) {
		// TODO: duplicate http-header sent?
		s.deleteResource(s.url2path(r.URL), w, r, false)
	}
}

func (s *Server) copyResource(w http.ResponseWriter, r *http.Request) bool {
	dest := r.Header.Get("Destination")
	if dest == "" {
		w.WriteHeader(StatusBadRequest)
		return false
	}

	d, err := url.Parse(dest)
	if err != nil {
		w.WriteHeader(StatusBadRequest)
		return false
	}
	// TODO: normalize dest?
	dest = s.url2path(d)
	source := s.url2path(r.URL)
	log.Println("DAV:", "copy", dest, source)

	// source equals destination
	if source == dest {
		w.WriteHeader(StatusForbidden)
		return false
	}

	// quick'n'dirty destination must be same server/namespace as source
	if d.Host != r.Host ||
		!strings.HasPrefix(d.Path, s.TrimPrefix) ||
		!strings.HasPrefix(r.URL.Path, s.TrimPrefix) {

		w.WriteHeader(StatusBadGateway)
		return false
	}

	overwrite := r.Header.Get("Overwrite") != "F"
	exists := s.pathExists(dest)

	log.Println("DAV:", "overwrite", r.Header.Get("Overwrite"), overwrite)

	if overwrite {
		if exists {
			if !s.deleteResource(dest, w, r, true) {
				// TODO: http status code?

				log.Println("DAV:", "overwrite existing", "failed")
				return false
			}
			log.Println("DAV:", "overwrite existing", "success")
		}
	} else {
		if exists {
			w.WriteHeader(StatusPreconditionFailed)
			return false
		}
	}

	// TODO: copy resource
	if !s.pathIsDirectory(source) || r.Header.Get("Depth") == "0" {
		if err := s.CopyFile(source, dest); err != nil {
			w.WriteHeader(StatusInternalServerError)
			return false
		}
	} else {
		// http://www.webdav.org/specs/rfc4918.html#copy.for.collections
		errors := map[string]int{}
		s.copyCollection(source, dest, w, r, errors)

		if err := s.CopyFile(source, dest); err != nil {
			errors[source] = StatusInternalServerError
		}

		if len(errors) != 0 {
			// send multistatus
			abs := r.RequestURI

			buf := new(bytes.Buffer)
			buf.WriteString(`<?xml version="1.0" encoding="utf-8"?>`)
			buf.WriteString(`<multistatus xmlns='DAV:'>`)

			for p, e := range errors {
				buf.WriteString(`<response>`)
				buf.WriteString(`<href>` + abs + p + `</href>`)
				buf.WriteString(`<status>HTTP/1.1 ` + string(e) + ` ` + StatusText(e) + `</status>`)
				buf.WriteString(`</response>`)
			}

			buf.WriteString(`</multistatus>`)

			w.WriteHeader(StatusMulti)
			w.Header().Set("Content-Length", string(buf.Len()))
			w.Header().Set("Content-Type", "application/xml; charset=utf-8")
			buf.WriteTo(w)

			return false
		}
	}

	// copy was successful
	if exists {
		w.WriteHeader(StatusNoContent)

	} else {
		w.WriteHeader(StatusCreated)
	}

	s.unlockResource(dest)
	return true
}

func (s *Server) CopyFile(source, dest string) error {
	fs, err := s.Fs.Open(source)
	if err != nil {
		return err
	}

	fd, err := s.Fs.Create(dest)
	if err != nil {
		return err
	}

	if _, err := io.Copy(fd, fs); err != nil {
		return err
	}

	return nil
}

func (s *Server) copyCollection(source, dest string, w http.ResponseWriter, r *http.Request, errors map[string]int) {
	ifHeader := r.Header.Get("If")
	lockToken := r.Header.Get("Lock-Token")

	for _, sub := range s.directoryContents(source) {
		ssub := source + "/" + sub
		dsub := dest + "/" + sub

		if s.isLocked(ssub, ifHeader+lockToken) {
			errors[ssub] = StatusLocked
		} else {
			if s.pathIsDirectory(ssub) {
				s.copyCollection(ssub, dsub, w, r, errors)
			}

			if err := s.CopyFile(ssub, dsub); err != nil {
				errors[ssub] = StatusInternalServerError
			}
		}
	}

}

func (s *Server) doLock(w http.ResponseWriter, r *http.Request) {
	if s.ReadOnly {
		w.WriteHeader(StatusForbidden)
		return
	}

	if s.isLockedRequest(r) {
		w.WriteHeader(StatusLocked)
		return
	}

	// TODO: lock
	w.WriteHeader(StatusNotImplemented)
}

func (s *Server) doUnlock(w http.ResponseWriter, r *http.Request) {
	if s.ReadOnly {
		w.WriteHeader(StatusForbidden)
		return
	}

	if s.isLockedRequest(r) {
		w.WriteHeader(StatusLocked)
		return
	}

	// TODO: unlock
	w.WriteHeader(StatusNotImplemented)
}

func (s *Server) doOptions(w http.ResponseWriter, r *http.Request) {
	// http://www.webdav.org/specs/rfc4918.html#dav.compliance.classes
	w.Header().Set("DAV", "1, 2")

	w.Header().Set("Allow", s.methodsAllowed(s.url2path(r.URL)))
	w.Header().Set("MS-Author-Via", "DAV")
}
