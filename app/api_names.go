package app

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/names"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

func (a *App) GetNamesPage(
	payload *types.Payload,
	first, pageSize int,
	sort sdk.SortSpec,
	filter string,
) (*names.NamesPage, error) {
	collection := names.GetNamesCollection(payload)
	return getCollectionPage[*names.NamesPage](collection, payload, first, pageSize, sort, filter)
}

func (a *App) NamesCrud(
	payload *types.Payload,
	op crud.Operation,
	item *names.Name,
) error {
	collection := names.GetNamesCollection(payload)
	return collection.Crud(payload, op, item)
}

func (a *App) GetNamesSummary(payload *types.Payload) types.Summary {
	collection := names.GetNamesCollection(payload)
	return collection.GetSummary()
}

func (a *App) ReloadNames(payload *types.Payload) error {
	collection := names.GetNamesCollection(payload)
	collection.Reset(payload.DataFacet)
	collection.LoadData(payload.DataFacet)
	return nil
}
