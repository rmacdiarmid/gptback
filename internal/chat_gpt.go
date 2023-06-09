package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/rmacdiarmid/GPTSite/logger"
	"gopkg.in/yaml.v3"
)

func LoadAPIKey() (string, error) {
	var config struct {
		OpenAI_API_Key string `yaml:"openai_api_key"`
	}

	// Get the absolute path to the config file
	absPath, err := filepath.Abs("/Users/ryanmacdiarmid/Dropbox (Personal)/_github/GoProjects/GPTSite/config.yaml")
	if err != nil {
		logger.DualLog.Printf("Error getting absolute path to config file: %s", err)
		return "", err
	}

	// Log message for file reading start
	logger.DualLog.Println("Starting to read config file...")

	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		logger.DualLog.Printf("Error reading config file: %s", err)
		return "", err
	}

	// Log message for file reading success
	logger.DualLog.Println("Config file read successfully.")

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		logger.DualLog.Printf("Error unmarshalling config data: %s", err)
		return "", err
	}

	// Log message for unmarshalling success
	logger.DualLog.Printf("Config data unmarshalled successfully. API key: %s", config.OpenAI_API_Key)

	return config.OpenAI_API_Key, nil
}

func ChatGPTRequest(prompt string) (string, error) {
	apiKey, err := LoadAPIKey()
	apiURL := "https://api.openai.com/v1/chat/completions"

	client := &http.Client{}

	data := map[string]interface{}{
		"model":       "gpt-3.5-turbo",
		"messages":    []map[string]string{{"role": "user", "content": prompt}},
		"temperature": 0.7,
	}
	jsonData, _ := json.Marshal(data)

	req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	// Log message for request start
	logger.DualLog.Println("Starting ChatGPT request...")

	resp, err := client.Do(req)
	if err != nil {
		logger.DualLog.Printf("Error sending ChatGPT request: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	choices := result["choices"].([]interface{})
	message := choices[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)

	// Log message for request success
	logger.DualLog.Println("ChatGPT request completed successfully.")

	return message, nil
}

func GenerateArticle(prompt string) (string, string, string, error) {
	logger.DualLog.Printf("GenerateArticle: Starting GenerateArticle function...")
	defer logger.DualLog.Printf("GenerateArticle: Exiting GenerateArticle function.")

	// Call ChatGPTRequest with the prompt to generate the article text
	articleText, err := ChatGPTRequest(prompt)
	if err != nil {
		logger.DualLog.Panicf("GenerateArticle: Error generating article: %v", err)
		return "", "", "", err
	}

	// Log message for article generation success
	logger.DualLog.Printf("GenerateArticle: Article generated successfully.")

	// Generate the title and image URL here
	// For now, we'll use placeholders
	title := "Example Title"
	imageURL := "https://example.com/image.jpg"

	return title, imageURL, articleText, nil
}
