package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/joho/godotenv"
)

type App struct {
	ctx       context.Context
	apiKey    string
	databases map[string][]string
	pipe2Chan chan *DalleDress
	pipe6Chan chan *DalleDress
	pTemplate *template.Template
	dTemplate *template.Template
	tTemplate *template.Template
	Series    Series `json:"series"`
}

type Adjustment int

const (
	NoAdjustment Adjustment = 0
	Filtered                = 1 << iota
	Noneable
)

func NewApp() *App {
	app := App{
		pipe2Chan: make(chan *DalleDress),
		pipe6Chan: make(chan *DalleDress),
		databases: make(map[string][]string),
	}

	var err error
	if app.Series, err = GetSeries(); err != nil {
		logger.Fatal(err)
	}

	for _, db := range databaseNames {
		if app.databases[db] == nil {
			aspect := Adjustment(Filtered)
			if db == "litstyles" || db == "occupations" || db == "colors" {
				aspect |= Noneable
			}
			if lines, err := app.toLines(db, aspect); err != nil {
				logger.Fatal(err)
			} else {
				app.databases[db] = lines
				for i := 0; i < len(app.databases[db]); i++ {
					app.databases[db][i] = strings.Replace(app.databases[db][i], "v0.1.0,", "", -1)
				}
			}
		}
	}

	if app.pTemplate, err = template.New("prompt").Parse(promptTemplate); err != nil {
		logger.Fatal("could not create prompt template:", err)
	}
	if app.dTemplate, err = template.New("data").Parse(dataTemplate); err != nil {
		logger.Fatal("could not create data template:", err)
	}
	if app.tTemplate, err = template.New("terse").Parse(terseTemplate); err != nil {
		logger.Fatal("could not create terse template:", err)
	}

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	} else if app.apiKey = os.Getenv("OPENAI_API_KEY"); app.apiKey == "" {
		log.Fatal("No OPENAI_API_KEY key found")
	}

	return &app
}

func (a App) String() string {
	bytes, _ := json.MarshalIndent(a, "", "  ")
	return string(bytes)
}

func (app *App) ReportOn(loc, address, ft, value string) {
	path := filepath.Join("./output/", app.Series.Suffix, strings.ToLower(loc))
	file.EstablishFolder(path)
	file.StringToAsciiFile(filepath.Join(path, address+"."+ft), value)
}

func (app *App) ReportDone(address string) {
	logger.Info("Done", address)
}

func (a *App) toLines(db string, adjust Adjustment) ([]string, error) {
	filename := "./databases/" + db + ".csv"
	lines := file.AsciiFileToLines(filename)
	var err error
	if len(lines) == 0 {
		err = fmt.Errorf("could not load %s", filename)
	} else {
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
