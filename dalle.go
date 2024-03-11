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
	HexIn           string    `json:"hexIn"`
	Seed            string    `json:"seed"`
	Adverb          Attribute `json:"adverb"`
	Adjective       Attribute `json:"adjective"`
	Noun            Attribute `json:"noun"`
	Emotion         Attribute `json:"emotion"`
	EmotionShort    Attribute `json:"emotionShort"`
	Occupation      Attribute `json:"occupation"`
	OccupationShort Attribute `json:"occupationShort"`
	Gerunds         Attribute `json:"gerunds"`
	ArtStyle        Attribute `json:"artstyle"`
	ArtStyleShort   string    `json:"artStyleShort"`
	ArtStyle2       Attribute `json:"artstyle2"`
	LitStyle        Attribute `json:"litstyle"`
	Color1          Attribute `json:"color1"`
	Color2          Attribute `json:"color2"`
	Color3          Attribute `json:"color3"`
	Background      Attribute `json:"background"`
	Orientation     Attribute `json:"orientation"`
	Prompt          Value     `json:"prompt"`
	Data            Value     `json:"data"`
	Terse           Value     `json:"terse"`
}

var promptTemplate = `
Draw a {{.Adverb.Val}} {{.Adjective.Val}} {{.Noun.Val}} who works as a 
{{.Occupation.Val}} and is {{.Gerunds.Val}} and feeling {{.Emotion.Val}}.
Noun: {{.Noun.Val}}.
Emotion: {{.Emotion.Val}}.
Occupation: {{.Occupation.Val}}.
Gerunds: {{.Gerunds.Val}}.
Artistic style: {{.ArtStyle.Val}}.
Use only the colors {{.Color1.Val}} and {{.Color2.Val}}.
{{.Orientation.Val}}.
{{.Background.Val}}.
Expand upon the most relevant connotative meanings of {{.Noun.Val}}, {{.Emotion.Val}}, {{.Adjective.Val}}, and {{.Adverb.Val}}.
Find the representation that most closely matches the description.
Focus on the noun, the occupation, the emotion, and literary style.{{.LitStyle.Val}}`

var dataTemplate = `
HexIn:       {{.HexIn}}.
Seed:        {{.Seed}}.
Adverb:      {{.Adverb.Val}} Adjective: {{.Adjective.Val}} Noun: {{.Noun.Val}}.
Emotion:     {{.EmotionShort.Val}}.
Occupation:  {{.OccupationShort.Val}}.
Gerunds:	 {{.Gerunds.Val}}.
LitStyle:    {{.LitStyle.Val}}.
ArtStyle:    {{.ArtStyle.Val}}.
Color1:      {{.Color1.Val}} Color2: {{.Color2.Val}} Color3: {{.Color3.Val}}.
Background:  {{.Background.Val}}.
Orientation: {{.Orientation.Val}}.`

var terseTemplate = `{{.Adverb.Val}} {{.Adjective.Val}} {{.Noun.Val}} {{.OccupationShort.Val}} {{.Gerunds.Val}} and feeling {{.EmotionShort.Val}} in the style of {{.ArtStyleShort}}`

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

