// Copyright 2016, 2026 The TrueBlocks Authors. All rights reserved.
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
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/preferences"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/store"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
	// EXISTING_CODE
)

// EXISTING_CODE
// EXISTING_CODE

type Abi = sdk.Abi
type Function = sdk.Function

var (
	abisStore   *store.Store[Abi]
	abisStoreMu sync.Mutex

	functionsStore   *store.Store[Function]
	functionsStoreMu sync.Mutex
)

func (c *AbisCollection) getAbisStore() *store.Store[Abi] {
	abisStoreMu.Lock()
	defer abisStoreMu.Unlock()

	theStore := abisStore
	if theStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			// EXISTING_CODE
			chain := preferences.GetLastChain()
			listOpts := sdk.AbisOptions{
				Globals:   sdk.Globals{Cache: true, Verbose: true, Chain: chain},
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

		processFunc := func(item interface{}) *Abi {
			// EXISTING_CODE
			// EXISTING_CODE
			if it, ok := item.(*Abi); ok {
				return it
			}
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
		theStore = store.NewStore(storeName, queryFunc, processFunc, mappingFunc)
		abisStore = theStore
	}

	return theStore
}

func (c *AbisCollection) getFunctionsStore() *store.Store[Function] {
	functionsStoreMu.Lock()
	defer functionsStoreMu.Unlock()

	theStore := functionsStore
	if theStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			// EXISTING_CODE
			chain := preferences.GetLastChain()
			detailOpts := sdk.AbisOptions{
				Globals:   sdk.Globals{Cache: true, Chain: chain},
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

		processFunc := func(item interface{}) *Function {
			// EXISTING_CODE
			// EXISTING_CODE
			if it, ok := item.(*Function); ok {
				return it
			}
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
		theStore = store.NewStore(storeName, queryFunc, processFunc, mappingFunc)
		functionsStore = theStore
	}

	return theStore
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

// TODO: THIS SHOULD BE PER STORE - SEE EXPORT COMMENTS
func GetAbisCount(payload *types.Payload) (int, error) {
	chain := preferences.GetLastChain()
	countOpts := sdk.AbisOptions{
		Globals: sdk.Globals{Cache: true, Chain: chain},
	}
	if countResult, _, err := countOpts.AbisCount(); err != nil {
		return 0, fmt.Errorf("AbisCount query error: %v", err)
	} else if len(countResult) > 0 {
		return int(countResult[0].Count), nil
	}
	return 0, nil
}

var (
	collections   = make(map[store.CollectionKey]*AbisCollection)
	collectionsMu sync.Mutex
)

func GetAbisCollection(payload *types.Payload) *AbisCollection {
	collectionsMu.Lock()
	defer collectionsMu.Unlock()

	pl := *payload
	pl.Address = ""

	key := store.GetCollectionKey(&pl)
	if collection, exists := collections[key]; exists {
		return collection
	}

	collection := NewAbisCollection()
	collections[key] = collection
	return collection
}

// EXISTING_CODE
// EXISTING_CODE
