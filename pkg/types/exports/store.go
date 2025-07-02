// Copyright 2016, 2026 The TrueBlocks Authors. All rights reserved.
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
	// EXISTING_CODE
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/logging"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/preferences"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/store"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

// EXISTING_CODE
// EXISTING_CODE

type Asset = sdk.Asset
type Balance = sdk.Balance
type Log = sdk.Log
type Receipt = sdk.Receipt
type Statement = sdk.Statement
type Trace = sdk.Trace
type Transaction = sdk.Transaction
type Transfer = sdk.Transfer
type Withdrawal = sdk.Withdrawal

var (
	assetsStore   = make(map[string]*store.Store[Asset])
	assetsStoreMu sync.Mutex

	balancesStore   = make(map[string]*store.Store[Balance])
	balancesStoreMu sync.Mutex

	logsStore   = make(map[string]*store.Store[Log])
	logsStoreMu sync.Mutex

	receiptsStore   = make(map[string]*store.Store[Receipt])
	receiptsStoreMu sync.Mutex

	statementsStore   = make(map[string]*store.Store[Statement])
	statementsStoreMu sync.Mutex

	tracesStore   = make(map[string]*store.Store[Trace])
	tracesStoreMu sync.Mutex

	transactionsStore   = make(map[string]*store.Store[Transaction])
	transactionsStoreMu sync.Mutex

	transfersStore   = make(map[string]*store.Store[Transfer])
	transfersStoreMu sync.Mutex

	withdrawalsStore   = make(map[string]*store.Store[Withdrawal])
	withdrawalsStoreMu sync.Mutex
)

func (c *ExportsCollection) getAssetsStore() *store.Store[Asset] {
	assetsStoreMu.Lock()
	defer assetsStoreMu.Unlock()

	chain := preferences.GetLastChain()
	address := preferences.GetLastAddress()
	storeKey := getStoreKey(chain, address)
	theStore := assetsStore[storeKey]
	if theStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			// EXISTING_CODE
			exportOpts := sdk.ExportOptions{
				Globals:   sdk.Globals{Cache: true, Verbose: true, Chain: chain},
				RenderCtx: ctx,
				Addrs:     []string{address},
			}
			if _, _, err := exportOpts.ExportAssets(); err != nil {
				wrappedErr := types.NewSDKError("exports", ExportsAssets, "fetch", err)
				logging.LogBackend(fmt.Sprintf("Exports assets SDK query error: %v", wrappedErr))
				return wrappedErr
			}
			// EXISTING_CODE
			return nil
		}

		processFunc := func(item interface{}) *Asset {
			// EXISTING_CODE
			// EXISTING_CODE
			if it, ok := item.(*Asset); ok {
				return it
			}
			return nil
		}

		mappingFunc := func(item *Asset) (key interface{}, includeInMap bool) {
			// EXISTING_CODE
			// EXISTING_CODE
			return nil, false
		}

		// EXISTING_CODE
		storeName := fmt.Sprintf("exports_assets_%s_%s", chain, address)
		// EXISTING_CODE
		theStore = store.NewStore(storeName, queryFunc, processFunc, mappingFunc)
		assetsStore[storeKey] = theStore
	}

	return theStore
}

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
				logging.LogBackend(fmt.Sprintf("Exports balances SDK query error: %v", wrappedErr))
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

func (c *ExportsCollection) getLogsStore() *store.Store[Log] {
	logsStoreMu.Lock()
	defer logsStoreMu.Unlock()

	chain := preferences.GetLastChain()
	address := preferences.GetLastAddress()
	storeKey := getStoreKey(chain, address)
	theStore := logsStore[storeKey]
	if theStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			// EXISTING_CODE
			exportOpts := sdk.ExportOptions{
				Globals:    sdk.Globals{Cache: true, Verbose: true, Chain: chain},
				RenderCtx:  ctx,
				Addrs:      []string{address},
				Articulate: true,
			}
			if _, _, err := exportOpts.ExportLogs(); err != nil {
				wrappedErr := types.NewSDKError("exports", ExportsLogs, "fetch", err)
				logging.LogBackend(fmt.Sprintf("Exports logs SDK query error: %v", wrappedErr))
				return wrappedErr
			}
			// EXISTING_CODE
			return nil
		}

		processFunc := func(item interface{}) *Log {
			// EXISTING_CODE
			// EXISTING_CODE
			if it, ok := item.(*Log); ok {
				return it
			}
			return nil
		}

		mappingFunc := func(item *Log) (key interface{}, includeInMap bool) {
			// EXISTING_CODE
			// EXISTING_CODE
			return nil, false
		}

		// EXISTING_CODE
		storeName := fmt.Sprintf("exports_logs_%s_%s", chain, address)
		// EXISTING_CODE
		theStore = store.NewStore(storeName, queryFunc, processFunc, mappingFunc)
		logsStore[storeKey] = theStore
	}

	return theStore
}

