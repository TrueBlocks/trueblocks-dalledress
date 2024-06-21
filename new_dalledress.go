package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"
)

type DalleDress struct {
	Num         int                  `json:"num"`
	Orig        string               `json:"orig"`
	Seed        string               `json:"seed"`
	Prompt      string               `json:"prompt,omitempty"`
	DataPrompt  string               `json:"dataPrompt,omitempty"`
	TersePrompt string               `json:"tersePrompt,omitempty"`
	Attribs     []Attribute          `json:"attributes"`
	AttribMap   map[string]Attribute `json:"-"`
}

func NewDalleDress(i int, address string) *DalleDress {
	reverse := func(s string) string {
		runes := []rune(s)
		n := len(runes)
		for i := 0; i < n/2; i++ {
			runes[i], runes[n-1-i] = runes[n-1-i], runes[i]
		}
		return string(runes)
	}
	return &DalleDress{
		Num:       i,
		Orig:      address,
		Seed:      address + reverse(address),
		AttribMap: make(map[string]Attribute),
	}
}

func (d *DalleDress) String() string {
	jsonData, _ := json.MarshalIndent(d, "", "  ")
	return string(jsonData)
}

func (dalleDress *DalleDress) generatePrompt(t *template.Template, f func(s string) string) (string, error) {
	var buffer bytes.Buffer
	if err := t.Execute(&buffer, dalleDress); err != nil {
		return "", err
	}
	if f == nil {
		return buffer.String(), nil
	}
	return f(buffer.String()), nil
}

var databaseNames = []string{
	"adverbs",
	"adjectives",
	"nouns",
	"emotions",
	"occupations",
	"actions",
	"artstyles",
	"artstyles",
	"litstyles",
	"colors",
	"colors",
	"colors",
	"orientations",
	"gazes",
	"backstyles",
}

var attributeNames = []string{
	"adverb",
	"adjective",
	"noun",
	"emotion",
	"occupation",
	"action",
	"artStyle1",
	"artStyle2",
	"litStyle",
	"color1",
	"color2",
	"color3",
	"orientation",
	"gaze",
	"backStyle",
}

func (dalleDress *DalleDress) Adverb(short bool) string {
	val := dalleDress.AttribMap["adverb"].Value
	parts := strings.Split(val, ",")
	if short {
		return parts[0]
	}
	return parts[0] + " (" + parts[1] + ")"
}

func (dalleDress *DalleDress) Adjective(short bool) string {
	val := dalleDress.AttribMap["adjective"].Value
	parts := strings.Split(val, ",")
	if short {
		return parts[0]
	}
	return parts[0] + " (" + parts[1] + ")"
}

func (dalleDress *DalleDress) Noun(short bool) string {
	val := dalleDress.AttribMap["noun"].Value
	parts := strings.Split(val, ",")
	if short {
		return parts[0]
	}
	return parts[0] + " (" + parts[1] + ", " + parts[2] + ")"
}

func (dalleDress *DalleDress) Emotion(short bool) string {
	val := dalleDress.AttribMap["emotion"].Value
	parts := strings.Split(val, ",")
	if short {
		return parts[0]
	}
	return parts[0] + " (" + parts[1] + ", " + parts[4] + ")"
}

func (dalleDress *DalleDress) Occupation(short bool) string {
	val := dalleDress.AttribMap["occupation"].Value
	parts := strings.Split(val, ",")
	if short {
		return parts[0]
	}
	return parts[0] + " (" + parts[1] + ")"
}

func (dalleDress *DalleDress) Action(short bool) string {
	val := dalleDress.AttribMap["action"].Value
	parts := strings.Split(val, ",")
	if short {
		return parts[0]
	}
	return parts[0] + " (" + parts[1] + ")"
}

func (dalleDress *DalleDress) ArtStyle(short bool, which int) string {
	val := dalleDress.AttribMap["artStyle"+fmt.Sprintf("%d", which)].Value
	parts := strings.Split(val, ",")
	if short {
		return parts[0]
	}
	if strings.HasPrefix(parts[2], parts[0]+" ") {
		parts[2] = strings.Replace(parts[2], (parts[0] + " "), "", 1)
	}
	return parts[0] + " (" + parts[2] + ")"
}

func (dalleDress *DalleDress) LitStyle(short bool) string {
	val := dalleDress.AttribMap["litStyle"].Value
	parts := strings.Split(val, ",")
	if short {
		return parts[0]
	}
	if strings.HasPrefix(parts[1], parts[0]+" ") {
		parts[1] = strings.Replace(parts[1], (parts[0] + " "), "", 1)
	}
	return parts[0] + " (" + parts[1] + ")"
}

func (dalleDress *DalleDress) Color(short bool, which int) string {
	val := dalleDress.AttribMap["color"+fmt.Sprintf("%d", which)].Value
	parts := strings.Split(val, ",")
	if short {
		return parts[1]
	}
	return parts[1] + " (" + parts[0] + ")"
}

func (dalleDress *DalleDress) Orientation(short bool) string {
	val := dalleDress.AttribMap["orientation"].Value
	if short {
		parts := strings.Split(val, ",")
		return parts[0]
	}
	ret := `Orient the scene [{ORI}] and make sure the [{NOUN}] is facing [{GAZE}]`
	ret = strings.ReplaceAll(ret, "[{ORI}]", strings.ReplaceAll(val, ",", " and "))
	ret = strings.ReplaceAll(ret, "[{NOUN}]", dalleDress.Noun(true))
	ret = strings.ReplaceAll(ret, "[{GAZE}]", dalleDress.Gaze(true))
	return ret
}

func (dalleDress *DalleDress) Gaze(short bool) string {
	val := dalleDress.AttribMap["gaze"].Value
	if short {
		parts := strings.Split(val, ",")
		return parts[0]
	}
	return strings.ReplaceAll(val, ",", ", ")
}

func (dalleDress *DalleDress) BackStyle(short bool) string {
	val := dalleDress.AttribMap["backStyle"].Value
	val = strings.ReplaceAll(val, "[{Color3}]", dalleDress.Color(true, 3))
	val = strings.ReplaceAll(val, "[{ArtStyle2}]", dalleDress.ArtStyle(false, 2))
	return val
}
