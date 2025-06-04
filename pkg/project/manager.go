// package project contains the data structures and methods for managing project files
package project

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
)

const (
	ProjectActivated  = "project_activated"
	ProjectCreated    = "project_created"
	ProjectSwitched   = "project_switched"
	ProjectOpened     = "project_opened"
	ProjectClosed     = "project_closed"
	AllProjectsClosed = "all_projects_closed"
	ProjectSaved      = "project_saved"
	ProjectSavedAs    = "project_saved_as"
)

// Manager handles multiple Project instances, maintaining a collection of open projects
// and tracking which one is currently active.
type Manager struct {
	OpenProjects map[string]*Project
	ActiveID     string
}

// NewManager creates a new project manager with no open projects
func NewManager() *Manager {
	return &Manager{
		OpenProjects: make(map[string]*Project),
		ActiveID:     "",
	}
}

// Active returns the currently active project, or nil if no project is active
func (m *Manager) Active() *Project {
	if m.ActiveID == "" {
		return nil
	}
	return m.OpenProjects[m.ActiveID]
}

// SetActive sets the active project by ID
func (m *Manager) SetActive(id string) error {
	if m.ActiveID == id {
		return nil
	}

	if _, exists := m.OpenProjects[id]; !exists {
		return fmt.Errorf("no project with ID %s exists", id)
	}

	if m.ActiveID != "" {
		prevProject := m.OpenProjects[m.ActiveID]
		if prevProject != nil {
			m.minimizeInactiveProject(prevProject)
		}
	}

	m.ActiveID = id
	msgs.EmitMessage(msgs.EventManager, ProjectActivated)
	return nil
}

// minimizeInactiveProject reduces memory usage of inactive projects
// This method provides a framework for memory optimization that can be
// extended in the future as project complexity grows
func (m *Manager) minimizeInactiveProject(project *Project) {
	// This is a hook for future memory optimization
	// For now, we don't have specific memory-intensive structures to clear
	// but having this framework in place allows for easy extension later

	// Example memory optimization patterns that could be implemented:
	// 1. Clear any cached computation results
	// 2. Clear any non-essential data structures that can be rebuilt
	// 3. If T implements a ClearCache() method, call it
}

// New creates a new project with the given name and makes it the active project
func (m *Manager) New(name string) *Project {
	project := New(name)
	id := name
	m.OpenProjects[id] = project
	m.ActiveID = id
	msgs.EmitMessage(msgs.EventManager, ProjectCreated)
	return project
}

// Open loads a project from the specified path and makes it the active project
func (m *Manager) Open(path string) (*Project, error) {
	for id, proj := range m.OpenProjects {
		if proj.Path == path {
			m.ActiveID = id
			msgs.EmitMessage(msgs.EventManager, ProjectSwitched)
			return proj, nil
		}
	}

	project, err := Load(path)
	if err != nil {
		return nil, err
	}

	id := filepath.Base(path)
	m.OpenProjects[id] = project
	m.ActiveID = id

	project.LastOpened = time.Now().Format(time.RFC3339)

	msgs.EmitMessage(msgs.EventManager, ProjectOpened)
	return project, nil
}

// Close closes the project with the given ID. If it's the active project,
// the active project becomes nil.
func (m *Manager) Close(id string) error {
	if _, exists := m.OpenProjects[id]; !exists {
		return fmt.Errorf("no project with ID %s exists", id)
	}

	delete(m.OpenProjects, id)

	if m.ActiveID == id {
		m.ActiveID = ""

		for newID := range m.OpenProjects {
			m.ActiveID = newID
			break
		}
	}

	msgs.EmitMessage(msgs.EventManager, ProjectClosed)
	return nil
}

// CloseAll closes all open projects
func (m *Manager) CloseAll() {
	m.OpenProjects = make(map[string]*Project)
	m.ActiveID = ""
	msgs.EmitMessage(msgs.EventManager, AllProjectsClosed)
}

// SaveActive saves the currently active project
func (m *Manager) SaveActive() error {
	if m.ActiveID == "" {
		return errors.New("no active project to save")
	}

	project := m.OpenProjects[m.ActiveID]

	if strings.HasPrefix(project.Name, "New Project") {
		return m.SaveActiveAs("")
	}

	err := project.Save()
	if err == nil {
		msgs.EmitMessage(msgs.EventManager, ProjectSaved)
	} else {
		msgs.EmitMessage(msgs.EventError, "Save project: "+err.Error())
	}
	return err
}

// SaveActiveAs saves the currently active project to a new path
func (m *Manager) SaveActiveAs(path string) error {
	if m.ActiveID == "" {
		return errors.New("no active project to save")
	}

	project := m.OpenProjects[m.ActiveID]

	if strings.HasPrefix(project.Name, "New Project") {
		wasDirty := project.Dirty
		project.Dirty = true

		err := project.SaveAs(path)

		if !wasDirty && err == nil {
			project.Dirty = false
		}

		if err == nil {
			msgs.EmitMessage(msgs.EventManager, ProjectSavedAs)
		} else {
			msgs.EmitMessage(msgs.EventError, "Save project as: "+err.Error())
		}
		return err
	}

	err := project.SaveAs(path)
	if err == nil {
		msgs.EmitMessage(msgs.EventManager, ProjectSavedAs)
	} else {
		msgs.EmitMessage(msgs.EventError, "Save project as: "+err.Error())
	}
	return err
}

// GetOpenProjectIDs returns a slice of IDs for all open projects
func (m *Manager) GetOpenProjectIDs() []string {
	ids := make([]string, 0, len(m.OpenProjects))
	for id := range m.OpenProjects {
		ids = append(ids, id)
	}
	return ids
}

// HasUnsavedChanges returns true if any open project has unsaved changes
func (m *Manager) HasUnsavedChanges() bool {
	for _, project := range m.OpenProjects {
		if project.IsDirty() {
			return true
		}
	}
	return false
}

// GetProjectByID returns the project with the given ID, or nil if it doesn't exist
func (m *Manager) GetProjectByID(id string) *Project {
	return m.OpenProjects[id]
}

// GetProjectByPath returns the project with the given path, or nil if it doesn't exist
func (m *Manager) GetProjectByPath(path string) *Project {
	for _, project := range m.OpenProjects {
		if project.Path == path {
			return project
		}
	}
	return nil
}
