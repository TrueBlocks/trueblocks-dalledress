package enhancedfacet

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/enhancedstore"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

// LoadState represents the state of a facet
type LoadState int

const (
	StateStale LoadState = iota
	StateFetching
	StateLoaded
	StateError
	StatePartial
)

// FilterFunc defines a function that filters items
type FilterFunc[T any] func(item *T) bool

// StreamingResult is the result of a Load operation
type StreamingResult struct {
	Payload types.DataLoadedPayload
	Error   error
}

// PageResult is the result of a GetPage operation
type PageResult[T any] struct {
	Items      []T
	TotalItems int
	State      LoadState
}

// ErrAlreadyLoading is returned when a Load operation is already in progress
var ErrAlreadyLoading = errors.New("already loading")

// BaseFacet provides facet functionality using a store for data fetching
type BaseFacet[T any] struct {
	state       atomic.Value
	store       *enhancedstore.Store[T]
	view        []*T // Direct pointers to store data
	expectedCnt int
	listKind    types.ListKind
	filterFunc  FilterFunc[T]
	isDupFunc   func(existing []*T, newItem *T) bool // Works with pointers
	mutex       sync.RWMutex
}

// NewBaseFacet creates a new facet that uses a store for data fetching
func NewBaseFacet[T any](
	listKind types.ListKind,
	filterFunc FilterFunc[T],
	isDupFunc func(existing []*T, newItem *T) bool,
	store *enhancedstore.Store[T],
) *BaseFacet[T] {
	facet := &BaseFacet[T]{
		store:       store,
		view:        make([]*T, 0),
		expectedCnt: 0,
		listKind:    listKind,
		filterFunc:  filterFunc,
		isDupFunc:   isDupFunc,
	}
	facet.state.Store(StateStale)

	// Register as observer with the store
	store.RegisterObserver(facet)

	return facet
}

// OnNewItem is called when a new item is added to the store
// This implements the enhancedstore.FacetObserver interface
func (r *BaseFacet[T]) OnNewItem(item *T, index int) {
	// Apply filter
	if r.filterFunc != nil && !r.filterFunc(item) {
		return
	}

	// Check for duplicates if needed
	if r.isDupFunc != nil {
		r.mutex.RLock()
		isDuplicate := r.isDupFunc(r.view, item)
		r.mutex.RUnlock()

		if isDuplicate {
			return
		}
	}

	// Add to view
	r.mutex.Lock()
	r.view = append(r.view, item)
	currentCount := len(r.view)
	r.mutex.Unlock()

	// Emit progress updates periodically
	if currentCount%10 == 0 || currentCount <= 10 {
		r.SetPartial()
		msgs.EmitLoaded("streaming", types.DataLoadedPayload{
			CurrentCount:  currentCount,
			ExpectedTotal: r.store.GetExpectedTotal(), // Use store's expected total
			ListKind:      r.listKind,
		})
		msgs.EmitStatus(fmt.Sprintf("Loaded %d items...", currentCount))
	}
}