func (c *ExportsCollection) getReceiptsStore() *store.Store[Receipt] {
	receiptsStoreMu.Lock()
	defer receiptsStoreMu.Unlock()

	chain := preferences.GetLastChain()
	address := preferences.GetLastAddress()
	storeKey := getStoreKey(chain, address)
	theStore := receiptsStore[storeKey]
	if theStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			// EXISTING_CODE
			exportOpts := sdk.ExportOptions{
				Globals:   sdk.Globals{Cache: true, Verbose: true, Chain: chain},
				RenderCtx: ctx,
				Addrs:     []string{address},
			}
			if _, _, err := exportOpts.ExportReceipts(); err != nil {
				wrappedErr := types.NewSDKError("exports", ExportsReceipts, "fetch", err)
				logging.LogBackend(fmt.Sprintf("Exports receipts SDK query error: %v", wrappedErr))
				return wrappedErr
			}
			// EXISTING_CODE
			return nil
		}

		processFunc := func(item interface{}) *Receipt {
			// EXISTING_CODE
			// EXISTING_CODE
			if it, ok := item.(*Receipt); ok {
				return it
			}
			return nil
		}

		mappingFunc := func(item *Receipt) (key interface{}, includeInMap bool) {
			// EXISTING_CODE
			// EXISTING_CODE
			return nil, false
		}

		// EXISTING_CODE
		storeName := fmt.Sprintf("exports_receipts_%s_%s", chain, address)
		// EXISTING_CODE
		theStore = store.NewStore(storeName, queryFunc, processFunc, mappingFunc)
		receiptsStore[storeKey] = theStore
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
				logging.LogBackend(fmt.Sprintf("Exports statements SDK query error: %v", wrappedErr))
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

func (c *ExportsCollection) getTracesStore() *store.Store[Trace] {
	tracesStoreMu.Lock()
	defer tracesStoreMu.Unlock()

	chain := preferences.GetLastChain()
	address := preferences.GetLastAddress()
	storeKey := getStoreKey(chain, address)
	theStore := tracesStore[storeKey]
	if theStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			// EXISTING_CODE
			exportOpts := sdk.ExportOptions{
				Globals:   sdk.Globals{Cache: true, Verbose: true, Chain: chain},
				RenderCtx: ctx,
				Addrs:     []string{address},
			}
			if _, _, err := exportOpts.ExportTraces(); err != nil {
				wrappedErr := types.NewSDKError("exports", ExportsTraces, "fetch", err)
				logging.LogBackend(fmt.Sprintf("Exports traces SDK query error: %v", wrappedErr))
				return wrappedErr
			}
			// EXISTING_CODE
			return nil
		}

		processFunc := func(item interface{}) *Trace {
			// EXISTING_CODE
			// EXISTING_CODE
			if it, ok := item.(*Trace); ok {
				return it
			}
			return nil
		}

		mappingFunc := func(item *Trace) (key interface{}, includeInMap bool) {
			// EXISTING_CODE
			// EXISTING_CODE
			return nil, false
		}

		// EXISTING_CODE
		storeName := fmt.Sprintf("exports_traces_%s_%s", chain, address)
		// EXISTING_CODE
		theStore = store.NewStore(storeName, queryFunc, processFunc, mappingFunc)
		tracesStore[storeKey] = theStore
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
				logging.LogBackend(fmt.Sprintf("Exports transactions SDK query error: %v", wrappedErr))
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
				logging.LogBackend(fmt.Sprintf("Exports transfers SDK query error: %v", wrappedErr))
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

