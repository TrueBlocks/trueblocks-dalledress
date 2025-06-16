// ADD_ROUTE
package app

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/chunks"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

func (a *App) GetChunksPage(
	listKind types.ListKind,
	first, pageSize int,
	sortSpec sdk.SortSpec,
	filter string,
) (*chunks.ChunksPage, error) {
	return getCollectionPage[*chunks.ChunksPage](a.chunks, listKind, first, pageSize, sortSpec, filter)
}

func (a *App) ChunksCrud(
	listKind types.ListKind,
	op crud.Operation,
	item interface{},
) error {
	return a.chunks.Crud(listKind, op, item)
}

// ADD_ROUTE
