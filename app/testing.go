package app

import (
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/preferences"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/abis"
)

func setupTestApp() *App {
	app := &App{
		Preferences: &preferences.Preferences{
			Org:  preferences.OrgPreferences{},
			User: preferences.UserPreferences{},
			App:  preferences.AppPreferences{},
		},
	}

	// ADD_ROUTE
	app.names = types.NewNamesCollection()
	app.abis = abis.NewAbisCollection()
	// ADD_ROUTE

	_ = app.Reload()
	return app
}
