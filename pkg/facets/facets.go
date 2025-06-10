package facets

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/logging"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/progress"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/store"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

type LoadState string

const (
	StateStale    LoadState = "stale"
	StateFetching LoadState = "fetching"
	StatePartial  LoadState = "partial"
	StateLoaded   LoadState = "loaded"
	StatePending  LoadState = "pending"
	StateError    LoadState = "error"
)

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

type LoadState1 int

const (
	StateStale1 LoadState1 = iota
	StateFetching1
	StateLoaded1
	StateError1
	StatePartial1
)

type FilterFunc[T any] func(item *T) bool

type StreamingResult struct {
	Payload types.DataLoadedPayload
	Error   error
}

type PageResult[T any] struct {
	Items      []T
	TotalItems int
	State      LoadState1
}

type Facet[T any] struct {
	state       atomic.Value
	store       *store.Store[T]
	view        []*T
	expectedCnt int
	listKind    types.ListKind
	filterFunc  FilterFunc[T]
	isDupFunc   func(existing []*T, newItem *T) bool
	mutex       sync.RWMutex
	progress    *progress.Progress
}

// NewFacet creates a new facet that uses a store for data fetching
func NewFacet[T any](
	listKind types.ListKind,
	filterFunc FilterFunc[T],
	isDupFunc func(existing []*T, newItem *T) bool,
	store *store.Store[T],
) *Facet[T] {
	facet := &Facet[T]{
		store:       store,
		view:        make([]*T, 0),
		expectedCnt: 0,
		listKind:    listKind,
		filterFunc:  filterFunc,
		isDupFunc:   isDupFunc,
		progress:    progress.NewProgress(listKind, nil),
	}
	facet.state.Store(StateStale1)
	store.RegisterObserver(facet)

	return facet
}

func (r *Facet[T]) IsLoaded() bool {
	return r.GetState() == StateLoaded1
}

func (r *Facet[T]) IsFetching() bool {
	return r.GetState() == StateFetching1
}

func (r *Facet[T]) GetState() LoadState1 {
	if state := r.state.Load(); state != nil {
		return state.(LoadState1)
	}
	return StateStale1
}

func (r *Facet[T]) NeedsUpdate() bool {
	state := r.GetState()
	return state == StateStale1
}

func (r *Facet[T]) StartFetching() bool {
	currentState := r.GetState()
	if currentState == StateFetching1 {
		return false
	}
	r.state.Store(StateFetching1)
	return true
}

func (r *Facet[T]) SetPartial() {
	if r.GetState() == StateFetching1 {
		r.state.Store(StatePartial1)
	}
}

func (r *Facet[T]) Reset() {
	r.mutex.Lock()
	r.view = r.view[:0]
	r.expectedCnt = 0
	r.state.Store(StateStale1)
	storeToReset := r.store
	r.mutex.Unlock()

	if storeToReset != nil {
		storeToReset.Reset()
	}
}

func (r *Facet[T]) ExpectedCount() int {
	return r.expectedCnt
}

func (r *Facet[T]) Count() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return len(r.view)
}

var ErrAlreadyLoading = errors.New("already loading")

func (r *Facet[T]) Load() (*StreamingResult, error) {
	if !r.NeedsUpdate() {
		cachedPayload := r.getCachedResult()
		msgs.EmitStatus(fmt.Sprintf("cached: %d items", len(r.view)))
		return cachedPayload, nil
	}

	if !r.StartFetching() {
		return nil, ErrAlreadyLoading
	}

	go func() {
		ticker := time.NewTicker(progress.MaxWaitTime / 2)
		defer ticker.Stop()

		done := make(chan error, 1)

		go func() {
			err := r.store.Fetch()
			done <- err
		}()

		for {
			select {
			case err := <-done:
				if err != nil && err != store.ErrStaleFetch() {
					logging.LogBackend(fmt.Sprintf("Failed fetch: %s", err.Error()))
				}
				return

			case <-ticker.C:
				r.mutex.RLock()
				currentCount := len(r.view)
				expectedTotal := r.store.GetExpectedTotal()
				r.mutex.RUnlock()

				r.progress.Heartbeat(currentCount, expectedTotal)
			}
		}
	}()

	return &StreamingResult{
		Payload: types.DataLoadedPayload{
			ListKind: r.listKind,
		},
	}, nil
}

