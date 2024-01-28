package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main_wails() {
	app := NewApp()
	opts := options.App{
		Title:            app.GetSettings().Title,
		Width:            app.GetSettings().Width,
		Height:           app.GetSettings().Height,
		BackgroundColour: nil,
		Bind:             app.GetTypes(),
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
}
