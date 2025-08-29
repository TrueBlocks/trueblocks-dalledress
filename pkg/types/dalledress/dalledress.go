// Copyright 2016, 2026 The Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package dalledress

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"sync"
	"time"

	// EXISTING_CODE
	dallev2 "github.com/TrueBlocks/trueblocks-dalle/v2"
	// EXISTING_CODE
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/facets"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/logging"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

const (
	DalleDressGenerator types.DataFacet = "generator"
	DalleDressSeries    types.DataFacet = "series"
	DalleDressDatabases types.DataFacet = "databases"
	DalleDressEvents    types.DataFacet = "events"
	DalleDressGallery   types.DataFacet = "gallery"
)

func init() {
	types.RegisterDataFacet(DalleDressGenerator)
	types.RegisterDataFacet(DalleDressSeries)
	types.RegisterDataFacet(DalleDressDatabases)
	types.RegisterDataFacet(DalleDressEvents)
	types.RegisterDataFacet(DalleDressGallery)
}

type DalleDressCollection struct {
	generatorFacet *facets.Facet[DalleDress]
	seriesFacet    *facets.Facet[Series]
	databasesFacet *facets.Facet[Database]
	eventsFacet    *facets.Facet[Log]
	galleryFacet   *facets.Facet[DalleDress]
	summary        types.Summary
	summaryMutex   sync.RWMutex
	//
	galleryCache      []*GalleryItem
	gallerySeriesInfo map[string]int64
	galleryCacheMux   sync.RWMutex
}

func NewDalleDressCollection(payload *types.Payload) *DalleDressCollection {
	c := &DalleDressCollection{}
	c.ResetSummary()
	c.gallerySeriesInfo = make(map[string]int64)
	c.initializeFacets(payload)
	return c
}

func (c *DalleDressCollection) initializeFacets(payload *types.Payload) {
	c.generatorFacet = facets.NewFacet(
		DalleDressGenerator,
		isGenerator,
		isDupDalleDress(),
		c.getDalleDressStore(payload, DalleDressGenerator),
		"dalledress",
		c,
	)

	c.seriesFacet = facets.NewFacet(
		DalleDressSeries,
		isSeries,
		isDupSeries(),
		c.getSeriesStore(payload, DalleDressSeries),
		"dalledress",
		c,
	)

	c.databasesFacet = facets.NewFacet(
		DalleDressDatabases,
		isDatabase,
		isDupDatabase(),
		c.getDatabasesStore(payload, DalleDressDatabases),
		"dalledress",
		c,
	)

	c.eventsFacet = facets.NewFacet(
		DalleDressEvents,
		isEvent,
		isDupLog(),
		c.getLogsStore(payload, DalleDressEvents),
		"dalledress",
		c,
	)

	c.galleryFacet = facets.NewFacet(
		DalleDressGallery,
		isGallery,
		isDupDalleDress(),
		c.getDalleDressStore(payload, DalleDressGallery),
		"dalledress",
		c,
	)
}

func isGenerator(item *DalleDress) bool {
	// EXISTING_CODE
	return true
	// EXISTING_CODE
}

func isSeries(item *Series) bool {
	// EXISTING_CODE
	return true
	// EXISTING_CODE
}

func isDatabase(item *Database) bool {
	// EXISTING_CODE
	return true
	// EXISTING_CODE
}

func isEvent(item *Log) bool {
	// EXISTING_CODE
	return true
	// EXISTING_CODE
}

func isGallery(item *DalleDress) bool {
	// EXISTING_CODE
	return true
	// EXISTING_CODE
}

func isDupLog() func(existing []*Log, newItem *Log) bool {
	// EXISTING_CODE
	return nil
	// EXISTING_CODE
}

func isDupDalleDress() func(existing []*DalleDress, newItem *DalleDress) bool {
	// EXISTING_CODE
	return nil
	// EXISTING_CODE
}

func isDupDatabase() func(existing []*Database, newItem *Database) bool {
	// EXISTING_CODE
	return nil
	// EXISTING_CODE
}

func isDupSeries() func(existing []*Series, newItem *Series) bool {
	// EXISTING_CODE
	return nil
	// EXISTING_CODE
}

