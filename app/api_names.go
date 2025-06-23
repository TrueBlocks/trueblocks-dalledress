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
	collection := names.GetNamesCollection()
	return getCollectionPage[*names.NamesPage](collection, payload, first, pageSize, sort, filter)
}

func (a *App) NamesCrud(
	payload types.Payload,
	op crud.Operation,
	item interface{},
) error {
	collection := names.GetNamesCollection()
	return collection.Crud(payload, op, item)
}

func (a *App) GetNamesSummary() types.Summary {
	collection := names.GetNamesCollection()
	return collection.GetSummary()
}

// EXISTING_CODE
func (a *App) NameFromAddress(address string) (*names.Name, bool) {
	return names.GetNamesCollection().NameFromAddress(base.HexToAddress(address))
}

// EXISTING_CODE
