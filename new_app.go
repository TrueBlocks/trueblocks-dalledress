package main

import (
	"path/filepath"
	"strings"
	"text/template"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
)

type App2 struct {
	pipe1Chan chan *DalleDress
	pipe2Chan chan *DalleDress
	pipe6Chan chan *DalleDress
	databases map[string][]string
	pTemplate *template.Template `json:"-"`
	dTemplate *template.Template `json:"-"`
	tTemplate *template.Template `json:"-"`
}

func NewApp2() *App2 {
	app := App2{
		pipe1Chan: make(chan *DalleDress),
		pipe2Chan: make(chan *DalleDress),
		pipe6Chan: make(chan *DalleDress),
		databases: make(map[string][]string),
	}

	for _, db := range databaseNames {
		if app.databases[db] == nil {
			app.databases[db] = file.AsciiFileToLines("./databases/" + db + ".csv")[1:]
			for i := 0; i < len(app.databases[db]); i++ {
				app.databases[db][i] = strings.Replace(app.databases[db][i], "v0.1.0,", "", -1)
			}
		}
	}

	var err error
	if app.pTemplate, err = template.New("prompt").Parse(promptTemplate); err != nil {
		logger.Fatal("could not create prompt template:", err)
	}
	if app.dTemplate, err = template.New("data").Parse(dataTemplate); err != nil {
		logger.Fatal("could not create data template:", err)
	}
	if app.tTemplate, err = template.New("terse").Parse(terseTemplate); err != nil {
		logger.Fatal("could not create terse template:", err)
	}

	return &app
}

func (app *App2) ReportOn(loc, address, value string) {
	path := filepath.Join("./output/", strings.ToLower(loc))
	file.EstablishFolder(path)
	file.StringToAsciiFile(filepath.Join(path, address+".json"), value)
}
