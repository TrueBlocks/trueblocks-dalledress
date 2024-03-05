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

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
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
	Gerunds      Attribute `json:"gerunds"`
	Literary     Attribute `json:"literary"`
	Noun         Attribute `json:"noun"`
	Style        Attribute `json:"style"`
	ShortStyle   string    `json:"shortStyle"`
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

var promptTemplate = `Draw a {{.Adverb.Val}} {{.Adjective.Val}} {{.Noun.Val}} {{.Gerunds.Val}} and feeling {{.EmotionShort.Val}}{{.Ens}}.
Noun: {{.Noun.Val}}.
Emotion: {{.Emotion.Val}}.
Gerunds: {{.Gerunds.Val}}.
Primary style: {{.Style.Val}}.
Use only the colors {{.Color1.Val}} and {{.Color2.Val}}.
{{.Orientation.Val}}.
{{.Background.Val}}.
Expand upon the most relevant connotative meanings of {{.Noun.Val}}, {{.Emotion.Val}}, {{.Adjective.Val}}, and {{.Adverb.Val}}.
Find the representation that most closely matches the description.
Focus on the Noun, the Gerunds, the Emotion, and Primary style.{{.Literary.Val}}
DO NOT PUT ANY TEXT IN THE IMAGE.`

var dataTemplate = `
Address:     {{.Addr}} Ens: {{.Ens}}.
Seed:        {{.Seed}}.
Adverb:      {{.Adverb.Val}} Adjective: {{.Adjective.Val}} Noun: {{.Noun.Val}}.
Emotion:     {{.EmotionShort.Val}}.
Gerunds:	 {{.Gerunds.Val}}.
Literary:    {{.Literary.Val}}.
Style:       {{.Style.Val}}.
Color1:      {{.Color1.Val}} Color2: {{.Color2.Val}} Color3: {{.Color3.Val}}.
Variant1:    {{.Variant1.Val}} Variant2: {{.Variant2.Val}} Variant3: {{.Variant3.Val}}.
Background:  {{.Background.Val}}.
Orientation: {{.Orientation.Val}}.`

