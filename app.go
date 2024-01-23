package main

import (
	"context"
	"encoding/json"
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
	if a.conn = rpc.NewConnection("mainnet", true, map[string]bool{
		"blocks":       true,
		"receipts":     true,
		"transactions": true,
		"traces":       true,
		"logs":         true,
		"statements":   true,
		"state":        true,
		"tokens":       true,
		"results":      true,
	}); a.conn == nil {
		logger.Error("Could not find rpc server.")
	}

	var err error
	if a.adverbs, err = toLines("/Users/jrush/Desktop/Animals.1/adverbs.csv"); err != nil {
		logger.Fatal(err)
	}
	if a.adjectives, err = toLines("/Users/jrush/Desktop/Animals.1/adjectives.csv"); err != nil {
		logger.Fatal(err)
	}
	if a.nouns, err = toLines("/Users/jrush/Desktop/Animals.1/nouns.csv"); err != nil {
		logger.Fatal(err)
	}
	if a.colors, err = toLines("/Users/jrush/Desktop/Animals.1/colors.csv"); err != nil {
		logger.Fatal(err)
	}
	if a.styles, err = toLines("/Users/jrush/Desktop/Animals.1/styles.csv"); err != nil {
		logger.Fatal(err)
	}
	x := make([]string, 0, len(a.styles))
	for _, s := range a.styles {
		if !strings.Contains(s, "sensitive") { // remove culturally sensitive styles
			x = append(x, s)
		}
	}
	a.styles = x

	pT := "I NEED to test how the tool works with extremely simple prompts. DO NOT add any detail, just use it AS-IS: " + promptTemplate
	if a.pTemplate, err = template.New("prompt").Parse(pT); err != nil {
		logger.Fatal("could not create prompt template:", err)
	}
	dT := "I NEED to test how the tool works with extremely simple prompts. DO NOT add any detail, just use it AS-IS: " + dataTemplate
	if a.dTemplate, err = template.New("data").Parse(dT); err != nil {
		logger.Fatal("could not create data template:", err)
	}

	if err = godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	} else if a.apiKey = os.Getenv("OPENAI_API_KEY"); a.apiKey == "" {
		log.Fatal("No API key found in .env")
	}
}

func (a *App) domReady(ctx context.Context) {
	runtime.WindowSetPosition(a.ctx, a.GetSettings().Width, a.GetSettings().Height)
	runtime.WindowShow(ctx)
}

func (a *App) shutdown(ctx context.Context) {
	//a.config.Load()
	//a.config.WindowState.Width, a.config.WindowState.Height = runtime.WindowGetSize(ctx)
	//a.config.WindowState.X, a.config.WindowState.Y = runtime.WindowGetPosition(ctx)
	//a.config.Save()
}

func (a *App) GetTypes() []interface{} {
	ret := make([]interface{}, 0)
	ret = append(ret, a)
	return ret
}

type Settings struct {
	Title  string `json:"title"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	loaded bool   `json:"-"`
}

var settings = Settings{
	Title:  "TrueBlocks Account Browser",
	Width:  800,
	Height: 600,
	loaded: false,
}

func (a *App) GetSettings() Settings {
	if !settings.loaded {
		contents := file.AsciiFileToString("./settings.json")
		if contents != "" {
			json.Unmarshal([]byte(contents), &settings)
		}
		settings.loaded = true
	}
	return settings
}

func (a *App) SaveSettings() {
	bytes, _ := json.Marshal(settings)
	file.StringToAsciiFile("./settings.json", string(bytes))
}
