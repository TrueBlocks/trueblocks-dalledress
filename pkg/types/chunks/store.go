// Copyright 2016, 2025 The TrueBlocks Authors. All rights reserved.
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

	// EXISTING_CODE
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/logging"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/preferences"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/store"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
	// EXISTING_CODE
)

type Bloom = coreTypes.ChunkBloom
type Index = coreTypes.ChunkIndex
type Manifest = coreTypes.ChunkManifest
type Stats = coreTypes.ChunkStats

var (
	bloomsStore   *store.Store[Bloom]
	bloomsStoreMu sync.Mutex

	indexStore   *store.Store[Index]
	indexStoreMu sync.Mutex

	manifestStore   *store.Store[Manifest]
	manifestStoreMu sync.Mutex

	statsStore   *store.Store[Stats]
	statsStoreMu sync.Mutex
)

func (c *ChunksCollection) getBloomsStore() *store.Store[Bloom] {
	bloomsStoreMu.Lock()
	defer bloomsStoreMu.Unlock()

	if bloomsStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			// EXISTING_CODE
			chainName := preferences.GetChain()
			opts := sdk.ChunksOptions{
				Globals: sdk.Globals{
					Verbose: false, // Set to false to avoid weird output issues
					Chain:   chainName,
				},
				RenderCtx: ctx,
			}

			if _, _, err := opts.ChunksBlooms(); err != nil {
				// Create structured error with proper context
				wrappedErr := types.NewSDKError("chunks", ChunksBlooms, "fetch", err)
				logging.LogBackend(fmt.Sprintf("Chunks blooms SDK query error: %v", wrappedErr))
				return wrappedErr
			}
			// EXISTING_CODE
			return nil
		}

		processFunc := func(item interface{}) *Bloom {
			// EXISTING_CODE
			if bloom, ok := item.(*Bloom); ok {
				return bloom
			}
			// EXISTING_CODE
			return nil
		}

		mappingFunc := func(item *Bloom) (key interface{}, includeInMap bool) {
			// EXISTING_CODE
			if item != nil {
				return item.Range, true
			}
			// EXISTING_CODE
			return nil, false
		}

		// EXISTING_CODE
		storeName := c.GetStoreName(ChunksBlooms)
		// EXISTING_CODE
		bloomsStore = store.NewStore(storeName, queryFunc, processFunc, mappingFunc)
	}

	return bloomsStore
}

func (c *ChunksCollection) getIndexStore() *store.Store[Index] {
	indexStoreMu.Lock()
	defer indexStoreMu.Unlock()

	if indexStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			// EXISTING_CODE
			chainName := preferences.GetChain()
			opts := sdk.ChunksOptions{
				Globals: sdk.Globals{
					Verbose: true,
					Chain:   chainName,
				},
				RenderCtx: ctx,
			}

			if _, _, err := opts.ChunksIndex(); err != nil {
				// Create structured error with proper context
				wrappedErr := types.NewSDKError("chunks", ChunksIndex, "fetch", err)
				logging.LogBackend(fmt.Sprintf("Chunks index SDK query error: %v", wrappedErr))
				return wrappedErr
			}
			// EXISTING_CODE
			return nil
		}

		processFunc := func(item interface{}) *Index {
			// EXISTING_CODE
			if index, ok := item.(*Index); ok {
				return index
			}
			// EXISTING_CODE
			return nil
		}

		mappingFunc := func(item *Index) (key interface{}, includeInMap bool) {
			// EXISTING_CODE
			if item != nil {
				return nil, false
			}
			// EXISTING_CODE
			return item.Range, true
		}

		// EXISTING_CODE
		storeName := c.GetStoreName(ChunksIndex)
		// EXISTING_CODE
		indexStore = store.NewStore(storeName, queryFunc, processFunc, mappingFunc)
	}

	return indexStore
}

