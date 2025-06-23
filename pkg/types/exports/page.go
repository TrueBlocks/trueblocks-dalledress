// Copyright 2016, 2025 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package exports

import (
	// EXISTING_CODE
	"fmt"
	"strings"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
	// EXISTING_CODE
)

// TODO: The slices should be slices to pointers
type ExportsPage struct {
	Facet         types.DataFacet `json:"facet"`
	Balances      []Balance       `json:"balances"`
	Statements    []Statement     `json:"statements"`
	Transactions  []Transaction   `json:"transactions"`
	Transfers     []Transfer      `json:"transfers"`
	TotalItems    int             `json:"totalItems"`
	ExpectedTotal int             `json:"expectedTotal"`
	IsFetching    bool            `json:"isFetching"`
	State         types.LoadState `json:"state"`
}

func (p *ExportsPage) GetFacet() types.DataFacet {
	return p.Facet
}

func (p *ExportsPage) GetTotalItems() int {
	return p.TotalItems
}

func (p *ExportsPage) GetExpectedTotal() int {
	return p.ExpectedTotal
}

func (p *ExportsPage) GetIsFetching() bool {
	return p.IsFetching
}

func (p *ExportsPage) GetState() types.LoadState {
	return p.State
}

// EXISTING_CODE
func (ec *ExportsCollection) GetPage(
	payload types.Payload,
	first, pageSize int,
	sortSpec sdk.SortSpec,
	filter string,
) (types.Page, error) {
	dataFacet := payload.DataFacet

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

	var filterFunc = func(item *Balance) bool {
		if filter == "" {
			return true
		}
		// Filter based on balance/state fields
		return strings.Contains(strings.ToLower(item.Address.Hex()), filter)
	}

	var sortFunc = func(items []Balance, sort sdk.SortSpec) error {
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

// EXISTING_CODE
