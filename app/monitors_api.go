// ADD_ROUTE
package app

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/monitors"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

func (a *App) GetMonitorsPage(
	kind types.ListKind,
	first, pageSize int,
	sort sdk.SortSpec,
	filter string,
) (*monitors.MonitorsPage, error) {
	return a.monitors.GetPage(kind, first, pageSize, sort, filter)
}

func (a *App) MonitorsCrud(
	kind types.ListKind,
	op crud.Operation,
	monitor *coreTypes.Monitor,
	address string,
) error {
	if address != "" && (monitor == nil || monitor.Address.IsZero()) {
		monitor = &coreTypes.Monitor{Address: base.HexToAddress(address)}
	}
	return a.monitors.Crud(kind, op, monitor)
}

func (a *App) MonitorsClean(addresses []string) error {
	return a.monitors.Clean(addresses)
}

// ADD_ROUTE
