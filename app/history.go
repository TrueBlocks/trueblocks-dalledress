package app

import (
	"sync"

	"github.com/TrueBlocks/trueblocks-core/sdk/v3"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
)

type TransactionEx struct {
	BlockNumber      base.Blknum    `json:"blockNumber"`
	TransactionIndex base.Txnum     `json:"transactionIndex"`
	Timestamp        base.Timestamp `json:"timestamp"`
	Date             string         `json:"date"`
	From             base.Address   `json:"from"`
	FromName         string         `json:"fromName"`
	To               base.Address   `json:"to"`
	ToName           string         `json:"toName"`
	Wei              base.Wei       `json:"wei"`
	Ether            string         `json:"ether"`
	Function         string         `json:"function"`
	HasToken         bool           `json:"hasToken"`
	IsError          bool           `json:"isError"`
}

func NewTransactionEx(a *App, tx *types.Transaction) *TransactionEx {
	fromName := a.namesMap[tx.From].Name
	if len(fromName) == 0 {
		fromName = tx.From.String()
	} else if len(fromName) > 39 {
		fromName = fromName[:39] + "..."
	}
	toName := a.namesMap[tx.To].Name
	if len(toName) == 0 {
		toName = tx.To.String()
	} else if len(toName) > 39 {
		toName = toName[:39] + "..."
	}
	ether := tx.Value.ToEtherStr(18)
	if tx.Value.IsZero() {
		ether = "-"
	} else if len(ether) > 5 {
		ether = ether[:5]
	}
	return &TransactionEx{
		BlockNumber:      tx.BlockNumber,
		TransactionIndex: tx.TransactionIndex,
		Timestamp:        tx.Timestamp,
		Date:             tx.Date(),
		From:             tx.From,
		FromName:         fromName,
		To:               tx.To,
		ToName:           toName,
		Wei:              tx.Value,
		Ether:            ether,
		HasToken:         tx.HasToken,
		IsError:          tx.IsError,
		// Function:         tx.Function(),
	}
}

var addrToHistoryMap = map[base.Address][]TransactionEx{}
var m = sync.Mutex{}

func (a *App) GetHistory(addr string, first, pageSize int) []TransactionEx {
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
					txEx := NewTransactionEx(a, tx)
					// if strings.HasPrefix(txEx.ToName, "0x") {
					addrToHistoryMap[address] = append(addrToHistoryMap[address], *txEx)
					if len(addrToHistoryMap[address])%pageSize == 0 {
						a.SendMessage(address, Progress, &ProgressMsg{
							Have: int64(len(addrToHistoryMap[address])),
							Want: nItems,
						})
					}
					// }
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
