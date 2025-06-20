// CHUNKS_ROUTE
package chunks

import (
	"fmt"
	"strings"
	"sync"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/facets"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/logging"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
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

// Page structures for frontend consumption
type ChunksStatsPage struct {
	Facet         types.DataFacet         `json:"facet"`
	ChunksStats   []*coreTypes.ChunkStats `json:"chunksStats"`
	TotalItems    int                     `json:"totalItems"`
	ExpectedTotal int                     `json:"expectedTotal"`
	IsFetching    bool                    `json:"isFetching"`
	State         types.LoadState         `json:"state"`
}

type ChunksIndexPage struct {
	Facet         types.DataFacet         `json:"facet"`
	ChunksIndex   []*coreTypes.ChunkIndex `json:"chunksIndex"`
	TotalItems    int                     `json:"totalItems"`
	ExpectedTotal int                     `json:"expectedTotal"`
	IsFetching    bool                    `json:"isFetching"`
	State         types.LoadState         `json:"state"`
}

type ChunksBloomsPage struct {
	Facet         types.DataFacet         `json:"facet"`
	ChunksBlooms  []*coreTypes.ChunkBloom `json:"chunksBlooms"`
	TotalItems    int                     `json:"totalItems"`
	ExpectedTotal int                     `json:"expectedTotal"`
	IsFetching    bool                    `json:"isFetching"`
	State         types.LoadState         `json:"state"`
}

type ChunksManifestPage struct {
	Facet          types.DataFacet            `json:"facet"`
	ChunksManifest []*coreTypes.ChunkManifest `json:"chunksManifest"`
	TotalItems     int                        `json:"totalItems"`
	ExpectedTotal  int                        `json:"expectedTotal"`
	IsFetching     bool                       `json:"isFetching"`
	State          types.LoadState            `json:"state"`
}

// Implement Page interface for all page types
func (csp *ChunksStatsPage) GetFacet() types.DataFacet { return csp.Facet }
func (csp *ChunksStatsPage) GetTotalItems() int        { return csp.TotalItems }
func (csp *ChunksStatsPage) GetExpectedTotal() int     { return csp.ExpectedTotal }
func (csp *ChunksStatsPage) GetIsFetching() bool       { return csp.IsFetching }
func (csp *ChunksStatsPage) GetState() types.LoadState { return csp.State }

func (cip *ChunksIndexPage) GetFacet() types.DataFacet { return cip.Facet }
func (cip *ChunksIndexPage) GetTotalItems() int        { return cip.TotalItems }
func (cip *ChunksIndexPage) GetExpectedTotal() int     { return cip.ExpectedTotal }
func (cip *ChunksIndexPage) GetIsFetching() bool       { return cip.IsFetching }
func (cip *ChunksIndexPage) GetState() types.LoadState { return cip.State }

func (cbp *ChunksBloomsPage) GetFacet() types.DataFacet { return cbp.Facet }
func (cbp *ChunksBloomsPage) GetTotalItems() int        { return cbp.TotalItems }
func (cbp *ChunksBloomsPage) GetExpectedTotal() int     { return cbp.ExpectedTotal }
func (cbp *ChunksBloomsPage) GetIsFetching() bool       { return cbp.IsFetching }
func (cbp *ChunksBloomsPage) GetState() types.LoadState { return cbp.State }

func (cmp *ChunksManifestPage) GetFacet() types.DataFacet { return cmp.Facet }
func (cmp *ChunksManifestPage) GetTotalItems() int        { return cmp.TotalItems }
func (cmp *ChunksManifestPage) GetExpectedTotal() int     { return cmp.ExpectedTotal }
func (cmp *ChunksManifestPage) GetIsFetching() bool       { return cmp.IsFetching }
func (cmp *ChunksManifestPage) GetState() types.LoadState { return cmp.State }

// ChunksPage is a union type that can represent any chunks page type
type ChunksPage struct {
	Facet          types.DataFacet            `json:"facet"`
	ChunksStats    []*coreTypes.ChunkStats    `json:"chunksStats,omitempty"`
	ChunksIndex    []*coreTypes.ChunkIndex    `json:"chunksIndex,omitempty"`
	ChunksBlooms   []*coreTypes.ChunkBloom    `json:"chunksBlooms,omitempty"`
	ChunksManifest []*coreTypes.ChunkManifest `json:"chunksManifest,omitempty"`
	TotalItems     int                        `json:"totalItems"`
	ExpectedTotal  int                        `json:"expectedTotal"`
	IsFetching     bool                       `json:"isFetching"`
	State          types.LoadState            `json:"state"`
}

func (cp *ChunksPage) GetFacet() types.DataFacet {
	return cp.Facet
}

func (cp *ChunksPage) GetTotalItems() int {
	return cp.TotalItems
}

func (cp *ChunksPage) GetExpectedTotal() int {
	return cp.ExpectedTotal
}

func (cp *ChunksPage) GetIsFetching() bool {
	return cp.IsFetching
}

func (cp *ChunksPage) GetState() types.LoadState {
	return cp.State
}

// Collection structure with multiple facets
type ChunksCollection struct {
	// Facets for different chunk views
	statsFacet    *facets.Facet[coreTypes.ChunkStats]
	indexFacet    *facets.Facet[coreTypes.ChunkIndex]
	bloomsFacet   *facets.Facet[coreTypes.ChunkBloom]
	manifestFacet *facets.Facet[coreTypes.ChunkManifest]
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

// GetPage implements the Collection interface
func (cc *ChunksCollection) GetPage(
	dataFacet types.DataFacet,
	first, pageSize int,
	sortSpec sdk.SortSpec,
	filter string,
) (types.Page, error) {
	page := &ChunksPage{
		Facet: dataFacet,
	}
	filter = strings.ToLower(filter)

	switch dataFacet {
	case ChunksStats:
		var filterFunc func(*coreTypes.ChunkStats) bool
		if filter != "" {
			filterFunc = func(stat *coreTypes.ChunkStats) bool {
				return cc.matchesStatsFilter(stat, filter)
			}
		}

		sortFunc := func(items []coreTypes.ChunkStats, sort sdk.SortSpec) error {
			// Placeholder - implement actual sorting if needed
			return nil
		}

		if result, err := cc.statsFacet.GetPage(first, pageSize, filterFunc, sortSpec, sortFunc); err != nil {
			// This is likely an SDK or store error, not a validation error
			return nil, types.NewStoreError("chunks", dataFacet, "GetPage", err)
		} else {
			chunksStats := make([]*coreTypes.ChunkStats, 0, len(result.Items))
			for i := range result.Items {
				chunksStats = append(chunksStats, &result.Items[i])
			}
			page.ChunksStats, page.TotalItems, page.State = chunksStats, result.TotalItems, result.State
		}
		page.IsFetching = cc.statsFacet.IsFetching()
		page.ExpectedTotal = cc.statsFacet.ExpectedCount()

	case ChunksIndex:
		var filterFunc func(*coreTypes.ChunkIndex) bool
		if filter != "" {
			filterFunc = func(index *coreTypes.ChunkIndex) bool {
				return cc.matchesIndexFilter(index, filter)
			}
		}

		sortFunc := func(items []coreTypes.ChunkIndex, sort sdk.SortSpec) error {
			// Placeholder - implement actual sorting if needed
			return nil
		}

		if result, err := cc.indexFacet.GetPage(first, pageSize, filterFunc, sortSpec, sortFunc); err != nil {
			// This is likely an SDK or store error, not a validation error
			return nil, types.NewStoreError("chunks", dataFacet, "GetPage", err)
		} else {
			chunksIndex := make([]*coreTypes.ChunkIndex, 0, len(result.Items))
			for i := range result.Items {
				chunksIndex = append(chunksIndex, &result.Items[i])
			}
			page.ChunksIndex, page.TotalItems, page.State = chunksIndex, result.TotalItems, result.State
		}
		page.IsFetching = cc.indexFacet.IsFetching()
		page.ExpectedTotal = cc.indexFacet.ExpectedCount()

	case ChunksBlooms:
		var filterFunc func(*coreTypes.ChunkBloom) bool
		if filter != "" {
			filterFunc = func(bloom *coreTypes.ChunkBloom) bool {
				return cc.matchesBloomsFilter(bloom, filter)
			}
		}

		sortFunc := func(items []coreTypes.ChunkBloom, sort sdk.SortSpec) error {
			// Placeholder - implement actual sorting if needed
			return nil
		}

		if result, err := cc.bloomsFacet.GetPage(first, pageSize, filterFunc, sortSpec, sortFunc); err != nil {
			// This is likely an SDK or store error, not a validation error
			return nil, types.NewStoreError("chunks", dataFacet, "GetPage", err)
		} else {
			chunksBlooms := make([]*coreTypes.ChunkBloom, 0, len(result.Items))
			for i := range result.Items {
				chunksBlooms = append(chunksBlooms, &result.Items[i])
			}
			page.ChunksBlooms, page.TotalItems, page.State = chunksBlooms, result.TotalItems, result.State
		}
		page.IsFetching = cc.bloomsFacet.IsFetching()
		page.ExpectedTotal = cc.bloomsFacet.ExpectedCount()

	case ChunksManifest:
		var filterFunc func(*coreTypes.ChunkManifest) bool
		if filter != "" {
			filterFunc = func(manifest *coreTypes.ChunkManifest) bool {
				return cc.matchesManifestFilter(manifest, filter)
			}
		}

		sortFunc := func(items []coreTypes.ChunkManifest, sort sdk.SortSpec) error {
			// Placeholder - implement actual sorting if needed
			return nil
		}

		if result, err := cc.manifestFacet.GetPage(first, pageSize, filterFunc, sortSpec, sortFunc); err != nil {
			// This is likely an SDK or store error, not a validation error
			return nil, types.NewStoreError("chunks", dataFacet, "GetPage", err)
		} else {
			chunksManifest := make([]*coreTypes.ChunkManifest, 0, len(result.Items))
			for i := range result.Items {
				chunksManifest = append(chunksManifest, &result.Items[i])
			}
			page.ChunksManifest, page.TotalItems, page.State = chunksManifest, result.TotalItems, result.State
		}
		page.IsFetching = cc.manifestFacet.IsFetching()
		page.ExpectedTotal = cc.manifestFacet.ExpectedCount()

	default:
		// This is truly a validation error - invalid DataFacet for this collection
		return nil, types.NewValidationError("chunks", dataFacet, "GetPage",
			fmt.Errorf("unsupported data facet: %v", dataFacet))
	}

	return page, nil
}

// Filter functions for each facet
func (cc *ChunksCollection) matchesStatsFilter(stat *coreTypes.ChunkStats, filter string) bool {
	filterLower := strings.ToLower(filter)

	// Filter by various fields in ChunkStats
	if strings.Contains(strings.ToLower(stat.Range), filterLower) {
		return true
	}
	if strings.Contains(strings.ToLower(fmt.Sprintf("%d", stat.NAddrs)), filterLower) {
		return true
	}
	if strings.Contains(strings.ToLower(fmt.Sprintf("%d", stat.NApps)), filterLower) {
		return true
	}
	if strings.Contains(strings.ToLower(fmt.Sprintf("%d", stat.NBlocks)), filterLower) {
		return true
	}

	return false
}

func (cc *ChunksCollection) matchesIndexFilter(index *coreTypes.ChunkIndex, filter string) bool {
	filterLower := strings.ToLower(filter)

	// Filter by various fields in ChunkIndex
	if strings.Contains(strings.ToLower(index.Range), filterLower) {
		return true
	}
	if strings.Contains(strings.ToLower(index.Hash.String()), filterLower) {
		return true
	}
	if strings.Contains(strings.ToLower(fmt.Sprintf("%d", index.NAddresses)), filterLower) {
		return true
	}
	if strings.Contains(strings.ToLower(fmt.Sprintf("%d", index.NAppearances)), filterLower) {
		return true
	}
	if strings.Contains(strings.ToLower(index.Magic), filterLower) {
		return true
	}

	return false
}

func (cc *ChunksCollection) matchesBloomsFilter(bloom *coreTypes.ChunkBloom, filter string) bool {
	filterLower := strings.ToLower(filter)

	// Filter by various fields in ChunkBloom
	if strings.Contains(strings.ToLower(bloom.Range), filterLower) {
		return true
	}
	if strings.Contains(strings.ToLower(bloom.Hash.String()), filterLower) {
		return true
	}
	if strings.Contains(strings.ToLower(fmt.Sprintf("%d", bloom.NBlooms)), filterLower) {
		return true
	}
	if strings.Contains(strings.ToLower(fmt.Sprintf("%d", bloom.NInserted)), filterLower) {
		return true
	}
	if strings.Contains(strings.ToLower(bloom.Magic), filterLower) {
		return true
	}

	return false
}

func (cc *ChunksCollection) matchesManifestFilter(manifest *coreTypes.ChunkManifest, filter string) bool {
	filterLower := strings.ToLower(filter)

	// Filter by various fields in ChunkManifest
	if strings.Contains(strings.ToLower(manifest.Version), filterLower) {
		return true
	}
	if strings.Contains(strings.ToLower(manifest.Chain), filterLower) {
		return true
	}
	if strings.Contains(strings.ToLower(manifest.Specification.String()), filterLower) {
		return true
	}

	return false
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
	case *coreTypes.ChunkStats:
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

	case *coreTypes.ChunkIndex:
		summary.TotalCount++
		summary.FacetCounts[ChunksIndex]++
		if summary.CustomData == nil {
			summary.CustomData = make(map[string]interface{})
		}

		indexCount, _ := summary.CustomData["indexCount"].(int)
		indexCount++
		summary.CustomData["indexCount"] = indexCount

	case *coreTypes.ChunkBloom:
		summary.TotalCount++
		summary.FacetCounts[ChunksBlooms]++
		if summary.CustomData == nil {
			summary.CustomData = make(map[string]interface{})
		}

		bloomsCount, _ := summary.CustomData["bloomsCount"].(int)
		bloomsCount++
		summary.CustomData["bloomsCount"] = bloomsCount

	case *coreTypes.ChunkManifest:
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

func (cc *ChunksCollection) GetCurrentSummary() types.Summary {
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

func (cc *ChunksCollection) GetSummary() types.Summary {
	return cc.GetCurrentSummary()
}
