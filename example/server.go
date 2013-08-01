package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/der-antikeks/go-webdav"
)

var (
	path = "./tmp"
)

func main() {
	os.Mkdir(path, os.ModeDir)

	// http.StripPrefix is not working, webdav.Server has no knowledge
	// of stripped component, but needs for COPY/MOVE methods.
	// Destination path is supplied as header and needs to be stripped.
	http.Handle("/webdav/", &webdav.Server{
		Fs:         webdav.Dir(path),
		TrimPrefix: "/webdav/",
		Listings:   true,
	})

	http.HandleFunc("/", index)

	log.Println("Listening on http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q\n", r.URL.Path)
}
