package app

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/messages"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) SystemAbout(cd *menu.CallbackData) {
	logger.Info("System About")
}

func (a *App) SystemQuit(cd *menu.CallbackData) {
	if a.CurrentDoc.Dirty {
		messages.Send(a.ctx, messages.Error, messages.NewErrorMsg(fmt.Errorf("you have unsaved changes, save before quitting")))
		a.CurrentDoc.Dirty = false
		return
	}
	logger.Info("System Quit")
	runtime.Quit(a.ctx)
}
