package app

import (
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/preferences"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
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
	app.names = types.NewNamesCollection(app)
	app.abis = types.NewAbisCollection(app)
	_ = app.Reload()
	return app
}
