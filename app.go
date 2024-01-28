package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/rpc"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/paths"
	"github.com/joho/godotenv"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx           context.Context
	conn          *rpc.Connection
	addresses     []string
	adverbs       []string
	adjectives    []string
	emotionsShort []string
	emotions      []string
	literary      []string
	nouns         []string
	styles        []string
	pTemplate     *template.Template
	dTemplate     *template.Template
	tTemplate     *template.Template
	apiKey        string
}

func doOne(wg *sync.WaitGroup, app *App, arg string) {
	defer wg.Done()
	app.GetImage(arg, false)
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

	dbFolder, _ := paths.GetConfigDir("TrueBlocks/dalledress/databases")

	var err error
	if a.adverbs, err = toLines(filepath.Join(dbFolder, "adverbs.csv")); err != nil {
		logger.Fatal(err)
	}
	if a.adjectives, err = toLines(filepath.Join(dbFolder, "adjectives.csv")); err != nil {
		logger.Fatal(err)
	}
	if a.emotions, err = toLines(filepath.Join(dbFolder, "emotions.csv")); err != nil {
		logger.Fatal(err)
	}
	for i := 0; i < len(a.emotions); i++ {
		e2 := a.emotions[i]
		parts := strings.Split(e2, ",")
		if len(parts) > 2 {
			a.emotions[i] = parts[0] + " (" + strings.Replace(parts[2], ".", "", -1) + ")"
			a.emotionsShort = append(a.emotionsShort, parts[0])
		}
	}
	if lines, err := toLines(filepath.Join(dbFolder, "literarystyles.csv")); err != nil {
		logger.Fatal(err)
	} else {
		for _, line := range lines {
			parts := strings.Split(line, ",")
			if len(parts) > 1 {
				a.literary = append(a.literary, parts[0]+" ("+strings.Replace(parts[1], ".", "", -1)+")")
			}
		}
	}
	if a.nouns, err = toLines(filepath.Join(dbFolder, "nouns.csv")); err != nil {
		logger.Fatal(err)
	}
	if a.styles, err = toLines(filepath.Join(dbFolder, "styles.csv")); err != nil {
		logger.Fatal(err)
	}
	for i := 0; i < len(a.styles); i++ {
		a.styles[i] = strings.Replace(a.styles[i], ",", " from the ", -1)
	}

	if a.pTemplate, err = template.New("prompt").Parse(promptTemplate); err != nil {
		logger.Fatal("could not create prompt template:", err)
	}
	if a.dTemplate, err = template.New("data").Parse(dataTemplate); err != nil {
		logger.Fatal("could not create data template:", err)
	}
	if a.tTemplate, err = template.New("terse").Parse(terseTemplate); err != nil {
		logger.Fatal("could not create terse template:", err)
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

func (a *App) SaveSettings() {
	bytes, _ := json.Marshal(settings)
	file.StringToAsciiFile("./settings.json", string(bytes))
}
