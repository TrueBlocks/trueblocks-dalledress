package repository

import (
	"errors"
	"fmt"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/streaming"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

var ErrorAlreadyLoading = errors.New("already loading")

// Load implements Repository.Load using the streaming system
func (r *BaseRepository[T]) Load(opts LoadOptions) (*StreamingResult, error) {
	if !r.state.StartLoading() {
		return nil, ErrorAlreadyLoading
	}
	defer r.state.StopLoading()

	if !opts.ForceReload && r.state.IsLoaded() {
		msgs.EmitStatus(fmt.Sprintf("cached: %d items", len(r.data)))
		cachedPayload := r.getCachedResult()
		return cachedPayload, nil
	}

	r.mutex.Lock()
	r.data = r.data[:0]
	r.loaded = false
	r.mutex.Unlock()

	contextKey := fmt.Sprintf("repo-%s", r.listKind)
	finalPayload, err := streaming.StreamData(
		contextKey,
		r.queryFunc,
		r.filterFunc,
		r.processFunc,
		r.dedupeFunc,
		&r.data,
		&r.expectedCount,
		&r.loaded,
		r.listKind,
		&r.mutex,
	)
	return &StreamingResult{
		Payload: finalPayload,
		Error:   err,
	}, nil
}

func (r *BaseRepository[T]) getCachedResult() *StreamingResult {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	dataCopy := make([]T, len(r.data))
	copy(dataCopy, r.data)
	return &StreamingResult{
		Payload: types.DataLoadedPayload{
			CurrentCount:  len(r.data),
			ExpectedTotal: len(r.data),
			IsLoaded:      true,
			ListKind:      r.listKind,
			Reason:        "initial",
		},
		Error: nil,
	}
}
