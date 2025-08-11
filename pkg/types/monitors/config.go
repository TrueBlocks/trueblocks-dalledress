package monitors

import (
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

// GetConfig returns the ViewConfig for the Monitors view
func (c *MonitorsCollection) GetConfig() (*types.ViewConfig, error) {
	facets := map[string]types.FacetConfig{
		"monitors": {
			Name:         "Monitors",
			Store:        "monitors",
			Columns:      getColumns(),
			DetailPanels: getDetailPanels(),
			Actions:      []string{"delete", "remove"},
		},
	}

	return &types.ViewConfig{
		ViewName: "monitors",
		Facets:   facets,
		Actions: map[string]types.ActionConfig{
			"delete": {
				Name:         "delete",
				Label:        "Delete",
				Icon:         "delete",
				Confirmation: true,
				Facets:       []string{"monitors"},
			},
			"remove": {
				Name:         "remove",
				Label:        "Remove",
				Icon:         "remove",
				Confirmation: true,
				Facets:       []string{"monitors"},
			},
		},
	}, nil
}

func getColumns() []types.ColumnConfig {
	return []types.ColumnConfig{
		{
			Key:        "address",
			Header:     "Address",
			Accessor:   "address",
			Width:      340,
			Sortable:   true,
			Filterable: true,
			Formatter:  "address",
		},
		{
			Key:        "name",
			Header:     "Name",
			Accessor:   "name",
			Width:      200,
			Sortable:   true,
			Filterable: true,
			Formatter:  "",
		},
		{
			Key:        "nRecords",
			Header:     "Records",
			Accessor:   "nRecords",
			Width:      120,
			Sortable:   true,
			Filterable: false,
			Formatter:  "number",
		},
		{
			Key:        "fileSize",
			Header:     "File Size",
			Accessor:   "fileSize",
			Width:      120,
			Sortable:   true,
			Filterable: false,
			Formatter:  "fileSize",
		},
		{
			Key:        "isEmpty",
			Header:     "Empty",
			Accessor:   "isEmpty",
			Width:      80,
			Sortable:   true,
			Filterable: false,
			Formatter:  "boolean",
		},
		{
			Key:        "lastScanned",
			Header:     "Last Scanned",
			Accessor:   "lastScanned",
			Width:      120,
			Sortable:   true,
			Filterable: false,
			Formatter:  "timestamp",
		},
		{
			Key:        "actions",
			Header:     "Actions",
			Accessor:   "actions",
			Width:      80,
			Sortable:   false,
			Filterable: false,
			Formatter:  "actions",
		},
	}
}

func getDetailPanels() []types.DetailPanelConfig {
	// Static panel configuration

	return []types.DetailPanelConfig{
		{
			Title: "Overview",
			Fields: []types.DetailFieldConfig{
				{
					Key:       "name",
					Label:     "Name",
					Formatter: "",
				},
				{
					Key:       "address",
					Label:     "Address",
					Formatter: "address",
				},
			},
		},
		{
			Title: "Statistics",
			Fields: []types.DetailFieldConfig{
				{
					Key:       "nRecords",
					Label:     "Records",
					Formatter: "number",
				},
				{
					Key:       "fileSize",
					Label:     "File Size",
					Formatter: "fileSize",
				},
				{
					Key:       "isEmpty",
					Label:     "Empty",
					Formatter: "boolean",
				},
			},
		},
		{
			Title: "Scanning",
			Fields: []types.DetailFieldConfig{
				{
					Key:       "lastScanned",
					Label:     "Last Scanned",
					Formatter: "timestamp",
				},
				{
					Key:       "lastScanAge",
					Label:     "Age",
					Formatter: "computed",
				},
			},
		},
	}
}
