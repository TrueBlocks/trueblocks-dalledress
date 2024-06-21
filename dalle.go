package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"log"
// 	"math/big"
// 	"net/http"
// 	"os"
// 	"path/filepath"
// 	"strings"
// 	"text/template"
// 	"time"

// 	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
// 	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
// 	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
// 	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
// 	"github.com/ethereum/go-ethereum/common/hexutil"
// 	"github.com/ethereum/go-ethereum/crypto"
// 	"github.com/lucasb-eyer/go-colorful"
// )

// type Value struct {
// 	Val   string `json:"val"`
// 	Error error  `json:"error"`
// }

// type Dalledress struct {
// 	Orig            string    `json:"orig"`
// 	Seed            string    `json:"seed"`
// 	Adverb          Attribute `json:"adverb"`
// 	Adjective       Attribute `json:"adjective"`
// 	Noun            Attribute `json:"noun"`
// 	Emotion         Attribute `json:"emotion"`
// 	E motionShort    Attribute `json:"emotionShort"`
// 	Occupation      Attribute `json:"occupation"`
// 	Occupat ionShort Attribute `json:"occupationShort"`
// 	Gerunds         Attribute `json:"gerunds"`
// 	A rtStyle        Attribute `json:"artstyle"`
// 	A rtStyleShort   string    `json:"artStyleShort"`
// 	A rtStyle2       Attribute `json:"artstyle2"`
// 	L itStyle        Attribute `json:"litstyle"`
// 	Color1          Attribute `json:"color1"`
// 	Color2          Attribute `json:"color2"`
// 	Color3          Attribute `json:"color3"`
// 	Background      Attribute `json:"background"`
// 	Orientation     Attribute `json:"orientation"`
// 	P rompt          Value     `json:"prompt"`
// 	Data            Value     `json:"data"`
// 	Terse           Value     `json:"terse"`
// }

// func (d *Dalledress) gener atePrompt(t *template.Template, f func(s string) string) (string, error) {
// 	var buffer bytes.Buffer
// 	if err := t.Execute(&buffer, d); err != nil {
// 		return "", err
// 	}
// 	return f(buffer.String()), nil
// }

// func (a *App) GetPrompt(ensOrAddr string) string {
// 	clean := func(s string) string {
// 		return strings.Replace(strings.Replace(strings.Replace(s, ". ", ".\n", -1), "|", "\t", -1), "+\n", " ", -1)
// 	}
// 	if dd, err := a.GetDalledress(ensOrAddr); err != nil {
// 		return fmt.Sprintf("Error generating prompt: %s\n", err)
// 	} else {
// 		if dd.P rompt.Val, dd.P rompt.Error = dd.gener atePrompt(a.p Template, clean); dd.P rompt.Error != nil {
// 			return dd.P rompt.Error.Error()
// 		}
// 		return dd.P rompt.Val
// 	}
// }

// func (a *App) GetTerse(ensOrAddr string) string {
// 	clean := func(s string) string {
// 		return strings.Replace(strings.Replace(strings.Replace(s, ". ", "\n", -1), "|", "\t", -1), "+\n", " ", -1)
// 	}
// 	if dd, err := a.GetDalledress(ensOrAddr); err != nil {
// 		return fmt.Sprintf("Error generating data: %s\n", err)
// 	} else {
// 		if dd.Terse.Val, dd.Terse.Error = dd.genera tePrompt(a.t Template, clean); dd.Terse.Error != nil {
// 			return dd.Terse.Error.Error()
// 		}
// 		return dd.Terse.Val
// 	}
// }

