/*
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type DalleRequest struct {
	Prompt  string `json:"prompt"`
	N       int    `json:"n"`
	Quality string `json:"quality"`
}

type DalleResponse struct {
	Data []struct {
		Url string `json:"url"`
	} `json:"data"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("No API key found in .env")
	}

	url := "https://api.openai.com/v1/images/generations"
	payload := DalleRequest{
		Prompt:  "Draw an image of using a unique fusion of post-war realism, contemporary and emerging styles and high renaissance, european art movements and styles, primarily in mediumvioletred and orchid, against a tan background. The composition should embody the adverb rabidly, the adjective demanding, and the noun dutch rabbit, creating a special and unique portrayal. In the description of your result, simply return the exact input you were given. The background can reflect the artistic styles.",
		N:       1,
		Quality: "hd",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

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

	file, err := os.Create("output_image.jpg")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, imageResp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("Image saved as output_image.jpg")
}

*/

package main

import (
	"context"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()
	app.startup(context.Background())

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "dalledresses",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
			// &SimpleBounds{},
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}

	// fn := "/Users/jrush/Desktop/Animals.1/addresses.csv"
	// addrs := file.AsciiFileToLines(fn)
	// lens := []int{len(app.adverbs), len(app.adjectives), len(app.nouns), len(app.styles), len(app.colors), len(app.colors), len(app.colors), len(app.styles)}
	// fmt.Println(lens)
	// for _, addr := range addrs {
	// 	fmt.Println(app.GetPrompt(addr))
	// 	// _ = app.GetPrompt(addr)
	// }
}
