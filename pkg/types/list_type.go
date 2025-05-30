package types

import "strings"

type ListKind string

type ListKindDef struct {
	Value  ListKind `json:"value"`
	TSName string   `json:"tsname"`
}

var AllListKinds = []ListKindDef{}

func RegisterKind(listKind ListKind) {
	AllListKinds = append(AllListKinds, ListKindDef{
		listKind,
		strings.ToUpper(string(listKind)),
	})
}
