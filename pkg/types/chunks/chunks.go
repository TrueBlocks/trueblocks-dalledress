// CHUNKS_ROUTE
package chunks

import (
	"fmt"
	"sync"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/facets"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/logging"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

// Define DataFacet constants for different chunk views
const (
	ChunksStats    types.DataFacet = "stats"
	ChunksIndex    types.DataFacet = "index"
	ChunksBlooms   types.DataFacet = "blooms"
	ChunksManifest types.DataFacet = "manifest"
)

// REQUIRED: Register DataFacet values for validation
func init() {
	types.RegisterDataFacet(ChunksStats)
	types.RegisterDataFacet(ChunksIndex)
	types.RegisterDataFacet(ChunksBlooms)
	types.RegisterDataFacet(ChunksManifest)
}

// Collection structure with multiple facets
type ChunksCollection struct {
	// Facets for different chunk views
	statsFacet    *facets.Facet[Stats]
	indexFacet    *facets.Facet[Index]
	bloomsFacet   *facets.Facet[Bloom]
	manifestFacet *facets.Facet[Manifest]
	summary       types.Summary
	summaryMutex  sync.RWMutex
}

// NewChunksCollection constructor function to initialize the chunks collection with all its facets
func NewChunksCollection() *ChunksCollection {
	cc := &ChunksCollection{
		summary: types.Summary{
			TotalCount:  0,
			FacetCounts: make(map[types.DataFacet]int),
			CustomData:  make(map[string]interface{}),
		},
	}

	cc.initializeFacets()
	return cc
}

// Initialize facets with appropriate filters
func (cc *ChunksCollection) initializeFacets() {
	// Stats facet - shows chunk statistics
	cc.statsFacet = facets.NewFacetWithSummary(
		ChunksStats,
		nil,
		nil,
		GetChunksStatsStore(),
		"chunks",
		cc,
	)

	// Index facet - shows index chunk information
	cc.indexFacet = facets.NewFacetWithSummary(
		ChunksIndex,
		nil,
		nil,
		GetChunksIndexStore(),
		"chunks",
		cc,
	)

	// Blooms facet - shows bloom filter information
	cc.bloomsFacet = facets.NewFacetWithSummary(
		ChunksBlooms,
		nil,
		nil,
		GetChunksBloomsStore(),
		"chunks",
		cc,
	)

	// Manifest facet - shows manifest information
	cc.manifestFacet = facets.NewFacetWithSummary(
		ChunksManifest,
		nil,
		nil,
		GetChunksManifestStore(),
		"chunks",
		cc,
	)
}

// Implement Collection interface methods
func (cc *ChunksCollection) LoadData(dataFacet types.DataFacet) {
	if !cc.NeedsUpdate(dataFacet) {
		return
	}

	var facet interface {
		Load() error
	}
	var facetName string

	switch dataFacet {
	case ChunksStats:
		facet = cc.statsFacet
		facetName = "chunks.stats"
	case ChunksIndex:
		facet = cc.indexFacet
		facetName = "chunks.index"
	case ChunksBlooms:
		facet = cc.bloomsFacet
		facetName = "chunks.blooms"
	case ChunksManifest:
		facet = cc.manifestFacet
		facetName = "chunks.manifest"
	default:
		logging.LogError("LoadData: unexpected data facet: %v", fmt.Errorf("invalid data facet: %s", dataFacet), nil)
		return
	}

	go func() {
		if err := facet.Load(); err != nil {
			logging.LogError(fmt.Sprintf("LoadData.%s from store: %%v", facetName), err, facets.ErrAlreadyLoading)
		}
	}()
}

func (cc *ChunksCollection) Reset(dataFacet types.DataFacet) {
	switch dataFacet {
	case ChunksStats:
		GetChunksStatsStore().Reset()
		cc.statsFacet.Reset()
	case ChunksIndex:
		GetChunksIndexStore().Reset()
		cc.indexFacet.Reset()
	case ChunksBlooms:
		GetChunksBloomsStore().Reset()
		cc.bloomsFacet.Reset()
	case ChunksManifest:
		GetChunksManifestStore().Reset()
		cc.manifestFacet.Reset()
	}
}

func (cc *ChunksCollection) NeedsUpdate(dataFacet types.DataFacet) bool {
	switch dataFacet {
	case ChunksStats:
		return cc.statsFacet.NeedsUpdate()
	case ChunksIndex:
		return cc.indexFacet.NeedsUpdate()
	case ChunksBlooms:
		return cc.bloomsFacet.NeedsUpdate()
	case ChunksManifest:
		return cc.manifestFacet.NeedsUpdate()
	default:
		return false
	}
}

func (cc *ChunksCollection) GetSupportedFacets() []types.DataFacet {
	return []types.DataFacet{ChunksStats, ChunksIndex, ChunksBlooms, ChunksManifest}
}

func (cc *ChunksCollection) GetStoreForFacet(dataFacet types.DataFacet) string {
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

func (cc *ChunksCollection) GetCollectionName() string {
	return "chunks"
}

// Crud implements the Collection interface with no-op operations since chunks data is immutable
func (cc *ChunksCollection) Crud(dataFacet types.DataFacet, op crud.Operation, item interface{}) error {
	// All CRUD operations are no-ops for chunks since the data is immutable blockchain data
	return nil
}

func (cc *ChunksCollection) AccumulateItem(item interface{}, summary *types.Summary) {
	cc.summaryMutex.Lock()
	defer cc.summaryMutex.Unlock()

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
}

func (cc *ChunksCollection) GetSummary() types.Summary {
	cc.summaryMutex.RLock()
	defer cc.summaryMutex.RUnlock()

	summary := cc.summary
	summary.FacetCounts = make(map[types.DataFacet]int)
	for k, v := range cc.summary.FacetCounts {
		summary.FacetCounts[k] = v
	}

	if cc.summary.CustomData != nil {
		summary.CustomData = make(map[string]interface{})
		for k, v := range cc.summary.CustomData {
			summary.CustomData[k] = v
		}
	}

	return summary
}

func (cc *ChunksCollection) ResetSummary() {
	cc.summaryMutex.Lock()
	defer cc.summaryMutex.Unlock()
	cc.summary = types.Summary{
		TotalCount:  0,
		FacetCounts: make(map[types.DataFacet]int),
		CustomData:  make(map[string]interface{}),
		LastUpdated: 0,
	}
}
