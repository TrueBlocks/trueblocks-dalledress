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
	payload types.Payload,
	first, pageSize int,
	sort sdk.SortSpec,
	filter string,
) (*monitors.MonitorsPage, error) {
	collection := monitors.GetMonitorsCollection()
	return getCollectionPage[*monitors.MonitorsPage](collection, payload, first, pageSize, sort, filter)
}

func (a *App) MonitorsCrud(
	payload types.Payload,
	op crud.Operation,
	item interface{},
) error {
	collection := monitors.GetMonitorsCollection()
	return collection.Crud(payload, op, item)
}

func (a *App) GetMonitorsSummary() types.Summary {
	collection := monitors.GetMonitorsCollection()
	return collection.GetSummary()
}

// EXISTING_CODE
func (a *App) MonitorsClean(addresses []string) error {
	return monitors.GetMonitorsCollection().Clean(addresses)
}

// EXISTING_CODE
