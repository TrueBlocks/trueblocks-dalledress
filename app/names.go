package app

import "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"

func (a *App) SaveName(name *types.Name) error {
	if name == nil {
		return nil
	}

	a.Logger("From SaveName in backend: " + name.String())
	return nil
}
