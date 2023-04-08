package logger

import (
	"io"
	"log"
	"os"
)

type MyDualWriter struct {
	W1 io.Writer
	W2 io.Writer
}

func (d *MyDualWriter) Write(p []byte) (n int, err error) {
	n, err = d.W1.Write(p)
	if err != nil {
		return n, err
	}
	return d.W2.Write(p)
}

var DualLog *log.Logger

func InitLogger(logFile io.Writer) {
	DualLog = log.New(&MyDualWriter{W1: logFile, W2: os.Stdout}, "", log.LstdFlags|log.Lshortfile)
}
