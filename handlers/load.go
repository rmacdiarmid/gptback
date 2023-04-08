package handlers

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/rmacdiarmid/GPTSite/logger"
)

func LoadEnvFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Error opening .env file: %s", err)
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
			log.Printf("Error setting environment variable %s: %s", key, err)
			return err
		}
		logger.DualLog.Printf("Environment variable %s set successfully.", key)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading .env file: %s", err)
		return err
	}
	logger.DualLog.Println(".env file read successfully.")

	return nil
}
