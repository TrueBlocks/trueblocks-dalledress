// Copyright 2016, 2026 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package contracts

import (
	"fmt"
	"sync"

	// EXISTING_CODE
	// EXISTING_CODE
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/logging"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/preferences"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/store"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

// EXISTING_CODE
type Abi = sdk.Abi
type Function = sdk.Function

// EXISTING_CODE

type Contract = sdk.Contract
type Log = sdk.Log

var (
	contractsStore   *store.Store[Contract]
	contractsStoreMu sync.Mutex

	logsStore   *store.Store[Log]
	logsStoreMu sync.Mutex
)

func (c *ContractsCollection) getContractsStore(facet types.DataFacet) *store.Store[Contract] {
	contractsStoreMu.Lock()
	defer contractsStoreMu.Unlock()

	chain := preferences.GetLastChain()
	address := preferences.GetLastAddress()
	contract := preferences.GetLastContract()
	logging.LogBackend(fmt.Sprintf("getContractsStore: %s %s %s", chain, address, contract))
	theStore := contractsStore
	if theStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			// EXISTING_CODE
			// EXISTING_CODE
			return nil
		}

		processFunc := func(item interface{}) *Contract {
			if it, ok := item.(*Contract); ok {
				return it
			}
			return nil
		}

		mappingFunc := func(item *Contract) (key interface{}, includeInMap bool) {
			// EXISTING_CODE
			// EXISTING_CODE
			return nil, false
		}

		storeName := c.GetStoreName(facet, chain, address)
		theStore = store.NewStore(storeName, queryFunc, processFunc, mappingFunc)
		contractsStore = theStore
	}

	return theStore
}

func (c *ContractsCollection) getLogsStore(facet types.DataFacet) *store.Store[Log] {
	logsStoreMu.Lock()
	defer logsStoreMu.Unlock()

	chain := preferences.GetLastChain()
	address := preferences.GetLastAddress()
	theStore := logsStore
	if theStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			// EXISTING_CODE
			// EXISTING_CODE
			return nil
		}

		processFunc := func(item interface{}) *Log {
			if it, ok := item.(*Log); ok {
				return it
			}
			return nil
		}

		mappingFunc := func(item *Log) (key interface{}, includeInMap bool) {
			// EXISTING_CODE
			// EXISTING_CODE
			return nil, false
		}

		storeName := c.GetStoreName(facet, chain, address)
		theStore = store.NewStore(storeName, queryFunc, processFunc, mappingFunc)
		logsStore = theStore
	}

	return theStore
}

func (c *ContractsCollection) GetStoreName(dataFacet types.DataFacet, chain, address string) string {
	_ = chain
	_ = address
	name := ""
	switch dataFacet {
	case ContractsDashboard:
		name = "contracts-contracts"
	case ContractsDynamic:
		name = "contracts-contracts"
	case ContractsEvents:
		name = "contracts-logs"
	default:
		return ""
	}
	return name
}

// TODO: THIS SHOULD BE PER STORE - SEE EXPORT COMMENTS
func GetContractsCount(payload *types.Payload) (int, error) {
	// chain := preferences.GetLastChain()
	// countOpts := sdk.ContractsOptions{
	// 	Globals: sdk.Globals{Cache: true, Chain: chain},
	// }
	// if countResult, _, err := countOpts.ContractsCount(); err != nil {
	// 	return 0, fmt.Errorf("ContractsCount query error: %v", err)
	// } else if len(countResult) > 0 {
	// 	return int(countResult[0].Count), nil
	// }
	return 0, nil
}

var (
	collections   = make(map[store.CollectionKey]*ContractsCollection)
	collectionsMu sync.Mutex
)

func GetContractsCollection(payload *types.Payload) *ContractsCollection {
	collectionsMu.Lock()
	defer collectionsMu.Unlock()

	pl := *payload
	pl.Address = ""

	key := store.GetCollectionKey(&pl)
	if collection, exists := collections[key]; exists {
		return collection
	}

	collection := NewContractsCollection()
	collections[key] = collection
	return collection
}

// EXISTING_CODE
// EXISTING_CODE
