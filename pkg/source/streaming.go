package source

import (
	"time" // Added for time-based updates

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

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
	onFirstDataCallback func(), // Renamed to avoid confusion, passed to Progress
) (types.DataLoadedPayload, error) {
	renderCtx := RegisterContext(contextKey)
	defer UnregisterContext(contextKey)

	done := make(chan struct{})
	streamErr := make(chan error, 1)

	reporter := NewProgress(listKind, onFirstDataCallback)
	go func() {
		defer close(done)
		ticker := time.NewTicker(maxWaitTime / 2) // expected time 1/2 of maxWaitTime
		defer ticker.Stop()

		fetchDone := make(chan error)
		go func() {
			err := source.Fetch(renderCtx, func(itemPtr *T) bool {
				if filterFunc(itemPtr) {
					var itemAdded bool
					var currentTotalCount int
					var currentExpectedCount int // To store dereferenced value of expectedCnt

					m.Lock()
					// Check for duplicates and add item if not a duplicate
					if isDupFunc == nil || !isDupFunc(*targetSlice, itemPtr) {
						*targetSlice = append(*targetSlice, *itemPtr)
						itemAdded = true
					}
					currentTotalCount = len(*targetSlice)
					if expectedCnt != nil {
						currentExpectedCount = *expectedCnt // Dereference under lock
					}
					m.Unlock()

					if !itemAdded {
						return true // Continue if item was a duplicate or filtered out by isDupFunc
					}

					reporter.Tick(currentTotalCount, currentExpectedCount)
				}
				return true
			})
			fetchDone <- err
		}()

		// Main loop for this goroutine: listen to fetch completion or ticker
		for {
			select {
			case err := <-fetchDone:
				if err != nil {
					streamErr <- err
				}
				return // Exit the select loop and the goroutine
			case <-ticker.C:
				var currentTotalCount int
				var currentExpectedCount int

				m.RLock() // Use RLock for reading counts
				currentTotalCount = len(*targetSlice)
				if expectedCnt != nil {
					currentExpectedCount = *expectedCnt
				}
				m.RUnlock()

				// Check for heartbeat update using the reporter
				reporter.HeartbeatUpdate(currentTotalCount, currentExpectedCount)
			}
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

	m.RLock() // Use RLock for final count read
	itemCount := len(*targetSlice)
	m.RUnlock()

	if itemCount > 0 {
		payload := reporter.Tick(itemCount, itemCount)
		if streamingError != nil && itemCount > 0 {
			return payload, nil
		}
		return payload, streamingError
	}

	// Only return error if we have no data at all and an error occurred
	if streamingError != nil {
		return types.DataLoadedPayload{}, streamingError
	}

	payload := reporter.Tick(itemCount, itemCount)
	return payload, nil
}
