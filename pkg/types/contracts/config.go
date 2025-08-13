package contracts

import (
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

// GetConfig returns the ViewConfig for the Contracts view
func (c *ContractsCollection) GetConfig() (*types.ViewConfig, error) {
	cfg := &types.ViewConfig{
		ViewName: "contracts",
		Facets: map[string]types.FacetConfig{
			"dashboard": {
				Name:          "Dashboard",
				Store:         "contracts",
				IsForm:        true,
				Fields:        getDashboardFields(),
				Actions:       []string{"create", "update", "delete"},
				HeaderActions: []string{"create", "update", "delete"},
			},
			"execute": {
				Name:          "Execute",
				Store:         "contracts",
				IsForm:        true,
				Fields:        getExecuteFields(),
				Actions:       []string{"execute"},
				HeaderActions: []string{"execute"},
			},
			"events": {
				Name:          "Events",
				Store:         "events",
				Fields:        getEventsFields(),
				Actions:       []string{"view", "export"},
				HeaderActions: []string{"view", "export"},
			},
		},
		FacetOrder: []string{"dashboard", "execute", "events"},
		Actions: map[string]types.ActionConfig{
			"create": {
				Name:  "create",
				Label: "Create Contract",
				Icon:  "plus",
			},
			"update": {
				Name:  "update",
				Label: "Update Contract",
				Icon:  "edit",
			},
			"delete": {
				Name:  "delete",
				Label: "Delete Contract",
				Icon:  "trash",
			},
			"execute": {
				Name:  "execute",
				Label: "Execute Function",
				Icon:  "play",
			},
			"view": {
				Name:  "view",
				Label: "View Details",
				Icon:  "eye",
			},
			"export": {
				Name:  "export",
				Label: "Export Data",
				Icon:  "download",
			},
		},
	}
	types.DeriveFacetFromFields(cfg)
	types.NormalizeOrders(cfg)
	return cfg, nil
}

func getDashboardFields() []types.FieldConfig {
	return []types.FieldConfig{
		{Key: "address", Label: "Address", ColumnLabel: "Address", DetailLabel: "Address", Section: "General", InTable: true, InDetail: true, Width: 240, Order: 1, DetailOrder: 1},
		{Key: "name", Label: "Name", ColumnLabel: "Name", DetailLabel: "Name", Section: "General", InTable: true, InDetail: true, Width: 200, Order: 2, DetailOrder: 2},
		{Key: "symbol", Label: "Symbol", ColumnLabel: "Symbol", DetailLabel: "Symbol", Section: "General", InTable: true, InDetail: true, Width: 100, Order: 3, DetailOrder: 3},
		{Key: "decimals", Label: "Decimals", ColumnLabel: "Decimals", DetailLabel: "Decimals", Section: "General", InTable: true, InDetail: true, Width: 100, Order: 4, DetailOrder: 4},
		{Key: "source", Label: "Source", ColumnLabel: "Source", DetailLabel: "Source", Section: "Source", InTable: true, InDetail: true, Width: 150, Order: 5, DetailOrder: 5},
	}
}

func getExecuteFields() []types.FieldConfig {
	return []types.FieldConfig{
		{Key: "address", Label: "Address", ColumnLabel: "Address", DetailLabel: "Address", Section: "General", InTable: true, InDetail: true, Width: 240, Order: 1, DetailOrder: 1},
		{Key: "name", Label: "Name", ColumnLabel: "Name", DetailLabel: "Name", Section: "General", InTable: true, InDetail: true, Width: 200, Order: 2, DetailOrder: 2},
		{Key: "symbol", Label: "Symbol", ColumnLabel: "Symbol", DetailLabel: "Symbol", Section: "General", InTable: true, InDetail: true, Width: 100, Order: 3, DetailOrder: 3},
		{Key: "decimals", Label: "Decimals", ColumnLabel: "Decimals", DetailLabel: "Decimals", Section: "General", InTable: true, InDetail: true, Width: 100, Order: 4, DetailOrder: 4},
		{Key: "source", Label: "Source", ColumnLabel: "Source", DetailLabel: "Source", Section: "Execution", InTable: true, InDetail: true, Width: 150, Order: 5, DetailOrder: 5},
	}
}

func getEventsFields() []types.FieldConfig {
	return []types.FieldConfig{
		{Key: "date", Label: "Date", ColumnLabel: "Date", DetailLabel: "Date", Section: "Block/Tx", InTable: true, InDetail: true, Width: 120, Order: 1, DetailOrder: 1},
		{Key: "address", Label: "Address", ColumnLabel: "Address", DetailLabel: "Address", Section: "Event", InTable: true, InDetail: true, Width: 340, Order: 2, DetailOrder: 5},
		{Key: "name", Label: "Event Name", ColumnLabel: "Name", DetailLabel: "Event Name", Section: "Event", InTable: true, InDetail: true, Width: 200, Order: 3, DetailOrder: 6},
		{Key: "articulatedLog", Label: "Articulated Log", ColumnLabel: "Log Details", DetailLabel: "Articulated Log", Section: "Event", InTable: true, InDetail: true, Width: 120, Order: 4, DetailOrder: 8},
		{Key: "blockNumber", Label: "Block Number", ColumnLabel: "", DetailLabel: "Block Number", Section: "Block/Tx", InTable: false, InDetail: true, DetailOrder: 2},
		{Key: "transactionIndex", Label: "Transaction Index", ColumnLabel: "", DetailLabel: "Transaction Index", Section: "Block/Tx", InTable: false, InDetail: true, DetailOrder: 3},
		{Key: "transactionHash", Label: "Transaction Hash", ColumnLabel: "", DetailLabel: "Transaction Hash", Section: "Block/Tx", InTable: false, InDetail: true, DetailOrder: 4},
		{Key: "signature", Label: "Signature", ColumnLabel: "", DetailLabel: "Signature", Section: "Event", InTable: false, InDetail: true, DetailOrder: 7},
	}
}
