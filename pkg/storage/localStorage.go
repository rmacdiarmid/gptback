// storage/local_file_storage.go
package storage

import (
	"io"
	"os"
	"path/filepath"
)

type LocalFileStorage struct {
	BasePath string
}

func (l *LocalFileStorage) GetFile(path string) (io.ReadSeekCloser, error) {
	absPath := filepath.Join(l.BasePath, path)
	file, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	return file, nil
}
