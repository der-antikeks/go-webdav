package main

import (
	"log"

	"github.com/der-antikeks/go-webdav"
)

func main() {
	fs, err := webdav.Dial("http://localhost:8080/webdav")
	if err != nil {
		log.Fatal(err)
	}
	defer fs.Close()

	f, err := fs.Open(".")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fi, err := f.Readdir(0)
	if err != nil {
		log.Fatal(err)
	}

	for _, i := range fi {
		name := i.Name()
		if i.IsDir() {
			name += "/"
		}

		log.Println(name)
	}
}
