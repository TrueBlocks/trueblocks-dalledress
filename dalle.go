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
	"sync"
	"text/template"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/joho/godotenv"
	"github.com/lucasb-eyer/go-colorful"
)

type Attribute struct {
	Seed string `json:"seed"`
	Num  int    `json:"num"`
	Val  string `json:"val"`
}

type Dalledress struct {
	Ens             string    `json:"ens"`
	Addr            string    `json:"addr"`
	Seed            string    `json:"seed"`
	Adverb          Attribute `json:"adverb"`
	Adjective       string    `json:"adjective"`
	AdjectiveSeed   string    `json:"-"`
	AdjectiveNum    int       `json:"adjectiveNum"`
	Emotion1        string    `json:"emotion1"`
	Emotion1Seed    string    `json:"-"`
	Emotion1Num     int       `json:"emotion1Num"`
	Emotion2        string    `json:"emotion2"`
	Emotion2Seed    string    `json:"-"`
	Emotion2Num     int       `json:"emotion2Num"`
	Literary        string    `json:"literary"`
	LiterarySeed    string    `json:"-"`
	LiteraryNum     int       `json:"literaryNum"`
	Noun            string    `json:"noun"`
	NounSeed        string    `json:"-"`
	NounNum         int       `json:"nounNum"`
	Style           string    `json:"style"`
	StyleSeed       string    `json:"-"`
	StyleNum        int       `json:"styleNum"`
	Color1          string    `json:"color1"`
	Color1Seed      string    `json:"-"`
	Color1Num       int       `json:"color1Num"`
	Color2          string    `json:"color2"`
	Color2Seed      string    `json:"-"`
	Color2Num       int       `json:"color2Num"`
	Color3          string    `json:"color3"`
	Color3Seed      string    `json:"-"`
	Color3Num       int       `json:"color3Num"`
	Variant1        string    `json:"variant1"`
	Variant1Seed    string    `json:"-"`
	Variant1Num     int       `json:"variant1Num"`
	Variant2        string    `json:"variant2"`
	Variant2Seed    string    `json:"-"`
	Variant2Num     int       `json:"variant2Num"`
	Variant3        string    `json:"variant3"`
	Variant3Seed    string    `json:"-"`
	Variant3Num     int       `json:"variant3Num"`
	Style2          string    `json:"style2"`
	Style2Seed      string    `json:"-"`
	Style2Num       int       `json:"style2Num"`
	Background      string    `json:"background"`
	BackgroundSeed  string    `json:"-"`
	BackgroundNum   int       `json:"backgroundNum"`
	Orientation     string    `json:"orientation"`
	OrientationSeed string    `json:"-"`
	OrientationNum  int       `json:"orientationNum"`
	Prompt          string    `json:"prompt"`
	PromptError     error     `json:"promptError"`
	Data            string    `json:"data"`
	DataError       error     `json:"dataError"`
	Terse           string    `json:"terse"`
	TerseError      error     `json:"terseError"`
}

// I NEED to test how the tool works with extremely simple prompts. DO NOT add any detail,  just use it AS-IS:
// Please draw an image for the following data:

// var promptTemplate1 = `
// Petname: {{.Adverb.Val}}-{{.Adjective}}-{{.Noun}} feeling {{.Emotion1}}{{.Ens}}
// Adverb: {{.Adverb.Val}}
// Adjective: {{.Adjective}}
// Emotion: {{.Emotion2}}
// Noun: {{.Noun}}
// Style: {{.Style}}
// Use only the colors {{.Color1}} and {{.Color2}}.
// {{.Orientation}}.
// {{.Background}}.
// In addition to the denotative meaning of the noun {{.Noun}}, the emotion {{.Emotion2}}, the
// adjective {{.Adjective}}, and the adverb {{.Adverb.Val}}; consider the two most obvious connotations
// of each of these words and incorporate those ideas in the image as well. Allow your mind to roam as
// deeply into the world as possible. Find something unique. Find the object that most closely matches
// the data. Focus on the petname {{.Adverb.Val}}-{{.Adjective}}-{{.Noun}} feeling {{.Emotion1}}{{.Ens}}
// and the first style {{.Style}}. It is very important that you DO NOT PUT ANY TEXT ON THE IMAGE.
// Write a detailed description of the image in the literary style {{.Literary}}.
// `

