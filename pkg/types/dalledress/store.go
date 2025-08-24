// Copyright 2016, 2026 The Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package dalledress

import (
	"path/filepath"
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

type Generator = dalle.Generator
type Database = dalle.Database
type Log = sdk.Log
type Series = dalle.Series

var (
	databasesStore   *store.Store[Database]
	databasesStoreMu sync.Mutex

	generatorStore   *store.Store[Generator]
	generatorStoreMu sync.Mutex

	logsStore   *store.Store[Log]
	logsStoreMu sync.Mutex

	seriesStore   *store.Store[Series]
	seriesStoreMu sync.Mutex
)

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

func (c *DalleDressCollection) getGeneratorStore(payload *types.Payload, facet types.DataFacet) *store.Store[Generator] {
	generatorStoreMu.Lock()
	defer generatorStoreMu.Unlock()

	// EXISTING_CODE
	// EXISTING_CODE

	chain := payload.Chain
	address := payload.Address
	theStore := generatorStore
	if theStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			// EXISTING_CODE
			// EXISTING_CODE
			return nil
		}

		processFunc := func(item interface{}) *Generator {
			if it, ok := item.(*Generator); ok {
				return it
			}
			return nil
		}

		mappingFunc := func(item *Generator) (key interface{}, includeInMap bool) {
			// EXISTING_CODE
			// EXISTING_CODE
			return nil, false
		}

		storeName := c.GetStoreName(facet, chain, address)
		theStore = store.NewStore(storeName, queryFunc, processFunc, mappingFunc)

		// EXISTING_CODE
		// EXISTING_CODE

		generatorStore = theStore
	}

	return theStore
}

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
			// EXISTING_CODE
			seriesDir := filepath.Join(dalle.DataDir(), "series")
			models, _ := dalle.LoadSeriesModels(seriesDir)
			dalle.SortSeries(models, "suffix", true)
			for i, m := range models {
				theStore.AddItem(m, i)
			}
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
		name = "dalledress-generator"
	case DalleDressSeries:
		name = "dalledress-series"
	case DalleDressDatabases:
		name = "dalledress-databases"
	case DalleDressEvents:
		name = "dalledress-logs"
	case DalleDressGallery:
		name = "dalledress-logs"
	default:
		return ""
	}
	return name
}

// TODO: THIS SHOULD BE PER STORE - SEE EXPORT COMMENTS
func GetDalleDressCount(payload *types.Payload) (int, error) {
	// chain := payload.Chain
	// countOpts := sdk.DalleDressOptions{
	// 	Globals: sdk.Globals{Cache: true, Chain: chain},
	// }
	// if countResult, _, err := countOpts.DalleDressCount(); err != nil {
	// 	return 0, fmt.Errorf("DalleDressCount query error: %v", err)
	// } else if len(countResult) > 0 {
	// 	return int(countResult[0].Count), nil
	// }
	return 0, nil
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
// EXISTING_CODE
