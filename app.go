package main

import (
	"context"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/rpc"
	"github.com/joho/godotenv"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx        context.Context
	conn       *rpc.Connection
	addresses  []string
	adverbs    []string
	adjectives []string
	emotions1  []string
	emotions2  []string
	nouns      []string
	colors     []string
	styles     []string
	pTemplate  *template.Template
	dTemplate  *template.Template
	apiKey     string
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.conn = rpc.NewConnection("mainnet", false, map[string]bool{})
	a.adverbs = file.AsciiFileToLines("/Users/jrush/Desktop/Animals.1/adverbs.csv")
	a.adjectives = file.AsciiFileToLines("/Users/jrush/Desktop/Animals.1/adjectives.csv")
	a.emotions1 = file.AsciiFileToLines("/Users/jrush/Desktop/Animals.1/emotions.csv")
x	a.emotions2 = a.emotions1
	a.nouns = file.AsciiFileToLines("/Users/jrush/Desktop/Animals.1/nouns.csv")
	a.colors = file.AsciiFileToLines("/Users/jrush/Desktop/Animals.1/colors.csv")
	a.styles = file.AsciiFileToLines("/Users/jrush/Desktop/Animals.1/styles.csv")
	x := make([]string, 0, len(a.styles))
	for _, s := range a.styles {
		if !strings.Contains(s, "sensitive") { // remove culturally sensitive styles
			x = append(x, s)
		}
	}
	a.styles = x
	var err error
	if a.pTemplate, err = template.New("prompt").Parse(promptTemplate); err != nil {
		logger.Fatal("could not create prompt template:", err)
	}
	if a.dTemplate, err = template.New("data").Parse(dataTemplate); err != nil {
		logger.Fatal("could not create data template:", err)
	}
	if err = godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	} else if a.apiKey = os.Getenv("OPENAI_API_KEY"); a.apiKey == "" {
		log.Fatal("No API key found in .env")
	}
}

type Settings struct {
	Title  string `json:"title"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
}

var settings = Settings{
	Title:  "Chifra",
	Width:  800,
	Height: 600,
	X:      0,
	Y:      0,
}

func (a *App) GetSettings() *Settings {
	return &settings
}

func (a *App) GetTypes() []interface{} {
	return []interface{}{
		a,
	}
}

func (a *App) domReady(ctx context.Context) {
	runtime.WindowSetPosition(a.ctx, settings.X, settings.Y)
	runtime.WindowShow(ctx)
}

func (a *App) shutdown(ctx context.Context) {
	settings.Width, settings.Height = runtime.WindowGetSize(ctx)
	settings.X, settings.Y = runtime.WindowGetPosition(ctx)
}
