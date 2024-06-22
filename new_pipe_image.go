package main

import "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"

func (app *App2) pipe6_handleImage() {
	for dd := range app.pipe6Chan {
		app.ReportDone(dd.Orig)
		if err := dd.getAiImage(); err != nil {
			logger.Error(err.Error())
		}
	}
}
