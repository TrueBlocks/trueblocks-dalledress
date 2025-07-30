package app

import (
	"fmt"
	"strings"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/project"
)

// GetLastView returns the last visited view/route in the active project.
func (a *App) GetLastView() string {
	if active := a.GetActiveProject(); active != nil {
		return active.GetLastView()
	}
	return ""
}

// SetLastView sets the last visited view/route in the active project
func (a *App) SetLastView(view string) error {
	if active := a.GetActiveProject(); active != nil {
		return active.SetLastView(view)
	}
	return fmt.Errorf("no active project")
}

// SetLastFacet sets the last visited facet for a specific view in the active project
func (a *App) SetLastFacet(view, facet string) error {
	if active := a.GetActiveProject(); active != nil {
		return active.SetLastFacet(view, facet)
	}
	return fmt.Errorf("no active project")
}

// GetLastFacet returns the last visited facet for a specific view from the active project
func (a *App) GetLastFacet(view string) string {
	if active := a.GetActiveProject(); active != nil {
		return active.GetLastFacet(view)
	}
	return ""
}

// GetFilterState retrieves view state for a given key from the active project
func (a *App) GetFilterState(key project.ViewStateKey) (project.FilterState, error) {
	if active := a.GetActiveProject(); active != nil {
		if state, exists := active.GetFilterState(key); exists {
			return state, nil
		}
	}
	return project.FilterState{}, fmt.Errorf("no active project")
}

// SetFilterState sets view state for a given key in the active project
func (a *App) SetFilterState(key project.ViewStateKey, state project.FilterState) error {
	if active := a.GetActiveProject(); active != nil {
		return active.SetFilterState(key, state)
	}
	return fmt.Errorf("no active project")
}

// ClearFilterState removes filter state for a given key from the active project
func (a *App) ClearFilterState(key project.ViewStateKey) error {
	if active := a.GetActiveProject(); active != nil {
		return active.ClearFilterState(key)
	}
	return fmt.Errorf("no active project")
}

// GetWizardReturn returns the last view without the "/wizard" prefix
func (a *App) GetWizardReturn() string {
	return strings.Replace(a.GetLastView(), "/wizard", "/", -1)
}
