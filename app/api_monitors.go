package app

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/monitors"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

func (a *App) GetMonitorsPage(
	payload *types.Payload,
	first, pageSize int,
	sort sdk.SortSpec,
	filter string,
) (*monitors.MonitorsPage, error) {
	collection := monitors.GetMonitorsCollection(payload)
	return getCollectionPage[*monitors.MonitorsPage](collection, payload, first, pageSize, sort, filter)
}

func (a *App) MonitorsCrud(
	payload *types.Payload,
	op crud.Operation,
	item *monitors.Monitor,
) error {
	collection := monitors.GetMonitorsCollection(payload)
	return collection.Crud(payload, op, item)
}

func (a *App) GetMonitorsSummary(payload *types.Payload) types.Summary {
	collection := monitors.GetMonitorsCollection(payload)
	return collection.GetSummary()
}

func (a *App) ReloadMonitors(payload *types.Payload) error {
	collection := monitors.GetMonitorsCollection(payload)
	collection.Reset(payload.DataFacet)
	collection.LoadData(payload.DataFacet)
	return nil
}

func (a *App) MonitorsClean(payload *types.Payload, addresses []string) error {
	collection := monitors.GetMonitorsCollection(payload)
	return collection.Clean(payload, addresses)
}
