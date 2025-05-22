package app

import (
	"sync"
	"sync/atomic"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

var namesLock atomic.Int32
var namesMutex sync.Mutex

func (a *App) ModifyName(operation string, nameToModify *coreTypes.Name) error {
	if !namesLock.CompareAndSwap(0, 1) {
		// types.LogBackend("ModifyName already in progress")
		return nil
	}
	defer namesLock.Store(0)

	// User can only modify custom names - editing an existing non-custom name makes it custom
	// types.LogBackend("ModifyName: " + operation + " " + nameToModify.Name + " " + nameToModify.Address.Hex())
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
		// types.LogBackend("ModifyName error " + err.Error())
		msgs.EmitError("ModifyName", err)
		return err
	}

	// types.LogBackend("ModifyName calling ReloadNames")
	// types.LogBackend(fmt.Sprintf("ReloadNames before reload %d", len(a.names.List)))
	// for _, name := range a.names.List {
	// if name.Address.Hex() == "0x00a3819199113fc6a6e6ba1298afde7377e2009b" {
	// 	types.LogBackend("ReloadNames found 0x00a3819199113fc6a6e6ba1298afde7377e2009b:" + name.Name)
	// }
	// }

	a.names = a.names.ReloadNames()

	// types.LogBackend(fmt.Sprintf("ReloadNames after reload %d", len(a.names.List)))
	// for _, name := range a.names.List {
	// 	if name.Address.Hex() == "0x00a3819199113fc6a6e6ba1298afde7377e2009b" {
	// 		types.LogBackend("ReloadNames found 0x00a3819199113fc6a6e6ba1298afde7377e2009b:" + name.Name)
	// 	}
	// }
	//
	// a.names.List = make([]*coreTypes.Name, 0, len(a.names.List))
	//
	// types.LogBackend(fmt.Sprintf("ReloadNames after empty List %d", len(a.names.List)))
	// for _, name := range a.names.List {
	// 	if name.Address.Hex() == "0x00a3819199113fc6a6e6ba1298afde7377e2009b" {
	// 		types.LogBackend("WTF? 0x00a3819199113fc6a6e6ba1298afde7377e2009b:" + name.Name)
	// 	}
	// }

	return nil
}
