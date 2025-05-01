package app

import (
	"fmt"
	"sort"

	"github.com/TrueBlocks/trueblocks-core/sdk/v3"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/version"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

func (a *App) GetStatus(first, pageSize int) types.SummaryStatus {
	first = max(0, min(first, len(a.status.Caches)-1))
	last := min(len(a.status.Caches), first+pageSize)
	copy := a.status.ShallowCopy()
	copy.Caches = a.status.Caches[first:last]
	return copy
}

func (a *App) GetStatusCnt() int {
	return len(a.status.Caches)
}

func (a *App) loadStatus() error {
	opts := sdk.StatusOptions{}
	if statusArray, _, err := opts.StatusAll(); err != nil {
		return err
	} else if (statusArray == nil) || (len(statusArray) == 0) {
		return fmt.Errorf("no status found")
	} else {
		a.status.Status = statusArray[0]
		// TODO: This is a hack. We need to get the version from the core
		a.status.Version = version.LibraryVersion
		sort.Slice(a.status.Caches, func(i, j int) bool {
			return a.status.Caches[i].SizeInBytes > a.status.Caches[j].SizeInBytes
		})
		a.status.Summarize()
	}
	return nil
}
