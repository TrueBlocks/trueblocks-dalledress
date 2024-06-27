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
	pTemplate *template.Template
	dTemplate *template.Template
	tTemplate *template.Template
	aTemplate *template.Template
	Series    Series `json:"series"`
}

func NewApp() *App {
	app := App{
		databases: make(map[string][]string),
	}

	var err error
	if app.Series, err = GetSeries(); err != nil {
		logger.Fatal(err)
	}

	for _, db := range databaseNames {
		if app.databases[db] == nil {
			if lines, err := app.toLines(db); err != nil {
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
	if app.aTemplate, err = template.New("author").Parse(authorTemplate); err != nil {
		logger.Fatal("could not create prompt template:", err)
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

func ReportOn(loc, address, ft, value string) {
	path := filepath.Join("./output/", strings.ToLower(loc))
	file.EstablishFolder(path)
	file.StringToAsciiFile(filepath.Join(path, address+"."+ft), value)
}

func (a *App) toLines(db string) ([]string, error) {
	filename := "./databases/" + db + ".csv"
	lines := file.AsciiFileToLines(filename)
	var err error
	if len(lines) == 0 {
		err = fmt.Errorf("could not load %s", filename)
	} else {
		fn := strings.ToUpper(db[:1]) + db[1:]
		if filter, err := a.Series.GetFilter(fn); err != nil {
			return lines, err

		} else {
			if len(filter) == 0 {
				return lines, nil
			}

			filtered := make([]string, 0, len(lines))
			for _, line := range lines {
				for _, f := range filter {
					if strings.Contains(line, f) {
						filtered = append(filtered, line)
					}
				}
			}
			lines = filtered
		}
	}

	if len(lines) == 0 {
		lines = append(lines, "none")
	}

	return lines, err
}
