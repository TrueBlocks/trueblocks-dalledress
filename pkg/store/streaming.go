package store

// Added for time-based updates

// // ProcessStream streams data from a Store with all the existing streaming features:
// // filtering, deduplication, progress report, and real-time messaging.
// func ProcessStream[T any](
// 	contextKey string,
// 	store *Store[T],
// 	filterFunc func(item *T) bool,
// 	isDupFunc func(existing []T, newItem *T) bool, // Returns true if item should be added (not a duplicate)
// 	targetSlice *[]T,
// 	expectedCnt *int,
// 	listKind types.ListKind,
// 	m interface {
// 		Lock()
// 		Unlock()
// 		RLock()
// 		RUnlock()
// 	},
// 	onFirstData func(),
// ) (types.DataLoadedPayload, error) {
// 	renderCtx := RegisterContext(contextKey)
// 	defer UnregisterContext(contextKey)

// 	done := make(chan struct{})
// 	streamErr := make(chan error, 1)

// 	reporter := progress.NewProgress(listKind, onFirstData)
// 	go func() {
// 		defer close(done)
// 		ticker := time.NewTicker(progress.MaxWaitTime / 2)
// 		defer ticker.Stop()

// 		fetchDone := make(chan error)
// 		go func() {
// 			err := store.Fetch(renderCtx, func(itemPtr *T) bool {
// 				if filterFunc(itemPtr) {
// 					var itemAdded bool
// 					var currentTotalCount int
// 					var currentExpectedCount int

// 					m.Lock()
// 					if isDupFunc == nil || !isDupFunc(*targetSlice, itemPtr) {
// 						*targetSlice = append(*targetSlice, *itemPtr)
// 						itemAdded = true
// 					}
// 					currentTotalCount = len(*targetSlice)
// 					if expectedCnt != nil {
// 						currentExpectedCount = *expectedCnt // Dereference under lock
// 					}
// 					m.Unlock()

// 					if !itemAdded {
// 						return true // Continue if item was a duplicate or filtered out by isDupFunc
// 					}

// 					reporter.Tick(currentTotalCount, currentExpectedCount)
// 				}
// 				return true
// 			})
// 			fetchDone <- err
// 		}()

// 		// Main loop for this goroutine: listen to fetch completion or ticker
// 		for {
// 			select {
// 			case err := <-fetchDone:
// 				if err != nil {
// 					streamErr <- err
// 				}
// 				return // Exit the select loop and the goroutine
// 			case <-ticker.C:
// 				var currentTotalCount int
// 				var currentExpectedCount int

// 				m.RLock()
// 				currentTotalCount = len(*targetSlice)
// 				if expectedCnt != nil {
// 					currentExpectedCount = *expectedCnt
// 				}
// 				m.RUnlock()

// 				reporter.HeartbeatUpdate(currentTotalCount, currentExpectedCount)
// 			}
// 		}
// 	}()

// 	var streamingError error
// 	select {
// 	case <-done:
// 		// Streaming completed successfully
// 	case err := <-streamErr:
// 		streamingError = err
// 		// Don't return immediately - we want to preserve partial data
// 	}

// 	m.RLock()
// 	itemCount := len(*targetSlice)
// 	m.RUnlock()

// 	if itemCount > 0 {
// 		payload := reporter.Tick(itemCount, itemCount)
// 		if streamingError != nil && itemCount > 0 {
// 			return payload, nil
// 		}
// 		return payload, streamingError
// 	}

// 	if streamingError != nil {
// 		return types.DataLoadedPayload{}, streamingError
// 	}

// 	payload := reporter.Tick(itemCount, itemCount)
// 	return payload, nil
// }
