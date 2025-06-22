package app

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/names"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

func (a *App) GetNamesPage(
	dataFacet types.DataFacet,
	first, pageSize int,
	sortSpec sdk.SortSpec,
	filter string,
) (*names.NamesPage, error) {
	return getCollectionPage[*names.NamesPage](a.names, dataFacet, first, pageSize, sortSpec, filter)
}

func (a *App) NamesCrud(
	dataFacet types.DataFacet,
	op crud.Operation,
	name *names.Name,
	address string,
) error {
	return a.names.Crud(dataFacet, op, name)
}

func (a *App) GetNamesSummary() types.Summary {
	return a.names.GetSummary()
}

// EXISTING_CODE
func (a *App) NameFromAddress(address string) (*names.Name, bool) {
	return a.names.NameFromAddress(base.HexToAddress(address))
}

// EXISTING_CODE