// func (a *App) GetData(ensOrAddr string) string {
// 	clean := func(s string) string {
// 		return strings.Replace(strings.Replace(strings.Replace(s, ". ", ".\n", -1), "|", "\t", -1), "+\n", " ", -1)
// 	}
// 	if dd, err := a.GetDalledress(ensOrAddr); err != nil {
// 		return fmt.Sprintf("Error generating data: %s\n", err)
// 	} else {
// 		if dd.Data.Val, dd.Data.Error = dd.generate P rompt(a.d Template, clean); dd.Data.Error != nil {
// 			return dd.Data.Error.Error()
// 		}
// 		return dd.Data.Val
// 	}
// }

// func (a *App) GetJson(ensOrAddr string) string {
// 	if dd, err := a.GetDalledress(ensOrAddr); err != nil {
// 		return fmt.Sprintf("{\"error\": \"%s\"}", err)
// 	} else {
// 		bytes, _ := json.Marshal(dd)
// 		return strings.Replace(string(bytes), ",", ", ", -1)
// 	}
// }

// func (a *App) GetAddress(index int) string {
// 	if len(a.addresses) == 0 {
// 		a.addresses = file.AsciiFileToLines("/Users/jrush/Desktop/Animals.1/addresses.csv")
// 	}
// 	return a.addresses[index%len(a.addresses)]
// }

// func hexStringToBigIntModulo(hexString string, seedBump, modulo int) uint64 {
// 	intValue := new(big.Int)
// 	seedValue := new(big.Int)
// 	intValue.SetString(hexString, 16)
// 	seedValue.SetInt64(int64(seedBump))
// 	modValue := big.NewInt(int64(modulo))
// 	if modValue == big.NewInt(0) {
// 		modValue = big.NewInt(1)
// 	}
// 	intValue = intValue.Add(intValue, seedValue)
// 	return uint64(intValue.Mod(intValue, modValue).Int64())
// }

// // func clip(s string) string {
// // 	parts := strings.Split(s, " from the ")
// // 	return parts[0]
// // }

// func (a *App) GetDalledress(hexIn string) (Dalledress, error) {
// 	if len(a.Adverbs) == 0 {
// 		return Dalledress{}, fmt.Errorf("adverbs not loaded")
// 	}

// 	hash1 := hexutil.Encode(crypto.Keccak256([]byte(hexIn)))
// 	hash2 := revers eString(hash1)
// 	hexIn2 := revers eString(hexIn)
// 	seed := hash1[2:] + hash2[2:] + hexIn[2:] + hexIn2[2:] // 64 + 64 + 32 + 32 = 192
// 	if len(seed) < 165 {
// 		return Dalledress{}, fmt.Errorf("invalid seed (too short): %s", seed)
// 	}

// 	dd := Dalledress{
// 		Orig:            hexIn,
// 		Seed:            seed,
// 		Adverb:          Attribute{Bytes: seed[0:12]},    // 12-0 = 12
// 		Adjective:       Attribute{Bytes: seed[12:24]},   // 24-12 = 12
// 		Noun:            Attribute{Bytes: seed[24:36]},   // 36-24 = 12
// 		Emotion:         Attribute{Bytes: seed[36:48]},   // 48-36 = 12
// 		E motionShort:    Attribute{Bytes: seed[48:60]},   // 60-48 = 12
// 		Occupation:      Attribute{Bytes: seed[60:72]},   // 72-60 = 12
// 		Occupati onShort: Attribute{Bytes: seed[72:84]},   // 84-72 = 12
// 		Gerunds:         Attribute{Bytes: seed[84:96]},   // 96-84 = 12
// 		A rtStyle:        Attribute{Bytes: seed[96:108]},  // 108-96 = 12
// 		A rtStyle2:       Attribute{Bytes: seed[108:120]}, // 120-108 = 12
// 		L itStyle:        Attribute{Bytes: seed[120:132]}, // 132-120 = 12
// 		Color1:          Attribute{Bytes: seed[132:144]}, // 144-132 = 12
// 		Color2:          Attribute{Bytes: seed[144:156]}, // 156-144 = 12
// 		Color3:          Attribute{Bytes: seed[156:168]}, // 168-156 = 12
// 		Background:      Attribute{Bytes: seed[168:180]}, // 180-168 = 12
// 		Orientation:     Attribute{Bytes: seed[180:192]}, // 192-180 = 12
// 	}

