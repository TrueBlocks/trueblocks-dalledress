package contracts

import (
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

// GetConfig returns the ViewConfig for the Contracts view
func (c *ContractsCollection) GetConfig() (*types.ViewConfig, error) {
	return &types.ViewConfig{
		ViewName: "contracts",
		Facets: map[string]types.FacetConfig{
			"dashboard": {
				Name:          "Dashboard",
				Store:         "contracts",
				IsForm:        true,
				Columns:       getDashboardColumns(),
				DetailPanels:  getDashboardDetailPanels(),
				Actions:       []string{"create", "update", "delete"},
				HeaderActions: []string{"create", "update", "delete"},
			},
			"execute": {
				Name:          "Execute",
				Store:         "contracts",
				IsForm:        true,
				Columns:       getExecuteColumns(),
				DetailPanels:  getExecuteDetailPanels(),
				Actions:       []string{"execute"},
				HeaderActions: []string{"execute"},
			},
			"events": {
				Name:          "Events",
				Store:         "events",
				Columns:       getEventsColumns(),
				DetailPanels:  getEventsDetailPanels(),
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
	}, nil
}

func getDashboardColumns() []types.ColumnConfig {
	return []types.ColumnConfig{
		{Key: "address", Header: "Address", Accessor: "address", Width: 240},
		{Key: "name", Header: "Name", Accessor: "name", Width: 200},
		{Key: "symbol", Header: "Symbol", Accessor: "symbol", Width: 100},
		{Key: "decimals", Header: "Decimals", Accessor: "decimals", Width: 100},
		{Key: "source", Header: "Source", Accessor: "source", Width: 150},
	}
}

func getDashboardDetailPanels() []types.DetailPanelConfig {
	return []types.DetailPanelConfig{
		{
			Title: "General",
			Fields: []types.DetailFieldConfig{
				{Key: "address", Label: "Address"},
				{Key: "name", Label: "Name"},
				{Key: "symbol", Label: "Symbol"},
				{Key: "decimals", Label: "Decimals"},
			},
		},
		{
			Title: "Source",
			Fields: []types.DetailFieldConfig{
				{Key: "source", Label: "Source"},
			},
		},
	}
}

func getExecuteColumns() []types.ColumnConfig {
	return []types.ColumnConfig{
		{Key: "address", Header: "Address", Accessor: "address", Width: 240},
		{Key: "name", Header: "Name", Accessor: "name", Width: 200},
		{Key: "symbol", Header: "Symbol", Accessor: "symbol", Width: 100},
		{Key: "decimals", Header: "Decimals", Accessor: "decimals", Width: 100},
		{Key: "source", Header: "Source", Accessor: "source", Width: 150},
	}
}

func getExecuteDetailPanels() []types.DetailPanelConfig {
	return []types.DetailPanelConfig{
		{
			Title: "General",
			Fields: []types.DetailFieldConfig{
				{Key: "address", Label: "Address"},
				{Key: "name", Label: "Name"},
				{Key: "symbol", Label: "Symbol"},
				{Key: "decimals", Label: "Decimals"},
			},
		},
		{
			Title: "Execution",
			Fields: []types.DetailFieldConfig{
				{Key: "source", Label: "Source"},
			},
		},
	}
}

func getEventsColumns() []types.ColumnConfig {
	return []types.ColumnConfig{
		{Key: "date", Header: "Date", Accessor: "date", Width: 120},
		{Key: "address", Header: "Address", Accessor: "address", Width: 340},
		{Key: "name", Header: "Name", Accessor: "name", Width: 200},
		{Key: "articulatedLog", Header: "Log Details", Accessor: "articulatedLog", Width: 120},
	}
}

func getEventsDetailPanels() []types.DetailPanelConfig {
	return []types.DetailPanelConfig{
		{
			Title: "Block/Tx",
			Fields: []types.DetailFieldConfig{
				{Key: "date", Label: "Date"},
				{Key: "blockNumber", Label: "Block Number"},
				{Key: "transactionIndex", Label: "Transaction Index"},
				{Key: "transactionHash", Label: "Transaction Hash"},
			},
		},
		{
			Title: "Event",
			Fields: []types.DetailFieldConfig{
				{Key: "address", Label: "Address"},
				{Key: "name", Label: "Event Name"},
				{Key: "signature", Label: "Signature"},
				{Key: "articulatedLog", Label: "Articulated Log"},
			},
		},
	}
}