// OnStateChanged is called when the store state changes
// This implements the enhancedstore.FacetObserver interface
func (r *BaseFacet[T]) OnStateChanged(state enhancedstore.StoreState, reason string) {
	// Map store states to facet states
	switch state {
	case enhancedstore.StateStale:
		r.state.Store(StateStale)
		r.expectedCnt = 0 // Reset expected count on stale
		// Clear the view as the data is stale
		r.mutex.Lock()
		r.view = r.view[:0]
		r.mutex.Unlock()
		msgs.EmitStatus(fmt.Sprintf("Data outdated: %s", reason))

	case enhancedstore.StateFetching:
		r.state.Store(StateFetching)
		r.expectedCnt = 0 // Reset expected count on fetching
		r.mutex.Lock()
		r.view = r.view[:0] // Clear existing view
		r.mutex.Unlock()

	case enhancedstore.StateLoaded:
		r.SyncWithStore() // Ensure view is up-to-date with the store
		r.state.Store(StateLoaded)
		r.mutex.RLock()
		currentCount := len(r.view)
		r.expectedCnt = r.store.GetExpectedTotal() // Update expected count from store
		r.mutex.RUnlock()
		msgs.EmitLoaded("loaded", types.DataLoadedPayload{ // Emit final loaded event
			CurrentCount:  currentCount,
			ExpectedTotal: r.expectedCnt, // Use the facet's updated expectedCnt
			ListKind:      r.listKind,
		})
		msgs.EmitStatus(fmt.Sprintf("Loaded %d items", currentCount))

	case enhancedstore.StateError:
		// Handle error state based on whether we have partial data
		r.mutex.RLock()
		hasData := len(r.view) > 0
		currentCount := len(r.view)
		r.mutex.RUnlock()

		if hasData {
			r.state.Store(StatePartial)
			msgs.EmitStatus(fmt.Sprintf("Partial load: %d items (error: %s)", currentCount, reason))
		} else {
			r.state.Store(StateError)
			msgs.EmitError(fmt.Sprintf("Load failed: %s", reason), errors.New(reason))
		}

	case enhancedstore.StateCanceled:
		r.state.Store(StateStale) // Cancelled fetches are considered stale
		msgs.EmitStatus("Loading canceled")
	}
}

// SetPartial sets the state to partial if it's fetching
func (r *BaseFacet[T]) SetPartial() {
	if r.state.Load() == StateFetching {
		r.state.Store(StatePartial)
	}
}

// GetState returns the current state of the facet
func (r *BaseFacet[T]) GetState() LoadState {
	if state := r.state.Load(); state != nil {
		return state.(LoadState)
	}
	return StateStale
}

// IsFetching returns true if the facet is currently fetching data
func (r *BaseFacet[T]) IsFetching() bool {
	return r.GetState() == StateFetching
}

// IsLoaded returns true if the facet has loaded data
func (r *BaseFacet[T]) IsLoaded() bool {
	return r.GetState() == StateLoaded
}

// ExpectedCount returns the expected number of items
func (r *BaseFacet[T]) ExpectedCount() int {
	// This should reflect the facet's understanding, updated from store on StateLoaded
	return r.expectedCnt
}

// Count returns the current number of items in the view
func (r *BaseFacet[T]) Count() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return len(r.view)
}

// Reset clears the facet data and sets state to stale
func (r *BaseFacet[T]) Reset() {
	r.mutex.Lock()
	r.view = r.view[:0]
	r.expectedCnt = 0
	r.state.Store(StateStale)
	storeToReset := r.store
	r.mutex.Unlock()

	if storeToReset != nil {
		storeToReset.Reset() // This will trigger ContextManager.UnregisterContext
	}
}

// Load implements loading using a store instead of a queryFunc
// It now accepts variadic queryArgs to be passed to the underlying store's Fetch method.
func (r *BaseFacet[T]) Load(queryArgs ...interface{}) (*StreamingResult, error) {
	if !r.NeedsUpdate() {
		msgs.EmitStatus(fmt.Sprintf("cached: %d items", len(r.view)))
		cachedPayload := r.getCachedResult()
		return cachedPayload, nil
	}

	if !r.StartFetching() {
		return nil, ErrAlreadyLoading
	}

	go func() {
		// Pass queryArgs to the store's Fetch method
		// The store will obtain its RenderCtx from the ContextManager.
		err := r.store.Fetch(queryArgs...)
		if err != nil && err != enhancedstore.ErrStaleFetch() { // Check for actual errors, not just stale fetch
			// If Fetch returns an error, the store should have already set its state to StateError.
			// The facet's OnStateChanged will handle updating the facet's state.
			// We might log it here or let OnStateChanged handle all user-facing messages.
			// For now, rely on OnStateChanged.
		}
		// If err is nil or errStaleFetch, the store has handled its state (Loaded or Stale/Canceled).
		// OnStateChanged will reflect this in the facet.
	}()

	return &StreamingResult{
		Payload: types.DataLoadedPayload{
			CurrentCount:  0,
			ExpectedTotal: 0, // This will be updated by OnStateChanged
			ListKind:      r.listKind,
		},
	}, nil
}

