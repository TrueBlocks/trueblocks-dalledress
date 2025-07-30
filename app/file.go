package app

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	"github.com/wailsapp/wails/v2/pkg/menu"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// FileNew creates a new project file with default settings
func (a *App) FileNew(_ *menu.CallbackData) {
	if err := a.fileNew(base.ZeroAddr); err != nil {
		msgs.EmitError("File → New failed", err)
		return
	}
	activeProject := a.GetActiveProject()
	msgs.EmitStatus("new file created " + activeProject.GetPath())
}

// FileOpen clears the current project and opens the project selection dialog
func (a *App) FileOpen(_ *menu.CallbackData) {
	if err := a.ClearActiveProject(); err != nil {
		return
	}

	msgs.EmitManager("show_project_modal")
	msgs.EmitStatus("project selection dialog opened")
}

// FileSave saves the active project to its current file path
func (a *App) FileSave(_ *menu.CallbackData) {
	if err := a.SaveProject(); err != nil {
		msgs.EmitError("save failed", err)
		return
	}
	msgs.EmitStatus("file saved")
}

// FileSaveAs opens a save dialog and saves the project to a new file path
func (a *App) FileSaveAs(_ *menu.CallbackData) {
	path, err := wailsRuntime.SaveFileDialog(a.ctx, wailsRuntime.SaveDialogOptions{
		Title: "Save Project As",
		Filters: []wailsRuntime.FileFilter{
			{
				DisplayName: "TrueBlocks Project Files (*.tbx)",
				Pattern:     "*.tbx",
			},
		},
	})
	if err != nil || path == "" {
		msgs.EmitStatus("save As canceled")
		return
	}

	if err := a.fileSaveAs(path, true); err != nil {
		msgs.EmitError("save As failed", err)
		return
	}

	msgs.EmitStatus("file saved as")
}

// FileQuit shuts down the application after saving if needed
func (a *App) FileQuit(_ *menu.CallbackData) {
	if a.Projects.HasUnsavedChanges() {
		response, err := wailsRuntime.MessageDialog(a.ctx, wailsRuntime.MessageDialogOptions{
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

	msgs.EmitStatus("quitting application")
	os.Exit(0)
}

// Generic errors
var ErrEmptyFilePath = errors.New("empty file path")
var ErrUnsavedChanges = errors.New("unsaved changes")
var ErrFileNotFound = errors.New("file not found")
var ErrOverwriteNotConfirmed = errors.New("file exists, overwrite not confirmed")

// File operation errors
var ErrReadFileFailed = errors.New("failed to read file")
var ErrWriteFileFailed = errors.New("failed to write file")
var ErrSerializeFailed = errors.New("failed to serialize data")
var ErrDeserializeFailed = errors.New("failed to deserialize data")

// fileNew creates a new project with the given address and default settings
func (a *App) fileNew(address base.Address) error {
	a.Projects.NewProject(a.uniqueProjectName("New Project"), address, []string{"mainnet"})

	// Ensure the newly created project is not marked as dirty
	activeProject := a.GetActiveProject()
	if activeProject != nil {
		activeProject.SetDirty(false)
	}

	a.updateRecentProjects()
	return nil
}

// fileSave saves the active project to its current file path
func (a *App) fileSave() error {
	project := a.GetActiveProject()
	if project == nil {
		return errors.New("no active project")
	}

	projectPath := project.GetPath()
	needsSaveAs := projectPath == "" || strings.Contains(filepath.Base(projectPath), "Unknown")
	if needsSaveAs {
		return ErrEmptyFilePath
	}

	if !project.IsDirty() {
		return nil
	}

	if err := a.Projects.SaveActive(); err != nil {
		return err
	}

	a.updateRecentProjects()
	return nil
}

// fileSaveAs saves the active project to a new file path with overwrite handling
func (a *App) fileSaveAs(newPath string, overwriteConfirmed bool) error {
	project := a.GetActiveProject()
	if project == nil {
		return errors.New("no active project")
	}

	if newPath == "" {
		return ErrEmptyFilePath
	}

	if file.FileExists(newPath) && !overwriteConfirmed {
		return ErrOverwriteNotConfirmed
	}

	if err := a.Projects.SaveActiveAs(newPath); err != nil {
		return err
	}

	a.updateRecentProjects()
	return nil
}

// fileOpen opens a project file from the specified path
func (a *App) fileOpen(path string) error {
	if path == "" {
		return ErrEmptyFilePath
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return ErrFileNotFound
	}

	_, err := a.Projects.Open(path)
	if err != nil {
		return err
	}

	a.SetActiveProjectPath(path)
	a.updateRecentProjects()
	return nil
}
