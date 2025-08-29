// Copyright 2016, 2026 The Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package dalledress

import (
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

// GetConfig returns the ViewConfig for the DalleDress view
func (c *DalleDressCollection) GetConfig() (*types.ViewConfig, error) {
	facets := map[string]types.FacetConfig{
		"generator": {
			Name:          "Generator",
			Store:         "dalledresses",
			IsForm:        true,
			DividerBefore: false,
			Fields:        getDalleDressFields(),
			Actions:       []string{},
			HeaderActions: []string{"speak"},
			RendererTypes: "facet",
		},
		"series": {
			Name:          "Series",
			Store:         "series",
			IsForm:        false,
			DividerBefore: false,
			Fields:        getSeriesFields(),
			Actions:       []string{"edit", "delete", "duplicate"},
			HeaderActions: []string{"create", "export"},
			RendererTypes: "",
		},
		"databases": {
			Name:          "Databases",
			Store:         "databases",
			IsForm:        false,
			DividerBefore: false,
			Fields:        getDatabasesFields(),
			Actions:       []string{},
			HeaderActions: []string{"export"},
			RendererTypes: "",
		},
		"events": {
			Name:          "Events",
			Store:         "logs",
			IsForm:        false,
			DividerBefore: false,
			Fields:        getLogsFields(),
			Actions:       []string{},
			HeaderActions: []string{"export"},
			RendererTypes: "",
		},
		"gallery": {
			Name:          "Gallery",
			Store:         "dalledresses",
			IsForm:        true,
			DividerBefore: false,
			Fields:        getDalleDressFields(),
			Actions:       []string{},
			HeaderActions: []string{},
			RendererTypes: "facet",
		},
	}

	cfg := &types.ViewConfig{
		ViewName:   "dalledress",
		Facets:     facets,
		FacetOrder: []string{"generator", "series", "databases", "events", "gallery"},
		Actions: map[string]types.ActionConfig{
			"create":    {Name: "create", Label: "Create", Icon: "Add"},
			"export":    {Name: "export", Label: "Export", Icon: "Export"},
			"edit":      {Name: "edit", Label: "Edit", Icon: "Edit"},
			"delete":    {Name: "delete", Label: "Delete", Icon: "Delete"},
			"duplicate": {Name: "duplicate", Label: "Duplicate", Icon: "Duplicate"},
			"speak":     {Name: "speak", Label: "Speak", Icon: "Speak"},
		},
	}
	types.DeriveFacets(cfg)
	types.NormalizeOrders(cfg)
	return cfg, nil
}

func getLogsFields() []types.FieldConfig {
	return []types.FieldConfig{
		// EXISTING_CODE
		{Key: "date", Label: "Date", ColumnLabel: "Date", DetailLabel: "Date", Section: "Block/Tx", Width: 120, Order: 1, DetailOrder: 1},
		{Key: "address", Label: "Address", ColumnLabel: "Address", DetailLabel: "Address", Section: "Event", Width: 340, Order: 2, DetailOrder: 5},
		{Key: "name", Label: "Event Name", ColumnLabel: "Name", DetailLabel: "Event Name", Section: "Event", Width: 200, Order: 3, DetailOrder: 6},
		{Key: "articulatedLog", Label: "Articulated Log", ColumnLabel: "Log Details", DetailLabel: "Articulated Log", Section: "Event", Width: 120, Order: 4, DetailOrder: 8},
		{Key: "blockNumber", Label: "Block Number", ColumnLabel: "", DetailLabel: "Block Number", Section: "Block/Tx", NoTable: true, DetailOrder: 2},
		{Key: "transactionIndex", Label: "Transaction Index", ColumnLabel: "", DetailLabel: "Transaction Index", Section: "Block/Tx", NoTable: true, DetailOrder: 3},
		{Key: "transactionHash", Label: "Transaction Hash", ColumnLabel: "", DetailLabel: "Transaction Hash", Section: "Block/Tx", NoTable: true, DetailOrder: 4},
		{Key: "signature", Label: "Signature", ColumnLabel: "", DetailLabel: "Signature", Section: "Event", NoTable: true, DetailOrder: 7},
		// EXISTING_CODE
	}
}

func getDalleDressFields() []types.FieldConfig {
	return []types.FieldConfig{
		// MISSING
		// EXISTING_CODE
		// EXISTING_CODE
	}
}

func getDatabasesFields() []types.FieldConfig {
	return []types.FieldConfig{
		// EXISTING_CODE
		{Key: "databaseName", Label: "Database Name", ColumnLabel: "Database", DetailLabel: "Database Name", Section: "General", Width: 200, Order: 1, DetailOrder: 1},
		{Key: "count", Label: "Count", ColumnLabel: "Count", DetailLabel: "Count", Section: "General", Width: 100, Order: 2, DetailOrder: 2},
		{Key: "sample", Label: "Sample", ColumnLabel: "Sample", DetailLabel: "Sample", Section: "General", Width: 300, Order: 3, DetailOrder: 3},
		{Key: "filtered", Label: "Filtered", ColumnLabel: "Filtered", DetailLabel: "Filtered", Section: "General", Width: 80, Order: 4, DetailOrder: 4},
		// EXISTING_CODE
	}
}

func getSeriesFields() []types.FieldConfig {
	return []types.FieldConfig{
		// EXISTING_CODE
		{Key: "suffix", Label: "Series Name", ColumnLabel: "Series", DetailLabel: "Series Name", Section: "General", Width: 150, Order: 1, DetailOrder: 1},
		{Key: "last", Label: "Last Index", ColumnLabel: "Last", DetailLabel: "Last Index", Section: "General", Width: 80, Order: 2, DetailOrder: 2},
		{Key: "adverbs", Label: "Adverbs", ColumnLabel: "", DetailLabel: "Adverbs", Section: "Content", NoTable: true, DetailOrder: 3},
		{Key: "adjectives", Label: "Adjectives", ColumnLabel: "", DetailLabel: "Adjectives", Section: "Content", NoTable: true, DetailOrder: 4},
		{Key: "nouns", Label: "Nouns", ColumnLabel: "", DetailLabel: "Nouns", Section: "Content", NoTable: true, DetailOrder: 5},
		{Key: "emotions", Label: "Emotions", ColumnLabel: "", DetailLabel: "Emotions", Section: "Content", NoTable: true, DetailOrder: 6},
		{Key: "artstyles", Label: "Art Styles", ColumnLabel: "", DetailLabel: "Art Styles", Section: "Style", NoTable: true, DetailOrder: 7},
		{Key: "colors", Label: "Colors", ColumnLabel: "", DetailLabel: "Colors", Section: "Style", NoTable: true, DetailOrder: 8},
		{Key: "orientations", Label: "Orientations", ColumnLabel: "", DetailLabel: "Orientations", Section: "Style", NoTable: true, DetailOrder: 9},
		{Key: "gazes", Label: "Gazes", ColumnLabel: "", DetailLabel: "Gazes", Section: "Style", NoTable: true, DetailOrder: 10},
		{Key: "backstyles", Label: "Back Styles", ColumnLabel: "", DetailLabel: "Back Styles", Section: "Style", NoTable: true, DetailOrder: 11},
		{Key: "modifiedAt", Label: "Modified At", ColumnLabel: "Modified", DetailLabel: "Modified At", Section: "General", Width: 120, Order: 3, DetailOrder: 12},
		// EXISTING_CODE
	}
}

// EXISTING_CODE
// EXISTING_CODE
