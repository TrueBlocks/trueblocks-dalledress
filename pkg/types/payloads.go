package types

type DataLoadedPayload struct {
	CurrentCount  int    `json:"currentCount"`
	ExpectedTotal int    `json:"expectedTotal"`
	IsLoaded      bool   `json:"isLoaded"`
	ListKind      string `json:"listKind,omitempty"`
}
