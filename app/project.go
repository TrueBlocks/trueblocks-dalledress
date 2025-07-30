package app

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/project"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// NewProject creates a new project with the given name and optional address
func (a *App) NewProject(name string, currentAddress string) error {
	var addr base.Address
	if currentAddress != "" {
		convertedAddr, ok := a.ConvertToAddress(currentAddress)
		if !ok {
			return fmt.Errorf("invalid address: %s", currentAddress)
		}
		addr = convertedAddr
	}

	// Create new project with current address
	newProject := a.Projects.NewProject(name, addr, []string{"mainnet"})
	if newProject == nil {
		return fmt.Errorf("failed to create project")
	}
	return nil
}

// OpenProjectFile opens a project file, optionally showing file picker if no path provided
func (a *App) OpenProjectFile(path string) error {
	if path == "" {
		// Trigger file picker
		selectedPath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
			Title: "Open Project File",
			Filters: []runtime.FileFilter{
				{
					DisplayName: "TrueBlocks Project Files (*.tbx)",
					Pattern:     "*.tbx",
				},
			},
		})
		if err != nil || selectedPath == "" {
			return fmt.Errorf("no file selected")
		}
		path = selectedPath
	}

	if err := a.fileOpen(path); err != nil {
		return err
	}

	// Get the opened project to restore full context
	active := a.GetActiveProject()
	if active == nil {
		return fmt.Errorf("failed to get active project after opening")
	}

	// Get the active project ID from the manager
	activeProjectID := a.Projects.ActiveID
	if activeProjectID == "" {
		return fmt.Errorf("no active project ID available")
	}

	// Trigger analytical state restoration by restoring project context
	if err := a.RestoreProjectContext(activeProjectID); err != nil {
		return fmt.Errorf("failed to restore project context: %w", err)
	}

	// Emit project opened event for frontend navigation and pagination reset
	lastView := active.GetLastView()
	if lastView == "" {
		lastView = "/" // Default to home if no last view
	}
	msgs.EmitProjectOpened(lastView, map[string]interface{}{
		"projectId":     activeProjectID,
		"projectName":   active.GetName(),
		"activeAddress": active.GetActiveAddress().String(),
		"activeChain":   active.GetActiveChain(),
	})

	msgs.EmitStatus(fmt.Sprintf("project opened: %s - restoring context", active.GetName()))
	return nil
}

// SaveProject saves the active project, showing file dialog if not previously saved
func (a *App) SaveProject() error {
	project := a.GetActiveProject()
	if project == nil {
		return fmt.Errorf("no active project")
	}

	if project.GetPath() == "" {
		// Project hasn't been saved before, need to use SaveAs with file dialog
		path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
			Title: "Save Project",
			Filters: []runtime.FileFilter{
				{
					DisplayName: "TrueBlocks Project Files (*.tbx)",
					Pattern:     "*.tbx",
				},
			},
			DefaultFilename: project.GetName() + ".tbx",
		})
		if err != nil {
			return fmt.Errorf("file dialog error: %w", err)
		}
		if path == "" {
			return fmt.Errorf("save canceled")
		}

		if err := a.Projects.SaveActiveAs(path); err != nil {
			return err
		}
	} else {
		// Project has a path, use normal save
		if err := a.Projects.SaveActive(); err != nil {
			return err
		}
	}

	a.updateRecentProjects()
	msgs.EmitStatus(fmt.Sprintf("project saved: %s", project.GetName()))
	return nil
}

// HasActiveProject checks if there is an active project loaded
func (a *App) HasActiveProject() bool {
	if a.Projects == nil {
		return false
	}
	return a.Projects.ActiveID != "" && a.GetActiveProject() != nil
}

// ValidateActiveProject checks if active project has valid addresses configured
func (a *App) ValidateActiveProject() bool {
	if !a.HasActiveProject() {
		return false
	}

	project := a.GetActiveProject()
	if project == nil {
		return false
	}

	// Check if project has at least one address and a valid current address
	addresses := project.GetAddresses()
	activeAddr := project.GetActiveAddress()
	return len(addresses) > 0 && activeAddr != base.ZeroAddr
}

