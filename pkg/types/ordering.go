package types

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/preferences"
)

//go:embed disablements.json
var disablementsJSON string

// EnableFacets enables or disables facets based on the view's enablement file.
type DisablementsConfig struct {
	Views map[string]struct {
		Disabled bool            `json:"disabled"`
		Facets   map[string]bool `json:"facets"`
	} `json:"views"`
}

func SetDisablements(vc *ViewConfig) {
	if vc == nil || vc.Facets == nil {
		return
	}

	var disablements DisablementsConfig
	dec := json.NewDecoder(strings.NewReader(disablementsJSON))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&disablements); err != nil {
		fmt.Printf("Failed to decode embedded disablements: %v\n", err)
		return
	}

	disablement, ok := disablements.Views[vc.ViewName]
	if !ok {
		vc.Disabled = false
		for key, facet := range vc.Facets {
			facet.Disabled = false
			vc.Facets[key] = facet
		}
		return
	}
	vc.Disabled = disablement.Disabled
	for key, facet := range vc.Facets {
		if disabled, exists := disablement.Facets[key]; exists {
			facet.Disabled = disabled
		} else {
			facet.Disabled = false
		}
		vc.Facets[key] = facet
	}
}

// NormalizeOrders sorts columns and detail fields by their explicit order values.
// It does not assign defaults; ordering must be provided in config.
func NormalizeOrders(vc *ViewConfig) {
	if vc == nil || vc.Facets == nil {
		return
	}
	for key, facet := range vc.Facets {
		// Columns: sort by Order when both orders are positive; otherwise keep input order
		sort.SliceStable(facet.Columns, func(i, j int) bool {
			oi, oj := facet.Columns[i].Order, facet.Columns[j].Order
			if oi > 0 && oj > 0 {
				return oi < oj
			}
			return false
		})

		// Detail panels: within each panel, sort by DetailOrder when both orders are positive
		for p := range facet.DetailPanels {
			sort.SliceStable(facet.DetailPanels[p].Fields, func(i, j int) bool {
				oi, oj := facet.DetailPanels[p].Fields[i].DetailOrder, facet.DetailPanels[p].Fields[j].DetailOrder
				if oi > 0 && oj > 0 {
					return oi < oj
				}
				return false
			})
		}

		vc.Facets[key] = facet
	}
}

// SetMenuOrder applies menu ordering and facet configurations from .create-local-app.json to ViewConfig
func SetMenuOrder(vc *ViewConfig) {
	if vc == nil {
		return
	}

	// Load app config
	config, err := preferences.LoadAppConfig()
	if err != nil {
		// If config fails to load, use default order (999)
		vc.MenuOrder = 999
		return
	}

	// Check if this view has configuration
	if viewConfig, exists := config.ViewConfig[vc.ViewName]; exists {
		// Apply menu order
		if viewConfig.MenuOrder > 0 {
			vc.MenuOrder = viewConfig.MenuOrder
		} else {
			vc.MenuOrder = 999 // Default order for views without explicit order
		}

		// Apply view-level disabled state
		vc.Disabled = viewConfig.Disabled

		// Apply facet configurations if both exist
		if len(viewConfig.DisabledFacets) > 0 && vc.Facets != nil {
			for facetName, facetConfig := range vc.Facets {
				if disabledState, facetExists := viewConfig.DisabledFacets[facetName]; facetExists {
					// Apply the configured disabled state directly (true = disabled, false = enabled)
					facetConfig.Disabled = disabledState
					vc.Facets[facetName] = facetConfig
				}
			}
		}
	} else {
		vc.MenuOrder = 999 // Default order for views not in config
	}
}
