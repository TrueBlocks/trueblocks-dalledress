package main

import (
	"context"
	"embed"
	"os"
	"sync"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()
	if len(os.Args) > 0 {
		ctx := context.Background()
		app.startup(ctx)
		app.domReady(ctx)
		wg := sync.WaitGroup{}
		for _, arg := range os.Args[1:] {
			wg.Add(1)
			go doOne(&wg, app, arg)
		}
		wg.Wait()
	} else {
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
}
