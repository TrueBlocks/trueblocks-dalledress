// ADD_ROUTE
package abis

import (
	"fmt"
	"sync/atomic"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
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

	finalStatus := ""
	finalPayload := types.DataLoadedPayload{
		Category: "abis",
		DataType: "functions-events",
		IsLoaded: true,
	}

	switch listKind {
	case AbisDownloaded:
		listOpts := sdk.AbisOptions{
			Globals: sdk.Globals{Cache: true, Verbose: true},
		}
		if sdkAbis, _, err := listOpts.AbisList(); err != nil || sdkAbis == nil {
			msgs.EmitError("AbisCollection.loadInternal: error fetching ABI list", err)
			return
		} else {
			ac.mutex.Lock()
			ac.downloadedAbis = make([]coreTypes.Abi, 0, len(ac.downloadedAbis))
			for _, abi := range sdkAbis {
				if !abi.IsKnown {
					ac.downloadedAbis = append(ac.downloadedAbis, abi)
				}
			}
			ac.isDownloadedLoaded = true
			ac.mutex.Unlock()

			finalStatus = fmt.Sprintf("Loaded %d downloaded abis", len(ac.downloadedAbis))
			finalPayload.CurrentCount = len(ac.downloadedAbis)
			finalPayload.ExpectedTotal = len(ac.downloadedAbis)
		}
	case AbisKnown:
		listOpts := sdk.AbisOptions{
			Globals: sdk.Globals{Cache: true, Verbose: true},
			Known:   true,
		}
		if sdkAbis, _, err := listOpts.AbisList(); err != nil || sdkAbis == nil {
			msgs.EmitError("AbisCollection.loadInternal: error fetching ABI list", err)
			return
		} else {
			ac.mutex.Lock()
			ac.knownAbis = make([]coreTypes.Abi, 0, len(ac.knownAbis))
			for _, abi := range sdkAbis {
				if abi.IsKnown {
					ac.knownAbis = append(ac.knownAbis, abi)
				}
			}
			ac.isKnownLoaded = true
			ac.mutex.Unlock()

			finalStatus = fmt.Sprintf("Loaded %d known abis", len(ac.knownAbis))
			finalPayload.CurrentCount = len(ac.knownAbis)
			finalPayload.ExpectedTotal = len(ac.knownAbis)
		}
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
// loadStreamingData is a generic method that handles streaming data of any type T
func loadStreamingData[T any](
	ac *AbisCollection,
	contextKey string,
	queryFunc func(*output.RenderCtx),
	filterFunc func(item *T) bool,
	processItemFunc func(itemIntf interface{}) *T,
	targetSlice *[]T,
	expectedCount *int,
	loadedFlag *bool,
	dataTypeName string,
) (string, types.DataLoadedPayload, error) {
	ac.App.Cancel(contextKey)
	defer func() {
		ac.App.Cancel(contextKey)
	}()

	renderCtx := ac.App.RegisterCtx(contextKey)
	done := make(chan struct{})

	go func() {
		defer func() {
			if renderCtx.ModelChan != nil {
				close(renderCtx.ModelChan)
			}
			if renderCtx.ErrorChan != nil {
				close(renderCtx.ErrorChan)
			}
			close(done)
		}()

		queryFunc(renderCtx)
	}()

	modelChanClosed := false
	errorChanClosed := false

	for !modelChanClosed || !errorChanClosed {
		select {
		case itemIntf, ok := <-renderCtx.ModelChan:
			if !ok {
				modelChanClosed = true
				continue
			}

			itemPtr := processItemFunc(itemIntf)
			if itemPtr == nil {
				logger.Info(fmt.Sprintf("AbisCollection.loadStreamingData: unexpected item type: %T", itemIntf))
				continue
			}

			if filterFunc(itemPtr) {
				ac.mutex.Lock()
				*targetSlice = append(*targetSlice, *itemPtr)
				ac.mutex.Unlock()

				if len(*targetSlice)%refreshRate == 0 {
					ac.mutex.RLock()
					isLoaded := len(*targetSlice) >= *expectedCount
					payload := types.DataLoadedPayload{
						DataType:      dataTypeName,
						CurrentCount:  len(*targetSlice),
						ExpectedTotal: *expectedCount,
						IsLoaded:      isLoaded,
						Category:      "abis",
					}
					ac.App.EmitEvent(msgs.EventDataLoaded, payload)
					statusMsg := fmt.Sprintf("Loading %s: %d processed.", dataTypeName, len(*targetSlice))
					if *expectedCount > 0 {
						statusMsg = fmt.Sprintf("Loading %s: %d of %d processed.", dataTypeName, len(*targetSlice), *expectedCount)
					}
					msgs.EmitStatus(statusMsg)
					ac.mutex.RUnlock()
				}
			}

		case streamErr, ok := <-renderCtx.ErrorChan:
			if !ok {
				errorChanClosed = true
				continue
			}
			msgs.EmitError("AbisCollection.loadStreamingData: streaming error", streamErr)

		case <-done:
			// Stream initialization completed
		}
	}

	ac.mutex.Lock()
	*loadedFlag = true
	ac.mutex.Unlock()

	finalStatus := fmt.Sprintf("%s loaded: %d items.", dataTypeName, len(*targetSlice))
	finalPayload := types.DataLoadedPayload{
		Category:      "abis",
		DataType:      dataTypeName,
		IsLoaded:      true,
		CurrentCount:  len(*targetSlice),
		ExpectedTotal: len(*targetSlice),
	}

	return finalStatus, finalPayload, nil
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

	finalStatus, finalPayload, _ := loadStreamingData(
		ac,
		contextKey,
		queryFunc,
		filterFunc,
		processItemFunc,
		&ac.allFunctions,
		&ac.expectedFunctions,
		&ac.isFuncsLoaded,
		"functions-events",
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

	finalStatus, finalPayload, _ := loadStreamingData(
		ac,
		contextKey,
		queryFunc,
		filterFunc,
		processItemFunc,
		&ac.allEvents,
		&ac.expectedEvents,
		&ac.isEventsLoaded,
		"functions-events",
	)

	return finalStatus, finalPayload
}

// ADD_ROUTE
