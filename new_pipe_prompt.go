package main

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
)

func (app *App2) pipe2_handlePrompt() {
	var err error
	for dd := range app.pipe2Chan {
		if dd.Prompt, err = dd.generatePrompt(app.pTemplate, nil); err != nil {
			logger.Fatal(err.Error())
		}
		app.ReportOn("Prompt", dd.Orig, dd.Prompt)
		if dd.DataPrompt, err = dd.generatePrompt(app.dTemplate, nil); err != nil {
			logger.Fatal(err.Error())
		}
		app.ReportOn("DataPrompt", dd.Orig, dd.DataPrompt)
		if dd.TersePrompt, err = dd.generatePrompt(app.tTemplate, nil); err != nil {
			logger.Fatal(err.Error())
		}
		app.pipe6Chan <- dd
	}
	close(app.pipe6Chan)
}
