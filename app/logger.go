package app

import (
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/logging"
)

func (a *App) LogBackend(msg string) {
	logging.LogBackend(msg)
}

func (a *App) LogFrontend(msg string) {
	logging.LogFrontend(msg)
}
