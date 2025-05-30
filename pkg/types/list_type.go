package types

type ListKind string

const (
	AbisDownloaded ListKind = "Downloaded"
	AbisKnown      ListKind = "Known"
	AbisFunctions  ListKind = "Functions"
	AbisEvents     ListKind = "Events"
)

var AllListKinds = []struct {
	Value  ListKind `json:"value"`
	TSName string   `json:"tsname"`
}{
	{AbisDownloaded, "DOWNLOADED"},
	{AbisKnown, "KNOWN"},
	{AbisFunctions, "FUNCTIONS"},
	{AbisEvents, "EVENTS"},
}
