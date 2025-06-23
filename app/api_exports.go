package app

import (
	// EXISTING_CODE
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/exports"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
	// EXISTING_CODE
)

func (a *App) GetExportsPage(
	payload types.Payload,
	first, pageSize int,
	sort sdk.SortSpec,
	filter string,
) (*exports.ExportsPage, error) {
	collection := exports.GetExportsCollection(payload.Chain, payload.Address)
	return getCollectionPage[*exports.ExportsPage](collection, payload, first, pageSize, sort, filter)
}

func (a *App) ExportsCrud(
	payload types.Payload,
	op crud.Operation,
	item interface{},
) error {
	collection := exports.GetExportsCollection(payload.Chain, payload.Address)
	return collection.Crud(payload, op, item)
}

func (a *App) GetExportsSummary(payload types.Payload) types.Summary {
	collection := exports.GetExportsCollection(payload.Chain, payload.Address)
	return collection.GetSummary()
}

// EXISTING_CODE
// EXISTING_CODE
