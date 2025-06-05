package types

import "strings"

type ListKind string

var AllListKinds = []struct {
	Value  ListKind `json:"value"`
	TSName string   `json:"tsname"`
}{}

func RegisterKind(listKind ListKind) {
	AllListKinds = append(AllListKinds, struct {
		Value  ListKind `json:"value"`
		TSName string   `json:"tsname"`
	}{
		listKind,
		strings.ToUpper(string(listKind)),
	})
}
