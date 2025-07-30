package app

import (
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/chunks"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
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

func (a *App) GetChunksSummary(payload *types.Payload) types.Summary {
	collection := chunks.GetChunksCollection(payload)
	return collection.GetSummary()
}

func (a *App) ReloadChunks(payload *types.Payload) error {
	collection := chunks.GetChunksCollection(payload)
	collection.Reset(payload.DataFacet)
	collection.LoadData(payload.DataFacet)
	return nil
}
