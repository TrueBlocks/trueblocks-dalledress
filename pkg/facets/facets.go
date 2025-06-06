package facets

import (
	"sync"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

// Facet defines the contract for data access
// Provides loading, paging, filtering, sorting, and cache management for type T
type Facet[T any] interface {
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

// BaseFacet provides common facet functionality for type T
type BaseFacet[T any] struct {
	listKind      types.ListKind
	filterFunc    FilterFunc[T]
	processFunc   func(interface{}) *T
	queryFunc     func(*output.RenderCtx)
	dedupeFunc    func(existing []T, newItem *T) bool
	state         *LoadState
	cache         *Cache[T]
	data          []T
	expectedCount int
	loaded        bool
	mutex         sync.RWMutex
}

func NewBaseFacet[T any](
	listKind types.ListKind,
	filterFunc FilterFunc[T],
	processFunc func(interface{}) *T,
	queryFunc func(*output.RenderCtx),
	dedupeFunc func(existing []T, newItem *T) bool,
) *BaseFacet[T] {
	return &BaseFacet[T]{
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

func (r *BaseFacet[T]) IsLoading() bool { return r.state.IsLoading() }
func (r *BaseFacet[T]) IsLoaded() bool  { return r.state.IsLoaded() }
func (r *BaseFacet[T]) Count() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return len(r.data)
}
func (r *BaseFacet[T]) ExpectedCount() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.expectedCount
}
func (r *BaseFacet[T]) NeedsUpdate() bool { return !r.state.IsLoaded() }
func (r *BaseFacet[T]) Clear() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.data = r.data[:0]
	r.loaded = false
	r.state.Reset()
}
