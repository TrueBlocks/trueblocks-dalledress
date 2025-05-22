package app

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
)

func (a *App) UpdateName(nameToEdit *coreTypes.Name) error {
	if err := a.ModifyName("update", nameToEdit); err != nil {
		// types.LogBackend("UpdateName error " + err.Error())
		msgs.EmitError("UpdateName", err)
		return err
	}

	// for _, name := range a.names.List {
	// 	if name.Address == nameToEdit.Address {
	// 		types.LogBackend("UpdateName success " + name.Name + " " + name.Address.Hex())
	// 	}
	// }
	// types.LogBackend(fmt.Sprintf("UpdateName Pointer value of n: %p\n", &a.names))
	msgs.EmitMessage(msgs.EventRefresh, fmt.Sprintf("Address %s was updated", nameToEdit.Address.Hex()))
	return nil
}

func (a *App) DeleteName(address string) error {
	nameToEdit := &coreTypes.Name{Address: base.HexToAddress(address)}
	if err := a.ModifyName("delete", nameToEdit); err != nil {
		msgs.EmitError("DeleteName", err)
		return err
	}
	msgs.EmitMessage(msgs.EventRefresh, fmt.Sprintf("Address %s was deleted", nameToEdit.Address.Hex()))
	return nil
}

func (a *App) UndeleteName(address string) error {
	nameToEdit := &coreTypes.Name{Address: base.HexToAddress(address)}
	if err := a.ModifyName("undelete", nameToEdit); err != nil {
		msgs.EmitError("UndeleteName", err)
		return err
	}
	msgs.EmitMessage(msgs.EventRefresh, fmt.Sprintf("Address %s was undeleted", nameToEdit.Address.Hex()))
	return nil
}

func (a *App) RemoveName(address string) error {
	nameToEdit := &coreTypes.Name{Address: base.HexToAddress(address)}
	if err := a.ModifyName("remove", nameToEdit); err != nil {
		msgs.EmitError("RemoveName", err)
		return err
	}
	msgs.EmitMessage(msgs.EventRefresh, fmt.Sprintf("Address %s was removed", nameToEdit.Address.Hex()))
	return nil
}

func (a *App) AutonameName(address string) error {
	nameToEdit := &coreTypes.Name{Address: base.HexToAddress(address)}
	if err := a.ModifyName("autoname", nameToEdit); err != nil {
		msgs.EmitError("AutonameName", err)
		return err
	}
	msgs.EmitMessage(msgs.EventRefresh, fmt.Sprintf("Address %s was autonamed", nameToEdit.Address.Hex()))
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

	msgs.EmitMessage(msgs.EventRefresh, "/names")
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

	msgs.EmitMessage(msgs.EventRefresh, "/names")
	return nil
}
