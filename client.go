package webdav

import ()

type FileSystemCloser interface {
	FileSystem
	Close() error
}

func Dial(url string) (FileSystemCloser, error) {
	// TODO
	return nil, ErrNotImplemented
}
