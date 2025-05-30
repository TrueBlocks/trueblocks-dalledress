// ADD_ROUTE
package app

import (
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/sorting"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types" // Import your new types package
)

// We need these functions so that the App's abis records get handled

func (a *App) GetAbisPage(
	kind types.ListKind,
	first, pageSize int,
	sort *sorting.SortDef,
	filter string,
) (types.AbisPage, error) {
	return a.abis.GetPage(kind, first, pageSize, sort, filter)
}

func (a *App) DeleteAbi(address string) error {
	return a.abis.Delete(address)
}

// ADD_ROUTE
