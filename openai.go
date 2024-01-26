package main

type Category struct {
	Harassment            bool `json:"harassment"`
	HarassmentThreatening bool `json:"harassment/threatening"`
	Hate                  bool `json:"hate"`
	HateThreatening       bool `json:"hate/threatening"`
	Selfharm              bool `json:"self-harm"`
	SelfharmInstructions  bool `json:"self-harm/instructions"`
	SelfharmIntent        bool `json:"self-harm/intent"`
	Sexual                bool `json:"sexual"`
	SexualMinors          bool `json:"sexual/minors"`
	Violence              bool `json:"violence"`
	ViolenceGraphic       bool `json:"violence/graphic"`
}

type Score struct {
	Harassment            float64 `json:"harassment"`
	HarassmentThreatening float64 `json:"harassment/threatening"`
	Hate                  float64 `json:"hate"`
	HateThreatening       float64 `json:"hate/threatening"`
	Selfharm              float64 `json:"self-harm"`
	SelfharmInstructions  float64 `json:"self-harm/instructions"`
	SelfharmIntent        float64 `json:"self-harm/intent"`
	Sexual                float64 `json:"sexual"`
	SexualMinors          float64 `json:"sexual/minors"`
	Violence              float64 `json:"violence"`
	ViolenceGraphic       float64 `json:"violence/graphic"`
}

type Results struct {
	Flagged    bool       `json:"flagged"`
	Categories []Category `json:"categories"`
	Scores     []Score    `json:"scores"`
}

type ModerationObject struct {
	ID      string    `json:"id"`
	Model   string    `json:"model"`
	Results []Results `json:"results"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string        `json:"model"`
	Seed     int           `json:"seed"`
	Messages []ChatMessage `json:"messages"`
}
