package app

import (
	"testing"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/preferences"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/project"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

// TestViewConfigFacetOrderIntegrity ensures every backend ViewConfig supplies a non-empty
// FacetOrder slice whose entries all correspond to actual keys in the Facets map.
// It will also fail if any facet present in the Facets map is omitted from FacetOrder
// (strong form – adjust if you later allow hidden facets).
func TestViewConfigFacetOrderIntegrity(t *testing.T) {
	a := &App{
		Projects:    project.NewManager(),
		Preferences: &preferences.Preferences{},
		apiKeys:     map[string]string{},
	}

	// Minimal baseline payload (fields not used by Get*Config implementations).
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
			if len(cfg.FacetOrder) == 0 {
				t.Fatalf("%s: FacetOrder is empty", tc.name)
			}
			// Track seen facets
			seen := make(map[string]bool, len(cfg.FacetOrder))
			for i, id := range cfg.FacetOrder {
				if id == "" {
					t.Errorf("%s: empty facet id at position %d", tc.name, i)
					continue
				}
				if seen[id] {
					t.Errorf("%s: duplicate facet id %q in FacetOrder", tc.name, id)
				}
				seen[id] = true
				if _, ok := cfg.Facets[id]; !ok {
					t.Errorf("%s: facet %q in FacetOrder not found in Facets map", tc.name, id)
				}
			}
			// Strong form: ensure no facet omitted from ordering.
			for k := range cfg.Facets {
				if !seen[k] {
					t.Errorf("%s: facet %q present in Facets but missing from FacetOrder", tc.name, k)
				}
			}
		})
	}
}
