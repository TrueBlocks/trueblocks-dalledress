// Copyright 2016, 2025 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package abis

import (
	// EXISTING_CODE
	"fmt"
	"strings"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/facets"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
	// EXISTING_CODE
)

// TODO: The slices should be slices to pointers
type AbisPage struct {
	Facet         types.DataFacet `json:"facet"`
	Abis          []Abi           `json:"abis"`
	Functions     []Function      `json:"functions"`
	TotalItems    int             `json:"totalItems"`
	ExpectedTotal int             `json:"expectedTotal"`
	IsFetching    bool            `json:"isFetching"`
	State         types.LoadState `json:"state"`
}

func (p *AbisPage) GetFacet() types.DataFacet {
	return p.Facet
}

func (p *AbisPage) GetTotalItems() int {
	return p.TotalItems
}

func (p *AbisPage) GetExpectedTotal() int {
	return p.ExpectedTotal
}

func (p *AbisPage) GetIsFetching() bool {
	return p.IsFetching
}

func (p *AbisPage) GetState() types.LoadState {
	return p.State
}

// EXISTING_CODE
func (c *AbisCollection) GetPage(
	dataFacet types.DataFacet,
	first, pageSize int,
	sortSpec sdk.SortSpec,
	filter string,
) (types.Page, error) {
	page := &AbisPage{
		Facet: dataFacet,
	}
	filter = strings.ToLower(filter)

	var listFacet *facets.Facet[Abi]
	var detailFacet *facets.Facet[Function]

	switch dataFacet {
	case AbisDownloaded:
		listFacet = c.downloadedFacet
	case AbisKnown:
		listFacet = c.knownFacet
	case AbisFunctions:
		detailFacet = c.functionsFacet
	case AbisEvents:
		detailFacet = c.eventsFacet
	default:
		// This is truly a validation error - invalid DataFacet for this collection
		return nil, types.NewValidationError("abis", dataFacet, "GetPage",
			fmt.Errorf("unsupported dataFacet: %v", dataFacet))
	}

	if listFacet != nil {
		var listFilterFunc = func(item *Abi) bool {
			return strings.Contains(strings.ToLower(item.Name), filter)
		}
		var listSortFunc = func(items []Abi, sort sdk.SortSpec) error {
			return sdk.SortAbis(items, sort)
		}
		if result, err := listFacet.GetPage(first, pageSize, listFilterFunc, sortSpec, listSortFunc); err != nil {
			// This is likely an SDK or store error, not a validation error
			return nil, types.NewStoreError("abis", dataFacet, "GetPage", err)
		} else {
			page.Abis, page.TotalItems, page.State = result.Items, result.TotalItems, result.State
		}
		page.IsFetching = listFacet.IsFetching()
		page.ExpectedTotal = listFacet.ExpectedCount()

	} else if detailFacet != nil {
		var detailFilter = func(item *Function) bool {
			return strings.Contains(strings.ToLower(item.Name), filter) ||
				strings.Contains(strings.ToLower(item.Encoding), filter)
		}
		var detailSortFunc = func(items []Function, sort sdk.SortSpec) error {
			return sdk.SortFunctions(items, sort)
		}
		if result, err := detailFacet.GetPage(first, pageSize, detailFilter, sortSpec, detailSortFunc); err != nil {
			// This is likely an SDK or store error, not a validation error
			return nil, types.NewStoreError("abis", dataFacet, "GetPage", err)
		} else {
			page.Functions, page.TotalItems, page.State = result.Items, result.TotalItems, result.State
		}
		page.IsFetching = detailFacet.IsFetching()
		page.ExpectedTotal = detailFacet.ExpectedCount()

	} else {
		// This should not happen since we validated dataFacet above
		return nil, types.NewValidationError("abis", dataFacet, "GetPage",
			fmt.Errorf("no facet found for dataFacet: %v", dataFacet))
	}

	return page, nil
}

// EXISTING_CODE