// 	lengths := map[string]int{
// 		"adverbs":     len(a.Adverbs),
// 		"adjectives":  len(a.Adjectives),
// 		"nouns":       len(a.Nouns),
// 		"emotions":    len(a.Emotions),
// 		"occupations": len(a.Occupations),
// 		"gerunds":     len(a.Gerunds),
// 		"artstyles":   len(a.Artstyles),
// 		"litstyles":   len(a.Litstyles),
// 		"colors":      len(a.Colors),
// 		"other":       8,
// 	}

// 	dd.A dverb.Selector = hexStringToBigIntModulo(dd.A dverb.Bytes, SeedBump, lengths["adverbs"])
// 	dd.A djective.Selector = hexStringToBigIntModulo(dd.A djective.Bytes, SeedBump, lengths["adjectives"])
// 	dd.N oun.Selector = hexStringToBigIntModulo(dd.N oun.Bytes, 0, lengths["nouns"])
// 	dd.E motion.Selector = hexStringToBigIntModulo(dd.E motion.Bytes, 0, lengths["emotions"])
// 	dd.E motionShort.Selector = hexStringToBigIntModulo(dd.E motionShort.Bytes, 0, lengths["emotions"])
// 	dd.O ccupation.Selector = hexStringToBigIntModulo(dd.O ccupation.Bytes, 0, lengths["occupations"])
// 	dd.O ccupationShort.Selector = hexStringToBigIntModulo(dd.O ccupationShort.Bytes, 0, lengths["occupations"])
// 	dd.Gerunds.Selector = hexStringToBigIntModulo(dd.Gerunds.Bytes, 0, lengths["gerunds"])
// 	dd.A rtStyle.Selector = hexStringToBigIntModulo(dd.A rtStyle.Bytes, 0, lengths["artstyles"])
// 	dd.A rtStyle2.Selector = hexStringToBigIntModulo(dd.A rtStyle2.Bytes, 0, lengths["artstyles"])
// 	dd.L itStyle.Selector = hexStringToBigIntModulo(dd.L itStyle.Bytes, 0, lengths["litstyles"])
// 	dd.C olor1.Selector = hexStringToBigIntModulo(dd.C olor1.Bytes, 0, lengths["colors"])
// 	dd.C olor2.Selector = hexStringToBigIntModulo(dd.C olor2.Bytes, 0, lengths["colors"])
// 	dd.C olor3.Selector = hexStringToBigIntModulo(dd.C olor3.Bytes, 0, lengths["colors"])
// 	dd.Background.Selector = hexStringToBigIntModulo(dd.Background.Bytes, 0, lengths["other"])
// 	dd.O rientation.Selector = hexStringToBigIntModulo(dd.O rientation.Bytes, 0, lengths["other"])

// 	dd.A dverb.Value = a.Adverbs[dd.A dverb.Selector]
// 	dd.A djective.Value = a.Adjectives[dd.A djective.Selector]
// 	dd.N oun.Value = a.Nouns[dd.N oun.Selector] + " with human-like characteristics"
// 	dd.E motion.Value = a.Emotions[dd.E motion.Selector]
// 	dd.O ccupation.Value = a.Occupations[dd.O ccupation.Selector]
// 	dd.Gerunds.Value = a.Gerunds[dd.Gerunds.Selector] + " and " + a.Gerunds[lengths["gerunds"]-int(dd.Gerunds.Selector)-1]
// 	dd.A rtStyle.Value = a.Artstyles[dd.A rtStyle.Selector]
// 	dd.A rtStyle2.Value = a.Artstyles[dd.A rtStyle2.Selector]
// 	dd.L itStyle.Value = " Write in the literary style " + a.Litstyles[dd.L itStyle.Selector]
// 	dd.Color1.Value = strings.Replace(a.Colors[dd.Color1.Selector], ",", " (", -1) + ")"
// 	dd.Color2.Value = strings.Replace(a.Colors[dd.Color2.Selector], ",", " (", -1) + ")"
// 	dd.Color3.Value = strings.Replace(a.Colors[dd.Color2.Selector], ",", " (", -1) + ")"

