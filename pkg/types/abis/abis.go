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
	"sync"
	"time"

	// EXISTING_CODE
	// EXISTING_CODE
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/facets"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/logging"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

const (
	AbisDownloaded types.DataFacet = "downloaded"
	AbisKnown      types.DataFacet = "known"
	AbisFunctions  types.DataFacet = "functions"
	AbisEvents     types.DataFacet = "events"
)

func init() {
	types.RegisterDataFacet(AbisDownloaded)
	types.RegisterDataFacet(AbisKnown)
	types.RegisterDataFacet(AbisFunctions)
	types.RegisterDataFacet(AbisEvents)
}

type AbisCollection struct {
	downloadedFacet *facets.Facet[Abi]
	knownFacet      *facets.Facet[Abi]
	functionsFacet  *facets.Facet[Function]
	eventsFacet     *facets.Facet[Function]
	summary         types.Summary
	summaryMutex    sync.RWMutex
}

func NewAbisCollection() *AbisCollection {
	c := &AbisCollection{}
	c.ResetSummary()
	c.initializeFacets()
	return c
}

func (c *AbisCollection) initializeFacets() {
	c.downloadedFacet = facets.NewFacetWithSummary(
		AbisDownloaded,
		isDownloaded,
		isDupAbi(),
		c.getAbisStore(),
		"abis",
		c,
	)

	c.knownFacet = facets.NewFacetWithSummary(
		AbisKnown,
		isKnown,
		isDupAbi(),
		c.getAbisStore(),
		"abis",
		c,
	)

	c.functionsFacet = facets.NewFacetWithSummary(
		AbisFunctions,
		isFunction,
		isDupFunction(),
		c.getFunctionsStore(),
		"abis",
		c,
	)

	c.eventsFacet = facets.NewFacetWithSummary(
		AbisEvents,
		isEvent,
		isDupFunction(),
		c.getFunctionsStore(),
		"abis",
		c,
	)
}

func isDownloaded(item *Abi) bool {
	// EXISTING_CODE
	return !item.IsKnown
	// EXISTING_CODE
}

func isKnown(item *Abi) bool {
	// EXISTING_CODE
	return item.IsKnown
	// EXISTING_CODE
}

func isFunction(item *Function) bool {
	// EXISTING_CODE
	return item.FunctionType != "event"
	// EXISTING_CODE
}

func isEvent(item *Function) bool {
	// EXISTING_CODE
	return item.FunctionType == "event"
	// EXISTING_CODE
}

func isDupAbi() func(existing []*Abi, newItem *Abi) bool {
	// EXISTING_CODE
	return nil
	// EXISTING_CODE
}

func isDupFunction() func(existing []*Function, newItem *Function) bool {
	// EXISTING_CODE
	seen := make(map[string]bool)
	lastExistingLen := 0

	return func(existing []*Function, newItem *Function) bool {
		if newItem == nil {
			return false
		}

		if len(existing) == 0 && lastExistingLen > 0 {
			seen = make(map[string]bool)
		}
		lastExistingLen = len(existing)

		if seen[newItem.Encoding] {
			return true
		}
		seen[newItem.Encoding] = true
		return false
	}
	// EXISTING_CODE
}

func (c *AbisCollection) LoadData(dataFacet types.DataFacet) {
	if !c.NeedsUpdate(dataFacet) {
		return
	}

	go func() {
		switch dataFacet {
		case AbisDownloaded:
			if err := c.downloadedFacet.Load(); err != nil {
				logging.LogError(fmt.Sprintf("LoadData.%s from store: %%v", dataFacet), err, facets.ErrAlreadyLoading)
			}
		case AbisKnown:
			if err := c.knownFacet.Load(); err != nil {
				logging.LogError(fmt.Sprintf("LoadData.%s from store: %%v", dataFacet), err, facets.ErrAlreadyLoading)
			}
		case AbisFunctions:
			if err := c.functionsFacet.Load(); err != nil {
				logging.LogError(fmt.Sprintf("LoadData.%s from store: %%v", dataFacet), err, facets.ErrAlreadyLoading)
			}
		case AbisEvents:
			if err := c.eventsFacet.Load(); err != nil {
				logging.LogError(fmt.Sprintf("LoadData.%s from store: %%v", dataFacet), err, facets.ErrAlreadyLoading)
			}
		default:
			logging.LogError("LoadData: unexpected dataFacet: %v", fmt.Errorf("invalid dataFacet: %s", dataFacet), nil)
			return
		}
	}()
}

func (c *AbisCollection) Reset(dataFacet types.DataFacet) {
	switch dataFacet {
	case AbisDownloaded:
		c.downloadedFacet.GetStore().Reset()
	case AbisKnown:
		c.knownFacet.GetStore().Reset()
	case AbisFunctions:
		c.functionsFacet.GetStore().Reset()
	case AbisEvents:
		c.eventsFacet.GetStore().Reset()
	default:
		return
	}
}

func (c *AbisCollection) NeedsUpdate(dataFacet types.DataFacet) bool {
	switch dataFacet {
	case AbisDownloaded:
		return c.downloadedFacet.NeedsUpdate()
	case AbisKnown:
		return c.knownFacet.NeedsUpdate()
	case AbisFunctions:
		return c.functionsFacet.NeedsUpdate()
	case AbisEvents:
		return c.eventsFacet.NeedsUpdate()
	default:
		return false
	}
}

func (c *AbisCollection) GetSupportedFacets() []types.DataFacet {
	return []types.DataFacet{
		AbisDownloaded,
		AbisKnown,
		AbisFunctions,
		AbisEvents,
	}
}

func (c *AbisCollection) AccumulateItem(item interface{}, summary *types.Summary) {
	// EXISTING_CODE
	c.summaryMutex.Lock()
	defer c.summaryMutex.Unlock()

	if summary.FacetCounts == nil {
		summary.FacetCounts = make(map[types.DataFacet]int)
	}

	switch v := item.(type) {
	case *Abi:
		summary.TotalCount++

		if v.IsKnown {
			summary.FacetCounts[AbisKnown]++
		} else {
			summary.FacetCounts[AbisDownloaded]++
		}

		if summary.CustomData == nil {
			summary.CustomData = make(map[string]interface{})
		}

		knownCount, _ := summary.CustomData["knownCount"].(int)
		downloadedCount, _ := summary.CustomData["downloadedCount"].(int)

		if v.IsKnown {
			knownCount++
		} else {
			downloadedCount++
		}

		summary.CustomData["knownCount"] = knownCount
		summary.CustomData["downloadedCount"] = downloadedCount

	case *Function:
		summary.TotalCount++

		if v.FunctionType == "event" {
			summary.FacetCounts[AbisEvents]++
		} else {
			summary.FacetCounts[AbisFunctions]++
		}

		if summary.CustomData == nil {
			summary.CustomData = make(map[string]interface{})
		}

		functionsCount, _ := summary.CustomData["functionsCount"].(int)
		eventsCount, _ := summary.CustomData["eventsCount"].(int)

		if v.FunctionType == "event" {
			eventsCount++
		} else {
			functionsCount++
		}

		summary.CustomData["functionsCount"] = functionsCount
		summary.CustomData["eventsCount"] = eventsCount
	}
	// EXISTING_CODE
}

func (c *AbisCollection) GetSummary() types.Summary {
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

func (c *AbisCollection) ResetSummary() {
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
