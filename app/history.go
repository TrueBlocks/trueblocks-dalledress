package app

import (
	"strings"
	"sync"

	"github.com/TrueBlocks/trueblocks-core/sdk/v3"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
)

// TODO: This should be on the App and it should be a sync.Map because it
// TODO: has the attributes described in the library's comments.
var addrToHistoryMap = map[base.Address][]TransactionEx{}
var m = sync.Mutex{}

func (a *App) GetHistoryPage(addr string, first, pageSize int) []TransactionEx {
	address := base.HexToAddress(addr)

	m.Lock()
	_, exists := addrToHistoryMap[address]
	m.Unlock()

	if !exists {
		rCtx := a.RegisterCtx(address)
		opts := sdk.ExportOptions{
			Addrs:     []string{addr},
			RenderCtx: rCtx,
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
					txEx := NewTransactionEx(a, tx)
					m.Lock()
					addrToHistoryMap[address] = append(addrToHistoryMap[address], *txEx)
					if len(addrToHistoryMap[address])%pageSize == 0 {
						a.SendMessage(address, Progress, &ProgressMsg{
							Have: int64(len(addrToHistoryMap[address])),
							Want: nItems,
						})
					}
					m.Unlock()
				case err := <-opts.RenderCtx.ErrorChan:
					a.SendMessage(address, Error, err.Error())
				default:
					if opts.RenderCtx.WasCanceled() {
						return
					}
				}
			}
		}()

		_, _, err := opts.Export()
		if err != nil {
			a.SendMessage(address, Error, err.Error())
			return []TransactionEx{}
		}

		a.SendMessage(address, Completed, "")
	}

	m.Lock()
	defer m.Unlock()
	first = base.Max(0, base.Min(first, len(addrToHistoryMap[address])-1))
	last := base.Min(len(addrToHistoryMap[address]), first+pageSize)
	return addrToHistoryMap[address][first:last]
}

func (a *App) GetHistoryCnt(addr string) int64 {
	address := base.HexToAddress(addr)

	opts := sdk.ListOptions{
		Addrs: []string{addr},
	}
	monitors, _, err := opts.ListCount()
	if err != nil {
		a.SendMessage(address, Error, err.Error())
		return 0
	} else if len(monitors) == 0 {
		return 0
	}
	return monitors[0].NRecords
}

func (a *App) ConvertToAddress(addr string) (base.Address, bool) {
	if !strings.HasSuffix(addr, ".eth") {
		ret := base.HexToAddress(addr)
		return ret, ret != base.ZeroAddr
	}

	m.Lock()
	defer m.Unlock()
	if ensAddr, exists := a.ensMap[addr]; exists {
		return ensAddr, true
	}

	// Try to get an ENS or return the same input
	opts := sdk.NamesOptions{
		Terms: []string{addr},
	}
	if names, _, err := opts.Names(); err != nil {
		a.SendMessage(base.ZeroAddr, Error, err.Error())
		return base.ZeroAddr, false
	} else {
		if len(names) > 0 {
			a.ensMap[addr] = names[0].Address
			return names[0].Address, true
		} else {
			ret := base.HexToAddress(addr)
			return ret, ret != base.ZeroAddr
		}
	}
}
