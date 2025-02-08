package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

type AIHandler struct{}

func NewAIHandler() *AIHandler {
	return &AIHandler{}
}

func (h *AIHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.HandlePrompt)
	return r
}

type AIRequest struct {
	Prompt string `json:"prompt"`
	Model  string `json:"model"`
}

type AIResponse struct {
	Response string `json:"response"`
}

type OpenRouterRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenRouterResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (h *AIHandler) HandlePrompt(w http.ResponseWriter, r *http.Request) {
	var req AIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		http.Error(w, "API key not configured", http.StatusInternalServerError)
		return
	}

	if req.Model == "" {
		req.Model = "meta-llama/llama-3-8b-instruct:free"
	}

	openRouterReq := OpenRouterRequest{
		Model: req.Model,
		Messages: []ChatMessage{
			{
				Role:    "user",
				Content: req.Prompt,
			},
		},
	}

	jsonData, err := json.Marshal(openRouterReq)
	if err != nil {
		http.Error(w, "Failed to create request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	request, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, "Failed to create request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+apiKey)
	request.Header.Set("HTTP-Referer", os.Getenv("APP_URL"))
	request.Header.Set("X-Title", os.Getenv("APP_NAME"))

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		http.Error(w, "Failed to send request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "OpenRouter API error: "+string(body), resp.StatusCode)
		return
	}

	var openRouterResp OpenRouterResponse
	if err := json.Unmarshal(body, &openRouterResp); err != nil {
		http.Error(w, "Failed to parse response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(openRouterResp.Choices) == 0 {
		http.Error(w, "No response from OpenRouter", http.StatusInternalServerError)
		return
	}

	response := AIResponse{
		Response: openRouterResp.Choices[0].Message.Content,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
