package facets

import (
	"sync"
	"sync/atomic"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/sources"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

// Facet defines the contract for data access
// Provides loading, paging, filtering, sorting, and cache management for type T
type Facet[T any] interface {
	Load(opts LoadOptions) (*StreamingResult, error)
	GetPage(first, pageSize int, filter FilterFunc[T], sortSpec interface{}, sortFunc func([]T, interface{}) error) (*PageResult[T], error)
	IsFetching() bool
	IsLoaded() bool
	NeedsUpdate() bool
	Remove(predicate func(*T) bool) bool
	Clear()
	ExpectedCount() int
	Count() int
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

// BaseFacet provides facet functionality using a source for data fetching
type BaseFacet[T any] struct {
	source      sources.Source[T]
	data        []T
	fetching    int32
	loaded      int32
	expectedCnt int
	listKind    types.ListKind
	filterFunc  FilterFunc[T]
	isDupFunc   func(existing []T, newItem *T) bool
	cache       *Cache[T]
	mutex       sync.RWMutex
}

// NewBaseFacet creates a new facet that uses a source for data fetching
func NewBaseFacet[T any](
	listKind types.ListKind,
	filterFunc FilterFunc[T],
	isDupFunc func(existing []T, newItem *T) bool,
	source sources.Source[T],
) *BaseFacet[T] {
	return &BaseFacet[T]{
		source:      source,
		data:        make([]T, 0),
		fetching:    0,
		loaded:      0,
		expectedCnt: 0,
		listKind:    listKind,
		filterFunc:  filterFunc,
		isDupFunc:   isDupFunc,
		cache:       NewCache[T](),
	}
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
	return atomic.CompareAndSwapInt32(&r.fetching, 0, 1)
}

func (r *BaseFacet[T]) StopFetching() {
	atomic.StoreInt32(&r.fetching, 0)
	atomic.StoreInt32(&r.loaded, 1)
}

func (r *BaseFacet[T]) IsFetching() bool {
	return atomic.LoadInt32(&r.fetching) == 1
}

func (r *BaseFacet[T]) NeedsUpdate() bool {
	return !r.IsLoaded()
}

func (r *BaseFacet[T]) IsLoaded() bool {
	return atomic.LoadInt32(&r.loaded) == 1
}

func (r *BaseFacet[T]) Reset() {
	atomic.StoreInt32(&r.fetching, 0)
	atomic.StoreInt32(&r.loaded, 0)
}
