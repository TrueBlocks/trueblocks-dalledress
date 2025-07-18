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
	"sync"
	"time"

	// EXISTING_CODE
	// EXISTING_CODE
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/facets"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/logging"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

const (
	ChunksStats    types.DataFacet = "stats"
	ChunksIndex    types.DataFacet = "index"
	ChunksBlooms   types.DataFacet = "blooms"
	ChunksManifest types.DataFacet = "manifest"
)

func init() {
	types.RegisterDataFacet(ChunksStats)
	types.RegisterDataFacet(ChunksIndex)
	types.RegisterDataFacet(ChunksBlooms)
	types.RegisterDataFacet(ChunksManifest)
}

type ChunksCollection struct {
	statsFacet    *facets.Facet[Stats]
	indexFacet    *facets.Facet[Index]
	bloomsFacet   *facets.Facet[Bloom]
	manifestFacet *facets.Facet[Manifest]
	summary       types.Summary
	summaryMutex  sync.RWMutex
}

func NewChunksCollection() *ChunksCollection {
	c := &ChunksCollection{}
	c.ResetSummary()
	c.initializeFacets()
	return c
}

func (c *ChunksCollection) initializeFacets() {
	c.statsFacet = facets.NewFacet(
		ChunksStats,
		isStats,
		isDupStats(),
		c.getStatsStore(ChunksStats),
		"chunks",
		c,
	)

	c.indexFacet = facets.NewFacet(
		ChunksIndex,
		isIndex,
		isDupIndex(),
		c.getIndexStore(ChunksIndex),
		"chunks",
		c,
	)

	c.bloomsFacet = facets.NewFacet(
		ChunksBlooms,
		isBloom,
		isDupBloom(),
		c.getBloomsStore(ChunksBlooms),
		"chunks",
		c,
	)

	c.manifestFacet = facets.NewFacet(
		ChunksManifest,
		isManifest,
		isDupManifest(),
		c.getManifestStore(ChunksManifest),
		"chunks",
		c,
	)
}

func isStats(item *Stats) bool {
	// EXISTING_CODE
	return true
	// EXISTING_CODE
}

func isIndex(item *Index) bool {
	// EXISTING_CODE
	return true
	// EXISTING_CODE
}

func isBloom(item *Bloom) bool {
	// EXISTING_CODE
	return true
	// EXISTING_CODE
}

func isManifest(item *Manifest) bool {
	// EXISTING_CODE
	return true
	// EXISTING_CODE
}

func isDupBloom() func(existing []*Bloom, newItem *Bloom) bool {
	// EXISTING_CODE
	return nil
	// EXISTING_CODE
}

func isDupIndex() func(existing []*Index, newItem *Index) bool {
	// EXISTING_CODE
	return nil
	// EXISTING_CODE
}

func isDupManifest() func(existing []*Manifest, newItem *Manifest) bool {
	// EXISTING_CODE
	return nil
	// EXISTING_CODE
}

func isDupStats() func(existing []*Stats, newItem *Stats) bool {
	// EXISTING_CODE
	return nil
	// EXISTING_CODE
}

func (c *ChunksCollection) LoadData(dataFacet types.DataFacet) {
	if !c.NeedsUpdate(dataFacet) {
		return
	}

	go func() {
		switch dataFacet {
		case ChunksStats:
			if err := c.statsFacet.Load(); err != nil {
				logging.LogError(fmt.Sprintf("LoadData.%s from store: %%v", dataFacet), err, facets.ErrAlreadyLoading)
			}
		case ChunksIndex:
			if err := c.indexFacet.Load(); err != nil {
				logging.LogError(fmt.Sprintf("LoadData.%s from store: %%v", dataFacet), err, facets.ErrAlreadyLoading)
			}
		case ChunksBlooms:
			if err := c.bloomsFacet.Load(); err != nil {
				logging.LogError(fmt.Sprintf("LoadData.%s from store: %%v", dataFacet), err, facets.ErrAlreadyLoading)
			}
		case ChunksManifest:
			if err := c.manifestFacet.Load(); err != nil {
				logging.LogError(fmt.Sprintf("LoadData.%s from store: %%v", dataFacet), err, facets.ErrAlreadyLoading)
			}
		default:
			logging.LogError("LoadData: unexpected dataFacet: %v", fmt.Errorf("invalid dataFacet: %s", dataFacet), nil)
			return
		}
	}()
}

func (c *ChunksCollection) Reset(dataFacet types.DataFacet) {
	switch dataFacet {
	case ChunksStats:
		c.statsFacet.GetStore().Reset()
	case ChunksIndex:
		c.indexFacet.GetStore().Reset()
	case ChunksBlooms:
		c.bloomsFacet.GetStore().Reset()
	case ChunksManifest:
		c.manifestFacet.GetStore().Reset()
	default:
		return
	}
}

func (c *ChunksCollection) NeedsUpdate(dataFacet types.DataFacet) bool {
	switch dataFacet {
	case ChunksStats:
		return c.statsFacet.NeedsUpdate()
	case ChunksIndex:
		return c.indexFacet.NeedsUpdate()
	case ChunksBlooms:
		return c.bloomsFacet.NeedsUpdate()
	case ChunksManifest:
		return c.manifestFacet.NeedsUpdate()
	default:
		return false
	}
}

func (c *ChunksCollection) GetSupportedFacets() []types.DataFacet {
	return []types.DataFacet{
		ChunksStats,
		ChunksIndex,
		ChunksBlooms,
		ChunksManifest,
	}
}

func (c *ChunksCollection) AccumulateItem(item interface{}, summary *types.Summary) {
	// EXISTING_CODE
	c.summaryMutex.Lock()
	defer c.summaryMutex.Unlock()

	if summary.FacetCounts == nil {
		summary.FacetCounts = make(map[types.DataFacet]int)
	}

	switch item.(type) {
	case *Stats:
		summary.TotalCount++
		summary.FacetCounts[ChunksStats]++
		if summary.CustomData == nil {
			summary.CustomData = make(map[string]interface{})
		}

		statsCount, _ := summary.CustomData["statsCount"].(int)
		totalBytes, _ := summary.CustomData["totalBytes"].(int64)

		statsCount++
		summary.CustomData["statsCount"] = statsCount
		summary.CustomData["totalBytes"] = totalBytes

	case *Index:
		summary.TotalCount++
		summary.FacetCounts[ChunksIndex]++
		if summary.CustomData == nil {
			summary.CustomData = make(map[string]interface{})
		}

		indexCount, _ := summary.CustomData["indexCount"].(int)
		indexCount++
		summary.CustomData["indexCount"] = indexCount

	case *Bloom:
		summary.TotalCount++
		summary.FacetCounts[ChunksBlooms]++
		if summary.CustomData == nil {
			summary.CustomData = make(map[string]interface{})
		}

		bloomsCount, _ := summary.CustomData["bloomsCount"].(int)
		bloomsCount++
		summary.CustomData["bloomsCount"] = bloomsCount

	case *Manifest:
		summary.TotalCount++
		summary.FacetCounts[ChunksManifest]++
		if summary.CustomData == nil {
			summary.CustomData = make(map[string]interface{})
		}

		manifestCount, _ := summary.CustomData["manifestCount"].(int)
		manifestCount++
		summary.CustomData["manifestCount"] = manifestCount
	}
	// EXISTING_CODE
}

func (c *ChunksCollection) GetSummary() types.Summary {
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

func (c *ChunksCollection) ResetSummary() {
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
