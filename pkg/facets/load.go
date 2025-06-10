package facets

// var ErrorAlreadyLoading = errors.New("already loading")

// // Load implements loading using a store instead of a queryFunc
// func (r *BaseFacet[T]) Load() (*StreamingResult, error) {
// 	if !r.NeedsUpdate() {
// 		msgs.EmitStatus(fmt.Sprintf("cached: %d items", len(r.data)))
// 		cachedPayload := r.getCachedResult()
// 		return cachedPayload, nil
// 	}

// 	if !r.StartFetching() {
// 		return nil, ErrorAlreadyLoading
// 	}

// 	r.mutex.Lock()
// 	r.data = r.data[:0]
// 	r.mutex.Unlock()

// 	// Start asynchronous loading in a goroutine
// 	go func() {
// 		contextKey := fmt.Sprintf("facet-%s-%s", r.listKind, r.store.GetStoreType())

// 		// We'll bypass ProcessStream and call store.Fetch directly to populate r.data
// 		renderCtx := store.RegisterContext(contextKey)
// 		defer store.UnregisterContext(contextKey)

// 		err := r.store.Fetch(renderCtx, func(itemPtr *T) bool {
// 			if r.filterFunc(itemPtr) {
// 				r.mutex.Lock()
// 				// Check for duplicates if needed
// 				if r.isDupFunc == nil || !r.isDupFunc(r.data, itemPtr) {
// 					r.data = append(r.data, itemPtr)
// 				}
// 				currentCount := len(r.data)
// 				r.mutex.Unlock()

// 				// Emit progress updates periodically
// 				if currentCount%10 == 0 || currentCount <= 10 {
// 					r.SetPartial()
// 					logging.LogBackend(fmt.Sprintf("BaseFacet.Load: kind: %s currentCount: %d expectedTotal: %d", r.listKind, currentCount, r.expectedCnt))
// 					msgs.EmitLoaded("streaming", types.DataLoadedPayload{
// 						CurrentCount:  currentCount,
// 						ExpectedTotal: r.expectedCnt,
// 						ListKind:      r.listKind,
// 					})
// 					// Also emit status message for the StatusBar
// 					msgs.EmitStatus(fmt.Sprintf("Loaded %d items...", currentCount))
// 				}
// 			}
// 			return true
// 		})

// 		// If we have partial data (even with an error), set state to partial instead of loaded
// 		if err == nil {
// 			// Successful completion
// 			r.state.Store(StateLoaded)
// 			finalCount := len(r.data) // No lock needed for final count
// 			msgs.EmitStatus(fmt.Sprintf("Loaded %d items", finalCount))
// 		} else {
// 			// Error occurred, but if we have data, mark as partial
// 			hasData := len(r.data) > 0 // No lock needed for interim check
// 			if hasData {
// 				r.state.Store(StatePartial)
// 				// currentCount := len(r.data) // No lock needed for interim count
// 				// msgs.EmitStatus(fmt.Sprintf("Partial load: %d items (error: %v)", currentCount, err))
// 			} else {
// 				r.state.Store(StateError)
// 				// msgs.EmitError(fmt.Sprintf("Load failed: %v", err), err)
// 			}
// 		}

// 		// Note: Final event emission is handled by the progress reporter in ProcessStream
// 	}()

// 	// Return immediately with initial state
// 	return &StreamingResult{
// 		Payload: types.DataLoadedPayload{
// 			CurrentCount:  0,
// 			ExpectedTotal: 0,
// 			ListKind:      r.listKind,
// 		},
// 		Error: nil,
// 	}, nil
// }

// func (r *BaseFacet[T]) getCachedResult() *StreamingResult {
// 	r.mutex.RLock()
// 	defer r.mutex.RUnlock()
// 	return &StreamingResult{
// 		Payload: types.DataLoadedPayload{
// 			CurrentCount:  len(r.data),
// 			ExpectedTotal: len(r.data),
// 			ListKind:      r.listKind,
// 		},
// 		Error: nil,
// 	}
// }
