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
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/preferences"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/store"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
	// EXISTING_CODE
)

// EXISTING_CODE
// EXISTING_CODE

type Balance = sdk.Balance
type Statement = sdk.Statement
type Transaction = sdk.Transaction
type Transfer = sdk.Transfer

var (
	balancesStore   = make(map[string]*store.Store[Balance])
	balancesStoreMu sync.Mutex

	statementsStore   = make(map[string]*store.Store[Statement])
	statementsStoreMu sync.Mutex

	transactionsStore   = make(map[string]*store.Store[Transaction])
	transactionsStoreMu sync.Mutex

	transfersStore   = make(map[string]*store.Store[Transfer])
	transfersStoreMu sync.Mutex
)

func (c *ExportsCollection) getBalancesStore() *store.Store[Balance] {
	balancesStoreMu.Lock()
	defer balancesStoreMu.Unlock()

	chain := preferences.GetLastChain()
	address := preferences.GetLastAddress()
	storeKey := getStoreKey(chain, address)
	theStore := balancesStore[storeKey]
	if theStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			// EXISTING_CODE
			exportOpts := sdk.ExportOptions{
				Globals:   sdk.Globals{Cache: true, Verbose: true, Chain: chain},
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

		processFunc := func(item interface{}) *Balance {
			// EXISTING_CODE
			// EXISTING_CODE
			if it, ok := item.(*Balance); ok {
				return it
			}
			return nil
		}

		mappingFunc := func(item *Balance) (key interface{}, includeInMap bool) {
			// EXISTING_CODE
			// EXISTING_CODE
			return nil, false
		}

		// EXISTING_CODE
		storeName := fmt.Sprintf("exports_balances_%s_%s", chain, address)
		// EXISTING_CODE
		theStore = store.NewStore(storeName, queryFunc, processFunc, mappingFunc)
		balancesStore[storeKey] = theStore
	}

	return theStore
}

func (c *ExportsCollection) getStatementsStore() *store.Store[Statement] {
	statementsStoreMu.Lock()
	defer statementsStoreMu.Unlock()

	chain := preferences.GetLastChain()
	address := preferences.GetLastAddress()
	storeKey := getStoreKey(chain, address)
	theStore := statementsStore[storeKey]
	if theStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			// EXISTING_CODE
			exportOpts := sdk.ExportOptions{
				Globals:    sdk.Globals{Cache: true, Verbose: true, Chain: chain},
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
		storeName := fmt.Sprintf("exports_statements_%s_%s", chain, address)
		// EXISTING_CODE
		theStore = store.NewStore(storeName, queryFunc, processFunc, mappingFunc)
		statementsStore[storeKey] = theStore
	}

	return theStore
}

func (c *ExportsCollection) getTransactionsStore() *store.Store[Transaction] {
	transactionsStoreMu.Lock()
	defer transactionsStoreMu.Unlock()

	chain := preferences.GetLastChain()
	address := preferences.GetLastAddress()
	storeKey := getStoreKey(chain, address)
	theStore := transactionsStore[storeKey]
	if theStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			// EXISTING_CODE
			exportOpts := sdk.ExportOptions{
				Globals:   sdk.Globals{Cache: true, Verbose: true, Chain: chain},
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
		storeName := fmt.Sprintf("exports_transactions_%s_%s", chain, address)
		// EXISTING_CODE
		theStore = store.NewStore(storeName, queryFunc, processFunc, mappingFunc)
		transactionsStore[storeKey] = theStore
	}

	return theStore
}

func (c *ExportsCollection) getTransfersStore() *store.Store[Transfer] {
	transfersStoreMu.Lock()
	defer transfersStoreMu.Unlock()

	chain := preferences.GetLastChain()
	address := preferences.GetLastAddress()
	storeKey := getStoreKey(chain, address)
	theStore := transfersStore[storeKey]
	if theStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			// EXISTING_CODE
			exportOpts := sdk.ExportOptions{
				Globals:    sdk.Globals{Cache: true, Verbose: true, Chain: chain},
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
		storeName := fmt.Sprintf("exports_transfers_%s_%s", chain, address)
		// EXISTING_CODE
		theStore = store.NewStore(storeName, queryFunc, processFunc, mappingFunc)
		transfersStore[storeKey] = theStore
	}

	return theStore
}

func (c *ExportsCollection) GetStoreName(dataFacet types.DataFacet) string {
	switch dataFacet {
	case ExportsStatements:
		return "exports-statements"
	case ExportsBalances:
		return "exports-balances"
	case ExportsTransfers:
		return "exports-transfers"
	case ExportsTransactions:
		return "exports-transactions"
	default:
		return ""
	}
}

