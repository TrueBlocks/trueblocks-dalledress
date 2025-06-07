package sources

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

const refreshRate = 31

// ProcessStream streams data from a Source with all the existing streaming features:
// filtering, deduplication, progress updates, and real-time messaging.
func ProcessStream[T any](
	contextKey string,
	source Source[T],
	filterFunc func(item *T) bool,
	isDupFunc func(existing []T, newItem *T) bool, // Returns true if item should be added (not a duplicate)
	targetSlice *[]T,
	expectedCnt *int,
	listKind types.ListKind,
	m interface {
		Lock()
		Unlock()
		RLock()
		RUnlock()
	},
	onFirstData func(),
) (types.DataLoadedPayload, error) {
	renderCtx := RegisterContext(contextKey)
	defer UnregisterContext(contextKey)

	done := make(chan struct{})
	streamErr := make(chan error, 1)
	go func() {
		defer close(done)

		err := source.Fetch(renderCtx, func(itemPtr *T) bool {
			if filterFunc(itemPtr) {
				m.Lock()
				// Apply deduplication if provided
				if isDupFunc == nil || !isDupFunc(*targetSlice, itemPtr) {
					*targetSlice = append(*targetSlice, *itemPtr)
				}
				m.Unlock()

				isFirstPage := len(*targetSlice) == 7
				if isFirstPage || len(*targetSlice)%refreshRate == 0 {
					if isFirstPage && onFirstData != nil {
						onFirstData()
					}

					m.RLock()
					reason := "partial"
					if isFirstPage {
						reason = "initial"
					}
					payload := types.DataLoadedPayload{
						CurrentCount:  len(*targetSlice),
						ExpectedTotal: *expectedCnt,
						ListKind:      listKind,
						Reason:        reason,
					}
					msgs.EmitLoaded(reason, payload)
					msgs.EmitStatus(fmt.Sprintf("loading %s: %d processed.", listKind, len(*targetSlice)))
					m.RUnlock()
				}
			}
			return true // Continue streaming
		})

		if err != nil {
			streamErr <- err
		}
	}()

	var streamingError error
	select {
	case <-done:
		// Streaming completed successfully
	case err := <-streamErr:
		streamingError = err
		// Don't return immediately - we want to preserve partial data
	}

	m.Lock()
	itemCount := len(*targetSlice)
	m.Unlock()

	// If we have partial data, treat it as successful partial load
	if itemCount > 0 {
		reason := "final"
		if streamingError != nil {
			reason = "partial"
		}

		msgs.EmitStatus(fmt.Sprintf("loaded: %d items.", itemCount))
		return types.DataLoadedPayload{
			CurrentCount:  itemCount,
			ExpectedTotal: itemCount,
			ListKind:      listKind,
			Reason:        reason,
		}, nil
	}

	// Only return error if we have no data at all
	if streamingError != nil {
		return types.DataLoadedPayload{}, streamingError
	}

	msgs.EmitStatus(fmt.Sprintf("loaded: %d items.", itemCount))
	return types.DataLoadedPayload{
		CurrentCount:  itemCount,
		ExpectedTotal: itemCount,
		ListKind:      listKind,
		Reason:        "final",
	}, nil
}
