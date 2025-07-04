// Copyright 2016, 2026 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package abis

import (
	"fmt"
	"strings"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
	// EXISTING_CODE
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

func (c *AbisCollection) GetPage(
	payload *types.Payload,
	first, pageSize int,
	sortSpec sdk.SortSpec,
	filter string,
) (types.Page, error) {
	dataFacet := payload.DataFacet

	page := &AbisPage{
		Facet: dataFacet,
	}
	filter = strings.ToLower(filter)

	switch dataFacet {
	case AbisDownloaded:
		facet := c.downloadedFacet
		var filterFunc func(*Abi) bool
		if filter != "" {
			filterFunc = func(item *Abi) bool {
				return c.matchesDownloadedFilter(item, filter)
			}
		}
		sortFunc := func(items []Abi, sort sdk.SortSpec) error {
			return sdk.SortAbis(items, sort)
		}
		if result, err := facet.GetPage(first, pageSize, filterFunc, sortSpec, sortFunc); err != nil {
			return nil, types.NewStoreError("abis", dataFacet, "GetPage", err)
		} else {

			page.Abis, page.TotalItems, page.State = result.Items, result.TotalItems, result.State
		}
		page.IsFetching = facet.IsFetching()
		page.ExpectedTotal = facet.ExpectedCount()
	case AbisKnown:
		facet := c.knownFacet
		var filterFunc func(*Abi) bool
		if filter != "" {
			filterFunc = func(item *Abi) bool {
				return c.matchesKnownFilter(item, filter)
			}
		}
		sortFunc := func(items []Abi, sort sdk.SortSpec) error {
			return sdk.SortAbis(items, sort)
		}
		if result, err := facet.GetPage(first, pageSize, filterFunc, sortSpec, sortFunc); err != nil {
			return nil, types.NewStoreError("abis", dataFacet, "GetPage", err)
		} else {

			page.Abis, page.TotalItems, page.State = result.Items, result.TotalItems, result.State
		}
		page.IsFetching = facet.IsFetching()
		page.ExpectedTotal = facet.ExpectedCount()
	case AbisFunctions:
		facet := c.functionsFacet
		var filterFunc func(*Function) bool
		if filter != "" {
			filterFunc = func(item *Function) bool {
				return c.matchesFunctionFilter(item, filter)
			}
		}
		sortFunc := func(items []Function, sort sdk.SortSpec) error {
			return sdk.SortFunctions(items, sort)
		}
		if result, err := facet.GetPage(first, pageSize, filterFunc, sortSpec, sortFunc); err != nil {
			return nil, types.NewStoreError("abis", dataFacet, "GetPage", err)
		} else {

			page.Functions, page.TotalItems, page.State = result.Items, result.TotalItems, result.State
		}
		page.IsFetching = facet.IsFetching()
		page.ExpectedTotal = facet.ExpectedCount()
	case AbisEvents:
		facet := c.eventsFacet
		var filterFunc func(*Function) bool
		if filter != "" {
			filterFunc = func(item *Function) bool {
				return c.matchesEventFilter(item, filter)
			}
		}
		sortFunc := func(items []Function, sort sdk.SortSpec) error {
			return sdk.SortFunctions(items, sort)
		}
		if result, err := facet.GetPage(first, pageSize, filterFunc, sortSpec, sortFunc); err != nil {
			return nil, types.NewStoreError("abis", dataFacet, "GetPage", err)
		} else {

			page.Functions, page.TotalItems, page.State = result.Items, result.TotalItems, result.State
		}
		page.IsFetching = facet.IsFetching()
		page.ExpectedTotal = facet.ExpectedCount()
	default:
		return nil, types.NewValidationError("abis", dataFacet, "GetPage",
			fmt.Errorf("unsupported dataFacet: %v", dataFacet))
	}

	return page, nil
}

// EXISTING_CODE
func (c *AbisCollection) matchesDownloadedFilter(item *Abi, filter string) bool {
	return strings.Contains(strings.ToLower(item.Name), filter)
}

func (c *AbisCollection) matchesKnownFilter(item *Abi, filter string) bool {
	return strings.Contains(strings.ToLower(item.Name), filter)
}

func (c *AbisCollection) matchesFunctionFilter(fn *Function, filter string) bool {
	if filter == "" {
		return true
	}
	filterLower := strings.ToLower(filter)
	return strings.Contains(strings.ToLower(fn.Name), filterLower) ||
		strings.Contains(strings.ToLower(fn.Encoding), filterLower)
}

func (c *AbisCollection) matchesEventFilter(item *Function, filter string) bool {
	return strings.Contains(strings.ToLower(item.Name), filter) ||
		strings.Contains(strings.ToLower(item.Encoding), filter)
}

// EXISTING_CODE
