package application

import "io"

type FileStore interface {
	Get(fileID string) (io.Reader, error)
	Save(file io.Reader) (string, error)
}
