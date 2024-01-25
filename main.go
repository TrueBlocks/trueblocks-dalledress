package main

import (
	"embed"

	"github.com/TrueBlocks/trueblocks-browse/app"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := app.NewApp()
	opts := options.App{
		Title:            "trueblocks-browse",
		Width:            1024,
		Height:           768,
		BackgroundColour: &options.RGBA{R: 33, G: 37, B: 41, A: 1},
		OnStartup:        app.Startup,
		OnDomReady:       app.DomReady,
		OnShutdown:       app.Shutdown,
		Bind: []interface{}{
			app,
		},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
	}

	if err := wails.Run(&opts); err != nil {
		println("Error:", err.Error())
	}
}
