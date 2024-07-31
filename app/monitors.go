package app

import (
	"github.com/TrueBlocks/trueblocks-core/sdk/v3"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

func (a *App) GetMonitorsPage(first, pageSize int) []types.MonitorEx {
	if len(a.monitors) == 0 {
		return a.monitors
	}

	first = base.Max(0, base.Min(first, len(a.monitors)-1))
	last := base.Min(len(a.monitors), first+pageSize)
	return a.monitors[first:last]
}

func (a *App) GetMonitorsCnt() int {
	return len(a.monitors)
}

func (a *App) loadMonitors() error {
	opts := sdk.MonitorsOptions{}
	if monitors, _, err := opts.MonitorsList(); err != nil {
		return err
	} else {
		for _, monitor := range monitors {
			monitorEx := types.NewMonitorEx(a.namesMap, &monitor)
			a.monitors = append(a.monitors, monitorEx)
			a.monitorsMap[monitorEx.Address] = monitorEx
		}
	}
	return nil
}
