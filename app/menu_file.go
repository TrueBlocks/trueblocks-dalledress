package app

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/messages"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) FileNew(cd *menu.CallbackData) {
	logger.Info("File New")
}

func (a *App) FileOpen(cd *menu.CallbackData) {
	file, _ := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		DefaultDirectory:           "/Users/jrush/Documents/",
		Title:                      "Open File",
		CanCreateDirectories:       true,
		ShowHiddenFiles:            false,
		ResolvesAliases:            false,
		TreatPackagesAsDirectories: false,
		Filters: []runtime.FileFilter{
			{DisplayName: "Monitor Groups", Pattern: "*.tbx"},
		},
	})
	a.CurrentDoc.Filename = file
	// a.CurrentDoc.Load()
	messages.Send(a.ctx, messages.Document, messages.NewDocumentMsg(a.CurrentDoc.Filename, "Opened"))
}

func (a *App) getFilename() error {
	filename, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		DefaultDirectory:           "/Users/jrush/Documents/",
		DefaultFilename:            a.CurrentDoc.Filename,
		Title:                      "Save File",
		CanCreateDirectories:       true,
		ShowHiddenFiles:            false,
		TreatPackagesAsDirectories: false,
		Filters: []runtime.FileFilter{
			{DisplayName: "Monitor Groups", Pattern: "*.tbx"},
		},
	})
	if err != nil {
		return err
	}
	if filename == "" {
		return fmt.Errorf("user hit escape")
	}
	a.CurrentDoc.Filename = filename
	return nil
}

func (a *App) FileSave(cd *menu.CallbackData) {
	if a.CurrentDoc.Filename == "Untitled" {
		if err := a.getFilename(); err != nil {
			messages.Send(a.ctx, messages.Error, messages.NewErrorMsg(err))
			return
		}
	}
	a.CurrentDoc.Save()
	messages.Send(a.ctx, messages.Document, messages.NewDocumentMsg(a.CurrentDoc.Filename, "Saved"))
}

func (a *App) FileSaveAs(cd *menu.CallbackData) {
	if err := a.getFilename(); err != nil {
		messages.Send(a.ctx, messages.Error, messages.NewErrorMsg(err))
		return
	}
	a.CurrentDoc.Save()
	messages.Send(a.ctx, messages.Document, messages.NewDocumentMsg(a.CurrentDoc.Filename, "Saved"))
}
