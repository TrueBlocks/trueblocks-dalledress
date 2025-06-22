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
	dataFacet types.DataFacet,
	first, pageSize int,
	sort sdk.SortSpec,
	filter string,
) (*names.NamesPage, error) {
	return getCollectionPage[*names.NamesPage](names.GetNamesCollection(), dataFacet, first, pageSize, sort, filter)
}

func (a *App) NamesCrud(
	dataFacet types.DataFacet,
	op crud.Operation,
	name *names.Name,
	address string,
) error {
	if address != "" && (name == nil || name.Address.IsZero()) {
		name = &names.Name{Address: base.HexToAddress(address)}
	}
	return names.GetNamesCollection().Crud(dataFacet, op, name)
}

func (a *App) GetNamesSummary() types.Summary {
	return names.GetNamesCollection().GetSummary()
}

// EXISTING_CODE
func (a *App) NameFromAddress(address string) (*names.Name, bool) {
	return a.names.NameFromAddress(base.HexToAddress(address))
}

// EXISTING_CODE
