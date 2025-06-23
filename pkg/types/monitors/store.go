// Copyright 2016, 2025 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package monitors

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

type Monitor = sdk.Monitor

var (
	monitorsStore   *store.Store[Monitor]
	monitorsStoreMu sync.Mutex
)

func (c *MonitorsCollection) getMonitorsStore() *store.Store[Monitor] {
	monitorsStoreMu.Lock()
	defer monitorsStoreMu.Unlock()

	if monitorsStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			// EXISTING_CODE
			chainName := preferences.GetChain()
			listOpts := sdk.MonitorsOptions{
				Globals:   sdk.Globals{Cache: true, Verbose: true, Chain: chainName},
				RenderCtx: ctx,
			}
			if _, _, err := listOpts.MonitorsList(); err != nil {
				// Create structured error with proper context
				wrappedErr := types.NewSDKError("monitors", MonitorsMonitors, "fetch", err)
				logger.Error(fmt.Sprintf("Monitors SDK query error: %v", wrappedErr))
				return wrappedErr
			}
			// EXISTING_CODE
			return nil
		}

		processFunc := func(item interface{}) *Monitor {
			// EXISTING_CODE
			// EXISTING_CODE
			if it, ok := item.(*Monitor); ok {
				return it
			}
			return nil
		}

		mappingFunc := func(item *Monitor) (key interface{}, includeInMap bool) {
			// EXISTING_CODE
			// EXISTING_CODE
			return nil, false
		}

		// EXISTING_CODE
		storeName := c.GetStoreName(MonitorsMonitors)
		// EXISTING_CODE
		monitorsStore = store.NewStore(storeName, queryFunc, processFunc, mappingFunc)
	}

	return monitorsStore
}

func (c *MonitorsCollection) GetStoreName(dataFacet types.DataFacet) string {
	switch dataFacet {
	case MonitorsMonitors:
		return "monitors-monitors"
	default:
		return ""
	}
}

// TODO: THIS SHOULD BE PER STORE - SEE EXPORT COMMENTS
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

var (
	collections   = make(map[store.CollectionKey]*MonitorsCollection)
	collectionsMu sync.Mutex
)

func GetMonitorsCollection(payload types.Payload) *MonitorsCollection {
	collectionsMu.Lock()
	defer collectionsMu.Unlock()

	key := store.GetCollectionKey("", "")
	if collection, exists := collections[key]; exists {
		return collection
	}

	collection := NewMonitorsCollection()
	collections[key] = collection
	return collection
}

// EXISTING_CODE
// EXISTING_CODE
