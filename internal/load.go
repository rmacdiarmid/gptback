package internal

import (
	"bufio"
	"os"
	"strings"

	"github.com/rmacdiarmid/gptback/logger"
)

func LoadEnvFile(path string) error {
	logger.DualLog.Printf("Starting LoadEnvFile function with path %s...", path)

	file, err := os.Open(path)
	if err != nil {
		logger.DualLog.Fatalf("Error opening .env file: %s", err)
		return err
	}
	defer file.Close()
	logger.DualLog.Println(".env file opened successfully.")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			continue
		}
		key := parts[0]
		value := parts[1]
		err := os.Setenv(key, value)
		if err != nil {
			logger.DualLog.Printf("Error setting environment variable %s: %s", key, err)
			return err
		}
		logger.DualLog.Printf("Environment variable %s set successfully.", key)
	}

	if err := scanner.Err(); err != nil {
		logger.DualLog.Fatalf("Error reading .env file: %s", err)
		return err
	}
	logger.DualLog.Println(".env file read successfully.")

	logger.DualLog.Println("Exiting LoadEnvFile function.")
	return nil
}
