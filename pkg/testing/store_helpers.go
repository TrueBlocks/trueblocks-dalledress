package testing

import (
	"sync"
	"testing"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
)

type TestData struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func (t *TestData) Model(chain, format string, verbose bool, extraOptions map[string]any) types.Model {
	return types.Model{
		Data: map[string]any{
			"id":    t.ID,
			"name":  t.Name,
			"value": t.Value,
		},
		Order: []string{"id", "name", "value"},
	}
}

type MockObserver[T any] struct {
	newItems     []*T
	stateChanges []StateChange
	mutex        sync.Mutex
}

type StateChange struct {
	State  interface{}
	Reason string
}

func NewMockObserver[T any]() *MockObserver[T] {
	return &MockObserver[T]{
		newItems:     make([]*T, 0),
		stateChanges: make([]StateChange, 0),
	}
}

func (m *MockObserver[T]) OnNewItem(item *T, index int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.newItems = append(m.newItems, item)
}

func (m *MockObserver[T]) OnStateChanged(state interface{}, reason string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.stateChanges = append(m.stateChanges, StateChange{
		State:  state,
		Reason: reason,
	})
}

func (m *MockObserver[T]) GetNewItems() []*T {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	result := make([]*T, len(m.newItems))
	copy(result, m.newItems)
	return result
}

func (m *MockObserver[T]) GetStateChanges() []StateChange {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	result := make([]StateChange, len(m.stateChanges))
	copy(result, m.stateChanges)
	return result
}

func (m *MockObserver[T]) Reset() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.newItems = nil
	m.stateChanges = nil
}

func CreateTestDataItems() []*TestData {
	return []*TestData{
		{ID: 1, Name: "Item 1", Value: 100},
		{ID: 2, Name: "Item 2", Value: 200},
		{ID: 3, Name: "Item 3", Value: 300},
	}
}

func CreateStreamFunc(items []*TestData, streamError error) func(*output.RenderCtx) error {
	return func(ctx *output.RenderCtx) error {
		go func() {
			defer close(ctx.ModelChan)
			defer close(ctx.ErrorChan)

			if streamError != nil {
				ctx.ErrorChan <- streamError
				return
			}

			for _, item := range items {
				ctx.ModelChan <- item
			}
		}()
		return nil
	}
}

func CreateTestStore(t *testing.T, name string) interface{} {
	t.Helper()
	return name
}
