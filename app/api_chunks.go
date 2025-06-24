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
	payload *types.Payload,
	first, pageSize int,
	sort sdk.SortSpec,
	filter string,
) (*chunks.ChunksPage, error) {
	collection := chunks.GetChunksCollection(payload)
	return getCollectionPage[*chunks.ChunksPage](collection, payload, first, pageSize, sort, filter)
}

func (a *App) ChunksCrud(
	payload *types.Payload,
	op crud.Operation,
	item interface{},
) error {
	collection := chunks.GetChunksCollection(payload)
	return collection.Crud(payload, op, item)
}

func (a *App) GetChunksSummary(payload *types.Payload) types.Summary {
	collection := chunks.GetChunksCollection(payload)
	return collection.GetSummary()
}

// EXISTING_CODE
// EXISTING_CODE
