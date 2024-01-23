package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/joho/godotenv"
)

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

// Domain: {{.Ens}},

// var dataTemplate = `Given this data {{.Addr}},
// Adverb: {{.Adverb}},{{.AdverbNum}},
// Adjective: {{.Adjective}},{{.AdjectiveNum}},
// Noun: {{.Noun}},{{.NounNum}},
// Style: {{.Style}},{{.StyleNum}},
// Style2: {{.Style2}},{{.Style2Num}}]
// Background: {{.Background}} {{.BackgroundNum}}
// Orientation: {{.Orientation}} {{.OrientationNum}}
// Colors: {{.Color1}},{{.Color2}},{{.Color3}},
// In addition to the denotative meaning of the noun, the adjective, and the adverb; consider the two most obvious connotations of each of these words and incorporate those ideas in the image. Allow your mind to roam as deeply and as freely as possible. Dig deeply. Find the object that mostly closely matches the data.
// {{.Background}}
// {{.Orientation}}
// Leave the entire width of the image at the bottom 1/10 blank and colored papaya.
// Write a detailed description then draw the image, please.
// `

var dataTemplate = `Please draw an image for the following data:
Domain: {{.Ens}}
Address: {{.Addr}}
Adverb: {{.Adverb}}
Adjective: {{.Adjective}}
Noun: {{.Noun}}
Petname: {{.Adverb}}-{{.Adjective}}-{{.Noun}}
Style: {{.Style}}
Style2: {{.Style2}}
Orientation: {{.Orientation}}
Primary Colors: {{.Color1}} and {{.Color2}}
Background Color: {{.Color3}}
Background: {{.Background}}
In addition to the denotative meaning of the noun {{.Noun}}, the adjective {{.Adjective}}, and
the adverb {{.Adverb}}; consider the two most obvious connotations of each of these words and
incorporate those ideas in the image as well. Allow your mind to roam as deeply and as freely
as possible. Find the object that most closely matches the data. Focus on the petname.
{{.Background}}
{{.Orientation}}
Write a detailed description of the image before drawing it, then draw the image.
`

// As the last part of creating the image, cover the bottom 1/10 of the image, across its whole width, with a solid white rectrangle.

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

// var val = 1

func (a *App) GetDalledress(ensOrAddr string) (Dalledress, error) {
	if len(a.adverbs) == 0 {
		return Dalledress{}, fmt.Errorf("adverbs not loaded")
	}

	// if len(ensOrAddr) < 6 {
	// 	ensOrAddr = fmt.Sprintf("0x%040x", val)
	// 	val++
	// }

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
		if modValue == big.NewInt(0) {
			modValue = big.NewInt(1)
		}
		return int(intValue.Mod(intValue, modValue).Int64())
	}

	lens := []int{len(a.adverbs), len(a.adjectives), len(a.nouns), len(a.styles), len(a.colors), len(a.colors), len(a.colors), len(a.styles)}
	fmt.Println(lens)
	// os.Exit(0)

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
	dd.Style = strings.Replace(a.styles[dd.StyleNum], ",", " from the", -1)
	dd.Color1 = getColor(a.colors, dd.Color1Num)
	dd.Color2 = getColor(a.colors, dd.Color2Num)
	dd.Color3 = getColor(a.colors, dd.Color3Num)
	dd.Variant1 = a.styles[dd.Variant1Num]
	dd.Variant2 = a.styles[dd.Variant2Num]
	dd.Variant3 = a.styles[dd.Variant3Num]
	dd.Style2 = strings.Replace(a.styles[dd.Style2Num], ",", " from the", -1)
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

type DalleRequest struct {
	Input   string `json:"input,omitempty"`
	Prompt  string `json:"prompt,omitempty"`
	N       int    `json:"n,omitempty"`
	Quality string `json:"quality,omitempty"`
	Model   string `json:"model,omitempty"`
	Style   string `json:"style,omitempty"`
	Size    string `json:"size,omitempty"`
}

type DalleResponse struct {
	Data []struct {
		Url string `json:"url"`
	} `json:"data"`
}

