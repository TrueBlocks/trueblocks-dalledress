// ADD_ROUTE
package app

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/names"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

func (a *App) GetNamesPage(
	listKind string,
	first, pageSize int,
	sortSpec sdk.SortSpec,
	filter string,
) (names.NamesPage, error) {
	return a.names.GetPage(listKind, first, pageSize, sortSpec, filter)
}

func (a *App) NamesCrud(op crud.Operation, name *coreTypes.Name, address string) error {
	if address != "" && name == nil || name.Address.IsZero() {
		name = &coreTypes.Name{Address: base.HexToAddress(address)}
	}

	return a.names.Crud(op, name)
}

func (a *App) CleanNames(tabName string) error {
	// opts := sdk.NamesOptions{
	// 	Globals: sdk.Globals{
	// 		Chain: "mainnet",
	// 	},
	// 	Clean: true,
	// }
	// if _, _, err := opts.CleanNames(); err != nil {
	// 	msgs.Emit Message(msgs. EventError, "CleanNames: "+err.Error())
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
	// 	msgs.Emit Message(msgs. EventError, "PublishNames: "+err.Error())
	// 	return err
	// }
	return nil
}

// ADD_ROUTE
