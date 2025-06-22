// Copyright 2016, 2025 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package exports

import (
	"fmt"
	"sync"

	// EXISTING_CODE
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/preferences"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/store"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
	// EXISTING_CODE
)

type Token = coreTypes.Token
type Statement = coreTypes.Statement
type Transaction = coreTypes.Transaction
type Transfer = coreTypes.Transfer

var (
	balancesStoreMap = make(map[string]*store.Store[Token])
	balancesStoreMu  sync.Mutex

	statementsStoreMap = make(map[string]*store.Store[Statement])
	statementsStoreMu  sync.Mutex

	transactionsStoreMap = make(map[string]*store.Store[Transaction])
	transactionsStoreMu  sync.Mutex

	transfersStoreMap = make(map[string]*store.Store[Transfer])
	transfersStoreMu  sync.Mutex
)

func (c *ExportsCollection) getBalancesStore(chain string, address string) *store.Store[Token] {
	balancesStoreMu.Lock()
	defer balancesStoreMu.Unlock()

	chainName := chain
	if chainName == "" {
		chainName = preferences.GetChain()
	}

	storeKey := getStoreKey(chainName, address)
	theStore, _ := balancesStoreMap[storeKey]

	if theStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			// EXISTING_CODE
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
			// EXISTING_CODE
			return nil
		}

		processFunc := func(item interface{}) *Token {
			// EXISTING_CODE
			// EXISTING_CODE
			if it, ok := item.(*Token); ok {
				return it
			}
			return nil
		}

		mappingFunc := func(item *Token) (key interface{}, includeInMap bool) {
			// EXISTING_CODE
			// EXISTING_CODE
			return nil, false
		}

		// EXISTING_CODE
		storeName := fmt.Sprintf("exports_balances_%s_%s", chainName, address)
		// EXISTING_CODE
		theStore = store.NewStore(storeName, queryFunc, processFunc, mappingFunc)
		balancesStoreMap[storeKey] = theStore
	}

	return theStore
}

func (c *ExportsCollection) getStatementsStore(chain string, address string) *store.Store[Statement] {
	statementsStoreMu.Lock()
	defer statementsStoreMu.Unlock()

	chainName := chain
	if chainName == "" {
		chainName = preferences.GetChain()
	}

	storeKey := getStoreKey(chainName, address)
	theStore, _ := statementsStoreMap[storeKey]

	if theStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			// EXISTING_CODE
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
			// EXISTING_CODE
			return nil
		}

		processFunc := func(item interface{}) *Statement {
			// EXISTING_CODE
			// EXISTING_CODE
			if it, ok := item.(*Statement); ok {
				return it
			}
			return nil
		}

		mappingFunc := func(item *Statement) (key interface{}, includeInMap bool) {
			// EXISTING_CODE
			// EXISTING_CODE
			return nil, false
		}

		// EXISTING_CODE
		storeName := fmt.Sprintf("exports_statements_%s_%s", chainName, address)
		// EXISTING_CODE
		theStore = store.NewStore(storeName, queryFunc, processFunc, mappingFunc)
		statementsStoreMap[storeKey] = theStore
	}

	return theStore
}

func (c *ExportsCollection) getTransactionsStore(chain string, address string) *store.Store[Transaction] {
	transactionsStoreMu.Lock()
	defer transactionsStoreMu.Unlock()

	chainName := chain
	if chainName == "" {
		chainName = preferences.GetChain()
	}

	storeKey := getStoreKey(chainName, address)
	theStore, _ := transactionsStoreMap[storeKey]

	if theStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			// EXISTING_CODE
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
			// EXISTING_CODE
			return nil
		}

		processFunc := func(item interface{}) *Transaction {
			// EXISTING_CODE
			// EXISTING_CODE
			if it, ok := item.(*Transaction); ok {
				return it
			}
			return nil
		}

		mappingFunc := func(item *Transaction) (key interface{}, includeInMap bool) {
			// EXISTING_CODE
			// EXISTING_CODE
			return nil, false
		}

		// EXISTING_CODE
		storeName := fmt.Sprintf("exports_transactions_%s_%s", chainName, address)
		// EXISTING_CODE
		theStore = store.NewStore(storeName, queryFunc, processFunc, mappingFunc)
		transactionsStoreMap[storeKey] = theStore
	}

	return theStore
}