// GetPage returns a page of items from the facet
// It implements the decoupled design by returning immediately with current data
func (r *BaseFacet[T]) GetPage(first, pageSize int, filter FilterFunc[T],
	sortSpec interface{}, sortFunc func([]T, interface{}) error, queryArgs ...interface{}) (*PageResult[T], error) { // Added queryArgs

	// Just return whatever data we have RIGHT NOW - never wait for loading
	r.mutex.RLock()

	// Convert view pointers to values for pagination logic
	data := make([]T, len(r.view))
	for i, ptr := range r.view {
		data[i] = *ptr
	}
	state := r.state.Load().(LoadState)
	r.mutex.RUnlock() // Unlock before potentially long operations like filtering/sorting or Load

	// If we have no data but should have some, trigger background fetch
	if len(data) == 0 && r.NeedsUpdate() {
		// Start async load WITHOUT waiting for it
		go func() {
			// Pass queryArgs to Load if this facet type requires them for its store.
			// This is crucial for facets that use stores needing specific arguments (e.g., address for detail store).
			r.Load(queryArgs...) // Pass queryArgs here
		}()
	}

	// Apply filtering
	var filteredData []T
	if filter != nil {
		for i := range data { // data is a copy, so this is safe
			ptr := &data[i] // filterFunc expects *T
			if filter(ptr) {
				filteredData = append(filteredData, data[i])
			}
		}
	} else {
		filteredData = data
	}

	// Apply sorting if needed
	if sortFunc != nil && len(filteredData) > 0 {
		sortFunc(filteredData, sortSpec)
	}

	// Apply pagination
	start := first
	end := first + pageSize
	if start >= len(filteredData) {
		start = 0 // Reset to 0 if out of bounds, effectively showing no items or first page
		end = 0
	}
	if end > len(filteredData) {
		end = len(filteredData)
	}

	var paginatedData []T
	if start < end { // Ensure start is less than end to avoid panic
		paginatedData = filteredData[start:end]
	} else {
		paginatedData = []T{} // Return empty slice if pagination is invalid or no items
	}

	return &PageResult[T]{
		Items:      paginatedData,
		TotalItems: len(filteredData),
		State:      state,
	}, nil
}

// NeedsUpdate returns true if the facet needs to be updated
func (r *BaseFacet[T]) NeedsUpdate() bool {
	state := r.state.Load().(LoadState)
	return state == StateStale
}

// StartFetching tries to set the state to fetching
func (r *BaseFacet[T]) StartFetching() bool {
	currentState := r.state.Load().(LoadState)
	if currentState == StateFetching {
		return false
	}
	r.state.Store(StateFetching)
	return true
}

// getCachedResult returns a result with cached data
func (r *BaseFacet[T]) getCachedResult() *StreamingResult {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return &StreamingResult{
		Payload: types.DataLoadedPayload{
			CurrentCount:  len(r.view),
			ExpectedTotal: len(r.view), // For cached, expected is what we have.
			ListKind:      r.listKind,
		},
	}
}

// GetStore returns the underlying store
// This is used for testing and direct access when needed
func (r *BaseFacet[T]) GetStore() *enhancedstore.Store[T] {
	return r.store
}

// SyncWithStore ensures the facet's view is in sync with the store's data
func (r *BaseFacet[T]) SyncWithStore() {
	store := r.store
	if store == nil {
		return
	}

	storeItems := store.GetItems()

	r.mutex.Lock()
	r.view = make([]*T, 0, len(storeItems))
	for i := range storeItems {
		itemPtr := storeItems[i]
		if r.filterFunc == nil || r.filterFunc(itemPtr) {
			if r.isDupFunc == nil || !r.isDupFunc(r.view, itemPtr) {
				r.view = append(r.view, itemPtr)
			}
		}
	}
	r.mutex.Unlock()
}
