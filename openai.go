package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"git.sr.ht/~sbinet/gg"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
)

// type Category struct {
// 	Harassment            bool `json:"harassment"`
// 	HarassmentThreatening bool `json:"harassment/threatening"`
// 	Hate                  bool `json:"hate"`
// 	HateThreatening       bool `json:"hate/threatening"`
// 	Selfharm              bool `json:"self-harm"`
// 	SelfharmInstructions  bool `json:"self-harm/instructions"`
// 	SelfharmIntent        bool `json:"self-harm/intent"`
// 	Sexual                bool `json:"sexual"`
// 	SexualMinors          bool `json:"sexual/minors"`
// 	Violence              bool `json:"violence"`
// 	ViolenceGraphic       bool `json:"violence/graphic"`
// }

// type Score struct {
// 	Harassment            float64 `json:"harassment"`
// 	HarassmentThreatening float64 `json:"harassment/threatening"`
// 	Hate                  float64 `json:"hate"`
// 	HateThreatening       float64 `json:"hate/threatening"`
// 	Selfharm              float64 `json:"self-harm"`
// 	SelfharmInstructions  float64 `json:"self-harm/instructions"`
// 	SelfharmIntent        float64 `json:"self-harm/intent"`
// 	Sexual                float64 `json:"sexual"`
// 	SexualMinors          float64 `json:"sexual/minors"`
// 	Violence              float64 `json:"violence"`
// 	ViolenceGraphic       float64 `json:"violence/graphic"`
// }

// type Results struct {
// 	Flagged    bool       `json:"flagged"`
// 	Categories []Category `json:"categories"`
// 	Scores     []Score    `json:"scores"`
// }

// type ModerationObject struct {
// 	ID      string    `json:"id"`
// 	Model   string    `json:"model"`
// 	Results []Results `json:"results"`
// }

// type ChatMessage struct {
// 	Role    string `json:"role"`
// 	Content string `json:"content"`
// }

// type ChatRequest struct {
// 	Model    string        `json:"model"`
// 	Seed     int           `json:"seed"`
// 	Messages []ChatMessage `json:"messages"`
// }

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

func (dd *DalleDress) enhancePrompt() (string, error) {
	prompt := dd.Prompt
	url := "https://api.openai.com/v1/chat/completions"
	payload := DalleRequest{
		Model: "gpt-3.5-turbo",
		Seed:  1337,
	}
	payload.Messages = append(payload.Messages, Message{Role: "system", Content: prompt})

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	apiKey := os.Getenv("OPENAI_API_KEY")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (dd *DalleDress) getAiImage() error {
	folder := "./output/generated/"
	file.EstablishFolder(folder)
	file.EstablishFolder(strings.Replace(folder, "/generated", "/annotated", -1))
	file.EstablishFolder(strings.Replace(folder, "/generated", "/stitched", -1))

	// addr := ensOrAddr
	fn := filepath.Join(folder, fmt.Sprintf("%s-%s.png", dd.Orig, "")) //a.Series.Suffix))
	// annoName := strings.Replace(fn, "/generated", "/annotated", -1)
	// if file.FileExists(annoName) {
	// 	logger.Info(colors.Yellow+"Image already exists: ", fn, which, colors.Off)
	// 	time.Sleep(250 * time.Millisecond)
	// 	utils.System("open " + annoName)
	// 	return
	// }
	// a.nMade++

	logger.Info(colors.Cyan, dd.Orig, colors.Yellow, "- improving the prompt...", colors.Off)

	prompt := dd.EnhancedPrompt // , orig := a.GetImprovedPrompt(ensOrAddr)
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
		logger.Fatal("No OPENAI_API_KEY key found")
	}

	logger.Info(colors.Cyan, dd.Orig, colors.Yellow, "- generating the image...", colors.Off)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	bodyStr := string(body)
	bodyStr = strings.Replace(bodyStr, "\"revised_prompt\": \"", "\"revised_prompt\": \"NO TEXT. ", -1)
	bodyStr = strings.Replace(bodyStr, ".\",", ". NO TEXT.\",", 1)
	body = []byte(bodyStr)

	// logger.Info("DalleResponse: ", string(body))
	var dalleResp DalleResponse
	err = json.Unmarshal(body, &dalleResp)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Error:", resp.Status, resp.StatusCode, string(body))
	}

	if len(dalleResp.Data) == 0 {
		return fmt.Errorf("No images returned")
	}

	imageURL := dalleResp.Data[0].Url

	// Download the image
	imageResp, err := http.Get(imageURL)
	if err != nil {
		return err
	}
	defer imageResp.Body.Close()

	// txtFn := strings.Replace(strings.Replace(fn, ".png", ".txt", -1), "generated", "txt-generated", -1)
	// file.StringToAsciiFile(txtFn, prompt)
	// promptFn := strings.Replace(strings.Replace(fn, ".png", ".txt", -1), "generated", "txt-prompt", -1)
	// file.StringToAsciiFile(promptFn, dd.Orig)
	// utils.System("open " + txtFn)

	os.Remove(fn)
	file, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("Failed to open output file: ", fn)
	}
	defer file.Close()

	_, err = io.Copy(file, imageResp.Body)
	if err != nil {
		return err
	}

	path, err := dd.annotate(fn, "bottom", 0.2)
	if err != nil {
		return fmt.Errorf("Error annotating image:", err.Error())
	}
	logger.Info(colors.Cyan, dd.Orig, colors.Green, "- image saved as", colors.White+strings.Trim(path, " "), colors.Off)
	utils.System("open " + path)
	return nil
}

