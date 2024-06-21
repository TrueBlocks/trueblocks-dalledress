package main

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
)

func (app *App2) pipe0_handleAddrs(addresses []string) {
	for i, address := range addresses {
		if !base.IsValidAddress(address) {
			continue
		}
		dalleDress := NewDalleDress(i, address)
		if len(dalleDress.Seed) > 66 {
			dalleDress.Seed = dalleDress.Seed[2:66]
			app.pipe1Chan <- dalleDress
		} else {
			logger.Fatal("Seed length is less than 66", dalleDress.Seed)
		}
	}
	close(app.pipe1Chan)
}
