package repository

import (
	"errors"
	"fmt"
	"sync"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/streaming"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

// Repository defines the contract for data access
type Repository[T any] interface {
	// Load data with streaming progress updates
	Load(opts LoadOptions) (*StreamingResult, error)

	// Get paginated data with sorting (assumes data is loaded)
	GetPage(first, pageSize int, filter FilterFunc[T], sortSpec interface{}, sortFunc func([]T, interface{}) error) (*PageResult[T], error)

	// State queries
	IsLoading() bool
	IsLoaded() bool
	NeedsUpdate() bool

	// Cache management
	Clear()
	Remove(predicate func(*T) bool) bool

	// Metadata
	Count() int
	ExpectedCount() int
}

// FilterFunc defines a filter function for type T
type FilterFunc[T any] func(*T) bool

// StreamingResult provides access to completed streaming data
type StreamingResult struct {
	Status  string                  // Status message from streaming
	Payload types.DataLoadedPayload // Payload with counts and completion
	Error   error                   // Any error that occurred
}

// LoadOptions configures how data is loaded
type LoadOptions struct {
	ForceReload bool
	Filter      FilterFunc[any]
}

// PageResult contains paginated data
type PageResult[T any] struct {
	Items      []T
	TotalItems int
	HasMore    bool
	IsLoaded   bool
}

// BaseRepository provides common repository functionality
type BaseRepository[T any] struct {
	listKind    types.ListKind
	filterFunc  FilterFunc[T]
	processFunc func(interface{}) *T
	queryFunc   func(*output.RenderCtx)
	dedupeFunc  func(existing []T, newItem *T) bool // Optional deduplication function

	// State management
	state *LoadState
	cache *Cache[T]

	// Data storage
	data          []T
	expectedCount int
	loaded        bool
	mutex         sync.RWMutex
}

// NewBaseRepository creates a new base repository
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

// Load implements Repository.Load using your existing streaming system
func (r *BaseRepository[T]) Load(opts LoadOptions) (*StreamingResult, error) {
	if !r.state.StartLoading() {
		return nil, ErrorAlreadyLoading
	}
	defer r.state.StopLoading()

	if !opts.ForceReload && r.state.IsLoaded() {
		return r.getCachedResult(), nil
	}

	r.mutex.Lock()
	r.data = r.data[:0]
	r.loaded = false
	r.mutex.Unlock()

	contextKey := fmt.Sprintf("repo-%s", r.listKind)
	finalStatus, finalPayload, err := streaming.LoadStreamingData(
		contextKey,       // ← Escape key cancels using this key
		r.queryFunc,      // ← Uses your rendering context
		r.filterFunc,     // ← Your existing filter function
		r.processFunc,    // ← Process each item
		r.dedupeFunc,     // ← Optional deduplication function
		&r.data,          // ← Collect results
		&r.expectedCount, // ← Track expected count
		&r.loaded,        // ← Track completion
		r.listKind,       // ← For progress messages
		&r.mutex,         // ← Thread safety
	)

	return &StreamingResult{
		Status:  finalStatus,
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
		Status: fmt.Sprintf("cached: %d items", len(r.data)),
		Payload: types.DataLoadedPayload{
			CurrentCount:  len(r.data),
			ExpectedTotal: len(r.data),
			IsLoaded:      true,
			ListKind:      r.listKind,
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
func (r *BaseRepository[T]) GetPage(first, pageSize int, filter FilterFunc[T], sortSpec interface{}, sortFunc func([]T, interface{}) error) (*PageResult[T], error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	data := r.data

	// Apply additional filter if provided
	if filter != nil {
		filtered := make([]T, 0, len(data))
		for _, item := range data {
			if filter(&item) {
				filtered = append(filtered, item)
			}
		}
		data = filtered
	}

	// Apply sorting if sortFunc is provided
	if sortFunc != nil {
		// Make a copy of data to avoid modifying the original slice during sorting
		dataCopy := make([]T, len(data))
		copy(dataCopy, data)

		if err := sortFunc(dataCopy, sortSpec); err != nil {
			return nil, fmt.Errorf("error sorting data: %w", err)
		}
		data = dataCopy
	}

	// Validate pagination params
	if first < 0 || pageSize <= 0 {
		return nil, fmt.Errorf("invalid pagination parameters")
	}

	// Calculate pagination
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
func (r *BaseRepository[T]) NeedsUpdate() bool {
	return !r.state.IsLoaded()
}
func (r *BaseRepository[T]) Clear() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.data = r.data[:0]
	r.loaded = false
	r.state.Reset()
}
