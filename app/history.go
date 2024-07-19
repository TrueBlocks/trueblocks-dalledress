package app

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
)

var addrToHistoryMap = map[base.Address][]types.Transaction{}

func (a *App) GetHistory(addr string, first, pageSize int) []types.Transaction {
	address := base.HexToAddress(addr)
	if address.IsZero() {
		return []types.Transaction{
			{
				TransactionIndex: 1,
				BlockNumber:      1,
				BlockHash:        base.HexToHash("0x730724cb08a6eb17bf6b3296359d261570d343ea7944a17a9d7287d77900db01"),
			},
			{
				TransactionIndex: 2,
				BlockNumber:      2,
				BlockHash:        base.HexToHash("0x730724cb08a6eb17bf6b3296359d261570d343ea7944a17a9d7287d77900db02"),
			},
			{
				TransactionIndex: 3,
				BlockNumber:      3,
				BlockHash:        base.HexToHash("0x730724cb08a6eb17bf6b3296359d261570d343ea7944a17a9d7287d77900db03"),
			},
		}
	}

	var ret []types.Transaction
	if len(addrToHistoryMap[address]) == 0 {
		return ret
	}
	first = base.Max(0, base.Min(first, len(addrToHistoryMap[address])-1))
	last := base.Min(len(addrToHistoryMap[address]), first+pageSize)
	return addrToHistoryMap[address][first:last]
}

func (a *App) GetHistoryCnt(addr string) int {
	address := base.HexToAddress(addr)
	if address.IsZero() {
		return 3
	}
	return len(addrToHistoryMap[address])
}
