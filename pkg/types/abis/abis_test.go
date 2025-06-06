// ADD_ROUTE
package abis

import (
	"testing"

	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

func TestNewAbisCollection(t *testing.T) {
	ac := NewAbisCollection()

	// Check that all facets are properly initialized and not loaded
	if ac.downloadedFacet == nil {
		t.Error("NewAbisCollection should initialize downloadedFacet")
	}
	if ac.downloadedFacet.IsLoaded() {
		t.Error("NewAbisCollection downloadedFacet should not be loaded initially")
	}
	if ac.downloadedFacet.IsLoading() {
		t.Error("NewAbisCollection downloadedFacet should not be loading initially")
	}

	if ac.knownFacet == nil {
		t.Error("NewAbisCollection should initialize knownFacet")
	}
	if ac.knownFacet.IsLoaded() {
		t.Error("NewAbisCollection knownFacet should not be loaded initially")
	}
	if ac.knownFacet.IsLoading() {
		t.Error("NewAbisCollection knownFacet should not be loading initially")
	}

	if ac.functionsFacet == nil {
		t.Error("NewAbisCollection should initialize functionsFacet")
	}
	if ac.functionsFacet.IsLoaded() {
		t.Error("NewAbisCollection functionsFacet should not be loaded initially")
	}
	if ac.functionsFacet.IsLoading() {
		t.Error("NewAbisCollection functionsFacet should not be loading initially")
	}

	if ac.eventsFacet == nil {
		t.Error("NewAbisCollection should initialize eventsFacet")
	}
	if ac.eventsFacet.IsLoaded() {
		t.Error("NewAbisCollection eventsFacet should not be loaded initially")
	}
	if ac.eventsFacet.IsLoading() {
		t.Error("NewAbisCollection eventsFacet should not be loading initially")
	}
}

func TestAbisCollection_ClearCache(t *testing.T) {
	ac := NewAbisCollection()

	// For all data types now using Facet pattern
	// We can't easily simulate initial data without loading, so we test the clear operation
	ac.ClearCache(AbisDownloaded)
	if ac.downloadedFacet.IsLoaded() {
		t.Error("ClearCache(AbisDownloaded) should clear the facet loaded state")
	}

	ac.ClearCache(AbisKnown)
	if ac.knownFacet.IsLoaded() {
		t.Error("ClearCache(AbisKnown) should clear the facet loaded state")
	}

	ac.ClearCache(AbisFunctions)
	if ac.functionsFacet.IsLoaded() {
		t.Error("ClearCache(AbisFunctions) should clear the facet loaded state")
	}

	ac.ClearCache(AbisEvents)
	if ac.eventsFacet.IsLoaded() {
		t.Error("ClearCache(AbisEvents) should clear the facet loaded state")
	}
}

func TestAbisCollection_NeedsUpdate(t *testing.T) {
	ac := NewAbisCollection()

	// Initially all should need update (facets not loaded)
	if !ac.NeedsUpdate(AbisDownloaded) {
		t.Error("NeedsUpdate(AbisDownloaded) should return true when facet is not loaded")
	}
	if !ac.NeedsUpdate(AbisKnown) {
		t.Error("NeedsUpdate(AbisKnown) should return true when facet is not loaded")
	}
	if !ac.NeedsUpdate(AbisFunctions) {
		t.Error("NeedsUpdate(AbisFunctions) should return true when facet is not loaded")
	}
	if !ac.NeedsUpdate(AbisEvents) {
		t.Error("NeedsUpdate(AbisEvents) should return true when facet is not loaded")
	}

	// Note: We can't easily simulate loaded facets in tests without actual data loading
	// so this test primarily verifies the initial state behavior
}

func TestAbisCollection_NeedsUpdate_UnknownListKind(t *testing.T) {
	ac := NewAbisCollection()

	// Test with unknown ListKind - should return true
	unknownKind := types.ListKind("Unknown")
	if !ac.NeedsUpdate(unknownKind) {
		t.Error("NeedsUpdate with unknown ListKind should return true")
	}
}

func TestAbisCollection_ThreadSafety(t *testing.T) {
	ac := NewAbisCollection()

	// Test concurrent access to NeedsUpdate and ClearCache across all slice types
	done := make(chan bool, 4)

	// Goroutine 1: repeatedly check NeedsUpdate for downloaded and known
	go func() {
		for i := 0; i < 100; i++ {
			ac.NeedsUpdate(AbisDownloaded)
			ac.NeedsUpdate(AbisKnown)
		}
		done <- true
	}()

	// Goroutine 2: repeatedly clear cache for downloaded and known
	go func() {
		for i := 0; i < 100; i++ {
			ac.ClearCache(AbisDownloaded)
			ac.ClearCache(AbisKnown)
		}
		done <- true
	}()

	// Goroutine 3: repeatedly check NeedsUpdate for functions and events
	go func() {
		for i := 0; i < 100; i++ {
			ac.NeedsUpdate(AbisFunctions)
			ac.NeedsUpdate(AbisEvents)
		}
		done <- true
	}()

	// Goroutine 4: repeatedly clear cache for functions and events
	go func() {
		for i := 0; i < 100; i++ {
			ac.ClearCache(AbisFunctions)
			ac.ClearCache(AbisEvents)
		}
		done <- true
	}()

	// Wait for all goroutines to complete
	<-done
	<-done
	<-done
	<-done

	// If we get here without panicking, the per-slice mutexes are working
	t.Log("Thread safety test completed successfully")
}

func TestAbisFacet(t *testing.T) {
	downloadedFacet := NewAbisFacet("Downloaded", func(abi *coreTypes.Abi) bool {
		return !abi.IsKnown
	})
	knownFacet := NewAbisFacet("Known", func(abi *coreTypes.Abi) bool {
		return abi.IsKnown
	})
	functionsFacet := NewFunctionsFacet("Functions", func(item *coreTypes.Function) bool {
		return item.FunctionType != "event"
	})
	eventsFacet := NewFunctionsFacet("Events", func(item *coreTypes.Function) bool {
		return item.FunctionType == "event"
	})

	if downloadedFacet == nil {
		t.Error("Downloaded facet should not be nil")
	}

	if knownFacet == nil {
		t.Error("Known facet should not be nil")
	}

	if functionsFacet == nil {
		t.Error("Functions facet should not be nil")
	}

	if eventsFacet == nil {
		t.Error("Events facet should not be nil")
	}

	// Test initial states
	facets := []interface {
		IsLoaded() bool
		NeedsUpdate() bool
		Count() int
	}{downloadedFacet, knownFacet, functionsFacet, eventsFacet}

	for i, facet := range facets {
		if facet.IsLoaded() {
			t.Errorf("Facet %d should not be loaded initially", i)
		}

		if !facet.NeedsUpdate() {
			t.Errorf("Facet %d should need update initially", i)
		}

		if facet.Count() != 0 {
			t.Errorf("Facet %d: expected count 0, got %d", i, facet.Count())
		}
	}

	t.Log("Abis facet test completed successfully")
}

// ADD_ROUTE
