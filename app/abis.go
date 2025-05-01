package app

import (
	"fmt"
	"sync"

	"github.com/TrueBlocks/trueblocks-core/sdk/v3"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

// Find: NewViews
func (a *App) GetAbis(first, pageSize int) types.SummaryAbis {
	first = max(0, min(first, len(a.abis.Files)-1))
	last := min(len(a.abis.Files), first+pageSize)
	copy := a.abis.ShallowCopy()
	copy.Files = a.abis.Files[first:last]
	return copy
}

func (a *App) GetAbisCnt() int {
	return len(a.abis.Files)
}

func (a *App) loadAbis(wg *sync.WaitGroup) error {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()

	opts := sdk.AbisOptions{
		Globals: sdk.Globals{
			Verbose: true,
		},
	}
	if abis, _, err := opts.AbisList(); err != nil {
		return err
	} else if (abis == nil) || (len(abis) == 0) {
		return fmt.Errorf("no status found")
	} else {
		if len(a.abis.Files) == len(abis) {
			return nil
		}
		a.abis = types.SummaryAbis{}
		a.abis.Files = append(a.abis.Files, abis...)
		a.abis.Summarize()
	}
	return nil
}
