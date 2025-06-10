package store

import (
	"errors"
	"sync"
	"sync/atomic"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
)

// StoreState represents the current state of a store
type StoreState int

const (
	StateStale    StoreState = iota // Needs refresh
	StateFetching                   // Currently loading
	StateLoaded                     // Complete data
	StateError                      // Error occurred
	StateCanceled                   // User canceled
)

// FacetObserver defines the interface for objects that observe store changes
type FacetObserver[T any] interface {
	OnNewItem(item *T, index int)
	OnStateChanged(state StoreState, reason string)
}

// Store handle the low-level data fetching and streaming from external systems
type Store[T any] struct {
	data               []T
	observers          []FacetObserver[T]
	queryFunc          func(*output.RenderCtx) error
	processFunc        func(interface{}) *T
	storeType          string
	contextKey         string // Key for ContextManager
	state              StoreState
	stateReason        string
	expectedTotalItems atomic.Int64
	dataGeneration     atomic.Uint64
	mutex              sync.RWMutex
}

var errStaleFetch = errors.New("stale fetch: store was reset")

// ErrStaleFetch returns the error indicating a stale fetch attempt.
func ErrStaleFetch() error {
	return errStaleFetch
}

// NewStore creates a new SDK-based store
func NewStore[T any](
	contextKey string,
	queryFunc func(*output.RenderCtx) error,
	processFunc func(interface{}) *T,
) *Store[T] {
	s := &Store[T]{
		data:        make([]T, 0),
		observers:   make([]FacetObserver[T], 0),
		queryFunc:   queryFunc,
		processFunc: processFunc,
		storeType:   "sdk",
		contextKey:  contextKey,
		state:       StateStale,
	}
	s.expectedTotalItems.Store(0)
	s.dataGeneration.Store(0)
	return s
}

// GetStoreType returns the store's type
func (s *Store[T]) GetStoreType() string {
	return s.storeType
}

// RegisterObserver registers a facet as an observer of this store
func (s *Store[T]) RegisterObserver(observer FacetObserver[T]) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if observer is already registered
	for _, obs := range s.observers {
		if obs == observer {
			return // Already registered
		}
	}

	// Initialize observers slice if nil
	if s.observers == nil {
		s.observers = make([]FacetObserver[T], 0)
	}

	s.observers = append(s.observers, observer)
}

// UnregisterObserver removes a facet as an observer
func (s *Store[T]) UnregisterObserver(observer FacetObserver[T]) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Find and remove the observer
	for i, obs := range s.observers {
		if obs == observer {
			// Remove by swapping with last element and truncating
			lastIndex := len(s.observers) - 1
			s.observers[i] = s.observers[lastIndex]
			s.observers = s.observers[:lastIndex]
			return
		}
	}
}

// ChangeState updates the store state and notifies all observers
// expectedGeneration is used to ensure state changes for Loaded/Error are not for stale fetches.
// For other states, or if expectedGeneration is 0, the check is skipped.
func (s *Store[T]) ChangeState(expectedGeneration uint64, newState StoreState, reason string) {
	s.mutex.Lock()
	if expectedGeneration != 0 && (newState == StateLoaded || newState == StateError) {
		if s.dataGeneration.Load() != expectedGeneration {
			s.mutex.Unlock()
			// Do not change state or notify if the generation is stale for these terminal states.
			// The fetch operation should detect this and return an error or handle it.
			return
		}
	}

	s.state = newState
	s.stateReason = reason

	currentObservers := make([]FacetObserver[T], len(s.observers))
	copy(currentObservers, s.observers)
	stateToSend := s.state
	messageToSend := s.stateReason
	s.mutex.Unlock()

	for _, observer := range currentObservers {
		observer.OnStateChanged(stateToSend, messageToSend)
	}
}

// GetState returns the current state of the store
func (s *Store[T]) GetState() StoreState {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.state
}

// MarkStale marks the store as stale (needs refresh)
// This is typically called externally, not tied to a specific fetch generation.
func (s *Store[T]) MarkStale(reason string) {
	// Pass 0 for expectedGeneration as this is a direct state change not tied to a fetch op generation.
	s.ChangeState(0, StateStale, reason)
}

// GetItem returns an item at the specified index
func (s *Store[T]) GetItem(index int) *T {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if index < 0 || index >= len(s.data) {
		return nil
	}

	return &s.data[index]
}

// GetItems returns all items in the store
func (s *Store[T]) GetItems() []*T {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	result := make([]*T, len(s.data))
	for i := range s.data {
		result[i] = &s.data[i]
	}

	return result
}

// Count returns the number of items in the store
func (s *Store[T]) Count() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.data)
}

