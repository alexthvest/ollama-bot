package ollama

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type SingleRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type SingleResponse struct {
	Response string `json:"response"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ChatResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

func GetCompletion(prompt string, model string) (SingleResponse, error) {
	payloadJson, err := json.Marshal(SingleRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	})

	if err != nil {
		return SingleResponse{}, err
	}

	payload := bytes.NewBuffer(payloadJson)
	req, err := http.NewRequest("POST", "http://localhost:11434/api/generate", payload)
	if err != nil {
		return SingleResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return SingleResponse{}, err
	}
	defer resp.Body.Close()

	var response SingleResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return SingleResponse{}, err
	}

	return response, nil
}

func GetChatCompletion(history []Message, model string) (ChatResponse, error) {
	payloadJson, err := json.Marshal(ChatRequest{
		Model:    model,
		Messages: history,
	})
	if err != nil {
		return ChatResponse{}, err
	}

	payload := bytes.NewBuffer(payloadJson)
	req, err := http.NewRequest("POST", "http://localhost:11434/v1/chat/completions", payload)
	if err != nil {
		return ChatResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ChatResponse{}, err
	}
	defer resp.Body.Close()

	var response ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return ChatResponse{}, err
	}

	return response, nil
}
