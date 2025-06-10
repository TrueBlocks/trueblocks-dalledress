// ADD_ROUTE
package app

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/enhancedcollection" // Import the new enhanced collection
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"              // Import your new types package
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

func (a *App) GetAbisPage(
	kind types.ListKind,
	first, pageSize int,
	sort sdk.SortSpec,
	filter string,
) (*enhancedcollection.AbisPage, error) { // Use *AbisPage from enhancedcollection
	return a.abis.GetPage(first, pageSize, "", filter, kind)
}

func (a *App) AbisCrud(
	kind types.ListKind,
	op crud.Operation,
	abi *coreTypes.Abi,
	address string) error {
	// if address != "" && (abi == nil || abi.Address.IsZero()) {
	// 	abi = &coreTypes.Abi{Address: base.HexToAddress(address)}
	// }
	// return a.abis.ForEvery(kind, op, abi)
	return nil
}

// ADD_ROUTE
