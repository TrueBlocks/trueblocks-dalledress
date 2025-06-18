package types

type Summary struct {
	TotalCount  int                    `json:"totalCount"`
	FacetCounts map[ListKind]int       `json:"facetCounts"`
	CustomData  map[string]interface{} `json:"customData,omitempty"`
	LastUpdated int64                  `json:"lastUpdated"`
}

type SummaryAccumulator interface {
	AccumulateItem(item interface{}, summary *Summary)
	GetCurrentSummary() Summary
	ResetSummary()
}

type DataLoadedPayload struct {
	Collection    string    `json:"collection"`
	ListKind      ListKind  `json:"listKind"`
	CurrentCount  int       `json:"currentCount"`
	ExpectedTotal int       `json:"expectedTotal"`
	State         LoadState `json:"state"`
	Summary       Summary   `json:"summary"`
	Error         string    `json:"error,omitempty"`
	Timestamp     int64     `json:"timestamp"`
	EventPhase    string    `json:"eventPhase"`
}
