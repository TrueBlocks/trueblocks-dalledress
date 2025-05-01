package app

import (
	"fmt"
	"sort"
	"sync"

	"github.com/TrueBlocks/trueblocks-core/sdk/v3"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/messages"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

var historyMutex sync.Mutex

func (a *App) GetHistory(addr string, first, pageSize int) types.SummaryTransaction {
	address, ok := a.ConvertToAddress(addr)
	if !ok {
		messages.Send(a.ctx, messages.Error, messages.NewErrorMsg(fmt.Errorf("Invalid address: "+addr)))
		return types.SummaryTransaction{}
	}

	historyMutex.Lock()
	_, exists := a.historyMap[address]
	historyMutex.Unlock()

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
					tx, ok := model.(*coreTypes.Transaction)
					if !ok {
						continue
					}
					txEx := tx //types.NewTransactionEx(tx)
					// if name, ok := a.names.NamesMap[tx.From]; ok {
					// 	txEx.FromName = name.Name
					// }
					// if name, ok := a.names.NamesMap[tx.To]; ok {
					// 	txEx.ToName = name.Name
					// }
					historyMutex.Lock()
					summary := a.historyMap[address]
					summary.Address = address
					summary.Name = a.names.NamesMap[address].Name
					summary.Transactions = append(summary.Transactions, *txEx)
					a.historyMap[address] = summary
					if len(a.historyMap[address].Transactions)%pageSize == 0 {
						messages.Send(a.ctx,
							messages.Progress,
							messages.NewProgressMsg(int64(len(a.historyMap[address].Transactions)), nItems, address),
						)
					}
					historyMutex.Unlock()
				case err := <-opts.RenderCtx.ErrorChan:
					messages.Send(a.ctx, messages.Error, messages.NewErrorMsg(err, address))
				default:
					if opts.RenderCtx.WasCanceled() {
						return
					}
				}
			}
		}()

		_, _, err := opts.Export()
		if err != nil {
			messages.Send(a.ctx, messages.Error, messages.NewErrorMsg(err, address))
			return types.SummaryTransaction{}
		}
		historyMutex.Lock()
		sort.Slice(a.historyMap[address].Transactions, func(i, j int) bool {
			if a.historyMap[address].Transactions[i].BlockNumber == a.historyMap[address].Transactions[j].BlockNumber {
				return a.historyMap[address].Transactions[i].TransactionIndex > a.historyMap[address].Transactions[j].TransactionIndex
			}
			return a.historyMap[address].Transactions[i].BlockNumber > a.historyMap[address].Transactions[j].BlockNumber
		})
		historyMutex.Unlock()

		messages.Send(a.ctx,
			messages.Completed,
			messages.NewProgressMsg(int64(len(a.historyMap[address].Transactions)), int64(len(a.historyMap[address].Transactions)), address),
		)
	}

	historyMutex.Lock()
	defer historyMutex.Unlock()

	first = max(0, min(first, len(a.historyMap[address].Transactions)-1))
	last := min(len(a.historyMap[address].Transactions), first+pageSize)
	sum := a.historyMap[address]
	sum.Summarize()
	copy := sum.ShallowCopy()
	copy.Balance = a.getBalance(address)
	copy.Transactions = a.historyMap[address].Transactions[first:last]
	return copy
}

func (a *App) GetHistoryCnt(addr string) int64 {
	address, ok := a.ConvertToAddress(addr)
	if !ok {
		messages.Send(a.ctx, messages.Error, messages.NewErrorMsg(fmt.Errorf("Invalid address: "+addr)))
		return 0
	}

	opts := sdk.ListOptions{
		Addrs: []string{addr},
	}
	appearances, _, err := opts.ListCount()
	if err != nil {
		messages.Send(a.ctx, messages.Error, messages.NewErrorMsg(err, address))
		return 0
	} else if len(appearances) == 0 {
		return 0
	}
	return appearances[0].NRecords
}

var bMutex sync.Mutex

func (a *App) getBalance(address base.Address) string {
	bMutex.Lock()
	_, exists := a.balanceMap[address]
	bMutex.Unlock()

	if exists {
		bMutex.Lock()
		defer bMutex.Unlock()
		return a.balanceMap[address]
	}

	opts := sdk.StateOptions{
		Addrs: []string{address.Hex()},
		Globals: sdk.Globals{
			Ether: true,
			Cache: true,
		},
	}
	if balances, _, err := opts.State(); err != nil {
		return "0"
	} else {
		bMutex.Lock()
		defer bMutex.Unlock()
		a.balanceMap[address] = balances[0].Balance.ToFloatString(18)
		return a.balanceMap[address]
	}
}
