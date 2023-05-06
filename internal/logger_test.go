package internal

import (
	"bytes"
	"io"
	"os"

	"github.com/rmacdiarmid/gptback/logger"
)

func initTestLogger() {
	// Create a buffer to store logs in the test environment
	var testLogBuf bytes.Buffer

	// Create a multiwriter to write logs to both stdout and the buffer
	multiWriter := io.MultiWriter(os.Stdout, &testLogBuf)

	// Initialize the test logger
	logger.InitLogger(multiWriter)
}