func (c *ChunksCollection) getManifestStore() *store.Store[Manifest] {
	manifestStoreMu.Lock()
	defer manifestStoreMu.Unlock()

	if manifestStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			// EXISTING_CODE
			chainName := preferences.GetChain()
			opts := sdk.ChunksOptions{
				Globals: sdk.Globals{
					Verbose: true,
					Chain:   chainName,
				},
				RenderCtx: ctx,
			}

			if _, _, err := opts.ChunksManifest(); err != nil {
				// Create structured error with proper context
				wrappedErr := types.NewSDKError("chunks", ChunksManifest, "fetch", err)
				logging.LogBackend(fmt.Sprintf("Chunks manifest SDK query error: %v", wrappedErr))
				return wrappedErr
			}
			// EXISTING_CODE
			return nil
		}

		processFunc := func(item interface{}) *Manifest {
			// EXISTING_CODE
			if manifest, ok := item.(*Manifest); ok {
				return manifest
			}
			// EXISTING_CODE
			return nil
		}

		mappingFunc := func(item *Manifest) (key interface{}, includeInMap bool) {
			// EXISTING_CODE
			if item != nil {
				return item.Chain, true
			}
			// EXISTING_CODE
			return nil, false
		}

		// EXISTING_CODE
		storeName := c.GetStoreName(ChunksManifest)
		// EXISTING_CODE
		manifestStore = store.NewStore(storeName, queryFunc, processFunc, mappingFunc)
	}

	return manifestStore
}

func (c *ChunksCollection) getStatsStore() *store.Store[Stats] {
	statsStoreMu.Lock()
	defer statsStoreMu.Unlock()

	if statsStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			// EXISTING_CODE
			chainName := preferences.GetChain()
			opts := sdk.ChunksOptions{
				Globals: sdk.Globals{
					Verbose: true,
					Chain:   chainName,
				},
				RenderCtx: ctx,
			}

			if _, _, err := opts.ChunksStats(); err != nil {
				// Create structured error with proper context
				wrappedErr := types.NewSDKError("chunks", ChunksStats, "fetch", err)
				logging.LogBackend(fmt.Sprintf("Chunks stats SDK query error: %v", wrappedErr))
				return wrappedErr
			}
			// EXISTING_CODE
			return nil
		}

		processFunc := func(item interface{}) *Stats {
			// EXISTING_CODE
			if stats, ok := item.(*Stats); ok {
				return stats
			}
			// EXISTING_CODE
			return nil
		}

		mappingFunc := func(item *Stats) (key interface{}, includeInMap bool) {
			// EXISTING_CODE
			if item != nil {
				return item.Range, true
			}
			// EXISTING_CODE
			return nil, false
		}

		// EXISTING_CODE
		storeName := c.GetStoreName(ChunksStats)
		// EXISTING_CODE
		statsStore = store.NewStore(storeName, queryFunc, processFunc, mappingFunc)
	}

	return statsStore
}

func (c *ChunksCollection) GetStoreName(dataFacet types.DataFacet) string {
	switch dataFacet {
	case ChunksStats:
		return "chunks-stats"
	case ChunksIndex:
		return "chunks-index"
	case ChunksBlooms:
		return "chunks-blooms"
	case ChunksManifest:
		return "chunks-manifest"
	default:
		return ""
	}
}

// EXISTING_CODE
func GetStatsCount() (int, error) {
	chainName := preferences.GetChain()

	// Use dedicated count method if available in SDK
	opts := sdk.ChunksOptions{
		Globals: sdk.Globals{Chain: chainName},
	}

	// Try to use count functionality if available
	countResults, _, err := opts.ChunksCount()
	if err == nil && len(countResults) > 0 {
		// Assuming the count result contains the total count
		return int(countResults[0].Count), nil
	}

	// Fallback to store count
	return statsStore.Count(), nil
}

func GetIndexCount() (int, error) {
	chainName := preferences.GetChain()

	opts := sdk.ChunksOptions{
		Globals: sdk.Globals{Chain: chainName},
	}

	// Try to use count functionality if available
	countResults, _, err := opts.ChunksCount()
	if err == nil && len(countResults) > 0 {
		return int(countResults[0].Count), nil
	}

	// Fallback to store count
	return indexStore.Count(), nil
}

func GetBloomsCount() (int, error) {
	chainName := preferences.GetChain()

	opts := sdk.ChunksOptions{
		Globals: sdk.Globals{Chain: chainName},
	}

	// Try to use count functionality if available
	countResults, _, err := opts.ChunksCount()
	if err == nil && len(countResults) > 0 {
		return int(countResults[0].Count), nil
	}

	// Fallback to store count
	return bloomsStore.Count(), nil
}

func GetManifestCount() (int, error) {
	chainName := preferences.GetChain()

	opts := sdk.ChunksOptions{
		Globals: sdk.Globals{Chain: chainName},
	}

	// Try to use count functionality if available
	countResults, _, err := opts.ChunksCount()
	if err == nil && len(countResults) > 0 {
		return int(countResults[0].Count), nil
	}

	// Fallback to store count
	return manifestStore.Count(), nil
}

// EXISTING_CODE
