package exports

import (
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

// GetConfig returns the ViewConfig for the Exports view
func (c *ExportsCollection) GetConfig() (*types.ViewConfig, error) {
	// Create facet configurations - Exports has 9 different facets with unique columns for each
	facets := map[string]types.FacetConfig{
		"statements": {
			Name:          "Statements",
			Store:         "exports",
			IsForm:        false,
			Columns:       getStatementsColumns(),
			DetailPanels:  getStatementsDetailPanels(),
			HeaderActions: []string{"export"},
		},
		"balances": {
			Name:          "Balances",
			Store:         "exports",
			IsForm:        false,
			Columns:       getBalancesColumns(),
			DetailPanels:  getBalancesDetailPanels(),
			HeaderActions: []string{"export"},
		},
		"transfers": {
			Name:          "Transfers",
			Store:         "exports",
			IsForm:        false,
			Columns:       getTransfersColumns(),
			DetailPanels:  getTransfersDetailPanels(),
			HeaderActions: []string{"export"},
		},
		"transactions": {
			Name:          "Transactions",
			Store:         "exports",
			IsForm:        false,
			Columns:       getTransactionsColumns(),
			DetailPanels:  getTransactionsDetailPanels(),
			HeaderActions: []string{"export"},
		},
		"withdrawals": {
			Name:          "Withdrawals",
			Store:         "exports",
			IsForm:        false,
			Columns:       getWithdrawalsColumns(),
			DetailPanels:  getWithdrawalsDetailPanels(),
			HeaderActions: []string{"export"},
		},
		"assets": {
			Name:          "Assets",
			Store:         "exports",
			IsForm:        false,
			Columns:       getAssetsColumns(),
			DetailPanels:  getAssetsDetailPanels(),
			HeaderActions: []string{"export"},
		},
		"logs": {
			Name:          "Logs",
			Store:         "exports",
			IsForm:        false,
			Columns:       getLogsColumns(),
			DetailPanels:  getLogsDetailPanels(),
			HeaderActions: []string{"export"},
		},
		"traces": {
			Name:          "Traces",
			Store:         "exports",
			IsForm:        false,
			Columns:       getTracesColumns(),
			DetailPanels:  getTracesDetailPanels(),
			HeaderActions: []string{"export"},
		},
		"receipts": {
			Name:          "Receipts",
			Store:         "exports",
			IsForm:        false,
			Columns:       getReceiptsColumns(),
			DetailPanels:  getReceiptsDetailPanels(),
			HeaderActions: []string{"export"},
		},
	}

	return &types.ViewConfig{
		ViewName:   "exports",
		Facets:     facets,
		FacetOrder: []string{"statements", "balances", "transfers", "transactions", "withdrawals", "assets", "logs", "traces", "receipts"},
		Actions: map[string]types.ActionConfig{
			"export": {Name: "export", Label: "Export Data", Icon: "Export"},
		}, // Exports uses export actions handled by ExportFormatModal
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
			Title: "Statement Summary",
			Fields: []types.DetailFieldConfig{
				{Key: "date", Label: "Date", Formatter: "datetime"},
				{Key: "accountedFor", Label: "Account", Formatter: "address"},
				{Key: "reconciled", Label: "Reconciled"},
				{Key: "reconciliationType", Label: "Type"},
			},
		},
		{
			Title: "Asset Information",
			Fields: []types.DetailFieldConfig{
				{Key: "asset", Label: "Asset", Formatter: "address"},
				{Key: "symbol", Label: "Symbol"},
				{Key: "decimals", Label: "Decimals"},
				{Key: "spotPrice", Label: "Spot Price"},
				{Key: "priceSource", Label: "Price Source"},
			},
		},
		{
			Title: "Balance Reconciliation",
			Fields: []types.DetailFieldConfig{
				{Key: "begBal", Label: "Begin Balance", Formatter: "wei"},
				{Key: "totalIn", Label: "Total In", Formatter: "wei"},
				{Key: "totalOut", Label: "Total Out", Formatter: "wei"},
				{Key: "amountNet", Label: "Net Amount", Formatter: "wei"},
				{Key: "endBal", Label: "End Balance", Formatter: "wei"},
				{Key: "endBalCalc", Label: "Calculated End Balance", Formatter: "wei"},
			},
		},
		{
			Title: "Inflow Details",
			Fields: []types.DetailFieldConfig{
				{Key: "amountIn", Label: "Amount In", Formatter: "wei"},
				{Key: "internalIn", Label: "Internal In", Formatter: "wei"},
				{Key: "selfDestructIn", Label: "Self Destruct In", Formatter: "wei"},
				{Key: "minerBaseRewardIn", Label: "Base Reward In", Formatter: "wei"},
				{Key: "minerTxFeeIn", Label: "Tx Fee In", Formatter: "wei"},
				{Key: "prefundIn", Label: "Prefund In", Formatter: "wei"},
			},
		},
		{
			Title: "Outflow Details",
			Fields: []types.DetailFieldConfig{
				{Key: "amountOut", Label: "Amount Out", Formatter: "wei"},
				{Key: "internalOut", Label: "Internal Out", Formatter: "wei"},
				{Key: "selfDestructOut", Label: "Self Destruct Out", Formatter: "wei"},
				{Key: "gasOut", Label: "Gas Out", Formatter: "wei"},
			},
		},
		{
			Title: "Transaction Details",
			Fields: []types.DetailFieldConfig{
				{Key: "blockNumber", Label: "Block Number"},
				{Key: "transactionIndex", Label: "Transaction Index"},
				{Key: "logIndex", Label: "Log Index"},
				{Key: "transactionHash", Label: "Transaction Hash", Formatter: "hash"},
				{Key: "sender", Label: "Sender", Formatter: "address"},
				{Key: "recipient", Label: "Recipient", Formatter: "address"},
			},
		},
		{
			Title: "Reconciliation Analysis",
			Fields: []types.DetailFieldConfig{
				{Key: "prevBal", Label: "Previous Balance", Formatter: "wei"},
				{Key: "begBalDiff", Label: "Begin Balance Diff", Formatter: "wei"},
				{Key: "endBalDiff", Label: "End Balance Diff", Formatter: "wei"},
				{Key: "correctingReasons", Label: "Correcting Reasons"},
			},
		},
		{
			Title: "Correction Entries",
			Fields: []types.DetailFieldConfig{
				{Key: "correctBegBalIn", Label: "Correct Begin Bal In", Formatter: "wei"},
				{Key: "correctAmountIn", Label: "Correct Amount In", Formatter: "wei"},
				{Key: "correctEndBalIn", Label: "Correct End Bal In", Formatter: "wei"},
				{Key: "correctBegBalOut", Label: "Correct Begin Bal Out", Formatter: "wei"},
				{Key: "correctAmountOut", Label: "Correct Amount Out", Formatter: "wei"},
				{Key: "correctEndBalOut", Label: "Correct End Bal Out", Formatter: "wei"},
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
			Title: "Balance Information",
			Fields: []types.DetailFieldConfig{
				{Key: "balance", Label: "Balance", Formatter: "wei"},
				{Key: "priorBalance", Label: "Prior Balance", Formatter: "wei"},
				{Key: "totalSupply", Label: "Total Supply", Formatter: "wei"},
				{Key: "type", Label: "Token Type"},
			},
		},
		{
			Title: "Token Details",
			Fields: []types.DetailFieldConfig{
				{Key: "address", Label: "Token Address", Formatter: "address"},
				{Key: "holder", Label: "Holder", Formatter: "address"},
				{Key: "symbol", Label: "Symbol"},
				{Key: "name", Label: "Token Name"},
				{Key: "decimals", Label: "Decimals"},
			},
		},
		{
			Title: "Transaction Context",
			Fields: []types.DetailFieldConfig{
				{Key: "blockNumber", Label: "Block Number"},
				{Key: "transactionIndex", Label: "Transaction Index"},
				{Key: "timestamp", Label: "Timestamp", Formatter: "datetime"},
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
			Title: "Transfer Details",
			Fields: []types.DetailFieldConfig{
				{Key: "sender", Label: "Sender", Formatter: "address"},
				{Key: "recipient", Label: "Recipient", Formatter: "address"},
				{Key: "holder", Label: "Holder", Formatter: "address"},
				{Key: "asset", Label: "Asset", Formatter: "address"},
				{Key: "decimals", Label: "Decimals"},
			},
		},
		{
			Title: "Amount Breakdown",
			Fields: []types.DetailFieldConfig{
				{Key: "amountIn", Label: "Amount In", Formatter: "wei"},
				{Key: "amountOut", Label: "Amount Out", Formatter: "wei"},
				{Key: "internalIn", Label: "Internal In", Formatter: "wei"},
				{Key: "internalOut", Label: "Internal Out", Formatter: "wei"},
				{Key: "gasOut", Label: "Gas Out", Formatter: "wei"},
			},
		},
		{
			Title: "Mining Rewards",
			Fields: []types.DetailFieldConfig{
				{Key: "minerBaseRewardIn", Label: "Base Reward In", Formatter: "wei"},
				{Key: "minerNephewRewardIn", Label: "Nephew Reward In", Formatter: "wei"},
				{Key: "minerTxFeeIn", Label: "Tx Fee In", Formatter: "wei"},
				{Key: "minerUncleRewardIn", Label: "Uncle Reward In", Formatter: "wei"},
			},
		},
		{
			Title: "Special Transfers",
			Fields: []types.DetailFieldConfig{
				{Key: "selfDestructIn", Label: "Self Destruct In", Formatter: "wei"},
				{Key: "selfDestructOut", Label: "Self Destruct Out", Formatter: "wei"},
				{Key: "prefundIn", Label: "Prefund In", Formatter: "wei"},
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
			Title: "Transaction Overview",
			Fields: []types.DetailFieldConfig{
				{Key: "hash", Label: "Hash", Formatter: "hash"},
				{Key: "from", Label: "From", Formatter: "address"},
				{Key: "to", Label: "To", Formatter: "address"},
				{Key: "value", Label: "Value", Formatter: "wei"},
				{Key: "isError", Label: "Error Status"},
				{Key: "hasToken", Label: "Has Token"},
			},
		},
		{
			Title: "Gas Information",
			Fields: []types.DetailFieldConfig{
				{Key: "gas", Label: "Gas Limit"},
				{Key: "gasUsed", Label: "Gas Used", Formatter: "gas"},
				{Key: "gasPrice", Label: "Gas Price", Formatter: "gas"},
				{Key: "maxFeePerGas", Label: "Max Fee Per Gas", Formatter: "gas"},
				{Key: "maxPriorityFeePerGas", Label: "Max Priority Fee", Formatter: "gas"},
			},
		},
		{
			Title: "Block Context",
			Fields: []types.DetailFieldConfig{
				{Key: "blockNumber", Label: "Block Number"},
				{Key: "blockHash", Label: "Block Hash", Formatter: "hash"},
				{Key: "transactionIndex", Label: "Transaction Index"},
				{Key: "timestamp", Label: "Timestamp", Formatter: "datetime"},
			},
		},
		{
			Title: "Transaction Details",
			Fields: []types.DetailFieldConfig{
				{Key: "nonce", Label: "Nonce"},
				{Key: "type", Label: "Transaction Type"},
				{Key: "input", Label: "Input Data"},
				{Key: "articulatedTx", Label: "Articulated Transaction", Formatter: "json"},
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
			Title: "Withdrawal Details",
			Fields: []types.DetailFieldConfig{
				{Key: "address", Label: "Recipient Address", Formatter: "address"},
				{Key: "amount", Label: "Amount", Formatter: "wei"},
				{Key: "validatorIndex", Label: "Validator Index"},
				{Key: "index", Label: "Withdrawal Index"},
			},
		},
		{
			Title: "Block Information",
			Fields: []types.DetailFieldConfig{
				{Key: "blockNumber", Label: "Block Number"},
				{Key: "timestamp", Label: "Timestamp", Formatter: "datetime"},
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
			Title: "Asset Information",
			Fields: []types.DetailFieldConfig{
				{Key: "address", Label: "Contract Address", Formatter: "address"},
				{Key: "name", Label: "Token Name"},
				{Key: "symbol", Label: "Symbol"},
				{Key: "decimals", Label: "Decimals"},
				{Key: "totalSupply", Label: "Total Supply", Formatter: "wei"},
			},
		},
		{
			Title: "Asset Classification",
			Fields: []types.DetailFieldConfig{
				{Key: "source", Label: "Source"},
				{Key: "tags", Label: "Tags"},
				{Key: "isContract", Label: "Is Contract"},
				{Key: "isCustom", Label: "Is Custom"},
				{Key: "isErc20", Label: "Is ERC20"},
				{Key: "isErc721", Label: "Is ERC721"},
				{Key: "isPrefund", Label: "Is Prefund"},
				{Key: "deleted", Label: "Deleted"},
			},
		},
		{
			Title: "Additional Data",
			Fields: []types.DetailFieldConfig{
				{Key: "parts", Label: "Parts"},
				{Key: "prefund", Label: "Prefund Amount", Formatter: "wei"},
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
			Title: "Log Overview",
			Fields: []types.DetailFieldConfig{
				{Key: "address", Label: "Contract Address", Formatter: "address"},
				{Key: "data", Label: "Data"},
				{Key: "logIndex", Label: "Log Index"},
			},
		},
		{
			Title: "Topics",
			Fields: []types.DetailFieldConfig{
				{Key: "topics.0", Label: "Topic 0 (Event Signature)", Formatter: "hash"},
				{Key: "topics.1", Label: "Topic 1", Formatter: "hash"},
				{Key: "topics.2", Label: "Topic 2", Formatter: "hash"},
				{Key: "topics.3", Label: "Topic 3", Formatter: "hash"},
			},
		},
		{
			Title: "Transaction Context",
			Fields: []types.DetailFieldConfig{
				{Key: "blockNumber", Label: "Block Number"},
				{Key: "blockHash", Label: "Block Hash", Formatter: "hash"},
				{Key: "transactionHash", Label: "Transaction Hash", Formatter: "hash"},
				{Key: "transactionIndex", Label: "Transaction Index"},
				{Key: "timestamp", Label: "Timestamp", Formatter: "datetime"},
			},
		},
		{
			Title: "Articulated Information",
			Fields: []types.DetailFieldConfig{
				{Key: "articulatedLog", Label: "Articulated Log", Formatter: "json"},
				{Key: "compressedLog", Label: "Compressed Log"},
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
			Title: "Trace Overview",
			Fields: []types.DetailFieldConfig{
				{Key: "type", Label: "Trace Type"},
				{Key: "error", Label: "Error"},
				{Key: "subtraces", Label: "Subtraces Count"},
				{Key: "traceAddress", Label: "Trace Address"},
			},
		},
		{
			Title: "Trace Action",
			Fields: []types.DetailFieldConfig{
				{Key: "from", Label: "From", Formatter: "address"},
				{Key: "to", Label: "To", Formatter: "address"},
				{Key: "value", Label: "Value", Formatter: "wei"},
				{Key: "gas", Label: "Gas Limit"},
				{Key: "callType", Label: "Call Type"},
				{Key: "input", Label: "Input Data"},
			},
		},
		{
			Title: "Trace Result",
			Fields: []types.DetailFieldConfig{
				{Key: "gasUsed", Label: "Gas Used", Formatter: "gas"},
				{Key: "output", Label: "Output"},
				{Key: "address", Label: "Created Address", Formatter: "address"},
				{Key: "code", Label: "Code"},
			},
		},
		{
			Title: "Transaction Context",
			Fields: []types.DetailFieldConfig{
				{Key: "blockNumber", Label: "Block Number"},
				{Key: "blockHash", Label: "Block Hash", Formatter: "hash"},
				{Key: "transactionHash", Label: "Transaction Hash", Formatter: "hash"},
				{Key: "transactionIndex", Label: "Transaction Index"},
				{Key: "timestamp", Label: "Timestamp", Formatter: "datetime"},
			},
		},
		{
			Title: "Articulated Information",
			Fields: []types.DetailFieldConfig{
				{Key: "articulatedTrace", Label: "Articulated Trace", Formatter: "json"},
				{Key: "compressedTrace", Label: "Compressed Trace"},
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
			Title: "Receipt Overview",
			Fields: []types.DetailFieldConfig{
				{Key: "transactionHash", Label: "Transaction Hash", Formatter: "hash"},
				{Key: "status", Label: "Status"},
				{Key: "isError", Label: "Is Error"},
				{Key: "contractAddress", Label: "Contract Address", Formatter: "address"},
			},
		},
		{
			Title: "Transaction Details",
			Fields: []types.DetailFieldConfig{
				{Key: "from", Label: "From", Formatter: "address"},
				{Key: "to", Label: "To", Formatter: "address"},
				{Key: "transactionIndex", Label: "Transaction Index"},
			},
		},
		{
			Title: "Gas Information",
			Fields: []types.DetailFieldConfig{
				{Key: "gasUsed", Label: "Gas Used", Formatter: "gas"},
				{Key: "cumulativeGasUsed", Label: "Cumulative Gas Used", Formatter: "gas"},
				{Key: "effectiveGasPrice", Label: "Effective Gas Price", Formatter: "gas"},
			},
		},
		{
			Title: "Block Context",
			Fields: []types.DetailFieldConfig{
				{Key: "blockNumber", Label: "Block Number"},
				{Key: "blockHash", Label: "Block Hash", Formatter: "hash"},
			},
		},
		{
			Title: "Additional Information",
			Fields: []types.DetailFieldConfig{
				{Key: "logsBloom", Label: "Logs Bloom"},
				{Key: "logs", Label: "Log Count"},
			},
		},
	}
}
