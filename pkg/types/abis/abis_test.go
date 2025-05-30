// ADD_ROUTE
package abis

import (
	"testing"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
)

// MockApp is a mock implementation of the App interface for testing.
type MockApp struct {
	RegisteredCtxs map[base.Address]*output.RenderCtx
	Events         []struct {
		EventType msgs.EventType
		Payload   interface{}
	}
}

func (m *MockApp) LogBackend(msg string) {
}

func (m *MockApp) RegisterCtx(addr base.Address) *output.RenderCtx {
	if m.RegisteredCtxs == nil {
		m.RegisteredCtxs = make(map[base.Address]*output.RenderCtx)
	}
	renderCtx := output.NewStreamingContext()
	m.RegisteredCtxs[addr] = renderCtx
	return renderCtx
}

func (m *MockApp) Cancel(addr base.Address) (int, bool) {
	if m.RegisteredCtxs == nil {
		return 0, false
	}
	if ctx, exists := m.RegisteredCtxs[addr]; exists {
		ctx.Cancel()
		delete(m.RegisteredCtxs, addr)
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

func NewMockApp() *MockApp {
	return &MockApp{
		RegisteredCtxs: make(map[base.Address]*output.RenderCtx),
		Events: make([]struct {
			EventType msgs.EventType
			Payload   interface{}
		}, 0),
	}
}

func TestNewAbisCollection(t *testing.T) {
	mockApp := NewMockApp()
	ac := NewAbisCollection(mockApp)
	if ac.isDownloadedLoaded || ac.isKnownLoaded || ac.isFuncsLoaded || ac.isEventsLoaded {
		t.Error("NewAbisCollection should not have any loaded flags set")
	}
	if ac.isLoading == 1 {
		t.Error("NewAbisCollection should not be loading")
	}
	// Check that slices are initialized (non-nil and empty)
	if ac.downloadedAbis == nil || len(ac.downloadedAbis) != 0 {
		t.Errorf("NewAbisCollection.downloadedAbis not initialized correctly: got %v, want empty non-nil slice", ac.downloadedAbis)
	}
	if ac.knownAbis == nil || len(ac.knownAbis) != 0 {
		t.Errorf("NewAbisCollection.knownAbis not initialized correctly: got %v, want empty non-nil slice", ac.knownAbis)
	}
	if ac.allFunctions == nil || len(ac.allFunctions) != 0 {
		t.Errorf("NewAbisCollection.allFunctions not initialized correctly: got %v", ac.allFunctions)
	}
	if ac.allEvents == nil || len(ac.allEvents) != 0 {
		t.Errorf("NewAbisCollection.allEvents not initialized correctly: got %v", ac.allEvents)
	}
}

// ADD_ROUTE
