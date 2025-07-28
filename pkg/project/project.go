// package project contains the data structures and methods for managing project files
package project

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
)

// Project represents a single project with its metadata and data.
type Project struct {
	Version    string       `json:"version"`
	Name       string       `json:"name"`
	LastOpened string       `json:"last_opened"`
	Dirty      bool         `json:"dirty"`
	Address    base.Address `json:"address"`
	Path       string       `json:"-"`
}

// New creates a new project with default values
func New(name string, address base.Address) *Project {
	return &Project{
		Version:    "1.0",
		Name:       name,
		LastOpened: time.Now().Format(time.RFC3339),
		Dirty:      true,
		Address:    address,
	}
}

var ErrProjectRecoveryIncomplete = fmt.Errorf("failed to parse project file, recovery attempted but may not be complete")

// Load loads a project from the specified file path with optimized deserialization
func Load(path string) (*Project, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("project file does not exist: %s", path)
	}

	// Using a buffered read approach for better performance
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open project file: %w", err)
	}
	defer file.Close()

	// For small files like our projects, ReadAll is actually quite efficient
	// It avoids multiple small reads and allocations
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read project file: %w", err)
	}

	var project Project
	if err := json.Unmarshal(data, &project); err != nil {
		projectPtr := New("Recovered Project", base.ZeroAddr)
		projectPtr.Dirty = true // Mark as dirty so user knows it was recovered
		projectPtr.Path = path
		if saveErr := projectPtr.Save(); saveErr != nil {
			return nil, fmt.Errorf("failed to parse project file and could not save recovered version: %w (original error: %v)", saveErr, err)
		}
		return nil, ErrProjectRecoveryIncomplete
	}

	// Set in-memory fields
	project.Path = path
	project.Dirty = false

	return &project, nil
}

// Save persists the project to its file path
func (p *Project) Save() error {
	if p.Path == "" {
		return fmt.Errorf("cannot save project with empty path")
	}
	return p.SaveAs(p.Path)
}

// SaveAs saves the project to a new file path and updates the project's path
// with optimized serialization for better performance
func (p *Project) SaveAs(path string) error {
	// Ensure the directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Update last opened timestamp
	p.LastOpened = time.Now().Format(time.RFC3339)
	p.Dirty = false

	// Create a temporary file for safe writing
	tempPath := path + ".tmp"

	// Optimize serialization with a single marshal operation
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize project: %w", err)
	}

	// Write to temporary file first
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		// Clean up temporary file if write fails
		os.Remove(tempPath)
		return fmt.Errorf("failed to write project file: %w", err)
	}

	// Atomically rename temporary file to final path for better durability
	if err := os.Rename(tempPath, path); err != nil {
		// Clean up temporary file if rename fails
		os.Remove(tempPath)
		return fmt.Errorf("failed to finalize project file: %w", err)
	}

	// Update in-memory state
	p.Path = path
	p.Dirty = false

	return nil
}

// IsDirty returns whether the project has unsaved changes
func (p *Project) IsDirty() bool {
	return p.Dirty
}

// SetDirty marks the project as having unsaved changes
func (p *Project) SetDirty(dirty bool) {
	p.Dirty = dirty
}

// GetPath returns the file path of the project
func (p *Project) GetPath() string {
	return p.Path
}

// GetName returns the name of the project
func (p *Project) GetName() string {
	return p.Name
}

// SetName updates the project name and marks it as dirty
func (p *Project) SetName(name string) {
	if p.Name != name {
		p.Name = name
		p.Dirty = true
	}
}

// GetAddress returns the project's main address
func (p *Project) GetAddress() base.Address {
	return p.Address
}

// SetAddress sets the project's main address and marks as dirty
func (p *Project) SetAddress(addr base.Address) {
	if p.Address != addr {
		p.Address = addr
		p.Dirty = true
	}
}
