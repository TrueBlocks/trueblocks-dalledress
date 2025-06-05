// ADD_ROUTE
package app

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types" // Import your new types package
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/abis"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

func (a *App) GetAbisPage(
	kind types.ListKind,
	first, pageSize int,
	sort sdk.SortSpec,
	filter string,
) (abis.AbisPage, error) {
	return a.abis.GetPage(kind, first, pageSize, sort, filter)
}

func (a *App) AbisCrud(op crud.Operation, abi *coreTypes.Abi, address string) error {
	if address != "" && abi == nil {
		abi = &coreTypes.Abi{Address: base.HexToAddress(address)}
	}
	return a.abis.Crud(op, abi)
}

// ADD_ROUTE
