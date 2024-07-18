package app

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
)

var historyMap = map[base.Address][]types.Transaction{}

func (a *App) GetHistory(addr string, first, pageSize int) []types.Transaction {
	address := base.HexToAddress(addr)
	if address.IsZero() {
		return []types.Transaction{
			{
				TransactionIndex: 1,
				BlockNumber:      1,
			},
			{
				TransactionIndex: 2,
				BlockNumber:      2,
			},
			{
				TransactionIndex: 3,
				BlockNumber:      3,
			},
		}
	}

	var ret []types.Transaction
	if len(historyMap[address]) == 0 {
		return ret
	}
	first = base.Max(0, base.Min(first, len(historyMap[address])-1))
	last := base.Min(len(historyMap[address]), first+pageSize)
	return historyMap[address][first:last]
}

func (a *App) GetHistoryCnt(addr string) int {
	address := base.HexToAddress(addr)
	if address.IsZero() {
		return 3
	}
	return len(historyMap[address])
}
