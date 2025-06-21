package exports

import (
	"fmt"
	"sync"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/preferences"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/store"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

type Statement = coreTypes.Statement
type Transfer = coreTypes.Transfer
type Token = coreTypes.Token
type Transaction = coreTypes.Transaction

var (
	transactionsStores = make(map[string]*store.Store[Transaction])
	statementsStores   = make(map[string]*store.Store[Statement])
	transfersStores    = make(map[string]*store.Store[Transfer])
	balancesStores     = make(map[string]*store.Store[Token])
	transactionsMu     sync.Mutex
	statementsMu       sync.Mutex
	transfersMu        sync.Mutex
	balancesMu         sync.Mutex
)

// getStoreKey creates a unique key for caching stores per (chain, address)
func getStoreKey(chain, address string) string {
	return fmt.Sprintf("%s_%s", chain, address)
}

// GetExportsStatementsStore returns a singleton store for statements export data per (chain, address)
func GetExportsStatementsStore(chain string, address string) *store.Store[Statement] {
	statementsMu.Lock()
	defer statementsMu.Unlock()

	// Use provided chain or fall back to preferences
	chainName := chain
	if chainName == "" {
		chainName = preferences.GetChain()
	}

	storeKey := getStoreKey(chainName, address)
	if existingStore, exists := statementsStores[storeKey]; exists {
		return existingStore
	}

	queryFunc := func(ctx *output.RenderCtx) error {
		exportOpts := sdk.ExportOptions{
			Globals:    sdk.Globals{Cache: true, Verbose: true, Chain: chainName},
			RenderCtx:  ctx,
			Addrs:      []string{address},
			Accounting: true, // Enable accounting for statements
		}

		if _, _, err := exportOpts.ExportStatements(); err != nil {
			wrappedErr := types.NewSDKError("exports", ExportsStatements, "fetch", err)
			logger.Error(fmt.Sprintf("Exports statements SDK query error: %v", wrappedErr))
			return wrappedErr
		}
		return nil
	}

	processFunc := func(itemIntf interface{}) *Statement {
		if stmt, ok := itemIntf.(*Statement); ok {
			return stmt
		}
		return nil
	}

	newStore := store.NewStore(
		fmt.Sprintf("exports_statements_%s_%s", chainName, address),
		queryFunc,
		processFunc,
		nil,
	)

	statementsStores[storeKey] = newStore
	return newStore
}

// GetExportsTransfersStore returns a singleton store for transfers export data per (chain, address)
func GetExportsTransfersStore(chain string, address string) *store.Store[Transfer] {
	transfersMu.Lock()
	defer transfersMu.Unlock()

	// Use provided chain or fall back to preferences
	chainName := chain
	if chainName == "" {
		chainName = preferences.GetChain()
	}

	storeKey := getStoreKey(chainName, address)
	if existingStore, exists := transfersStores[storeKey]; exists {
		return existingStore
	}

	queryFunc := func(ctx *output.RenderCtx) error {
		exportOpts := sdk.ExportOptions{
			Globals:    sdk.Globals{Cache: true, Verbose: true, Chain: chainName},
			RenderCtx:  ctx,
			Addrs:      []string{address},
			Accounting: true, // Enable accounting for transfers
		}

		if _, _, err := exportOpts.ExportTransfers(); err != nil {
			wrappedErr := types.NewSDKError("exports", ExportsTransfers, "fetch", err)
			logger.Error(fmt.Sprintf("Exports transfers SDK query error: %v", wrappedErr))
			return wrappedErr
		}
		return nil
	}

	processFunc := func(itemIntf interface{}) *Transfer {
		if transfer, ok := itemIntf.(*Transfer); ok {
			return transfer
		}
		return nil
	}

	newStore := store.NewStore(
		fmt.Sprintf("exports_transfers_%s_%s", chainName, address),
		queryFunc,
		processFunc,
		nil,
	)

	transfersStores[storeKey] = newStore
	return newStore
}

// GetExportsBalancesStore returns a singleton store for balances export data per (chain, address)
func GetExportsBalancesStore(chain string, address string) *store.Store[Token] {
	balancesMu.Lock()
	defer balancesMu.Unlock()

	// Use provided chain or fall back to preferences
	chainName := chain
	if chainName == "" {
		chainName = preferences.GetChain()
	}

	storeKey := getStoreKey(chainName, address)
	if existingStore, exists := balancesStores[storeKey]; exists {
		return existingStore
	}

	queryFunc := func(ctx *output.RenderCtx) error {
		exportOpts := sdk.ExportOptions{
			Globals:   sdk.Globals{Cache: true, Verbose: true, Chain: chainName},
			RenderCtx: ctx,
			Addrs:     []string{address},
		}

		if _, _, err := exportOpts.ExportBalances(); err != nil {
			wrappedErr := types.NewSDKError("exports", ExportsBalances, "fetch", err)
			logger.Error(fmt.Sprintf("Exports balances SDK query error: %v", wrappedErr))
			return wrappedErr
		}

		return nil
	}

	processFunc := func(itemIntf interface{}) *Token {
		if state, ok := itemIntf.(*Token); ok {
			return state
		}
		return nil
	}

	newStore := store.NewStore(
		fmt.Sprintf("exports_balances_%s_%s", chainName, address),
		queryFunc,
		processFunc,
		nil,
	)

	balancesStores[storeKey] = newStore
	return newStore
}

// GetExportsTransactionsStore returns a singleton store for transactions export data per (chain, address)
func GetExportsTransactionsStore(chain string, address string) *store.Store[Transaction] {
	transactionsMu.Lock()
	defer transactionsMu.Unlock()

	chainName := chain
	if chainName == "" {
		chainName = preferences.GetChain()
	}

	storeKey := getStoreKey(chainName, address)
	if existingStore, exists := transactionsStores[storeKey]; exists {
		return existingStore
	}

	queryFunc := func(ctx *output.RenderCtx) error {
		exportOpts := sdk.ExportOptions{
			Globals:   sdk.Globals{Cache: true, Verbose: true, Chain: chainName},
			RenderCtx: ctx,
			Addrs:     []string{address},
		}

		if _, _, err := exportOpts.Export(); err != nil {
			wrappedErr := types.NewSDKError("exports", ExportsTransactions, "fetch", err)
			logger.Error(fmt.Sprintf("Exports transactions SDK query error: %v", wrappedErr))
			return wrappedErr
		}
		return nil
	}

	processFunc := func(itemIntf interface{}) *Transaction {
		if tx, ok := itemIntf.(*Transaction); ok {
			return tx
		}
		return nil
	}

	newStore := store.NewStore(
		fmt.Sprintf("exports_transactions_%s_%s", chainName, address),
		queryFunc,
		processFunc,
		nil, // No mapping function needed
	)

	transactionsStores[storeKey] = newStore
	return newStore
}

// GetExportsCount returns the count of items for a specific export facet
func GetExportsCount(chain string, address string, dataFacet string) (uint64, error) {
	switch dataFacet {
	case string(ExportsTransactions):
		return getExportsTransactionsCount(chain, address)
	case string(ExportsStatements):
		return getExportsStatementsCount(chain, address)
	case string(ExportsTransfers):
		return getExportsTransfersCount(chain, address)
	case string(ExportsBalances):
		return getExportsBalancesCount(chain, address)
	default:
		return 0, fmt.Errorf("unknown dataFacet: %s", dataFacet)
	}
}

// Count optimization functions for each facet
func getExportsTransactionsCount(chain string, address string) (uint64, error) {
	chainName := chain
	if chainName == "" {
		chainName = preferences.GetChain()
	}

	exportOpts := sdk.ExportOptions{
		Globals: sdk.Globals{Cache: true, Chain: chainName},
		Addrs:   []string{address},
	}

	// Use ExportCount for optimized counting
	if results, _, err := exportOpts.ExportCount(); err != nil {
		return 0, fmt.Errorf("failed to get exports transactions count: %w", err)
	} else {
		// Count the relevant transactions for this address
		count := uint64(0)
		for _, monitor := range results {
			count += uint64(monitor.NRecords)
		}
		return count, nil
	}
}

func getExportsStatementsCount(chain string, address string) (uint64, error) {
	// For statements, we can use the same base count since statements are derived from transactions
	return getExportsTransactionsCount(chain, address)
}

func getExportsTransfersCount(chain string, address string) (uint64, error) {
	// For transfers, we can use the same base count since transfers are derived from transactions
	return getExportsTransactionsCount(chain, address)
}

func getExportsBalancesCount(chain string, address string) (uint64, error) {
	// For balances, we can use the same base count since balances are derived from transactions
	return getExportsTransactionsCount(chain, address)
}

// ResetExportsStore resets a specific store for a given chain, address, and dataFacet
func ResetExportsStore(chain, address string, dataFacet types.DataFacet) {
	chainName := chain
	if chainName == "" {
		chainName = preferences.GetChain()
	}

	storeKey := getStoreKey(chainName, address)

	switch dataFacet {
	case ExportsTransactions:
		transactionsMu.Lock()
		if store, exists := transactionsStores[storeKey]; exists {
			store.Reset()
		}
		transactionsMu.Unlock()
	case ExportsStatements:
		statementsMu.Lock()
		if store, exists := statementsStores[storeKey]; exists {
			store.Reset()
		}
		statementsMu.Unlock()
	case ExportsTransfers:
		transfersMu.Lock()
		if store, exists := transfersStores[storeKey]; exists {
			store.Reset()
		}
		transfersMu.Unlock()
	case ExportsBalances:
		balancesMu.Lock()
		if store, exists := balancesStores[storeKey]; exists {
			store.Reset()
		}
		balancesMu.Unlock()
	}
}

// ClearExportsStores clears all cached stores for a given chain and address
func ClearExportsStores(chain, address string) {
	chainName := chain
	if chainName == "" {
		chainName = preferences.GetChain()
	}

	storeKey := getStoreKey(chainName, address)

	transactionsMu.Lock()
	delete(transactionsStores, storeKey)
	transactionsMu.Unlock()

	statementsMu.Lock()
	delete(statementsStores, storeKey)
	statementsMu.Unlock()

	transfersMu.Lock()
	delete(transfersStores, storeKey)
	transfersMu.Unlock()

	balancesMu.Lock()
	delete(balancesStores, storeKey)
	balancesMu.Unlock()
}

// ClearAllExportsStores clears all cached stores (useful for global reset)
func ClearAllExportsStores() {
	transactionsMu.Lock()
	transactionsStores = make(map[string]*store.Store[Transaction])
	transactionsMu.Unlock()

	statementsMu.Lock()
	statementsStores = make(map[string]*store.Store[Statement])
	statementsMu.Unlock()

	transfersMu.Lock()
	transfersStores = make(map[string]*store.Store[Transfer])
	transfersMu.Unlock()

	balancesMu.Lock()
	balancesStores = make(map[string]*store.Store[Token])
	balancesMu.Unlock()
}