func (r *Facet[T]) getCachedResult() *StreamingResult {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return &StreamingResult{
		Payload: types.DataLoadedPayload{
			CurrentCount:  len(r.view),
			ExpectedTotal: len(r.view),
			ListKind:      r.listKind,
		},
	}
}

func (r *Facet[T]) GetPage(
	first, pageSize int,
	filter FilterFunc[T],
	sortSpec sdk.SortSpec,
	sortFunc func([]T, sdk.SortSpec) error) (*PageResult[T], error,
) {
	r.mutex.RLock()
	data := make([]T, len(r.view))
	for i, ptr := range r.view {
		data[i] = *ptr
	}
	state := r.GetState()
	r.mutex.RUnlock()

	if len(data) == 0 && r.NeedsUpdate() {
		go func() {
			_, _ = r.Load()
		}()
	}

	var filteredData []T
	if filter != nil {
		for i := range data {
			ptr := &data[i]
			if filter(ptr) {
				filteredData = append(filteredData, data[i])
			}
		}
	} else {
		filteredData = data
	}

	if sortFunc != nil && len(filteredData) > 0 {
		if err := sortFunc(filteredData, sortSpec); err != nil {
			return nil, fmt.Errorf("error sorting data: %w", err)
		}
	}

	start := first
	end := first + pageSize
	if start >= len(filteredData) {
		start = 0
		end = 0
	}
	if end > len(filteredData) {
		end = len(filteredData)
	}

	var paginatedData []T
	if start < end {
		paginatedData = filteredData[start:end]
	} else {
		paginatedData = []T{}
	}

	return &PageResult[T]{
		Items:      paginatedData,
		TotalItems: len(filteredData),
		State:      state,
	}, nil
}

func (r *Facet[T]) GetStore() *store.Store[T] {
	return r.store
}

func (r *Facet[T]) SyncWithStore() {
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

func (r *Facet[T]) OnNewItem(item *T, index int) {
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
	expectedTotal := r.store.GetExpectedTotal()
	r.mutex.Unlock()

	payload := r.progress.Tick(currentCount, expectedTotal)
	if payload.CurrentCount > 0 {
		r.SetPartial()
	}
}

func (r *Facet[T]) OnStateChanged(state store.StoreState, reason string) {
	// Map store states to facet states
	switch state {
	case store.StateStale:
		r.state.Store(StateStale1)
		r.expectedCnt = 0
		r.mutex.Lock()
		r.view = r.view[:0]
		r.mutex.Unlock()
		msgs.EmitStatus(fmt.Sprintf("Data outdated: %s", reason))

	case store.StateFetching:
		r.state.Store(StateFetching1)
		r.expectedCnt = 0
		r.mutex.Lock()
		r.view = r.view[:0]
		r.mutex.Unlock()

	case store.StateLoaded:
		r.SyncWithStore()
		r.state.Store(StateLoaded1)
		r.mutex.RLock()
		currentCount := len(r.view)
		r.expectedCnt = r.store.GetExpectedTotal()
		r.mutex.RUnlock()
		_ = r.progress.Tick(currentCount, currentCount)

	case store.StateError:
		r.mutex.RLock()
		hasData := len(r.view) > 0
		currentCount := len(r.view)
		r.mutex.RUnlock()
		if hasData {
			r.state.Store(StatePartial1)
			msgs.EmitStatus(fmt.Sprintf("Partial load: %d items (error: %s)", currentCount, reason))
		} else {
			r.state.Store(StateError1)
			msgs.EmitError(fmt.Sprintf("Load failed: %s", reason), errors.New(reason))
		}

	case store.StateCanceled:
		r.state.Store(StateStale1)
		msgs.EmitStatus("Loading canceled")
	}
}
