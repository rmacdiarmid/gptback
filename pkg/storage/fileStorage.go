package storage

import (
	"io"
)

type FileStorage interface {
	GetFile(path string) (io.ReadSeekCloser, error)
	// ... other methods
}
