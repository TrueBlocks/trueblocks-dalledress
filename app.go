package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/rpc"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/walk"
	"github.com/joho/godotenv"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx  context.Context `json:"-"`
	conn *rpc.Connection `json:"-"`
	// addresses   []string           `json:"-"`
	Adverbs     []string `json:"adverbs"`
	Adjectives  []string `json:"adjectives"`
	Nouns       []string `json:"nouns"`
	Emotions    []string `json:"emotions"`
	Occupations []string `json:"occupations"`
	Gerunds     []string `json:"gerunds"`
	Litstyles   []string `json:"litstyles"`
	Artstyles   []string `json:"artstyles"`
	Colors      []string `json:"colors"`
	// p Template   *template.Template `json:"-"`
	// d Template   *template.Template `json:"-"`
	// t Template *template.Template `json:"-"`
	apiKey string `json:"-"`
	// nMade       int                `json:"-"`
	Series Series `json:"series"`
}

func (a App) String() string {
	bytes, _ := json.MarshalIndent(a, "", "  ")
	return string(bytes)
}

// func doOne(which int, wg *sync.WaitGroup, app *App, arg string) {
// 	defer wg.Done()
// 	app.GetImage(which, arg)
// }

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	if a.conn = rpc.NewConnection("mainnet", true, map[walk.CacheType]bool{
		walk.Cache_Blocks:       true,
		walk.Cache_Receipts:     true,
		walk.Cache_Transactions: true,
		walk.Cache_Traces:       true,
		walk.Cache_Logs:         true,
		walk.Cache_Statements:   true,
		walk.Cache_State:        true,
		walk.Cache_Tokens:       true,
		walk.Cache_Results:      true,
	}); a.conn == nil {
		logger.Error("Could not find rpc server.")
	}

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	} else if a.apiKey = os.Getenv("OPENAI_API_KEY"); a.apiKey == "" {
		log.Fatal("No OPENAI_API_KEY key found")
	}

	var err error
	// Make the templates
	// if a.p Template, err = template.New("prompt").Parse(prompt Template); err != nil {
	// 	logger.Fatal("could not create prompt template:", err)
	// }
	// if a.d Template, err = template.New("data").Parse(dataTemplate); err != nil {
	// 	logger.Fatal("could not create data template:", err)
	// }
	// if a.t Template, err = template.New("terse").Parse(terseTemplate); err != nil {
	// 	logger.Fatal("could not create terse template:", err)
	// }

	if a.Series, err = GetSeries(); err != nil {
		logger.Fatal(err)
	}

	// Read (and possibly filter) the databases
	if a.Adverbs, err = a.toLines("adverbs", Filtered); err != nil {
		logger.Fatal(err)
	}

	if a.Adjectives, err = a.toLines("adjectives", Filtered); err != nil {
		logger.Fatal(err)
	}

	if a.Nouns, err = a.toLines("nouns", Filtered); err != nil {
		logger.Fatal(err)
	}

	if a.Emotions, err = a.toLines("emotions", CommaToParens|Filtered); err != nil {
		logger.Fatal(err)
	}

	if a.Occupations, err = a.toLines("occupations", CommaToParens|Filtered|Noneable); err != nil {
		logger.Fatal(err)
	}

	if a.Gerunds, err = a.toLines("gerunds", Filtered); err != nil {
		logger.Fatal(err)
	}

	if a.Artstyles, err = a.toLines("artstyles", Filtered); err != nil {
		logger.Fatal(err)
	}

	if a.Litstyles, err = a.toLines("litstyles", CommaToParens|Filtered); err != nil {
		logger.Fatal(err)
	}

	if a.Colors, err = a.toLines("colors", Filtered|Noneable); err != nil {
		logger.Fatal(err)
	}
}

type Adjustment int

const (
	NoAdjustment  Adjustment = 0
	CommaToParens            = 1 << iota
	Filtered
	Noneable
)

type Series struct {
	Last        int      `json:"last"`
	Suffix      string   `json:"suffix"`
	Artstyles   []string `json:"artstyles"`
	Emotions    []string `json:"emotions"`
	Occupations []string `json:"occupations"`
	Nouns       []string `json:"nouns"`
	Litstyles   []string `json:"litstyles"`
	Adverbs     []string `json:"adverbs"`
	Adjectives  []string `json:"adjectives"`
	Gerunds     []string `json:"gerunds"`
	Colors      []string `json:"colors"`
}

func (s *Series) String() string {
	bytes, _ := json.MarshalIndent(s, "", "  ")
	return string(bytes)
}

func (s *Series) Save() {
	file.StringToAsciiFile("series.json", s.String())
}

func GetSeries() (Series, error) {
	if str := file.AsciiFileToString("series.json"); len(str) == 0 {
		return Series{}, fmt.Errorf("could not load series.json")
	} else {
		bytes := []byte(str)
		var s Series
		if err := json.Unmarshal(bytes, &s); err != nil {
			logger.Error("could not unmarshal series:", err)
			return Series{}, err
		}
		return s, nil
	}
}

func (s *Series) GetFilter(fieldName string) ([]string, error) {
	reflectedT := reflect.ValueOf(s)
	field := reflect.Indirect(reflectedT).FieldByName(fieldName)
	if !field.IsValid() {
		return nil, fmt.Errorf("field %s not valid", fieldName)
	}
	if field.Kind() != reflect.Slice {
		return nil, fmt.Errorf("field %s not a slice", fieldName)
	}
	if field.Type().Elem().Kind() != reflect.String {
		return nil, fmt.Errorf("field %s not a string slice", fieldName)
	}
	return field.Interface().([]string), nil
}

func (a *App) toLines(db string, adjust Adjustment) ([]string, error) {
	dbFolder, _ := GetConfigDir("TrueBlocks/dalledress/databases")
	filename := filepath.Join(dbFolder, db+".csv")

	lines := file.AsciiFileToLines(filename)
	var err error
	if len(lines) == 0 {
		err = fmt.Errorf("could not load %s", filename)
	} else {
		if adjust&CommaToParens != 0 {
			for i, line := range lines {
				lines[i] = strings.Replace(line, ",", " (", 1) + ")"
			}
		}
		if adjust&Filtered != 0 {
			fn := strings.ToUpper(db[:1]) + db[1:]
			if filter, err := a.Series.GetFilter(fn); err != nil {
				return lines, err

			} else {
				if len(filter) == 0 {
					// logger.Info("says filtered, but no filter for", db)
					return lines, nil
				}
				// logger.Info("found filter for", db, filter)

				filtered := make([]string, 0, len(lines))
				for _, line := range lines {
					for _, f := range filter {
						// logger.Info(line, f, strings.Contains(line, f))
						if strings.Contains(line, f) {
							filtered = append(filtered, line)
						}
					}
				}
				lines = filtered
			}
		}
	}

	if len(lines) == 0 {
		if adjust&Noneable != 0 {
			lines = append(lines, "none")
		} else {
			return lines, fmt.Errorf("no lines match query in %s", filename)
		}
	}

	return lines, err
}

// GetConfigDir returns the operating system's configuration folder for the current user.
// If the folder does not exist, it is created.
func GetConfigDir(appDir string) (string, error) {
	if configPath, err := os.UserConfigDir(); err != nil {
		return "", err
	} else {
		path := filepath.Join(configPath, appDir)
		file.EstablishFolder(path)
		return path, nil
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
