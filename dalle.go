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
