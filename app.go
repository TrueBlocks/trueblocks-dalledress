package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"text/template"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/rpc"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/joho/godotenv"
)

// App struct
type App struct {
	ctx        context.Context
	conn       *rpc.Connection
	addresses  []string
	adverbs    []string
	adjectives []string
	nouns      []string
	colors     []string
	styles     []string
	pTemplate  *template.Template
	dTemplate  *template.Template
	apiKey     string
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.conn = rpc.NewConnection("mainnet", false, map[string]bool{})
	a.adverbs = file.AsciiFileToLines("/Users/jrush/Desktop/Animals.1/adverbs.csv")
	a.adjectives = file.AsciiFileToLines("/Users/jrush/Desktop/Animals.1/adjectives.csv")
	a.nouns = file.AsciiFileToLines("/Users/jrush/Desktop/Animals.1/nouns.csv")
	a.colors = file.AsciiFileToLines("/Users/jrush/Desktop/Animals.1/colors.csv")
	a.styles = file.AsciiFileToLines("/Users/jrush/Desktop/Animals.1/styles.csv")
	x := make([]string, 0, len(a.styles))
	for _, s := range a.styles {
		if !strings.Contains(s, "sensitive") { // remove culturally sensitive styles
			x = append(x, s)
		}
	}
	a.styles = x
	var err error
	if a.pTemplate, err = template.New("prompt").Parse(promptTemplate); err != nil {
		logger.Fatal("could not create prompt template:", err)
	}
	if a.dTemplate, err = template.New("data").Parse(dataTemplate); err != nil {
		logger.Fatal("could not create data template:", err)
	}
	if err = godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	} else if a.apiKey = os.Getenv("OPENAI_API_KEY"); a.apiKey == "" {
		log.Fatal("No API key found in .env")
	}
}

type Dalledress struct {
	Ens             string `json:"ens"`
	Addr            string `json:"addr"`
	Seed            string `json:"seed"`
	Adverb          string `json:"adverb"`
	AdverbSeed      string `json:"-"`
	AdverbNum       int    `json:"adverbNum"`
	Adjective       string `json:"adjective"`
	AdjectiveSeed   string `json:"-"`
	AdjectiveNum    int    `json:"adjectiveNum"`
	Noun            string `json:"noun"`
	NounSeed        string `json:"-"`
	NounNum         int    `json:"nounNum"`
	Style           string `json:"style"`
	StyleSeed       string `json:"-"`
	StyleNum        int    `json:"styleNum"`
	Color1          string `json:"color1"`
	Color1Seed      string `json:"-"`
	Color1Num       int    `json:"color1Num"`
	Color2          string `json:"color2"`
	Color2Seed      string `json:"-"`
	Color2Num       int    `json:"color2Num"`
	Color3          string `json:"color3"`
	Color3Seed      string `json:"-"`
	Color3Num       int    `json:"color3Num"`
	Variant1        string `json:"variant1"`
	Variant1Seed    string `json:"-"`
	Variant1Num     int    `json:"variant1Num"`
	Variant2        string `json:"variant2"`
	Variant2Seed    string `json:"-"`
	Variant2Num     int    `json:"variant2Num"`
	Variant3        string `json:"variant3"`
	Variant3Seed    string `json:"-"`
	Variant3Num     int    `json:"variant3Num"`
	Style2          string `json:"style2"`
	Style2Seed      string `json:"-"`
	Style2Num       int    `json:"style2Num"`
	Background      string `json:"background"`
	BackgroundSeed  string `json:"-"`
	BackgroundNum   int    `json:"backgroundNum"`
	Orientation     string `json:"orientation"`
	OrientationSeed string `json:"-"`
	OrientationNum  int    `json:"orientationNum"`
}

// var promptTemplate = `For {{.Addr}} ({{.Ens}}), draw an image of a primarily {{.Color1}} and
// {{.Color2}} {{.Adverb}} {{.Adjective}} {{.Noun}} in the {{.Style}} style.
// Make the background {{.Color3}} in the {{.Style2}} style.`

// var promptTemplate = `For {{.Addr}} {{.Ens}}, draw an image of a primarily {{.Color1}} and
// {{.Color2}} {{.Adverb}} {{.Adjective}} {{.Noun}} in the {{.Style}} style.
// Color the background {{.Color3}} and use {{.Variant3}} style. Try very hard to make the image unique and special.`

