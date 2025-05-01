package app

import (
	"fmt"
	"sort"
	"sync"

	"github.com/TrueBlocks/trueblocks-core/sdk/v3"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

// Find: NewViews
func (a *App) GetIndex(first, pageSize int) types.SummaryIndex {
	first = max(0, min(first, len(a.index.Chunks)-1))
	last := min(len(a.index.Chunks), first+pageSize)
	copy := a.index.ShallowCopy()
	copy.Chunks = a.index.Chunks[first:last]
	return copy
}

func (a *App) GetIndexCnt() int {
	return len(a.index.Chunks)
}

func (a *App) loadIndex(wg *sync.WaitGroup) error {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()

	opts := sdk.ChunksOptions{}
	if chunks, _, err := opts.ChunksStats(); err != nil {
		return err
	} else if (chunks == nil) || (len(chunks) == 0) {
		return fmt.Errorf("no index chunks found")
	} else {
		if len(a.index.Chunks) == len(chunks) {
			return nil
		}
		a.index = types.SummaryIndex{Chunks: chunks}
		sort.Slice(a.index.Chunks, func(i, j int) bool {
			// reverse order
			return a.index.Chunks[i].Range > a.index.Chunks[j].Range
		})
		a.index.Summarize()
	}
	return nil
}
