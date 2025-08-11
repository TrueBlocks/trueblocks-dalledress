package abis

import (
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

// GetConfig returns the ViewConfig for the Abis view
func (c *AbisCollection) GetConfig() (*types.ViewConfig, error) {
	return &types.ViewConfig{
		ViewName: "abis",
		Facets: map[string]types.FacetConfig{
			"downloaded": {
				Name:         "Downloaded",
				Store:        "downloaded",
				IsForm:       false,
				Columns:      getAbisColumns(),
				DetailPanels: getAbisDetailPanels(),
				Actions:      nil, // No special actions for downloaded facet
			},
			"known": {
				Name:         "Known",
				Store:        "known",
				IsForm:       false,
				Columns:      getAbisColumns(),
				DetailPanels: getAbisDetailPanels(),
				Actions:      nil,
			},
			"functions": {
				Name:         "Functions",
				Store:        "functions",
				IsForm:       false,
				Columns:      getFunctionsColumns(),
				DetailPanels: getFunctionsDetailPanels(),
				Actions:      nil,
			},
			"events": {
				Name:         "Events",
				Store:        "events",
				IsForm:       false,
				Columns:      getFunctionsColumns(), // Events use same structure as functions
				DetailPanels: getFunctionsDetailPanels(),
				Actions:      nil,
			},
		},
		Actions: make(map[string]types.ActionConfig),
	}, nil
}

func getAbisColumns() []types.ColumnConfig {
	return []types.ColumnConfig{
		{Key: "address", Header: "Address", Accessor: "address", Width: 340, Formatter: "address"},
		{Key: "name", Header: "Name", Accessor: "name", Width: 200},
		{Key: "nFunctions", Header: "Functions", Accessor: "nFunctions", Width: 100, Formatter: "number"},
		{Key: "nEvents", Header: "Events", Accessor: "nEvents", Width: 100, Formatter: "number"},
		{Key: "fileSize", Header: "File Size", Accessor: "fileSize", Width: 120, Formatter: "fileSize"},
		{Key: "isEmpty", Header: "Empty", Accessor: "isEmpty", Width: 80, Formatter: "boolean"},
		{Key: "isKnown", Header: "Known", Accessor: "isKnown", Width: 80, Formatter: "boolean"},
		{Key: "lastModDate", Header: "Last Modified", Accessor: "lastModDate", Width: 150, Formatter: "timestamp"},
		{Key: "actions", Header: "Actions", Accessor: "actions", Width: 80},
	}
}

func getFunctionsColumns() []types.ColumnConfig {
	return []types.ColumnConfig{
		{Key: "name", Header: "Name", Accessor: "name", Width: 200},
		{Key: "type", Header: "Type", Accessor: "type", Width: 100},
		{Key: "encoding", Header: "Encoding", Accessor: "encoding", Width: 250},
		{Key: "signature", Header: "Signature", Accessor: "signature", Width: 300},
		{Key: "stateMutability", Header: "State Mutability", Accessor: "stateMutability", Width: 150},
		{Key: "constant", Header: "Constant", Accessor: "constant", Width: 100, Formatter: "boolean"},
		{Key: "anonymous", Header: "Anonymous", Accessor: "anonymous", Width: 100, Formatter: "boolean"},
	}
}

func getAbisDetailPanels() []types.DetailPanelConfig {
	return []types.DetailPanelConfig{
		{
			Title: "ABI Identity",
			Fields: []types.DetailFieldConfig{
				{Key: "address", Label: "Address", Formatter: "address"},
				{Key: "name", Label: "Name"},
				{Key: "path", Label: "Path"},
			},
		},
		{
			Title: "Statistics",
			Fields: []types.DetailFieldConfig{
				{Key: "nFunctions", Label: "Number of Functions", Formatter: "number"},
				{Key: "nEvents", Label: "Number of Events", Formatter: "number"},
				{Key: "fileSize", Label: "File Size", Formatter: "fileSize"},
			},
		},
		{
			Title: "Properties",
			Fields: []types.DetailFieldConfig{
				{Key: "isEmpty", Label: "Is Empty", Formatter: "boolean"},
				{Key: "isKnown", Label: "Is Known", Formatter: "boolean"},
				{Key: "hasConstructor", Label: "Has Constructor", Formatter: "boolean"},
				{Key: "hasFallback", Label: "Has Fallback", Formatter: "boolean"},
			},
		},
		{
			Title:     "Metadata",
			Collapsed: true, // This section starts collapsed
			Fields: []types.DetailFieldConfig{
				{Key: "lastModDate", Label: "Last Modified", Formatter: "timestamp"},
			},
		},
	}
}

func getFunctionsDetailPanels() []types.DetailPanelConfig {
	return []types.DetailPanelConfig{
		{
			Title: "Function Details",
			Fields: []types.DetailFieldConfig{
				{Key: "name", Label: "Function Name"},
				{Key: "type", Label: "Type"},
				{Key: "encoding", Label: "Encoding"},
				{Key: "signature", Label: "Signature"},
			},
		},
		{
			Title: "Function Properties",
			Fields: []types.DetailFieldConfig{
				{Key: "stateMutability", Label: "State Mutability"},
				{Key: "constant", Label: "Constant", Formatter: "boolean"},
				{Key: "anonymous", Label: "Anonymous", Formatter: "boolean"},
			},
		},
		{
			Title:     "Parameters",
			Collapsed: true, // This section starts collapsed
			Fields: []types.DetailFieldConfig{
				{Key: "inputs", Label: "Inputs"},
				{Key: "outputs", Label: "Outputs"},
				{Key: "message", Label: "Message"},
			},
		},
	}
}
