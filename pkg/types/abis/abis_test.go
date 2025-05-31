// ADD_ROUTE
package abis

import (
	"testing"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
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

func (m *MockApp) LogBackend(msg string) {
	m.LogMessages = append(m.LogMessages, msg)
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

func TestAbisCollection_ClearCache(t *testing.T) {
	mockApp := NewMockApp()
	ac := NewAbisCollection(mockApp)

	// Simulate some initial data
	ac.downloadedAbis = []coreTypes.Abi{{}}
	ac.isDownloadedLoaded = true
	ac.knownAbis = []coreTypes.Abi{{}, {}}
	ac.isKnownLoaded = true
	ac.allFunctions = []coreTypes.Function{{}}
	ac.isFuncsLoaded = true
	ac.allEvents = []coreTypes.Function{{}, {}, {}}
	ac.isEventsLoaded = true

	// Test clearing downloaded ABIs
	ac.ClearCache(AbisDownloaded)
	if ac.isDownloadedLoaded {
		t.Error("ClearCache(AbisDownloaded) should set isDownloadedLoaded to false")
	}
	if ac.expectedDownloaded != 1 {
		t.Errorf("ClearCache(AbisDownloaded) should set expectedDownloaded to 1, got %d", ac.expectedDownloaded)
	}
	if len(ac.downloadedAbis) != 0 {
		t.Errorf("ClearCache(AbisDownloaded) should clear downloadedAbis, got length %d", len(ac.downloadedAbis))
	}

	// Test clearing known ABIs
	ac.ClearCache(AbisKnown)
	if ac.isKnownLoaded {
		t.Error("ClearCache(AbisKnown) should set isKnownLoaded to false")
	}
	if ac.expectedKnown != 2 {
		t.Errorf("ClearCache(AbisKnown) should set expectedKnown to 2, got %d", ac.expectedKnown)
	}
	if len(ac.knownAbis) != 0 {
		t.Errorf("ClearCache(AbisKnown) should clear knownAbis, got length %d", len(ac.knownAbis))
	}

	// Test clearing functions
	ac.ClearCache(AbisFunctions)
	if ac.isFuncsLoaded {
		t.Error("ClearCache(AbisFunctions) should set isFuncsLoaded to false")
	}
	if ac.expectedFunctions != 1 {
		t.Errorf("ClearCache(AbisFunctions) should set expectedFunctions to 1, got %d", ac.expectedFunctions)
	}
	if len(ac.allFunctions) != 0 {
		t.Errorf("ClearCache(AbisFunctions) should clear allFunctions, got length %d", len(ac.allFunctions))
	}

	// Test clearing events
	ac.ClearCache(AbisEvents)
	if ac.isEventsLoaded {
		t.Error("ClearCache(AbisEvents) should set isEventsLoaded to false")
	}
	if ac.expectedEvents != 3 {
		t.Errorf("ClearCache(AbisEvents) should set expectedEvents to 3, got %d", ac.expectedEvents)
	}
	if len(ac.allEvents) != 0 {
		t.Errorf("ClearCache(AbisEvents) should clear allEvents, got length %d", len(ac.allEvents))
	}

	// Verify deduper is cleared
	if len(ac.deduper) != 0 {
		t.Error("ClearCache should clear deduper map")
	}
}

func TestAbisCollection_ClearCache_UnknownListKind(t *testing.T) {
	mockApp := NewMockApp()
	ac := NewAbisCollection(mockApp)

	// Test with unknown ListKind
	unknownKind := types.ListKind("Unknown")
	ac.ClearCache(unknownKind)

	// Should have logged an error
	if len(mockApp.LogMessages) == 0 {
		t.Error("ClearCache with unknown ListKind should log an error message")
	}
}

func TestAbisCollection_NeedsUpdate(t *testing.T) {
	mockApp := NewMockApp()
	ac := NewAbisCollection(mockApp)

	// Initially all should need update (all loaded flags are false)
	if !ac.NeedsUpdate(AbisDownloaded) {
		t.Error("NeedsUpdate(AbisDownloaded) should return true when isDownloadedLoaded is false")
	}
	if !ac.NeedsUpdate(AbisKnown) {
		t.Error("NeedsUpdate(AbisKnown) should return true when isKnownLoaded is false")
	}
	if !ac.NeedsUpdate(AbisFunctions) {
		t.Error("NeedsUpdate(AbisFunctions) should return true when isFuncsLoaded is false")
	}
	if !ac.NeedsUpdate(AbisEvents) {
		t.Error("NeedsUpdate(AbisEvents) should return true when isEventsLoaded is false")
	}

	// Set loaded flags to true
	ac.isDownloadedLoaded = true
	ac.isKnownLoaded = true
	ac.isFuncsLoaded = true
	ac.isEventsLoaded = true

	// Now all should not need update
	if ac.NeedsUpdate(AbisDownloaded) {
		t.Error("NeedsUpdate(AbisDownloaded) should return false when isDownloadedLoaded is true")
	}
	if ac.NeedsUpdate(AbisKnown) {
		t.Error("NeedsUpdate(AbisKnown) should return false when isKnownLoaded is true")
	}
	if ac.NeedsUpdate(AbisFunctions) {
		t.Error("NeedsUpdate(AbisFunctions) should return false when isFuncsLoaded is true")
	}
	if ac.NeedsUpdate(AbisEvents) {
		t.Error("NeedsUpdate(AbisEvents) should return false when isEventsLoaded is true")
	}
}

func TestAbisCollection_NeedsUpdate_UnknownListKind(t *testing.T) {
	mockApp := NewMockApp()
	ac := NewAbisCollection(mockApp)

	// Test with unknown ListKind - should return true
	unknownKind := types.ListKind("Unknown")
	if !ac.NeedsUpdate(unknownKind) {
		t.Error("NeedsUpdate with unknown ListKind should return true")
	}
}

func TestAbisCollection_ThreadSafety(t *testing.T) {
	mockApp := NewMockApp()
	ac := NewAbisCollection(mockApp)

	// Test concurrent access to NeedsUpdate and ClearCache
	done := make(chan bool, 2)

	// Goroutine 1: repeatedly check NeedsUpdate
	go func() {
		for i := 0; i < 100; i++ {
			ac.NeedsUpdate(AbisDownloaded)
			ac.NeedsUpdate(AbisKnown)
		}
		done <- true
	}()

	// Goroutine 2: repeatedly clear cache
	go func() {
		for i := 0; i < 100; i++ {
			ac.ClearCache(AbisDownloaded)
			ac.ClearCache(AbisKnown)
		}
		done <- true
	}()

	// Wait for both goroutines to complete
	<-done
	<-done

	// If we get here without panicking, the mutex is working
	t.Log("Thread safety test completed successfully")
}

// ADD_ROUTE
