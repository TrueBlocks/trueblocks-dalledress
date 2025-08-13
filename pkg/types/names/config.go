package names

import (
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

// GetConfig returns the ViewConfig for the Names view
func (c *NamesCollection) GetConfig() (*types.ViewConfig, error) {
	// Create facet configurations - Names uses tables for all facets (no forms)
	facets := map[string]types.FacetConfig{
		"all": {
			Name:          "All",
			Store:         "names",
			IsForm:        false,
			Fields:        getNamesFields(),
			HeaderActions: []string{"add", "export"},
			Actions:       []string{"update"},
		},
		"custom": {
			Name:          "Custom",
			Store:         "names",
			IsForm:        false,
			Fields:        getNamesFields(),
			HeaderActions: []string{"add", "publish", "pin", "export"},
			Actions:       []string{"delete", "remove", "update"},
		},
		"prefund": {
			Name:          "Prefund",
			Store:         "names",
			IsForm:        false,
			DividerBefore: true, // Divider appears before this facet (after Custom)
			Fields:        getNamesFields(),
			HeaderActions: []string{"add", "export"},
			Actions:       []string{"update"},
		},
		"regular": {
			Name:          "Regular",
			Store:         "names",
			IsForm:        false,
			Fields:        getNamesFields(),
			HeaderActions: []string{"add", "export"},
			Actions:       []string{"update"},
		},
		"baddress": {
			Name:          "Baddress",
			Store:         "names",
			IsForm:        false,
			Fields:        getNamesFields(),
			HeaderActions: []string{"add", "export"},
			Actions:       []string{"update"},
		},
	}

	cfg := &types.ViewConfig{
		ViewName:   "names",
		Facets:     facets,
		FacetOrder: []string{"all", "custom", "prefund", "regular", "baddress"},
		Actions: map[string]types.ActionConfig{
			"add":     {Name: "add", Label: "Add Name", Icon: "Add"},
			"publish": {Name: "publish", Label: "Publish", Icon: "Publish"},
			"pin":     {Name: "pin", Label: "Pin", Icon: "Pin"},
			"export":  {Name: "export", Label: "Export Names", Icon: "Export"},
			"delete":  {Name: "delete", Label: "Delete", Icon: "Delete"},
			"remove":  {Name: "remove", Label: "Remove", Icon: "Remove"},
			"update":  {Name: "update", Label: "Modify", Icon: "Edit"},
		}, // Names uses CRUD actions handled by useActions hook
	}
	types.DeriveFacetFromFields(cfg)
	types.NormalizeOrders(cfg)
	return cfg, nil
}

func getNamesFields() []types.FieldConfig {
	return []types.FieldConfig{
		// Name Identity section
		{Key: "address", Label: "Address", Section: "Name Identity", InTable: true, InDetail: true, Width: 340, Formatter: "address", Order: 1, DetailOrder: 1},
		{Key: "name", Label: "Name", Section: "Name Identity", InTable: true, InDetail: true, Width: 200, Order: 2, DetailOrder: 2},
		{Key: "symbol", Label: "Symbol", Section: "Name Identity", InTable: true, InDetail: true, Width: 100, Order: 5, DetailOrder: 3},
		{Key: "decimals", Label: "Decimals", Section: "Name Identity", InTable: true, InDetail: true, Width: 100, Order: 6, DetailOrder: 4},

		// Classification section
		{Key: "source", Label: "Source", Section: "Classification", InTable: true, InDetail: true, Width: 120, Order: 4, DetailOrder: 5},
		{Key: "tags", Label: "Tags", Section: "Classification", InTable: true, InDetail: true, Width: 150, Order: 3, DetailOrder: 6},
		{Key: "deleted", Label: "Deleted", Section: "Classification", InTable: false, InDetail: true, Formatter: "boolean", DetailOrder: 7},

		// Contract Properties section
		{Key: "isContract", Label: "Is Contract", Section: "Contract Properties", InTable: false, InDetail: true, Formatter: "boolean", DetailOrder: 8},
		{Key: "isCustom", Label: "Is Custom", Section: "Contract Properties", InTable: false, InDetail: true, Formatter: "boolean", DetailOrder: 9},
		{Key: "isErc20", Label: "Is ERC20", Section: "Contract Properties", InTable: false, InDetail: true, Formatter: "boolean", DetailOrder: 10},
		{Key: "isErc721", Label: "Is ERC721", Section: "Contract Properties", InTable: false, InDetail: true, Formatter: "boolean", DetailOrder: 11},
		{Key: "isPrefund", Label: "Is Prefund", Section: "Contract Properties", InTable: false, InDetail: true, Formatter: "boolean", DetailOrder: 12},

		// Prefund Information section
		{Key: "prefund", Label: "Prefund Amount", Section: "Prefund Information", InTable: false, InDetail: true, Formatter: "wei", DetailOrder: 13},
		{Key: "parts", Label: "Parts", Section: "Prefund Information", InTable: false, InDetail: true, DetailOrder: 14},

		// Synthetic actions column for table only
		{Key: "actions", Label: "Actions", Section: "Name Identity", InTable: true, InDetail: false, Width: 80, Order: 7},
	}
}
