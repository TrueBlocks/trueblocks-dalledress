package app

import (
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
)

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
	viewMenu.AddText("Names", keys.CmdOrCtrl("5"), a.ViewNames)
	viewMenu.AddText("Servers", keys.CmdOrCtrl("6"), a.ViewServers)
	viewMenu.AddText("Settings", keys.CmdOrCtrl("7"), a.ViewSettings)

	return appMenu
}
