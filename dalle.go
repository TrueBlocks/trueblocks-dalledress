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
	"path/filepath"
	"strings"
	"text/template"
	"time"

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
	Emotion1        string `json:"emotion1"`
	Emotion1Seed    string `json:"-"`
	Emotion1Num     int    `json:"emotion1Num"`
	Emotion2        string `json:"emotion2"`
	Emotion2Seed    string `json:"-"`
	Emotion2Num     int    `json:"emotion2Num"`
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

var promptTemplate = `Draw an image of a {{.Noun}} using a unique fusion of {{.Style}} and {{.Style2}},
primarily in {{.Color1}} and {{.Color2}}, against a {{.Color3}} background. The composition
should embody the adverb {{.Adverb}}, the adjective {{.Adjective}}, and the noun {{.Noun}},
creating a special and unique portrayal. DO NOT put any words in the image.{{.Background}}{{.Orientation}}`

var dataTemplate = `Please draw an image for the following data:
Petname: {{.Ens}}-{{.Adverb}}-{{.Adjective}}-{{.Noun}} feeling {{.Emotion1}}
Adverb: {{.Adverb}}
Adjective: {{.Adjective}}
Emotion: {{.Emotion2}}
Noun: {{.Noun}}
Style: {{.Style}}
Orientation: {{.Orientation}}
Background: {{.Background}}
Colors: {{.Color1}} and {{.Color2}}
BGStyle: {{.Style2}}
BgColor: {{.Color3}}
Address: {{.Addr}}
In addition to the denotative meaning of the noun {{.Noun}}, the emotion {{.Emotion2}}, the adjective {{.Adjective}}, and
the adverb {{.Adverb}}; consider the two most obvious connotations of each of these words and
incorporate those ideas in the image as well. Allow your mind to roam as deeply and as freely
as possible. Find the object that most closely matches the data. Focus on the petname
"{{.Ens}}-{{.Adverb}}-{{.Adjective}}-{{.Noun}} feeling {{.Emotion1}}" and the first style "{{.Style}}".
It's very important that you DO NOT PUT ANY TEXT ON THE IMAGE.
Write a detailed description of the image before drawing it, then draw the image.
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
	if len(a.adverbs) == 0 {
		return Dalledress{}, fmt.Errorf("adverbs not loaded")
	}

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
		Emotion1Seed:    seed[24:40],
		Emotion2Seed:    seed[24:40],
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

	lens := []int{len(a.adverbs), len(a.adjectives), len(a.emotions1), len(a.nouns), len(a.styles), len(a.colors), len(a.colors), len(a.colors), len(a.styles)}

	dd.AdverbNum = hexStringToBigIntModulo(dd.AdverbSeed, lens[0])
	dd.AdjectiveNum = hexStringToBigIntModulo(dd.AdjectiveSeed, lens[1])
	dd.Emotion1Num = hexStringToBigIntModulo(dd.Emotion1Seed, lens[2])
	dd.Emotion2Num = hexStringToBigIntModulo(dd.Emotion2Seed, lens[2])
	dd.NounNum = hexStringToBigIntModulo(dd.NounSeed, lens[3])
	dd.StyleNum = hexStringToBigIntModulo(dd.StyleSeed, lens[4])
	dd.Color1Num = hexStringToBigIntModulo(dd.Color1Seed, lens[5])
	dd.Color2Num = hexStringToBigIntModulo(dd.Color2Seed, lens[6])
	dd.Color3Num = hexStringToBigIntModulo(dd.Color3Seed, lens[7])
	dd.Variant1Num = hexStringToBigIntModulo(dd.Variant1Seed, lens[8])
	dd.Variant2Num = hexStringToBigIntModulo(dd.Variant2Seed, lens[8])
	dd.Variant3Num = hexStringToBigIntModulo(dd.Variant3Seed, lens[8])
	dd.Style2Num = hexStringToBigIntModulo(dd.Style2Seed, lens[8])
	dd.BackgroundNum = hexStringToBigIntModulo(dd.BackgroundSeed, 4)
	dd.OrientationNum = hexStringToBigIntModulo(dd.OrientationSeed, 4)

	dd.Adverb = a.adverbs[dd.AdverbNum]
	dd.Adjective = a.adjectives[dd.AdjectiveNum]
	dd.Emotion1 = a.emotions1[dd.Emotion1Num]
	dd.Emotion2 = a.emotions2[dd.Emotion2Num]
	dd.Noun = a.nouns[dd.NounNum]
	dd.Style = strings.Replace(a.styles[dd.StyleNum], ",", " from the", -1)
	dd.Color1 = "#" + dd.Color1Seed[:8] //getColor(a.colors, dd.Color1Num)
	dd.Color2 = "#" + dd.Color2Seed[:8] // getColor(a.colors, dd.Color2Num)
	dd.Color3 = "#" + dd.Color3Seed[:8] // getColor(a.colors, dd.Color3Num)
	dd.Variant1 = a.styles[dd.Variant1Num]
	dd.Variant2 = a.styles[dd.Variant2Num]
	dd.Variant3 = a.styles[dd.Variant3Num]
	dd.Style2 = strings.Replace(a.styles[dd.Style2Num], ",", " from the", -1)
	switch dd.BackgroundNum {
	case 0:
		dd.Background = " The background should be transparent. "
	case 1:
		dd.Background = " The background should reflect the value of the BG Style. "
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

func (a *App) GetImage(ensOrAddr string, replace bool) {
	if addr, _ := a.conn.GetEnsAddress(ensOrAddr); base.HexToAddress(addr) == base.ZeroAddr || !base.IsValidAddress(addr) {
		logger.Error(fmt.Errorf("ENS not registered: %s", ensOrAddr))
		return
	} else {
		folder := "./generated/"
		fn := filepath.Join(folder, addr+".png")
		file.EstablishFolder(folder)
		msg := fmt.Sprintf("%s,%s,%s,image\n", utils.FormattedDate(time.Now().Unix()), addr, strings.ToLower(ensOrAddr))
		if replace {
			os.Remove(fn)
			msg = fmt.Sprintf("%s,%s,%s,replace image\n", utils.FormattedDate(time.Now().Unix()), addr, strings.ToLower(ensOrAddr))
		}
		file.AppendToAsciiFile("dalledress.csv", msg)
		if file.FileExists(fn) {
			utils.System("open " + fn)
			return
		}

		prompt := a.GetData(ensOrAddr)
		size := "1024x1024"
		if strings.Contains(prompt, "horizontal") {
			size = "1792x1024"
		} else if strings.Contains(prompt, "vertical") {
			size = "1024x1792"
		}

		url := "https://api.openai.com/v1/images/generations"
		payload := DalleRequest{
			Prompt:  prompt,
			N:       1,
			Quality: "hd",
			Style:   "vivid",
			Model:   "dall-e-3",
			Size:    size,
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			panic(err)
		}

		err = godotenv.Load()
		if err != nil {
			logger.Fatal("Error loading .env file")
		}
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			logger.Fatal("No API key found in .env")
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
		// fmt.Println("imageURL:", imageURL)

		// Download the image
		imageResp, err := http.Get(imageURL)
		if err != nil {
			panic(err)
		}
		defer imageResp.Body.Close()

		file, err := os.Create(fn)
		if err != nil {
			logger.Error("Failed to open output file: ", fn)
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
		Input: prompt,
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
