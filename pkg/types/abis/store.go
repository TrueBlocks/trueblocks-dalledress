// Copyright 2016, 2025 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package abis

import (
	"fmt"
	"sync"

	// EXISTING_CODE
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/preferences"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/store"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
	// EXISTING_CODE
)

type Abi = coreTypes.Abi
type Function = coreTypes.Function

var (
	abisStore   *store.Store[Abi]
	abisStoreMu sync.Mutex

	functionsStore   *store.Store[Function]
	functionsStoreMu sync.Mutex
)

func (c *AbisCollection) getAbisStore() *store.Store[Abi] {
	abisStoreMu.Lock()
	defer abisStoreMu.Unlock()

	if abisStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			// EXISTING_CODE
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
			// EXISTING_CODE
			return nil
		}

		processFunc := func(itemIntf interface{}) *Abi {
			// EXISTING_CODE
			if abi, ok := itemIntf.(*Abi); ok {
				return abi
			}
			// EXISTING_CODE
			return nil
		}

		mappingFunc := func(item *Abi) (key interface{}, includeInMap bool) {
			// EXISTING_CODE
			// EXISTING_CODE
			return nil, false
		}

		// EXISTING_CODE
		storeName := c.GetStoreName(AbisDownloaded)
		// EXISTING_CODE
		abisStore = store.NewStore(storeName, queryFunc, processFunc, mappingFunc)
	}

	return abisStore
}

func (c *AbisCollection) getFunctionsStore() *store.Store[Function] {
	functionsStoreMu.Lock()
	defer functionsStoreMu.Unlock()

	if functionsStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			// EXISTING_CODE
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
			// EXISTING_CODE
			return nil
		}

		processFunc := func(itemIntf interface{}) *Function {
			// EXISTING_CODE
			if fn, ok := itemIntf.(*Function); ok {
				return fn
			}
			// EXISTING_CODE
			return nil
		}

		mappingFunc := func(item *Function) (key interface{}, includeInMap bool) {
			// EXISTING_CODE
			// EXISTING_CODE
			return nil, false
		}

		// EXISTING_CODE
		storeName := c.GetStoreName(AbisFunctions)
		// EXISTING_CODE
		functionsStore = store.NewStore(storeName, queryFunc, processFunc, mappingFunc)
	}

	return functionsStore
}

func (c *AbisCollection) GetStoreName(dataFacet types.DataFacet) string {
	switch dataFacet {
	case AbisDownloaded:
		return "abis-abis"
	case AbisKnown:
		return "abis-abis"
	case AbisFunctions:
		return "abis-functions"
	case AbisEvents:
		return "abis-functions"
	default:
		return ""
	}
}

func GetAbisCount() (int, error) {
	chainName := preferences.GetChain()
	countOpts := sdk.AbisOptions{
		Globals: sdk.Globals{Cache: true, Chain: chainName},
	}
	if countResult, _, err := countOpts.AbisCount(); err != nil {
		return 0, fmt.Errorf("AbisCount query error: %v", err)
	} else if len(countResult) > 0 {
		return int(countResult[0].Count), nil
	}
	return 0, nil
}

// EXISTING_CODE
// EXISTING_CODE
