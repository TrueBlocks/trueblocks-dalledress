package types

type DataLoadedPayload struct {
	CurrentCount  int      `json:"currentCount"`
	ExpectedTotal int      `json:"expectedTotal"`
	IsLoaded      bool     `json:"isLoaded"`
	ListKind      ListKind `json:"listKind,omitempty"`
}
