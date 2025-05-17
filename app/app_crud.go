package app

import (
	"sync"
	"sync/atomic"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

var namesLock atomic.Int32
var namesMutex sync.Mutex

type ModifyData struct {
	Operation string       `json:"operation"`
	Address   base.Address `json:"address"`
	Value     string       `json:"value"`
}

func (a *App) ModifyName(modData *ModifyData) error {
	if !namesLock.CompareAndSwap(0, 1) {
		return nil
	}
	defer namesLock.Store(0)

	op := modData.Operation
	newName := types.Name{
		Address:  modData.Address,
		Name:     modData.Value,
		IsCustom: true,
		Source:   "TrueBlocks Browse",
		Tags:     "99-User-Defined",
	}
	if existing, ok := a.names.Map[modData.Address]; ok {
		if existing.IsCustom {
			// We preserve the tags if it's already customized
			newName.Tags = existing.Tags
		}
	}

	cd := crud.CrudFromName(newName)
	opts := sdk.NamesOptions{
		Globals: sdk.Globals{
			Chain: "mainnet", // a.getChain(),
		},
	}
	opts.Globals.Chain = "mainnet" // namesChain

	if _, _, err := opts.ModifyName(crud.OpFromString(op), cd); err != nil {
		msgs.EmitError("ModifyName", err)
		return err
	}

	newArray := []*types.Name{}
	for _, name := range a.names.List {
		if name.Address == modData.Address {
			switch crud.OpFromString(op) {
			case crud.Update:
				*name = newName
			default:
				if name.IsCustom {
					// we can only delete if it's custom already
					switch crud.OpFromString(op) {
					case crud.Delete:
						name.Deleted = true
					case crud.Undelete:
						name.Deleted = false
					case crud.Remove:
						continue
					}
				}
			}
			namesMutex.Lock()
			a.names.Map[modData.Address] = *name
			namesMutex.Unlock()
		}
		newArray = append(newArray, name)
	}
	namesMutex.Lock()
	a.names.List = newArray
	namesMutex.Unlock()

	return nil
}