var promptTemplate = `Petname: {{.Adverb.Val}} {{.Adjective}} {{.Noun}} feeling {{.Emotion1}}{{.Ens}}.
Noun: {{.Noun}}.
Emotion: {{.Emotion2}}.
Primary style: {{.Style}}.
Use only the colors {{.Color1}} and {{.Color2}}.
{{.Orientation}}.
{{.Background}}.
Find the two most relevant connotative meanings of the noun {{.Noun}},+
the emotion {{.Emotion2}}, the adjective {{.Adjective}}, and the adverb {{.Adverb.Val}}.
Allow your mind to roam as deeply into the language model as possible.
Find the single object that most closely matches the description.
Focus on the Petname, the Emotion, and the Primary style.
DO NOT PUT ANY TEXT ON THE IMAGE.
Write in the literary style {{.Literary}}.`

// BGStyle: {{.Style2}}

// Throw in slight hints of one or more of these additional artistic styles {{.Variant1}}, {{.Variant2}}, {{.Variant3}}.
// Address: {{.Addr}}

var dataTemplate = `
Address:     {{.Addr}} Ens: {{.Ens}}.
Seed:        {{.Seed}}.
Adverb:      {{.Adverb.Val}} Adjective: {{.Adjective}} Noun: {{.Noun}}.
Emotion1:    {{.Emotion1}} Literary: {{.Literary}}.
Style1:      {{.Style}}.
Color1:      {{.Color1}} Color2: {{.Color2}} Color3: {{.Color3}}.
Variant1:    {{.Variant1}} Variant2: {{.Variant2}} Variant3: {{.Variant3}}.
Style2:      {{.Style2}}.
Background:  {{.Background}}.
Orientation: {{.Orientation}}.`

var terseTemplate = `{{.Adverb.Val}} {{.Adjective}} {{.Noun}} feeling {{.Emotion1}} {{.Orientation}}`

func (d *Dalledress) generatePrompt(t *template.Template, f func(s string) string) (string, error) {
	var buffer bytes.Buffer
	if err := t.Execute(&buffer, d); err != nil {
		return "", err
	}
	return f(buffer.String()), nil
}

func (a *App) GetPrompt(ensOrAddr string) string {
	clean := func(s string) string {
		return strings.Replace(strings.Replace(strings.Replace(s, ". ", ".\n", -1), "|", "\t", -1), "+\n", " ", -1)
	}
	if dd, err := a.GetDalledress(ensOrAddr); err != nil {
		return fmt.Sprintf("Error generating prompt: %s\n", err)
	} else {
		if dd.Prompt, dd.PromptError = dd.generatePrompt(a.pTemplate, clean); dd.PromptError != nil {
			return dd.PromptError.Error()
		}
		return dd.Prompt
	}
}

func (a *App) GetTerse(ensOrAddr string) string {
	clean := func(s string) string {
		return strings.Replace(strings.Replace(strings.Replace(s, ". ", "\n", -1), "|", "\t", -1), "+\n", " ", -1)
	}
	if dd, err := a.GetDalledress(ensOrAddr); err != nil {
		return fmt.Sprintf("Error generating data: %s\n", err)
	} else {
		if dd.Terse, dd.TerseError = dd.generatePrompt(a.tTemplate, clean); dd.TerseError != nil {
			return dd.TerseError.Error()
		}
		return dd.Terse
	}
}

