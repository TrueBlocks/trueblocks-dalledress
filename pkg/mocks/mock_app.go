package mocks

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
)

// MockApp is a mock implementation of the App interface for testing.
type MockApp struct {
	RegisteredCtxs map[string]*output.RenderCtx
	Events         []struct {
		EventType msgs.EventType
		Payload   interface{}
	}
	LogMessages []string
}

func NewMockApp() *MockApp {
	return &MockApp{
		RegisteredCtxs: make(map[string]*output.RenderCtx),
		Events: make([]struct {
			EventType msgs.EventType
			Payload   interface{}
		}, 0),
		LogMessages: make([]string, 0),
	}
}

func (m *MockApp) RegisterCtx(key string) *output.RenderCtx {
	if m.RegisteredCtxs == nil {
		m.RegisteredCtxs = make(map[string]*output.RenderCtx)
	}
	renderCtx := output.NewStreamingContext()
	m.RegisteredCtxs[key] = renderCtx
	return renderCtx
}

func (m *MockApp) Cancel(key string) (int, bool) {
	if m.RegisteredCtxs == nil {
		return 0, false
	}
	if ctx, exists := m.RegisteredCtxs[key]; exists {
		ctx.Cancel()
		// Safely close channels - check if they're already closed
		select {
		case <-ctx.ModelChan:
			// Channel is already closed
		default:
			close(ctx.ModelChan)
		}
		select {
		case <-ctx.ErrorChan:
			// Channel is already closed
		default:
			close(ctx.ErrorChan)
		}
		delete(m.RegisteredCtxs, key)
		return 1, true
	}
	return 0, false
}

func (m *MockApp) EmitEvent(eventType msgs.EventType, payload interface{}) {
	m.Events = append(m.Events, struct {
		EventType msgs.EventType
		Payload   interface{}
	}{EventType: eventType, Payload: payload})
}

func (m *MockApp) LogBackend(msg string) {
	m.LogMessages = append(m.LogMessages, msg)
}
