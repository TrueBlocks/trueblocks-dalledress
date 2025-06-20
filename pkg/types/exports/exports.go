package exports

import (
	"fmt"
	"strings"
	"sync"

	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/facets"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/logging"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

const (
	ExportsStatements   types.ListKind = "statements"
	ExportsTransfers    types.ListKind = "transfers"
	ExportsBalances     types.ListKind = "balances"
	ExportsTransactions types.ListKind = "transactions"
)

const (
	DataFacetStatements   types.DataFacet = "statements"
	DataFacetTransfers    types.DataFacet = "transfers"
	DataFacetBalances     types.DataFacet = "balances"
	DataFacetTransactions types.DataFacet = "transactions"
)

func init() {
	types.RegisterKind(ExportsStatements)
	types.RegisterKind(ExportsTransfers)
	types.RegisterKind(ExportsBalances)
	types.RegisterKind(ExportsTransactions)

	types.RegisterDataFacet(DataFacetStatements)
	types.RegisterDataFacet(DataFacetTransfers)
	types.RegisterDataFacet(DataFacetBalances)
	types.RegisterDataFacet(DataFacetTransactions)
}

type ExportsPage struct {
	Kind          types.ListKind          `json:"kind"`
	Statements    []coreTypes.Statement   `json:"statements,omitempty"`
	Transfers     []coreTypes.Transfer    `json:"transfers,omitempty"`
	Balances      []coreTypes.Token       `json:"balances,omitempty"`
	Transactions  []coreTypes.Transaction `json:"transactions,omitempty"`
	TotalItems    int                     `json:"totalItems"`
	ExpectedTotal int                     `json:"expectedTotal"`
	IsFetching    bool                    `json:"isFetching"`
	State         types.LoadState         `json:"state"`
}

func (ep *ExportsPage) GetKind() types.ListKind   { return ep.Kind }
func (ep *ExportsPage) GetTotalItems() int        { return ep.TotalItems }
func (ep *ExportsPage) GetExpectedTotal() int     { return ep.ExpectedTotal }
func (ep *ExportsPage) GetIsFetching() bool       { return ep.IsFetching }
func (ep *ExportsPage) GetState() types.LoadState { return ep.State }

type ExportsCollection struct {
	chain             string
	address           string
	statementsFacet   *facets.Facet[coreTypes.Statement]
	transfersFacet    *facets.Facet[coreTypes.Transfer]
	balancesFacet     *facets.Facet[coreTypes.Token]
	transactionsFacet *facets.Facet[coreTypes.Transaction]
	summary           types.Summary
	summaryMutex      sync.RWMutex
}

func NewExportsCollection(chain, address string) *ExportsCollection {
	instance := &ExportsCollection{
		chain:   chain,
		address: address,
		summary: types.Summary{
			TotalCount:  0,
			FacetCounts: make(map[types.ListKind]int),
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

func (ec *ExportsCollection) GetPage(
	listKind types.ListKind,
	first, pageSize int,
	sortSpec sdk.SortSpec,
	filter string,
) (types.Page, error) {
	switch listKind {
	case ExportsStatements:
		return ec.getStatementsPage(first, pageSize, sortSpec, filter)
	case ExportsTransfers:
		return ec.getTransfersPage(first, pageSize, sortSpec, filter)
	case ExportsBalances:
		return ec.getBalancesPage(first, pageSize, sortSpec, filter)
	case ExportsTransactions:
		return ec.getTransactionsPage(first, pageSize, sortSpec, filter)
	default:
		return nil, fmt.Errorf("GetPage: unexpected list kind: %v", listKind)
	}
}

func (ec *ExportsCollection) getStatementsPage(first, pageSize int, sortSpec sdk.SortSpec, filter string) (*ExportsPage, error) {
	page := &ExportsPage{
		Kind: ExportsStatements,
	}
	filter = strings.ToLower(filter)

	var filterFunc = func(item *coreTypes.Statement) bool {
		if filter == "" {
			return true
		}
		// Filter based on statement fields
		return strings.Contains(strings.ToLower(item.AccountedFor.Hex()), filter) ||
			strings.Contains(strings.ToLower(item.Asset.Hex()), filter)
	}

	var sortFunc = func(items []coreTypes.Statement, sort sdk.SortSpec) error {
		// TODO: Implement proper sorting when SDK methods are available
		return nil
	}

	if result, err := ec.statementsFacet.GetPage(first, pageSize, filterFunc, sortSpec, sortFunc); err != nil {
		return nil, types.NewStoreError("exports", ExportsStatements, "GetPage", err)
	} else {
		page.Statements, page.TotalItems, page.State = result.Items, result.TotalItems, result.State
	}

	page.IsFetching = ec.statementsFacet.IsFetching()
	page.ExpectedTotal = ec.statementsFacet.ExpectedCount()
	return page, nil
}

func (ec *ExportsCollection) getTransfersPage(first, pageSize int, sortSpec sdk.SortSpec, filter string) (*ExportsPage, error) {
	page := &ExportsPage{
		Kind: ExportsTransfers,
	}
	filter = strings.ToLower(filter)

	var filterFunc = func(item *coreTypes.Transfer) bool {
		if filter == "" {
			return true
		}
		// Filter based on transfer fields
		return strings.Contains(strings.ToLower(item.Asset.Hex()), filter) ||
			strings.Contains(strings.ToLower(item.Sender.Hex()), filter) ||
			strings.Contains(strings.ToLower(item.Recipient.Hex()), filter)
	}

	var sortFunc = func(items []coreTypes.Transfer, sort sdk.SortSpec) error {
		// TODO: Implement proper sorting when SDK methods are available
		return nil
	}

	if result, err := ec.transfersFacet.GetPage(first, pageSize, filterFunc, sortSpec, sortFunc); err != nil {
		return nil, types.NewStoreError("exports", ExportsTransfers, "GetPage", err)
	} else {
		page.Transfers, page.TotalItems, page.State = result.Items, result.TotalItems, result.State
	}

	page.IsFetching = ec.transfersFacet.IsFetching()
	page.ExpectedTotal = ec.transfersFacet.ExpectedCount()
	return page, nil
}

func (ec *ExportsCollection) getBalancesPage(first, pageSize int, sortSpec sdk.SortSpec, filter string) (*ExportsPage, error) {
	page := &ExportsPage{
		Kind: ExportsBalances,
	}
	filter = strings.ToLower(filter)

	var filterFunc = func(item *coreTypes.Token) bool {
		if filter == "" {
			return true
		}
		// Filter based on balance/state fields
		return strings.Contains(strings.ToLower(item.Address.Hex()), filter)
	}

	var sortFunc = func(items []coreTypes.Token, sort sdk.SortSpec) error {
		// TODO: Implement proper sorting when SDK methods are available
		return nil
	}

	if result, err := ec.balancesFacet.GetPage(first, pageSize, filterFunc, sortSpec, sortFunc); err != nil {
		return nil, types.NewStoreError("exports", ExportsBalances, "GetPage", err)
	} else {
		page.Balances, page.TotalItems, page.State = result.Items, result.TotalItems, result.State
	}

	page.IsFetching = ec.balancesFacet.IsFetching()
	page.ExpectedTotal = ec.balancesFacet.ExpectedCount()
	return page, nil
}

func (ec *ExportsCollection) getTransactionsPage(first, pageSize int, sortSpec sdk.SortSpec, filter string) (*ExportsPage, error) {
	page := &ExportsPage{
		Kind: ExportsTransactions,
	}
	filter = strings.ToLower(filter)

	var filterFunc = func(item *coreTypes.Transaction) bool {
		if filter == "" {
			return true
		}
		// TODO: Implement proper filtering based on transaction fields
		return strings.Contains(strings.ToLower(item.Hash.Hex()), filter)
	}

	var sortFunc = func(items []coreTypes.Transaction, sort sdk.SortSpec) error {
		// TODO: Implement proper sorting when SDK methods are available
		return nil
	}

	if result, err := ec.transactionsFacet.GetPage(first, pageSize, filterFunc, sortSpec, sortFunc); err != nil {
		return nil, types.NewStoreError("exports", ExportsTransactions, "GetPage", err)
	} else {
		page.Transactions, page.TotalItems, page.State = result.Items, result.TotalItems, result.State
	}

	page.IsFetching = ec.transactionsFacet.IsFetching()
	page.ExpectedTotal = ec.transactionsFacet.ExpectedCount()
	return page, nil
}

func (ec *ExportsCollection) LoadData(listKind types.ListKind) {
	if !ec.NeedsUpdate(listKind) {
		return
	}

	var facetName string

	switch listKind {
	case ExportsStatements:
		facetName = string(ExportsStatements)
	case ExportsTransfers:
		facetName = string(ExportsTransfers)
	case ExportsBalances:
		facetName = string(ExportsBalances)
	case ExportsTransactions:
		facetName = string(ExportsTransactions)
	default:
		logging.LogError("LoadData: unexpected list kind: %v", fmt.Errorf("invalid list kind: %s", listKind), nil)
		return
	}

	go func() {
		// Handle each facet type specifically since they're different types
		switch listKind {
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

func (ec *ExportsCollection) Reset(listKind types.ListKind) {
	switch listKind {
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

func (ec *ExportsCollection) NeedsUpdate(listKind types.ListKind) bool {
	switch listKind {
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

func (ec *ExportsCollection) GetSupportedKinds() []types.ListKind {
	return []types.ListKind{
		ExportsStatements,
		ExportsTransfers,
		ExportsBalances,
		ExportsTransactions,
	}
}

func (ec *ExportsCollection) GetStoreForKind(kind types.ListKind) string {
	switch kind {
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
		summary.FacetCounts = make(map[types.ListKind]int)
	}

	switch v := item.(type) {
	case *coreTypes.Statement:
		summary.TotalCount++
		summary.FacetCounts[ExportsStatements]++
		if summary.CustomData == nil {
			summary.CustomData = make(map[string]interface{})
		}

		stmtCount, _ := summary.CustomData["statementsCount"].(int)
		stmtCount++
		summary.CustomData["statementsCount"] = stmtCount

	case *coreTypes.Transfer:
		summary.TotalCount++
		summary.FacetCounts[ExportsTransfers]++
		if summary.CustomData == nil {
			summary.CustomData = make(map[string]interface{})
		}

		transferCount, _ := summary.CustomData["transfersCount"].(int)
		transferCount++
		summary.CustomData["transfersCount"] = transferCount

	case *coreTypes.Token:
		summary.TotalCount++
		summary.FacetCounts[ExportsBalances]++
		if summary.CustomData == nil {
			summary.CustomData = make(map[string]interface{})
		}

		balanceCount, _ := summary.CustomData["balancesCount"].(int)
		balanceCount++
		summary.CustomData["balancesCount"] = balanceCount

	case *coreTypes.Transaction:
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

func (ec *ExportsCollection) GetCurrentSummary() types.Summary {
	ec.summaryMutex.RLock()
	defer ec.summaryMutex.RUnlock()

	summary := ec.summary
	summary.FacetCounts = make(map[types.ListKind]int)
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
		FacetCounts: make(map[types.ListKind]int),
		CustomData:  make(map[string]interface{}),
		LastUpdated: 0,
	}
}

func (ec *ExportsCollection) GetSummary() types.Summary {
	return ec.GetCurrentSummary()
}
