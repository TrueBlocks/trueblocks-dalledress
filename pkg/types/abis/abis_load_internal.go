// ADD_ROUTE
package abis

import (
	"fmt"
	"sync/atomic"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
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

	// At the start, reset flags and slices under a write lock
	ac.mutex.Lock()
	switch listKind {
	case AbisDownloaded:
		ac.downloadedAbis = make([]coreTypes.Abi, 0)
		ac.isDownloadedLoaded = false
		ac.expectedDownloaded = 0
	case AbisKnown:
		ac.knownAbis = make([]coreTypes.Abi, 0)
		ac.isKnownLoaded = false
		ac.expectedKnown = 0
	case AbisFunctions:
		ac.allFunctions = make([]coreTypes.Function, 0)
		ac.isFuncsLoaded = false
		ac.expectedFunctions = 0
	case AbisEvents:
		ac.allEvents = make([]coreTypes.Function, 0)
		ac.isEventsLoaded = false
		ac.expectedEvents = 0
	}
	ac.deduper = make(map[string]struct{})
	ac.mutex.Unlock()

	addrMap := make(map[base.Address]bool)

	// Always fetch all ABIs first to ensure we have the addresses for functions/events
	listOpts := sdk.AbisOptions{
		Globals: sdk.Globals{Cache: true, Verbose: true},
	}
	sdkAbis, _, err := listOpts.AbisList()
	if err != nil {
		msgs.EmitError("AbisCollection.loadInternal: error fetching ABI list", err)
		return
	}

	if sdkAbis == nil {
		sdkAbis = make([]coreTypes.Abi, 0)
	}

	ac.mutex.Lock()
	for _, abi := range sdkAbis {
		if abi.IsKnown {
			ac.knownAbis = append(ac.knownAbis, abi)
		} else {
			if _, exists := addrMap[abi.Address]; !exists {
				ac.downloadedAbis = append(ac.downloadedAbis, abi)
				ac.expectedFunctions += int(abi.NFunctions)
				ac.expectedEvents += int(abi.NEvents)
				addrMap[abi.Address] = true
			}
		}
	}
	ac.expectedDownloaded = len(ac.downloadedAbis)
	ac.expectedKnown = len(ac.knownAbis)
	currentCount := 0
	msgs.EmitStatus(fmt.Sprintf("Loading ABIs: %d, %d. Fetching details...", len(ac.downloadedAbis), len(ac.knownAbis)))
	ac.mutex.Unlock()

	if len(addrMap) == 0 {
		logger.Info(fmt.Sprintln("AbisCollection.loadInternal: no non-known ABI addresses to fetch detailed functions/events."))
		ac.mutex.Lock()
		switch listKind {
		case AbisDownloaded:
			ac.isDownloadedLoaded = true
		case AbisKnown:
			ac.isKnownLoaded = true
		case AbisFunctions:
			ac.isFuncsLoaded = true
		case AbisEvents:
			ac.isEventsLoaded = true
		}
		ac.mutex.Unlock()

		statusMsg := "ABI details loaded: No new items to fetch."
		msgs.EmitStatus(statusMsg)
		ac.App.EmitEvent(msgs.EventDataLoaded, types.DataLoadedPayload{
			DataType:      "functions-events",
			CurrentCount:  currentCount,
			ExpectedTotal: ac.expectedFunctions + ac.expectedEvents,
			IsLoaded:      true,
			Category:      "abis",
		})
		return
	}

	abisAddr := base.ZeroAddr
	renderCtx := ac.App.RegisterCtx(abisAddr)

	// Ensure context is cancelled/cleaned up when we're done
	defer func() {
		ac.App.Cancel(abisAddr)
	}()

	terms := make([]string, 0, len(addrMap))
	for addr := range addrMap {
		terms = append(terms, addr.Hex())
	}
	detailOpts := sdk.AbisOptions{
		Globals:   sdk.Globals{Cache: true},
		Addrs:     terms,
		RenderCtx: renderCtx,
	}

	go func() {
		defer func() {
			if renderCtx.ModelChan != nil {
				close(renderCtx.ModelChan)
			}
			if renderCtx.ErrorChan != nil {
				close(renderCtx.ErrorChan)
			}
		}()

		if listKind == AbisFunctions {
			_, _, streamInitiationErr := detailOpts.AbisListFuncs()
			if streamInitiationErr != nil {
				logger.Info(fmt.Sprintf("AbisCollection.loadInternal: error initiating stream: %v", streamInitiationErr))
			}
		} else if listKind == AbisEvents {
			_, _, streamInitiationErr := detailOpts.AbisListEvents()
			if streamInitiationErr != nil {
				logger.Info(fmt.Sprintf("AbisCollection.loadInternal: error initiating stream: %v", streamInitiationErr))
			}
		} else {
			logger.Error(fmt.Sprintf("AbisCollection.loadInternal: unexpected list kind: %v", listKind))
			return
		}
	}()

ProcessingLoop:
	for {
		select {
		case itemIntf, ok := <-renderCtx.ModelChan:
			if !ok {
				renderCtx.ModelChan = nil
				if renderCtx.ErrorChan == nil {
					break ProcessingLoop
				}
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
			} else {
				ac.allFunctions = append(ac.allFunctions, *itemPtr)
			}
			currentCount := len(ac.allFunctions) + len(ac.allEvents)
			ac.mutex.Unlock()

			if currentCount%refreshRate == 0 {
				ac.mutex.RLock()
				isLoaded := false
				switch listKind {
				case AbisDownloaded, AbisKnown:
					isLoaded = currentCount >= ac.expectedDownloaded
				case AbisFunctions, AbisEvents:
					isLoaded = currentCount >= (ac.expectedFunctions + ac.expectedEvents)
				}
				payload := types.DataLoadedPayload{
					DataType:      "functions-events",
					CurrentCount:  currentCount,
					ExpectedTotal: ac.expectedFunctions + ac.expectedEvents,
					IsLoaded:      isLoaded,
					Category:      "abis",
				}
				ac.App.EmitEvent(msgs.EventDataLoaded, payload)
				statusMsg := fmt.Sprintf("Loading ABI details: %d processed.", currentCount)
				if (ac.expectedFunctions + ac.expectedEvents) > 0 {
					statusMsg = fmt.Sprintf("Loading ABI details: %d of %d processed.", currentCount, ac.expectedFunctions+ac.expectedEvents)
				}
				msgs.EmitStatus(statusMsg)
				ac.mutex.RUnlock()
			}

		case streamErr, ok := <-renderCtx.ErrorChan:
			if !ok {
				renderCtx.ErrorChan = nil
				if renderCtx.ModelChan == nil {
					break ProcessingLoop
				}
				continue
			}
			msgs.EmitError("AbisCollection.loadInternal: streaming error", streamErr)
		}
	}

	ac.mutex.Lock()
	switch listKind {
	case AbisDownloaded:
		ac.isDownloadedLoaded = true
	case AbisKnown:
		ac.isKnownLoaded = true
	case AbisFunctions:
		ac.isFuncsLoaded = true
	case AbisEvents:
		ac.isEventsLoaded = true
	}
	finalCount := len(ac.allFunctions) + len(ac.allEvents)
	finalStatus := fmt.Sprintf("ABI details fully loaded: %d functions/events.", finalCount)
	msgs.EmitStatus(finalStatus)
	ac.App.EmitEvent(msgs.EventDataLoaded, types.DataLoadedPayload{
		DataType:      "functions-events",
		CurrentCount:  finalCount,
		ExpectedTotal: ac.expectedFunctions + ac.expectedEvents,
		IsLoaded:      true,
		Category:      "abis",
	})
	ac.mutex.Unlock()
}

// ADD_ROUTE
