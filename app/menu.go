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
	viewMenu.AddText("View 1", keys.CmdOrCtrl("1"), a.ViewHome)
	viewMenu.AddText("View 2", keys.CmdOrCtrl("2"), a.ViewDalle)
	viewMenu.AddText("View 3", keys.CmdOrCtrl("3"), a.ViewSeries)
	viewMenu.AddText("View 4", keys.CmdOrCtrl("4"), a.ViewHistory)
	viewMenu.AddText("View 5", keys.CmdOrCtrl("5"), a.ViewNames)
	viewMenu.AddText("View 6", keys.CmdOrCtrl("6"), a.ViewServers)
	viewMenu.AddText("View 7", keys.CmdOrCtrl("7"), a.ViewSettings)
	return appMenu
}
