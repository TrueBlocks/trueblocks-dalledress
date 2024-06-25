package main

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
)

func (app *App) pipe0_handleAddrs(addresses []string) {
	for _, address := range addresses {
		if len(address) != 66 && !base.IsValidAddress(address) {
			continue
		}
		if dd, err := NewDalleDress(app.databases, address); err != nil {
			logger.Fatal(err.Error())
		} else {
			app.ReportOn("PostSelect", dd.Orig, "json", dd.String())
			app.pipe2Chan <- dd
		}
	}
	close(app.pipe2Chan)
}
