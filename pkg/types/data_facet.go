package types

import "strings"

type DataFacet string

var AllDataFacets = []struct {
	Value  DataFacet `json:"value"`
	TSName string    `json:"tsname"`
}{}

func RegisterDataFacet(df DataFacet) {
	AllDataFacets = append(AllDataFacets, struct {
		Value  DataFacet `json:"value"`
		TSName string    `json:"tsname"`
	}{
		df,
		strings.ToUpper(string(df)),
	})
}
