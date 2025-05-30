// ADD_ROUTE
package abis

import (
	"fmt"
	"sync/atomic"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

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

// loadFunctions loads ABI functions asynchronously and returns status, payload, and error
func (ac *AbisCollection) loadFunctions() (string, types.DataLoadedPayload) {
	contextKey := "abis-load-internal-functions"
	ac.App.Cancel(contextKey) // cancel any previous context for this key
	defer func() {
		ac.App.Cancel(contextKey) // ensure we clean up the context
	}()

	renderCtx := ac.App.RegisterCtx(contextKey)
	detailOpts := sdk.AbisOptions{
		Globals:   sdk.Globals{Cache: true},
		RenderCtx: renderCtx,
	}

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

		_, _, streamInitiationErr := detailOpts.AbisDetails()
		if streamInitiationErr != nil {
			logger.Info(fmt.Sprintf("AbisCollection.loadInternal: error initiating stream: %v", streamInitiationErr))
		}
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

			itemPtr, okAssert := itemIntf.(*coreTypes.Function)
			if !okAssert {
				logger.Info(fmt.Sprintf("AbisCollection.loadInternal: unexpected item type: %T", itemIntf))
				continue
			}

			ac.mutex.Lock()
			if _, exists := ac.deduper[itemPtr.Encoding]; exists {
				ac.mutex.Unlock()
				continue
			}

			ac.deduper[itemPtr.Encoding] = struct{}{}
			if itemPtr.FunctionType != "event" {
				ac.allFunctions = append(ac.allFunctions, *itemPtr)
			}
			ac.mutex.Unlock()

			if len(ac.allFunctions)%refreshRate == 0 {
				ac.mutex.RLock()
				isLoaded := len(ac.allFunctions) >= (ac.expectedFunctions)
				payload := types.DataLoadedPayload{
					DataType:      "functions-events",
					CurrentCount:  len(ac.allFunctions),
					ExpectedTotal: ac.expectedFunctions,
					IsLoaded:      isLoaded,
					Category:      "abis",
				}
				ac.App.EmitEvent(msgs.EventDataLoaded, payload)
				statusMsg := fmt.Sprintf("Loading ABI details: %d processed.", len(ac.allFunctions))
				if (ac.expectedFunctions) > 0 {
					statusMsg = fmt.Sprintf("Loading ABI details: %d of %d processed.", len(ac.allFunctions), ac.expectedFunctions)
				}
				msgs.EmitStatus(statusMsg)
				ac.mutex.RUnlock()
			}

		case streamErr, ok := <-renderCtx.ErrorChan:
			if !ok {
				errorChanClosed = true
				continue
			}
			msgs.EmitError("AbisCollection.loadInternal: streaming error", streamErr)

		case <-done:
			// Stream initialization completed
		}
	}

	ac.mutex.Lock()
	ac.isFuncsLoaded = true
	ac.mutex.Unlock()

	finalStatus := fmt.Sprintf("ABI function details loaded: %d functions.", len(ac.allFunctions))
	finalPayload := types.DataLoadedPayload{
		Category:      "abis",
		DataType:      "functions-events",
		IsLoaded:      true,
		CurrentCount:  len(ac.allFunctions),
		ExpectedTotal: len(ac.allFunctions),
	}

	return finalStatus, finalPayload
}

// loadEvents loads ABI events asynchronously and returns status, payload
func (ac *AbisCollection) loadEvents() (string, types.DataLoadedPayload) {
	contextKey := "abis-load-internal-events"
	ac.App.Cancel(contextKey) // cancel any previous context for this key
	defer func() {
		ac.App.Cancel(contextKey) // ensure we clean up the context
	}()

	renderCtx := ac.App.RegisterCtx(contextKey)
	detailOpts := sdk.AbisOptions{
		Globals:   sdk.Globals{Cache: true},
		RenderCtx: renderCtx,
	}

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

		_, _, streamInitiationErr := detailOpts.AbisDetails()
		if streamInitiationErr != nil {
			logger.Info(fmt.Sprintf("AbisCollection.loadInternal: error initiating stream: %v", streamInitiationErr))
		}
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

			itemPtr, okAssert := itemIntf.(*coreTypes.Function)
			if !okAssert {
				logger.Info(fmt.Sprintf("AbisCollection.loadInternal: unexpected item type: %T", itemIntf))
				continue
			}

			ac.mutex.Lock()
			if _, exists := ac.deduper[itemPtr.Encoding]; exists {
				ac.mutex.Unlock()
				continue
			}

			ac.deduper[itemPtr.Encoding] = struct{}{}
			if itemPtr.FunctionType == "event" {
				ac.allEvents = append(ac.allEvents, *itemPtr)
			}
			ac.mutex.Unlock()

			if len(ac.allEvents)%refreshRate == 0 {
				ac.mutex.RLock()
				isLoaded := len(ac.allEvents) >= (ac.expectedEvents)
				payload := types.DataLoadedPayload{
					DataType:      "functions-events",
					CurrentCount:  len(ac.allEvents),
					ExpectedTotal: ac.expectedEvents,
					IsLoaded:      isLoaded,
					Category:      "abis",
				}
				ac.App.EmitEvent(msgs.EventDataLoaded, payload)
				statusMsg := fmt.Sprintf("Loading ABI details: %d processed.", len(ac.allEvents))
				if (ac.expectedEvents) > 0 {
					statusMsg = fmt.Sprintf("Loading ABI details: %d of %d processed.", len(ac.allEvents), ac.expectedEvents)
				}
				msgs.EmitStatus(statusMsg)
				ac.mutex.RUnlock()
			}

		case streamErr, ok := <-renderCtx.ErrorChan:
			if !ok {
				errorChanClosed = true
				continue
			}
			msgs.EmitError("AbisCollection.loadInternal: streaming error", streamErr)

		case <-done:
			// Stream initialization completed
		}
	}

	ac.mutex.Lock()
	ac.isEventsLoaded = true
	ac.mutex.Unlock()

	finalStatus := fmt.Sprintf("ABI event details loaded: %d events.", len(ac.allEvents))
	finalPayload := types.DataLoadedPayload{
		Category:      "abis",
		DataType:      "functions-events",
		IsLoaded:      true,
		CurrentCount:  len(ac.allEvents),
		ExpectedTotal: len(ac.allEvents),
	}

	return finalStatus, finalPayload
}

// ADD_ROUTE
