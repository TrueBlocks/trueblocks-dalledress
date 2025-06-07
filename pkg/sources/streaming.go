package sources

import (
	"fmt"
	"strings"
	"time" // Added for time-based updates

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

// ProcessStream streams data from a Source with all the existing streaming features:
// filtering, deduplication, progress updates, and real-time messaging.
func ProcessStream[T any](
	contextKey string,
	source Source[T],
	filterFunc func(item *T) bool,
	isDupFunc func(existing []T, newItem *T) bool, // Returns true if item should be added (not a duplicate)
	targetSlice *[]T,
	expectedCnt *int,
	listKind types.ListKind,
	m interface {
		Lock()
		Unlock()
		RLock()
		RUnlock()
	},
	onFirstDataCallback func(), // Renamed to avoid confusion, passed to Progress
) (types.DataLoadedPayload, error) {
	renderCtx := RegisterContext(contextKey)
	defer UnregisterContext(contextKey)

	done := make(chan struct{})
	streamErr := make(chan error, 1)

	reporter := NewProgress(listKind, onFirstDataCallback)
	go func() {
		defer close(done)
		ticker := time.NewTicker(maxWaitTime / 2) // expected time 1/2 of maxWaitTime
		defer ticker.Stop()

		fetchDone := make(chan error)
		go func() {
			err := source.Fetch(renderCtx, func(itemPtr *T) bool {
				if filterFunc(itemPtr) {
					var itemAdded bool
					var currentTotalCount int
					var currentExpectedCount int // To store dereferenced value of expectedCnt

					m.Lock()
					// Check for duplicates and add item if not a duplicate
					if isDupFunc == nil || !isDupFunc(*targetSlice, itemPtr) {
						*targetSlice = append(*targetSlice, *itemPtr)
						itemAdded = true
					}
					currentTotalCount = len(*targetSlice)
					if expectedCnt != nil {
						currentExpectedCount = *expectedCnt // Dereference under lock
					}
					m.Unlock()

					if !itemAdded {
						return true // Continue if item was a duplicate or filtered out by isDupFunc
					}

					reporter.Tick(currentTotalCount, currentExpectedCount)
				}
				return true
			})
			fetchDone <- err
		}()

		// Main loop for this goroutine: listen to fetch completion or ticker
		for {
			select {
			case err := <-fetchDone:
				if err != nil {
					streamErr <- err
				}
				return // Exit the select loop and the goroutine
			case <-ticker.C:
				var currentTotalCount int
				var currentExpectedCount int

				m.RLock() // Use RLock for reading counts
				currentTotalCount = len(*targetSlice)
				if expectedCnt != nil {
					currentExpectedCount = *expectedCnt
				}
				m.RUnlock()

				// Check for heartbeat update using the reporter
				reporter.HeartbeatUpdate(currentTotalCount, currentExpectedCount)
			}
		}
	}()

	var streamingError error
	select {
	case <-done:
		// Streaming completed successfully
	case err := <-streamErr:
		streamingError = err
		// Don't return immediately - we want to preserve partial data
	}

	m.RLock() // Use RLock for final count read
	itemCount := len(*targetSlice)
	m.RUnlock()

	if itemCount > 0 {
		payload := reporter.Tick(itemCount, itemCount)
		if streamingError != nil && itemCount > 0 {
			return payload, nil
		}
		return payload, streamingError
	}

	// Only return error if we have no data at all and an error occurred
	if streamingError != nil {
		return types.DataLoadedPayload{}, streamingError
	}

	payload := reporter.Tick(itemCount, itemCount)
	return payload, nil
}

func progress(cnt int, kind types.ListKind, heartbeat bool) string {
	k := strings.Trim(strings.ToLower(string(kind)), " ")
	if heartbeat {
		return fmt.Sprintf("Loaded %d %s...", cnt, k)
	}
	return fmt.Sprintf("Loaded %d %s", cnt, k)
}

const (
	firstPageCount   = 7
	initialIncrement = 10
	incrementGrowth  = 10
	maxWaitTime      = 125 * time.Millisecond
)

// Progress manages the logic for sending progress updates.
type Progress struct {
	lastUpdate        time.Time
	nItemsSinceUpdate int
	nextThreshold     int
	currentIncrement  int
	listKind          types.ListKind
	onFirstDataFunc   func()
	firstDataSent     bool
}

// NewProgress creates and initializes a Progress.
func NewProgress(
	listKindCfg types.ListKind,
	onFirstDataCallback func(), // Can be nil
) *Progress {
	pr := &Progress{
		listKind:        listKindCfg,
		onFirstDataFunc: onFirstDataCallback,
	}
	// Initialize internal state
	pr.lastUpdate = time.Now()
	pr.nItemsSinceUpdate = 0
	pr.nextThreshold = firstPageCount + initialIncrement
	pr.currentIncrement = initialIncrement
	pr.firstDataSent = false
	return pr
}

func (pr *Progress) Tick(currentTotalCount, expectedTotal int) types.DataLoadedPayload {
	pr.nItemsSinceUpdate++
	shouldUpdate := false

	if !pr.firstDataSent && currentTotalCount == firstPageCount {
		shouldUpdate = true
		if pr.onFirstDataFunc != nil {
			go pr.onFirstDataFunc()
		}
		pr.firstDataSent = true
	} else if currentTotalCount >= pr.nextThreshold && currentTotalCount > firstPageCount {
		shouldUpdate = true
		pr.currentIncrement += incrementGrowth
		pr.nextThreshold = currentTotalCount + pr.currentIncrement
	}

	payload := types.DataLoadedPayload{
		CurrentCount:  currentTotalCount,
		ExpectedTotal: expectedTotal,
		ListKind:      pr.listKind,
	}
	if shouldUpdate {
		msgs.EmitLoaded("streaming", payload)
		msgs.EmitStatus(progress(currentTotalCount, pr.listKind, false))
		pr.nItemsSinceUpdate = 0
		pr.lastUpdate = time.Now()
	}

	return payload
}

func (pr *Progress) HeartbeatUpdate(currentTotalCount, expectedTotal int) types.DataLoadedPayload {
	payload := types.DataLoadedPayload{
		CurrentCount:  currentTotalCount,
		ExpectedTotal: expectedTotal,
		ListKind:      pr.listKind,
	}

	if time.Since(pr.lastUpdate) >= maxWaitTime && pr.nItemsSinceUpdate > 0 {
		msgs.EmitLoaded("partial", payload)
		msgs.EmitStatus(progress(currentTotalCount, pr.listKind, true))

		pr.nItemsSinceUpdate = 0
		pr.lastUpdate = time.Now()
	}

	return payload
}
