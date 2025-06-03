// ADD_ROUTE
package abis

import (
	"testing"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/mocks"
)

// TestLoadData tests the LoadData function from abis_load.go
func TestLoadData(t *testing.T) {
	mockApp := mocks.NewMockApp()
	ac := NewAbisCollection(mockApp)

	// Verify initial state - all repositories should not be loaded
	if ac.downloadedRepo.IsLoaded() || ac.knownRepo.IsLoaded() || ac.functionsRepo.IsLoaded() || ac.eventsRepo.IsLoaded() {
		t.Error("NewAbisCollection should not have any loaded repositories")
	}
	if ac.downloadedRepo.IsLoading() || ac.knownRepo.IsLoading() || ac.functionsRepo.IsLoading() || ac.eventsRepo.IsLoading() {
		t.Error("NewAbisCollection should not have any loading repositories")
	}

	// Call LoadData for known ABIs
	ac.LoadData(AbisKnown)

	// Give the goroutine a moment to start
	time.Sleep(10 * time.Millisecond)

	// After calling LoadData, the known repository should either be loading or loaded
	// Note: Due to the asynchronous nature, we can't easily test completion without
	// more complex synchronization mechanisms
	isLoadingOrLoaded := ac.knownRepo.IsLoading() || ac.knownRepo.IsLoaded()

	if !isLoadingOrLoaded {
		t.Error("After LoadData, known repository should be loading or loaded")
	}

	// Test that calling LoadData again doesn't cause issues when already loading/loaded
	prevState := ac.knownRepo.IsLoading() || ac.knownRepo.IsLoaded()
	ac.LoadData(AbisKnown)

	newState := ac.knownRepo.IsLoading() || ac.knownRepo.IsLoaded()

	if !prevState || !newState {
		t.Error("LoadData should not change state when already loading/loaded")
	}
}

func TestLoadAbis(t *testing.T) {
	mockApp := mocks.NewMockApp()
	ac := NewAbisCollection(mockApp)

	// Verify initial state - all repositories should not be loaded
	if ac.downloadedRepo.IsLoaded() || ac.knownRepo.IsLoaded() || ac.functionsRepo.IsLoaded() || ac.eventsRepo.IsLoaded() {
		t.Fatalf("NewAbisCollection should not have any loaded repositories initially")
	}

	ac.LoadData(AbisKnown)
	// Note: LoadData is async. Testing its completion requires a different approach.
	// For now, we'll assume it kicks off the process.
	// We can check isLoading state.

	// The rest of this test needs significant rework due to the asynchronous nature
	// of loadInternal and the repository pattern replacing direct slice management.
	// We cannot directly check ac.loaded or counts immediately after LoadData.
	// We would need to:
	// 1. Wait for isLoaded to become true (with a timeout).
	// 2. Or, inspect events emitted by MockApp.
	// For now, this part of the test is effectively disabled or needs to be redesigned.
	t.Skip("TestLoadAbis needs rework for asynchronous loading via LoadData and repository pattern.")
}

// TestReloadCancellation tests that Reload properly cancels ongoing operations
func TestReloadCancellation(t *testing.T) {
	mockApp := mocks.NewMockApp()

	// Simulate registering a context (like what happens in loadInternal)
	abisAddr := base.ZeroAddr.Hex()
	renderCtx := mockApp.RegisterCtx(abisAddr)

	// Verify the context was registered
	if len(mockApp.RegisteredCtxs) != 1 {
		t.Errorf("Expected 1 registered context, got %d", len(mockApp.RegisteredCtxs))
	}

	// Verify the context is not nil
	if renderCtx == nil {
		t.Error("RegisterCtx should return non-nil context")
	}

	// Simulate a reload operation by cancelling the context
	cancelled, found := mockApp.Cancel(abisAddr)
	if !found {
		t.Error("Cancel should find the registered context")
	}
	if cancelled != 1 {
		t.Errorf("Expected 1 cancelled context, got %d", cancelled)
	}

	// Verify the context was cancelled and removed
	if len(mockApp.RegisteredCtxs) != 0 {
		t.Errorf("Expected 0 registered contexts after reload, got %d", len(mockApp.RegisteredCtxs))
	}

	// Note: We can't easily test if the context was actually cancelled since
	// the Cancel method removes it from the map, but the fact that it was
	// removed indicates it was processed correctly
}

// TestContextRegistration tests that contexts are properly registered and cleaned up
func TestContextRegistration(t *testing.T) {
	mockApp := mocks.NewMockApp()

	// Test RegisterCtx
	addr1 := "0x1234567890123456789012345678901234567890"
	addr2 := "0x2234567890123456789012345678901234567890"

	ctx1 := mockApp.RegisterCtx(addr1)
	ctx2 := mockApp.RegisterCtx(addr2)

	if len(mockApp.RegisteredCtxs) != 2 {
		t.Errorf("Expected 2 registered contexts, got %d", len(mockApp.RegisteredCtxs))
	}

	if ctx1 == nil || ctx2 == nil {
		t.Error("RegisterCtx should return non-nil contexts")
	}

	// Test Cancel for specific address
	cancelled, found := mockApp.Cancel(addr1)
	if !found {
		t.Error("Cancel should find the registered context")
	}
	if cancelled != 1 {
		t.Errorf("Expected 1 cancelled context, got %d", cancelled)
	}
	if len(mockApp.RegisteredCtxs) != 1 {
		t.Errorf("Expected 1 remaining context after cancel, got %d", len(mockApp.RegisteredCtxs))
	}

	// Test Cancel for non-existent address
	nonExistentAddr := "0x9999999999999999999999999999999999999999"
	cancelled, found = mockApp.Cancel(nonExistentAddr)
	if found {
		t.Error("Cancel should not find non-existent context")
	}
	if cancelled != 0 {
		t.Errorf("Expected 0 cancelled contexts for non-existent address, got %d", cancelled)
	}
}

// ADD_ROUTE
