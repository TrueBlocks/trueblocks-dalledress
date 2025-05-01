package app

import (
	"errors"
	"fmt"
	"os"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/project"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

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

func (a *App) fileNew() error {
	a.Projects.New(a.uniqueProjectName("New Project"))

	// Ensure the newly created project is not marked as dirty
	activeProject := a.Projects.Active()
	if activeProject != nil {
		activeProject.SetDirty(false)
	}

	a.updateRecentProjects()
	return nil
}

func (a *App) fileSave() error {
	project := a.Projects.Active()
	if project == nil {
		return errors.New("no active project")
	}

	if project.GetPath() == "" {
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

func (a *App) fileSaveAs(newPath string, overwriteConfirmed bool) error {
	project := a.Projects.Active()
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

	a.updateRecentProjects()
	return nil
}

func (a *App) updateRecentProjects() {
	activeProject := a.Projects.Active()
	if activeProject == nil || activeProject.GetPath() == "" {
		return
	}

	path := activeProject.GetPath()

	if err := a.Preferences.AddRecentProject(path); err != nil {
		msgs.EmitError("add recent project failed", err)
		return
	}

	msgs.EmitMessage(msgs.EventProjectsUpdated, "")
}

func (a *App) GetFilename() *project.Project {
	return a.Projects.Active()
}

func (a *App) uniqueProjectName(baseName string) string {
	projectExists := func(name string) bool {
		for _, project := range a.Projects.GetOpenProjectIDs() {
			projectObj := a.Projects.GetProjectByID(project)
			if projectObj.Name == name {
				return true
			}
		}
		return false
	}

	count := 1
	uniqueName := baseName
	for {
		if !projectExists(uniqueName) {
			break
		}
		uniqueName = baseName + " " + fmt.Sprintf("%d", count)
		count++
	}

	return uniqueName
}

func (a *App) SwitchToProject(id string) error {
	if a.Projects.GetProjectByID(id) == nil {
		return fmt.Errorf("no project with ID %s exists", id)
	}
	return a.Projects.SetActive(id)
}

func (a *App) CloseProject(id string) error {
	project := a.Projects.GetProjectByID(id)
	if project == nil {
		return fmt.Errorf("no project with ID %s exists", id)
	}

	// Check if project has unsaved changes
	if project.IsDirty() {
		response, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Unsaved Changes",
			Message: fmt.Sprintf("Do you want to save changes to project '%s' before closing?", project.GetName()),
			Buttons: []string{"Yes", "No", "Cancel"},
		})

		if err != nil {
			return err
		}

		switch response {
		case "Yes":
			// Save the project before closing
			if project.GetPath() == "" {
				// Project hasn't been saved before, need to use SaveAs
				path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
					Title: "Save Project Before Closing",
				})
				if err != nil || path == "" {
					return fmt.Errorf("save canceled")
				}

				// Use project's SaveAs method instead of SaveProjectAs
				if err := project.SaveAs(path); err != nil {
					return err
				}
			} else {
				// Project has a path, use normal save
				// Use project's Save method instead of SaveProject
				if err := project.Save(); err != nil {
					return err
				}
			}
		case "Cancel":
			return fmt.Errorf("close canceled")
		}
	}

	return a.Projects.Close(id)
}
