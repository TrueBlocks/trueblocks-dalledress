package types

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

type Page interface {
	GetFacet() DataFacet
	GetTotalItems() int
	GetExpectedTotal() int
	GetIsFetching() bool
	GetState() LoadState
}

type Collection interface {
	GetPage(facet DataFacet, first, pageSize int, sort sdk.SortSpec, filter string) (Page, error)
	LoadData(facet DataFacet)
	Reset(facet DataFacet)
	NeedsUpdate(facet DataFacet) bool
	Crud(facet DataFacet, op crud.Operation, item interface{}, itemStr string) error
	GetSupportedFacets() []DataFacet
	GetStoreName(facet DataFacet) string
	GetSummary() Summary
	SummaryAccumulator
}