func (a *App) GetDalledress(hexIn string) (Dalledress, error) {
	if len(a.Adverbs) == 0 {
		return Dalledress{}, fmt.Errorf("adverbs not loaded")
	}

	hash1 := hexutil.Encode(crypto.Keccak256([]byte(hexIn)))
	hash2 := reverseString(hash1)
	hexIn2 := reverseString(hexIn)
	seed := hash1[2:] + hash2[2:] + hexIn[2:] + hexIn2[2:] // 64 + 64 + 32 + 32 = 192
	if len(seed) < 165 {
		return Dalledress{}, fmt.Errorf("invalid seed (too short): %s", seed)
	}

	dd := Dalledress{
		HexIn:           hexIn,
		Seed:            seed,
		Adverb:          Attribute{Seed: seed[0:12]},    // 12-0 = 12
		Adjective:       Attribute{Seed: seed[12:24]},   // 24-12 = 12
		Noun:            Attribute{Seed: seed[24:36]},   // 36-24 = 12
		Emotion:         Attribute{Seed: seed[36:48]},   // 48-36 = 12
		EmotionShort:    Attribute{Seed: seed[48:60]},   // 60-48 = 12
		Occupation:      Attribute{Seed: seed[60:72]},   // 72-60 = 12
		OccupationShort: Attribute{Seed: seed[72:84]},   // 84-72 = 12
		Gerunds:         Attribute{Seed: seed[84:96]},   // 96-84 = 12
		ArtStyle:        Attribute{Seed: seed[96:108]},  // 108-96 = 12
		ArtStyle2:       Attribute{Seed: seed[108:120]}, // 120-108 = 12
		LitStyle:        Attribute{Seed: seed[120:132]}, // 132-120 = 12
		Color1:          Attribute{Seed: seed[132:144]}, // 144-132 = 12
		Color2:          Attribute{Seed: seed[144:156]}, // 156-144 = 12
		Color3:          Attribute{Seed: seed[156:168]}, // 168-156 = 12
		Background:      Attribute{Seed: seed[168:180]}, // 180-168 = 12
		Orientation:     Attribute{Seed: seed[180:192]}, // 192-180 = 12
	}

	lengths := map[string]int{
		"adverbs":     len(a.Adverbs),
		"adjectives":  len(a.Adjectives),
		"nouns":       len(a.Nouns),
		"emotions":    len(a.Emotions),
		"occupations": len(a.Occupations),
		"gerunds":     len(a.Gerunds),
		"artstyles":   len(a.Artstyles),
		"litstyles":   len(a.Litstyles),
		"colors":      len(a.Colors),
		"other":       8,
	}

	dd.Adverb.Num = hexStringToBigIntModulo(dd.Adverb.Seed, SeedBump, lengths["adverbs"])
	dd.Adjective.Num = hexStringToBigIntModulo(dd.Adjective.Seed, SeedBump, lengths["adjectives"])
	dd.Noun.Num = hexStringToBigIntModulo(dd.Noun.Seed, 0, lengths["nouns"])
	dd.Emotion.Num = hexStringToBigIntModulo(dd.Emotion.Seed, 0, lengths["emotions"])
	dd.EmotionShort.Num = hexStringToBigIntModulo(dd.EmotionShort.Seed, 0, lengths["emotions"])
	dd.Occupation.Num = hexStringToBigIntModulo(dd.Occupation.Seed, 0, lengths["occupations"])
	dd.OccupationShort.Num = hexStringToBigIntModulo(dd.OccupationShort.Seed, 0, lengths["occupations"])
	dd.Gerunds.Num = hexStringToBigIntModulo(dd.Gerunds.Seed, 0, lengths["gerunds"])
	dd.ArtStyle.Num = hexStringToBigIntModulo(dd.ArtStyle.Seed, 0, lengths["artstyles"])
	dd.ArtStyle2.Num = hexStringToBigIntModulo(dd.ArtStyle2.Seed, 0, lengths["artstyles"])
	dd.LitStyle.Num = hexStringToBigIntModulo(dd.LitStyle.Seed, 0, lengths["litstyles"])
	dd.Color1.Num = hexStringToBigIntModulo(dd.Color1.Seed, 0, lengths["colors"])
	dd.Color2.Num = hexStringToBigIntModulo(dd.Color2.Seed, 0, lengths["colors"])
	dd.Color3.Num = hexStringToBigIntModulo(dd.Color3.Seed, 0, lengths["colors"])
	dd.Background.Num = hexStringToBigIntModulo(dd.Background.Seed, 0, lengths["other"])
	dd.Orientation.Num = hexStringToBigIntModulo(dd.Orientation.Seed, 0, lengths["other"])

	dd.Adverb.Val = a.Adverbs[dd.Adverb.Num]
	dd.Adjective.Val = a.Adjectives[dd.Adjective.Num]
	dd.Noun.Val = a.Nouns[dd.Noun.Num] + " with human-like characteristics"
	dd.Emotion.Val = a.Emotions[dd.Emotion.Num]
	dd.Occupation.Val = a.Occupations[dd.Occupation.Num]
	dd.Gerunds.Val = a.Gerunds[dd.Gerunds.Num] + " and " + a.Gerunds[lengths["gerunds"]-dd.Gerunds.Num]
	dd.ArtStyle.Val = a.Artstyles[dd.ArtStyle.Num]
	dd.ArtStyle2.Val = a.Artstyles[dd.ArtStyle2.Num]
	dd.LitStyle.Val = " Write in the literary style " + a.Litstyles[dd.LitStyle.Num]
	dd.Color1.Val = strings.Replace(a.Colors[dd.Color1.Num], ",", " (", -1) + ")"
	dd.Color2.Val = strings.Replace(a.Colors[dd.Color2.Num], ",", " (", -1) + ")"
	dd.Color3.Val = "#" + muteColor(dd.Color3.Seed[:8], dd.Color3.Seed[4:4])

	// Noneable
	if dd.Occupation.Val == "none" {
		dd.Occupation.Val = ""
	}
	if dd.Color1.Val == "none" {
		dd.Color1.Val = ""
	}
	if dd.Color2.Val == "none" {
		dd.Color2.Val = ""
	}
	if dd.Color3.Val == "none" {
		dd.Color3.Val = ""
	}

	parts := strings.Split(dd.Emotion.Val, " (")
	if len(parts) > 0 {
		dd.EmotionShort.Val = parts[0]
	} else {
		dd.EmotionShort.Val = dd.Emotion.Val
	}

	parts = strings.Split(dd.ArtStyle.Val, ",")
	if len(parts) > 0 {
		dd.ArtStyleShort = parts[0]
	} else {
		dd.ArtStyleShort = dd.ArtStyle.Val
	}

	parts = strings.Split(dd.Occupation.Val, " (")
	if len(parts) > 0 {
		dd.OccupationShort.Val = parts[0]
	} else {
		dd.OccupationShort.Val = dd.Occupation.Val
	}

	switch dd.Background.Num {
	case 0:
		dd.Background.Val = "The background should be transparent"
	case 1:
		dd.Background.Val = "The background should be this color {{.Color3.Val}} and pay homage to this style {{.ArtStyle2.Val}}"
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

	e := os.Getenv("DALLE_BACKGROUND")
	if e != "" {
		switch e {
		case "solid":
			dd.Background.Val = "Make the image's background a solid color: {{.Color3.Val}}"
		default:
			logger.Fatal("Invalid DALLE_BACKGROUND: ", e)
		}
	}
	dd.Background.Val = strings.Replace(dd.Background.Val, "{{.Color3.Val}}", dd.Color3.Val, -1)
	dd.Background.Val = strings.Replace(dd.Background.Val, "{{.ArtStyle2.Val}}", dd.ArtStyle2.Val, -1)

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
	folder := "./generated/"
	file.EstablishFolder(folder)
	file.EstablishFolder(strings.Replace(folder, "/generated", "/txt-prompt", -1))
	file.EstablishFolder(strings.Replace(folder, "/generated", "/txt-generated", -1))
	file.EstablishFolder(strings.Replace(folder, "/generated", "/annotated", -1))
	file.EstablishFolder(strings.Replace(folder, "/generated", "/stitched", -1))

	addr := ensOrAddr
	fn := filepath.Join(folder, fmt.Sprintf("%s-%s.png", addr, a.Series.Suffix))
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
	bodyStr := string(body)
	bodyStr = strings.Replace(bodyStr, "\"revised_prompt\": \"", "\"revised_prompt\": \"NO TEXT. ", -1)
	bodyStr = strings.Replace(bodyStr, ".\",", ". NO TEXT.\",", 1)
	body = []byte(bodyStr)

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

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
