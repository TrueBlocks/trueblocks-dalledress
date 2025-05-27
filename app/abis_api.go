// ADD_ABIS_CODE
package app

import (
	"fmt"

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
	return a.abis.GetPage(kind, first, pageSize, sort, filter)
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
