package testing

import (
	"time"
)

type MockSDKBase[T any] struct {
	items []T
	count int
	err   error
	delay time.Duration
}

func NewMockSDKBase[T any](items []T) *MockSDKBase[T] {
	return &MockSDKBase[T]{
		items: items,
		count: len(items),
		err:   nil,
		delay: 0,
	}
}

func (m *MockSDKBase[T]) SetError(err error) {
	m.err = err
}

func (m *MockSDKBase[T]) SetDelay(delay time.Duration) {
	m.delay = delay
}

func (m *MockSDKBase[T]) SetItems(items []T) {
	m.items = items
	m.count = len(items)
}

func (m *MockSDKBase[T]) GetItems() []T {
	if m.delay > 0 {
		time.Sleep(m.delay)
	}
	if m.err != nil {
		return nil
	}
	return m.items
}

func (m *MockSDKBase[T]) GetCount() int {
	return m.count
}

func (m *MockSDKBase[T]) GetError() error {
	return m.err
}
