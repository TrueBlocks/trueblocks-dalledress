package facets

import (
	"errors"
	"fmt"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/sources"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

var ErrorAlreadyLoading = errors.New("already loading")

// Load implements loading using a source instead of a queryFunc
func (r *BaseFacet[T]) Load(opts LoadOptions) (*StreamingResult, error) {
	if !r.StartFetching() {
		return nil, ErrorAlreadyLoading
	}
	defer r.StopFetching()

	if !opts.ForceReload && r.IsLoaded() {
		msgs.EmitStatus(fmt.Sprintf("cached: %d items", len(r.data)))
		cachedPayload := r.getCachedResult()
		return cachedPayload, nil
	}

	r.mutex.Lock()
	r.data = r.data[:0]
	r.mutex.Unlock()

	contextKey := fmt.Sprintf("facet-%s-%s", r.listKind, r.source.GetSourceType())
	finalPayload, err := sources.ProcessStream(
		contextKey,
		r.source,
		r.filterFunc,
		r.isDupFunc,
		&r.data,
		&r.expectedCnt,
		&r.loaded,
		r.listKind,
		&r.mutex,
	)
	return &StreamingResult{
		Payload: finalPayload,
		Error:   err,
	}, nil
}

func (r *BaseFacet[T]) getCachedResult() *StreamingResult {
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
