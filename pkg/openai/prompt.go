package openai

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func EnhancePrompt(prompt string) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"
	payload := dalleRequest{
		Model: "gpt-3.5-turbo",
		Seed:  1337,
	}
	payload.Messages = append(payload.Messages, message{Role: "system", Content: prompt})

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	apiKey := os.Getenv("OPENAI_API_KEY")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
