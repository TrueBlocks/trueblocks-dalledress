package app

import (
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
)

// NewViews
func (a *App) GetMenus() *menu.Menu {
	appMenu := menu.NewMenu()

	systemMenu := appMenu.AddSubmenu("DalleDress")
	systemMenu.AddText("About", nil, a.SystemAbout)
	systemMenu.AddText("Quit", keys.CmdOrCtrl("q"), a.SystemQuit)

	fileMenu := appMenu.AddSubmenu("File")
	fileMenu.AddText("New File", keys.CmdOrCtrl("n"), a.FileNew)
	fileMenu.AddText("Open File", keys.CmdOrCtrl("o"), a.FileOpen)
	fileMenu.AddText("Save File", keys.CmdOrCtrl("s"), a.FileSave)
	fileMenu.AddText("Save As File", keys.CmdOrCtrl("a"), a.FileSaveAs)

	viewMenu := appMenu.AddSubmenu("View")
	viewMenu.AddText("Home", keys.CmdOrCtrl("1"), a.ViewHome)
	viewMenu.AddText("Dalle", keys.CmdOrCtrl("2"), a.ViewDalle)
	viewMenu.AddText("Series", keys.CmdOrCtrl("3"), a.ViewSeries)
	viewMenu.AddText("History", keys.CmdOrCtrl("4"), a.ViewHistory)
	viewMenu.AddText("Monitors", keys.CmdOrCtrl("5"), a.ViewMonitors)
	viewMenu.AddText("Names", keys.CmdOrCtrl("6"), a.ViewNames)
	viewMenu.AddText("Indexes", keys.CmdOrCtrl("7"), a.ViewIndexes)
	viewMenu.AddText("Manifest", keys.CmdOrCtrl("8"), a.ViewManifest)
	viewMenu.AddText("Abis", keys.CmdOrCtrl("9"), a.ViewAbis)
	viewMenu.AddText("Status", keys.CmdOrCtrl("0"), a.ViewStatus)
	viewMenu.AddText("Daemons", keys.CmdOrCtrl("b"), a.ViewDaemons)
	viewMenu.AddText("Settings", keys.CmdOrCtrl("c"), a.ViewSettings)

	helpMenu := appMenu.AddSubmenu("Help")
	helpMenu.AddText("Show Help", keys.CmdOrCtrl("h"), a.HelpToggle)

	return appMenu
}
