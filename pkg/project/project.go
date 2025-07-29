// package project contains the data structures and methods for managing project files
package project

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
)

// ------------------------------------------------------------------------------------
// Project represents a single project with its metadata and data.
type Project struct {
	mu             sync.RWMutex                 `json:"-"`
	Dirty          bool                         `json:"dirty,omitempty"`
	Version        string                       `json:"version"`
	Name           string                       `json:"name"`
	LastOpened     string                       `json:"last_opened"`
	LastView       string                       `json:"lastView"`
	LastFacetMap   map[string]string            `json:"lastFacetMap"`
	Addresses      []base.Address               `json:"addresses"`
	ActiveAddress  base.Address                 `json:"activeAddress"`
	Chains         []string                     `json:"chains"`
	ActiveChain    string                       `json:"activeChain"`
	Contracts      []string                     `json:"contracts"`
	ActiveContract string                       `json:"activeContract"`
	FilterStates   map[ViewStateKey]FilterState `json:"filterStates"`
	Path           string                       `json:"-"`
}

// ------------------------------------------------------------------------------------
// NewProject creates a new project with default values and required active address
func NewProject(name string, activeAddress base.Address, chains []string) *Project {
	addresses := []base.Address{}
	if activeAddress != base.ZeroAddr {
		addresses = append(addresses, activeAddress)
	}
	return &Project{
		Version:        "1.0",
		Name:           name,
		LastOpened:     time.Now().Format(time.RFC3339),
		Dirty:          true,
		LastView:       "",
		LastFacetMap:   map[string]string{},
		ActiveAddress:  activeAddress,
		Addresses:      addresses,
		ActiveChain:    chains[0],
		Chains:         chains,
		ActiveContract: "",
		Contracts:      []string{},
		FilterStates:   make(map[ViewStateKey]FilterState),
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
	_ = file.EstablishFolder(filepath.Dir(path))

	p.LastOpened = time.Now().Format(time.RFC3339)
	p.Dirty = false

	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize project: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write project file: %w", err)
	}

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
	found := false
	for _, existingAddr := range p.Addresses {
		if existingAddr == addr {
			found = true
			break
		}
	}
	if !found {
		// TODO: Bound this at 10?
		p.Addresses = append([]base.Address{addr}, p.Addresses...)
		p.Dirty = true
	}

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
	found := false
	for _, existingAddr := range p.Addresses {
		if existingAddr == addr {
			found = true
			break
		}
	}

	if !found {
		// TODO: Why do we require it to already be seen?
		return fmt.Errorf("address %s not found in project", addr.Hex())
	}

	if p.ActiveAddress != addr {
		p.ActiveAddress = addr
		p.Dirty = true
	}

	return nil
}

// ------------------------------------------------------------------------------------
// AddAddress adds a new address to the project
func (p *Project) AddAddress(addr base.Address) error {
	for _, existingAddr := range p.Addresses {
		if existingAddr == addr {
			return fmt.Errorf("address %s already exists in project", addr.Hex())
		}
	}

	p.Addresses = append(p.Addresses, addr)
	p.Dirty = true
	return nil
}

// ------------------------------------------------------------------------------------
// GetAddresses returns all addresses in the project
func (p *Project) GetAddresses() []base.Address {
	return p.Addresses
}

// ------------------------------------------------------------------------------------
// GetLastView returns the last visited view/route
func (p *Project) GetLastView() string {
	return p.LastView
}

// ------------------------------------------------------------------------------------
// SetLastView updates the last visited view/route and saves immediately (session state)
func (p *Project) SetLastView(view string) error {
	if p.LastView != view {
		p.LastView = strings.Trim(view, "/")
		return p.Save()
	}
	return nil
}

// ------------------------------------------------------------------------------------
// GetLastFacet returns the last visited facet for a specific view
func (p *Project) GetLastFacet(view string) string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.LastFacetMap[view]
}

// ------------------------------------------------------------------------------------
// SetLastFacet updates the last visited facet for a specific view and saves immediately (session state)
func (p *Project) SetLastFacet(view, facet string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.LastFacetMap == nil {
		p.LastFacetMap = make(map[string]string)
	}
	currentFacet := p.LastFacetMap[view]
	if currentFacet != facet {
		p.LastFacetMap[view] = facet
		return p.Save()
	}
	return nil
}

// ------------------------------------------------------------------------------------
// RemoveAddress removes an address from the project
func (p *Project) RemoveAddress(addr base.Address) error {
	for i, existingAddr := range p.Addresses {
		if existingAddr == addr {
			p.Addresses = append(p.Addresses[:i], p.Addresses[i+1:]...)
			if p.ActiveAddress == addr {
				if len(p.Addresses) > 0 {
					p.ActiveAddress = p.Addresses[0]
				} else {
					p.ActiveAddress = base.ZeroAddr
				}
			}
			p.Dirty = true
			return nil
		}
	}
	return fmt.Errorf("address %s not found in project", addr.Hex())
}

