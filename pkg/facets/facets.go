package facets

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

// // Facet defines the contract for data access
// // Provides loading, paging, filtering, sorting, and cache management for type T
// type Facet[T any] interface {
// 	Load() (*StreamingResult, error)
// 	GetPage(first, pageSize int, filter FilterFunc[T], sortSpec interface{}, sortFunc func([]T, interface{}) error) (*PageResult[T], error)
// 	IsFetching() bool
// 	IsLoaded() bool
// 	GetState() LoadState
// 	NeedsUpdate() bool
// 	ForEvery(func(item *T) (error, bool), func(item *T) bool) (int, error)
// 	Reset()
// 	ExpectedCount() int
// 	Count() int
// }

// type FilterFunc[T any] func(*T) bool

// type StreamingResult struct {
// 	Payload types.DataLoadedPayload // Payload with counts and completion
// 	Error   error                   // Any error that occurred
// }

// type PageResult[T any] struct {
// 	Items      []T
// 	TotalItems int
// 	State      LoadState
// }

// // BaseFacet provides facet functionality using a store for data fetching
// type BaseFacet[T any] struct {
// 	state       atomic.Value
// 	store       *store.Store[T]
// 	data        []*T
// 	expectedCnt int
// 	listKind    types.ListKind
// 	filterFunc  FilterFunc[T]
// 	isDupFunc   func(existing []*T, newItem *T) bool // Changed to work with pointers
// 	mutex       sync.RWMutex
// }

// // NewBaseFacet creates a new facet that uses a store for data fetching
// func NewBaseFacet[T any](
// 	listKind types.ListKind,
// 	filterFunc FilterFunc[T],
// 	isDupFunc func(existing []*T, newItem *T) bool, // Changed to work with pointers
// 	store *store.Store[T],
// ) *BaseFacet[T] {
// 	facet := &BaseFacet[T]{
// 		store:       store,
// 		data:        make([]*T, 0),
// 		expectedCnt: 0,
// 		listKind:    listKind,
// 		filterFunc:  filterFunc,
// 		isDupFunc:   isDupFunc,
// 	}
// 	facet.state.Store(StateStale)
// 	return facet
// }

// // MarkStale sets the state to stale, indicating external changes have occurred
// // This can be called by background processes monitoring for data changes
// func (r *BaseFacet[T]) MarkStale() {
// 	r.state.Store(StateStale)
// }

// func (r *BaseFacet[T]) IsLoaded() bool {
// 	return r.GetState() == StateLoaded
// }

// func (r *BaseFacet[T]) IsFetching() bool {
// 	return r.GetState() == StateFetching
// }

// func (r *BaseFacet[T]) GetState() LoadState {
// 	if state := r.state.Load(); state != nil {
// 		return state.(LoadState)
// 	}
// 	return StateStale
// }

// func (r *BaseFacet[T]) NeedsUpdate() bool {
// 	state := r.GetState()
// 	return state == StateStale || state == StateError
// }

// func (r *BaseFacet[T]) StartFetching() bool {
// 	return r.state.CompareAndSwap(StateStale, StateFetching) ||
// 		r.state.CompareAndSwap(StateError, StateFetching)
// }

// func (r *BaseFacet[T]) SetPartial() {
// 	r.state.Store(StatePartial)
// }

// func (r *BaseFacet[T]) Reset() {
// 	r.mutex.Lock()
// 	defer r.mutex.Unlock()
// 	r.data = r.data[:0]
// 	r.expectedCnt = 0
// 	r.state.Store(StateStale)
// }

// func (r *BaseFacet[T]) ExpectedCount() int {
// 	return r.expectedCnt
// }

// func (r *BaseFacet[T]) Count() int {
// 	r.mutex.RLock()
// 	defer r.mutex.RUnlock()
// 	return len(r.data)
// }
