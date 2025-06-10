package enhancedcollection

import (
	"fmt"
	"sync"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/enhancedstore"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/preferences"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

// Singleton instances of the shared stores
var (
	enhancedAbisListStore   *enhancedstore.Store[coreTypes.Abi]
	enhancedAbisDetailStore *enhancedstore.Store[coreTypes.Function]
	listStoreMu             sync.Mutex
	detailStoreMu           sync.Mutex
)

// GetEnhancedListStore returns a singleton instance of the enhanced abis list store
func GetEnhancedListStore() *enhancedstore.Store[coreTypes.Abi] {
	listStoreMu.Lock()
	defer listStoreMu.Unlock()

	if enhancedAbisListStore == nil {
		logger.Info("Creating new enhanced ABI list store")
		enhancedAbisListStore = enhancedstore.NewStore(
			"abis-list",
			// Query function
			func(ctx *output.RenderCtx, args ...interface{}) error { // Signature updated
				chainName := preferences.GetChain()
				listOpts := sdk.AbisOptions{
					Globals:   sdk.Globals{Cache: true, Verbose: true, Chain: chainName},
					RenderCtx: ctx,
				}
				// AbisList now streams data via ctx.ModelChan and returns an error
				_, _, err := listOpts.AbisList()
				if err != nil {
					logger.Error(fmt.Sprintf("Enhanced AbisList source query error: %v", err))
					return err
				}
				return nil
			},
			// Process function
			func(itemIntf interface{}) *coreTypes.Abi {
				if abi, ok := itemIntf.(*coreTypes.Abi); ok {
					return abi
				}
				return nil
			},
		)
	}
	return enhancedAbisListStore
}

// GetEnhancedDetailStore returns a singleton instance of the enhanced abis detail store
func GetEnhancedDetailStore() *enhancedstore.Store[coreTypes.Function] {
	detailStoreMu.Lock()
	defer detailStoreMu.Unlock()

	if enhancedAbisDetailStore == nil {
		logger.Info("Creating new enhanced ABI detail store")
		enhancedAbisDetailStore = enhancedstore.NewStore(
			"abis-detail",
			// Query function
			func(ctx *output.RenderCtx, args ...interface{}) error { // Signature updated
				var addr string

				if len(args) > 0 && args[0] != nil {
					addrStr, ok := args[0].(string)
					if !ok {
						// Log the type for debugging
						logger.Warn(fmt.Sprintf("abisDetailStore.queryFunc: args[0] is non-nil but not a string. Value: %v, Type: %T", args[0], args[0]))
						return fmt.Errorf("address not provided or not a string for AbisDetails query (received type %T)", args[0])
					}
					if addrStr != "" {
						addr = addrStr
					}
				}

				var chainName string // Placeholder: This needs to be correctly sourced

				detailOpts := sdk.AbisOptions{
					Globals:   sdk.Globals{Cache: true, Chain: chainName},
					RenderCtx: ctx,
				}

				// Only set Addrs if a specific address is provided.
				if addr != "" {
					detailOpts.Addrs = []string{addr}
				} else {
					logger.Info("No specific address provided to AbisDetails query, attempting to fetch all functions/events.")
				}

				// AbisDetails now streams data via ctx.ModelChan and returns an error
				_, _, err := detailOpts.AbisDetails()
				if err != nil {
					logger.Error(fmt.Sprintf("Enhanced AbisDetails source query error: %v", err))
					return err
				}
				return nil
			},
			// Process function
			func(itemIntf interface{}) *coreTypes.Function {
				if fn, ok := itemIntf.(*coreTypes.Function); ok {
					return fn
				}
				return nil
			},
		)
	}
	return enhancedAbisDetailStore
}
