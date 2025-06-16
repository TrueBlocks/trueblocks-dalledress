// CHUNKS_ROUTE
package chunks

import (
	"fmt"
	"sync"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/logging"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/preferences"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/store"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

var (
	chunksStatsStore    *store.Store[coreTypes.ChunkStats]
	chunksIndexStore    *store.Store[coreTypes.ChunkIndex]
	chunksBloomsStore   *store.Store[coreTypes.ChunkBloom]
	chunksManifestStore *store.Store[coreTypes.ChunkManifest]

	statsStoreMu    sync.Mutex
	indexStoreMu    sync.Mutex
	bloomsStoreMu   sync.Mutex
	manifestStoreMu sync.Mutex
)

// GetChunksStatsStore returns singleton store instance for chunk statistics
func GetChunksStatsStore() *store.Store[coreTypes.ChunkStats] {
	statsStoreMu.Lock()
	defer statsStoreMu.Unlock()

	if chunksStatsStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			chainName := preferences.GetChain()
			opts := sdk.ChunksOptions{
				Globals: sdk.Globals{
					Verbose: true,
					Chain:   chainName,
				},
				RenderCtx: ctx,
			}

			if _, _, err := opts.ChunksStats(); err != nil {
				logging.LogBackend(fmt.Sprintf("ChunksStatsStore query error: %v", err))
				return err
			}

			return nil
		}

		processFunc := func(item interface{}) *coreTypes.ChunkStats {
			if stats, ok := item.(*coreTypes.ChunkStats); ok {
				return stats
			}
			return nil
		}

		mappingFunc := func(item *coreTypes.ChunkStats) (key interface{}, includeInMap bool) {
			if item == nil {
				return nil, false
			}
			// Use Range as unique key for chunk stats
			return item.Range, true
		}

		chunksStatsStore = store.NewStore("chunks-stats", queryFunc, processFunc, mappingFunc)
	}

	return chunksStatsStore
}

// GetChunksIndexStore returns singleton store instance for chunk index
func GetChunksIndexStore() *store.Store[coreTypes.ChunkIndex] {
	indexStoreMu.Lock()
	defer indexStoreMu.Unlock()

	if chunksIndexStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			chainName := preferences.GetChain()
			opts := sdk.ChunksOptions{
				Globals: sdk.Globals{
					Verbose: true,
					Chain:   chainName,
				},
				RenderCtx: ctx,
			}

			if _, _, err := opts.ChunksIndex(); err != nil {
				logging.LogBackend(fmt.Sprintf("ChunksIndexStore query error: %v", err))
				return err
			}

			return nil
		}

		processFunc := func(item interface{}) *coreTypes.ChunkIndex {
			if index, ok := item.(*coreTypes.ChunkIndex); ok {
				return index
			}
			return nil
		}

		mappingFunc := func(item *coreTypes.ChunkIndex) (key interface{}, includeInMap bool) {
			if item == nil {
				return nil, false
			}
			// Use Range as unique key for chunk index
			return item.Range, true
		}

		chunksIndexStore = store.NewStore("chunks-index", queryFunc, processFunc, mappingFunc)
	}

	return chunksIndexStore
}

// GetChunksBloomsStore returns singleton store instance for chunk blooms
func GetChunksBloomsStore() *store.Store[coreTypes.ChunkBloom] {
	bloomsStoreMu.Lock()
	defer bloomsStoreMu.Unlock()

	if chunksBloomsStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			chainName := preferences.GetChain()
			opts := sdk.ChunksOptions{
				Globals: sdk.Globals{
					Verbose: true,
					Chain:   chainName,
				},
				RenderCtx: ctx,
			}

			if _, _, err := opts.ChunksBlooms(); err != nil {
				logging.LogBackend(fmt.Sprintf("ChunksBloomsStore query error: %v", err))
				return err
			}

			return nil
		}

		processFunc := func(item interface{}) *coreTypes.ChunkBloom {
			if bloom, ok := item.(*coreTypes.ChunkBloom); ok {
				return bloom
			}
			return nil
		}

		mappingFunc := func(item *coreTypes.ChunkBloom) (key interface{}, includeInMap bool) {
			if item == nil {
				return nil, false
			}
			// Use Range as unique key for chunk bloom
			return item.Range, true
		}

		chunksBloomsStore = store.NewStore("chunks-blooms", queryFunc, processFunc, mappingFunc)
	}

	return chunksBloomsStore
}

// GetChunksManifestStore returns singleton store instance for chunk manifest
func GetChunksManifestStore() *store.Store[coreTypes.ChunkManifest] {
	manifestStoreMu.Lock()
	defer manifestStoreMu.Unlock()

	if chunksManifestStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			chainName := preferences.GetChain()
			opts := sdk.ChunksOptions{
				Globals: sdk.Globals{
					Verbose: true,
					Chain:   chainName,
				},
				RenderCtx: ctx,
			}

			if _, _, err := opts.ChunksManifest(); err != nil {
				logging.LogBackend(fmt.Sprintf("ChunksManifestStore query error: %v", err))
				return err
			}

			return nil
		}

		processFunc := func(item interface{}) *coreTypes.ChunkManifest {
			if manifest, ok := item.(*coreTypes.ChunkManifest); ok {
				return manifest
			}
			return nil
		}

		mappingFunc := func(item *coreTypes.ChunkManifest) (key interface{}, includeInMap bool) {
			if item == nil {
				return nil, false
			}
			// Use Chain as unique key for chunk manifest
			return item.Chain, true
		}

		chunksManifestStore = store.NewStore("chunks-manifest", queryFunc, processFunc, mappingFunc)
	}

	return chunksManifestStore
}

// GetChunksStatsCount provides optimized count for stats
func GetChunksStatsCount() (int, error) {
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
	return GetChunksStatsStore().Count(), nil
}

// GetChunksIndexCount provides optimized count for index
func GetChunksIndexCount() (int, error) {
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
	return GetChunksIndexStore().Count(), nil
}

// GetChunksBloomsCount provides optimized count for blooms
func GetChunksBloomsCount() (int, error) {
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
	return GetChunksBloomsStore().Count(), nil
}

// GetChunksManifestCount provides optimized count for manifest
func GetChunksManifestCount() (int, error) {
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
	return GetChunksManifestStore().Count(), nil
}

// CHUNKS_ROUTE
