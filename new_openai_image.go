package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
)

func getAiImage(prompt, tersePrompt, hash string) error {
	folder := "./output/generated/"
	file.EstablishFolder(folder)
	file.EstablishFolder(strings.Replace(folder, "/generated", "/annotated", -1))
	file.EstablishFolder(strings.Replace(folder, "/generated", "/stitched", -1))

	t := time.Now()
	s := t.Format("20060102150405")
	fn := filepath.Join(folder, fmt.Sprintf("%s-%s.png", hash, s)) //a.Series.Suffix))

	logger.Info(colors.Cyan, hash, colors.Yellow, "- improving the prompt...", colors.Off)

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

	logger.Info(colors.Cyan, hash, colors.Yellow, "- generating the image...", colors.Off)

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

	var dalleResp DalleResponse
	err = json.Unmarshal(body, &dalleResp)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("error: %s %d %s", resp.Status, resp.StatusCode, string(body))
	}

	if len(dalleResp.Data) == 0 {
		return fmt.Errorf("no images returned")
	}

	imageURL := dalleResp.Data[0].Url

	imageResp, err := http.Get(imageURL)
	if err != nil {
		return err
	}
	defer imageResp.Body.Close()

	os.Remove(fn)
	file, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("failed to open output file: %s", fn)
	}
	defer file.Close()

	_, err = io.Copy(file, imageResp.Body)
	if err != nil {
		return err
	}

	path, err := annotate(tersePrompt, fn, "bottom", 0.2)
	if err != nil {
		return fmt.Errorf("error annotating image: %v", err)
	}
	logger.Info(colors.Cyan, hash, colors.Green, "- image saved as", colors.White+strings.Trim(path, " "), colors.Off)
	utils.System("open " + path)
	return nil
}