func (c *ExportsCollection) getWithdrawalsStore() *store.Store[Withdrawal] {
	withdrawalsStoreMu.Lock()
	defer withdrawalsStoreMu.Unlock()

	chain := preferences.GetLastChain()
	address := preferences.GetLastAddress()
	storeKey := getStoreKey(chain, address)
	theStore := withdrawalsStore[storeKey]
	if theStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			// EXISTING_CODE
			exportOpts := sdk.ExportOptions{
				Globals:   sdk.Globals{Cache: true, Verbose: true, Chain: chain},
				RenderCtx: ctx,
				Addrs:     []string{address},
			}
			if _, _, err := exportOpts.ExportWithdrawals(); err != nil {
				wrappedErr := types.NewSDKError("exports", ExportsTransfers, "fetch", err)
				logging.LogBackend(fmt.Sprintf("Exports transfers SDK query error: %v", wrappedErr))
				return wrappedErr
			}
			// EXISTING_CODE
			return nil
		}

		processFunc := func(item interface{}) *Withdrawal {
			// EXISTING_CODE
			// EXISTING_CODE
			if it, ok := item.(*Withdrawal); ok {
				return it
			}
			return nil
		}

		mappingFunc := func(item *Withdrawal) (key interface{}, includeInMap bool) {
			// EXISTING_CODE
			// EXISTING_CODE
			return nil, false
		}

		// EXISTING_CODE
		storeName := fmt.Sprintf("exports_withdrawals_%s_%s", chain, address)
		// EXISTING_CODE
		theStore = store.NewStore(storeName, queryFunc, processFunc, mappingFunc)
		withdrawalsStore[storeKey] = theStore
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
	case ExportsWithdrawals:
		return "exports-withdrawals"
	case ExportsAssets:
		return "exports-assets"
	case ExportsLogs:
		return "exports-logs"
	case ExportsTraces:
		return "exports-traces"
	case ExportsReceipts:
		return "exports-receipts"
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
		return 0, fmt.Errorf("ExportCount query error: %v", err)
	} else if len(countResult) > 0 {
		return int(countResult[0].Count), nil
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

	pl := *payload

	key := store.GetCollectionKey(&pl)
	if collection, exists := collections[key]; exists {
		return collection
	}

	collection := NewExportsCollection()
	collections[key] = collection
	return collection
}

// EXISTING_CODE
func GetExportsCount2(dataFacet string) (int, error) {
	switch types.DataFacet(dataFacet) {
	case ExportsTransactions:
		return getExportsTransactionsCount()
	case ExportsStatements:
		return getExportsStatementsCount()
	case ExportsTransfers:
		return getExportsTransfersCount()
	case ExportsBalances:
		return getExportsBalancesCount()
	case ExportsWithdrawals:
		return getExportsWithdrawalsCount()
	default:
		return 0, fmt.Errorf("unknown dataFacet: %s", dataFacet)
	}
}

func getExportsTransactionsCount() (int, error) {
	chain := preferences.GetLastChain()
	address := preferences.GetLastAddress()
	listOpts := sdk.ListOptions{
		Globals: sdk.Globals{Cache: true, Chain: chain},
		Addrs:   []string{address},
	}

	// Use ExportCount for optimized counting
	if results, _, err := listOpts.ListCount(); err != nil {
		return 0, fmt.Errorf("failed to get exports transactions count: %w", err)
	} else {
		return int(results[0].Count), nil
	}
}

func getExportsStatementsCount() (int, error) {
	return getExportsTransactionsCount()
}

func getExportsTransfersCount() (int, error) {
	return getExportsTransactionsCount()
}

func getExportsBalancesCount() (int, error) {
	return getExportsTransactionsCount()
}

func getExportsWithdrawalsCount() (int, error) {
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
	case ExportsWithdrawals:
		withdrawalsStoreMu.Lock()
		if store, exists := withdrawalsStore[storeKey]; exists {
			store.Reset()
		}
		withdrawalsStoreMu.Unlock()
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

	withdrawalsStoreMu.Lock()
	delete(withdrawalsStore, storeKey)
	withdrawalsStoreMu.Unlock()
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

	withdrawalsStoreMu.Lock()
	withdrawalsStore = make(map[string]*store.Store[Withdrawal])
	withdrawalsStoreMu.Unlock()
}

func getStoreKey(chain, address string) string {
	return fmt.Sprintf("%s_%s", chain, address)
}

// EXISTING_CODE
