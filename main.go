package main

import (
	"context"
	"embed"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()
	logger.Info("nArgs: ", len(os.Args))
	if len(os.Args) > 0 {
		ctx := context.Background()
		app.startup(ctx)
		app.domReady(ctx)
		if len(os.Args) == 2 {
			for i := 0; i < 10; i++ {
				wg := sync.WaitGroup{}
				for j := 0; j < 5; j++ {
					logger.Info("Round", i, "run", j)
					wg.Add(1)
					go doOne(&wg, app, fmt.Sprintf("0x%040x", 10010010+(i*10)+j)) // os.Args[1])
					// SeedBump++
				}
				wg.Wait()
				logger.Info("Sleeping for 60 seconds")
				time.Sleep(time.Second * 60)
			}
		} else {
			wg := sync.WaitGroup{}
			for _, arg := range os.Args[1:] {
				wg.Add(1)
				go doOne(&wg, app, arg)
			}
			wg.Wait()
		}
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
