package dalle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/TrueBlocks/trueblocks-core/sdk"
)

type DalleDress struct {
	Orig           string               `json:"orig"`
	Seed           string               `json:"seed"`
	Error          string               `json:"error,omitempty"`
	Prompt         string               `json:"prompt,omitempty"`
	DataPrompt     string               `json:"dataPrompt,omitempty"`
	TersePrompt    string               `json:"tersePrompt,omitempty"`
	EnhancedPrompt string               `json:"enhancedPrompt,omitempty"`
	Attribs        []Attribute          `json:"attributes"`
	AttribMap      map[string]Attribute `json:"-"`
}

func NewDalleDress(databases map[string][]string, address string) (*DalleDress, error) {
	reverse := func(s string) string {
		runes := []rune(s)
		n := len(runes)
		for i := 0; i < n/2; i++ {
			runes[i], runes[n-1-i] = runes[n-1-i], runes[i]
		}
		return string(runes)
	}
	if strings.HasSuffix(address, ".eth") {
		opts := sdk.NamesOptions{
			Terms: []string{address},
		}
		if names, _, err := opts.Names(); err != nil {
			return nil, fmt.Errorf("error getting names for %s", address)
		} else {
			if len(names) > 0 {
				address = names[0].Address.Hex()
			}
		}
	}
	parts := strings.Split(address, ",")
	seed := parts[0] + reverse(parts[0])
	if len(seed) < 66 {
		return nil, fmt.Errorf("seed length is less than 66")
	}
	seed = seed[2:66]

	dd := DalleDress{
		Orig:      address,
		Seed:      seed,
		AttribMap: make(map[string]Attribute),
	}

	for i := 0; i < len(dd.Seed); i = i + 8 {
		index := len(dd.Attribs)
		attr := NewAttribute(databases, index, dd.Seed[i:i+6])
		dd.Attribs = append(dd.Attribs, attr)
		dd.AttribMap[attr.Name] = attr
		if i+4+6 < len(dd.Seed) {
			index = len(dd.Attribs)
			attr = NewAttribute(databases, index, dd.Seed[i+4:i+4+6])
			dd.Attribs = append(dd.Attribs, attr)
			dd.AttribMap[attr.Name] = attr
		}
	}

	return &dd, nil
}

func (d *DalleDress) String() string {
	jsonData, _ := json.MarshalIndent(d, "", "  ")
	return string(jsonData)
}

func (dd *DalleDress) GeneratePrompt(t *template.Template, f func(s string) string) (string, error) {
	var buffer bytes.Buffer
	if err := t.Execute(&buffer, dd); err != nil {
		return "", err
	}
	if f == nil {
		return buffer.String(), nil
	}
	return f(buffer.String()), nil
}

var DatabaseNames = []string{
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

func (dd *DalleDress) Adverb(short bool) string {
	val := dd.AttribMap["adverb"].Value
	parts := strings.Split(val, ",")
	if short {
		return parts[0]
	}
	return parts[0] + " (" + parts[1] + ")"
}

func (dd *DalleDress) Adjective(short bool) string {
	val := dd.AttribMap["adjective"].Value
	parts := strings.Split(val, ",")
	if short {
		return parts[0]
	}
	return parts[0] + " (" + parts[1] + ")"
}

func (dd *DalleDress) Noun(short bool) string {
	val := dd.AttribMap["noun"].Value
	parts := strings.Split(val, ",")
	if short {
		return parts[0]
	}
	return parts[0] + " (" + parts[1] + ", " + parts[2] + ")"
}

func (dd *DalleDress) Emotion(short bool) string {
	val := dd.AttribMap["emotion"].Value
	parts := strings.Split(val, ",")
	if short {
		return parts[0]
	}
	return parts[0] + " (" + parts[1] + ", " + parts[4] + ")"
}

func (dd *DalleDress) Occupation(short bool) string {
	val := dd.AttribMap["occupation"].Value
	if val == "none" {
		return ""
	}
	parts := strings.Split(val, ",")
	if short {
		return parts[0]
	}
	return " who works as a " + parts[0] + " (" + parts[1] + ")"
}

func (dd *DalleDress) Action(short bool) string {
	val := dd.AttribMap["action"].Value
	parts := strings.Split(val, ",")
	if short {
		return parts[0]
	}
	return parts[0] + " (" + parts[1] + ")"
}

func (dd *DalleDress) ArtStyle(short bool, which int) string {
	val := dd.AttribMap["artStyle"+fmt.Sprintf("%d", which)].Value
	parts := strings.Split(val, ",")
	if short {
		return parts[0]
	}
	if strings.HasPrefix(parts[2], parts[0]+" ") {
		parts[2] = strings.Replace(parts[2], (parts[0] + " "), "", 1)
	}
	return parts[0] + " (" + parts[2] + ")"
}

func (dd *DalleDress) HasLitStyle() bool {
	ret := dd.AttribMap["litStyle"].Value
	return ret != "none" && ret != ""
}

func (dd *DalleDress) LitStyle(short bool) string {
	val := dd.AttribMap["litStyle"].Value
	if val == "none" {
		return ""
	}
	parts := strings.Split(val, ",")
	if short {
		return parts[0]
	}
	if strings.HasPrefix(parts[1], parts[0]+" ") {
		parts[1] = strings.Replace(parts[1], (parts[0] + " "), "", 1)
	}
	return parts[0] + " (" + parts[1] + ")"
}

func (dd *DalleDress) LitStyleDescr() string {
	val := dd.AttribMap["litStyle"].Value
	if val == "none" {
		return ""
	}
	parts := strings.Split(val, ",")
	if strings.HasPrefix(parts[1], parts[0]+" ") {
		parts[1] = strings.Replace(parts[1], (parts[0] + " "), "", 1)
	}
	return parts[1]
}

func (dd *DalleDress) Color(short bool, which int) string {
	val := dd.AttribMap["color"+fmt.Sprintf("%d", which)].Value
	parts := strings.Split(val, ",")
	if short {
		return parts[1]
	}
	return parts[1] + " (" + parts[0] + ")"
}

func (dd *DalleDress) Orientation(short bool) string {
	val := dd.AttribMap["orientation"].Value
	if short {
		parts := strings.Split(val, ",")
		return parts[0]
	}
	ret := `Orient the scene [{ORI}] and make sure the [{NOUN}] is facing [{GAZE}]`
	ret = strings.ReplaceAll(ret, "[{ORI}]", strings.ReplaceAll(val, ",", " and "))
	ret = strings.ReplaceAll(ret, "[{NOUN}]", dd.Noun(true))
	ret = strings.ReplaceAll(ret, "[{GAZE}]", dd.Gaze(true))
	return ret
}

func (dd *DalleDress) Gaze(short bool) string {
	val := dd.AttribMap["gaze"].Value
	if short {
		parts := strings.Split(val, ",")
		return parts[0]
	}
	return strings.ReplaceAll(val, ",", ", ")
}

func (dd *DalleDress) BackStyle(short bool) string {
	val := dd.AttribMap["backStyle"].Value
	val = strings.ReplaceAll(val, "[{Color3}]", dd.Color(true, 3))
	val = strings.ReplaceAll(val, "[{ArtStyle2}]", dd.ArtStyle(false, 2))
	return val
}

func (dd *DalleDress) LitPrompt(short bool) string {
	val := dd.AttribMap["litStyle"].Value
	if val == "none" {
		return ""
	}
	text := `Please give me a detailed rewrite of the following
	prompt in the literary style ` + dd.LitStyle(short) + `. 
	Be imaginative, creative, and complete.
`
	return text
}
