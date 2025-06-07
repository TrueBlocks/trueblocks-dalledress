package facets

import (
	"sync"
	"sync/atomic"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/sources"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

// LoadState represents the current state of data loading
type LoadState string

const (
	StateStale    LoadState = "stale"
	StateFetching LoadState = "fetching"
	StatePartial  LoadState = "partial"
	StateLoaded   LoadState = "loaded"
	StatePending  LoadState = "pending"
	StateError    LoadState = "error"
)

// AllStates contains all possible load states for frontend binding
var AllStates = []struct {
	Value  LoadState `json:"value"`
	TSName string    `json:"tsname"`
}{
	{StateStale, "STALE"},
	{StateFetching, "FETCHING"},
	{StatePartial, "PARTIAL"},
	{StateLoaded, "LOADED"},
	{StatePending, "PENDING"},
	{StateError, "ERROR"},
}

// Facet defines the contract for data access
// Provides loading, paging, filtering, sorting, and cache management for type T
type Facet[T any] interface {
	Load() (*StreamingResult, error)
	GetPage(first, pageSize int, filter FilterFunc[T], sortSpec interface{}, sortFunc func([]T, interface{}) error) (*PageResult[T], error)
	IsFetching() bool
	IsLoaded() bool
	GetState() LoadState
	NeedsUpdate() bool
	Remove(predicate func(*T) bool) bool
	Reset()
	ExpectedCount() int
	Count() int
}

type FilterFunc[T any] func(*T) bool

type StreamingResult struct {
	Payload types.DataLoadedPayload // Payload with counts and completion
	Error   error                   // Any error that occurred
}

type PageResult[T any] struct {
	Items      []T
	TotalItems int
	State      LoadState
}

// BaseFacet provides facet functionality using a source for data fetching
type BaseFacet[T any] struct {
	state       atomic.Value
	source      sources.Source[T]
	data        []T
	expectedCnt int
	listKind    types.ListKind
	filterFunc  FilterFunc[T]
	isDupFunc   func(existing []T, newItem *T) bool
	mutex       sync.RWMutex
}

// NewBaseFacet creates a new facet that uses a source for data fetching
func NewBaseFacet[T any](
	listKind types.ListKind,
	filterFunc FilterFunc[T],
	isDupFunc func(existing []T, newItem *T) bool,
	source sources.Source[T],
) *BaseFacet[T] {
	facet := &BaseFacet[T]{
		source:      source,
		data:        make([]T, 0),
		expectedCnt: 0,
		listKind:    listKind,
		filterFunc:  filterFunc,
		isDupFunc:   isDupFunc,
	}
	facet.state.Store(StateStale)
	return facet
}

func (r *BaseFacet[T]) Count() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return len(r.data)
}
func (r *BaseFacet[T]) ExpectedCount() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.expectedCnt
}
func (r *BaseFacet[T]) Clear() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.data = r.data[:0]
	r.Reset()
}

func (r *BaseFacet[T]) StartFetching() bool {
	return r.state.CompareAndSwap(StateStale, StateFetching)
}

func (r *BaseFacet[T]) SetPartial() {
	r.state.CompareAndSwap(StateFetching, StatePartial)
}

func (r *BaseFacet[T]) StopFetching() {
	r.state.Store(StateLoaded)
}

// MarkStale sets the state to stale, indicating external changes have occurred
// This can be called by background processes monitoring for data changes
func (r *BaseFacet[T]) MarkStale() {
	// Only transition to stale if we're in a "complete" state
	// Don't interrupt ongoing fetches - let them complete first

	// Retry loop to handle concurrent state changes
	for {
		current := r.state.Load().(LoadState)

		switch current {
		case StateLoaded, StatePartial:
			// Attempt atomic transition to stale
			if r.state.CompareAndSwap(current, StateStale) {
				return // Success
			}
			// CAS failed due to concurrent change, retry
			continue

		case StateFetching:
			// Don't interrupt ongoing fetch - it will complete and then
			// the next access will check staleness again
			return

		case StateStale, StatePending, StateError:
			// Already in appropriate state or transitional state
			return

		default:
			// Unknown state, don't change it
			return
		}
	}
}

func (r *BaseFacet[T]) IsFetching() bool {
	state := r.state.Load()
	return state == StateFetching || state == StatePartial
}

func (r *BaseFacet[T]) NeedsUpdate() bool {
	state := r.state.Load().(LoadState)
	return state == StateStale
}

func (r *BaseFacet[T]) IsLoaded() bool {
	return r.state.Load() == StateLoaded
}

func (r *BaseFacet[T]) GetState() LoadState {
	return r.state.Load().(LoadState)
}

func (r *BaseFacet[T]) Reset() {
	r.state.Store(StateStale)
}
