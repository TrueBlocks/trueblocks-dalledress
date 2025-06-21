package exports

import (
	"fmt"
	"strings"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

type ExportsPage struct {
	Facet         types.DataFacet `json:"facet"`
	Statements    []Statement     `json:"statements"`
	Transfers     []Transfer      `json:"transfers"`
	Balances      []Token         `json:"balances"`
	Transactions  []Transaction   `json:"transactions"`
	TotalItems    int             `json:"totalItems"`
	ExpectedTotal int             `json:"expectedTotal"`
	IsFetching    bool            `json:"isFetching"`
	State         types.LoadState `json:"state"`
}

func (ep *ExportsPage) GetFacet() types.DataFacet { return ep.Facet }
func (ep *ExportsPage) GetTotalItems() int        { return ep.TotalItems }
func (ep *ExportsPage) GetExpectedTotal() int     { return ep.ExpectedTotal }
func (ep *ExportsPage) GetIsFetching() bool       { return ep.IsFetching }
func (ep *ExportsPage) GetState() types.LoadState { return ep.State }

func (ec *ExportsCollection) GetPage(
	dataFacet types.DataFacet,
	first, pageSize int,
	sortSpec sdk.SortSpec,
	filter string,
) (types.Page, error) {
	switch dataFacet {
	case ExportsStatements:
		return ec.getStatementsPage(first, pageSize, sortSpec, filter)
	case ExportsTransfers:
		return ec.getTransfersPage(first, pageSize, sortSpec, filter)
	case ExportsBalances:
		return ec.getBalancesPage(first, pageSize, sortSpec, filter)
	case ExportsTransactions:
		return ec.getTransactionsPage(first, pageSize, sortSpec, filter)
	default:
		return nil, fmt.Errorf("GetPage: unexpected dataFacet: %v", dataFacet)
	}
}

func (ec *ExportsCollection) getStatementsPage(first, pageSize int, sortSpec sdk.SortSpec, filter string) (*ExportsPage, error) {
	page := &ExportsPage{
		Facet: ExportsStatements,
	}
	filter = strings.ToLower(filter)

	var filterFunc = func(item *Statement) bool {
		if filter == "" {
			return true
		}
		// Filter based on statement fields
		return strings.Contains(strings.ToLower(item.AccountedFor.Hex()), filter) ||
			strings.Contains(strings.ToLower(item.Asset.Hex()), filter)
	}

	var sortFunc = func(items []Statement, sort sdk.SortSpec) error {
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
		Facet: ExportsTransfers,
	}
	filter = strings.ToLower(filter)

	var filterFunc = func(item *Transfer) bool {
		if filter == "" {
			return true
		}
		// Filter based on transfer fields
		return strings.Contains(strings.ToLower(item.Asset.Hex()), filter) ||
			strings.Contains(strings.ToLower(item.Sender.Hex()), filter) ||
			strings.Contains(strings.ToLower(item.Recipient.Hex()), filter)
	}

	var sortFunc = func(items []Transfer, sort sdk.SortSpec) error {
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
		Facet: ExportsBalances,
	}
	filter = strings.ToLower(filter)

	var filterFunc = func(item *Token) bool {
		if filter == "" {
			return true
		}
		// Filter based on balance/state fields
		return strings.Contains(strings.ToLower(item.Address.Hex()), filter)
	}

	var sortFunc = func(items []Token, sort sdk.SortSpec) error {
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
		Facet: ExportsTransactions,
	}
	filter = strings.ToLower(filter)

	var filterFunc = func(item *Transaction) bool {
		if filter == "" {
			return true
		}
		// TODO: Implement proper filtering based on transaction fields
		return strings.Contains(strings.ToLower(item.Hash.Hex()), filter)
	}

	var sortFunc = func(items []Transaction, sort sdk.SortSpec) error {
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