var promptTemplate = `Draw an image of a {{.Noun}} using a unique fusion of {{.Style}} and {{.Style2}},
primarily in {{.Color1}} and {{.Color2}}, against a {{.Color3}} background. The composition
should embody the adverb {{.Adverb}}, the adjective {{.Adjective}}, and the noun {{.Noun}},
creating a special and unique portrayal. DO NOT put any words in the image.{{.Background}}{{.Orientation}}`

// var dataTemplate = `Given this data [{{.Addr}},{{.Ens}},{{.Seed}},
// {{.Adverb}},{{.AdverbSeed}},{{.AdverbNum}},
// {{.Adjective}},{{.AdjectiveSeed}},{{.AdjectiveNum}},
// {{.Noun}},{{.NounSeed}},{{.NounNum}},
// {{.Style}},{{.StyleSeed}},{{.StyleNum}},
// {{.Color1}},{{.Color1Seed}},{{.Color1Num}},
// {{.Color2}},{{.Color2Seed}},{{.Color2Num}},
// {{.Color3}},{{.Color3Seed}},{{.Color3Num}},
// {{.Variant1}},{{.Variant1Seed}},{{.Variant1Num}},
// {{.Variant2}},{{.Variant2Seed}},{{.Variant2Num}},
// {{.Variant3}},{{.Variant3Seed}},{{.Variant3Num}},
// {{.Style2}},{{.Style2Seed}},{{.Style2Num}}] dig as deeply as you can in your memory to find the object
// that mostly closely matches the data. Allow your mind to roam as deeply and as freely as possible.
// Find the nearest thing. Draw it.{{.Background}}{{.Orientation}}
// `

var dataTemplate = `Given this data {{.Addr}},
Domain: {{.Ens}},
Adverb: {{.Adverb}},{{.AdverbNum}},
Adjective: {{.Adjective}},{{.AdjectiveNum}},
Noun: {{.Noun}},{{.NounNum}},
Style: {{.Style}},{{.StyleNum}},
Style2: {{.Style2}},{{.Style2Num}}]
Background: {{.Background}} {{.BackgroundNum}}
Orientation: {{.Orientation}} {{.OrientationNum}}
Colors: {{.Color1}},{{.Color2}},{{.Color3}},
In addition to the denotative meaning of the noun, the adjective, and the adverb; consider the two most obvious connotations of each of these words and incorporate those ideas in the image. Allow your mind to roam as deeply and as freely as possible. Dig deeply. Find the object that mostly closely matches the data. 
{{.Background}}
{{.Orientation}}
Write a detailed description then draw the image, please.
`

func (a *App) generatePrompt(d Dalledress, t *template.Template, f func(s string) string) (string, error) {
	var buffer bytes.Buffer
	if err := t.Execute(&buffer, d); err != nil {
		return "", err
	}
	return f(buffer.String()), nil
}

func (a *App) GetPrompt(ensOrAddr string) string {
	if dd, err := a.GetDalledress(ensOrAddr); err != nil {
		return fmt.Sprintf("Error generating prompt: %s\n", err)
	} else {
		var s string
		var err error
		if s, err = a.generatePrompt(dd, a.pTemplate, func(s string) string { return strings.Replace(s, "\n", " ", -1) }); err != nil {
			s += fmt.Sprintf("Error generating prompt: %s\n", err)
		}
		return s
	}
}

func (a *App) GetData(ensOrAddr string) string {
	if dd, err := a.GetDalledress(ensOrAddr); err != nil {
		return fmt.Sprintf("%s", err)
	} else {
		var s string
		var err error
		if s, err = a.generatePrompt(dd, a.dTemplate, func(s string) string { return strings.Replace(s, ",", ", ", -1) }); err != nil {
			s += fmt.Sprintf("%s", err)
		}
		return s
	}
}

func (a *App) GetJson(ensOrAddr string) string {
	if dd, err := a.GetDalledress(ensOrAddr); err != nil {
		return fmt.Sprintf("{\"error\": \"%s\"}", err)
	} else {
		bytes, _ := json.Marshal(dd)
		return strings.Replace(string(bytes), ",", ", ", -1)
	}
}

func (a *App) GetAddress(index int) string {
	if len(a.addresses) == 0 {
		a.addresses = file.AsciiFileToLines("/Users/jrush/Desktop/Animals.1/addresses.csv")
	}
	return a.addresses[index%len(a.addresses)]
}

