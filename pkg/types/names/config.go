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
			Columns:       getNamesColumns(),
			DetailPanels:  getNamesDetailPanels(),
			HeaderActions: []string{"add", "export"},
			Actions:       []string{"update"},
		},
		"custom": {
			Name:          "Custom",
			Store:         "names",
			IsForm:        false,
			Columns:       getNamesColumns(),
			DetailPanels:  getNamesDetailPanels(),
			HeaderActions: []string{"add", "publish", "pin", "export"},
			Actions:       []string{"delete", "remove", "update"},
		},
		"prefund": {
			Name:          "Prefund",
			Store:         "names",
			IsForm:        false,
			DividerBefore: true, // Divider appears before this facet (after Custom)
			Columns:       getNamesColumns(),
			DetailPanels:  getNamesDetailPanels(),
			HeaderActions: []string{"add", "export"},
			Actions:       []string{"update"},
		},
		"regular": {
			Name:          "Regular",
			Store:         "names",
			IsForm:        false,
			Columns:       getNamesColumns(),
			DetailPanels:  getNamesDetailPanels(),
			HeaderActions: []string{"add", "export"},
			Actions:       []string{"update"},
		},
		"baddress": {
			Name:          "Baddress",
			Store:         "names",
			IsForm:        false,
			Columns:       getNamesColumns(),
			DetailPanels:  getNamesDetailPanels(),
			HeaderActions: []string{"add", "export"},
			Actions:       []string{"update"},
		},
	}

	return &types.ViewConfig{
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
	}, nil
}

func getNamesColumns() []types.ColumnConfig {
	return []types.ColumnConfig{
		{Key: "address", Header: "Address", Accessor: "address", Width: 340, Formatter: "address"},
		{Key: "name", Header: "Name", Accessor: "name", Width: 200},
		{Key: "tags", Header: "Tags", Accessor: "tags", Width: 150},
		{Key: "source", Header: "Source", Accessor: "source", Width: 120},
		{Key: "symbol", Header: "Symbol", Accessor: "symbol", Width: 100},
		{Key: "decimals", Header: "Decimals", Accessor: "decimals", Width: 100},
		{Key: "actions", Header: "Actions", Accessor: "actions", Width: 80},
	}
}

func getNamesDetailPanels() []types.DetailPanelConfig {
	return []types.DetailPanelConfig{
		{
			Title:     "Name Identity",
			Collapsed: false,
			Fields: []types.DetailFieldConfig{
				{Key: "address", Label: "Address", Formatter: "address"},
				{Key: "name", Label: "Name"},
				{Key: "symbol", Label: "Symbol"},
				{Key: "decimals", Label: "Decimals"},
			},
		},
		{
			Title:     "Classification",
			Collapsed: false,
			Fields: []types.DetailFieldConfig{
				{Key: "source", Label: "Source"},
				{Key: "tags", Label: "Tags"},
				{Key: "deleted", Label: "Deleted", Formatter: "boolean"},
			},
		},
		{
			Title:     "Contract Properties",
			Collapsed: false,
			Fields: []types.DetailFieldConfig{
				{Key: "isContract", Label: "Is Contract", Formatter: "boolean"},
				{Key: "isCustom", Label: "Is Custom", Formatter: "boolean"},
				{Key: "isErc20", Label: "Is ERC20", Formatter: "boolean"},
				{Key: "isErc721", Label: "Is ERC721", Formatter: "boolean"},
				{Key: "isPrefund", Label: "Is Prefund", Formatter: "boolean"},
			},
		},
		{
			Title:     "Prefund Information",
			Collapsed: false,
			Fields: []types.DetailFieldConfig{
				{Key: "prefund", Label: "Prefund Amount", Formatter: "wei"},
				{Key: "parts", Label: "Parts"},
			},
		},
	}
}
