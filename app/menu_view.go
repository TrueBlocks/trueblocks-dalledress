package app

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Find: NewViews
func (a *App) ViewHome(cd *menu.CallbackData) {
	logger.Info("ViewHome")
	runtime.EventsEmit(a.ctx, "navigate", "/")
	a.SetLast("route", "/")
}

func (a *App) ViewDalle(cd *menu.CallbackData) {
	logger.Info("ViewDalle")
	runtime.EventsEmit(a.ctx, "navigate", "/dalle")
	a.SetLast("route", "/dalle")
}

func (a *App) ViewSeries(cd *menu.CallbackData) {
	logger.Info("ViewSeries")
	runtime.EventsEmit(a.ctx, "navigate", "/series")
	a.SetLast("route", "/series")
}

func (a *App) ViewHistory(cd *menu.CallbackData) {
	logger.Info("ViewHistory")
	last := a.GetLast("address")
	if len(last) == 0 {
		last = "trueblocks.eth"
	}
	runtime.EventsEmit(a.ctx, "navigate", "/history/"+last)
	a.SetLast("route", "/history/"+last)
}

func (a *App) ViewMonitors(cd *menu.CallbackData) {
	logger.Info("ViewMonitors")
	runtime.EventsEmit(a.ctx, "navigate", "/monitors")
	a.SetLast("route", "/monitors")
}

func (a *App) ViewNames(cd *menu.CallbackData) {
	logger.Info("ViewNames")
	runtime.EventsEmit(a.ctx, "navigate", "/names")
	a.SetLast("route", "/names")
}

func (a *App) ViewIndexes(cd *menu.CallbackData) {
	logger.Info("ViewIndexes")
	runtime.EventsEmit(a.ctx, "navigate", "/indexes")
	a.SetLast("route", "/indexes")
}

func (a *App) ViewManifest(cd *menu.CallbackData) {
	logger.Info("ViewManifest")
	runtime.EventsEmit(a.ctx, "navigate", "/manifest")
	a.SetLast("route", "/manifest")
}

func (a *App) ViewAbis(cd *menu.CallbackData) {
	logger.Info("ViewAbis")
	runtime.EventsEmit(a.ctx, "navigate", "/abis")
	a.SetLast("route", "/abis")
}

func (a *App) ViewStatus(cd *menu.CallbackData) {
	logger.Info("ViewStatus")
	runtime.EventsEmit(a.ctx, "navigate", "/status")
	a.SetLast("route", "/status")
}

func (a *App) ViewDaemons(cd *menu.CallbackData) {
	logger.Info("ViewDaemons")
	runtime.EventsEmit(a.ctx, "navigate", "/daemons")
	a.SetLast("route", "/daemons")
}

func (a *App) ViewSettings(cd *menu.CallbackData) {
	logger.Info("ViewSettings")
	runtime.EventsEmit(a.ctx, "navigate", "/settings")
	a.SetLast("route", "/settings")
}