func (a *App) GetData(ensOrAddr string) string {
	clean := func(s string) string {
		return strings.Replace(strings.Replace(strings.Replace(s, ". ", ".\n", -1), "|", "\t", -1), "+\n", " ", -1)
	}
	if dd, err := a.GetDalledress(ensOrAddr); err != nil {
		return fmt.Sprintf("Error generating data: %s\n", err)
	} else {
		if dd.Data, dd.DataError = dd.generatePrompt(a.dTemplate, clean); dd.DataError != nil {
			return dd.DataError.Error()
		}
		return dd.Data
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

func hexStringToBigIntModulo(hexString string, seedBump, modulo int) int {
	intValue := new(big.Int)
	seedValue := new(big.Int)
	intValue.SetString(hexString, 16)
	seedValue.SetInt64(int64(seedBump))
	modValue := big.NewInt(int64(modulo))
	if modValue == big.NewInt(0) {
		modValue = big.NewInt(1)
	}
	intValue = intValue.Add(intValue, seedValue)
	return int(intValue.Mod(intValue, modValue).Int64())
}

func clip(s string) string {
	parts := strings.Split(s, " from the ")
	return parts[0]
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
	} else {
		ensOrAddr = strings.Replace(" named "+ensOrAddr, ".eth", "", -1)
	}

	dd := Dalledress{
		Ens:             ensOrAddr,
		Addr:            addr,
		Seed:            seed,
		Adverb:          Attribute{Seed: seed[0:12]},
		AdjectiveSeed:   seed[12:24],
		Emotion1Seed:    seed[24:36],
		Emotion2Seed:    seed[36:48],
		LiterarySeed:    seed[48:60],
		NounSeed:        seed[60:72],
		StyleSeed:       seed[72:84],
		Style2Seed:      seed[84:96],
		Color1Seed:      seed[92:104],
		Color2Seed:      seed[80:92],
		Color3Seed:      seed[68:80],
		Variant1Seed:    seed[56:68],
		Variant2Seed:    seed[44:56],
		Variant3Seed:    seed[32:44],
		BackgroundSeed:  hash[20:32],
		OrientationSeed: hash[8:20],
	}

	lengths := []int{
		len(a.adverbs),    // 0
		len(a.adjectives), // 1
		len(a.emotions1),  // 2
		len(a.literary),   // 3
		len(a.nouns),      // 4
		len(a.styles),     // 5
		8,                 // 6
	}

	dd.Adverb.Num = hexStringToBigIntModulo(dd.Adverb.Seed, SeedBump, lengths[0])
	dd.AdjectiveNum = hexStringToBigIntModulo(dd.AdjectiveSeed, SeedBump, lengths[1])
	dd.Emotion1Num = hexStringToBigIntModulo(dd.Emotion1Seed, 0, lengths[2])
	dd.Emotion2Num = hexStringToBigIntModulo(dd.Emotion2Seed, 0, lengths[2])
	dd.LiteraryNum = hexStringToBigIntModulo(dd.LiterarySeed, 0, lengths[3])
	dd.NounNum = hexStringToBigIntModulo(dd.NounSeed, 0, lengths[4])
	dd.StyleNum = hexStringToBigIntModulo(dd.StyleSeed, 0, lengths[5])
	dd.Variant1Num = hexStringToBigIntModulo(dd.Variant1Seed, 0, lengths[5])
	dd.Variant2Num = hexStringToBigIntModulo(dd.Variant2Seed, 0, lengths[5])
	dd.Variant3Num = hexStringToBigIntModulo(dd.Variant3Seed, 0, lengths[5])
	dd.Style2Num = hexStringToBigIntModulo(dd.Style2Seed, 0, lengths[5])
	dd.BackgroundNum = hexStringToBigIntModulo(dd.BackgroundSeed, 0, lengths[6])
	dd.OrientationNum = hexStringToBigIntModulo(dd.OrientationSeed, 0, lengths[6])

	dd.Adverb.Val = a.adverbs[dd.Adverb.Num]
	dd.Adjective = a.adjectives[dd.AdjectiveNum]
	dd.Emotion1 = a.emotions1[dd.Emotion1Num]
	dd.Emotion2 = a.emotions2[dd.Emotion2Num]
	dd.Literary = a.literary[dd.LiteraryNum]
	dd.Noun = a.nouns[dd.NounNum]
	dd.Style = a.styles[dd.StyleNum]
	dd.Color1 = "#" + muteColor(dd.Color1Seed[:8], dd.Color1Seed[2:2])
	dd.Color2 = "#" + muteColor(dd.Color2Seed[:8], dd.Color2Seed[3:3])
	dd.Color3 = "#" + muteColor(dd.Color3Seed[:8], dd.Color3Seed[4:4])
	dd.Variant1 = clip(a.styles[dd.Variant1Num])
	dd.Variant2 = clip(a.styles[dd.Variant2Num])
	dd.Variant3 = clip(a.styles[dd.Variant3Num])
	dd.Style2 = a.styles[dd.Style2Num]

	switch dd.BackgroundNum {
	case 0:
		dd.Background = "The background should be transparent"
	case 1:
		dd.Background = "The background should be this color {{.Color3}} and pay homage to this style {{.Style2}}"
	case 2:
		dd.Background = " The background should be this color {{.Color3}} and subtly patterned"
	case 3:
		fallthrough
	case 4:
		fallthrough
	case 5:
		fallthrough
	case 6:
		fallthrough
	case 7:
		dd.Background = "The background should be solid and colored with this color: {{.Color3}}"
	default:
		logger.Fatal("Invalid background number: ", dd.BackgroundNum)
	}
	dd.Background = strings.Replace(dd.Background, "{{.Color3}}", dd.Color3, -1)
	dd.Background = strings.Replace(dd.Background, "{{.Style2}}", dd.Style2, -1)

	ori, gaze, sym := "", "", ""
	switch dd.OrientationNum {
	case 0:
		ori, gaze, sym = "vertically", "into the camera", "symmetrically"
	case 1:
		ori, gaze, sym = "horizontally", "to the left", "asymmetrically"
	case 2:
		ori, gaze, sym = "diagonally", "into the camera", "symmetrically"
	case 3:
		ori, gaze, sym = "in a unique way", "to the right", "asymmetrically"
	case 4:
		ori, gaze, sym = "vertically", "into the camera", "symmetrically"
	case 5:
		ori, gaze, sym = "horizontally", "to the right", "asymmetrically"
	case 6:
		ori, gaze, sym = "diagonally", "into the camera", "symmetrically"
	case 7:
		ori, gaze, sym = "in a unique way", "to the left", "asymmetrically"
	default:
		logger.Fatal("Invalid orientation number: ", dd.OrientationNum)
	}
	dd.Orientation = "Orient the scene {Ori} and {Sym} and make sure the {{.Noun}} is facing {Gaze}"
	dd.Orientation = strings.Replace(dd.Orientation, "{Ori}", ori, -1)
	dd.Orientation = strings.Replace(dd.Orientation, "{Sym}", sym, -1)
	dd.Orientation = strings.Replace(dd.Orientation, "{Gaze}", gaze, -1)
	dd.Orientation = strings.Replace(dd.Orientation, "{{.Noun}}", dd.Noun, -1)
	return dd, nil
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type DalleRequest struct {
	Input    string    `json:"input,omitempty"`
	Prompt   string    `json:"prompt,omitempty"`
	N        int       `json:"n,omitempty"`
	Quality  string    `json:"quality,omitempty"`
	Model    string    `json:"model,omitempty"`
	Style    string    `json:"style,omitempty"`
	Size     string    `json:"size,omitempty"`
	Seed     int       `json:"seed,omitempty"`
	Messages []Message `json:"messages,omitempty"`
}

type DalleResponse struct {
	Data []struct {
		Url string `json:"url"`
	} `json:"data"`
}

var SeedBump = int(0)

var fM sync.Mutex
var reserved = make(map[string]bool)

func (a *App) GetImage(ensOrAddr string, replace bool) {
	if addr, _ := a.conn.GetEnsAddress(ensOrAddr); base.HexToAddress(addr) == base.ZeroAddr || !base.IsValidAddress(addr) {
		logger.Error(fmt.Errorf("ENS not registered: %s", ensOrAddr))
		return
	} else {
		folder := "./generated/"
		file.EstablishFolder(folder)
		file.EstablishFolder(strings.Replace(folder, "/generated", "/txt-generated", -1))
		cnt := 0
		fn := ""
		for {
			fn = filepath.Join(folder, fmt.Sprintf("%s-%05d.png", addr, cnt))
			fM.Lock()
			if !file.FileExists(fn) && !reserved[fn] {
				reserved[fn] = true
				fM.Unlock()
				break
			}
			fM.Unlock()
			cnt++
		}
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

		logger.Info(colors.Cyan, addr, colors.Yellow, "- improving the prompt...", colors.Off)

		prompt := a.GetImprovedPrompt(ensOrAddr)
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

		logger.Info(colors.Cyan, addr, colors.Yellow, "- generating the image...", colors.Off)

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

		// logger.Info("DalleResponse: ", string(body))
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

		txtFn := strings.Replace(strings.Replace(fn, ".png", ".txt", -1), "generated", "txt-generated", -1)
		file.StringToAsciiFile(txtFn, prompt)
		utils.System("open " + txtFn)

		os.Remove(fn)
		file, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			logger.Error("Failed to open output file: ", fn)
			panic(err)
		}
		defer file.Close()

		_, err = io.Copy(file, imageResp.Body)
		if err != nil {
			panic(err)
		}

		logger.Info(colors.Cyan, addr, colors.Green, "- image saved as", colors.White+strings.Trim(fn, " "), colors.Off)
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
	// fmt.Println("API key found", apiKey)

	prompt := a.GetPrompt(ensOrAddr)
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

func (a *App) GetImprovedPrompt(ensOrAddr string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("No API key found in .env")
	}
	// fmt.Println("API key found", apiKey)

	prompt := a.GetPrompt(ensOrAddr)
	prompt = "Please give me a detailed, but terse, rewriting of this description of an image. Be imaginative. " + prompt
	url := "https://api.openai.com/v1/chat/completions"
	payload := DalleRequest{
		Model: "gpt-3.5-turbo",
		Seed:  1337,
	}
	payload.Messages = append(payload.Messages, Message{Role: "system", Content: prompt})

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
	// logger.Info("Reading", filename, len(lines), "lines")
	var err error
	if len(lines) == 0 {
		err = fmt.Errorf("could not load %s", filename)
	}
	return lines, err
}

func muteColor(hex, test string) string {
	if test < "b" {
		return hex
	}
	hexColor := hex[0:6]
	c, _ := colorful.Hex(hexColor)
	h, s, l := c.Hsl()
	s *= 0.5
	return colorful.Hsl(h, s, l).Hex() + hex[6:8]
}
