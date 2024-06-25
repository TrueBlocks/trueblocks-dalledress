package main

import (
	"context"
	"embed"
	"encoding/json"
	"log"
	"os"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/joho/godotenv"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	if os.Getenv("TB_WAILS") == "true" {
		app := NewApp()
		opts := options.App{
			Title:            app.GetSettings().Title,
			Width:            app.GetSettings().Width,
			Height:           app.GetSettings().Height,
			BackgroundColour: nil,
			Bind:             []interface{}{app},
			StartHidden:      true,
			OnStartup:        app.startup,
			OnDomReady:       app.domReady,
			OnShutdown:       app.shutdown,
			AssetServer: &assetserver.Options{
				Assets: assets,
			},
		}
		if err := wails.Run(&opts); err != nil {
			println("Error:", err.Error())
		}
	} else {
		app := NewApp()
		addresses := file.AsciiFileToLines("./addresses.txt")
		go app.pipe0_handleAddrs(addresses)
		go app.pipe2_handlePrompt()
		app.pipe6_handleImage()
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	} else if a.apiKey = os.Getenv("OPENAI_API_KEY"); a.apiKey == "" {
		log.Fatal("No OPENAI_API_KEY key found")
	}
}

func (a *App) domReady(ctx context.Context) {
	if a.ctx != context.Background() {
		runtime.WindowSetPosition(a.ctx, a.GetSettings().X, a.GetSettings().Y)
		runtime.WindowShow(ctx)
	}
}

func (a *App) shutdown(ctx context.Context) {
	if a.ctx != context.Background() {
		settings.Width, settings.Height = runtime.WindowGetSize(ctx)
		settings.X, settings.Y = runtime.WindowGetPosition(ctx)
		a.SaveSettings()
	}
}

type Settings struct {
	Title  string `json:"title"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	loaded bool   `json:"-"`
}

var settings = Settings{
	Title:  "TrueBlocks Dapplet",
	Width:  800,
	Height: 600,
	X:      0,
	Y:      0,
	loaded: false,
}

func (a *App) GetSettings() *Settings {
	if !settings.loaded {
		contents := file.AsciiFileToString("./settings.json")
		if contents != "" {
			json.Unmarshal([]byte(contents), &settings)
		}
		settings.loaded = true
	}
	return &settings
}

func (a *App) SaveSettings() {
	bytes, _ := json.Marshal(settings)
	file.StringToAsciiFile("./settings.json", string(bytes))
}
