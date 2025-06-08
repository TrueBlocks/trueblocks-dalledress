// ADD_ROUTE
package abis

import (
	"testing"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/source"
)

// TestLoadData tests the LoadData function from abis_load.go
func TestLoadData(t *testing.T) {
	// TODO: Turn this back on: Removed t.Parallel() to prevent concurrent SDK access causing race conditions
	ac := NewAbisCollection()

	// Verify initial state - all facets should not be loaded
	if ac.downloadedFacet.IsLoaded() || ac.knownFacet.IsLoaded() || ac.functionsFacet.IsLoaded() || ac.eventsFacet.IsLoaded() {
		t.Error("NewAbisCollection should not have any loaded facets")
	}
	if ac.downloadedFacet.IsFetching() || ac.knownFacet.IsFetching() || ac.functionsFacet.IsFetching() || ac.eventsFacet.IsFetching() {
		t.Error("NewAbisCollection should not have any fetching facets")
	}

	// Call LoadData for known ABIs
	ac.LoadData(AbisKnown)

	// Give the goroutine a moment to start
	time.Sleep(10 * time.Millisecond)

	// After calling LoadData, the known facet should either be fetching or loaded
	// Note: Due to the asynchronous nature, we can't easily test completion without
	// more complex synchronization mechanisms
	isFetchingOrLoaded := ac.knownFacet.IsFetching() || ac.knownFacet.IsLoaded()

	if !isFetchingOrLoaded {
		t.Error("After LoadData, known facet should be fetching or loaded")
	}

	// Test that calling LoadData again doesn't cause issues when already fetching/loaded
	prevState := ac.knownFacet.IsFetching() || ac.knownFacet.IsLoaded()
	ac.LoadData(AbisKnown)

	newState := ac.knownFacet.IsFetching() || ac.knownFacet.IsLoaded()

	if !prevState || !newState {
		t.Error("LoadData should not change state when already fetching/loaded")
	}
}

func TestLoadAbis(t *testing.T) {
	ac := NewAbisCollection()

	// Verify initial state - all facets should not be loaded
	if ac.downloadedFacet.IsLoaded() || ac.knownFacet.IsLoaded() || ac.functionsFacet.IsLoaded() || ac.eventsFacet.IsLoaded() {
		t.Fatalf("NewAbisCollection should not have any loaded facets initially")
	}

	ac.LoadData(AbisKnown)
	// Note: LoadData is async. Testing its completion requires a different approach.
	// For now, we'll assume it kicks off the process.
	// We can check isFetching state.

	// The rest of this test needs significant rework due to the asynchronous nature
	// of loadInternal and the facet pattern replacing direct slice management.
	// We cannot directly check ac.loaded or counts immediately after LoadData.
	// We would need to:
	// 1. Wait for isLoaded to become true (with a timeout).
	// 2. Or, inspect events emitted.
	// For now, this part of the test is effectively disabled or needs to be redesigned.
	t.Skip("TestLoadAbis needs rework for asynchronous loading via LoadData and facet pattern.")
}

// TestReloadCancellation tests that Reload properly cancels ongoing operations
func TestReloadCancellation(t *testing.T) {
	abisAddr := base.ZeroAddr.Hex()
	renderCtx := source.RegisterContext(abisAddr)

	// Verify the context was registered
	if source.CtxCount(abisAddr) != 1 {
		t.Errorf("Expected 1 registered context, got %d", source.CtxCount(abisAddr))
	}

	// Verify the context is not nil
	if renderCtx == nil {
		t.Error("RegisterContext should return non-nil context")
	}

	// Simulate a reload operation by cancelling the context
	cancelled, found := source.UnregisterContext(abisAddr)
	if !found {
		t.Error("Cancel should find the registered context")
	}
	if cancelled != 1 {
		t.Errorf("Expected 1 cancelled context, got %d", cancelled)
	}

	// Verify the context was cancelled and removed
	if source.CtxCount(abisAddr) != 0 {
		t.Errorf("Expected 0 registered contexts after reload, got %d", source.CtxCount(abisAddr))
	}

	// Note: We can't easily test if the context was actually cancelled since
	// the Cancel method removes it from the map, but the fact that it was
	// removed indicates it was processed correctly
}

// TestContextRegistration tests that contexts are properly registered and cleaned up
func TestContextRegistration(t *testing.T) {
	addr1 := "0x1234567890123456789012345678901234567890"
	addr2 := "0x2234567890123456789012345678901234567890"

	ctx1 := source.RegisterContext(addr1)
	ctx2 := source.RegisterContext(addr2)

	cnt := source.CtxCount(addr1) + source.CtxCount(addr1)
	if cnt != 2 {
		t.Errorf("Expected 2 registered contexts, got %d", cnt)
	}

	if ctx1 == nil || ctx2 == nil {
		t.Error("RegisterContext should return non-nil contexts")
	}

	// Test Cancel for specific address
	cancelled, found := source.UnregisterContext(addr1)
	if !found {
		t.Error("Cancel should find the registered context")
	}
	if cancelled != 1 {
		t.Errorf("Expected 1 cancelled context, got %d", cancelled)
	}
	cnt = source.CtxCount(addr1) + source.CtxCount(addr2)
	if cnt != 1 {
		t.Errorf("Expected 1 remaining context after cancel, got %d", cnt)
	}

	// Test Cancel for non-existent address
	nonExistentAddr := "0x9999999999999999999999999999999999999999"
	cancelled, found = source.UnregisterContext(nonExistentAddr)
	if found {
		t.Error("Cancel should not find non-existent context")
	}
	if cancelled != 0 {
		t.Errorf("Expected 0 cancelled contexts for non-existent address, got %d", cancelled)
	}
}

// ADD_ROUTE
