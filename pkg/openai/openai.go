package openai

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type dalleRequest struct {
	Input    string    `json:"input,omitempty"`
	Prompt   string    `json:"prompt,omitempty"`
	N        int       `json:"n,omitempty"`
	Quality  string    `json:"quality,omitempty"`
	Model    string    `json:"model,omitempty"`
	Style    string    `json:"style,omitempty"`
	Size     string    `json:"size,omitempty"`
	Seed     int       `json:"seed,omitempty"`
	Messages []message `json:"messages,omitempty"`
}

type dalleResponse struct {
	Data []struct {
		Url string `json:"url"`
	} `json:"data"`
}
