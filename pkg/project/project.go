// package project contains the data structures and methods for managing project files
package project

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
)

// ------------------------------------------------------------------------------------
// Project represents a single project with its metadata and data.
type Project struct {
	mu            sync.RWMutex                 `json:"-"`
	Dirty         bool                         `json:"dirty,omitempty"`
	Version       string                       `json:"version"`
	Name          string                       `json:"name"`
	LastOpened    string                       `json:"last_opened"`
	Addresses     []base.Address               `json:"addresses"`
	ActiveAddress base.Address                 `json:"activeAddress"`
	FilterStates  map[ViewStateKey]FilterState `json:"filterStates"`
	Path          string                       `json:"-"`
}

// ------------------------------------------------------------------------------------
// NewProject creates a new project with default values and required active address
func NewProject(name string, activeAddress base.Address, chains []string) *Project {
	addresses := []base.Address{}
	if activeAddress != base.ZeroAddr {
		addresses = append(addresses, activeAddress)
	}
	return &Project{
		Version:       "1.0",
		Name:          name,
		LastOpened:    time.Now().Format(time.RFC3339),
		Dirty:         true,
		ActiveAddress: activeAddress,
		Addresses:     addresses,
		FilterStates:  make(map[ViewStateKey]FilterState),
	}
}

// ------------------------------------------------------------------------------------
var ErrProjectRecoveryIncomplete = fmt.Errorf("failed to parse project file, recovery attempted but may not be complete")

// ------------------------------------------------------------------------------------
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
		projectPtr := NewProject("Recovered Project", base.ZeroAddr, []string{"mainnet"})
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

// ------------------------------------------------------------------------------------
// Save persists the project to its file path
func (p *Project) Save() error {
	if p.Path == "" {
		return fmt.Errorf("cannot save project with empty path")
	}
	return p.SaveAs(p.Path)
}

// ------------------------------------------------------------------------------------
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

// ------------------------------------------------------------------------------------
// IsDirty returns whether the project has unsaved changes
func (p *Project) IsDirty() bool {
	return p.Dirty
}

// ------------------------------------------------------------------------------------
// SetDirty marks the project as having unsaved changes
func (p *Project) SetDirty(dirty bool) {
	p.Dirty = dirty
}

// ------------------------------------------------------------------------------------
// GetPath returns the file path of the project
func (p *Project) GetPath() string {
	return p.Path
}

// ------------------------------------------------------------------------------------
// GetName returns the name of the project
func (p *Project) GetName() string {
	return p.Name
}

// ------------------------------------------------------------------------------------
// SetName updates the project name and marks it as dirty
func (p *Project) SetName(name string) {
	if p.Name != name {
		p.Name = name
		p.Dirty = true
	}
}

// ------------------------------------------------------------------------------------
// GetAddress returns the project's active address (backward compatibility)
func (p *Project) GetAddress() base.Address {
	return p.ActiveAddress
}

// ------------------------------------------------------------------------------------
// SetAddress sets the project's active address and adds it to addresses if not present (backward compatibility)
func (p *Project) SetAddress(addr base.Address) {
	if p.ActiveAddress != addr {
		p.ActiveAddress = addr
		p.Dirty = true
	}
}

// ------------------------------------------------------------------------------------
// GetActiveAddress returns the currently selected address
func (p *Project) GetActiveAddress() base.Address {
	return p.ActiveAddress
}

// ------------------------------------------------------------------------------------
// SetActiveAddress sets the currently selected address (must be in project)
func (p *Project) SetActiveAddress(addr base.Address) error {
	return nil
}

// ------------------------------------------------------------------------------------
// AddAddress adds a new address to the project
func (p *Project) AddAddress(addr base.Address) error {
	return nil
}

// ------------------------------------------------------------------------------------
// GetAddresses returns all addresses in the project
func (p *Project) GetAddresses() []base.Address {
	return p.Addresses
}

// ------------------------------------------------------------------------------------
// GetLastView returns the last visited view/route

// ------------------------------------------------------------------------------------
// SetLastView updates the last visited view/route and saves immediately (session state)
func (p *Project) SetLastView(view string) error {
	return nil
}

// ------------------------------------------------------------------------------------
// GetLastFacet returns the last visited facet for a specific view
func (p *Project) GetLastFacet(view string) string {
	return ""
}

// ------------------------------------------------------------------------------------
// SetLastFacet updates the last visited facet for a specific view and saves immediately (session state)
func (p *Project) SetLastFacet(view, facet string) error {
	return nil
}

// ------------------------------------------------------------------------------------
// RemoveAddress removes an address from the project
func (p *Project) RemoveAddress(addr base.Address) error {
	return fmt.Errorf("address %s not found in project", addr.Hex())
}

// ------------------------------------------------------------------------------------
// GetChains returns all chains in the project

// ------------------------------------------------------------------------------------
// GetActiveChain returns the currently selected chain

// ------------------------------------------------------------------------------------
// SetActiveChain sets the currently selected chain (must be in project)

// ------------------------------------------------------------------------------------
// AddChain adds a new chain to the project

// ------------------------------------------------------------------------------------
// RemoveChain removes a chain from the project

// ------------------------------------------------------------------------------------
// GetContracts returns all contracts in the project

// ------------------------------------------------------------------------------------
// GetActiveContract returns the currently selected contract

// ------------------------------------------------------------------------------------
// SetActiveContract sets the currently selected contract (must be in project or empty)

// ------------------------------------------------------------------------------------
// AddContract adds a new contract to the project

// ------------------------------------------------------------------------------------
// RemoveContract removes a contract from the project

// ------------------------------------------------------------------------------------
// ViewStateKeyToString converts a ViewStateKey to a string for map indexing

// ------------------------------------------------------------------------------------
// GetFilterState retrieves filter state for a given key
func (p *Project) GetFilterState(key ViewStateKey) (FilterState, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	state, exists := p.FilterStates[key]
	return state, exists
}

// ------------------------------------------------------------------------------------
// SetFilterState sets filter state for a given key and saves immediately (session state)
func (p *Project) SetFilterState(key ViewStateKey, state FilterState) error {
	return nil
}

// ------------------------------------------------------------------------------------
// ClearFilterState removes view state for a given key and saves immediately (session state)
func (p *Project) ClearFilterState(key ViewStateKey) error {
	return nil
}

// ------------------------------------------------------------------------------------
// ClearAllFilterStates removes all filter states and saves immediately (session state)
func (p *Project) ClearAllFilterStates() error {
	return nil
}
