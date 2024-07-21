package app

import (
	"sync"

	"github.com/TrueBlocks/trueblocks-core/sdk/v3"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var addrToHistoryMap = map[base.Address][]types.Transaction{}
var m = sync.Mutex{}

func (a *App) GetHistory(addr string, first, pageSize int) []types.Transaction {
	address := base.HexToAddress(addr)
	m.Lock()
	defer m.Unlock()

	if len(addrToHistoryMap[address]) == 0 {
		opts := sdk.ExportOptions{
			Addrs:     []string{addr},
			RenderCtx: output.NewStreamingContext(),
			Globals: sdk.Globals{
				Cache: true,
				Ether: true,
			},
		}

		go func() {
			nItems := a.GetHistoryCnt(addr)
			for {
				select {
				case model := <-opts.RenderCtx.ModelChan:
					tx, ok := model.(*types.Transaction)
					if !ok {
						continue
					}
					addrToHistoryMap[address] = append(addrToHistoryMap[address], *tx)
					if len(addrToHistoryMap[address])%pageSize == 0 {
						runtime.EventsEmit(a.ctx, "Progress", len(addrToHistoryMap[address]), nItems)
					}
				case err := <-opts.RenderCtx.ErrorChan:
					runtime.EventsEmit(a.ctx, "Error", err.Error())
				default:
					if opts.RenderCtx.WasCanceled() {
						return
					}
				}
			}
		}()

		_, _, err := opts.Export()
		if err != nil {
			runtime.EventsEmit(a.ctx, "Error", err.Error())
			return []types.Transaction{}
		}

		runtime.EventsEmit(a.ctx, "Done")
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
