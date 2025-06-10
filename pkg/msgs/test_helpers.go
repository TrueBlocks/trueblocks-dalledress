package msgs

import "fmt"

// TestHelpers contains functions that are only useful for testing
// They should not be used in production code
type TestHelpers struct{}

// NewTestHelpers returns a new TestHelpers instance
func NewTestHelpers() *TestHelpers {
	// Enable test mode
	SetTestMode(true)
	return &TestHelpers{}
}

// Cleanup disables test mode and clears all listeners
func (t *TestHelpers) Cleanup() {
	// Clear all listeners
	listenersLock.Lock()
	listeners = make(map[EventType][]func(optionalData ...interface{}))
	listenersLock.Unlock()

	// Disable test mode
	SetTestMode(false)
}

// SimulateLoaded simulates a data loaded event
func (t *TestHelpers) SimulateLoaded(listKind string, currentCount, expectedTotal int) {
	payload := map[string]interface{}{
		"listKind":      listKind,
		"currentCount":  currentCount,
		"expectedTotal": expectedTotal,
	}

	EmitLoaded("test", payload)
}

// SimulateError simulates an error event
func (t *TestHelpers) SimulateError(message string) {
	EmitError(message, fmt.Errorf("test error"))
}

// SimulateStatus simulates a status event
func (t *TestHelpers) SimulateStatus(message string) {
	EmitStatus(message)
}
