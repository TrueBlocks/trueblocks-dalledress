package abis

import (
	"fmt"
	"sync"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/preferences"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/store"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

var (
	abisListStore *store.Store[coreTypes.Abi]
	listStoreMu   sync.Mutex

	abisDetailStore *store.Store[coreTypes.Function]
	detailStoreMu   sync.Mutex
)

func GetAbisListStore() *store.Store[coreTypes.Abi] {
	listStoreMu.Lock()
	defer listStoreMu.Unlock()

	if abisListStore == nil {
		abisListStore = store.NewStore(
			// TODO: Must we hard code store names? We know it's view-query
			"abis-list",
			func(ctx *output.RenderCtx) error {
				chainName := preferences.GetChain()
				listOpts := sdk.AbisOptions{
					Globals:   sdk.Globals{Cache: true, Verbose: true, Chain: chainName},
					RenderCtx: ctx,
				}
				if _, _, err := listOpts.AbisList(); err != nil {
					logger.Error(fmt.Sprintf("Shared AbisList source query error: %v", err))
					return err
				}
				return nil
			},
			func(itemIntf interface{}) *coreTypes.Abi {
				if abi, ok := itemIntf.(*coreTypes.Abi); ok {
					return abi
				}
				return nil
			},
		)
	}
	return abisListStore
}

func GetAbisDetailStore() *store.Store[coreTypes.Function] {
	detailStoreMu.Lock()
	defer detailStoreMu.Unlock()

	if abisDetailStore == nil {
		abisDetailStore = store.NewStore(
			// TODO: Must we hard code store names? We know it's view-query
			"abis-detail",
			func(ctx *output.RenderCtx) error {
				chainName := preferences.GetChain()
				detailOpts := sdk.AbisOptions{
					Globals:   sdk.Globals{Cache: true, Chain: chainName},
					RenderCtx: ctx,
				}
				if _, _, err := detailOpts.AbisDetails(); err != nil {
					logger.Error(fmt.Sprintf("Shared AbisDetails source query error: %v", err))
					return err
				}
				return nil
			},
			func(itemIntf interface{}) *coreTypes.Function {
				if fn, ok := itemIntf.(*coreTypes.Function); ok {
					return fn
				}
				return nil
			},
		)
	}
	return abisDetailStore
}
