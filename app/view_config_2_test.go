package app

import (
	"testing"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/preferences"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/project"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

// TestViewConfigHeaderActionsExport ensures that:
// 1) FacetConfig.HeaderActions is never nil (must be an empty slice when none)
// 2) Every non-form facet (IsForm == false) includes "export" in HeaderActions
func TestViewConfigHeaderActionsExport(t *testing.T) {
	a := &App{
		Projects:    project.NewManager(),
		Preferences: &preferences.Preferences{},
		apiKeys:     map[string]string{},
	}

	base := types.Payload{Chain: "mainnet", Address: "0x0", Period: types.PeriodBlockly}

	tests := []struct {
		name        string
		getter      func(types.Payload) (*types.ViewConfig, error)
		collection  string
		sampleFacet types.DataFacet
	}{
		{"abis", a.GetAbisConfig, "abis", types.DataFacet("downloaded")},
		{"chunks", a.GetChunksConfig, "chunks", types.DataFacet("stats")},
		{"contracts", a.GetContractsConfig, "contracts", types.DataFacet("dashboard")},
		{"exports", a.GetExportsConfig, "exports", types.DataFacet("statements")},
		{"monitors", a.GetMonitorsConfig, "monitors", types.DataFacet("monitors")},
		{"names", a.GetNamesConfig, "names", types.DataFacet("all")},
		{"status", a.GetStatusConfig, "status", types.DataFacet("status")},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			payload := base
			payload.Collection = tc.collection
			payload.DataFacet = tc.sampleFacet

			cfg, err := tc.getter(payload)
			if err != nil {
				t.Fatalf("%s: error retrieving config: %v", tc.name, err)
			}
			if cfg == nil {
				t.Fatalf("%s: nil config returned", tc.name)
			}

			for facetKey, facet := range cfg.Facets {
				if facet.HeaderActions == nil {
					t.Errorf("%s: facet %q has nil HeaderActions (must be empty slice when none)", tc.name, facetKey)
				}
				if !facet.IsForm {
					if !containsString(facet.HeaderActions, "export") {
						t.Errorf("%s: non-form facet %q missing required 'export' in HeaderActions", tc.name, facetKey)
					}
				}
			}
		})
	}
}

func containsString(list []string, needle string) bool {
	for _, s := range list {
		if s == needle {
			return true
		}
	}
	return false
}
