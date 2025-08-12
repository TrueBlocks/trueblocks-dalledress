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
				Name:          "Downloaded",
				Store:         "downloaded",
				IsForm:        false,
				Columns:       getAbisColumns(),
				DetailPanels:  getAbisDetailPanels(),
				Actions:       []string{"remove"},
				HeaderActions: []string{"export"},
			},
			"known": {
				Name:          "Known",
				Store:         "known",
				IsForm:        false,
				Columns:       getAbisColumns(),
				DetailPanels:  getAbisDetailPanels(),
				Actions:       nil,
				HeaderActions: []string{"export"},
			},
			"functions": {
				Name:          "Functions",
				Store:         "functions",
				IsForm:        false,
				Columns:       getFunctionsColumns(),
				DetailPanels:  getFunctionsDetailPanels(),
				Actions:       nil,
				HeaderActions: []string{"export"},
			},
			"events": {
				Name:          "Events",
				Store:         "events",
				IsForm:        false,
				Columns:       getFunctionsColumns(), // Events use same structure as functions
				DetailPanels:  getFunctionsDetailPanels(),
				Actions:       nil,
				HeaderActions: []string{"export"},
			},
		},
		FacetOrder: []string{"downloaded", "known", "functions", "events"},
		Actions: map[string]types.ActionConfig{
			"export": {Name: "export", Label: "Export Data", Icon: "Export"},
			"remove": {Name: "remove", Label: "Remove", Icon: "Remove"},
		},
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
			Title:     "ABI Identity",
			Collapsed: false,
			Fields: []types.DetailFieldConfig{
				{Key: "address", Label: "Address", Formatter: "address"},
				{Key: "name", Label: "Name"},
				{Key: "path", Label: "File Path"},
			},
		},
		{
			Title:     "Content Statistics",
			Collapsed: false,
			Fields: []types.DetailFieldConfig{
				{Key: "nFunctions", Label: "Number of Functions", Formatter: "number"},
				{Key: "nEvents", Label: "Number of Events", Formatter: "number"},
				{Key: "fileSize", Label: "File Size", Formatter: "fileSize"},
			},
		},
		{
			Title:     "ABI Properties",
			Collapsed: false,
			Fields: []types.DetailFieldConfig{
				{Key: "isEmpty", Label: "Is Empty", Formatter: "boolean"},
				{Key: "isKnown", Label: "Is Known", Formatter: "boolean"},
				{Key: "hasConstructor", Label: "Has Constructor", Formatter: "boolean"},
				{Key: "hasFallback", Label: "Has Fallback Function", Formatter: "boolean"},
			},
		},
		{
			Title:     "File Metadata",
			Collapsed: false,
			Fields: []types.DetailFieldConfig{
				{Key: "lastModDate", Label: "Last Modified", Formatter: "timestamp"},
			},
		},
		{
			Title:     "Functions List",
			Collapsed: false,
			Fields: []types.DetailFieldConfig{
				{Key: "functions", Label: "Available Functions", Formatter: "json"},
			},
		},
	}
}

func getFunctionsDetailPanels() []types.DetailPanelConfig {
	return []types.DetailPanelConfig{
		{
			Title:     "Function Overview",
			Collapsed: false,
			Fields: []types.DetailFieldConfig{
				{Key: "name", Label: "Function Name"},
				{Key: "type", Label: "Function Type"},
				{Key: "encoding", Label: "Encoding Hash"},
				{Key: "signature", Label: "Function Signature"},
			},
		},
		{
			Title:     "Function Properties",
			Collapsed: false,
			Fields: []types.DetailFieldConfig{
				{Key: "stateMutability", Label: "State Mutability"},
				{Key: "constant", Label: "Constant", Formatter: "boolean"},
				{Key: "anonymous", Label: "Anonymous", Formatter: "boolean"},
			},
		},
		{
			Title:     "Input Parameters",
			Collapsed: false,
			Fields: []types.DetailFieldConfig{
				{Key: "inputs", Label: "Input Parameters", Formatter: "json"},
			},
		},
		{
			Title:     "Output Parameters",
			Collapsed: false,
			Fields: []types.DetailFieldConfig{
				{Key: "outputs", Label: "Output Parameters", Formatter: "json"},
			},
		},
		{
			Title:     "Additional Information",
			Collapsed: false,
			Fields: []types.DetailFieldConfig{
				{Key: "message", Label: "Error Message"},
			},
		},
	}
}
