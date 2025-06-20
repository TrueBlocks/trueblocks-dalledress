// ADD_ROUTE
package app

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/chunks"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

func (a *App) GetChunksPage(
	dataFacet types.DataFacet,
	first, pageSize int,
	sortSpec sdk.SortSpec,
	filter string,
) (*chunks.ChunksPage, error) {
	return getCollectionPage[*chunks.ChunksPage](a.chunks, dataFacet, first, pageSize, sortSpec, filter)
}

func (a *App) ChunksCrud(
	dataFacet types.DataFacet,
	op crud.Operation,
	item interface{},
) error {
	return a.chunks.Crud(dataFacet, op, item)
}

func (a *App) GetChunksSummary() types.Summary {
	return a.chunks.GetSummary()
}

// ADD_ROUTE
