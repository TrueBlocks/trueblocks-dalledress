package types

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

type Page interface {
	GetKind() ListKind
	GetTotalItems() int
	GetExpectedTotal() int
	GetIsFetching() bool
	GetState() LoadState
}

type Collection interface {
	GetPage(kind ListKind, first, pageSize int, sort sdk.SortSpec, filter string) (Page, error)
	LoadData(kind ListKind)
	Reset(kind ListKind)
	NeedsUpdate(kind ListKind) bool
	Crud(kind ListKind, op crud.Operation, item interface{}) error
	GetSupportedKinds() []ListKind
	GetStoreForKind(kind ListKind) string
	GetCollectionName() string
	GetSummary() Summary
	SummaryAccumulator
}
