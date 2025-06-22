package app

import (
	// EXISTING_CODE
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/abis"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
	// EXISTING_CODE
)

func (a *App) GetAbisPage(
	dataFacet types.DataFacet,
	first, pageSize int,
	sort sdk.SortSpec,
	filter string,
) (*abis.AbisPage, error) {
	return getCollectionPage[*abis.AbisPage](abis.GetAbisCollection(), dataFacet, first, pageSize, sort, filter)
}

func (a *App) AbisCrud(
	dataFacet types.DataFacet,
	op crud.Operation,
	abi *abis.Abi,
	address string,
) error {
	if address != "" && (abi == nil || abi.Address.IsZero()) {
		abi = &abis.Abi{Address: base.HexToAddress(address)}
	}
	return abis.GetAbisCollection().Crud(dataFacet, op, abi)
}

func (a *App) GetAbisSummary() types.Summary {
	return abis.GetAbisCollection().GetSummary()
}

// EXISTING_CODE
// EXISTING_CODE
