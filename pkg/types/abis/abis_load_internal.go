// ADD_ROUTE
package abis

import (
	"fmt"
	"sync/atomic"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/streaming"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

// ----------------------------------------------------------------
// loadInternal is the core logic for loading ABIs, functions, and events. It runs as a
// goroutine and updates the AbisCollection state as it processes the data reporting its
// progress through the msgs system.
func (ac *AbisCollection) loadInternal(listKind types.ListKind) {
	defer func() {
		atomic.StoreInt32(&ac.isLoading, 0)
	}()

	ac.ClearCache(listKind)

	var finalStatus string
	var finalPayload types.DataLoadedPayload

	switch listKind {
	case AbisDownloaded:
		finalStatus, finalPayload = ac.loadDownloadedAbis()
	case AbisKnown:
		finalStatus, finalPayload = ac.loadKnownAbis()
	case AbisFunctions:
		finalStatus, finalPayload = ac.loadFunctions()
	case AbisEvents:
		finalStatus, finalPayload = ac.loadEvents()

	default:
		logger.Error(fmt.Sprintf("AbisCollection.loadInternal: unexpected list kind: %v", listKind))
		return
	}

	msgs.EmitStatus(finalStatus)
	ac.App.EmitEvent(msgs.EventDataLoaded, finalPayload)
}

// ----------------------------------------------------------------
// loadDownloadedAbis loads ABI that are downloaded (and not known) asynchronously and returns status, payload
func (ac *AbisCollection) loadDownloadedAbis() (string, types.DataLoadedPayload) {
	contextKey := "abis-load-internal-downloaded"

	queryFunc := func(renderCtx *output.RenderCtx) {
		listOpts := sdk.AbisOptions{
			Globals:   sdk.Globals{Cache: true, Verbose: true},
			RenderCtx: renderCtx,
		}
		listOpts.AbisList()
	}

	filterFunc := func(item *coreTypes.Abi) bool {
		return !item.IsKnown // Filter for downloaded (not known) ABIs
	}

	processItemFunc := func(itemIntf interface{}) *coreTypes.Abi {
		itemPtr, ok := itemIntf.(*coreTypes.Abi)
		if !ok {
			return nil
		}
		return itemPtr
	}

	finalStatus, finalPayload, _ := streaming.LoadStreamingData(
		ac.App,
		contextKey,
		queryFunc,
		filterFunc,
		processItemFunc,
		&ac.downloadedAbis,
		&ac.expectedDownloaded,
		&ac.isDownloadedLoaded,
		"functions-events",
		&ac.mutex,
	)

	return finalStatus, finalPayload
}

// ----------------------------------------------------------------
func (ac *AbisCollection) loadKnownAbis() (string, types.DataLoadedPayload) {
	contextKey := "abis-load-internal-known"

	queryFunc := func(renderCtx *output.RenderCtx) {
		listOpts := sdk.AbisOptions{
			Globals:   sdk.Globals{Cache: true, Verbose: true},
			Known:     true,
			RenderCtx: renderCtx,
		}
		listOpts.AbisList()
	}

	filterFunc := func(item *coreTypes.Abi) bool {
		return item.IsKnown
	}

	processItemFunc := func(itemIntf interface{}) *coreTypes.Abi {
		itemPtr, ok := itemIntf.(*coreTypes.Abi)
		if !ok {
			return nil
		}
		return itemPtr
	}

	finalStatus, finalPayload, _ := streaming.LoadStreamingData(
		ac.App,
		contextKey,
		queryFunc,
		filterFunc,
		processItemFunc,
		&ac.knownAbis,
		&ac.expectedKnown,
		&ac.isKnownLoaded,
		"functions-events",
		&ac.mutex,
	)

	return finalStatus, finalPayload
}

// ----------------------------------------------------------------
// loadFunctions loads ABI functions asynchronously and returns status, payload, and error
func (ac *AbisCollection) loadFunctions() (string, types.DataLoadedPayload) {
	contextKey := "abis-load-internal-functions"

	queryFunc := func(renderCtx *output.RenderCtx) {
		detailOpts := sdk.AbisOptions{
			Globals:   sdk.Globals{Cache: true},
			RenderCtx: renderCtx,
		}
		detailOpts.AbisDetails()
	}

	filterFunc := func(item *coreTypes.Function) bool {
		if item.FunctionType == "event" {
			return false
		}

		ac.mutex.Lock()
		defer ac.mutex.Unlock()
		if _, exists := ac.deduper[item.Encoding]; exists {
			return false
		}
		ac.deduper[item.Encoding] = struct{}{}
		return true
	}

	processItemFunc := func(itemIntf interface{}) *coreTypes.Function {
		itemPtr, ok := itemIntf.(*coreTypes.Function)
		if !ok {
			return nil
		}
		return itemPtr
	}

	finalStatus, finalPayload, _ := streaming.LoadStreamingData(
		ac.App,
		contextKey,
		queryFunc,
		filterFunc,
		processItemFunc,
		&ac.allFunctions,
		&ac.expectedFunctions,
		&ac.isFuncsLoaded,
		"functions-events",
		&ac.mutex,
	)

	return finalStatus, finalPayload
}

// ----------------------------------------------------------------
// loadEvents loads ABI events asynchronously and returns status, payload
func (ac *AbisCollection) loadEvents() (string, types.DataLoadedPayload) {
	contextKey := "abis-load-internal-events"

	queryFunc := func(renderCtx *output.RenderCtx) {
		detailOpts := sdk.AbisOptions{
			Globals:   sdk.Globals{Cache: true},
			RenderCtx: renderCtx,
		}
		detailOpts.AbisDetails()
	}

	filterFunc := func(item *coreTypes.Function) bool {
		// First check if it's an event (not a function)
		if item.FunctionType != "event" {
			return false
		}

		// Then check deduper for events only
		ac.mutex.Lock()
		defer ac.mutex.Unlock()
		if _, exists := ac.deduper[item.Encoding]; exists {
			return false
		}
		ac.deduper[item.Encoding] = struct{}{}
		return true
	}

	processItemFunc := func(itemIntf interface{}) *coreTypes.Function {
		itemPtr, ok := itemIntf.(*coreTypes.Function)
		if !ok {
			return nil
		}
		return itemPtr
	}

	finalStatus, finalPayload, _ := streaming.LoadStreamingData(
		ac.App,
		contextKey,
		queryFunc,
		filterFunc,
		processItemFunc,
		&ac.allEvents,
		&ac.expectedEvents,
		&ac.isEventsLoaded,
		"functions-events",
		&ac.mutex,
	)

	return finalStatus, finalPayload
}

// ADD_ROUTE
