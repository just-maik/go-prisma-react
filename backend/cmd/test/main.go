package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type AIRequest struct {
	Prompt string `json:"prompt"`
	Model  string `json:"model"`
}

type AIResponse struct {
	Response string `json:"response"`
}

func main() {
	// Create the request payload
	reqBody := AIRequest{
		Prompt: "Translate 'Prompt Testing' into these languages: en,fr,de and respond ONLY with a json object where the key is the language and the value is the translation",
		Model:  "google/gemini-2.0-flash-lite-preview-02-05:free",
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Printf("Failed to marshal request: %v\n", err)
		os.Exit(1)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", "http://localhost:8080/api/ai", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to send request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response: %v\n", err)
		os.Exit(1)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error response from server (status %d): %s\n", resp.StatusCode, string(body))
		os.Exit(1)
	}

	// Parse the response
	var aiResp AIResponse
	if err := json.Unmarshal(body, &aiResp); err != nil {
		fmt.Printf("Failed to parse response: %v\n", err)
		os.Exit(1)
	}

	// Print the result
	fmt.Println("AI Response:")
	fmt.Println(aiResp.Response)
}
