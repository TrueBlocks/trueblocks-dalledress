// Copyright 2016, 2026 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package status

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
	StatusStatus types.DataFacet = "status"
	StatusCaches types.DataFacet = "caches"
	StatusChains types.DataFacet = "chains"
)

func init() {
	types.RegisterDataFacet(StatusStatus)
	types.RegisterDataFacet(StatusCaches)
	types.RegisterDataFacet(StatusChains)
}

type StatusCollection struct {
	statusFacet  *facets.Facet[Status]
	cachesFacet  *facets.Facet[Cache]
	chainsFacet  *facets.Facet[Chain]
	summary      types.Summary
	summaryMutex sync.RWMutex
}

func NewStatusCollection() *StatusCollection {
	c := &StatusCollection{}
	c.ResetSummary()
	c.initializeFacets()
	return c
}

func (c *StatusCollection) initializeFacets() {
	c.statusFacet = facets.NewFacet(
		StatusStatus,
		isStatus,
		isDupStatus(),
		c.getStatusStore(StatusStatus),
		"status",
		c,
	)

	c.cachesFacet = facets.NewFacet(
		StatusCaches,
		isCache,
		isDupCache(),
		c.getCachesStore(StatusCaches),
		"status",
		c,
	)

	c.chainsFacet = facets.NewFacet(
		StatusChains,
		isChain,
		isDupChain(),
		c.getChainsStore(StatusChains),
		"status",
		c,
	)
}

func isStatus(item *Status) bool {
	// EXISTING_CODE
	return true
	// EXISTING_CODE
}

func isCache(item *Cache) bool {
	// EXISTING_CODE
	return true
	// EXISTING_CODE
}

func isChain(item *Chain) bool {
	// EXISTING_CODE
	return true
	// EXISTING_CODE
}

func isDupCache() func(existing []*Cache, newItem *Cache) bool {
	// EXISTING_CODE
	return func(existing []*Cache, newItem *Cache) bool {
		return false
	}
	// EXISTING_CODE
}

func isDupChain() func(existing []*Chain, newItem *Chain) bool {
	// EXISTING_CODE
	return func(existing []*Chain, newItem *Chain) bool {
		return false
	}
	// EXISTING_CODE
}

func isDupStatus() func(existing []*Status, newItem *Status) bool {
	// EXISTING_CODE
	return func(existing []*Status, newItem *Status) bool {
		return false
	}
	// EXISTING_CODE
}

func (c *StatusCollection) LoadData(dataFacet types.DataFacet) {
	if !c.NeedsUpdate(dataFacet) {
		return
	}

	go func() {
		switch dataFacet {
		case StatusStatus:
			if err := c.statusFacet.Load(); err != nil {
				logging.LogError(fmt.Sprintf("LoadData.%s from store: %%v", dataFacet), err, facets.ErrAlreadyLoading)
			}
		case StatusCaches:
			if err := c.cachesFacet.Load(); err != nil {
				logging.LogError(fmt.Sprintf("LoadData.%s from store: %%v", dataFacet), err, facets.ErrAlreadyLoading)
			}
		case StatusChains:
			if err := c.chainsFacet.Load(); err != nil {
				logging.LogError(fmt.Sprintf("LoadData.%s from store: %%v", dataFacet), err, facets.ErrAlreadyLoading)
			}
		default:
			logging.LogError("LoadData: unexpected dataFacet: %v", fmt.Errorf("invalid dataFacet: %s", dataFacet), nil)
			return
		}
	}()
}

func (c *StatusCollection) Reset(dataFacet types.DataFacet) {
	switch dataFacet {
	case StatusStatus:
		c.statusFacet.GetStore().Reset()
	case StatusCaches:
		c.cachesFacet.GetStore().Reset()
	case StatusChains:
		c.chainsFacet.GetStore().Reset()
	default:
		return
	}
}

func (c *StatusCollection) NeedsUpdate(dataFacet types.DataFacet) bool {
	switch dataFacet {
	case StatusStatus:
		return c.statusFacet.NeedsUpdate()
	case StatusCaches:
		return c.cachesFacet.NeedsUpdate()
	case StatusChains:
		return c.chainsFacet.NeedsUpdate()
	default:
		return false
	}
}

func (c *StatusCollection) GetSupportedFacets() []types.DataFacet {
	return []types.DataFacet{
		StatusStatus,
		StatusCaches,
		StatusChains,
	}
}

func (c *StatusCollection) AccumulateItem(item interface{}, summary *types.Summary) {
	// EXISTING_CODE
	c.summaryMutex.Lock()
	defer c.summaryMutex.Unlock()

	if summary.FacetCounts == nil {
		summary.FacetCounts = make(map[types.DataFacet]int)
	}

	switch item.(type) {
	case *Cache:
		summary.TotalCount++
		summary.FacetCounts[StatusCaches]++
		if summary.CustomData == nil {
			summary.CustomData = make(map[string]interface{})
		}

		cachesCount, _ := summary.CustomData["caches"].(int)
		cachesCount++
		summary.CustomData["cachesCount"] = cachesCount

	case *Chain:
		summary.TotalCount++
		summary.FacetCounts[StatusChains]++
		if summary.CustomData == nil {
			summary.CustomData = make(map[string]interface{})
		}

		chainsCount, _ := summary.CustomData["chains"].(int)
		chainsCount++

		summary.CustomData["chainsCount"] = chainsCount
	}
	// EXISTING_CODE
}

func (c *StatusCollection) GetSummary() types.Summary {
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

func (c *StatusCollection) ResetSummary() {
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
