package streaming

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

const refreshRate = 31

// LoadStreamingData is a generic method that handles streaming data of any type T
func LoadStreamingData[T any](
	contextKey string,
	queryFunc func(*output.RenderCtx),
	filterFunc func(item *T) bool,
	processItemFunc func(itemIntf interface{}) *T,
	dedupeFunc func(existing []T, newItem *T) bool, // Returns true if item should be added (not a duplicate)
	targetSlice *[]T,
	expectedCount *int,
	loadedFlag *bool,
	listKind types.ListKind,
	m interface {
		Lock()
		Unlock()
		RLock()
		RUnlock()
	},
) (string, types.DataLoadedPayload, error) {
	Cancel(contextKey)
	defer func() {
		Cancel(contextKey)
	}()

	renderCtx := RegisterCtx(contextKey)
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
				logger.Info(fmt.Sprintf("LoadStreamingData: unexpected item type: %T", itemIntf))
				continue
			}

			if filterFunc(itemPtr) {
				m.Lock()
				// Apply deduplication if provided
				if dedupeFunc == nil || !dedupeFunc(*targetSlice, itemPtr) {
					*targetSlice = append(*targetSlice, *itemPtr)
				}
				m.Unlock()

				if len(*targetSlice) == 7 || (len(*targetSlice) > 7 && len(*targetSlice)%refreshRate == 0) {
					m.RLock()
					isLoaded := len(*targetSlice) >= *expectedCount
					payload := types.DataLoadedPayload{
						CurrentCount:  len(*targetSlice),
						ExpectedTotal: *expectedCount,
						IsLoaded:      isLoaded,
						ListKind:      listKind,
					}
					msgs.EmitPayload(msgs.EventDataLoaded, payload)
					statusMsg := fmt.Sprintf("Loading %s: %d processed.", listKind, len(*targetSlice))
					if *expectedCount > 0 {
						statusMsg = fmt.Sprintf("Loading %s: %d of %d processed.", listKind, len(*targetSlice), *expectedCount)
					}
					msgs.EmitStatus(statusMsg)
					m.RUnlock()
				}
			}

		case streamErr, ok := <-renderCtx.ErrorChan:
			if !ok {
				errorChanClosed = true
				continue
			}
			msgs.EmitError("LoadStreamingData: streaming error", streamErr)

		case <-done:
			// Stream initialization completed
		}
	}

	m.Lock()
	*loadedFlag = true
	m.Unlock()

	finalStatus := fmt.Sprintf("loaded: %d items.", len(*targetSlice))
	finalPayload := types.DataLoadedPayload{
		IsLoaded:      true,
		CurrentCount:  len(*targetSlice),
		ExpectedTotal: len(*targetSlice),
		ListKind:      listKind,
	}

	return finalStatus, finalPayload, nil
}
