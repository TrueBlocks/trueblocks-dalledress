package app

import (
	// EXISTING_CODE
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/monitors"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
	// EXISTING_CODE
)

func (a *App) GetMonitorsPage(
	dataFacet types.DataFacet,
	first, pageSize int,
	sort sdk.SortSpec,
	filter string,
) (*monitors.MonitorsPage, error) {
	return getCollectionPage[*monitors.MonitorsPage](monitors.GetMonitorsCollection(), dataFacet, first, pageSize, sort, filter)
}

func (a *App) MonitorsCrud(
	dataFacet types.DataFacet,
	op crud.Operation,
	monitor *monitors.Monitor,
	address string,
) error {
	if address != "" && (monitor == nil || monitor.Address.IsZero()) {
		monitor = &monitors.Monitor{Address: base.HexToAddress(address)}
	}
	return monitors.GetMonitorsCollection().Crud(dataFacet, op, monitor)
}

func (a *App) GetMonitorsSummary() types.Summary {
	return monitors.GetMonitorsCollection().GetSummary()
}

// EXISTING_CODE
func (a *App) MonitorsClean(addresses []string) error {
	return a.monitors.Clean(addresses)
}

// EXISTING_CODE
