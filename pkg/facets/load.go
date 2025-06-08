package facets

import (
	"errors"
	"fmt"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/store"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

var ErrorAlreadyLoading = errors.New("already loading")

// Load implements loading using a store instead of a queryFunc
func (r *BaseFacet[T]) Load() (*StreamingResult, error) {
	if !r.NeedsUpdate() {
		msgs.EmitStatus(fmt.Sprintf("cached: %d items", len(r.data)))
		cachedPayload := r.getCachedResult()
		return cachedPayload, nil
	}

	if !r.StartFetching() {
		return nil, ErrorAlreadyLoading
	}

	r.mutex.Lock()
	r.data = r.data[:0]
	r.mutex.Unlock()

	contextKey := fmt.Sprintf("facet-%s-%s", r.listKind, r.store.GetStoreType())
	finalPayload, err := store.ProcessStream(
		contextKey,
		r.store,
		r.filterFunc,
		r.isDupFunc,
		&r.data,
		&r.expectedCnt,
		r.listKind,
		&r.mutex,
		func() { r.SetPartial() },
	)

	// If we have partial data (even with an error), set state to partial instead of loaded
	if err == nil {
		// Successful completion
		r.state.Store(StateLoaded)
	} else {
		// Error occurred, but if we have data, mark as partial
		r.mutex.RLock()
		hasData := len(r.data) > 0
		r.mutex.RUnlock()

		if hasData {
			r.state.Store(StatePartial)
			// Return success with partial data instead of error
			return &StreamingResult{
				Payload: types.DataLoadedPayload{
					CurrentCount:  finalPayload.CurrentCount,
					ExpectedTotal: r.expectedCnt,
					ListKind:      r.listKind,
				},
				Error: nil, // Don't return error if we have partial data
			}, nil
		} else {
			r.state.Store(StateError)
		}
	}

	return &StreamingResult{
		Payload: finalPayload,
		Error:   err,
	}, nil
}

func (r *BaseFacet[T]) getCachedResult() *StreamingResult {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return &StreamingResult{
		Payload: types.DataLoadedPayload{
			CurrentCount:  len(r.data),
			ExpectedTotal: len(r.data),
			ListKind:      r.listKind,
		},
		Error: nil,
	}
}
