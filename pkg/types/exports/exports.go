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
	"time"

	// EXISTING_CODE

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/facets"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/logging"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	// EXISTING_CODE
)

const (
	ExportsStatements   types.DataFacet = "statements"
	ExportsBalances     types.DataFacet = "balances"
	ExportsTransfers    types.DataFacet = "transfers"
	ExportsTransactions types.DataFacet = "transactions"
)

func init() {
	types.RegisterDataFacet(ExportsStatements)
	types.RegisterDataFacet(ExportsBalances)
	types.RegisterDataFacet(ExportsTransfers)
	types.RegisterDataFacet(ExportsTransactions)
}

type ExportsCollection struct {
	address           string
	statementsFacet   *facets.Facet[Statement]
	balancesFacet     *facets.Facet[Balance]
	transfersFacet    *facets.Facet[Transfer]
	transactionsFacet *facets.Facet[Transaction]
	summary           types.Summary
	summaryMutex      sync.RWMutex
}

func NewExportsCollection(address string) *ExportsCollection {
	c := &ExportsCollection{
		address: address,
	}
	c.ResetSummary()
	c.initializeFacets()
	return c
}

func (c *ExportsCollection) initializeFacets() {
	c.statementsFacet = facets.NewFacetWithSummary(
		ExportsStatements,
		isStatement,
		isDupStatement(),
		c.getStatementsStore(),
		"exports",
		c,
	)

	c.balancesFacet = facets.NewFacetWithSummary(
		ExportsBalances,
		isBalance,
		isDupBalance(),
		c.getBalancesStore(),
		"exports",
		c,
	)

	c.transfersFacet = facets.NewFacetWithSummary(
		ExportsTransfers,
		isTransfer,
		isDupTransfer(),
		c.getTransfersStore(),
		"exports",
		c,
	)

	c.transactionsFacet = facets.NewFacetWithSummary(
		ExportsTransactions,
		isTransaction,
		isDupTransaction(),
		c.getTransactionsStore(),
		"exports",
		c,
	)
}

func isStatement(item *Statement) bool {
	// EXISTING_CODE
	return true
	// EXISTING_CODE
}

func isBalance(item *Balance) bool {
	// EXISTING_CODE
	return true
	// EXISTING_CODE
}

func isTransfer(item *Transfer) bool {
	// EXISTING_CODE
	return true
	// EXISTING_CODE
}

func isTransaction(item *Transaction) bool {
	// EXISTING_CODE
	return true
	// EXISTING_CODE
}

func isDupBalance() func(existing []*Balance, newItem *Balance) bool {
	// EXISTING_CODE
	return nil
	// EXISTING_CODE
}

func isDupStatement() func(existing []*Statement, newItem *Statement) bool {
	// EXISTING_CODE
	return nil
	// EXISTING_CODE
}

func isDupTransaction() func(existing []*Transaction, newItem *Transaction) bool {
	// EXISTING_CODE
	return nil
	// EXISTING_CODE
}

func isDupTransfer() func(existing []*Transfer, newItem *Transfer) bool {
	// EXISTING_CODE
	return nil
	// EXISTING_CODE
}

func (c *ExportsCollection) LoadData(dataFacet types.DataFacet) {
	if !c.NeedsUpdate(dataFacet) {
		return
	}

	go func() {
		switch dataFacet {
		case ExportsStatements:
			if err := c.statementsFacet.Load(); err != nil {
				logging.LogError(fmt.Sprintf("LoadData.%s from store: %%v", dataFacet), err, facets.ErrAlreadyLoading)
			}
		case ExportsBalances:
			if err := c.balancesFacet.Load(); err != nil {
				logging.LogError(fmt.Sprintf("LoadData.%s from store: %%v", dataFacet), err, facets.ErrAlreadyLoading)
			}
		case ExportsTransfers:
			if err := c.transfersFacet.Load(); err != nil {
				logging.LogError(fmt.Sprintf("LoadData.%s from store: %%v", dataFacet), err, facets.ErrAlreadyLoading)
			}
		case ExportsTransactions:
			if err := c.transactionsFacet.Load(); err != nil {
				logging.LogError(fmt.Sprintf("LoadData.%s from store: %%v", dataFacet), err, facets.ErrAlreadyLoading)
			}
		default:
			logging.LogError("LoadData: unexpected dataFacet: %v", fmt.Errorf("invalid dataFacet: %s", dataFacet), nil)
			return
		}
	}()
}

func (c *ExportsCollection) Reset(dataFacet types.DataFacet) {
	switch dataFacet {
	case ExportsStatements:
		c.statementsFacet.GetStore().Reset()
	case ExportsBalances:
		c.balancesFacet.GetStore().Reset()
	case ExportsTransfers:
		c.transfersFacet.GetStore().Reset()
	case ExportsTransactions:
		c.transactionsFacet.GetStore().Reset()
	default:
		return
	}
}

func (c *ExportsCollection) NeedsUpdate(dataFacet types.DataFacet) bool {
	switch dataFacet {
	case ExportsStatements:
		return c.statementsFacet.NeedsUpdate()
	case ExportsBalances:
		return c.balancesFacet.NeedsUpdate()
	case ExportsTransfers:
		return c.transfersFacet.NeedsUpdate()
	case ExportsTransactions:
		return c.transactionsFacet.NeedsUpdate()
	default:
		return false
	}
}

func (c *ExportsCollection) GetSupportedFacets() []types.DataFacet {
	return []types.DataFacet{
		ExportsStatements,
		ExportsBalances,
		ExportsTransfers,
		ExportsTransactions,
	}
}

func (c *ExportsCollection) AccumulateItem(item interface{}, summary *types.Summary) {
	// EXISTING_CODE
	c.summaryMutex.Lock()
	defer c.summaryMutex.Unlock()

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

	case *Balance:
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
	// EXISTING_CODE
}

func (c *ExportsCollection) GetSummary() types.Summary {
	c.summaryMutex.RLock()
	defer c.summaryMutex.RUnlock()

	summary := c.summary
	summary.FacetCounts = make(map[types.DataFacet]int)
	for k, v := range c.summary.FacetCounts {
		summary.FacetCounts[k] = v
	}

	if c.summary.CustomData != nil {
		summary.CustomData = make(map[string]interface{})
		for k, v := range c.summary.CustomData {
			summary.CustomData[k] = v
		}
	}

	return summary
}

func (c *ExportsCollection) ResetSummary() {
	c.summaryMutex.Lock()
	defer c.summaryMutex.Unlock()
	c.summary = types.Summary{
		TotalCount:  0,
		FacetCounts: make(map[types.DataFacet]int),
		CustomData:  make(map[string]interface{}),
		LastUpdated: time.Now().Unix(),
	}
}

// EXISTING_CODE
// EXISTING_CODE