func (c *DalleDressCollection) LoadData(dataFacet types.DataFacet) {
	if !c.NeedsUpdate(dataFacet) {
		return
	}

	go func() {
		switch dataFacet {
		case DalleDressGenerator:
			if err := c.generatorFacet.Load(); err != nil {
				logging.LogError(fmt.Sprintf("LoadData.%s from store: %%v", dataFacet), err, facets.ErrAlreadyLoading)
			}
		case DalleDressSeries:
			if err := c.seriesFacet.Load(); err != nil {
				logging.LogError(fmt.Sprintf("LoadData.%s from store: %%v", dataFacet), err, facets.ErrAlreadyLoading)
			}
		case DalleDressDatabases:
			if err := c.databasesFacet.Load(); err != nil {
				logging.LogError(fmt.Sprintf("LoadData.%s from store: %%v", dataFacet), err, facets.ErrAlreadyLoading)
			}
		case DalleDressEvents:
			if err := c.eventsFacet.Load(); err != nil {
				logging.LogError(fmt.Sprintf("LoadData.%s from store: %%v", dataFacet), err, facets.ErrAlreadyLoading)
			}
		case DalleDressGallery:
			if err := c.galleryFacet.Load(); err != nil {
				logging.LogError(fmt.Sprintf("LoadData.%s from store: %%v", dataFacet), err, facets.ErrAlreadyLoading)
			}
		default:
			logging.LogError("LoadData: unexpected dataFacet: %v", fmt.Errorf("invalid dataFacet: %s", dataFacet), nil)
			return
		}
	}()
}

func (c *DalleDressCollection) Reset(dataFacet types.DataFacet) {
	switch dataFacet {
	case DalleDressGenerator:
		c.generatorFacet.GetStore().Reset()
	case DalleDressSeries:
		c.seriesFacet.GetStore().Reset()
	case DalleDressDatabases:
		c.databasesFacet.GetStore().Reset()
	case DalleDressEvents:
		c.eventsFacet.GetStore().Reset()
	case DalleDressGallery:
		c.galleryFacet.GetStore().Reset()
	default:
		return
	}
}

func (c *DalleDressCollection) NeedsUpdate(dataFacet types.DataFacet) bool {
	switch dataFacet {
	case DalleDressGenerator:
		return c.generatorFacet.NeedsUpdate()
	case DalleDressSeries:
		return c.seriesFacet.NeedsUpdate()
	case DalleDressDatabases:
		return c.databasesFacet.NeedsUpdate()
	case DalleDressEvents:
		return c.eventsFacet.NeedsUpdate()
	case DalleDressGallery:
		return c.galleryFacet.NeedsUpdate()
	default:
		return false
	}
}

func (c *DalleDressCollection) GetSupportedFacets() []types.DataFacet {
	return []types.DataFacet{
		DalleDressGenerator,
		DalleDressSeries,
		DalleDressDatabases,
		DalleDressEvents,
		DalleDressGallery,
	}
}

func (c *DalleDressCollection) AccumulateItem(item interface{}, summary *types.Summary) {
	// EXISTING_CODE
	// EXISTING_CODE
}