// ------------------------------------------------------------------------------------
// GetChains returns all chains in the project
func (p *Project) GetChains() []string {
	return p.Chains
}

// ------------------------------------------------------------------------------------
// GetActiveChain returns the currently selected chain
func (p *Project) GetActiveChain() string {
	return p.ActiveChain
}

// ------------------------------------------------------------------------------------
// SetActiveChain sets the currently selected chain (must be in project)
func (p *Project) SetActiveChain(chain string) error {
	found := false
	for _, existingChain := range p.Chains {
		if existingChain == chain {
			found = true
			break
		}
	}

	if !found {
		// TODO: Why does it already have to exist?
		return fmt.Errorf("chain %s not found in project", chain)
	}

	if p.ActiveChain != chain {
		p.ActiveChain = chain
		p.Dirty = true
	}
	return nil
}

// ------------------------------------------------------------------------------------
// AddChain adds a new chain to the project
func (p *Project) AddChain(chain string) error {
	for _, existingChain := range p.Chains {
		if existingChain == chain {
			return fmt.Errorf("chain %s already exists in project", chain)
		}
	}
	p.Chains = append(p.Chains, chain)
	p.Dirty = true
	return nil
}

// ------------------------------------------------------------------------------------
// RemoveChain removes a chain from the project
func (p *Project) RemoveChain(chain string) error {
	for i, existingChain := range p.Chains {
		if existingChain == chain {
			p.Chains = append(p.Chains[:i], p.Chains[i+1:]...)
			if p.ActiveChain == chain {
				if len(p.Chains) > 0 {
					p.ActiveChain = p.Chains[0]
				} else {
					p.ActiveChain = "mainnet"
					p.Chains = append(p.Chains, "mainnet")
				}
			}
			p.Dirty = true
			return nil
		}
	}
	return fmt.Errorf("chain %s not found in project", chain)
}

// ------------------------------------------------------------------------------------
// GetContracts returns all contracts in the project
func (p *Project) GetContracts() []string {
	return p.Contracts
}

// ------------------------------------------------------------------------------------
// GetActiveContract returns the currently selected contract
func (p *Project) GetActiveContract() string {
	return p.ActiveContract
}

// ------------------------------------------------------------------------------------
// SetActiveContract sets the currently selected contract (must be in project or empty)
func (p *Project) SetActiveContract(contract string) error {
	if contract == "" {
		if p.ActiveContract != contract {
			p.ActiveContract = contract
			p.Dirty = true
		}
		return nil
	}

	found := false
	for _, existingContract := range p.Contracts {
		if existingContract == contract {
			found = true
			break
		}
	}

	if !found {
		// TODO: Why does it have to exist?
		return fmt.Errorf("contract %s not found in project", contract)
	}

	if p.ActiveContract != contract {
		p.ActiveContract = contract
		p.Dirty = true
	}
	return nil
}

// ------------------------------------------------------------------------------------
// AddContract adds a new contract to the project
func (p *Project) AddContract(contract string) error {
	for _, existingContract := range p.Contracts {
		if existingContract == contract {
			return fmt.Errorf("contract %s already exists in project", contract)
		}
	}
	p.Contracts = append(p.Contracts, contract)
	p.Dirty = true
	return nil
}

// ------------------------------------------------------------------------------------
// RemoveContract removes a contract from the project
func (p *Project) RemoveContract(contract string) error {
	for i, existingContract := range p.Contracts {
		if existingContract == contract {
			p.Contracts = append(p.Contracts[:i], p.Contracts[i+1:]...)
			if p.ActiveContract == contract {
				p.ActiveContract = ""
			}
			p.Dirty = true
			return nil
		}
	}
	return fmt.Errorf("contract %s not found in project", contract)
}

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
	if p.FilterStates == nil {
		p.FilterStates = make(map[ViewStateKey]FilterState)
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.FilterStates[key] = state
	return p.Save() // Save immediately - no dirty flag set
}

// ------------------------------------------------------------------------------------
// ClearFilterState removes view state for a given key and saves immediately (session state)
func (p *Project) ClearFilterState(key ViewStateKey) error {
	if p.FilterStates != nil {
		p.mu.Lock()
		defer p.mu.Unlock()
		delete(p.FilterStates, key)
		return p.Save() // Save immediately - no dirty flag set
	}
	return nil
}

// ------------------------------------------------------------------------------------
// ClearAllFilterStates removes all filter states and saves immediately (session state)
func (p *Project) ClearAllFilterStates() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.FilterStates = make(map[ViewStateKey]FilterState)
	return p.Save()
}
