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
	batchSize := 5
	rateLimit := time.Second / 5
	sem := make(chan struct{}, batchSize)
	lines := func() []string {
		if len(os.Args) < 2 {
			return file.AsciiFileToLines("addresses.txt")
		}
		return os.Args[1:]
	}()

	// logger.Info("Starting at address ", app.Series.Last, " of ", len(lines))

	var wg sync.WaitGroup
	ticker := time.NewTicker(rateLimit)
	defer ticker.Stop()

	for i, addr := range lines {
		if app.Series.Last > 0 && i <= int(app.Series.Last) {
			continue
		}

		sem <- struct{}{}
		wg.Add(1)
		go func(index int, address string) {
			defer wg.Done()
			defer func() { <-sem }()
			backoff := time.Second
			maxRetries := 5
			for attempt := 0; attempt < maxRetries; attempt++ {
				<-ticker.C
				err := app.GetImage(address)
				if err == nil {
					return
				}
				if isContentPolicyViolation(err) {
					msg := fmt.Sprintf("Content policy violation, skipping retry for address: %s Error: %s", address, err)
					logger.Error(msg)
					return
				}
				msg := fmt.Sprintf("Error fetching image: %s Retry attempt: %d Sleeping: %d", err, attempt+1, backoff)
				logger.Error(msg)
				time.Sleep(backoff)
				backoff = time.Duration(float64(backoff) * (1 + rand.Float64()))
			}
			logger.Error("Failed to fetch image after max retries:", address)
		}(i, addr)

		app.Series.Last = i
		app.Series.SaveSeries("series.json", app.Series.Last)

		if (i+1)%batchSize == 0 {
			wg.Wait()
		}
	}
	wg.Wait()
}
