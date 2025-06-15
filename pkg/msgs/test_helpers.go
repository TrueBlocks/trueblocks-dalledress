package msgs

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

func SetTestMode(enabled bool) {
	testModeLock.Lock()
	defer testModeLock.Unlock()
	testMode = enabled
}

func IsTestMode() bool {
	testModeLock.RLock()
	defer testModeLock.RUnlock()
	return testMode
}

// func (t *TestHelpers) SimulateLoaded(listKind string, currentCount, expectedTotal int) {
// 	payload := map[string]interface{}{
// 		"listKind":      listKind,
// 		"currentCount":  currentCount,
// 		"expectedTotal": expectedTotal,
// 	}

// 	EmitLoaded("test", payload)
// }

// func (t *TestHelpers) SimulateError(message string) {
// 	EmitError(message, fmt.Errorf("test error"))
// }

// func (t *TestHelpers) SimulateStatus(message string) {
// 	EmitStatus(message)
// }
