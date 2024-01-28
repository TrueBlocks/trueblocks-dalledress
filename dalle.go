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

type Value struct {
	Val   string `json:"val"`
	Error error  `json:"error"`
}

type Dalledress struct {
	Ens          string    `json:"ens"`
	Addr         string    `json:"addr"`
	Seed         string    `json:"seed"`
	Adverb       Attribute `json:"adverb"`
	Adjective    Attribute `json:"adjective"`
	EmotionShort Attribute `json:"emotionShort"`
	Emotion      Attribute `json:"emotion"`
	Literary     Attribute `json:"literary"`
	Noun         Attribute `json:"noun"`
	Style        Attribute `json:"style"`
	Style2       Attribute `json:"style2"`
	Color1       Attribute `json:"color1"`
	Color2       Attribute `json:"color2"`
	Color3       Attribute `json:"color3"`
	Variant1     Attribute `json:"variant1"`
	Variant2     Attribute `json:"variant2"`
	Variant3     Attribute `json:"variant3"`
	Background   Attribute `json:"background"`
	Orientation  Attribute `json:"orientation"`
	Prompt       Value     `json:"prompt"`
	Data         Value     `json:"data"`
	Terse        Value     `json:"terse"`
}

var promptTemplate = `Petname: {{.Adverb.Val}} {{.Adjective.Val}} {{.Noun.Val}} feeling {{.EmotionShort.Val}}{{.Ens}}.
Noun: {{.Noun.Val}}.
Emotion: {{.Emotion.Val}}.
Primary style: {{.Style.Val}}.
Use only the colors {{.Color1.Val}} and {{.Color2.Val}}.
{{.Orientation.Val}}.
{{.Background.Val}}.
Find the two most relevant connotative meanings of the noun {{.Noun.Val}},+
the emotion {{.Emotion.Val}}, the adjective {{.Adjective.Val}}, and the adverb {{.Adverb.Val}}.
Allow your mind to roam as deeply into the language model as possible.
Find the single object that most closely matches the description.
Focus on the Petname, the Emotion, and the Primary style.
DO NOT PUT ANY TEXT ON THE IMAGE.
Write in the literary style {{.Literary.Val}}.`

// BGStyle: {{.Style2.Val}}

// Throw in slight hints of one or more of these additional artistic styles {{.Variant1.Val}}, {{.Variant2.Val}}, {{.Variant3.Val}}.
// Address: {{.Addr}}

var dataTemplate = `
Address:     {{.Addr}} Ens: {{.Ens}}.
Seed:        {{.Seed}}.
Adverb:      {{.Adverb.Val}} Adjective: {{.Adjective.Val}} Noun: {{.Noun.Val}}.
Emotion:     {{.EmotionShort.Val}} Literary: {{.Literary.Val}}.
Style1:      {{.Style.Val}}.
Color1:      {{.Color1.Val}} Color2: {{.Color2.Val}} Color3: {{.Color3.Val}}.
Variant1:    {{.Variant1.Val}} Variant2: {{.Variant2.Val}} Variant3: {{.Variant3.Val}}.
Style2:      {{.Style2.Val}}.
Background:  {{.Background.Val}}.
Orientation: {{.Orientation.Val}}.`

var terseTemplate = `{{.Adverb.Val}} {{.Adjective.Val}} {{.Noun.Val}} feeling {{.EmotionShort.Val}} {{.Orientation.Val}}`

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
		if dd.Prompt.Val, dd.Prompt.Error = dd.generatePrompt(a.pTemplate, clean); dd.Prompt.Error != nil {
			return dd.Prompt.Error.Error()
		}
		return dd.Prompt.Val
	}
}

func (a *App) GetTerse(ensOrAddr string) string {
	clean := func(s string) string {
		return strings.Replace(strings.Replace(strings.Replace(s, ". ", "\n", -1), "|", "\t", -1), "+\n", " ", -1)
	}
	if dd, err := a.GetDalledress(ensOrAddr); err != nil {
		return fmt.Sprintf("Error generating data: %s\n", err)
	} else {
		if dd.Terse.Val, dd.Terse.Error = dd.generatePrompt(a.tTemplate, clean); dd.Terse.Error != nil {
			return dd.Terse.Error.Error()
		}
		return dd.Terse.Val
	}
}

