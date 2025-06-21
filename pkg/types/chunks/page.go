package chunks

import (
	"fmt"
	"strings"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

type ChunksPage struct {
	Facet         types.DataFacet `json:"facet"`
	Stats         []*Stats        `json:"stats"`
	Index         []*Index        `json:"index"`
	Blooms        []*Bloom        `json:"blooms"`
	Manifest      []*Manifest     `json:"manifest"`
	TotalItems    int             `json:"totalItems"`
	ExpectedTotal int             `json:"expectedTotal"`
	IsFetching    bool            `json:"isFetching"`
	State         types.LoadState `json:"state"`
}

func (p *ChunksPage) GetFacet() types.DataFacet {
	return p.Facet
}

func (p *ChunksPage) GetTotalItems() int {
	return p.TotalItems
}

func (p *ChunksPage) GetExpectedTotal() int {
	return p.ExpectedTotal
}

func (p *ChunksPage) GetIsFetching() bool {
	return p.IsFetching
}

func (p *ChunksPage) GetState() types.LoadState {
	return p.State
}

func (c *ChunksCollection) GetPage(
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
		var filterFunc func(*Stats) bool
		if filter != "" {
			filterFunc = func(stat *Stats) bool {
				return c.matchesStatsFilter(stat, filter)
			}
		}

		sortFunc := func(items []Stats, sort sdk.SortSpec) error {
			// Placeholder - implement actual sorting if needed
			return nil
		}

		if result, err := c.statsFacet.GetPage(first, pageSize, filterFunc, sortSpec, sortFunc); err != nil {
			// This is likely an SDK or store error, not a validation error
			return nil, types.NewStoreError("chunks", dataFacet, "GetPage", err)
		} else {
			stats := make([]*Stats, 0, len(result.Items))
			for i := range result.Items {
				stats = append(stats, &result.Items[i])
			}
			page.Stats, page.TotalItems, page.State = stats, result.TotalItems, result.State
		}
		page.IsFetching = c.statsFacet.IsFetching()
		page.ExpectedTotal = c.statsFacet.ExpectedCount()

	case ChunksIndex:
		var filterFunc func(*Index) bool
		if filter != "" {
			filterFunc = func(index *Index) bool {
				return c.matchesIndexFilter(index, filter)
			}
		}

		sortFunc := func(items []Index, sort sdk.SortSpec) error {
			// Placeholder - implement actual sorting if needed
			return nil
		}

		if result, err := c.indexFacet.GetPage(first, pageSize, filterFunc, sortSpec, sortFunc); err != nil {
			// This is likely an SDK or store error, not a validation error
			return nil, types.NewStoreError("chunks", dataFacet, "GetPage", err)
		} else {
			index := make([]*Index, 0, len(result.Items))
			for i := range result.Items {
				index = append(index, &result.Items[i])
			}
			page.Index, page.TotalItems, page.State = index, result.TotalItems, result.State
		}
		page.IsFetching = c.indexFacet.IsFetching()
		page.ExpectedTotal = c.indexFacet.ExpectedCount()

	case ChunksBlooms:
		var filterFunc func(*Bloom) bool
		if filter != "" {
			filterFunc = func(bloom *Bloom) bool {
				return c.matchesBloomsFilter(bloom, filter)
			}
		}

		sortFunc := func(items []Bloom, sort sdk.SortSpec) error {
			// Placeholder - implement actual sorting if needed
			return nil
		}

		if result, err := c.bloomsFacet.GetPage(first, pageSize, filterFunc, sortSpec, sortFunc); err != nil {
			// This is likely an SDK or store error, not a validation error
			return nil, types.NewStoreError("chunks", dataFacet, "GetPage", err)
		} else {
			blooms := make([]*Bloom, 0, len(result.Items))
			for i := range result.Items {
				blooms = append(blooms, &result.Items[i])
			}
			page.Blooms, page.TotalItems, page.State = blooms, result.TotalItems, result.State
		}
		page.IsFetching = c.bloomsFacet.IsFetching()
		page.ExpectedTotal = c.bloomsFacet.ExpectedCount()

	case ChunksManifest:
		var filterFunc func(*Manifest) bool
		if filter != "" {
			filterFunc = func(manifest *Manifest) bool {
				return c.matchesManifestFilter(manifest, filter)
			}
		}

		sortFunc := func(items []Manifest, sort sdk.SortSpec) error {
			// Placeholder - implement actual sorting if needed
			return nil
		}

		if result, err := c.manifestFacet.GetPage(first, pageSize, filterFunc, sortSpec, sortFunc); err != nil {
			// This is likely an SDK or store error, not a validation error
			return nil, types.NewStoreError("chunks", dataFacet, "GetPage", err)
		} else {
			manifest := make([]*Manifest, 0, len(result.Items))
			for i := range result.Items {
				manifest = append(manifest, &result.Items[i])
			}
			page.Manifest, page.TotalItems, page.State = manifest, result.TotalItems, result.State
		}
		page.IsFetching = c.manifestFacet.IsFetching()
		page.ExpectedTotal = c.manifestFacet.ExpectedCount()

	default:
		// This is truly a validation error - invalid DataFacet for this collection
		return nil, types.NewValidationError("chunks", dataFacet, "GetPage",
			fmt.Errorf("unsupported data facet: %v", dataFacet))
	}

	return page, nil
}

func (c *ChunksCollection) matchesStatsFilter(stat *Stats, filter string) bool {
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

func (c *ChunksCollection) matchesIndexFilter(index *Index, filter string) bool {
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

func (c *ChunksCollection) matchesBloomsFilter(bloom *Bloom, filter string) bool {
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

func (c *ChunksCollection) matchesManifestFilter(manifest *Manifest, filter string) bool {
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

// TODO: THIS ISN'T NEEDED (FOR CHUNKS OR ANY OTHER -- SEE THE APP)
// func (c *ChunksCollection) Get ChunksPage(
// 	dataFacet types.DataFacet,
// 	first, pageSize int,
// 	sortSpec sdk.SortSpec,
// 	filter string,
// ) (*ChunksPage, error) {
// 	page, err := c.GetPage(dataFacet, first, pageSize, sortSpec, filter)
// 	if err != nil {
// 		return nil, err
// 	}
// 	chunksPage, ok := page.(*ChunksPage)
// 	if !ok {
// 		return nil, fmt.Errorf("internal error: GetPage returned unexpected type %T", page)
// 	}
// 	return chunksPage, nil
// }
