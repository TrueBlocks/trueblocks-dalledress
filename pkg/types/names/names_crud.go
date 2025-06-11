// NAMES_ROUTE
package names

import (
	"sync/atomic"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

var namesLock atomic.Int32

// Crud performs CRUD operations on names (moved from App)
func (n *NamesCollection) Crud(operation crud.Operation, nameToModify *coreTypes.Name) error {
	if !namesLock.CompareAndSwap(0, 1) {
		return nil
	}
	defer namesLock.Store(0)

	nameToModify.IsCustom = true

	cd := crud.CrudFromName(*nameToModify)
	opts := sdk.NamesOptions{
		Globals: sdk.Globals{
			Chain: "mainnet", // namesChain
		},
	}

	if _, _, err := opts.ModifyName(operation, cd); err != nil {
		msgs.EmitError("Crud", err)
		return err
	}

	*n = n.ClearCache()
	return nil
}

// NAMES_ROUTE
