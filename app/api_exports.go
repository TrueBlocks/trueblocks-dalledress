package app

import (
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/exports"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

func (a *App) GetExportsPage(
	payload *types.Payload,
	first, pageSize int,
	sort sdk.SortSpec,
	filter string,
) (*exports.ExportsPage, error) {
	collection := exports.GetExportsCollection(payload)
	return getCollectionPage[*exports.ExportsPage](collection, payload, first, pageSize, sort, filter)
}

func (a *App) GetExportsSummary(payload *types.Payload) types.Summary {
	collection := exports.GetExportsCollection(payload)
	return collection.GetSummary()
}

func (a *App) ReloadExports(payload *types.Payload) error {
	collection := exports.GetExportsCollection(payload)
	collection.Reset(payload.DataFacet)
	collection.LoadData(payload.DataFacet)
	return nil
}
