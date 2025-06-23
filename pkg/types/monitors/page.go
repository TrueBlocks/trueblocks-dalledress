// Copyright 2016, 2025 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package monitors

import (
	// EXISTING_CODE
	"fmt"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
	// EXISTING_CODE
)

// TODO: The slices should be slices to pointers
type MonitorsPage struct {
	Facet         types.DataFacet `json:"facet"`
	Monitors      []Monitor       `json:"monitors"`
	TotalItems    int             `json:"totalItems"`
	ExpectedTotal int             `json:"expectedTotal"`
	IsFetching    bool            `json:"isFetching"`
	State         types.LoadState `json:"state"`
}

func (p *MonitorsPage) GetFacet() types.DataFacet {
	return p.Facet
}

func (p *MonitorsPage) GetTotalItems() int {
	return p.TotalItems
}

func (p *MonitorsPage) GetExpectedTotal() int {
	return p.ExpectedTotal
}

func (p *MonitorsPage) GetIsFetching() bool {
	return p.IsFetching
}

func (p *MonitorsPage) GetState() types.LoadState {
	return p.State
}

// EXISTING_CODE
func (c *MonitorsCollection) GetPage(
	payload types.Payload,
	first, pageSize int,
	sortSpec sdk.SortSpec,
	filter string,
) (types.Page, error) {
	dataFacet := payload.DataFacet

	switch dataFacet {
	case MonitorsMonitors:
		var filterFunc func(*Monitor) bool
		if filter != "" {
			filterFunc = func(monitor *Monitor) bool {
				return c.matchesFilter(monitor, filter)
			}
		}

		var sortFunc func([]Monitor, sdk.SortSpec) error
		sortFunc = func(items []Monitor, sort sdk.SortSpec) error {
			return sdk.SortMonitors(items, sort)
		}

		pageResult, err := c.monitorsFacet.GetPage(
			first,
			pageSize,
			filterFunc,
			sortSpec,
			sortFunc,
		)
		if err != nil {
			// This is likely an SDK or store error, not a validation error
			return nil, types.NewStoreError("monitors", dataFacet, "GetPage", err)
		}

		return &MonitorsPage{
			Facet:         dataFacet,
			Monitors:      pageResult.Items,
			TotalItems:    pageResult.TotalItems,
			ExpectedTotal: c.getExpectedTotal(dataFacet),
			IsFetching:    c.monitorsFacet.IsFetching(),
			State:         pageResult.State,
		}, nil
	default:
		// This is truly a validation error - invalid DataFacet for this collection
		return nil, types.NewValidationError("monitors", dataFacet, "GetPage",
			fmt.Errorf("unsupported dataFacet: %s", dataFacet))
	}
}

// EXISTING_CODE
