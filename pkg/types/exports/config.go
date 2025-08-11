package exports

import (
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

// GetConfig returns the ViewConfig for the Exports view
func (c *ExportsCollection) GetConfig() (*types.ViewConfig, error) {
	// Create facet configurations - Exports has 9 different facets with unique columns for each
	facets := map[string]types.FacetConfig{
		"statements": {
			Name:         "Statements",
			Store:        "exports",
			IsForm:       false,
			Columns:      getStatementsColumns(),
			DetailPanels: getStatementsDetailPanels(),
		},
		"balances": {
			Name:         "Balances",
			Store:        "exports",
			IsForm:       false,
			Columns:      getBalancesColumns(),
			DetailPanels: getBalancesDetailPanels(),
		},
		"transfers": {
			Name:         "Transfers",
			Store:        "exports",
			IsForm:       false,
			Columns:      getTransfersColumns(),
			DetailPanels: getTransfersDetailPanels(),
		},
		"transactions": {
			Name:         "Transactions",
			Store:        "exports",
			IsForm:       false,
			Columns:      getTransactionsColumns(),
			DetailPanels: getTransactionsDetailPanels(),
		},
		"withdrawals": {
			Name:         "Withdrawals",
			Store:        "exports",
			IsForm:       false,
			Columns:      getWithdrawalsColumns(),
			DetailPanels: getWithdrawalsDetailPanels(),
		},
		"assets": {
			Name:         "Assets",
			Store:        "exports",
			IsForm:       false,
			Columns:      getAssetsColumns(),
			DetailPanels: getAssetsDetailPanels(),
		},
		"logs": {
			Name:         "Logs",
			Store:        "exports",
			IsForm:       false,
			Columns:      getLogsColumns(),
			DetailPanels: getLogsDetailPanels(),
		},
		"traces": {
			Name:         "Traces",
			Store:        "exports",
			IsForm:       false,
			Columns:      getTracesColumns(),
			DetailPanels: getTracesDetailPanels(),
		},
		"receipts": {
			Name:         "Receipts",
			Store:        "exports",
			IsForm:       false,
			Columns:      getReceiptsColumns(),
			DetailPanels: getReceiptsDetailPanels(),
		},
	}

	return &types.ViewConfig{
		ViewName: "exports",
		Facets:   facets,
		Actions:  make(map[string]types.ActionConfig), // Exports uses export actions handled by ExportFormatModal
	}, nil
}

// Helper functions for different export types
// Each facet has different column structures based on the data type

func getStatementsColumns() []types.ColumnConfig {
	return []types.ColumnConfig{
		{Key: "date", Header: "Date", Accessor: "date", Width: 120, Formatter: "datetime"},
		{Key: "asset", Header: "Asset", Accessor: "asset", Width: 340, Formatter: "address"},
		{Key: "symbol", Header: "Symbol", Accessor: "symbol", Width: 200},
		{Key: "decimals", Header: "Decimals", Accessor: "decimals", Width: 100},
		{Key: "begBal", Header: "Begin Balance", Accessor: "begBal", Width: 150, Formatter: "wei"},
		{Key: "amountIn", Header: "In", Accessor: "amountIn", Width: 150, Formatter: "wei"},
		{Key: "amountOut", Header: "Out", Accessor: "amountOut", Width: 150, Formatter: "wei"},
		{Key: "endBal", Header: "End Balance", Accessor: "endBal", Width: 150, Formatter: "wei"},
		{Key: "gasUsed", Header: "Gas Used", Accessor: "gasUsed", Width: 120, Formatter: "gas"},
		{Key: "reconciliationType", Header: "Type", Accessor: "reconciliationType", Width: 100},
		{Key: "actions", Header: "Actions", Accessor: "actions", Width: 80},
	}
}

func getStatementsDetailPanels() []types.DetailPanelConfig {
	return []types.DetailPanelConfig{
		{
			Title: "Statement",
			Fields: []types.DetailFieldConfig{
				{Key: "date", Label: "Timestamp", Formatter: "datetime"},
				{Key: "begBal", Label: "Begin Balance", Formatter: "wei"},
				{Key: "amountIn", Label: "In", Formatter: "wei"},
				{Key: "amountOut", Label: "Out", Formatter: "wei"},
				{Key: "endBal", Label: "End Balance", Formatter: "wei"},
			},
		},
		{
			Title: "Asset Info",
			Fields: []types.DetailFieldConfig{
				{Key: "asset", Label: "Asset", Formatter: "address"},
				{Key: "symbol", Label: "Symbol"},
				{Key: "decimals", Label: "Decimals"},
			},
		},
		{
			Title: "Transaction",
			Fields: []types.DetailFieldConfig{
				{Key: "gasUsed", Label: "Gas Used", Formatter: "gas"},
				{Key: "reconciliationType", Label: "Type"},
			},
		},
	}
}

func getBalancesColumns() []types.ColumnConfig {
	return []types.ColumnConfig{
		{Key: "date", Header: "Date", Accessor: "date", Width: 120, Formatter: "datetime"},
		{Key: "holder", Header: "Holder", Accessor: "holder", Width: 340, Formatter: "address"},
		{Key: "address", Header: "Address", Accessor: "address", Width: 340, Formatter: "address"},
		{Key: "symbol", Header: "Symbol", Accessor: "symbol", Width: 100},
		{Key: "balance", Header: "Balance", Accessor: "balance", Width: 150, Formatter: "wei"},
		{Key: "decimals", Header: "Decimals", Accessor: "decimals", Width: 100},
		{Key: "actions", Header: "Actions", Accessor: "actions", Width: 80},
	}
}

func getBalancesDetailPanels() []types.DetailPanelConfig {
	return []types.DetailPanelConfig{
		{
			Title: "Balance",
			Fields: []types.DetailFieldConfig{
				{Key: "date", Label: "Date", Formatter: "datetime"},
				{Key: "balance", Label: "Balance", Formatter: "wei"},
			},
		},
		{
			Title: "Addresses",
			Fields: []types.DetailFieldConfig{
				{Key: "holder", Label: "Holder", Formatter: "address"},
				{Key: "address", Label: "Address", Formatter: "address"},
			},
		},
		{
			Title: "Token Info",
			Fields: []types.DetailFieldConfig{
				{Key: "symbol", Label: "Symbol"},
				{Key: "decimals", Label: "Decimals"},
			},
		},
	}
}

func getTransfersColumns() []types.ColumnConfig {
	return []types.ColumnConfig{
		{Key: "blockNumber", Header: "Block", Accessor: "blockNumber", Width: 100},
		{Key: "transactionIndex", Header: "Tx Index", Accessor: "transactionIndex", Width: 80},
		{Key: "logIndex", Header: "Log Index", Accessor: "logIndex", Width: 80},
		{Key: "from", Header: "From", Accessor: "from", Width: 340, Formatter: "address"},
		{Key: "to", Header: "To", Accessor: "to", Width: 340, Formatter: "address"},
		{Key: "asset", Header: "Asset", Accessor: "asset", Width: 340, Formatter: "address"},
		{Key: "amount", Header: "Amount", Accessor: "amount", Width: 150, Formatter: "wei"},
		{Key: "actions", Header: "Actions", Accessor: "actions", Width: 80},
	}
}

func getTransfersDetailPanels() []types.DetailPanelConfig {
	return []types.DetailPanelConfig{
		{
			Title: "Transfer",
			Fields: []types.DetailFieldConfig{
				{Key: "from", Label: "From", Formatter: "address"},
				{Key: "to", Label: "To", Formatter: "address"},
				{Key: "amount", Label: "Amount", Formatter: "wei"},
				{Key: "asset", Label: "Asset", Formatter: "address"},
			},
		},
		{
			Title: "Transaction Info",
			Fields: []types.DetailFieldConfig{
				{Key: "blockNumber", Label: "Block"},
				{Key: "transactionIndex", Label: "Tx Index"},
				{Key: "logIndex", Label: "Log Index"},
			},
		},
	}
}

func getTransactionsColumns() []types.ColumnConfig {
	return []types.ColumnConfig{
		{Key: "blockNumber", Header: "Block", Accessor: "blockNumber", Width: 100},
		{Key: "transactionIndex", Header: "Index", Accessor: "transactionIndex", Width: 80},
		{Key: "hash", Header: "Hash", Accessor: "hash", Width: 340, Formatter: "hash"},
		{Key: "from", Header: "From", Accessor: "from", Width: 340, Formatter: "address"},
		{Key: "to", Header: "To", Accessor: "to", Width: 340, Formatter: "address"},
		{Key: "value", Header: "Value", Accessor: "value", Width: 150, Formatter: "wei"},
		{Key: "gasUsed", Header: "Gas Used", Accessor: "gasUsed", Width: 120, Formatter: "gas"},
		{Key: "actions", Header: "Actions", Accessor: "actions", Width: 80},
	}
}

func getTransactionsDetailPanels() []types.DetailPanelConfig {
	return []types.DetailPanelConfig{
		{
			Title: "Transaction",
			Fields: []types.DetailFieldConfig{
				{Key: "hash", Label: "Hash", Formatter: "hash"},
				{Key: "from", Label: "From", Formatter: "address"},
				{Key: "to", Label: "To", Formatter: "address"},
				{Key: "value", Label: "Value", Formatter: "wei"},
			},
		},
		{
			Title: "Block Info",
			Fields: []types.DetailFieldConfig{
				{Key: "blockNumber", Label: "Block"},
				{Key: "transactionIndex", Label: "Index"},
				{Key: "gasUsed", Label: "Gas Used", Formatter: "gas"},
			},
		},
	}
}

func getWithdrawalsColumns() []types.ColumnConfig {
	return []types.ColumnConfig{
		{Key: "blockNumber", Header: "Block", Accessor: "blockNumber", Width: 100},
		{Key: "index", Header: "Index", Accessor: "index", Width: 80},
		{Key: "validatorIndex", Header: "Validator", Accessor: "validatorIndex", Width: 100},
		{Key: "address", Header: "Address", Accessor: "address", Width: 340, Formatter: "address"},
		{Key: "amount", Header: "Amount", Accessor: "amount", Width: 150, Formatter: "wei"},
		{Key: "actions", Header: "Actions", Accessor: "actions", Width: 80},
	}
}

func getWithdrawalsDetailPanels() []types.DetailPanelConfig {
	return []types.DetailPanelConfig{
		{
			Title: "Withdrawal",
			Fields: []types.DetailFieldConfig{
				{Key: "address", Label: "Address", Formatter: "address"},
				{Key: "amount", Label: "Amount", Formatter: "wei"},
				{Key: "validatorIndex", Label: "Validator"},
			},
		},
		{
			Title: "Block Info",
			Fields: []types.DetailFieldConfig{
				{Key: "blockNumber", Label: "Block"},
				{Key: "index", Label: "Index"},
			},
		},
	}
}

func getAssetsColumns() []types.ColumnConfig {
	return []types.ColumnConfig{
		{Key: "address", Header: "Address", Accessor: "address", Width: 340, Formatter: "address"},
		{Key: "holder", Header: "Holder", Accessor: "holder", Width: 340, Formatter: "address"},
		{Key: "symbol", Header: "Symbol", Accessor: "symbol", Width: 100},
		{Key: "name", Header: "Name", Accessor: "name", Width: 200},
		{Key: "decimals", Header: "Decimals", Accessor: "decimals", Width: 100},
		{Key: "totalSupply", Header: "Total Supply", Accessor: "totalSupply", Width: 150, Formatter: "wei"},
		{Key: "actions", Header: "Actions", Accessor: "actions", Width: 80},
	}
}

func getAssetsDetailPanels() []types.DetailPanelConfig {
	return []types.DetailPanelConfig{
		{
			Title: "Asset",
			Fields: []types.DetailFieldConfig{
				{Key: "address", Label: "Address", Formatter: "address"},
				{Key: "holder", Label: "Holder", Formatter: "address"},
				{Key: "name", Label: "Name"},
				{Key: "symbol", Label: "Symbol"},
			},
		},
		{
			Title: "Token Info",
			Fields: []types.DetailFieldConfig{
				{Key: "decimals", Label: "Decimals"},
				{Key: "totalSupply", Label: "Total Supply", Formatter: "wei"},
			},
		},
	}
}

func getLogsColumns() []types.ColumnConfig {
	return []types.ColumnConfig{
		{Key: "blockNumber", Header: "Block", Accessor: "blockNumber", Width: 100},
		{Key: "transactionIndex", Header: "Tx Index", Accessor: "transactionIndex", Width: 80},
		{Key: "logIndex", Header: "Log Index", Accessor: "logIndex", Width: 80},
		{Key: "address", Header: "Address", Accessor: "address", Width: 340, Formatter: "address"},
		{Key: "topic0", Header: "Topic0", Accessor: "topic0", Width: 340, Formatter: "hash"},
		{Key: "topic1", Header: "Topic1", Accessor: "topic1", Width: 340, Formatter: "hash"},
		{Key: "actions", Header: "Actions", Accessor: "actions", Width: 80},
	}
}

func getLogsDetailPanels() []types.DetailPanelConfig {
	return []types.DetailPanelConfig{
		{
			Title: "Log",
			Fields: []types.DetailFieldConfig{
				{Key: "address", Label: "Address", Formatter: "address"},
				{Key: "topic0", Label: "Topic0", Formatter: "hash"},
				{Key: "topic1", Label: "Topic1", Formatter: "hash"},
				{Key: "topic2", Label: "Topic2", Formatter: "hash"},
				{Key: "topic3", Label: "Topic3", Formatter: "hash"},
			},
		},
		{
			Title: "Transaction Info",
			Fields: []types.DetailFieldConfig{
				{Key: "blockNumber", Label: "Block"},
				{Key: "transactionIndex", Label: "Tx Index"},
				{Key: "logIndex", Label: "Log Index"},
			},
		},
	}
}

func getTracesColumns() []types.ColumnConfig {
	return []types.ColumnConfig{
		{Key: "blockNumber", Header: "Block", Accessor: "blockNumber", Width: 100},
		{Key: "transactionIndex", Header: "Tx Index", Accessor: "transactionIndex", Width: 80},
		{Key: "traceIndex", Header: "Trace Index", Accessor: "traceIndex", Width: 80},
		{Key: "from", Header: "From", Accessor: "from", Width: 340, Formatter: "address"},
		{Key: "to", Header: "To", Accessor: "to", Width: 340, Formatter: "address"},
		{Key: "value", Header: "Value", Accessor: "value", Width: 150, Formatter: "wei"},
		{Key: "type", Header: "Type", Accessor: "type", Width: 100},
		{Key: "actions", Header: "Actions", Accessor: "actions", Width: 80},
	}
}

func getTracesDetailPanels() []types.DetailPanelConfig {
	return []types.DetailPanelConfig{
		{
			Title: "Trace",
			Fields: []types.DetailFieldConfig{
				{Key: "from", Label: "From", Formatter: "address"},
				{Key: "to", Label: "To", Formatter: "address"},
				{Key: "value", Label: "Value", Formatter: "wei"},
				{Key: "type", Label: "Type"},
			},
		},
		{
			Title: "Transaction Info",
			Fields: []types.DetailFieldConfig{
				{Key: "blockNumber", Label: "Block"},
				{Key: "transactionIndex", Label: "Tx Index"},
				{Key: "traceIndex", Label: "Trace Index"},
			},
		},
	}
}

func getReceiptsColumns() []types.ColumnConfig {
	return []types.ColumnConfig{
		{Key: "blockNumber", Header: "Block", Accessor: "blockNumber", Width: 100},
		{Key: "transactionIndex", Header: "Index", Accessor: "transactionIndex", Width: 80},
		{Key: "transactionHash", Header: "Tx Hash", Accessor: "transactionHash", Width: 340, Formatter: "hash"},
		{Key: "from", Header: "From", Accessor: "from", Width: 340, Formatter: "address"},
		{Key: "to", Header: "To", Accessor: "to", Width: 340, Formatter: "address"},
		{Key: "gasUsed", Header: "Gas Used", Accessor: "gasUsed", Width: 120, Formatter: "gas"},
		{Key: "status", Header: "Status", Accessor: "status", Width: 80},
		{Key: "actions", Header: "Actions", Accessor: "actions", Width: 80},
	}
}

func getReceiptsDetailPanels() []types.DetailPanelConfig {
	return []types.DetailPanelConfig{
		{
			Title: "Receipt",
			Fields: []types.DetailFieldConfig{
				{Key: "transactionHash", Label: "Tx Hash", Formatter: "hash"},
				{Key: "from", Label: "From", Formatter: "address"},
				{Key: "to", Label: "To", Formatter: "address"},
				{Key: "status", Label: "Status"},
			},
		},
		{
			Title: "Gas Info",
			Fields: []types.DetailFieldConfig{
				{Key: "gasUsed", Label: "Gas Used", Formatter: "gas"},
				{Key: "blockNumber", Label: "Block"},
				{Key: "transactionIndex", Label: "Index"},
			},
		},
	}
}
