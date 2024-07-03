package main

import (
	"context"
	"embed"
	"os"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-dalledress/app"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	if os.Getenv("TB_WAILS") == "true" {
		a := app.NewApp()
		opts := options.App{
			Title:            a.GetSession().Title,
			Width:            a.GetSession().Width,
			Height:           a.GetSession().Height,
			OnStartup:        a.Startup,
			OnDomReady:       a.DomReady,
			OnShutdown:       a.Shutdown,
			BackgroundColour: nil,
			Bind: []interface{}{
				a,
			},
			StartHidden: true,
			AssetServer: &assetserver.Options{
				Assets: assets,
			},
		}
		if err := wails.Run(&opts); err != nil {
			println("Error:", err.Error())
		}
	} else {
		app := app.NewApp()
		ctx := context.Background()
		app.Startup(ctx)
		app.DomReady(ctx)
		logger.Info("Running in console mode")
		app.HandleLines()
	}
}