// 	// Noneable
// 	if dd.O ccupation.Value == "none" {
// 		dd.O ccupation.Value = ""
// 	}
// 	if dd.Color1.Value == "none" {
// 		dd.Color1.Value = ""
// 	}
// 	if dd.Color2.Value == "none" {
// 		dd.Color2.Value = ""
// 	}
// 	if dd.Color3.Value == "none" {
// 		dd.Color3.Value = ""
// 	}

// 	parts := strings.Split(dd.E motion.Value, " (")
// 	if len(parts) > 0 {
// 		dd.E motionShort.Value = parts[0]
// 	} else {
// 		dd.E motionShort.Value = dd.E motion.Value
// 	}

// 	parts = strings.Split(dd.A rtStyle.Value, ",")
// 	if len(parts) > 0 {
// 		dd.Art StyleShort = parts[0]
// 	} else {
// 		dd.Art StyleShort = dd.A rtStyle.Value
// 	}

// 	parts = strings.Split(dd.O ccupation.Value, " (")
// 	if len(parts) > 0 {
// 		dd.Occupat ionShort.Value = parts[0]
// 	} else {
// 		dd.Occupat ionShort.Value = dd.O ccupation.Value
// 	}

// 	switch dd.Background.Selector {
// 	case 0:
// 		dd.Background.Value = "The background should be transparent"
// 	case 1:
// 		dd.Background.Value = "The background should be this color {{.Color3.Value}} and pay homage to this style {{.ArtStyle2.Value}}"
// 	case 2:
// 		dd.Background.Value = " The background should be this color {{.Color3.Value}} and subtly patterned"
// 	case 3:
// 		fallthrough
// 	case 4:
// 		fallthrough
// 	case 5:
// 		fallthrough
// 	case 6:
// 		fallthrough
// 	case 7:
// 		dd.Background.Value = "The background should be solid and colored with this color: {{.Color3.Value}}"
// 	default:
// 		logger.Fatal("Invalid background number: ", dd.Background.Selector)
// 	}

// 	e := os.Getenv("DALLE_BACKGROUND")
// 	if e != "" {
// 		switch e {
// 		case "solid":
// 			dd.Background.Value = "Make the image's background a solid color: {{.Color3.Value}}"
// 		default:
// 			logger.Fatal("Invalid DALLE_BACKGROUND: ", e)
// 		}
// 	}
// 	dd.Background.Value = strings.Replace(dd.Background.Value, "{{.Color3.Value}}", dd.Color3.Value, -1)
// 	dd.Background.Value = strings.Replace(dd.Background.Value, "{{.ArtStyle2.Value}}", dd.ArtStyle2.Value, -1)

