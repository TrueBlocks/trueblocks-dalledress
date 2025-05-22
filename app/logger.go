package app

import (
	"log"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
)

func (a *App) LogBackend(msg string) {
	log.Println(colors.BrightBlue+"BACKEND", msg, colors.Off)
}

func (a *App) LogFrontend(msg string) {
	log.Println(colors.Green+"FRONTEND", msg, colors.Off)
}
