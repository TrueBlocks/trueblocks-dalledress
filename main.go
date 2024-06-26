package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
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
		ctx := context.Background()
		app.startup(ctx)
		app.domReady(ctx)
		app.handleLines()
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

func isContentPolicyViolation(err error) bool {
	var apiErr struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Param   string `json:"param"`
		Type    string `json:"type"`
	}
	if jsonErr := json.Unmarshal([]byte(err.Error()), &apiErr); jsonErr == nil {
		return apiErr.Code == "content_policy_violation"
	}
	return false
}

func (app *App) handleLines() {
// sourceFile := "addresses.txt"
// lines := file.AsciiFileToLines(sourceFile)
// if len(lines) > 0 {
// 	wg := sync.WaitGroup{}
// 	logger.Info("Starting at address ", app.Series.Last, " of ", len(lines))
// 	app.nMade = 0
// 	for i := 0; i < len(lines); i++ {
// 		if lines[i][0] == '#' || len(lines[i]) < 42 {
// 			continue
// 		}
// 		if i > int(app.Series.Last) {
// 			if address, ok := app.validateInput(lines[i]); !ok {
// 				fmt.Println("Invalid address", lines[i])
// 				return
// 			} else {
// 				wg.Add(1)
// 				go doOne(i, &wg, app, address.Hex())
// 				app.Series.Last = i
// 				app.Series.Save()
// 				if (i+1)%5 == 0 {
// 					wg.Wait()
// 					if app.nMade > 4 {
// 						logger.Info("Sleeping for 60 seconds")
// 						time.Sleep(time.Second * 60)
// 						app.nMade = 0
// 					}
// 				}
// 			}
// 		}
// 	}
// 	wg.Wait()
// 	return
// }
// if len(os.Args) == 2 {
// 	for i := 0; i < 10; i++ {
// 		wg := sync.WaitGroup{}
// 		for j := 0; j < 5; j++ {
// 			logger.Info("Round", i, "run", j)
// 			wg.Add(1)
// 			go doOne(i, &wg, app, fmt.Sprintf("0x%040x", 10010010+(i*10)+j)) // os.Args[1])
// 		}
// 		wg.Wait()
// 		logger.Info("Sleeping for 60 seconds")
// 		time.Sleep(time.Second * 60)
// 	}
// } else {
// 	wg := sync.WaitGroup{}
// 	for i, arg := range os.Args[1:] {
// 		wg.Add(1)
// 		go doOne(i, &wg, app, arg)
// 	}
// 	wg.Wait()
// }
// }
}
