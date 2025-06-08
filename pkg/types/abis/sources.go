package abis

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/logging"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/preferences"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/source"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

// Global shared sources - these are the key to eliminating redundant SDK calls
var (
	sharedAbisListSource    source.Source[coreTypes.Abi]
	sharedAbisDetailsSource source.Source[coreTypes.Function]
)

// GetSharedAbisListSource returns the shared source for ABI list data
// Both Downloaded and Known facets use this SAME source
func GetSharedAbisListSource() source.Source[coreTypes.Abi] {
	if sharedAbisListSource == nil {
		logging.LogBackend("Creating new shared ABI list source")
		sharedAbisListSource = source.NewSDKSource(
			func(ctx *output.RenderCtx) error {
				chainName := preferences.GetPreferredChainName()
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
	return sharedAbisListSource
}

// GetSharedAbisDetailsSource returns the shared source for ABI details/functions data
// Both Functions and Events facets use this SAME source
func GetSharedAbisDetailsSource() source.Source[coreTypes.Function] {
	if sharedAbisDetailsSource == nil {
		logging.LogBackend("Creating new shared ABI list source")
		sharedAbisDetailsSource = source.NewSDKSource(
			func(ctx *output.RenderCtx) error {
				chainName := preferences.GetPreferredChainName()
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
	return sharedAbisDetailsSource
}
