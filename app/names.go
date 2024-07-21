package app

import (
	"github.com/TrueBlocks/trueblocks-core/sdk/v3"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/names"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
)

func (a *App) GetNames(first, pageSize int) []types.Name {
	var ret []types.Name
	if len(a.names) == 0 {
		return ret
	}
	first = base.Max(0, base.Min(first, len(a.names)-1))
	last := base.Min(len(a.names), first+pageSize)
	return a.names[first:last]
}

func (a *App) GetNamesCnt() int {
	return len(a.names)
}

func (a *App) loadNames() error {
	var err error

	types := names.Regular | names.Custom | names.Prefund | names.Baddress
	if a.namesMap, err = names.LoadNamesMap("mainnet", types, nil); err != nil {
		return err
	}

	opts := sdk.NamesOptions{
		Regular: true,
		// Custom:  true,
		// Prefund: true,
		Globals: sdk.Globals{
			Chain: "mainnet",
		},
	}

	if a.names, _, err = opts.Names(); err != nil {
		return err
	} else {
		return nil
	}
}
