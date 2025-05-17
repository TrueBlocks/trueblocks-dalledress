package app

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
)

func (a *App) SaveName(name *types.Name) error {
	if name == nil {
		return nil
	}

	a.Logger("From SaveName in backend: " + name.String())
	return nil
}

func (a *App) DeleteName(address string) error {
	modData := &ModifyData{
		Operation: "delete",
		Address:   base.HexToAddress(address),
		Value:     "",
	}
	if err := a.ModifyName(modData); err != nil {
		msgs.EmitError("DeleteName", err)
		return err
	}

	msgs.EmitMessage(msgs.EventRefresh, "/names")
	return nil
}

func (a *App) UndeleteName(address string) error {
	modData := &ModifyData{
		Operation: "undelete",
		Address:   base.HexToAddress(address),
		Value:     "",
	}
	if err := a.ModifyName(modData); err != nil {
		msgs.EmitError("UndeleteName", err)
		return err
	}

	msgs.EmitMessage(msgs.EventRefresh, "/names")
	return nil
}

func (a *App) RemoveName(address string) error {
	modData := &ModifyData{
		Operation: "remove",
		Address:   base.HexToAddress(address),
		Value:     "",
	}
	if err := a.ModifyName(modData); err != nil {
		msgs.EmitError("RemoveName", err)
		return err
	}

	msgs.EmitMessage(msgs.EventRefresh, "/names")
	return nil
}

func (a *App) AutonameName(address string) error {
	modData := &ModifyData{
		Operation: "autoname",
		Address:   base.HexToAddress(address),
		Value:     "",
	}
	if err := a.ModifyName(modData); err != nil {
		msgs.EmitError("AutonameName", err)
		return err
	}

	msgs.EmitMessage(msgs.EventRefresh, "/names")
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
