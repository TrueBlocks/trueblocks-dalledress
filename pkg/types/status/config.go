// Copyright 2016, 2026 The Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package status

import (
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

// GetConfig returns the ViewConfig for the Status view
func (c *StatusCollection) GetConfig() (*types.ViewConfig, error) {
	facets := map[string]types.FacetConfig{
		"status": {
			Name:          "Status",
			Store:         "status",
			IsForm:        true,
			DividerBefore: false,
			Fields:        getStatusFields(),
			Actions:       []string{},
			HeaderActions: []string{},
		},
		"caches": {
			Name:          "Caches",
			Store:         "caches",
			IsForm:        false,
			DividerBefore: false,
			Fields:        getCachesFields(),
			Actions:       []string{},
			HeaderActions: []string{"export"},
		},
		"chains": {
			Name:          "Chains",
			Store:         "chains",
			IsForm:        false,
			DividerBefore: false,
			Fields:        getChainsFields(),
			Actions:       []string{},
			HeaderActions: []string{"export"},
		},
	}

	cfg := &types.ViewConfig{
		ViewName:   "status",
		Facets:     facets,
		FacetOrder: []string{"status", "caches", "chains"},
		Actions: map[string]types.ActionConfig{
			"export": {Name: "export", Label: "Export", Icon: "Export"},
		},
	}
	types.DeriveFacetFromFields(cfg)
	types.NormalizeOrders(cfg)
	return cfg, nil
}

func getCachesFields() []types.FieldConfig {
	return []types.FieldConfig{
		// EXISTING_CODE
		{Key: "type", Label: "Type", Section: "General", InTable: true, InDetail: true, Width: 100, Order: 1, DetailOrder: 1},
		{Key: "path", Label: "Path", Section: "General", InTable: true, InDetail: true, Width: 300, Order: 2, DetailOrder: 2},
		{Key: "nFiles", Label: "Files", Section: "Statistics", InTable: true, InDetail: true, Width: 100, Order: 3, DetailOrder: 3},
		{Key: "nFolders", Label: "Folders", Section: "Statistics", InTable: true, InDetail: true, Width: 100, Order: 4, DetailOrder: 4},
		{Key: "sizeInBytes", Label: "Size (Bytes)", Section: "Statistics", InTable: true, InDetail: true, Width: 150, Order: 5, DetailOrder: 5},
		{Key: "lastCached", Label: "Last Cached", Section: "Timestamps", InTable: true, InDetail: true, Width: 150, Order: 6, DetailOrder: 6},
		// EXISTING_CODE
	}
}

func getChainsFields() []types.FieldConfig {
	return []types.FieldConfig{
		// EXISTING_CODE
		{Key: "chain", Label: "Chain", Section: "General", InTable: true, InDetail: true, Width: 100, Order: 1, DetailOrder: 1},
		{Key: "chainId", Label: "Chain ID", Section: "General", InTable: true, InDetail: true, Width: 100, Order: 2, DetailOrder: 2},
		{Key: "symbol", Label: "Symbol", Section: "General", InTable: true, InDetail: true, Width: 100, Order: 3, DetailOrder: 3},
		{Key: "rpcProvider", Label: "RPC Provider", Section: "Providers", InTable: true, InDetail: true, Width: 200, Order: 4, DetailOrder: 4},
		{Key: "ipfsGateway", Label: "IPFS Gateway", Section: "Providers", InTable: true, InDetail: true, Width: 200, Order: 5, DetailOrder: 5},
		{Key: "localExplorer", Label: "Local Explorer", Section: "Explorers", InTable: true, InDetail: true, Width: 200, Order: 6, DetailOrder: 6},
		{Key: "remoteExplorer", Label: "Remote Explorer", Section: "Explorers", InTable: true, InDetail: true, Width: 200, Order: 7, DetailOrder: 7},
		// EXISTING_CODE
	}
}

func getStatusFields() []types.FieldConfig {
	return []types.FieldConfig{
		// EXISTING_CODE
		{Key: "cachePath", Label: "Cache Path", Section: "Paths", InTable: true, InDetail: true, Width: 200, Order: 1, DetailOrder: 1},
		{Key: "indexPath", Label: "Index Path", Section: "Paths", InTable: true, InDetail: true, Width: 200, Order: 2, DetailOrder: 2},
		{Key: "chain", Label: "Chain", Section: "Chain", InTable: true, InDetail: true, Width: 100, Order: 3, DetailOrder: 3},
		{Key: "chainId", Label: "Chain ID", Section: "Chain", InTable: true, InDetail: true, Width: 100, Order: 4, DetailOrder: 4},
		{Key: "networkId", Label: "Network ID", Section: "Chain", InTable: true, InDetail: true, Width: 100, Order: 5, DetailOrder: 5},
		{Key: "chainConfig", Label: "Chain Config", Section: "Chain", InTable: true, InDetail: true, Width: 150, Order: 6, DetailOrder: 6},
		{Key: "rootConfig", Label: "Root Config", Section: "Config", InTable: true, InDetail: true, Width: 150, Order: 7, DetailOrder: 7},
		{Key: "clientVersion", Label: "Client Version", Section: "Config", InTable: true, InDetail: true, Width: 150, Order: 8, DetailOrder: 8},
		{Key: "version", Label: "Version", Section: "Config", InTable: true, InDetail: true, Width: 100, Order: 9, DetailOrder: 9},
		{Key: "progress", Label: "Progress", Section: "Progress", InTable: true, InDetail: true, Width: 100, Order: 10, DetailOrder: 10},
		{Key: "rpcProvider", Label: "RPC Provider", Section: "Providers", InTable: true, InDetail: true, Width: 200, Order: 11, DetailOrder: 11},
		{Key: "hasEsKey", Label: "Has ES Key", Section: "Flags", InTable: true, InDetail: true, Width: 100, Order: 12, DetailOrder: 12},
		{Key: "hasPinKey", Label: "Has Pin Key", Section: "Flags", InTable: true, InDetail: true, Width: 100, Order: 13, DetailOrder: 13},
		{Key: "isApi", Label: "Is API", Section: "Flags", InTable: true, InDetail: true, Width: 100, Order: 14, DetailOrder: 14},
		{Key: "isArchive", Label: "Is Archive", Section: "Flags", InTable: true, InDetail: true, Width: 100, Order: 15, DetailOrder: 15},
		{Key: "isScraping", Label: "Is Scraping", Section: "Flags", InTable: true, InDetail: true, Width: 100, Order: 16, DetailOrder: 16},
		{Key: "isTesting", Label: "Is Testing", Section: "Flags", InTable: true, InDetail: true, Width: 100, Order: 17, DetailOrder: 17},
		{Key: "isTracing", Label: "Is Tracing", Section: "Flags", InTable: true, InDetail: true, Width: 100, Order: 18, DetailOrder: 18},
		// EXISTING_CODE
	}
}

// EXISTING_CODE
// EXISTING_CODE
