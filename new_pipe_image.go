package main

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/openai"
)

func (app *App) pipe6_handleImage() {
	for dd := range app.pipe6Chan {
		app.ReportDone(dd.Orig)
		imageData := openai.ImageData{
			EnhancedPrompt: dd.EnhancedPrompt,
			TersePrompt:    dd.TersePrompt,
			Hash:           dd.Orig,
			SeriesName:     app.Series.Suffix,
		}
		if err := openai.GetImage(&imageData); err != nil {
			logger.Error(err.Error())
		}
	}
}
