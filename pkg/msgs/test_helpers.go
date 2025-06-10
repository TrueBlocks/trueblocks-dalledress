package msgs

import "fmt"

type TestHelpers struct{}

func NewTestHelpers() *TestHelpers {
	// Enable test mode
	SetTestMode(true)
	return &TestHelpers{}
}

func (t *TestHelpers) Cleanup() {
	listenersLock.Lock()
	listeners = make(map[EventType][]func(optionalData ...interface{}))
	listenersLock.Unlock()

	SetTestMode(false)
}

func (t *TestHelpers) SimulateLoaded(listKind string, currentCount, expectedTotal int) {
	payload := map[string]interface{}{
		"listKind":      listKind,
		"currentCount":  currentCount,
		"expectedTotal": expectedTotal,
	}

	EmitLoaded("test", payload)
}

func (t *TestHelpers) SimulateError(message string) {
	EmitError(message, fmt.Errorf("test error"))
}

func (t *TestHelpers) SimulateStatus(message string) {
	EmitStatus(message)
}
