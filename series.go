package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
)

type Series struct {
	Last         int      `json:"last,omitempty"`
	Suffix       string   `json:"suffix"`
	Adverbs      []string `json:"adverbs"`
	Adjectives   []string `json:"adjectives"`
	Nouns        []string `json:"nouns"`
	Emotions     []string `json:"emotions"`
	Occupations  []string `json:"occupations"`
	Actions      []string `json:"actions"`
	Artstyles    []string `json:"artstyles"`
	Litstyles    []string `json:"litstyles"`
	Colors       []string `json:"colors"`
	Orientations []string `json:"orientations"`
	Gazes        []string `json:"gazes"`
	Backstyles   []string `json:"backstyles"`
}

func (s *Series) String() string {
	bytes, _ := json.MarshalIndent(s, "", "  ")
	return string(bytes)
}

func (s *Series) Save(fn string) {
	ss := *s
	ss.Last = 0
	file.EstablishFolder(fn)
	file.StringToAsciiFile(fn, ss.String())
}

func GetSeries() (Series, error) {
	str := strings.TrimSpace(file.AsciiFileToString("series.json"))
	if len(str) == 0 {
		return Series{}, nil
	}

	bytes := []byte(str)
	var s Series
	if err := json.Unmarshal(bytes, &s); err != nil {
		logger.Error("could not unmarshal series:", err)
		return Series{}, err
	}
	s.Suffix = strings.ReplaceAll(s.Suffix, " ", "-")
	file.EstablishFolder("./series/")
	file.StringToAsciiFile(filepath.Join("./series", s.Suffix+".json"), s.String())
	return s, nil
}

func (s *Series) GetFilter(fieldName string) ([]string, error) {
	reflectedT := reflect.ValueOf(s)
	field := reflect.Indirect(reflectedT).FieldByName(fieldName)
	if !field.IsValid() {
		return nil, fmt.Errorf("field %s not valid", fieldName)
	}
	if field.Kind() != reflect.Slice {
		return nil, fmt.Errorf("field %s not a slice", fieldName)
	}
	if field.Type().Elem().Kind() != reflect.String {
		return nil, fmt.Errorf("field %s not a string slice", fieldName)
	}
	return field.Interface().([]string), nil
}