func (a *App) GetDalledress(ensOrAddr string) (Dalledress, error) {
	addr, _ := a.conn.GetEnsAddress(ensOrAddr)
	if base.HexToAddress(addr) == base.ZeroAddr || !base.IsValidAddress(addr) {
		return Dalledress{}, fmt.Errorf("ENS not registered: %s", ensOrAddr)
	}
	hash := hexutil.Encode(crypto.Keccak256([]byte(addr)))
	seed := hash[2:] + addr[2:]
	if len(seed) != 104 {
		return Dalledress{}, fmt.Errorf("invalid seed: %s", seed)
	}
	if ensOrAddr == addr {
		ensOrAddr = ""
	}

	dd := Dalledress{
		Ens:             strings.Replace(ensOrAddr, ".eth", "", -1),
		Addr:            addr,
		Seed:            seed,
		AdverbSeed:      seed[0:16],
		AdjectiveSeed:   seed[16:32],
		NounSeed:        seed[32:48],
		StyleSeed:       seed[48:64],
		Color1Seed:      seed[64:76],
		Color2Seed:      seed[76:88],
		Color3Seed:      seed[88:100],
		Variant1Seed:    seed[100:104],
		Variant2Seed:    seed[101:104] + seed[100:100],
		Variant3Seed:    seed[102:104] + seed[100:101],
		Style2Seed:      seed[103:104] + seed[100:102],
		BackgroundSeed:  hash[2:4],
		OrientationSeed: hash[4:8],
	}

	hexStringToBigIntModulo := func(hexString string, modulo int) int {
		intValue := new(big.Int)
		intValue.SetString(hexString, 16)
		modValue := big.NewInt(int64(modulo))
		return int(intValue.Mod(intValue, modValue).Int64())
	}

	lens := []int{len(a.adverbs), len(a.adjectives), len(a.nouns), len(a.styles), len(a.colors), len(a.colors), len(a.colors), len(a.styles)}
	dd.AdverbNum = hexStringToBigIntModulo(dd.AdverbSeed, lens[0])
	dd.AdjectiveNum = hexStringToBigIntModulo(dd.AdjectiveSeed, lens[1])
	dd.NounNum = hexStringToBigIntModulo(dd.NounSeed, lens[2])
	dd.StyleNum = hexStringToBigIntModulo(dd.StyleSeed, lens[3])
	dd.Color1Num = hexStringToBigIntModulo(dd.Color1Seed, lens[4])
	dd.Color2Num = hexStringToBigIntModulo(dd.Color2Seed, lens[5])
	dd.Color3Num = hexStringToBigIntModulo(dd.Color3Seed, lens[6])
	dd.Variant1Num = hexStringToBigIntModulo(dd.Variant1Seed, lens[7])
	dd.Variant2Num = hexStringToBigIntModulo(dd.Variant2Seed, lens[7])
	dd.Variant3Num = hexStringToBigIntModulo(dd.Variant3Seed, lens[7])
	dd.Style2Num = hexStringToBigIntModulo(dd.Style2Seed, lens[7])
	dd.BackgroundNum = hexStringToBigIntModulo(dd.BackgroundSeed, 4)
	dd.OrientationNum = hexStringToBigIntModulo(dd.OrientationSeed, 4)

	getColor := func(colors []string, index int) string {
		ret := colors[index]
		parts := strings.Split(ret, "|")
		return parts[0]
	}

	dd.Adverb = a.adverbs[dd.AdverbNum]
	dd.Adjective = a.adjectives[dd.AdjectiveNum]
	dd.Noun = a.nouns[dd.NounNum]
	dd.Style = a.styles[dd.StyleNum]
	dd.Color1 = getColor(a.colors, dd.Color1Num)
	dd.Color2 = getColor(a.colors, dd.Color2Num)
	dd.Color3 = getColor(a.colors, dd.Color3Num)
	dd.Variant1 = a.styles[dd.Variant1Num]
	dd.Variant2 = a.styles[dd.Variant2Num]
	dd.Variant3 = a.styles[dd.Variant3Num]
	dd.Style2 = a.styles[dd.Style2Num]
	switch dd.BackgroundNum {
	case 0:
		dd.Background = " The background should be transparent. "
	case 1:
		dd.Background = " The background should reflect the artistic style. "
	case 2:
		dd.Background = " The background should subtly stripped or checked. "
	default:
		dd.Background = " The background should be a solid color. "
	}
	switch dd.OrientationNum {
	case 0:
		dd.Orientation = " The image should be oriented vertically. "
	case 1:
		dd.Orientation = " The image should be oriented horizontally. "
	case 2:
		dd.Orientation = " The image should be oriented diagonally. "
	default:
		dd.Orientation = " The image should be oriented in a unique way. "
	}

	return dd, nil
}
