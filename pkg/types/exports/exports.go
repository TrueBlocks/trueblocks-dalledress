package exports

import (
	"fmt"
	"sync"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/facets"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/logging"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

const (
	ExportsStatements   types.DataFacet = "statements"
	ExportsTransfers    types.DataFacet = "transfers"
	ExportsBalances     types.DataFacet = "balances"
	ExportsTransactions types.DataFacet = "transactions"
)

func init() {
	types.RegisterDataFacet(ExportsStatements)
	types.RegisterDataFacet(ExportsTransfers)
	types.RegisterDataFacet(ExportsBalances)
	types.RegisterDataFacet(ExportsTransactions)
}

type ExportsCollection struct {
	chain             string
	address           string
	statementsFacet   *facets.Facet[Statement]
	transfersFacet    *facets.Facet[Transfer]
	balancesFacet     *facets.Facet[Token]
	transactionsFacet *facets.Facet[Transaction]
	summary           types.Summary
	summaryMutex      sync.RWMutex
}

func NewExportsCollection(chain, address string) *ExportsCollection {
	instance := &ExportsCollection{
		chain:   chain,
		address: address,
		summary: types.Summary{
			TotalCount:  0,
			FacetCounts: make(map[types.DataFacet]int),
			CustomData:  make(map[string]interface{}),
		},
	}
	instance.initializeFacets()
	return instance
}

func (ec *ExportsCollection) initializeFacets() {
	ec.statementsFacet = facets.NewFacetWithSummary(
		ExportsStatements,
		nil,
		nil,
		GetExportsStatementsStore(ec.chain, ec.address),
		"exports",
		ec,
	)
	ec.transfersFacet = facets.NewFacetWithSummary(
		ExportsTransfers,
		nil,
		nil,
		GetExportsTransfersStore(ec.chain, ec.address),
		"exports",
		ec,
	)
	ec.balancesFacet = facets.NewFacetWithSummary(
		ExportsBalances,
		nil,
		nil,
		GetExportsBalancesStore(ec.chain, ec.address),
		"exports",
		ec,
	)
	ec.transactionsFacet = facets.NewFacetWithSummary(
		ExportsTransactions,
		nil, // No filter function for now
		nil, // No deduplication function for now
		GetExportsTransactionsStore(ec.chain, ec.address),
		"exports",
		ec,
	)
}

func (ec *ExportsCollection) LoadData(dataFacet types.DataFacet) {
	if !ec.NeedsUpdate(dataFacet) {
		return
	}

	var facetName string

	switch dataFacet {
	case ExportsStatements:
		facetName = string(ExportsStatements)
	case ExportsTransfers:
		facetName = string(ExportsTransfers)
	case ExportsBalances:
		facetName = string(ExportsBalances)
	case ExportsTransactions:
		facetName = string(ExportsTransactions)
	default:
		logging.LogError("LoadData: unexpected dataFacet: %v", fmt.Errorf("invalid dataFacet: %s", dataFacet), nil)
		return
	}

	go func() {
		// Handle each facet type specifically since they're different types
		switch dataFacet {
		case ExportsStatements:
			if err := ec.statementsFacet.Load(); err != nil {
				logging.LogError(fmt.Sprintf("LoadData.%s from store: %%v", facetName), err, facets.ErrAlreadyLoading)
			}
		case ExportsTransfers:
			if err := ec.transfersFacet.Load(); err != nil {
				logging.LogError(fmt.Sprintf("LoadData.%s from store: %%v", facetName), err, facets.ErrAlreadyLoading)
			}
		case ExportsBalances:
			if err := ec.balancesFacet.Load(); err != nil {
				logging.LogError(fmt.Sprintf("LoadData.%s from store: %%v", facetName), err, facets.ErrAlreadyLoading)
			}
		case ExportsTransactions:
			if err := ec.transactionsFacet.Load(); err != nil {
				logging.LogError(fmt.Sprintf("LoadData.%s from store: %%v", facetName), err, facets.ErrAlreadyLoading)
			}
		}
	}()
}

func (ec *ExportsCollection) Reset(dataFacet types.DataFacet) {
	switch dataFacet {
	case ExportsStatements:
		ResetExportsStore(ec.chain, ec.address, ExportsStatements)
		ec.statementsFacet.Reset()
	case ExportsTransfers:
		ResetExportsStore(ec.chain, ec.address, ExportsTransfers)
		ec.transfersFacet.Reset()
	case ExportsBalances:
		ResetExportsStore(ec.chain, ec.address, ExportsBalances)
		ec.balancesFacet.Reset()
	case ExportsTransactions:
		ResetExportsStore(ec.chain, ec.address, ExportsTransactions)
		ec.transactionsFacet.Reset()
	}
}

func (ec *ExportsCollection) NeedsUpdate(dataFacet types.DataFacet) bool {
	switch dataFacet {
	case ExportsStatements:
		return ec.statementsFacet.NeedsUpdate()
	case ExportsTransfers:
		return ec.transfersFacet.NeedsUpdate()
	case ExportsBalances:
		return ec.balancesFacet.NeedsUpdate()
	case ExportsTransactions:
		return ec.transactionsFacet.NeedsUpdate()
	default:
		return false
	}
}

func (ec *ExportsCollection) GetSupportedFacets() []types.DataFacet {
	return []types.DataFacet{
		ExportsStatements,
		ExportsTransfers,
		ExportsBalances,
		ExportsTransactions,
	}
}

func (ec *ExportsCollection) GetStoreForFacet(dataFacet types.DataFacet) string {
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

func (ec *ExportsCollection) GetCollectionName() string {
	return "exports"
}

var (
	// Collection cache per (chain, address) combination
	collections   = make(map[string]*ExportsCollection)
	collectionsMu sync.Mutex
)

// getCollectionKey creates a unique key for caching collections per (chain, address)
func getCollectionKey(chain, address string) string {
	return fmt.Sprintf("%s_%s", chain, address)
}

// GetExportsCollection returns a singleton collection for the given chain and address
func GetExportsCollection(chain, address string) *ExportsCollection {
	collectionsMu.Lock()
	defer collectionsMu.Unlock()

	key := getCollectionKey(chain, address)
	if collection, exists := collections[key]; exists {
		return collection
	}

	collection := NewExportsCollection(chain, address)
	collections[key] = collection
	return collection
}

// ClearExportsCollection removes a cached collection for the given chain and address
func ClearExportsCollection(chain, address string) {
	collectionsMu.Lock()
	defer collectionsMu.Unlock()

	key := getCollectionKey(chain, address)
	delete(collections, key)
}

// ClearAllExportsCollections removes all cached collections
func ClearAllExportsCollections() {
	collectionsMu.Lock()
	defer collectionsMu.Unlock()

	collections = make(map[string]*ExportsCollection)
}

func (ec *ExportsCollection) AccumulateItem(item interface{}, summary *types.Summary) {
	ec.summaryMutex.Lock()
	defer ec.summaryMutex.Unlock()

	if summary.FacetCounts == nil {
		summary.FacetCounts = make(map[types.DataFacet]int)
	}

	switch v := item.(type) {
	case *Statement:
		summary.TotalCount++
		summary.FacetCounts[ExportsStatements]++
		if summary.CustomData == nil {
			summary.CustomData = make(map[string]interface{})
		}

		stmtCount, _ := summary.CustomData["statementsCount"].(int)
		stmtCount++
		summary.CustomData["statementsCount"] = stmtCount

	case *Transfer:
		summary.TotalCount++
		summary.FacetCounts[ExportsTransfers]++
		if summary.CustomData == nil {
			summary.CustomData = make(map[string]interface{})
		}

		transferCount, _ := summary.CustomData["transfersCount"].(int)
		transferCount++
		summary.CustomData["transfersCount"] = transferCount

	case *Token:
		summary.TotalCount++
		summary.FacetCounts[ExportsBalances]++
		if summary.CustomData == nil {
			summary.CustomData = make(map[string]interface{})
		}

		balanceCount, _ := summary.CustomData["balancesCount"].(int)
		balanceCount++
		summary.CustomData["balancesCount"] = balanceCount

	case *Transaction:
		summary.TotalCount++
		summary.FacetCounts[ExportsTransactions]++
		if summary.CustomData == nil {
			summary.CustomData = make(map[string]interface{})
		}

		txCount, _ := summary.CustomData["transactionsCount"].(int)
		totalValue, _ := summary.CustomData["totalValue"].(int64)
		totalGasUsed, _ := summary.CustomData["totalGasUsed"].(int64)

		txCount++
		totalValue += int64(v.Value.Uint64())
		totalGasUsed += int64(v.Receipt.GasUsed)

		summary.CustomData["transactionsCount"] = txCount
		summary.CustomData["totalValue"] = totalValue
		summary.CustomData["totalGasUsed"] = totalGasUsed

	}
}

func (ec *ExportsCollection) GetSummary() types.Summary {
	ec.summaryMutex.RLock()
	defer ec.summaryMutex.RUnlock()

	summary := ec.summary
	summary.FacetCounts = make(map[types.DataFacet]int)
	for k, v := range ec.summary.FacetCounts {
		summary.FacetCounts[k] = v
	}

	if ec.summary.CustomData != nil {
		summary.CustomData = make(map[string]interface{})
		for k, v := range ec.summary.CustomData {
			summary.CustomData[k] = v
		}
	}

	return summary
}

func (ec *ExportsCollection) ResetSummary() {
	ec.summaryMutex.Lock()
	defer ec.summaryMutex.Unlock()
	ec.summary = types.Summary{
		TotalCount:  0,
		FacetCounts: make(map[types.DataFacet]int),
		CustomData:  make(map[string]interface{}),
		LastUpdated: 0,
	}
}