// ClearActiveProject clears the active project, prompting to save if dirty
func (a *App) ClearActiveProject() error {
	if !a.HasActiveProject() {
		return nil
	}

	project := a.GetActiveProject()
	if project == nil {
		a.Projects.ActiveID = ""
		return nil
	}

	if project.IsDirty() {
		response, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:          runtime.QuestionDialog,
			Title:         "Unsaved Changes",
			Message:       fmt.Sprintf("Project '%s' has unsaved changes. What would you like to do?", project.GetName()),
			Buttons:       []string{"Save", "Discard Changes", "Cancel"},
			DefaultButton: "Save",
		})

		if err != nil {
			return fmt.Errorf("failed to show unsaved changes dialog: %w", err)
		}

		switch response {
		case "Save":
			if err := a.SaveProject(); err != nil {
				return fmt.Errorf("failed to save project: %w", err)
			}
		case "Cancel":
			return fmt.Errorf("operation canceled by user")
		case "Discard Changes":
			// Continue with clearing
		}
	}

	// Clear the active project
	a.Projects.ActiveID = ""
	msgs.EmitManager("active_project_cleared")
	return nil
}

// GetOpenProjects returns a list of all open projects with their metadata
func (a *App) GetOpenProjects() []map[string]interface{} {
	projectIDs := a.Projects.GetOpenProjectIDs()
	result := make([]map[string]interface{}, 0, len(projectIDs))

	for _, id := range projectIDs {
		project := a.Projects.GetProjectByID(id)
		if project == nil {
			continue
		}

		projectInfo := map[string]interface{}{
			"id":         id,
			"name":       project.GetName(),
			"path":       project.GetPath(),
			"isActive":   id == a.Projects.ActiveID,
			"isDirty":    project.IsDirty(),
			"lastOpened": project.LastOpened,
			// "createdAt":  project.CreatedAt,
		}

		result = append(result, projectInfo)
	}

	return result
}

func (a *App) GetActiveProject() *project.Project {
	return a.Projects.GetActiveProject()
}

// SwitchToProject sets the specified project as the active project
func (a *App) SwitchToProject(id string) error {
	if a.Projects.GetProjectByID(id) == nil {
		return fmt.Errorf("no project with ID %s exists", id)
	}
	return a.Projects.SetActiveProject(id)
}

// RestoreProjectContext restores the full context for the specified project
func (a *App) RestoreProjectContext(projectID string) error {
	if err := a.SwitchToProject(projectID); err != nil {
		return err
	}

	project := a.GetActiveProject()
	if project == nil {
		return fmt.Errorf("failed to get active project after switching to %s", projectID)
	}

	a.Preferences.App.LastProject = projectID
	if err := a.SetAppPreferences(&a.Preferences.App); err != nil {
		return fmt.Errorf("failed to update app preferences: %w", err)
	}

	return nil
}

// CloseProject closes a project, prompting to save if it has unsaved changes
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
					Filters: []runtime.FileFilter{
						{
							DisplayName: "TrueBlocks Project Files (*.tbx)",
							Pattern:     "*.tbx",
						},
					},
					DefaultFilename: project.GetName() + ".tbx",
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

// GetProjectAddress returns the address of the active project
func (a *App) GetProjectAddress() base.Address {
	active := a.GetActiveProject()
	if active == nil {
		return base.ZeroAddr
	}
	return active.GetAddress()
}

// SetProjectAddress sets the address for the active project
func (a *App) SetProjectAddress(addr base.Address) {
	active := a.GetActiveProject()
	if active != nil {
		active.SetAddress(addr)
	}
}

// GetFilename returns the active project (legacy method for compatibility)
func (a *App) GetFilename() *project.Project {
	return a.GetActiveProject()
}

// uniqueProjectName generates a unique project name by appending numbers if needed
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
