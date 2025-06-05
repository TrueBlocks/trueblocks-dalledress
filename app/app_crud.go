package app

import (
	"sync/atomic"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

var namesLock atomic.Int32

// TODO: Do we need this?
// var namesMutex sync.Mutex

func (a *App) ModifyName(operation string, nameToModify *coreTypes.Name) error {
	if !namesLock.CompareAndSwap(0, 1) {
		return nil
	}
	defer namesLock.Store(0)

	nameToModify.IsCustom = true

	op := crud.OpFromString(operation)
	cd := crud.CrudFromName(*nameToModify)
	opts := sdk.NamesOptions{
		Globals: sdk.Globals{
			Chain: "mainnet", // a.getChain(),
		},
	}
	opts.Globals.Chain = "mainnet" // namesChain

	if _, _, err := opts.ModifyName(op, cd); err != nil {
		msgs.EmitError("ModifyName", err)
		return err
	}

	a.names = a.names.ClearCaches()
	return nil
}