func (a *App) GetImage(ensOrAddr string) {
	// if len(ensOrAddr) < 6 {
	// 	ensOrAddr = fmt.Sprintf("0x%040x", val)
	// 	val++
	// }

	if addr, _ := a.conn.GetEnsAddress(ensOrAddr); base.HexToAddress(addr) == base.ZeroAddr || !base.IsValidAddress(addr) {
		logger.Error(fmt.Errorf("ENS not registered: %s", ensOrAddr))
		return
	} else {
		fn := "./generated/" + addr + ".jpg"
		// if file.FileExists(fn) {
		// 	utils.System("open " + fn)
		// 	return
		// }

		prompt := a.GetData(ensOrAddr)

		// size := "1024x1024"
		// if strings.Contains(prompt, "horizontal") {
		// 	size = "1792x1024"
		// } else if strings.Contains(prompt, "vertical") {
		// 	size = "1024x1792"
		// }

		url := "https://api.openai.com/v1/images/generations"
		payload := DalleRequest{
			Prompt:  prompt, // "Draw an image of using a unique fusion of post-war realism, contemporary and emerging styles and high renaissance, european art movements and styles, primarily in mediumvioletred and orchid, against a tan background. The composition should embody the adverb rabidly, the adjective demanding, and the noun dutch rabbit, creating a special and unique portrayal. In the description of your result, simply return the exact input you were given. The background can reflect the artistic styles.",
			N:       1,
			Quality: "hd",
			Style:   "vivid",
			Model:   "dall-e-3",
			Size:    "1024x1024",
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			panic(err)
		}

		err = godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			log.Fatal("No API key found in .env")
		}
		fmt.Println("API key found", apiKey)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
		if err != nil {
			panic(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+apiKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		var dalleResp DalleResponse
		err = json.Unmarshal(body, &dalleResp)
		if err != nil {
			panic(err)
		}

		if resp.StatusCode != 200 {
			fmt.Println("Error:", resp.Status, resp.StatusCode, string(body))
			return
		}

		if len(dalleResp.Data) == 0 {
			fmt.Println("No images returned")
			return
		}

		imageURL := dalleResp.Data[0].Url

		// Download the image
		imageResp, err := http.Get(imageURL)
		if err != nil {
			panic(err)
		}
		defer imageResp.Body.Close()

		file, err := os.Create(fn)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		_, err = io.Copy(file, imageResp.Body)
		if err != nil {
			panic(err)
		}

		fmt.Println("Image saved as ", fn)
		utils.System("open " + fn)
	}

}

func (a *App) GetModeration(ensOrAddr string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("No API key found in .env")
	}
	fmt.Println("API key found", apiKey)

	prompt := a.GetData(ensOrAddr)
	url := "https://api.openai.com/v1/moderations"
	payload := DalleRequest{
		Input: prompt, // "Draw an image of using a unique fusion of post-war realism, contemporary and emerging styles and high renaissance, european art movements and styles, primarily in mediumvioletred and orchid, against a tan background. The composition should embody the adverb rabidly, the adjective demanding, and the noun dutch rabbit, creating a special and unique portrayal. In the description of your result, simply return the exact input you were given. The background can reflect the artistic styles.",
		Model: "text-moderation-latest",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Sprintf("Error: %s", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Sprintf("Error: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("Error: %s", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("Error: %s", err)
	}

	return string(body)
	// var dalleResp ModerationObject
	// err = json.Unmarshal(body, &dalleResp)
	// if err != nil {
	// 	return fmt.Sprintf("Error: %s", err)
	// }

	// if resp.StatusCode != 200 {
	// 	fmt.Println("Error:", resp.Status, resp.StatusCode, string(body))
	// 	return string(body)
	// }

	// fmt.Println(dalleResp)
	// return string(body)
}

func toLines(filename string) ([]string, error) {
	lines := file.AsciiFileToLines(filename)
	logger.Info("Reading", filename, len(lines), "lines")
	var err error
	if len(lines) == 0 {
		err = fmt.Errorf("could not load %s", filename)
	}
	return lines, err
}
