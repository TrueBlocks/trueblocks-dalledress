package app

import (
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/status"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

func (a *App) GetStatusPage(
	payload *types.Payload,
	first, pageSize int,
	sort sdk.SortSpec,
	filter string,
) (*status.StatusPage, error) {
	collection := status.GetStatusCollection(payload)
	return getCollectionPage[*status.StatusPage](collection, payload, first, pageSize, sort, filter)
}

func (a *App) GetStatusSummary(payload *types.Payload) types.Summary {
	collection := status.GetStatusCollection(payload)
	return collection.GetSummary()
}

func (a *App) ReloadStatus(payload *types.Payload) error {
	collection := status.GetStatusCollection(payload)
	collection.Reset(payload.DataFacet)
	collection.LoadData(payload.DataFacet)
	return nil
}