func (c *DalleDressCollection) GetSummary() types.Summary {
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

func (c *DalleDressCollection) ResetSummary() {
	c.summaryMutex.Lock()
	defer c.summaryMutex.Unlock()
	c.summary = types.Summary{
		TotalCount:  0,
		FacetCounts: make(map[types.DataFacet]int),
		CustomData:  make(map[string]interface{}),
		LastUpdated: time.Now().Unix(),
	}
}

func (c *DalleDressCollection) ExportData(payload *types.Payload) (string, error) {
	switch payload.DataFacet {
	case DalleDressGenerator:
		return c.generatorFacet.ExportData(payload, string(DalleDressGenerator))
	case DalleDressSeries:
		return c.seriesFacet.ExportData(payload, string(DalleDressSeries))
	case DalleDressDatabases:
		return c.databasesFacet.ExportData(payload, string(DalleDressDatabases))
	case DalleDressEvents:
		return c.eventsFacet.ExportData(payload, string(DalleDressEvents))
	case DalleDressGallery:
		return c.galleryFacet.ExportData(payload, string(DalleDressGallery))
	default:
		return "", fmt.Errorf("[ExportData] unsupported dalledress facet: %s", payload.DataFacet)
	}
}

// EXISTING_CODE
// RefreshGallery clears the gallery cache forcing a full rescan on next access
func (c *DalleDressCollection) RefreshGallery() {
	c.galleryCacheMux.Lock()
	c.galleryCache = nil
	c.gallerySeriesInfo = make(map[string]int64)
	c.galleryCacheMux.Unlock()
}

// getGalleryItems returns cached gallery items performing incremental scan per series
func (c *DalleDressCollection) getGalleryItems() (items []*GalleryItem) {
	root := dallev2OutputDir()

	// snapshot existing cache state
	c.galleryCacheMux.RLock()
	cached := c.galleryCache
	prevSeriesInfo := make(map[string]int64, len(c.gallerySeriesInfo))
	for k, v := range c.gallerySeriesInfo {
		prevSeriesInfo[k] = v
	}
	c.galleryCacheMux.RUnlock()

	current := make(map[string]int64)
	entries, err := os.ReadDir(root)
	if err == nil {
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			annotatedPath := filepath.Join(root, e.Name(), "annotated")
			if info, ierr := os.Stat(annotatedPath); ierr == nil && info.IsDir() {
				current[e.Name()] = info.ModTime().Unix()
			}
		}
	}

	changedSeries := make([]string, 0, len(current))
	if len(current) != len(prevSeriesInfo) {
		for series := range current {
			if prevSeriesInfo[series] != current[series] {
				changedSeries = append(changedSeries, series)
			}
		}
		// removed series handled by absence
	} else {
		for series, m := range current {
			if prevSeriesInfo[series] != m {
				changedSeries = append(changedSeries, series)
			}
		}
	}
	if len(changedSeries) == 0 && cached != nil {
		return cached
	}

	// build existing map for unchanged reuse
	existingBySeries := make(map[string][]*GalleryItem)
	for _, it := range cached {
		existingBySeries[it.Series] = append(existingBySeries[it.Series], it)
	}

	// If only one series changed, keep it simple sequentially
	merged := make([]*GalleryItem, 0, 512)
	if len(changedSeries) == 1 {
		changed := changedSeries[0]
		// reuse other series directly
		for series, items := range existingBySeries {
			if series == changed {
				continue
			}
			merged = append(merged, items...)
		}
		if seriesItems, err := collectGalleryItemsForSeries(root, changed); err == nil && len(seriesItems) > 0 {
			merged = append(merged, seriesItems...)
		}
	} else {
		// multi-series change -> parallelize rescans
		changedSet := make(map[string]struct{}, len(changedSeries))
		for _, s := range changedSeries {
			changedSet[s] = struct{}{}
		}
		for series, items := range existingBySeries { // keep unchanged first
			if _, ok := changedSet[series]; !ok {
				merged = append(merged, items...)
			}
		}
		workerCount := runtime.NumCPU()
		if workerCount > len(changedSeries) {
			workerCount = len(changedSeries)
		}
		if workerCount < 2 {
			workerCount = 2
		}
		jobs := make(chan string, len(changedSeries))
		results := make(chan []*GalleryItem, len(changedSeries))
		var wg sync.WaitGroup
		worker := func() {
			defer wg.Done()
			for series := range jobs {
				if seriesItems, err := collectGalleryItemsForSeries(root, series); err == nil && len(seriesItems) > 0 {
					results <- seriesItems
				} else {
					results <- nil
				}
			}
		}
		wg.Add(workerCount)
		for i := 0; i < workerCount; i++ {
			go worker()
		}
		for _, s := range changedSeries {
			jobs <- s
		}
		close(jobs)
		wg.Wait()
		close(results)
		for r := range results {
			if len(r) > 0 {
				merged = append(merged, r...)
			}
		}
	}

	// final sort
	sort.SliceStable(merged, func(i, j int) bool {
		if merged[i].Series == merged[j].Series {
			if merged[i].Index == merged[j].Index {
				return merged[i].FileName < merged[j].FileName
			}
			return merged[i].Index < merged[j].Index
		}
		return merged[i].Series < merged[j].Series
	})

	c.galleryCacheMux.Lock()
	c.galleryCache = merged
	c.gallerySeriesInfo = current
	c.galleryCacheMux.Unlock()
	return merged
}

// dallev2OutputDir isolates external call for easier test mocking
func dallev2OutputDir() string { return dallev2.OutputDir() }

// EXISTING_CODE
