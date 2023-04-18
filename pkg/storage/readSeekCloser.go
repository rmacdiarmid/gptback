// storage/readseekcloser.go
package storage

import (
	"bytes"
	"io"
)

type ReadSeekCloser struct {
	*bytes.Reader
}

func (r *ReadSeekCloser) Close() error {
	return nil
}

func (r *ReadSeekCloser) Seek(offset int64, whence int) (int64, error) {
	return r.Reader.Seek(offset, whence)
}

func NewReadSeekCloser(b []byte) io.ReadSeekCloser {
	return &ReadSeekCloser{
		Reader: bytes.NewReader(b),
	}
}

func toReadSeekCloser(reader io.Reader) io.ReadSeekCloser {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	return NewReadSeekCloser(buf.Bytes())
}
