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
	if a.emotions2, err = toLines("/Users/jrush/Desktop/Animals.1/emotions.csv"); err != nil {
		logger.Fatal(err)
	}
	for _, e := range a.emotions2 {
		parts := strings.Split(e, ",")
		if len(parts) > 0 {
			a.emotions1 = append(a.emoptions1, parts[0])
		}
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

func (a *App) GetTypes() []interface{} {
	ret := make([]interface{}, 0)
	ret = append(ret, a)
	return ret
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

func (a *App) domReady(ctx context.Context) {
	runtime.WindowSetPosition(a.ctx, a.GetSettings().X, a.GetSettings().Y)
	runtime.WindowShow(ctx)
}

func (a *App) shutdown(ctx context.Context) {
	settings.Width, settings.Height = runtime.WindowGetSize(ctx)
	settings.X, settings.Y = runtime.WindowGetPosition(ctx)
	SaveSettings()
}

func (a *App) SaveSettings() {
	bytes, _ := json.Marshal(settings)
	file.StringToAsciiFile("./settings.json", string(bytes))
}
