package types

type DataLoadedPayload struct {
	DataType      string `json:"dataType"`
	CurrentCount  int    `json:"currentCount"`
	ExpectedTotal int    `json:"expectedTotal"`
	IsFullyLoaded bool   `json:"isFullyLoaded"`
	Category      string `json:"category,omitempty"`
}
