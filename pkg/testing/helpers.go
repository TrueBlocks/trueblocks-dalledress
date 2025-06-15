package testing

import (
	"testing"
	"time"
)

func WaitWithTimeout(t *testing.T, condition func() bool, timeout time.Duration, message string) bool {
	t.Helper()

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return true
		}
		time.Sleep(10 * time.Millisecond)
	}

	if message != "" {
		t.Errorf("Timeout waiting for condition: %s", message)
	}
	return false
}

func AssertEventually(t *testing.T, condition func() bool, timeout time.Duration, message string) {
	t.Helper()
	if !WaitWithTimeout(t, condition, timeout, message) {
		t.FailNow()
	}
}

func CreateSafeTimeout(baseTimeout time.Duration) time.Duration {
	return baseTimeout * 2
}

func GetTestTimeout() time.Duration {
	return 5 * time.Second
}
