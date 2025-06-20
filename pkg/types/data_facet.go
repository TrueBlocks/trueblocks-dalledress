package types

import "strings"

type DataFacet string

var AllDataFacets = []struct {
	Value  DataFacet `json:"value"`
	TSName string    `json:"tsname"`
}{}

func RegisterDataFacet(dataFacet DataFacet) {
	AllDataFacets = append(AllDataFacets, struct {
		Value  DataFacet `json:"value"`
		TSName string    `json:"tsname"`
	}{
		dataFacet,
		strings.ToUpper(string(dataFacet)),
	})
}
