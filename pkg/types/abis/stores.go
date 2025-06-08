package abis

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/logging"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/preferences"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/store"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

var (
	abisListStore   store.Store[coreTypes.Abi]
	abisDetailStore store.Store[coreTypes.Function]
)

func GetListStore() store.Store[coreTypes.Abi] {
	if abisListStore == nil {
		logging.LogBackend("Creating new shared ABI list source")
		abisListStore = store.NewSDKStore(
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

func GetDetailStore() store.Store[coreTypes.Function] {
	if abisDetailStore == nil {
		logging.LogBackend("Creating new shared ABI list source")
		abisDetailStore = store.NewSDKStore(
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
