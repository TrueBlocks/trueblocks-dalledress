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

type Monitor = coreTypes.Monitor

var (
	monitorsStore *store.Store[Monitor]
	monitorsMu    sync.Mutex
)

func GetMonitorsStore() *store.Store[Monitor] {
	monitorsMu.Lock()
	defer monitorsMu.Unlock()

	if monitorsStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			chainName := preferences.GetChain()
			listOpts := sdk.MonitorsOptions{
				Globals:   sdk.Globals{Cache: true, Verbose: true, Chain: chainName},
				RenderCtx: ctx,
			}
			if _, _, err := listOpts.MonitorsList(); err != nil {
				// Create structured error with proper context
				wrappedErr := types.NewSDKError("monitors", MonitorsList, "fetch", err)
				logger.Error(fmt.Sprintf("Monitors SDK query error: %v", wrappedErr))
				return wrappedErr
			}
			return nil
		}

		processFunc := func(itemIntf interface{}) *Monitor {
			if monitor, ok := itemIntf.(*Monitor); ok {
				return monitor
			}
			return nil
		}

		storeName := GetStoreName(MonitorsList)
		monitorsStore = store.NewStore(storeName, queryFunc, processFunc, nil)
	}

	return monitorsStore
}

func GetStoreName(dataFacet types.DataFacet) string {
	switch dataFacet {
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
