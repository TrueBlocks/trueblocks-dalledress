package repository

import (
	"errors"
	"fmt"
	"sync"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/streaming"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

// Repository defines the contract for data access
// Provides loading, paging, filtering, sorting, and cache management for type T
type Repository[T any] interface {
	Load(opts LoadOptions) (*StreamingResult, error)
	GetPage(first, pageSize int, filter FilterFunc[T], sortSpec interface{}, sortFunc func([]T, interface{}) error) (*PageResult[T], error)
	IsLoading() bool
	IsLoaded() bool
	NeedsUpdate() bool
	Clear()
	Remove(predicate func(*T) bool) bool
	Count() int
	ExpectedCount() int
}

type FilterFunc[T any] func(*T) bool

type StreamingResult struct {
	Payload types.DataLoadedPayload // Payload with counts and completion
	Error   error                   // Any error that occurred
}

type LoadOptions struct {
	ForceReload bool
	Filter      FilterFunc[any]
}

type PageResult[T any] struct {
	Items      []T
	TotalItems int
	HasMore    bool
	IsLoaded   bool
}

// BaseRepository provides common repository functionality for type T
type BaseRepository[T any] struct {
	listKind      types.ListKind
	filterFunc    FilterFunc[T]
	processFunc   func(interface{}) *T
	queryFunc     func(*output.RenderCtx)
	dedupeFunc    func(existing []T, newItem *T) bool // Optional deduplication function
	state         *LoadState
	cache         *Cache[T]
	data          []T
	expectedCount int
	loaded        bool
	mutex         sync.RWMutex
}

func NewBaseRepository[T any](
	listKind types.ListKind,
	filterFunc FilterFunc[T],
	processFunc func(interface{}) *T,
	queryFunc func(*output.RenderCtx),
	dedupeFunc func(existing []T, newItem *T) bool, // Optional deduplication function
) *BaseRepository[T] {
	return &BaseRepository[T]{
		listKind:    listKind,
		filterFunc:  filterFunc,
		processFunc: processFunc,
		queryFunc:   queryFunc,
		dedupeFunc:  dedupeFunc,
		state:       NewLoadState(),
		cache:       NewCache[T](),
		data:        make([]T, 0),
	}
}

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

/*
Performance Note:
When filtering is applied, the current implementation creates a new slice
containing copies of all matching items, which can be expensive for large datasets.
An alternative approach would be to use indices or iterators instead:

1. Create a slice of indices that point to matching items in the original data
2. Use these indices to access the original data when returning the page
3. Only copy the specific items needed for the returned page

This would avoid the full copy during filtering while still protecting the original data.
For example:

	matchingIndices := []int{}
	for i, item := range data {
	    if filter(&item) {
	        matchingIndices = append(matchingIndices, i)
	    }
	}

	Then use matchingIndices[first:end] to access only the needed items
	and make copies just of those for the returned page

This approach trades memory efficiency for slightly more complex access patterns.
*/
// GetPage returns a filtered, sorted, and paginated page of data
func (r *BaseRepository[T]) GetPage(first, pageSize int, filter FilterFunc[T], sortSpec interface{}, sortFunc func([]T, interface{}) error) (*PageResult[T], error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	data := r.data
	if filter != nil {
		filtered := make([]T, 0, len(data))
		for _, item := range data {
			if filter(&item) {
				filtered = append(filtered, item)
			}
		}
		data = filtered
	}
	if sortFunc != nil {
		dataCopy := make([]T, len(data))
		copy(dataCopy, data)
		if err := sortFunc(dataCopy, sortSpec); err != nil {
			return nil, fmt.Errorf("error sorting data: %w", err)
		}
		data = dataCopy
	}
	if first < 0 || pageSize <= 0 {
		return nil, fmt.Errorf("invalid pagination parameters")
	}
	end := first + pageSize
	if end > len(data) {
		end = len(data)
	}
	if first >= len(data) {
		return &PageResult[T]{
			Items:      []T{},
			TotalItems: len(data),
			HasMore:    false,
			IsLoaded:   r.state.IsLoaded(),
		}, nil
	}
	return &PageResult[T]{
		Items:      data[first:end],
		TotalItems: len(data),
		HasMore:    end < len(data),
		IsLoaded:   r.state.IsLoaded(),
	}, nil
}

func (r *BaseRepository[T]) IsLoading() bool { return r.state.IsLoading() }
func (r *BaseRepository[T]) IsLoaded() bool  { return r.state.IsLoaded() }
func (r *BaseRepository[T]) Count() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return len(r.data)
}
func (r *BaseRepository[T]) ExpectedCount() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.expectedCount
}
func (r *BaseRepository[T]) NeedsUpdate() bool { return !r.state.IsLoaded() }
func (r *BaseRepository[T]) Clear() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.data = r.data[:0]
	r.loaded = false
	r.state.Reset()
}
