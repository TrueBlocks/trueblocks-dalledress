package sources

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
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

	// renderCtx := RegisterCtx(contextKey)
	renderCtx := output.NewStreamingContext()
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
			return true // Continue streaming
		})

		if err != nil {
			streamErr <- err
		}
	}()

	select {
	case <-done:
		// Streaming completed successfully
	case err := <-streamErr:
		msgs.EmitError("ProcessStream", err)
		return types.DataLoadedPayload{}, err
	}

	m.Lock()
	*loadedFlag = true
	itemCount := len(*targetSlice)
	m.Unlock()

	msgs.EmitStatus(fmt.Sprintf("loaded: %d items.", itemCount))
	return types.DataLoadedPayload{
		IsLoaded:      true,
		CurrentCount:  itemCount,
		ExpectedTotal: itemCount,
		ListKind:      listKind,
		Reason:        "final",
	}, nil
}
