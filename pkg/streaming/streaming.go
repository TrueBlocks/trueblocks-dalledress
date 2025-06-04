package streaming

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

const refreshRate = 31

// StreamData streams data of any type T with filtering, deduplication, and progress updates.
func StreamData[T any](
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
) (types.DataLoadedPayload, error) {
	Cancel(contextKey)
	defer func() { Cancel(contextKey) }()

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
				logger.Info(fmt.Sprintf("StreamData: unexpected item type: %T", itemIntf))
				continue
			}

			if filterFunc(itemPtr) {
				m.Lock()
				// Apply deduplication if provided
				if dedupeFunc == nil || !dedupeFunc(*targetSlice, itemPtr) {
					*targetSlice = append(*targetSlice, *itemPtr)
				}
				m.Unlock()

				isFirstPage := len(*targetSlice) == 7
				if isFirstPage || len(*targetSlice)%refreshRate == 0 {
					m.RLock()
					isLoaded := len(*targetSlice) >= *expectedCount
					reason := "partial"
					if isFirstPage {
						reason = "initial"
					}
					payload := types.DataLoadedPayload{
						CurrentCount:  len(*targetSlice),
						ExpectedTotal: *expectedCount,
						IsLoaded:      isLoaded,
						ListKind:      listKind,
						Reason:        reason,
					}
					msgs.EmitLoaded(reason, payload)
					msgs.EmitStatus(fmt.Sprintf("loading %s: %d processed.", listKind, len(*targetSlice)))
					m.RUnlock()
				}
			}

		case streamErr, ok := <-renderCtx.ErrorChan:
			if !ok {
				errorChanClosed = true
				continue
			}
			msgs.EmitError("StreamData", streamErr)

		case <-done:
			// Stream initialization completed
		}
	}

	m.Lock()
	*loadedFlag = true
	m.Unlock()

	msgs.EmitStatus(fmt.Sprintf("loaded: %d items.", len(*targetSlice)))
	return types.DataLoadedPayload{
		IsLoaded:      true,
		CurrentCount:  len(*targetSlice),
		ExpectedTotal: len(*targetSlice),
		ListKind:      listKind,
		Reason:        "final",
	}, nil
}
