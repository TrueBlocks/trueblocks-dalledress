package app

import "github.com/TrueBlocks/trueblocks-dalledress/pkg/types"

func (a *App) Reload(payload *types.Payload) error {
	lastView := a.GetAppPreferences().LastView

	// ADD_ROUTE
	switch lastView {
	case "/names":
		return a.ReloadNames(payload)
	case "/abis":
		return a.ReloadAbis(payload)
	case "/exports":
		return a.ReloadExports(payload)
	case "/monitors":
		return a.ReloadMonitors(payload)
	case "/chunks":
		return a.ReloadChunks(payload)
	}
	// ADD_ROUTE

	return nil
}