// annotate reads an image and adds a text annotation to it.
func (dd *DalleDress) annotate(fileName, location string, annoPct float64) (ret string, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	text := dd.TersePrompt
	estimatedFontSize := 30. * (float64(width) / float64(len(text)*7.))

	textWidth := float64(width) * 0.95
	lines := math.Ceil(float64(len(text)) / (textWidth / estimatedFontSize))
	marginHeight := float64(height) * 0.025
	annoHeight := lines * estimatedFontSize * 1.5
	newHeight := height + int(annoHeight+marginHeight*2)

	newImg := image.NewRGBA(image.Rect(0, 0, width, newHeight))
	draw.Draw(newImg, newImg.Bounds(), img, bounds.Min, draw.Src)

	bgColor, _ := findAverageDominantColor(img)
	col, err := parseHexColor(bgColor)
	if err != nil {
		return "", err
	}

	bgRect := image.Rect(0, height, width, newHeight)
	if location == "top" {
		bgRect = image.Rect(0, 0, width, int(annoHeight+marginHeight*2))
		draw.Draw(newImg, bgRect, &image.Uniform{col}, image.Point{}, draw.Src)
	} else {
		draw.Draw(newImg, bgRect, &image.Uniform{col}, image.Point{}, draw.Src)
	}

	gc := gg.NewContextForImage(newImg)
	if err := gc.LoadFontFace("/System/Library/Fonts/Monaco.ttf", estimatedFontSize); err != nil {
		log.Fatalf("Error loading font: %v", err) // Handle the error appropriately
	}
	borderCol := darkenColor(col)
	gc.SetColor(borderCol)
	gc.SetLineWidth(2)
	if location == "top" {
		gc.DrawLine(0, float64(height)*annoPct, float64(width), float64(height)*annoPct)
	} else {
		gc.DrawLine(0, float64(height), float64(width), float64(height))
	}
	gc.Stroke()

	// Draw the text with adjusted margins.
	textColor, _ := contrastColor(col)
	// fmt.Println(l, textColor, col)

	gc.SetColor(textColor) // use the contrasting color for the text
	gc.DrawStringWrapped(text, float64(width)/2, float64(height)+marginHeight*2, 0.5, 0.35, textWidth, 1.5, gg.AlignLeft)

	// Save the new image.
	outputPath := strings.Replace(fileName, "generated/", "annotated/", -1)
	out, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	err = png.Encode(out, gc.Image())
	if err != nil {
		return "", err
	}

	return outputPath, nil
}
