package openai

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
)

func EnhancePrompt(prompt, authorType string) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"
	payload := dalleRequest{
		Model:     "gpt-4",
		Seed:      1337,
		Tempature: 0.2,
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

	var response dalleResponse
	json.Unmarshal(body, &response)
	ret := strings.ReplaceAll(response.Choices[0].Message.Content, "\n", " ")
	ret = strings.ReplaceAll(ret, "  ", " ")
	return ret, nil
}

type Message struct {
	Content string `json:"content"`
}
type dalleChoice struct {
	Message Message `json:"message"`
}

type dalleResponse struct {
	Choices []dalleChoice `json:"choices"`
}
