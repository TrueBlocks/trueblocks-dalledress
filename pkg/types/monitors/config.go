package monitors

import (
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

// GetConfig returns the ViewConfig for the Monitors view
func (c *MonitorsCollection) GetConfig() (*types.ViewConfig, error) {
	facets := map[string]types.FacetConfig{
		"monitors": {
			Name:          "Monitors",
			Store:         "monitors",
			Fields:        getFields(),
			Actions:       []string{"delete", "remove"},
			HeaderActions: []string{"export"},
		},
	}

	cfg := &types.ViewConfig{
		ViewName:   "monitors",
		Facets:     facets,
		FacetOrder: []string{"monitors"},
		Actions: map[string]types.ActionConfig{
			"export": {
				Name:  "export",
				Label: "Export Data",
				Icon:  "Export",
			},
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
	}
	types.DeriveFacetFromFields(cfg)
	types.NormalizeOrders(cfg)
	return cfg, nil
}

func getFields() []types.FieldConfig {
	return []types.FieldConfig{
		{Key: "address", Label: "Address", ColumnLabel: "Address", DetailLabel: "Address", Formatter: "address", Section: "Monitor Overview", InTable: true, InDetail: true, Width: 340, Sortable: true, Filterable: true, Order: 1, DetailOrder: 2},
		{Key: "name", Label: "Name", ColumnLabel: "Name", DetailLabel: "Name", Section: "Monitor Overview", InTable: true, InDetail: true, Width: 200, Sortable: true, Filterable: true, Order: 2, DetailOrder: 1},
		{Key: "deleted", Label: "Deleted", ColumnLabel: "Deleted", DetailLabel: "Deleted", Formatter: "boolean", Section: "Monitor Overview", InTable: false, InDetail: true, DetailOrder: 3},
		{Key: "isStaged", Label: "Staged", ColumnLabel: "Staged", DetailLabel: "Staged", Formatter: "boolean", Section: "Monitor Overview", InTable: false, InDetail: true, DetailOrder: 4},

		{Key: "nRecords", Label: "Records", ColumnLabel: "Records", DetailLabel: "Total Records", Formatter: "number", Section: "File Statistics", InTable: true, InDetail: true, Width: 120, Sortable: true, Filterable: false, Order: 3, DetailOrder: 5},
		{Key: "fileSize", Label: "File Size", ColumnLabel: "File Size", DetailLabel: "File Size", Formatter: "fileSize", Section: "File Statistics", InTable: true, InDetail: true, Width: 120, Sortable: true, Filterable: false, Order: 4, DetailOrder: 6},
		{Key: "isEmpty", Label: "Empty", ColumnLabel: "Empty", DetailLabel: "Is Empty", Formatter: "boolean", Section: "File Statistics", InTable: true, InDetail: true, Width: 80, Sortable: true, Filterable: false, Order: 5, DetailOrder: 7},

		{Key: "lastScanned", Label: "Last Scanned", ColumnLabel: "Last Scanned", DetailLabel: "Last Scanned", Formatter: "timestamp", Section: "Scanning Information", InTable: true, InDetail: true, Width: 120, Sortable: true, Filterable: false, Order: 6, DetailOrder: 8},

		{Key: "actions", Label: "Actions", ColumnLabel: "Actions", DetailLabel: "Actions", Section: "", InTable: true, InDetail: false, Width: 80, Sortable: false, Filterable: false, Order: 7},
	}
}