func (a *App) GetData(ensOrAddr string) string {
	clean := func(s string) string {
		return strings.Replace(strings.Replace(strings.Replace(s, ". ", ".\n", -1), "|", "\t", -1), "+\n", " ", -1)
	}
	if dd, err := a.GetDalledress(ensOrAddr); err != nil {
		return fmt.Sprintf("Error generating data: %s\n", err)
	} else {
		if dd.Data.Val, dd.Data.Error = dd.generatePrompt(a.dTemplate, clean); dd.Data.Error != nil {
			return dd.Data.Error.Error()
		}
		return dd.Data.Val
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
		Ens:          ensOrAddr,
		Addr:         addr,
		Seed:         seed,
		Adverb:       Attribute{Seed: seed[0:12]},
		Adjective:    Attribute{Seed: seed[12:24]},
		EmotionShort: Attribute{Seed: seed[24:36]},
		Emotion:      Attribute{Seed: seed[36:48]},
		Literary:     Attribute{Seed: seed[48:60]},
		Noun:         Attribute{Seed: seed[60:72]},
		Style:        Attribute{Seed: seed[72:84]},
		Style2:       Attribute{Seed: seed[84:96]},
		Color1:       Attribute{Seed: seed[92:104]},
		Color2:       Attribute{Seed: seed[80:92]},
		Color3:       Attribute{Seed: seed[68:80]},
		Variant1:     Attribute{Seed: seed[56:68]},
		Variant2:     Attribute{Seed: seed[44:56]},
		Variant3:     Attribute{Seed: seed[32:44]},
		Background:   Attribute{Seed: hash[20:32]},
		Orientation:  Attribute{Seed: hash[8:20]},
	}

	lengths := []int{
		len(a.adverbs),       // 0
		len(a.adjectives),    // 1
		len(a.emotionsShort), // 2
		len(a.literary),      // 3
		len(a.nouns),         // 4
		len(a.styles),        // 5
		8,                    // 6
	}

	dd.Adverb.Num = hexStringToBigIntModulo(dd.Adverb.Seed, SeedBump, lengths[0])
	dd.Adjective.Num = hexStringToBigIntModulo(dd.Adjective.Seed, SeedBump, lengths[1])
	dd.EmotionShort.Num = hexStringToBigIntModulo(dd.EmotionShort.Seed, 0, lengths[2])
	dd.Emotion.Num = hexStringToBigIntModulo(dd.Emotion.Seed, 0, lengths[2])
	dd.Literary.Num = hexStringToBigIntModulo(dd.Literary.Seed, 0, lengths[3])
	dd.Noun.Num = hexStringToBigIntModulo(dd.Noun.Seed, 0, lengths[4])
	dd.Style.Num = hexStringToBigIntModulo(dd.Style.Seed, 0, lengths[5])
	dd.Style2.Num = hexStringToBigIntModulo(dd.Style2.Seed, 0, lengths[5])
	dd.Variant1.Num = hexStringToBigIntModulo(dd.Variant1.Seed, 0, lengths[5])
	dd.Variant2.Num = hexStringToBigIntModulo(dd.Variant2.Seed, 0, lengths[5])
	dd.Variant3.Num = hexStringToBigIntModulo(dd.Variant3.Seed, 0, lengths[5])
	dd.Background.Num = hexStringToBigIntModulo(dd.Background.Seed, 0, lengths[6])
	dd.Orientation.Num = hexStringToBigIntModulo(dd.Orientation.Seed, 0, lengths[6])

	dd.Adverb.Val = a.adverbs[dd.Adverb.Num]
	dd.Adjective.Val = a.adjectives[dd.Adjective.Num]
	dd.EmotionShort.Val = a.emotionsShort[dd.EmotionShort.Num]
	dd.Emotion.Val = a.emotions[dd.Emotion.Num]
	dd.Literary.Val = a.literary[dd.Literary.Num]
	dd.Noun.Val = a.nouns[dd.Noun.Num]
	dd.Style.Val = a.styles[dd.Style.Num]
	dd.Style2.Val = a.styles[dd.Style2.Num]
	dd.Color1.Val = "#" + muteColor(dd.Color1.Seed[:8], dd.Color1.Seed[2:2])
	dd.Color2.Val = "#" + muteColor(dd.Color2.Seed[:8], dd.Color2.Seed[3:3])
	dd.Color3.Val = "#" + muteColor(dd.Color3.Seed[:8], dd.Color3.Seed[4:4])
	dd.Variant1.Val = clip(a.styles[dd.Variant1.Num])
	dd.Variant2.Val = clip(a.styles[dd.Variant2.Num])
	dd.Variant3.Val = clip(a.styles[dd.Variant3.Num])

	switch dd.Background.Num {
	case 0:
		dd.Background.Val = "The background should be transparent"
	case 1:
		dd.Background.Val = "The background should be this color {{.Color3.Val}} and pay homage to this style {{.Style2.Val}}"
	case 2:
		dd.Background.Val = " The background should be this color {{.Color3.Val}} and subtly patterned"
	case 3:
		fallthrough
	case 4:
		fallthrough
	case 5:
		fallthrough
	case 6:
		fallthrough
	case 7:
		dd.Background.Val = "The background should be solid and colored with this color: {{.Color3.Val}}"
	default:
		logger.Fatal("Invalid background number: ", dd.Background.Num)
	}
	dd.Background.Val = strings.Replace(dd.Background.Val, "{{.Color3.Val}}", dd.Color3.Val, -1)
	dd.Background.Val = strings.Replace(dd.Background.Val, "{{.Style2.Val}}", dd.Style2.Val, -1)

	ori, gaze, sym := "", "", ""
	switch dd.Orientation.Num {
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
		logger.Fatal("Invalid orientation number: ", dd.Orientation.Num)
	}
	dd.Orientation.Val = "Orient the scene {Ori} and {Sym} and make sure the {{.Noun.Val}} is facing {Gaze}"
	dd.Orientation.Val = strings.Replace(dd.Orientation.Val, "{Ori}", ori, -1)
	dd.Orientation.Val = strings.Replace(dd.Orientation.Val, "{Sym}", sym, -1)
	dd.Orientation.Val = strings.Replace(dd.Orientation.Val, "{Gaze}", gaze, -1)
	dd.Orientation.Val = strings.Replace(dd.Orientation.Val, "{{.Noun.Val}}", dd.Noun.Val, -1)
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
	prompt = "Please give me a detailed, but terse, rewriting of this description of an image. Be imaginative. " + prompt + " Mute the colors."
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
