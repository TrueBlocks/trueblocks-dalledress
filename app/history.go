package app

import (
	"sync"

	"github.com/TrueBlocks/trueblocks-core/sdk/v3"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
)

var addrToHistoryMap = map[base.Address][]types.Transaction{}
var m = sync.Mutex{}

func (a *App) GetHistory(addr string, first, pageSize int) []types.Transaction {
	address := base.HexToAddress(addr)
	m.Lock()
	defer m.Unlock()

	if len(addrToHistoryMap[address]) == 0 {
		opts := sdk.ExportOptions{
			Addrs: []string{addr},
			Globals: sdk.Globals{
				Cache: true,
			},
		}
		monitors, _, err := opts.Export()
		if err != nil {
			// EventEmitter.Emit("error", err)
			logger.Info(err)
			return []types.Transaction{}
		} else if len(monitors) == 0 {
			logger.Info("none")
			return []types.Transaction{}
		} else {
			logger.Info("got em", len(monitors))
			addrToHistoryMap[address] = monitors
		}
	}
	first = base.Max(0, base.Min(first, len(addrToHistoryMap[address])-1))
	last := base.Min(len(addrToHistoryMap[address]), first+pageSize)
	return addrToHistoryMap[address][first:last]
}

func (a *App) GetHistoryCnt(addr string) int64 {
	opts := sdk.ListOptions{
		Addrs: []string{addr},
	}
	monitors, _, err := opts.ListCount()
	if err != nil {
		// EventEmitter.Emit("error", err)
		return 0
	} else if len(monitors) == 0 {
		return 0
	}
	return monitors[0].NRecords
}
