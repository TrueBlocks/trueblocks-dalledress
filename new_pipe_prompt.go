package main

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/openai"
)

func (app *App2) pipe2_handlePrompt() {
	var err error
	for dd := range app.pipe2Chan {
		if dd.Prompt, err = dd.generatePrompt(app.pTemplate, nil); err != nil {
			logger.Fatal(err.Error())
		}
		app.ReportOn("Prompt", dd.Orig, "txt", dd.Prompt)
		if dd.DataPrompt, err = dd.generatePrompt(app.dTemplate, nil); err != nil {
			logger.Fatal(err.Error())
		}
		app.ReportOn("DataPrompt", dd.Orig, "txt", dd.DataPrompt)
		if dd.TersePrompt, err = dd.generatePrompt(app.tTemplate, nil); err != nil {
			logger.Fatal(err.Error())
		}
		app.ReportOn("TersePrompt", dd.Orig, "txt", dd.TersePrompt)
		if dd.EnhancedPrompt, err = openai.EnhancePrompt(dd.Prompt); err != nil {
			logger.Fatal(err.Error())
		}
		dd.EnhancedPrompt += " DO NOT PUT ANY TEXT IN THE IMAGE. NONE. ZERO. NADA."
		app.ReportOn("EnhancedPrompt", dd.Orig, "txt", dd.EnhancedPrompt)
		app.pipe6Chan <- dd
	}
	close(app.pipe6Chan)
}