// 	ori, gaze, sym := "", "", ""
// 	switch dd.O rientation.Selector {
// 	case 0:
// 		ori, gaze, sym = "vertically", "into the camera", "symmetrically"
// 	case 1:
// 		ori, gaze, sym = "horizontally", "to the left", "asymmetrically"
// 	case 2:
// 		ori, gaze, sym = "diagonally", "into the camera", "symmetrically"
// 	case 3:
// 		ori, gaze, sym = "in a unique way", "to the right", "asymmetrically"
// 	case 4:
// 		ori, gaze, sym = "vertically", "into the camera", "symmetrically"
// 	case 5:
// 		ori, gaze, sym = "horizontally", "to the right", "asymmetrically"
// 	case 6:
// 		ori, gaze, sym = "diagonally", "into the camera", "symmetrically"
// 	case 7:
// 		ori, gaze, sym = "in a unique way", "to the left", "asymmetrically"
// 	default:
// 		logger.Fatal("Invalid orientation number: ", dd.O rientation.Selector)
// 	}
// 	dd.O rientation.Value = "Orient the scene {Ori} and {Sym} and make sure the {{.N oun.Value}} is facing {Gaze}"
// 	e = os.Getenv("DALLE_ORIENTATION")
// 	if e != "" {
// 		dd.O rientation.Value = e
// 	}
// 	dd.O rientation.Value = strings.Replace(dd.O rientation.Value, "{Ori}", ori, -1)
// 	dd.O rientation.Value = strings.Replace(dd.O rientation.Value, "{Sym}", sym, -1)
// 	dd.O rientation.Value = strings.Replace(dd.O rientation.Value, "{Gaze}", gaze, -1)
// 	dd.O rientation.Value = strings.Replace(dd.O rientation.Value, "{{.N oun.Value}}", dd.N oun.Value, -1)
// 	return dd, nil
// }

// type Message struct {
// 	Role    string `json:"role"`
// 	Content string `json:"content"`
// }

// type DalleRequest struct {
// 	Input    string    `json:"input,omitempty"`
// 	P rompt   string    `json:"prompt,omitempty"`
// 	N        int       `json:"n,omitempty"`
// 	Quality  string    `json:"quality,omitempty"`
// 	Model    string    `json:"model,omitempty"`
// 	Style    string    `json:"style,omitempty"`
// 	Size     string    `json:"size,omitempty"`
// 	Seed     int       `json:"seed,omitempty"`
// 	Messages []Message `json:"messages,omitempty"`
// }

// type DalleResponse struct {
// 	Data []struct {
// 		Url string `json:"url"`
// 	} `json:"data"`
// }

// var SeedBump = int(0)

// func (a *App) GetImage(which int, ensOrAddr string) {
// 	folder := "./output/generated/"
// 	file.EstablishFolder(folder)
// 	file.EstablishFolder(strings.Replace(folder, "/generated", "/txt-prompt", -1))
// 	file.EstablishFolder(strings.Replace(folder, "/generated", "/txt-generated", -1))
// 	file.EstablishFolder(strings.Replace(folder, "/generated", "/annotated", -1))
// 	file.EstablishFolder(strings.Replace(folder, "/generated", "/stitched", -1))

// 	addr := ensOrAddr
// 	fn := filepath.Join(folder, fmt.Sprintf("%s-%s.png", addr, a.Series.Suffix))
// 	annoName := strings.Replace(fn, "/generated", "/annotated", -1)
// 	if file.FileExists(annoName) {
// 		logger.Info(colors.Yellow+"Image already exists: ", fn, which, colors.Off)
// 		time.Sleep(250 * time.Millisecond)
// 		utils.System("open " + annoName)
// 		return
// 	}
// 	a.nMade++

// 	logger.Info(colors.Cyan, addr, colors.Yellow, "- improving the prompt...", which, colors.Off)

// 	prompt, orig := a.GetImprovedPrompt(ensOrAddr)
// 	size := "1024x1024"
// 	if strings.Contains(prompt, "horizontal") {
// 		size = "1792x1024"
// 	} else if strings.Contains(prompt, "vertical") {
// 		size = "1024x1792"
// 	}

// 	quality := "standard"
// 	if os.Getenv("DALLE_QUALITY") != "" {
// 		quality = os.Getenv("DALLE_QUALITY")
// 	}

// 	url := "https://api.openai.com/v1/images/generations"
// 	payload := DalleRequest{
// 		P rompt:  prompt,
// 		N:       1,
// 		Quality: quality,
// 		Style:   "vivid",
// 		Model:   "dall-e-3",
// 		Size:    size,
// 	}

// 	payloadBytes, err := json.Marshal(payload)
// 	if err != nil {
// 		panic(err)
// 	}