var terseTemplate = `{{.Adverb.Val}} {{.Adjective.Val}} {{.Noun.Val}} {{.Gerunds.Val}} and feeling {{.Emotion.Val}} in the style of {{.ShortStyle}}`

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

	// addr, _ := a.conn.GetEnsAddress(ensOrAddr)
	addr := ensOrAddr
	// if base.HexToAddress(addr) == base.ZeroAddr || !base.IsValidAddress(addr) {
	// 	return Dalledress{}, fmt.Errorf("ENS not registered: %s", ensOrAddr)
	// }
	hash := hexutil.Encode(crypto.Keccak256([]byte(addr)))
	seed := hash[2:] + addr[2:]
	if len(seed) < 104 {
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
		Emotion:      Attribute{Seed: seed[24:36]},
		Gerunds:      Attribute{Seed: seed[36:48]},
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
		Background:   Attribute{Seed: seed[20:32]},
		Orientation:  Attribute{Seed: seed[8:20]},
	}

	lengths := []int{
		len(a.adverbs),       // 0
		len(a.adjectives),    // 1
		len(a.emotionsShort), // 2
		len(a.literary),      // 3
		len(a.nouns),         // 4
		len(a.styles),        // 5
		len(a.colors),        // 6
		8,                    // 7
		len(a.gerunds),       // 8
	}

	dd.Adverb.Num = hexStringToBigIntModulo(dd.Adverb.Seed, SeedBump, lengths[0])
	dd.Adjective.Num = hexStringToBigIntModulo(dd.Adjective.Seed, SeedBump, lengths[1])
	dd.EmotionShort.Num = hexStringToBigIntModulo(dd.EmotionShort.Seed, 0, lengths[2])
	dd.Emotion.Num = hexStringToBigIntModulo(dd.Emotion.Seed, 0, lengths[2])
	dd.Gerunds.Num = hexStringToBigIntModulo(dd.Gerunds.Seed, 0, lengths[8])
	dd.Literary.Num = hexStringToBigIntModulo(dd.Literary.Seed, 0, lengths[3])
	dd.Noun.Num = hexStringToBigIntModulo(dd.Noun.Seed, 0, lengths[4])
	dd.Style.Num = hexStringToBigIntModulo(dd.Style.Seed, 0, lengths[5])
	dd.Style2.Num = hexStringToBigIntModulo(dd.Style2.Seed, 0, lengths[5])
	dd.Variant1.Num = hexStringToBigIntModulo(dd.Variant1.Seed, 0, lengths[5])
	dd.Variant2.Num = hexStringToBigIntModulo(dd.Variant2.Seed, 0, lengths[5])
	dd.Variant3.Num = hexStringToBigIntModulo(dd.Variant3.Seed, 0, lengths[5])
	dd.Color1.Num = hexStringToBigIntModulo(dd.Color1.Seed, 0, lengths[6])
	dd.Color2.Num = hexStringToBigIntModulo(dd.Color2.Seed, 0, lengths[6])
	dd.Color3.Num = hexStringToBigIntModulo(dd.Color3.Seed, 0, lengths[6])
	dd.Background.Num = hexStringToBigIntModulo(dd.Background.Seed, 0, lengths[7])
	dd.Orientation.Num = hexStringToBigIntModulo(dd.Orientation.Seed, 0, lengths[7])

	series := strings.Trim(file.AsciiFileToString("series.txt"), "\n\r")

	dd.Adverb.Val = a.adverbs[dd.Adverb.Num]
	dd.Adjective.Val = a.adjectives[dd.Adjective.Num]
	dd.EmotionShort.Val = a.emotionsShort[dd.EmotionShort.Num]
	dd.Emotion.Val = a.emotions[dd.Emotion.Num]
	if series == "postal" {
		dd.EmotionShort.Val = "postal going"
		dd.Emotion.Val = "postal going (becoming extremely and uncontrollably angry often to the point of violence)"
	} else if series == "happiness" {
		dd.EmotionShort.Val = "happiness"
		dd.Emotion.Val = "happiness (state of well-being characterized by emotions ranging from contentment to intense joy)"
	} else if series == "fury" || series == "steam" {
		dd.EmotionShort.Val = "fury and postal going"
		dd.Emotion.Val = "fury (wild or violent anger) and postal going (becoming extremely and uncontrollably angry often to the point of violence)"
	} else if series == "love" || series == "solar" {
		dd.EmotionShort.Val = "love and compassion"
		dd.Emotion.Val = "love (a strong positive emotion of regard and affection) and compassion (sympathetic pity and concern for the sufferings or misfortunes of others)"
	}
	dd.Gerunds.Val = a.gerunds[dd.Gerunds.Num] + " and " + a.gerunds[lengths[8]-dd.Gerunds.Num]
	dd.Literary.Val = " Write in the literary style {{.Literary.Val}}."
	if series == "human" || series == "human2" {
		dd.Noun.Val = "human"
	} else if strings.Contains(series, "human-with") || series == "postal" || series == "happiness" {
		dd.Noun.Val = "human with " + a.nouns[dd.Noun.Num] + " characteristics"
	} else if strings.Contains(series, "human-like") || series == "fury" || series == "love" || series == "steam" || series == "solar" {
		dd.Noun.Val = a.nouns[dd.Noun.Num] + " with human-like characteristics"
	} else {
		dd.Noun.Val = a.nouns[dd.Noun.Num]
	}
	dd.Style.Val = a.styles[dd.Style.Num]
	dd.ShortStyle = a.shortStyles[dd.Style.Num]
	if series == "steam" {
		dd.Style.Val = "steam punk,modern western art movements"
		dd.ShortStyle = "steam punk"
	} else if series == "solar" {
		dd.Style.Val = "solar punk,modern western art movements"
		dd.ShortStyle = "solar punk"
	}
	dd.Style2.Val = a.styles[dd.Style2.Num]
	dd.Color1.Val = a.colors[dd.Color1.Num]
	dd.Color2.Val = a.colors[dd.Color2.Num]
	dd.Color3.Val = "#" + muteColor(dd.Color3.Seed[:8], dd.Color3.Seed[4:4])
	dd.Variant1.Val = clip(a.styles[dd.Variant1.Num])
	dd.Variant2.Val = clip(a.styles[dd.Variant2.Num])
	dd.Variant3.Val = clip(a.styles[dd.Variant3.Num])

	e := os.Getenv("DALLE_NO_LITERARY")
	if e != "" {
		dd.Literary.Val = ""
	}
	dd.Literary.Val = strings.Replace(dd.Literary.Val, "{{.Literary.Val}}", a.literary[dd.Literary.Num], -1)

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

	e = os.Getenv("DALLE_BACKGROUND")
	if e != "" {
		switch e {
		case "solid":
			dd.Background.Val = "Make the image's background a solid color: {{.Color3.Val}}"
		default:
			logger.Fatal("Invalid DALLE_BACKGROUND: ", e)
		}
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
	e = os.Getenv("DALLE_ORIENTATION")
	if e != "" {
		dd.Orientation.Val = e
	}
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

func (a *App) GetImage(which int, ensOrAddr string) {
	// logger.Info("Generating image for", ensOrAddr)
	addr := ensOrAddr
	// if addr, _ := a.conn.GetEnsAddress(ensOrAddr); len(addr) < 42 { // base.HexToAddress(addr) == base.ZeroAddr || !base.IsValidAddress(addr) {
	// 	logger.Error(fmt.Errorf("ENS not registered: %s", ensOrAddr))
	// 	return
	// } else {
	folder := "./generated/"
	file.EstablishFolder(folder)
	file.EstablishFolder(strings.Replace(folder, "/generated", "/txt-prompt", -1))
	file.EstablishFolder(strings.Replace(folder, "/generated", "/txt-generated", -1))
	file.EstablishFolder(strings.Replace(folder, "/generated", "/annotated", -1))
	file.EstablishFolder(strings.Replace(folder, "/generated", "/stitched", -1))

	series := strings.Trim(file.AsciiFileToString("series.txt"), "\n\r")

	fn := filepath.Join(folder, fmt.Sprintf("%s-%s.png", addr, series)) // cnt))
	annoName := strings.Replace(fn, "/generated", "/annotated", -1)
	if file.FileExists(annoName) {
		logger.Info(colors.Yellow+"Image already exists: ", fn, colors.Off)
		time.Sleep(250 * time.Millisecond)
		utils.System("open " + annoName)
		return
	}
	a.nMade++

	logger.Info(colors.Cyan, addr, colors.Yellow, "- improving the prompt...", colors.Off)

	prompt, orig := a.GetImprovedPrompt(ensOrAddr)
	size := "1024x1024"
	if strings.Contains(prompt, "horizontal") {
		size = "1792x1024"
	} else if strings.Contains(prompt, "vertical") {
		size = "1024x1792"
	}

	quality := "standard"
	if os.Getenv("DALLE_QUALITY") != "" {
		quality = os.Getenv("DALLE_QUALITY")
	}

	url := "https://api.openai.com/v1/images/generations"
	payload := DalleRequest{
		Prompt:  prompt,
		N:       1,
		Quality: quality,
		Style:   "vivid",
		Model:   "dall-e-3",
		Size:    size,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("No OPENAI_API_KEY key found")
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
	promptFn := strings.Replace(strings.Replace(fn, ".png", ".txt", -1), "generated", "txt-prompt", -1)
	file.StringToAsciiFile(promptFn, orig)
	// utils.System("open " + txtFn)

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

	path, err := annotate(fn, "bottom", 0.2)
	if err != nil {
		fmt.Println("Error annotating image:", err)
		return
	}
	logger.Info(colors.Cyan, addr, colors.Green, "- image saved as", colors.White+strings.Trim(path, " "), fmt.Sprintf("%d", which), colors.Off)
	utils.System("open " + path)
	// }
}

func (a *App) GetModeration(ensOrAddr string) string {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("No OPENAI_API_KEY key found")
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

func (a *App) GetImprovedPrompt(ensOrAddr string) (string, string) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("No OPENAI_API_KEY key found")
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
		return fmt.Sprintf("Error: %s", err), ""
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Sprintf("Error: %s", err), ""
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("Error: %s", err), ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("Error: %s", err), ""
	}

	return string(body), prompt
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
