// Copyright 2016, 2026 The Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package contracts

import (
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

// GetConfig returns the ViewConfig for the Contracts view
func (c *ContractsCollection) GetConfig() (*types.ViewConfig, error) {
	facets := map[string]types.FacetConfig{
		"dashboard": {
			Name:          "Dashboard",
			Store:         "contracts",
			IsForm:        true,
			DividerBefore: false,
			Fields:        getContractsFields(),
			Actions:       []string{},
			HeaderActions: []string{},
		},
		"execute": {
			Name:          "Execute",
			Store:         "contracts",
			IsForm:        true,
			DividerBefore: false,
			Fields:        getContractsFields(),
			Actions:       []string{},
			HeaderActions: []string{},
		},
		"events": {
			Name:          "Events",
			Store:         "logs",
			IsForm:        false,
			DividerBefore: false,
			Fields:        getLogsFields(),
			Actions:       []string{},
			HeaderActions: []string{"export"},
		},
	}

	cfg := &types.ViewConfig{
		ViewName:   "contracts",
		Facets:     facets,
		FacetOrder: []string{"dashboard", "execute", "events"},
		Actions: map[string]types.ActionConfig{
			"export": {Name: "export", Label: "Export", Icon: "Export"},
		},
	}
	types.DeriveFacetFromFields(cfg)
	types.NormalizeOrders(cfg)
	return cfg, nil
}

func getContractsFields() []types.FieldConfig {
	return []types.FieldConfig{
		// EXISTING_CODE
		{Key: "address", Label: "Address", ColumnLabel: "Address", DetailLabel: "Address", Section: "General", InTable: true, InDetail: true, Width: 240, Order: 1, DetailOrder: 1},
		{Key: "name", Label: "Name", ColumnLabel: "Name", DetailLabel: "Name", Section: "General", InTable: true, InDetail: true, Width: 200, Order: 2, DetailOrder: 2},
		{Key: "symbol", Label: "Symbol", ColumnLabel: "Symbol", DetailLabel: "Symbol", Section: "General", InTable: true, InDetail: true, Width: 100, Order: 3, DetailOrder: 3},
		{Key: "decimals", Label: "Decimals", ColumnLabel: "Decimals", DetailLabel: "Decimals", Section: "General", InTable: true, InDetail: true, Width: 100, Order: 4, DetailOrder: 4},
		{Key: "source", Label: "Source", ColumnLabel: "Source", DetailLabel: "Source", Section: "Source", InTable: true, InDetail: true, Width: 150, Order: 5, DetailOrder: 5},
		// EXISTING_CODE
	}
}

func getLogsFields() []types.FieldConfig {
	return []types.FieldConfig{
		// EXISTING_CODE
		{Key: "date", Label: "Date", ColumnLabel: "Date", DetailLabel: "Date", Section: "Block/Tx", InTable: true, InDetail: true, Width: 120, Order: 1, DetailOrder: 1},
		{Key: "address", Label: "Address", ColumnLabel: "Address", DetailLabel: "Address", Section: "Event", InTable: true, InDetail: true, Width: 340, Order: 2, DetailOrder: 5},
		{Key: "name", Label: "Event Name", ColumnLabel: "Name", DetailLabel: "Event Name", Section: "Event", InTable: true, InDetail: true, Width: 200, Order: 3, DetailOrder: 6},
		{Key: "articulatedLog", Label: "Articulated Log", ColumnLabel: "Log Details", DetailLabel: "Articulated Log", Section: "Event", InTable: true, InDetail: true, Width: 120, Order: 4, DetailOrder: 8},
		{Key: "blockNumber", Label: "Block Number", ColumnLabel: "", DetailLabel: "Block Number", Section: "Block/Tx", InTable: false, InDetail: true, DetailOrder: 2},
		{Key: "transactionIndex", Label: "Transaction Index", ColumnLabel: "", DetailLabel: "Transaction Index", Section: "Block/Tx", InTable: false, InDetail: true, DetailOrder: 3},
		{Key: "transactionHash", Label: "Transaction Hash", ColumnLabel: "", DetailLabel: "Transaction Hash", Section: "Block/Tx", InTable: false, InDetail: true, DetailOrder: 4},
		{Key: "signature", Label: "Signature", ColumnLabel: "", DetailLabel: "Signature", Section: "Event", InTable: false, InDetail: true, DetailOrder: 7},
		// EXISTING_CODE
	}
}

// EXISTING_CODE
// EXISTING_CODE
