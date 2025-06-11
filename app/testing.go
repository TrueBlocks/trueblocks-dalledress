package app

import (
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/preferences"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/abis"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/names"
)

func setupTestApp() *App {
	app := &App{
		Preferences: &preferences.Preferences{
			Org:  preferences.OrgPreferences{},
			User: preferences.UserPreferences{},
			App:  preferences.AppPreferences{},
		},
	}

	// NAMES_ROUTE
	app.names = names.NewNamesCollection()
	// NAMES_ROUTE
	// ABIS_ROUTE
	app.abis = abis.NewAbisCollection()
	// ABIS_ROUTE

	_ = app.Reload()
	return app
}
