// MONITORS_ROUTE
package monitors

import (
	"fmt"
	"sync"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/preferences"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/store"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

var (
	monitorsStore *store.Store[coreTypes.Monitor]
	storeMu       sync.Mutex
)

func GetMonitorsStore() *store.Store[coreTypes.Monitor] {
	storeMu.Lock()
	defer storeMu.Unlock()

	if monitorsStore == nil {
		monitorsStore = store.NewStore(
			GetMonitorsStoreName(MonitorsList),
			func(ctx *output.RenderCtx) error {
				chainName := preferences.GetChain()
				listOpts := sdk.MonitorsOptions{
					Globals:   sdk.Globals{Cache: true, Verbose: true, Chain: chainName},
					RenderCtx: ctx,
				}
				if _, _, err := listOpts.MonitorsList(); err != nil {
					logger.Error(fmt.Sprintf("Shared MonitorsList source query error: %v", err))
					return err
				}
				return nil
			},
			func(itemIntf interface{}) *coreTypes.Monitor {
				if monitor, ok := itemIntf.(*coreTypes.Monitor); ok {
					return monitor
				}
				return nil
			},
		)
	}
	return monitorsStore
}

func GetMonitorsStoreName(listKind types.ListKind) string {
	switch listKind {
	case MonitorsList:
		return "monitors-list"
	default:
		return ""
	}
}

func GetMonitorsCount() (int, error) {
	chainName := preferences.GetChain()
	countOpts := sdk.MonitorsOptions{
		Globals: sdk.Globals{Cache: true, Chain: chainName},
	}
	if countResult, _, err := countOpts.MonitorsCount(); err != nil {
		return 0, fmt.Errorf("MonitorsCount query error: %v", err)
	} else if len(countResult) > 0 {
		return int(countResult[0].Count), nil
	}
	return 0, nil
}

// MONITORS_ROUTE
