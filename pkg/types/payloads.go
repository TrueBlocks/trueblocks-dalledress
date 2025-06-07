package types

type DataLoadedPayload struct {
	CurrentCount  int      `json:"currentCount"`
	ExpectedTotal int      `json:"expectedTotal"`
	ListKind      ListKind `json:"listKind,omitempty"`
}
