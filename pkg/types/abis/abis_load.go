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

// LoadData kicks off a go routine that streams the requested data making it available
// to GetPage as soon as it becomes available (even partially).
func (ac *AbisCollection) LoadData(listKind types.ListKind) {
	if !atomic.CompareAndSwapInt32(&ac.isLoading, 0, 1) {
		return
	}

	needsUpdate := ac.NeedsUpdate(listKind)
	if needsUpdate {
		go ac.loadInternal(listKind)
	} else {
		atomic.StoreInt32(&ac.isLoading, 0)
	}
}

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
		finalStatus, finalPayload = ac.loadDownloadedAbis(listKind)
	case AbisKnown:
		finalStatus, finalPayload = ac.loadKnownAbis(listKind)
	case AbisFunctions:
		finalStatus, finalPayload = ac.loadFunctions(listKind)
	case AbisEvents:
		finalStatus, finalPayload = ac.loadEvents(listKind)

	default:
		logger.Error(fmt.Sprintf("AbisCollection.loadInternal: unexpected list kind: %v", listKind))
		return
	}

	msgs.EmitStatus(finalStatus)
	if finalPayload.ListKind == "" {
		finalPayload.ListKind = string(listKind)
	}
	ac.App.EmitEvent(msgs.EventDataLoaded, finalPayload)
}

// ----------------------------------------------------------------
// loadDownloadedAbis loads ABI that are downloaded (and not known) asynchronously and returns status, payload
func (ac *AbisCollection) loadDownloadedAbis(listKind types.ListKind) (string, types.DataLoadedPayload) {
	contextKey := "abis-load-internal-downloaded"

	queryFunc := func(renderCtx *output.RenderCtx) {
		listOpts := sdk.AbisOptions{
			Globals:   sdk.Globals{Cache: true, Verbose: true},
			RenderCtx: renderCtx,
		}
		if _, _, err := listOpts.AbisList(); err != nil {
			logger.Error(fmt.Sprintf("AbisCollection.loadDownloadedAbis: %v", err))
			return
		}
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
		listKind,
		&ac.mutex,
	)

	return finalStatus, finalPayload
}

// ----------------------------------------------------------------
func (ac *AbisCollection) loadKnownAbis(listKind types.ListKind) (string, types.DataLoadedPayload) {
	contextKey := "abis-load-internal-known"

	queryFunc := func(renderCtx *output.RenderCtx) {
		listOpts := sdk.AbisOptions{
			Globals:   sdk.Globals{Cache: true, Verbose: true},
			Known:     true,
			RenderCtx: renderCtx,
		}
		if _, _, err := listOpts.AbisList(); err != nil {
			logger.Error(fmt.Sprintf("AbisCollection.loadKnownAbis: %v", err))
			return
		}
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
		listKind,
		&ac.mutex,
	)

	return finalStatus, finalPayload
}

// ----------------------------------------------------------------
// loadFunctions loads ABI functions asynchronously and returns status, payload, and error
func (ac *AbisCollection) loadFunctions(listKind types.ListKind) (string, types.DataLoadedPayload) {
	contextKey := "abis-load-internal-functions"

	queryFunc := func(renderCtx *output.RenderCtx) {
		detailOpts := sdk.AbisOptions{
			Globals:   sdk.Globals{Cache: true},
			RenderCtx: renderCtx,
		}
		if _, _, err := detailOpts.AbisDetails(); err != nil {
			logger.Error(fmt.Sprintf("AbisCollection.loadFunctions: %v", err))
			return
		}
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
		listKind,
		&ac.mutex,
	)

	return finalStatus, finalPayload
}

// ----------------------------------------------------------------
// loadEvents loads ABI events asynchronously and returns status, payload
func (ac *AbisCollection) loadEvents(listKind types.ListKind) (string, types.DataLoadedPayload) {
	contextKey := "abis-load-internal-events"

	queryFunc := func(renderCtx *output.RenderCtx) {
		detailOpts := sdk.AbisOptions{
			Globals:   sdk.Globals{Cache: true},
			RenderCtx: renderCtx,
		}
		if _, _, err := detailOpts.AbisDetails(); err != nil {
			logger.Error(fmt.Sprintf("AbisCollection.loadEvents: %v", err))
			return
		}
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
		listKind,
		&ac.mutex,
	)

	return finalStatus, finalPayload
}

// ADD_ROUTE
