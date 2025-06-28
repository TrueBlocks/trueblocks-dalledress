package app

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/monitors"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/names"
)

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

func (a *App) NameFromAddress(address string) (*names.Name, bool) {
	collection := names.GetNamesCollection(&types.Payload{})
	return collection.NameFromAddress(base.HexToAddress(address))
}

func (a *App) MonitorsClean(payload *types.Payload, addresses []string) error {
	collection := monitors.GetMonitorsCollection(payload)
	return collection.Clean(addresses)
}