func (c *ExportsCollection) getTransfersStore(chain string, address string) *store.Store[Transfer] {
	transfersStoreMu.Lock()
	defer transfersStoreMu.Unlock()

	chainName := chain
	if chainName == "" {
		chainName = preferences.GetChain()
	}

	storeKey := getStoreKey(chainName, address)
	theStore, _ := transfersStoreMap[storeKey]

	if theStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			// EXISTING_CODE
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
			// EXISTING_CODE
			return nil
		}

		processFunc := func(item interface{}) *Transfer {
			// EXISTING_CODE
			// EXISTING_CODE
			if it, ok := item.(*Transfer); ok {
				return it
			}
			return nil
		}

		mappingFunc := func(item *Transfer) (key interface{}, includeInMap bool) {
			// EXISTING_CODE
			// EXISTING_CODE
			return nil, false
		}

		// EXISTING_CODE
		storeName := fmt.Sprintf("exports_transfers_%s_%s", chainName, address)
		// EXISTING_CODE
		theStore = store.NewStore(storeName, queryFunc, processFunc, mappingFunc)
		transfersStoreMap[storeKey] = theStore
	}

	return theStore
}

func (c *ExportsCollection) GetStoreName(dataFacet types.DataFacet) string {
	switch dataFacet {
	case ExportsStatements:
		return "exports-statements"
	case ExportsTransfers:
		return "exports-transfers"
	case ExportsBalances:
		return "exports-balances"
	case ExportsTransactions:
		return "exports-transactions"
	default:
		return ""
	}
}

// TODO: THIS SHOULD BE PER STORE - SEE EXPORT COMMENTS
func GetExportsCount(chain string, address string, dataFacet string) (int, error) {
	chainName := preferences.GetChain()
	countOpts := sdk.ExportOptions{
		Globals: sdk.Globals{Cache: true, Chain: chainName},
	}
	if countResult, _, err := countOpts.ExportCount(); err != nil {
		return 0, fmt.Errorf("ExportsCount query error: %v", err)
	} else if len(countResult) > 0 {
		return len(countResult), nil
	}
	return 0, nil
}

// EXISTING_CODE
func GetExportsCount2(chain string, address string, dataFacet string) (int, error) {
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

func getExportsTransactionsCount(chain string, address string) (int, error) {
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
		count := int(0)
		for _, monitor := range results {
			count += int(monitor.NRecords)
		}
		return count, nil
	}
}

func getExportsStatementsCount(chain string, address string) (int, error) {
	// For statements, we can use the same base count since statements are derived from transactions
	return getExportsTransactionsCount(chain, address)
}

func getExportsTransfersCount(chain string, address string) (int, error) {
	// For transfers, we can use the same base count since transfers are derived from transactions
	return getExportsTransactionsCount(chain, address)
}

func getExportsBalancesCount(chain string, address string) (int, error) {
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
		transactionsStoreMu.Lock()
		if store, exists := transactionsStoreMap[storeKey]; exists {
			store.Reset()
		}
		transactionsStoreMu.Unlock()
	case ExportsStatements:
		statementsStoreMu.Lock()
		if store, exists := statementsStoreMap[storeKey]; exists {
			store.Reset()
		}
		statementsStoreMu.Unlock()
	case ExportsTransfers:
		transfersStoreMu.Lock()
		if store, exists := transfersStoreMap[storeKey]; exists {
			store.Reset()
		}
		transfersStoreMu.Unlock()
	case ExportsBalances:
		balancesStoreMu.Lock()
		if store, exists := balancesStoreMap[storeKey]; exists {
			store.Reset()
		}
		balancesStoreMu.Unlock()
	}
}

// ClearExportsStores clears all cached stores for a given chain and address
func ClearExportsStores(chain, address string) {
	chainName := chain
	if chainName == "" {
		chainName = preferences.GetChain()
	}

	storeKey := getStoreKey(chainName, address)

	transactionsStoreMu.Lock()
	delete(transactionsStoreMap, storeKey)
	transactionsStoreMu.Unlock()

	statementsStoreMu.Lock()
	delete(statementsStoreMap, storeKey)
	statementsStoreMu.Unlock()

	transfersStoreMu.Lock()
	delete(transfersStoreMap, storeKey)
	transfersStoreMu.Unlock()

	balancesStoreMu.Lock()
	delete(balancesStoreMap, storeKey)
	balancesStoreMu.Unlock()
}

// ClearAllExportsStores clears all cached stores (useful for global reset)
func ClearAllExportsStores() {
	transactionsStoreMu.Lock()
	transactionsStoreMap = make(map[string]*store.Store[Transaction])
	transactionsStoreMu.Unlock()

	statementsStoreMu.Lock()
	statementsStoreMap = make(map[string]*store.Store[Statement])
	statementsStoreMu.Unlock()

	transfersStoreMu.Lock()
	transfersStoreMap = make(map[string]*store.Store[Transfer])
	transfersStoreMu.Unlock()

	balancesStoreMu.Lock()
	balancesStoreMap = make(map[string]*store.Store[Token])
	balancesStoreMu.Unlock()
}

func getStoreKey(chain, address string) string {
	return fmt.Sprintf("%s_%s", chain, address)
}
// EXISTING_CODE
