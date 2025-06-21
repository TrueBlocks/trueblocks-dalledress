package abis

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

type Abi = coreTypes.Abi
type Function = coreTypes.Function

var (
	abisListStore *store.Store[Abi]
	listStoreMu   sync.Mutex

	abisDetailStore *store.Store[Function]
	detailStoreMu   sync.Mutex
)

func GetAbisListStore() *store.Store[Abi] {
	listStoreMu.Lock()
	defer listStoreMu.Unlock()

	if abisListStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			chainName := preferences.GetChain()
			listOpts := sdk.AbisOptions{
				Globals:   sdk.Globals{Cache: true, Verbose: true, Chain: chainName},
				RenderCtx: ctx,
			}
			if _, _, err := listOpts.AbisList(); err != nil {
				// Create structured error with proper context
				wrappedErr := types.NewSDKError("abis", AbisDownloaded, "fetch", err)
				logger.Error(fmt.Sprintf("Abis SDK query error: %v", wrappedErr))
				return wrappedErr
			}
			return nil
		}

		processFunc := func(itemIntf interface{}) *Abi {
			if abi, ok := itemIntf.(*Abi); ok {
				return abi
			}
			return nil
		}

		storeName := GetStoreName(AbisDownloaded)
		abisListStore = store.NewStore(storeName, queryFunc, processFunc, nil)
	}

	return abisListStore
}

func GetAbisDetailStore() *store.Store[Function] {
	detailStoreMu.Lock()
	defer detailStoreMu.Unlock()

	if abisDetailStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			chainName := preferences.GetChain()
			detailOpts := sdk.AbisOptions{
				Globals:   sdk.Globals{Cache: true, Chain: chainName},
				RenderCtx: ctx,
			}
			if _, _, err := detailOpts.AbisDetails(); err != nil {
				// Create structured error with proper context
				wrappedErr := types.NewSDKError("abis", AbisFunctions, "fetch", err)
				logger.Error(fmt.Sprintf("Abis detail SDK query error: %v", wrappedErr))
				return wrappedErr
			}
			return nil
		}

		processFunc := func(itemIntf interface{}) *Function {
			if fn, ok := itemIntf.(*Function); ok {
				return fn
			}
			return nil
		}

		storeName := GetStoreName(AbisFunctions)
		abisDetailStore = store.NewStore(storeName, queryFunc, processFunc, nil)
	}

	return abisDetailStore
}

func GetStoreName(dataFacet types.DataFacet) string {
	switch dataFacet {
	case AbisDownloaded:
		fallthrough
	case AbisKnown:
		return "abis-list"
	case AbisFunctions:
		fallthrough
	case AbisEvents:
		return "abis-detail"
	default:
		return ""
	}
}
