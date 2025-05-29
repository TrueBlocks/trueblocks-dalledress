// ADD_ROUTE
package types

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

var refreshRate = 31

func (ac *AbisCollection) EnsureInitialLoad() {
	ac.mutex.Lock()
	if !ac.isLoaded && !ac.isLoading {
		ac.isLoading = true
		go ac.loadInternal()
	}
	ac.mutex.Unlock()
}

func (ac *AbisCollection) loadInternal() {
	defer func() {
		ac.mutex.Lock()
		ac.isLoading = false
		ac.mutex.Unlock()
	}()

	// At the start, reset flags and slices under a write lock
	ac.mutex.Lock()
	ac.downloadedAbis = make([]coreTypes.Abi, 0)
	ac.knownAbis = make([]coreTypes.Abi, 0)
	ac.allFunctions = make([]coreTypes.Function, 0)
	ac.allEvents = make([]coreTypes.Function, 0)
	ac.seenEncodings = make(map[string]bool) // Reset deduplication map
	ac.isLoaded = false
	ac.expectedFunctions = 0
	ac.expectedEvents = 0
	ac.expectedDownloaded = 0
	ac.expectedKnown = 0

	ac.mutex.Unlock()

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

	addrMap := make(map[base.Address]bool)
	ac.mutex.Lock()
	for _, abi := range sdkAbis {
		if abi.IsKnown {
			ac.knownAbis = append(ac.knownAbis, abi)
		} else {
			if _, exists := addrMap[abi.Address]; !exists {
				ac.downloadedAbis = append(ac.downloadedAbis, abi)
				ac.expectedFunctions += int(abi.NFunctions)
				ac.expectedEvents += int(abi.NEvents)
			}
			addrMap[abi.Address] = true
		}
	}
	ac.expectedDownloaded = len(ac.downloadedAbis)
	ac.expectedKnown = len(ac.knownAbis)
	currentCount := 0
	msgs.EmitStatus(fmt.Sprintf("Loading ABIs: %d downloaded, %d known. Fetching details...", len(ac.downloadedAbis), len(ac.knownAbis)))
	ac.mutex.Unlock()

	if len(addrMap) == 0 {
		logger.Info(fmt.Sprintln("AbisCollection.loadInternal: no non-known ABI addresses to fetch detailed functions/events."))
		ac.mutex.Lock()
		ac.isLoaded = true
		statusMsg := "ABI details loaded: No new items to fetch."
		msgs.EmitStatus(statusMsg)
		ac.App.EmitEvent(msgs.EventDataLoaded, DataLoadedPayload{
			DataType:      "functions-events",
			CurrentCount:  currentCount,
			ExpectedTotal: ac.expectedFunctions + ac.expectedEvents,
			IsLoaded:      true,
			Category:      "abis",
		})
		ac.mutex.Unlock()
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

		_, _, streamInitiationErr := detailOpts.AbisListItems()
		if streamInitiationErr != nil {
			logger.Info(fmt.Sprintf("AbisCollection.loadInternal: error initiating stream: %v", streamInitiationErr))
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
			if ac.seenEncodings[itemPtr.Encoding] {
				ac.mutex.Unlock()
				continue
			}

			ac.seenEncodings[itemPtr.Encoding] = true
			if itemPtr.FunctionType == "event" {
				ac.allEvents = append(ac.allEvents, *itemPtr)
			} else {
				ac.allFunctions = append(ac.allFunctions, *itemPtr)
			}
			currentCount = len(ac.allFunctions) + len(ac.allEvents)
			ac.mutex.Unlock()

			if currentCount%refreshRate == 0 {
				ac.mutex.RLock()
				payload := DataLoadedPayload{
					DataType:      "functions-events",
					CurrentCount:  currentCount,
					ExpectedTotal: ac.expectedFunctions + ac.expectedEvents,
					IsLoaded:      false,
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
	ac.isLoaded = true
	finalCount := len(ac.allFunctions) + len(ac.allEvents)
	finalStatus := fmt.Sprintf("ABI details fully loaded: %d functions/events.", finalCount)
	msgs.EmitStatus(finalStatus)
	ac.App.EmitEvent(msgs.EventDataLoaded, DataLoadedPayload{
		DataType:      "functions-events",
		CurrentCount:  finalCount,
		ExpectedTotal: ac.expectedFunctions + ac.expectedEvents,
		IsLoaded:      true,
		Category:      "abis",
	})
	ac.mutex.Unlock()
}

// ADD_ROUTE
