// Copyright 2016, 2026 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package chunks

import (
	"fmt"
	"strings"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
	// EXISTING_CODE
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

func (c *ChunksCollection) GetPage(
	payload *types.Payload,
	first, pageSize int,
	sortSpec sdk.SortSpec,
	filter string,
) (types.Page, error) {
	dataFacet := payload.DataFacet

	page := &ChunksPage{
		Facet: dataFacet,
	}
	filter = strings.ToLower(filter)

	switch dataFacet {
	case ChunksStats:
		facet := c.statsFacet
		var filterFunc func(*Stats) bool
		if filter != "" {
			filterFunc = func(item *Stats) bool {
				return c.matchesStatsFilter(item, filter)
			}
		}
		sortFunc := func(items []Stats, sort sdk.SortSpec) error {
			return sdk.SortStats(items, sort)
		}
		if result, err := facet.GetPage(first, pageSize, filterFunc, sortSpec, sortFunc); err != nil {
			return nil, types.NewStoreError("chunks", dataFacet, "GetPage", err)
		} else {
			stats := make([]*Stats, 0, len(result.Items))
			for i := range result.Items {
				stats = append(stats, &result.Items[i])
			}
			page.Stats, page.TotalItems, page.State = stats, result.TotalItems, result.State
		}
		page.IsFetching = facet.IsFetching()
		page.ExpectedTotal = facet.ExpectedCount()
	case ChunksIndex:
		facet := c.indexFacet
		var filterFunc func(*Index) bool
		if filter != "" {
			filterFunc = func(item *Index) bool {
				return c.matchesIndexFilter(item, filter)
			}
		}
		sortFunc := func(items []Index, sort sdk.SortSpec) error {
			return sdk.SortIndex(items, sort)
		}
		if result, err := facet.GetPage(first, pageSize, filterFunc, sortSpec, sortFunc); err != nil {
			return nil, types.NewStoreError("chunks", dataFacet, "GetPage", err)
		} else {
			index := make([]*Index, 0, len(result.Items))
			for i := range result.Items {
				index = append(index, &result.Items[i])
			}
			page.Index, page.TotalItems, page.State = index, result.TotalItems, result.State
		}
		page.IsFetching = facet.IsFetching()
		page.ExpectedTotal = facet.ExpectedCount()
	case ChunksBlooms:
		facet := c.bloomsFacet
		var filterFunc func(*Bloom) bool
		if filter != "" {
			filterFunc = func(item *Bloom) bool {
				return c.matchesBloomFilter(item, filter)
			}
		}
		sortFunc := func(items []Bloom, sort sdk.SortSpec) error {
			return sdk.SortBlooms(items, sort)
		}
		if result, err := facet.GetPage(first, pageSize, filterFunc, sortSpec, sortFunc); err != nil {
			return nil, types.NewStoreError("chunks", dataFacet, "GetPage", err)
		} else {
			bloom := make([]*Bloom, 0, len(result.Items))
			for i := range result.Items {
				bloom = append(bloom, &result.Items[i])
			}
			page.Blooms, page.TotalItems, page.State = bloom, result.TotalItems, result.State
		}
		page.IsFetching = facet.IsFetching()
		page.ExpectedTotal = facet.ExpectedCount()
	case ChunksManifest:
		facet := c.manifestFacet
		var filterFunc func(*Manifest) bool
		if filter != "" {
			filterFunc = func(item *Manifest) bool {
				return c.matchesManifestFilter(item, filter)
			}
		}
		sortFunc := func(items []Manifest, sort sdk.SortSpec) error {
			return sdk.SortManifest(items, sort)
		}
		if result, err := facet.GetPage(first, pageSize, filterFunc, sortSpec, sortFunc); err != nil {
			return nil, types.NewStoreError("chunks", dataFacet, "GetPage", err)
		} else {
			manifest := make([]*Manifest, 0, len(result.Items))
			for i := range result.Items {
				manifest = append(manifest, &result.Items[i])
			}
			page.Manifest, page.TotalItems, page.State = manifest, result.TotalItems, result.State
		}
		page.IsFetching = facet.IsFetching()
		page.ExpectedTotal = facet.ExpectedCount()
	default:
		return nil, types.NewValidationError("chunks", dataFacet, "GetPage",
			fmt.Errorf("unsupported dataFacet: %v", dataFacet))
	}

	return page, nil
}

// EXISTING_CODE
// EXISTING_CODE
