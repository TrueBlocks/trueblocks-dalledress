package app

import (
	"fmt"
	"sort"
	"sync"

	"github.com/TrueBlocks/trueblocks-core/sdk/v3"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

func (a *App) GetManifest(first, pageSize int) types.SummaryManifest {
	first = base.Max(0, base.Min(first, len(a.manifest.Chunks)-1))
	last := base.Min(len(a.manifest.Chunks), first+pageSize)
	copy := a.manifest.ShallowCopy()
	copy.Chunks = a.manifest.Chunks[first:last]
	return copy
}

func (a *App) GetManifestCnt() int {
	return len(a.manifest.Chunks)
}

func (a *App) loadManifest(wg *sync.WaitGroup) error {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()

	opts := sdk.ChunksOptions{}
	if manifests, _, err := opts.ChunksManifest(); err != nil {
		return err
	} else if (manifests == nil) || (len(manifests) == 0) {
		return fmt.Errorf("no manifest found")
	} else {
		if len(a.manifest.Chunks) == len(manifests[0].Chunks) {
			return nil
		}
		a.manifest = types.NewSummaryManifest(manifests[0])
		sort.Slice(a.manifest.Chunks, func(i, j int) bool {
			return a.manifest.Chunks[i].Range > a.manifest.Chunks[j].Range
		})
	}
	return nil
}
