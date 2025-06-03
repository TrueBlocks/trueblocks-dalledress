// ADD_ROUTE
package app

import (
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types" // Import your new types package
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/abis"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

func (a *App) GetAbisPage(
	kind types.ListKind,
	first, pageSize int,
	sort sdk.SortSpec,
	filter string,
) (abis.AbisPage, error) {
	return a.abis.GetPage(kind, first, pageSize, sort, filter)
}

func (a *App) RemoveAbi(address string) error {
	return a.abis.Remove(address)
}

// ADD_ROUTE
