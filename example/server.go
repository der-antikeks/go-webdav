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

	http.HandleFunc("/", index)
	http.Handle("/webdav/", http.StripPrefix("/webdav/",
		webdav.Handler(webdav.Dir(path))))

	log.Println("Listening on http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q\n", r.URL.Path)
}