// TODO: THIS SHOULD BE PER STORE - SEE EXPORT COMMENTS
func GetExportsCount(payload *types.Payload) (int, error) {
	chain := preferences.GetLastChain()
	address := preferences.GetLastAddress()
	countOpts := sdk.ExportOptions{
		Globals: sdk.Globals{Cache: true, Chain: chain},
		Addrs:   []string{address},
	}
	if countResult, _, err := countOpts.ExportCount(); err != nil {
		return 0, fmt.Errorf("ExportsCount query error: %v", err)
	} else if len(countResult) > 0 {
		return len(countResult), nil
	}
	return 0, nil
}

var (
	collections   = make(map[store.CollectionKey]*ExportsCollection)
	collectionsMu sync.Mutex
)

func GetExportsCollection(payload *types.Payload) *ExportsCollection {
	collectionsMu.Lock()
	defer collectionsMu.Unlock()

	key := store.GetCollectionKey(payload)
	if collection, exists := collections[key]; exists {
		return collection
	}

	collection := NewExportsCollection()
	collections[key] = collection
	return collection
}

// EXISTING_CODE
func GetExportsCount2(dataFacet string) (int, error) {
	switch dataFacet {
	case string(ExportsTransactions):
		return getExportsTransactionsCount()
	case string(ExportsStatements):
		return getExportsStatementsCount()
	case string(ExportsTransfers):
		return getExportsTransfersCount()
	case string(ExportsBalances):
		return getExportsBalancesCount()
	default:
		return 0, fmt.Errorf("unknown dataFacet: %s", dataFacet)
	}
}

func getExportsTransactionsCount() (int, error) {
	chain := preferences.GetLastChain()
	address := preferences.GetLastAddress()
	exportOpts := sdk.ExportOptions{
		Globals: sdk.Globals{Cache: true, Chain: chain},
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

func getExportsStatementsCount() (int, error) {
	// For statements, we can use the same base count since statements are derived from transactions
	return getExportsTransactionsCount()
}

func getExportsTransfersCount() (int, error) {
	// For transfers, we can use the same base count since transfers are derived from transactions
	return getExportsTransactionsCount()
}

func getExportsBalancesCount() (int, error) {
	// For balances, we can use the same base count since balances are derived from transactions
	return getExportsTransactionsCount()
}

// ResetExportsStore resets a specific store for a given chain, address, and dataFacet
func ResetExportsStore(dataFacet types.DataFacet) {
	chain := preferences.GetLastChain()
	address := preferences.GetLastAddress()
	storeKey := getStoreKey(chain, address)

	switch dataFacet {
	case ExportsTransactions:
		transactionsStoreMu.Lock()
		if store, exists := transactionsStore[storeKey]; exists {
			store.Reset()
		}
		transactionsStoreMu.Unlock()
	case ExportsStatements:
		statementsStoreMu.Lock()
		if store, exists := statementsStore[storeKey]; exists {
			store.Reset()
		}
		statementsStoreMu.Unlock()
	case ExportsTransfers:
		transfersStoreMu.Lock()
		if store, exists := transfersStore[storeKey]; exists {
			store.Reset()
		}
		transfersStoreMu.Unlock()
	case ExportsBalances:
		balancesStoreMu.Lock()
		if store, exists := balancesStore[storeKey]; exists {
			store.Reset()
		}
		balancesStoreMu.Unlock()
	}
}

// ClearExportsStores clears all cached stores for a given chain and address
func ClearExportsStores() {
	chain := preferences.GetLastChain()
	address := preferences.GetLastAddress()
	storeKey := getStoreKey(chain, address)

	transactionsStoreMu.Lock()
	delete(transactionsStore, storeKey)
	transactionsStoreMu.Unlock()

	statementsStoreMu.Lock()
	delete(statementsStore, storeKey)
	statementsStoreMu.Unlock()

	transfersStoreMu.Lock()
	delete(transfersStore, storeKey)
	transfersStoreMu.Unlock()

	balancesStoreMu.Lock()
	delete(balancesStore, storeKey)
	balancesStoreMu.Unlock()
}

// ClearAllExportsStores clears all cached stores (useful for global reset)
func ClearAllExportsStores() {
	transactionsStoreMu.Lock()
	transactionsStore = make(map[string]*store.Store[Transaction])
	transactionsStoreMu.Unlock()

	statementsStoreMu.Lock()
	statementsStore = make(map[string]*store.Store[Statement])
	statementsStoreMu.Unlock()

	transfersStoreMu.Lock()
	transfersStore = make(map[string]*store.Store[Transfer])
	transfersStoreMu.Unlock()

	balancesStoreMu.Lock()
	balancesStore = make(map[string]*store.Store[Balance])
	balancesStoreMu.Unlock()
}

func getStoreKey(chain, address string) string {
	return fmt.Sprintf("%s_%s", chain, address)
}

// EXISTING_CODE
