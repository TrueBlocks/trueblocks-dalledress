package app

import (
	// EXISTING_CODE

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/names"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
	// EXISTING_CODE
)

func (a *App) GetNamesPage(
	payload types.Payload,
	first, pageSize int,
	sort sdk.SortSpec,
	filter string,
) (*names.NamesPage, error) {
	collection := names.GetNamesCollection(payload)
	return getCollectionPage[*names.NamesPage](collection, payload, first, pageSize, sort, filter)
}

func (a *App) NamesCrud(
	payload types.Payload,
	op crud.Operation,
	item interface{},
) error {
	collection := names.GetNamesCollection(payload)
	return collection.Crud(payload, op, item)
}

func (a *App) GetNamesSummary(payload types.Payload) types.Summary {
	collection := names.GetNamesCollection(payload)
	return collection.GetSummary()
}

// EXISTING_CODE
func (a *App) NameFromAddress(address string) (*names.Name, bool) {
	collection := names.GetNamesCollection(types.Payload{})
	return collection.NameFromAddress(base.HexToAddress(address))
}

// EXISTING_CODE
