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

	if listKind == AbisDownloaded || listKind == AbisKnown {
		listOpts := sdk.AbisOptions{
			Globals: sdk.Globals{Cache: true, Verbose: true},
		}
		if sdkAbis, _, err := listOpts.AbisList(); err != nil || sdkAbis == nil {
			msgs.EmitError("AbisCollection.loadInternal: error fetching ABI list", err)
			return
		} else {
			ac.mutex.Lock()
			ac.knownAbis = make([]coreTypes.Abi, 0, len(ac.knownAbis))
			ac.downloadedAbis = make([]coreTypes.Abi, 0, len(ac.downloadedAbis))
			for _, abi := range sdkAbis {
				if abi.IsKnown {
					ac.knownAbis = append(ac.knownAbis, abi)
				} else {
					ac.downloadedAbis = append(ac.downloadedAbis, abi)
				}
			}
			ac.isDownloadedLoaded = true
			ac.isKnownLoaded = true

			msgs.EmitStatus(fmt.Sprintf("Loaded %d downloaded and %d known abis", len(ac.downloadedAbis), len(ac.knownAbis)))
			ac.App.EmitEvent(msgs.EventDataLoaded, types.DataLoadedPayload{
				Category:      "abis",
				DataType:      "functions-events",
				CurrentCount:  len(ac.downloadedAbis) + len(ac.knownAbis),
				ExpectedTotal: len(ac.downloadedAbis) + len(ac.knownAbis),
				IsLoaded:      true,
			})
			ac.mutex.Unlock()
		}

	} else {
		contextKey := "abis-load-internal"
		defer func() {
			ac.App.Cancel(contextKey)
		}()

		renderCtx := ac.App.RegisterCtx(contextKey)
		detailOpts := sdk.AbisOptions{
			Globals:   sdk.Globals{Cache: true},
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
				_, _, streamInitiationErr := detailOpts.AbisDetails()
				if streamInitiationErr != nil {
					logger.Info(fmt.Sprintf("AbisCollection.loadInternal: error initiating stream: %v", streamInitiationErr))
				}
			} else if listKind == AbisEvents {
				_, _, streamInitiationErr := detailOpts.AbisDetails()
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
					isLoaded := currentCount >= (ac.expectedFunctions + ac.expectedEvents)
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
}

// ADD_ROUTE
