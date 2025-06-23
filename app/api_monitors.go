package app

import (
	// EXISTING_CODE

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
	item interface{},
	itemStr string,
) error {
	return monitors.GetMonitorsCollection().Crud(dataFacet, op, item, itemStr)
}

func (a *App) GetMonitorsSummary() types.Summary {
	return monitors.GetMonitorsCollection().GetSummary()
}

// EXISTING_CODE
func (a *App) MonitorsClean(addresses []string) error {
	return monitors.GetMonitorsCollection().Clean(addresses)
}

// EXISTING_CODE
