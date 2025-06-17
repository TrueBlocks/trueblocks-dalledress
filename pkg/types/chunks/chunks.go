// CHUNKS_ROUTE
package chunks

import (
	"fmt"
	"strings"
	"sync"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/facets"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

// Define ListKind constants for different chunk views
const (
	ChunksStats    types.ListKind = "Stats"
	ChunksIndex    types.ListKind = "Index"
	ChunksBlooms   types.ListKind = "Blooms"
	ChunksManifest types.ListKind = "Manifest"
)

// REQUIRED: Register kinds for validation
func init() {
	types.RegisterKind(ChunksStats)
	types.RegisterKind(ChunksIndex)
	types.RegisterKind(ChunksBlooms)
	types.RegisterKind(ChunksManifest)
}

// Page structures for frontend consumption
type ChunksStatsPage struct {
	Kind          types.ListKind          `json:"kind"`
	ChunksStats   []*coreTypes.ChunkStats `json:"chunksStats"`
	TotalItems    int                     `json:"totalItems"`
	ExpectedTotal int                     `json:"expectedTotal"`
	IsFetching    bool                    `json:"isFetching"`
	State         types.LoadState         `json:"state"`
}

type ChunksIndexPage struct {
	Kind          types.ListKind          `json:"kind"`
	ChunksIndex   []*coreTypes.ChunkIndex `json:"chunksIndex"`
	TotalItems    int                     `json:"totalItems"`
	ExpectedTotal int                     `json:"expectedTotal"`
	IsFetching    bool                    `json:"isFetching"`
	State         types.LoadState         `json:"state"`
}

type ChunksBloomsPage struct {
	Kind          types.ListKind          `json:"kind"`
	ChunksBlooms  []*coreTypes.ChunkBloom `json:"chunksBlooms"`
	TotalItems    int                     `json:"totalItems"`
	ExpectedTotal int                     `json:"expectedTotal"`
	IsFetching    bool                    `json:"isFetching"`
	State         types.LoadState         `json:"state"`
}

type ChunksManifestPage struct {
	Kind           types.ListKind             `json:"kind"`
	ChunksManifest []*coreTypes.ChunkManifest `json:"chunksManifest"`
	TotalItems     int                        `json:"totalItems"`
	ExpectedTotal  int                        `json:"expectedTotal"`
	IsFetching     bool                       `json:"isFetching"`
	State          types.LoadState            `json:"state"`
}

// Implement Page interface for all page types
func (csp *ChunksStatsPage) GetKind() types.ListKind   { return csp.Kind }
func (csp *ChunksStatsPage) GetTotalItems() int        { return csp.TotalItems }
func (csp *ChunksStatsPage) GetExpectedTotal() int     { return csp.ExpectedTotal }
func (csp *ChunksStatsPage) GetIsFetching() bool       { return csp.IsFetching }
func (csp *ChunksStatsPage) GetState() types.LoadState { return csp.State }

func (cip *ChunksIndexPage) GetKind() types.ListKind   { return cip.Kind }
func (cip *ChunksIndexPage) GetTotalItems() int        { return cip.TotalItems }
func (cip *ChunksIndexPage) GetExpectedTotal() int     { return cip.ExpectedTotal }
func (cip *ChunksIndexPage) GetIsFetching() bool       { return cip.IsFetching }
func (cip *ChunksIndexPage) GetState() types.LoadState { return cip.State }

func (cbp *ChunksBloomsPage) GetKind() types.ListKind   { return cbp.Kind }
func (cbp *ChunksBloomsPage) GetTotalItems() int        { return cbp.TotalItems }
func (cbp *ChunksBloomsPage) GetExpectedTotal() int     { return cbp.ExpectedTotal }
func (cbp *ChunksBloomsPage) GetIsFetching() bool       { return cbp.IsFetching }
func (cbp *ChunksBloomsPage) GetState() types.LoadState { return cbp.State }

func (cmp *ChunksManifestPage) GetKind() types.ListKind   { return cmp.Kind }
func (cmp *ChunksManifestPage) GetTotalItems() int        { return cmp.TotalItems }
func (cmp *ChunksManifestPage) GetExpectedTotal() int     { return cmp.ExpectedTotal }
func (cmp *ChunksManifestPage) GetIsFetching() bool       { return cmp.IsFetching }
func (cmp *ChunksManifestPage) GetState() types.LoadState { return cmp.State }

// ChunksPage is a union type that can represent any chunks page type
type ChunksPage struct {
	Kind           types.ListKind             `json:"kind"`
	ChunksStats    []*coreTypes.ChunkStats    `json:"chunksStats,omitempty"`
	ChunksIndex    []*coreTypes.ChunkIndex    `json:"chunksIndex,omitempty"`
	ChunksBlooms   []*coreTypes.ChunkBloom    `json:"chunksBlooms,omitempty"`
	ChunksManifest []*coreTypes.ChunkManifest `json:"chunksManifest,omitempty"`
	TotalItems     int                        `json:"totalItems"`
	ExpectedTotal  int                        `json:"expectedTotal"`
	IsFetching     bool                       `json:"isFetching"`
	State          types.LoadState            `json:"state"`
}

func (cp *ChunksPage) GetKind() types.ListKind {
	return cp.Kind
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
}

var chunksInstance *ChunksCollection
var chunksOnce sync.Once

// GetChunksCollection returns singleton instance
func GetChunksCollection() *ChunksCollection {
	chunksOnce.Do(func() {
		chunksInstance = &ChunksCollection{}
		chunksInstance.initializeFacets()
	})
	return chunksInstance
}

// NewChunksCollection constructor function to initialize the chunks collection with all its facets
func NewChunksCollection() *ChunksCollection {
	statsStore := GetChunksStatsStore()
	indexStore := GetChunksIndexStore()
	bloomsStore := GetChunksBloomsStore()
	manifestStore := GetChunksManifestStore()

	statsFacet := facets.NewFacet(
		ChunksStats,
		nil,
		nil,
		statsStore,
	)

	indexFacet := facets.NewFacet(
		ChunksIndex,
		nil,
		nil,
		indexStore,
	)

	bloomsFacet := facets.NewFacet(
		ChunksBlooms,
		nil,
		nil,
		bloomsStore,
	)

	manifestFacet := facets.NewFacet(
		ChunksManifest,
		nil,
		nil,
		manifestStore,
	)

	return &ChunksCollection{
		statsFacet:    statsFacet,
		indexFacet:    indexFacet,
		bloomsFacet:   bloomsFacet,
		manifestFacet: manifestFacet,
	}
}

// Initialize facets with appropriate filters
func (cc *ChunksCollection) initializeFacets() {
	// Stats facet - shows chunk statistics
	cc.statsFacet = facets.NewFacet(
		ChunksStats,
		func(stat *coreTypes.ChunkStats) bool {
			return true // Show all stats
		},
		nil, // No deduplication needed
		GetChunksStatsStore(),
	)

	// Index facet - shows index chunk information
	cc.indexFacet = facets.NewFacet(
		ChunksIndex,
		func(index *coreTypes.ChunkIndex) bool {
			return true // Show all index chunks
		},
		nil, // No deduplication needed
		GetChunksIndexStore(),
	)

	// Blooms facet - shows bloom filter information
	cc.bloomsFacet = facets.NewFacet(
		ChunksBlooms,
		func(bloom *coreTypes.ChunkBloom) bool {
			return true // Show all bloom filters
		},
		nil, // No deduplication needed
		GetChunksBloomsStore(),
	)

	// Manifest facet - shows manifest information
	cc.manifestFacet = facets.NewFacet(
		ChunksManifest,
		func(manifest *coreTypes.ChunkManifest) bool {
			return true // Show all manifest entries
		},
		nil, // No deduplication needed
		GetChunksManifestStore(),
	)
}

// GetPage implements the Collection interface
func (cc *ChunksCollection) GetPage(
	listKind types.ListKind,
	first, pageSize int,
	sortSpec sdk.SortSpec,
	filter string,
) (types.Page, error) {
	page := &ChunksPage{
		Kind: listKind,
	}
	filter = strings.ToLower(filter)

	switch listKind {
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
			return nil, types.NewStoreError("chunks", listKind, "GetPage", err)
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
			return nil, types.NewStoreError("chunks", listKind, "GetPage", err)
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
			return nil, types.NewStoreError("chunks", listKind, "GetPage", err)
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
			return nil, types.NewStoreError("chunks", listKind, "GetPage", err)
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
		// This is truly a validation error - invalid ListKind for this collection
		return nil, types.NewValidationError("chunks", listKind, "GetPage",
			fmt.Errorf("unsupported list kind: %v", listKind))
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
func (cc *ChunksCollection) LoadData(listKind types.ListKind) {
	switch listKind {
	case ChunksStats:
		_, _ = cc.statsFacet.Load()
	case ChunksIndex:
		_, _ = cc.indexFacet.Load()
	case ChunksBlooms:
		_, _ = cc.bloomsFacet.Load()
	case ChunksManifest:
		_, _ = cc.manifestFacet.Load()
	}
}

func (cc *ChunksCollection) Reset(listKind types.ListKind) {
	switch listKind {
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

func (cc *ChunksCollection) NeedsUpdate(listKind types.ListKind) bool {
	switch listKind {
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

func (cc *ChunksCollection) GetSupportedKinds() []types.ListKind {
	return []types.ListKind{ChunksStats, ChunksIndex, ChunksBlooms, ChunksManifest}
}

func (cc *ChunksCollection) GetStoreForKind(kind types.ListKind) string {
	switch kind {
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
func (cc *ChunksCollection) Crud(kind types.ListKind, op crud.Operation, item interface{}) error {
	// All CRUD operations are no-ops for chunks since the data is immutable blockchain data
	return nil
}

// CHUNKS_ROUTE
