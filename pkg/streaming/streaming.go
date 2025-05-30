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
	app types.App,
	contextKey string,
	queryFunc func(*output.RenderCtx),
	filterFunc func(item *T) bool,
	processItemFunc func(itemIntf interface{}) *T,
	targetSlice *[]T,
	expectedCount *int,
	loadedFlag *bool,
	dataTypeName string,
	m interface {
		Lock()
		Unlock()
		RLock()
		RUnlock()
	},
) (string, types.DataLoadedPayload, error) {
	app.Cancel(contextKey)
	defer func() {
		app.Cancel(contextKey)
	}()

	renderCtx := app.RegisterCtx(contextKey)
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
				*targetSlice = append(*targetSlice, *itemPtr)
				m.Unlock()

				if len(*targetSlice)%refreshRate == 0 {
					m.RLock()
					isLoaded := len(*targetSlice) >= *expectedCount
					payload := types.DataLoadedPayload{
						DataType:      dataTypeName,
						CurrentCount:  len(*targetSlice),
						ExpectedTotal: *expectedCount,
						IsLoaded:      isLoaded,
						Category:      "abis",
					}
					app.EmitEvent(msgs.EventDataLoaded, payload)
					statusMsg := fmt.Sprintf("Loading %s: %d processed.", dataTypeName, len(*targetSlice))
					if *expectedCount > 0 {
						statusMsg = fmt.Sprintf("Loading %s: %d of %d processed.", dataTypeName, len(*targetSlice), *expectedCount)
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
