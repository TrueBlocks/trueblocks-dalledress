// Copyright 2016, 2025 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package chunks

import (
	// EXISTING_CODE
	"fmt"
	"strings"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
	// EXISTING_CODE
)

type ChunksPage struct {
	Facet         types.DataFacet `json:"facet"`
	Blooms        []*Bloom        `json:"blooms"`
	Index         []*Index        `json:"index"`
	Manifest      []*Manifest     `json:"manifest"`
	Stats         []*Stats        `json:"stats"`
	TotalItems    int             `json:"totalItems"`
	ExpectedTotal int             `json:"expectedTotal"`
	IsFetching    bool            `json:"isFetching"`
	State         types.LoadState `json:"state"`
}

func (p *ChunksPage) GetFacet() types.DataFacet {
	return p.Facet
}

func (p *ChunksPage) GetTotalItems() int {
	return p.TotalItems
}

func (p *ChunksPage) GetExpectedTotal() int {
	return p.ExpectedTotal
}

func (p *ChunksPage) GetIsFetching() bool {
	return p.IsFetching
}

func (p *ChunksPage) GetState() types.LoadState {
	return p.State
}

// EXISTING_CODE
func (c *ChunksCollection) GetPage(
	dataFacet types.DataFacet,
	first, pageSize int,
	sortSpec sdk.SortSpec,
	filter string,
) (types.Page, error) {
	page := &ChunksPage{
		Facet: dataFacet,
	}
	filter = strings.ToLower(filter)

	switch dataFacet {
	case ChunksStats:
		var filterFunc func(*Stats) bool
		if filter != "" {
			filterFunc = func(stat *Stats) bool {
				return c.matchesStatsFilter(stat, filter)
			}
		}

		sortFunc := func(items []Stats, sort sdk.SortSpec) error {
			// Placeholder - implement actual sorting if needed
			return nil
		}

		if result, err := c.statsFacet.GetPage(first, pageSize, filterFunc, sortSpec, sortFunc); err != nil {
			// This is likely an SDK or store error, not a validation error
			return nil, types.NewStoreError("chunks", dataFacet, "GetPage", err)
		} else {
			stats := make([]*Stats, 0, len(result.Items))
			for i := range result.Items {
				stats = append(stats, &result.Items[i])
			}
			page.Stats, page.TotalItems, page.State = stats, result.TotalItems, result.State
		}
		page.IsFetching = c.statsFacet.IsFetching()
		page.ExpectedTotal = c.statsFacet.ExpectedCount()

	case ChunksIndex:
		var filterFunc func(*Index) bool
		if filter != "" {
			filterFunc = func(index *Index) bool {
				return c.matchesIndexFilter(index, filter)
			}
		}

		sortFunc := func(items []Index, sort sdk.SortSpec) error {
			// Placeholder - implement actual sorting if needed
			return nil
		}

		if result, err := c.indexFacet.GetPage(first, pageSize, filterFunc, sortSpec, sortFunc); err != nil {
			// This is likely an SDK or store error, not a validation error
			return nil, types.NewStoreError("chunks", dataFacet, "GetPage", err)
		} else {
			index := make([]*Index, 0, len(result.Items))
			for i := range result.Items {
				index = append(index, &result.Items[i])
			}
			page.Index, page.TotalItems, page.State = index, result.TotalItems, result.State
		}
		page.IsFetching = c.indexFacet.IsFetching()
		page.ExpectedTotal = c.indexFacet.ExpectedCount()

	case ChunksBlooms:
		var filterFunc func(*Bloom) bool
		if filter != "" {
			filterFunc = func(bloom *Bloom) bool {
				return c.matchesBloomsFilter(bloom, filter)
			}
		}

		sortFunc := func(items []Bloom, sort sdk.SortSpec) error {
			// Placeholder - implement actual sorting if needed
			return nil
		}

		if result, err := c.bloomsFacet.GetPage(first, pageSize, filterFunc, sortSpec, sortFunc); err != nil {
			// This is likely an SDK or store error, not a validation error
			return nil, types.NewStoreError("chunks", dataFacet, "GetPage", err)
		} else {
			blooms := make([]*Bloom, 0, len(result.Items))
			for i := range result.Items {
				blooms = append(blooms, &result.Items[i])
			}
			page.Blooms, page.TotalItems, page.State = blooms, result.TotalItems, result.State
		}
		page.IsFetching = c.bloomsFacet.IsFetching()
		page.ExpectedTotal = c.bloomsFacet.ExpectedCount()

	case ChunksManifest:
		var filterFunc func(*Manifest) bool
		if filter != "" {
			filterFunc = func(manifest *Manifest) bool {
				return c.matchesManifestFilter(manifest, filter)
			}
		}

		sortFunc := func(items []Manifest, sort sdk.SortSpec) error {
			// Placeholder - implement actual sorting if needed
			return nil
		}

		if result, err := c.manifestFacet.GetPage(first, pageSize, filterFunc, sortSpec, sortFunc); err != nil {
			// This is likely an SDK or store error, not a validation error
			return nil, types.NewStoreError("chunks", dataFacet, "GetPage", err)
		} else {
			manifest := make([]*Manifest, 0, len(result.Items))
			for i := range result.Items {
				manifest = append(manifest, &result.Items[i])
			}
			page.Manifest, page.TotalItems, page.State = manifest, result.TotalItems, result.State
		}
		page.IsFetching = c.manifestFacet.IsFetching()
		page.ExpectedTotal = c.manifestFacet.ExpectedCount()

	default:
		// This is truly a validation error - invalid DataFacet for this collection
		return nil, types.NewValidationError("chunks", dataFacet, "GetPage",
			fmt.Errorf("unsupported data facet: %v", dataFacet))
	}

	return page, nil
}

// EXISTING_CODE
