// ADD_ROUTE
package app

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/sorting"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

func (a *App) GetNamesPage(listType string, first, pageSize int, sortKey sorting.SortDef, filter string) types.NamesPage {
	return a.names.GetPage(listType, first, pageSize, sortKey, filter)
}

func (a *App) UpdateName(nameToEdit *coreTypes.Name) error {
	if err := a.ModifyName("update", nameToEdit); err != nil {
		// types.LogBackend("UpdateName error " + err.Error())
		msgs.EmitError("UpdateName", err)
		return err
	}
	return nil
}

func (a *App) DeleteName(address string) error {
	nameToEdit := &coreTypes.Name{Address: base.HexToAddress(address)}
	if err := a.ModifyName("delete", nameToEdit); err != nil {
		msgs.EmitError("DeleteName", err)
		return err
	}
	return nil
}

func (a *App) UndeleteName(address string) error {
	nameToEdit := &coreTypes.Name{Address: base.HexToAddress(address)}
	if err := a.ModifyName("undelete", nameToEdit); err != nil {
		msgs.EmitError("UndeleteName", err)
		return err
	}
	return nil
}

func (a *App) RemoveName(address string) error {
	nameToEdit := &coreTypes.Name{Address: base.HexToAddress(address)}
	if err := a.ModifyName("remove", nameToEdit); err != nil {
		msgs.EmitError("RemoveName", err)
		return err
	}
	return nil
}

func (a *App) AutonameName(address string) error {
	nameToEdit := &coreTypes.Name{Address: base.HexToAddress(address)}
	if err := a.ModifyName("autoname", nameToEdit); err != nil {
		msgs.EmitError("AutonameName", err)
		return err
	}
	return nil
}

func (a *App) CleanNames(tabName string) error {
	// opts := sdk.NamesOptions{
	// 	Globals: sdk.Globals{
	// 		Chain: "mainnet",
	// 	},
	// 	Clean: true,
	// }
	// if _, _, err := opts.CleanNames(); err != nil {
	// 	msgs.EmitError("CleanNames", err)
	// 	return err
	// }
	return nil
}

func (a *App) PublishNames(tabName string) error {
	// opts := sdk.NamesOptions{
	// 	Globals: sdk.Globals{
	// 		Chain: "mainnet",
	// 	},
	// 	Publish: true,
	// }
	// if _, _, err := opts.PublishNames(); err != nil {
	// 	msgs.EmitError("PublishNames", err)
	// 	return err
	// }
	return nil
}

// ADD_ROUTE