// 	apiKey := os.Getenv("OPENAI_API_KEY")
// 	if apiKey == "" {
// 		log.Fatal("No OPENAI_API_KEY key found")
// 	}

// 	logger.Info(colors.Cyan, addr, colors.Yellow, "- generating the image...", which, colors.Off)

// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
// 	if err != nil {
// 		panic(err)
// 	}
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Authorization", "Bearer "+apiKey)

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		panic(err)
// 	}
// 	bodyStr := string(body)
// 	bodyStr = strings.Replace(bodyStr, "\"revised_prompt\": \"", "\"revised_prompt\": \"NO TEXT. ", -1)
// 	bodyStr = strings.Replace(bodyStr, ".\",", ". NO TEXT.\",", 1)
// 	body = []byte(bodyStr)

// 	// logger.Info("DalleResponse: ", string(body))
// 	var dalleResp DalleResponse
// 	err = json.Unmarshal(body, &dalleResp)
// 	if err != nil {
// 		panic(err)
// 	}

// 	if resp.StatusCode != 200 {
// 		fmt.Println("Error:", resp.Status, resp.StatusCode, string(body))
// 		return
// 	}

// 	if len(dalleResp.Data) == 0 {
// 		fmt.Println("No images returned")
// 		return
// 	}

// 	imageURL := dalleResp.Data[0].Url

// 	// Download the image
// 	imageResp, err := http.Get(imageURL)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer imageResp.Body.Close()

// 	txtFn := strings.Replace(strings.Replace(fn, ".png", ".txt", -1), "generated", "txt-generated", -1)
// 	file.StringToAsciiFile(txtFn, prompt)
// 	promptFn := strings.Replace(strings.Replace(fn, ".png", ".txt", -1), "generated", "txt-prompt", -1)
// 	file.StringToAsciiFile(promptFn, orig)
// 	// utils.System("open " + txtFn)

// 	os.Remove(fn)
// 	file, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
// 	if err != nil {
// 		logger.Error("Failed to open output file: ", fn)
// 		panic(err)
// 	}
// 	defer file.Close()

// 	_, err = io.Copy(file, imageResp.Body)
// 	if err != nil {
// 		panic(err)
// 	}

// 	path, err := annotate(fn, "bottom", 0.2)
// 	if err != nil {
// 		fmt.Println("Error annotating image:", err)
// 		return
// 	}
// 	logger.Info(colors.Cyan, addr, colors.Green, "- image saved as", colors.White+strings.Trim(path, " "), fmt.Sprintf("%d", which), colors.Off)
// 	utils.System("open " + path)
// 	// }
// }

// func (a *App) GetImprovedPrompt(ensOrAddr string) (string, string) {
// 	apiKey := os.Getenv("OPENAI_API_KEY")
// 	if apiKey == "" {
// 		log.Fatal("No OPENAI_API_KEY key found")
// 	}
// 	// fmt.Println("API key found", apiKey)

// 	prompt := a.GetPrompt(ensOrAddr)
// 	prompt = "Please give me a detailed, but terse, rewriting of this description of an image. Be imaginative. " + prompt + " Mute the colors."
// 	url := "https://api.openai.com/v1/chat/completions"
// 	payload := DalleRequest{
// 		Model: "gpt-3.5-turbo",
// 		Seed:  1337,
// 	}
// 	payload.Messages = append(payload.Messages, Message{Role: "system", Content: prompt})

// 	payloadBytes, err := json.Marshal(payload)
// 	if err != nil {
// 		return fmt.Sprintf("Error: %s", err), ""
// 	}

// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
// 	if err != nil {
// 		return fmt.Sprintf("Error: %s", err), ""
// 	}
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Authorization", "Bearer "+apiKey)

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return fmt.Sprintf("Error: %s", err), ""
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return fmt.Sprintf("Error: %s", err), ""
// 	}

// 	return string(body), prompt
// }
