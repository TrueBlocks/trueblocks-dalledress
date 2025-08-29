// Copyright 2016, 2026 The Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package dalledress

import (
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"sync"

	// EXISTING_CODE
	dalle "github.com/TrueBlocks/trueblocks-dalle/v2"
	// EXISTING_CODE

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/store"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

// EXISTING_CODE
// EXISTING_CODE

type Log = sdk.Log
type DalleDress = dalle.DalleDress
type Database = dalle.Database
type Series = dalle.Series

var (
	logsStore   *store.Store[Log]
	logsStoreMu sync.Mutex

	dalledressStore   *store.Store[DalleDress]
	dalledressStoreMu sync.Mutex

	databasesStore   *store.Store[Database]
	databasesStoreMu sync.Mutex

	seriesStore   *store.Store[Series]
	seriesStoreMu sync.Mutex
)

func (c *DalleDressCollection) getLogsStore(payload *types.Payload, facet types.DataFacet) *store.Store[Log] {
	logsStoreMu.Lock()
	defer logsStoreMu.Unlock()

	// EXISTING_CODE
	// EXISTING_CODE

	chain := payload.Chain
	address := payload.Address
	theStore := logsStore
	if theStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			// EXISTING_CODE
			// EXISTING_CODE
			return nil
		}

		processFunc := func(item interface{}) *Log {
			if it, ok := item.(*Log); ok {
				return it
			}
			return nil
		}

		mappingFunc := func(item *Log) (key interface{}, includeInMap bool) {
			// EXISTING_CODE
			// EXISTING_CODE
			return nil, false
		}

		storeName := c.GetStoreName(facet, chain, address)
		theStore = store.NewStore(storeName, queryFunc, processFunc, mappingFunc)

		// EXISTING_CODE
		// EXISTING_CODE

		logsStore = theStore
	}

	return theStore
}

func (c *DalleDressCollection) getDalleDressStore(payload *types.Payload, facet types.DataFacet) *store.Store[DalleDress] {
	dalledressStoreMu.Lock()
	defer dalledressStoreMu.Unlock()

	// EXISTING_CODE
	// EXISTING_CODE

	chain := payload.Chain
	address := payload.Address
	theStore := dalledressStore
	if theStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			// EXISTING_CODE
			// EXISTING_CODE
			// need query
			return nil
		}

		processFunc := func(item interface{}) *DalleDress {
			if it, ok := item.(*DalleDress); ok {
				return it
			}
			return nil
		}

		mappingFunc := func(item *DalleDress) (key interface{}, includeInMap bool) {
			// EXISTING_CODE
			if item != nil {
				k := item.Original + ":" + item.AnnotatedPath
				return k, true
			}
			// EXISTING_CODE
			return nil, false
		}

		storeName := c.GetStoreName(facet, chain, address)
		theStore = store.NewStore(storeName, queryFunc, processFunc, mappingFunc)

		// EXISTING_CODE
		// EXISTING_CODE

		dalledressStore = theStore
	}

	return theStore
}

func (c *DalleDressCollection) getDatabasesStore(payload *types.Payload, facet types.DataFacet) *store.Store[Database] {
	databasesStoreMu.Lock()
	defer databasesStoreMu.Unlock()

	// EXISTING_CODE
	// EXISTING_CODE

	chain := payload.Chain
	address := payload.Address
	theStore := databasesStore
	if theStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			// EXISTING_CODE
			// EXISTING_CODE
			return nil
		}

		processFunc := func(item interface{}) *Database {
			if it, ok := item.(*Database); ok {
				return it
			}
			return nil
		}

		mappingFunc := func(item *Database) (key interface{}, includeInMap bool) {
			// EXISTING_CODE
			// EXISTING_CODE
			return nil, false
		}

		storeName := c.GetStoreName(facet, chain, address)
		theStore = store.NewStore(storeName, queryFunc, processFunc, mappingFunc)

		// EXISTING_CODE
		// EXISTING_CODE

		databasesStore = theStore
	}

	return theStore
}

func (c *DalleDressCollection) getSeriesStore(payload *types.Payload, facet types.DataFacet) *store.Store[Series] {
	seriesStoreMu.Lock()
	defer seriesStoreMu.Unlock()

	// EXISTING_CODE
	// EXISTING_CODE

	chain := payload.Chain
	address := payload.Address
	theStore := seriesStore
	if theStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			// EXISTING_CODE
			seriesDir := filepath.Join(dalle.DataDir(), "series")
			models, _ := dalle.LoadSeriesModels(seriesDir)
			_ = dalle.SortSeries(models, sdk.SortSpec{
				Fields: []string{"suffix"},
				Order:  []sdk.SortOrder{sdk.Asc},
			},
			)
			for i, m := range models {
				theStore.AddItem(&m, i)
			}
			// EXISTING_CODE
			return nil
		}

		processFunc := func(item interface{}) *Series {
			if it, ok := item.(*Series); ok {
				return it
			}
			return nil
		}

		mappingFunc := func(item *Series) (key interface{}, includeInMap bool) {
			// EXISTING_CODE
			if item != nil && item.Suffix != "" {
				return item.Suffix, true
			}
			// EXISTING_CODE
			return nil, false
		}

		storeName := c.GetStoreName(facet, chain, address)
		theStore = store.NewStore(storeName, queryFunc, processFunc, mappingFunc)

		// EXISTING_CODE
		// EXISTING_CODE

		seriesStore = theStore
	}

	return theStore
}

func (c *DalleDressCollection) GetStoreName(dataFacet types.DataFacet, chain, address string) string {
	_ = chain
	_ = address
	name := ""
	switch dataFacet {
	case DalleDressGenerator:
		name = "dalledress-dalledress"
	case DalleDressSeries:
		name = "dalledress-series"
	case DalleDressDatabases:
		name = "dalledress-databases"
	case DalleDressEvents:
		name = "dalledress-logs"
	case DalleDressGallery:
		name = "dalledress-dalledress"
	default:
		return ""
	}
	return name
}

var (
	collections   = make(map[store.CollectionKey]*DalleDressCollection)
	collectionsMu sync.Mutex
)

func GetDalleDressCollection(payload *types.Payload) *DalleDressCollection {
	collectionsMu.Lock()
	defer collectionsMu.Unlock()

	pl := *payload
	pl.Address = ""

	key := store.GetCollectionKey(&pl)
	if collection, exists := collections[key]; exists {
		return collection
	}

	collection := NewDalleDressCollection(payload)
	collections[key] = collection
	return collection
}

// EXISTING_CODE
// getGalleryItems returns cached gallery items performing incremental scan per series
func (c *DalleDressCollection) getGalleryItems() (items []*GalleryItem) {
	root := dalle.OutputDir()

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

// EXISTING_CODE
