package app

import (
	"os"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/preferences"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) FileNew(_ *menu.CallbackData) {
	// Simplified version without save prompts
	if err := a.fileNew(); err != nil {
		msgs.EmitError("File → New failed", err)
		return
	}
	activeProject := a.Projects.Active()
	msgs.EmitStatus("New file created " + activeProject.GetPath())
}

func (a *App) FileOpen(_ *menu.CallbackData) {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Open Project File",
	})
	if err != nil || path == "" {
		msgs.EmitStatus("No file selected")
		return
	}

	if err := a.fileOpen(path); err != nil {
		msgs.EmitError("Open failed", err)
		return
	}

	msgs.EmitStatus("File opened")
}

func (a *App) FileSave(_ *menu.CallbackData) {
	if err := a.fileSave(); err != nil {
		msgs.EmitError("Save failed", err)
		return
	}
	msgs.EmitStatus("File saved")
}

func (a *App) FileSaveAs(_ *menu.CallbackData) {
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title: "Save Project As",
	})
	if err != nil || path == "" {
		msgs.EmitStatus("Save As canceled")
		return
	}

	if err := a.fileSaveAs(path, true); err != nil {
		msgs.EmitError("Save As failed", err)
		return
	}

	msgs.EmitStatus("File saved as")
}

func (a *App) FileQuit(_ *menu.CallbackData) {
	if a.Projects.HasUnsavedChanges() {
		response, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Unsaved Changes",
			Message: "Do you want to save changes before quitting?",
			Buttons: []string{"Yes", "No", "Cancel"},
		})

		if err != nil {
			msgs.EmitError("Dialog error", err)
			return
		}

		switch response {
		case "Yes":
			if err := a.fileSave(); err != nil {
				msgs.EmitError("Save failed", err)
				return // Don't quit if save fails
			}
			// Continue to quit after successful save
		case "Cancel":
			return // Don't quit if user cancels
		}
	}

	msgs.EmitStatus("Quitting application")
	os.Exit(0)
}

func (a *App) buildAppMenu() *menu.Menu {
	appMenu := menu.NewMenu()

	// System Menu (added before File menu)
	system := appMenu.AddSubmenu("System")
	system.AddText("Preferences...", keys.CmdOrCtrl("5"), func(_ *menu.CallbackData) {
		// This matches the cmd+5 keyboard shortcut
		// a.ShowPage("settings")
	})
	system.AddSeparator()
	// TODO: add applicastion name to this menu item
	system.AddText("Quit", keys.CmdOrCtrl("q"), a.FileQuit)

	// File Menu
	file := appMenu.AddSubmenu("File")
	file.AddText("New", keys.CmdOrCtrl("n"), a.FileNew)
	file.AddText("Open", keys.CmdOrCtrl("o"), a.FileOpen)
	file.AddText("Save", keys.CmdOrCtrl("s"), a.FileSave)
	file.AddText("Save As", keys.CmdOrCtrl("shift+s"), a.FileSaveAs)

	// Edit Menu
	edit := appMenu.AddSubmenu("Edit")
	edit.AddText("Cut", keys.CmdOrCtrl("x"), nil)        // menu.EditCut)
	edit.AddText("Copy", keys.CmdOrCtrl("c"), nil)       // menu.EditCopy)
	edit.AddText("Paste", keys.CmdOrCtrl("v"), nil)      // menu.EditPaste)
	edit.AddText("Select All", keys.CmdOrCtrl("a"), nil) // menu.EditSelectAll)

	// Window Menu
	window := appMenu.AddSubmenu("Window")
	window.AddText("Minimize", keys.CmdOrCtrl("m"), nil) // menu.WindowMinimize)
	window.AddText("Zoom", nil, nil)                     // menu.WindowZoom)

	// Help Menu
	help := appMenu.AddSubmenu("Help")
	// TODO: add applicastion name to this menu item
	aboutLink := "https://" + preferences.GetAppId().Domain + "/about"
	help.AddText("About", nil, func(_ *menu.CallbackData) {
		runtime.BrowserOpenURL(a.ctx, aboutLink)
	})
	help.AddText("Report Issue", nil, func(_ *menu.CallbackData) {
		runtime.BrowserOpenURL(a.ctx, preferences.GetAppId().Github+"/issues")
	})

	return appMenu
}
