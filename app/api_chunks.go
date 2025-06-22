package app

import (
	// EXISTING_CODE
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/chunks"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
	// EXISTING_CODE
)

func (a *App) GetChunksPage(
	dataFacet types.DataFacet,
	first, pageSize int,
	sort sdk.SortSpec,
	filter string,
) (*chunks.ChunksPage, error) {
	return getCollectionPage[*chunks.ChunksPage](chunks.GetChunksCollection(), dataFacet, first, pageSize, sort, filter)
}

func (a *App) ChunksCrud(
	dataFacet types.DataFacet,
	op crud.Operation,
	chunk *chunks.Index,
	address string,
) error {
	_ = address
	return chunks.GetChunksCollection().Crud(dataFacet, op, chunk)
}

func (a *App) GetChunksSummary() types.Summary {
	return chunks.GetChunksCollection().GetSummary()
}

// EXISTING_CODE
// EXISTING_CODE
