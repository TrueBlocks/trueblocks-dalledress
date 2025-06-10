package app

import (
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/enhancedcollection"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/preferences"
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

	// ADD_ROUTE
	app.names = names.NewNamesCollection()
	app.abis = enhancedcollection.NewEnhancedAbisCollection()
	// ADD_ROUTE

	_ = app.Reload()
	return app
}
