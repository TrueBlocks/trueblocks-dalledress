package types

type DataLoadedPayload struct {
	DataType      string `json:"dataType"`
	CurrentCount  int    `json:"currentCount"`
	ExpectedTotal int    `json:"expectedTotal"`
	IsLoaded      bool   `json:"isLoaded"`
}
