package app

import (
	"fmt"

	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/sorting"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types" // Import your new types package
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

func (a *App) GetAbisPage(
	kind string,
	first, pageSize int,
	sort *sorting.SortDef,
	filter string,
) (types.AbisPage, error) {
	page, err := a.abis.GetPage(kind, first, pageSize, sort, filter)
	if err != nil {
		a.LogBackend(fmt.Sprintf("Error getting ABI page for kind '%s': %v", kind, err))
		return types.AbisPage{}, err
	}
	return page, nil
}

func (a *App) GetOneAbi(address string) (coreTypes.Abi, error) {
	a.LogBackend(fmt.Sprintf("GetOneAbi called for address: %s (STUB)", address))
	return coreTypes.Abi{}, fmt.Errorf("GetOneAbi not implemented with AbisCollection yet")
}

func (a *App) DeleteAbi(address string) error {
	opts := sdk.AbisOptions{
		Addrs: []string{address},
		Globals: sdk.Globals{
			Decache: true,
		},
	}
	if _, _, err := opts.Abis(); err != nil {
		return err
	}

	a.LogBackend(fmt.Sprintf("Deleted ABI for address: %s", address))
	return a.abis.Delete(address)
}
