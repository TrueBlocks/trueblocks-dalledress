package app

import (
	"github.com/TrueBlocks/trueblocks-core/sdk/v3"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

func (a *App) GetMonitorsPage(first, pageSize int) []types.MonitorEx {
	if len(a.CurrentDoc.Monitors) == 0 {
		return a.CurrentDoc.Monitors
	}

	first = base.Max(0, base.Min(first, len(a.CurrentDoc.Monitors)-1))
	last := base.Min(len(a.CurrentDoc.Monitors), first+pageSize)
	return a.CurrentDoc.Monitors[first:last]
}

func (a *App) GetMonitorsCnt() int {
	return len(a.CurrentDoc.Monitors)
}

func (a *App) loadMonitors() error {
	opts := sdk.MonitorsOptions{}
	if monitors, _, err := opts.MonitorsList(); err != nil {
		return err
	} else {
		for _, monitor := range monitors {
			monitorEx := types.NewMonitorEx(a.namesMap, &monitor)
			a.CurrentDoc.Monitors = append(a.CurrentDoc.Monitors, monitorEx)
			a.CurrentDoc.MonitorsMap[monitorEx.Address] = monitorEx
			a.CurrentDoc.Dirty = true
		}
	}
	return nil
}
