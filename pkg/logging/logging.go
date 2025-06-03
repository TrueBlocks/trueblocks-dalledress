package logging

import (
	"log"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
)

func LogBackend(msg string) {
	log.Println(colors.BrightBlue+"BACKEND", msg, colors.Off)
}

func LogFrontend(msg string) {
	log.Println(colors.Green+"FRONTEND", msg, colors.Off)
}