// Fetch initiates the data retrieval process from the SDK.
func (s *Store[T]) Fetch() error { // Removed ctx *output.RenderCtx parameter
	fetchGeneration := s.dataGeneration.Load()
	s.ChangeState(fetchGeneration, StateFetching, "Starting data fetch")

	// Register context with ContextManager
	renderCtx := RegisterContext(s.contextKey)
	// Note: Unregistering is handled by Reset or by RegisterContext itself if called again for the same key.

	done := make(chan struct{})
	errChan := make(chan error, 1)
	var processingError error

	go func() {
		defer close(done)
		err := s.queryFunc(renderCtx) // Use renderCtx from ContextManager
		if err != nil {
			errChan <- err
			return
		}
	}()

	modelChanClosed := false
	errorChanClosed := false

	for !modelChanClosed || !errorChanClosed {
		select {
		case itemIntf, ok := <-renderCtx.ModelChan: // Use renderCtx
			if !ok {
				modelChanClosed = true
				if errorChanClosed && processingError == nil {
					s.ChangeState(fetchGeneration, StateLoaded, "Data loaded successfully")
				}
				continue
			}

			itemPtr := s.processFunc(itemIntf)
			if itemPtr == nil {
				continue
			}

			s.mutex.Lock()
			if s.dataGeneration.Load() != fetchGeneration {
				s.mutex.Unlock()
				return errStaleFetch
			}
			var newItem T = *itemPtr
			s.data = append(s.data, newItem)
			s.expectedTotalItems.Store(int64(len(s.data)))
			index := len(s.data) - 1
			currentObservers := make([]FacetObserver[T], len(s.observers))
			copy(currentObservers, s.observers)
			s.mutex.Unlock()

			for _, obs := range currentObservers {
				s.mutex.RLock()
				if index < len(s.data) {
					itemToSend := &s.data[index]
					s.mutex.RUnlock()
					obs.OnNewItem(itemToSend, index)
				} else {
					s.mutex.RUnlock()
				}
			}

		case streamErr, ok := <-renderCtx.ErrorChan: // Use renderCtx
			if !ok {
				errorChanClosed = true
				if modelChanClosed && processingError == nil {
					s.ChangeState(fetchGeneration, StateLoaded, "Data loaded successfully")
				}
				continue
			}
			processingError = streamErr
			s.ChangeState(fetchGeneration, StateError, streamErr.Error())

		case err := <-errChan:
			processingError = err
			s.ChangeState(fetchGeneration, StateError, err.Error())

		case <-done:
			if modelChanClosed && errorChanClosed && processingError == nil {
				s.ChangeState(fetchGeneration, StateLoaded, "Data loaded successfully")
			}
			return processingError

		case <-renderCtx.Ctx.Done(): // Use renderCtx
			s.ChangeState(0, StateCanceled, "User cancelled operation")
			return renderCtx.Ctx.Err()
		}
	}
	return processingError // Should be unreachable if loop condition is correct
}

// AddItem adds an item to the store and notifies observers
// This is primarily used for testing
func (s *Store[T]) AddItem(item T, index int) {
	s.mutex.Lock()
	s.data = append(s.data, item)
	newIndex := len(s.data) - 1
	itemPtr := &s.data[newIndex]

	observers := make([]FacetObserver[T], len(s.observers))
	copy(observers, s.observers)
	s.mutex.Unlock()

	for _, observer := range observers {
		observer.OnNewItem(itemPtr, newIndex)
	}
}

// GetExpectedTotal returns the expected total number of items.
func (s *Store[T]) GetExpectedTotal() int {
	return int(s.expectedTotalItems.Load())
}

// Reset clears the store's data, cancels any ongoing fetch, and sets the state to Stale.
func (s *Store[T]) Reset() {
	s.mutex.Lock()
	// Cancel and unregister the context using ContextManager
	// UnregisterContext cancels and removes. If only cancellation is needed without removal,
	// and RegisterContext handles replacing/cancelling old ones, CancelFetch might be an alternative.
	// Given the name "Reset", UnregisterContext seems appropriate.
	UnregisterContext(s.contextKey)

	s.data = s.data[:0] // Clear data
	s.expectedTotalItems.Store(0)
	s.dataGeneration.Add(1) // Increment generation to invalidate previous fetches
	newState := StateStale
	reason := "Store reset"
	s.mutex.Unlock() // Unlock before calling ChangeState to avoid deadlock if observer calls back into store

	s.ChangeState(0, newState, reason) // Pass 0 for expectedGeneration as this is a reset
}

// ExpectedTotalItems returns the expected total number of items.
func (s *Store[T]) ExpectedTotalItems() int64 {
	return s.expectedTotalItems.Load()
}
