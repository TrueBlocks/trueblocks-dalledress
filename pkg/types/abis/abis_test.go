// ABIS_ROUTE
package abis

import (
	"testing"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/facets"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
)

func TestEnhancedAbisCollection(t *testing.T) {
	msgs.SetTestMode(true)
	defer msgs.SetTestMode(false)

	collection := NewAbisCollection()

	if state := collection.downloadedFacet.GetState(); state != facets.StateStale {
		t.Errorf("Expected initial state to be Stale, got %v", state)
	}

	if collection.downloadedFacet.GetState() != facets.StateStale {
		t.Error("downloadedFacet should be in StateStale")
	}
	if collection.knownFacet.GetState() != facets.StateStale {
		t.Error("knownFacet should be in StateStale")
	}
	if collection.functionsFacet.GetState() != facets.StateStale {
		t.Error("functionsFacet should be in StateStale")
	}
	if collection.eventsFacet.GetState() != facets.StateStale {
		t.Error("eventsFacet should be in StateStale")
	}

	if !collection.NeedsUpdate(AbisDownloaded) {
		t.Error("downloadedFacet should need update")
	}
	if !collection.NeedsUpdate(AbisKnown) {
		t.Error("knownFacet should need update")
	}
	if !collection.NeedsUpdate(AbisFunctions) {
		t.Error("functionsFacet should need update")
	}
	if !collection.NeedsUpdate(AbisEvents) {
		t.Error("eventsFacet should need update")
	}

	// Test Reset function
	collection.Reset(AbisDownloaded)
	if !collection.NeedsUpdate(AbisDownloaded) {
		t.Error("downloadedFacet should need update after reset")
	}
}

// func TestNewAbisCollection(t *testing.T) {
// 	ac := NewAbisCollection()

// 	// Check that all facets are properly initialized and not loaded
// 	if ac.downloadedFacet == nil {
// 		t.Error("NewAbisCollection should initialize downloadedFacet")
// 	}
// 	if ac.downloadedFacet.IsLoaded() {
// 		t.Error("NewAbisCollection downloadedFacet should not be loaded initially")
// 	}
// 	if ac.downloadedFacet.IsFetching() {
// 		t.Error("NewAbisCollection downloadedFacet should not be fetching initially")
// 	}

// 	if ac.knownFacet == nil {
// 		t.Error("NewAbisCollection should initialize knownFacet")
// 	}
// 	if ac.knownFacet.IsLoaded() {
// 		t.Error("NewAbisCollection knownFacet should not be loaded initially")
// 	}
// 	if ac.knownFacet.IsFetching() {
// 		t.Error("NewAbisCollection knownFacet should not be fetching initially")
// 	}

// 	if ac.functionsFacet == nil {
// 		t.Error("NewAbisCollection should initialize functionsFacet")
// 	}
// 	if ac.functionsFacet.IsLoaded() {
// 		t.Error("NewAbisCollection functionsFacet should not be loaded initially")
// 	}
// 	if ac.functionsFacet.IsFetching() {
// 		t.Error("NewAbisCollection functionsFacet should not be fetching initially")
// 	}

// 	if ac.eventsFacet == nil {
// 		t.Error("NewAbisCollection should initialize eventsFacet")
// 	}
// 	if ac.eventsFacet.IsLoaded() {
// 		t.Error("NewAbisCollection eventsFacet should not be loaded initially")
// 	}
// 	if ac.eventsFacet.IsFetching() {
// 		t.Error("NewAbisCollection eventsFacet should not be fetching initially")
// 	}
// }

// func TestAbisCollection_ClearCache(t *testing.T) {
// 	ac := NewAbisCollection()

// 	// For all data types now using Facet pattern
// 	// We can't easily simulate initial data without loading, so we test the clear operation
// 	ac.Reset(AbisDownloaded)
// 	if ac.downloadedFacet.IsLoaded() {
// 		t.Error("ClearCache(AbisDownloaded) should clear the facet loaded state")
// 	}

// 	ac.Reset(AbisKnown)
// 	if ac.knownFacet.IsLoaded() {
// 		t.Error("ClearCache(AbisKnown) should clear the facet loaded state")
// 	}

// 	ac.Reset(AbisFunctions)
// 	if ac.functionsFacet.IsLoaded() {
// 		t.Error("ClearCache(AbisFunctions) should clear the facet loaded state")
// 	}

// 	ac.Reset(AbisEvents)
// 	if ac.eventsFacet.IsLoaded() {
// 		t.Error("ClearCache(AbisEvents) should clear the facet loaded state")
// 	}
// }

// func TestAbisCollection_NeedsUpdate(t *testing.T) {
// 	ac := NewAbisCollection()

// 	// Initially all should need update (facets not loaded)
// 	if !ac.NeedsUpdate(AbisDownloaded) {
// 		t.Error("NeedsUpdate(AbisDownloaded) should return true when facet is not loaded")
// 	}
// 	if !ac.NeedsUpdate(AbisKnown) {
// 		t.Error("NeedsUpdate(AbisKnown) should return true when facet is not loaded")
// 	}
// 	if !ac.NeedsUpdate(AbisFunctions) {
// 		t.Error("NeedsUpdate(AbisFunctions) should return true when facet is not loaded")
// 	}
// 	if !ac.NeedsUpdate(AbisEvents) {
// 		t.Error("NeedsUpdate(AbisEvents) should return true when facet is not loaded")
// 	}

// 	// Note: We can't easily simulate loaded facets in tests without actual data loading
// 	// so this test primarily verifies the initial state behavior
// }

// func TestAbisCollection_NeedsUpdate_UnknownListKind(t *testing.T) {
// 	ac := NewAbisCollection()

// 	// Test with unknown ListKind - should return true
// 	unknownKind := types.ListKind("Unknown")
// 	if !ac.NeedsUpdate(unknownKind) {
// 		t.Error("NeedsUpdate with unknown ListKind should return true")
// 	}
// }

// func TestAbisCollection_ThreadSafety(t *testing.T) {
// 	ac := NewAbisCollection()

// 	// Test concurrent access to NeedsUpdate and ClearCache across all slice types
// 	done := make(chan bool, 4)

// 	// Goroutine 1: repeatedly check NeedsUpdate for downloaded and known
// 	go func() {
// 		for i := 0; i < 100; i++ {
// 			ac.NeedsUpdate(AbisDownloaded)
// 			ac.NeedsUpdate(AbisKnown)
// 		}
// 		done <- true
// 	}()

// 	// Goroutine 2: repeatedly clear cache for downloaded and known
// 	go func() {
// 		for i := 0; i < 100; i++ {
// 			ac.Reset(AbisDownloaded)
// 			ac.Reset(AbisKnown)
// 		}
// 		done <- true
// 	}()

// 	// Goroutine 3: repeatedly check NeedsUpdate for functions and events
// 	go func() {
// 		for i := 0; i < 100; i++ {
// 			ac.NeedsUpdate(AbisFunctions)
// 			ac.NeedsUpdate(AbisEvents)
// 		}
// 		done <- true
// 	}()

// 	// Goroutine 4: repeatedly clear cache for functions and events
// 	go func() {
// 		for i := 0; i < 100; i++ {
// 			ac.Reset(AbisFunctions)
// 			ac.Reset(AbisEvents)
// 		}
// 		done <- true
// 	}()

// 	// Wait for all goroutines to complete
// 	<-done
// 	<-done
// 	<-done
// 	<-done

// 	// If we get here without panicking, the per-slice mutexes are working
// 	t.Log("Thread safety test completed successfully")
// }

// func TestAbisFacet(t *testing.T) {
// 	sharedAbisListSource := GetSharedAbisListSource()
// 	sharedAbisDetailsSource := GetSharedAbisDetailsSource()

// 	downloadedFacet := facets.NewBaseFacet(
// 		AbisDownloaded,
// 		func(abi *coreTypes.Abi) bool { return !abi.IsKnown },
// 		nil,
// 		sharedAbisListSource,
// 	)

// 	knownFacet := facets.NewBaseFacet(
// 		AbisKnown,
// 		func(abi *coreTypes.Abi) bool { return abi.IsKnown },
// 		nil,
// 		sharedAbisListSource,
// 	)

// 	functionsFacet := facets.NewBaseFacet(
// 		AbisFunctions,
// 		func(fn *coreTypes.Function) bool { return fn.FunctionType != "event" },
// 		IsDupFuncByEncoding(),
// 		sharedAbisDetailsSource,
// 	)

// 	eventsFacet := facets.NewBaseFacet(
// 		AbisEvents,
// 		func(fn *coreTypes.Function) bool { return fn.FunctionType == "event" },
// 		IsDupFuncByEncoding(),
// 		sharedAbisDetailsSource,
// 	)

// 	if downloadedFacet == nil {
// 		t.Error("Downloaded facet should not be nil")
// 	}

// 	if knownFacet == nil {
// 		t.Error("Known facet should not be nil")
// 	}

// 	if functionsFacet == nil {
// 		t.Error("Functions facet should not be nil")
// 	}

// 	if eventsFacet == nil {
// 		t.Error("Events facet should not be nil")
// 	}

// 	// Test initial states - all facets should implement the Facet interface
// 	facets := []interface {
// 		IsLoaded() bool
// 		NeedsUpdate() bool
// 		Count() int
// 	}{downloadedFacet, knownFacet, functionsFacet, eventsFacet}

// 	for i, facet := range facets {
// 		if facet.IsLoaded() {
// 			t.Errorf("Facet %d should not be loaded initially", i)
// 		}

// 		if !facet.NeedsUpdate() {
// 			t.Errorf("Facet %d should need update initially", i)
// 		}

// 		if facet.Count() != 0 {
// 			t.Errorf("Facet %d: expected count 0, got %d", i, facet.Count())
// 		}
// 	}

// 	if sharedAbisListSource == nil {
// 		t.Error("Shared ABI list source should not be nil")
// 	}

// 	if sharedAbisDetailsSource == nil {
// 		t.Error("Shared ABI details source should not be nil")
// 	}

// 	if GetSharedAbisListSource() != sharedAbisListSource {
// 		t.Error("GetSharedAbisListSource should return the same instance (singleton pattern)")
// 	}

// 	if GetSharedAbisDetailsSource() != sharedAbisDetailsSource {
// 		t.Error("GetSharedAbisDetailsSource should return the same instance (singleton pattern)")
// 	}

// 	t.Log("Abis facet test completed successfully - all facets created with shared sources")
// }

// func TestSourceBasedFacets(t *testing.T) {
// 	// Test creating shared sources
// 	abisListSource := GetSharedAbisListSource()
// 	if abisListSource == nil {
// 		t.Fatal("GetSharedAbisListSource should not return nil")
// 	}

// 	if abisListSource.GetSourceType() != "sdk" {
// 		t.Errorf("Expected source type 'sdk', got '%s'", abisListSource.GetSourceType())
// 	}

// 	abisDetailsSource := GetSharedAbisDetailsSource()
// 	if abisDetailsSource == nil {
// 		t.Fatal("GetSharedAbisDetailsSource should not return nil")
// 	}

// 	if abisDetailsSource.GetSourceType() != "sdk" {
// 		t.Errorf("Expected source type 'sdk', got '%s'", abisDetailsSource.GetSourceType())
// 	}
// }

// func TestAbisCollectionWithSources(t *testing.T) {
// 	// Test creating collection with sources
// 	ac := NewAbisCollection()

// 	// Verify all facets are properly initialized
// 	if ac.downloadedFacet == nil {
// 		t.Error("downloadedFacet should not be nil")
// 	}
// 	if ac.knownFacet == nil {
// 		t.Error("knownFacet should not be nil")
// 	}
// 	if ac.functionsFacet == nil {
// 		t.Error("functionsFacet should not be nil")
// 	}
// 	if ac.eventsFacet == nil {
// 		t.Error("eventsFacet should not be nil")
// 	}

// 	// Verify initial state
// 	if ac.downloadedFacet.IsLoaded() {
// 		t.Error("downloadedFacet should not be loaded initially")
// 	}
// 	if ac.downloadedFacet.IsFetching() {
// 		t.Error("downloadedFacet should not be fetching initially")
// 	}
// }

// func TestSourceBasedFacetIntegration(t *testing.T) {
// 	// Create a source-based facet for downloaded ABIs using the shared source
// 	abisListSource := GetSharedAbisListSource()
// 	downloadedFacet := facets.NewBaseFacet(
// 		AbisDownloaded,
// 		func(abi *coreTypes.Abi) bool { return !abi.IsKnown },
// 		nil, // No deduplication for this test
// 		abisListSource,
// 	)

// 	// Verify the facet implements the Facet interface
// 	if downloadedFacet == nil {
// 		t.Fatal("BaseFacet should not be nil")
// 	}

// 	// Test basic facet operations
// 	if downloadedFacet.IsLoaded() {
// 		t.Error("New facet should not be loaded")
// 	}
// 	if downloadedFacet.IsFetching() {
// 		t.Error("New facet should not be fetching")
// 	}
// 	if downloadedFacet.Count() != 0 {
// 		t.Error("New facet should have count of 0")
// 	}
// 	if !downloadedFacet.NeedsUpdate() {
// 		t.Error("New facet should need update")
// 	}

// 	// Test the Load method exists and can be called
// 	// Note: We can't easily test actual loading without mocking the SDK
// 	// but we can verify the method signature is correct
// 	result, err := downloadedFacet.Load()
// 	if err != facets.ErrorAlreadyLoading {
// 		// The method should be callable, though it may fail due to no actual SDK setup
// 		// We're primarily testing that the interface is correctly implemented
// 		t.Logf("Load called successfully (result: %v, err: %v)", result, err)
// 	}
// }

// func TestLoadDataFromSourceMethod(t *testing.T) {
// 	ac := NewAbisCollection()

// 	// Test that LoadData method exists and can be called
// 	// This tests the demonstration method we added
// 	ac.LoadData(AbisDownloaded)
// 	ac.LoadData(AbisKnown)
// 	ac.LoadData(AbisFunctions)
// 	ac.LoadData(AbisEvents)

// 	// Give any goroutines a moment to start
// 	time.Sleep(10 * time.Millisecond)

// 	// Test with invalid list kind
// 	ac.LoadData("InvalidKind")

// 	t.Log("LoadData method tested successfully")
// }

// func TestSourceInterfaceCompatibility(t *testing.T) {
// 	// Test that shared sources implement the expected interface
// 	sources := []interface{}{
// 		GetSharedAbisListSource(),
// 		GetSharedAbisDetailsSource(),
// 	}

// 	for i, source := range sources {
// 		// Check that each source has the expected methods
// 		// This is a compile-time check that verifies interface compatibility
// 		switch s := source.(type) {
// 		case interface{ GetSourceType() string }:
// 			sourceType := s.GetSourceType()
// 			if sourceType != "sdk" {
// 				t.Errorf("Source %d: expected type 'sdk', got '%s'", i, sourceType)
// 			}
// 		default:
// 			t.Errorf("Source %d does not implement GetSourceType method", i)
// 		}
// 	}
// }

// // TestSourceSharingDemo demonstrates the key benefit of the Source → Facet pattern:
// // Multiple facets sharing the same data sources to eliminate redundant SDK queries.
// func TestSourceSharingDemo(t *testing.T) {
// 	// Create a collection using the new source-based pattern
// 	sourceBasedCollection := NewAbisCollection()

// 	// Verify all facets are properly initialized
// 	if sourceBasedCollection.downloadedFacet == nil {
// 		t.Error("Source-based collection should have downloadedFacet")
// 	}
// 	if sourceBasedCollection.knownFacet == nil {
// 		t.Error("Source-based collection should have knownFacet")
// 	}
// 	if sourceBasedCollection.functionsFacet == nil {
// 		t.Error("Source-based collection should have functionsFacet")
// 	}
// 	if sourceBasedCollection.eventsFacet == nil {
// 		t.Error("Source-based collection should have eventsFacet")
// 	}

// 	// Test that the sources are properly shared
// 	// Note: This is a conceptual test since the sources are internal to the facets
// 	// In a real implementation, we might expose source IDs for verification

// 	// Test LoadData capability
// 	sourceBasedCollection.LoadData(AbisDownloaded)
// 	sourceBasedCollection.LoadData(AbisKnown)
// 	sourceBasedCollection.LoadData(AbisFunctions)
// 	sourceBasedCollection.LoadData(AbisEvents)

// 	// Give the goroutines a moment to start
// 	time.Sleep(10 * time.Millisecond)

// 	// For source-based facets, we expect them to either be loading or loaded
// 	// The key insight is that when Downloaded facet loads, Known facet should
// 	// benefit from the same data source (and vice versa)

// 	t.Log("Source sharing demo completed successfully")
// 	t.Log("Key benefit: Downloaded + Known facets share ONE AbisList source")
// 	t.Log("Key benefit: Functions + Events facets share ONE AbisDetails source")
// 	t.Log("Result: 4 facets powered by only 2 SDK queries instead of 4!")
// }

// // TestSharedSourcesReturnSameInstance verifies that our shared source functions
// // return the same instances (singleton pattern)
// func TestSharedSourcesReturnSameInstance(t *testing.T) {
// 	// Get shared sources multiple times
// 	source1 := GetSharedAbisListSource()
// 	source2 := GetSharedAbisListSource()
// 	source3 := GetSharedAbisDetailsSource()
// 	source4 := GetSharedAbisDetailsSource()

// 	// Verify they return the same instances (pointer equality)
// 	if source1 != source2 {
// 		t.Error("GetSharedAbisListSource should return the same instance")
// 	}
// 	if source3 != source4 {
// 		t.Error("GetSharedAbisDetailsSource should return the same instance")
// 	}

// 	t.Log("Shared sources properly implement singleton pattern")
// }

// // TestSourceSharingComparison compares the old vs new patterns
// func TestSourceSharingComparison(t *testing.T) {
// 	// Create collections using both patterns
// 	originalCollection := NewAbisCollection()    // Old pattern: 4 queries
// 	sourceBasedCollection := NewAbisCollection() // New pattern: 2 queries

// 	// Both should have the same structure
// 	if originalCollection.downloadedFacet == nil || sourceBasedCollection.downloadedFacet == nil {
// 		t.Error("Both collections should have downloadedFacet")
// 	}
// 	if originalCollection.knownFacet == nil || sourceBasedCollection.knownFacet == nil {
// 		t.Error("Both collections should have knownFacet")
// 	}
// 	if originalCollection.functionsFacet == nil || sourceBasedCollection.functionsFacet == nil {
// 		t.Error("Both collections should have functionsFacet")
// 	}
// 	if originalCollection.eventsFacet == nil || sourceBasedCollection.eventsFacet == nil {
// 		t.Error("Both collections should have eventsFacet")
// 	}

// 	// The key difference is in the implementation:
// 	// - Original: Each facet makes its own SDK query (4 total)
// 	// - Source-based: Facets share sources (2 total)

// 	t.Log("Comparison completed:")
// 	t.Log("  Original pattern: 4 facets → 4 SDK queries")
// 	t.Log("  Source-based pattern: 4 facets → 2 shared sources → 2 SDK queries")
// 	t.Log("  Efficiency improvement: 50% reduction in SDK calls!")
// }

// // MockSDKCallCounter tracks how many times SDK functions are called
// type MockSDKCallCounter struct {
// 	mu            sync.Mutex
// 	abisListCalls int
// 	abisDetsCalls int
// }

// func (m *MockSDKCallCounter) IncrementAbisList() {
// 	m.mu.Lock()
// 	defer m.mu.Unlock()
// 	m.abisListCalls++
// }

// func (m *MockSDKCallCounter) IncrementAbisDetails() {
// 	m.mu.Lock()
// 	defer m.mu.Unlock()
// 	m.abisDetsCalls++
// }

// func (m *MockSDKCallCounter) GetCounts() (int, int) {
// 	m.mu.Lock()
// 	defer m.mu.Unlock()
// 	return m.abisListCalls, m.abisDetsCalls
// }

// // var globalMockCounter = &MockSDKCallCounter{}

// // TestSourceSharingArchitecture verifies that the source-sharing architecture is set up correctly
// // without triggering the race conditions in the core library
// func TestSourceSharingArchitecture(t *testing.T) {
// 	t.Run("Old Architecture Uses Separate LoadData Methods", func(t *testing.T) {
// 		// Create old-style collection
// 		oldCollection := NewAbisCollection()

// 		// Verify sequential loading works (avoiding concurrent calls that trigger core race conditions)
// 		t.Log("Testing old architecture sequentially to avoid core library race conditions")

// 		// Test that the old collection supports the LoadData interface
// 		if oldCollection.downloadedFacet == nil {
// 			t.Error("downloadedFacet should be initialized")
// 		}
// 		if oldCollection.knownFacet == nil {
// 			t.Error("knownFacet should be initialized")
// 		}
// 		if oldCollection.functionsFacet == nil {
// 			t.Error("functionsFacet should be initialized")
// 		}
// 		if oldCollection.eventsFacet == nil {
// 			t.Error("eventsFacet should be initialized")
// 		}

// 		t.Log("Old architecture: Each facet makes separate SDK calls")
// 	})

// 	t.Run("New Architecture Uses Shared Sources", func(t *testing.T) {
// 		// Create new source-based collection
// 		newCollection := NewAbisCollection()

// 		// Verify that shared sources exist
// 		abisListSource := GetSharedAbisListSource()
// 		abisDetailsSource := GetSharedAbisDetailsSource()

// 		if abisListSource == nil {
// 			t.Error("Shared AbisList source should be available")
// 		}
// 		if abisDetailsSource == nil {
// 			t.Error("Shared AbisDetails source should be available")
// 		}

// 		// Verify that the LoadData method exists and gracefully handles both architectures
// 		if newCollection.downloadedFacet == nil {
// 			t.Error("downloadedFacet should be initialized")
// 		}

// 		t.Log("New architecture: Downloaded+Known share one source, Functions+Events share another")
// 	})
// }

// // TestSourceBasedCollectionArchitecture verifies the source-based collection structure
// func TestSourceBasedCollectionArchitecture(t *testing.T) {
// 	// Create source-based collection
// 	ac := NewAbisCollection()

// 	// Verify that facets are the right type for source-based operations
// 	t.Run("Downloaded Facet Type", func(t *testing.T) {
// 		if _, ok := ac.downloadedFacet.(*facets.BaseFacet[coreTypes.Abi]); !ok {
// 			t.Error("downloadedFacet should be a BaseFacet for source sharing")
// 		}
// 	})

// 	t.Run("Known Facet Type", func(t *testing.T) {
// 		if _, ok := ac.knownFacet.(*facets.BaseFacet[coreTypes.Abi]); !ok {
// 			t.Error("knownFacet should be a BaseFacet for source sharing")
// 		}
// 	})

// 	t.Run("Functions Facet Type", func(t *testing.T) {
// 		if _, ok := ac.functionsFacet.(*facets.BaseFacet[coreTypes.Function]); !ok {
// 			t.Error("functionsFacet should be a BaseFacet for source sharing")
// 		}
// 	})

// 	t.Run("Events Facet Type", func(t *testing.T) {
// 		if _, ok := ac.eventsFacet.(*facets.BaseFacet[coreTypes.Function]); !ok {
// 			t.Error("eventsFacet should be a BaseFacet for source sharing")
// 		}
// 	})
// }

// // TestLoadDataFromSourceFallback verifies fallback behavior
// func TestLoadDataFromSourceFallback(t *testing.T) {
// 	t.Run("Source-Based Collection", func(t *testing.T) {
// 		ac := NewAbisCollection()

// 		// This should use the source-based loading
// 		ac.LoadData(AbisDownloaded)

// 		// Give it a moment to start
// 		time.Sleep(10 * time.Millisecond)

// 		// Should be using source-based facets
// 		downloadedSourceFacet, ok := ac.downloadedFacet.(*facets.BaseFacet[coreTypes.Abi])
// 		if !ok {
// 			t.Error("Should be using BaseFacet for source-based collection")
// 		} else {
// 			t.Logf("✓ Using source-based facet: %T", downloadedSourceFacet)
// 		}
// 	})
// }

// // TestSharedSourceSingleton verifies that sources are properly shared
// func TestSharedSourceSingleton(t *testing.T) {
// 	t.Run("AbisList Source Sharing", func(t *testing.T) {
// 		// Get the shared source multiple times
// 		source1 := GetSharedAbisListSource()
// 		source2 := GetSharedAbisListSource()

// 		// Should be the same instance (singleton)
// 		if source1 != source2 {
// 			t.Error("GetSharedAbisListSource should return the same instance (singleton)")
// 		}

// 		// Verify source type
// 		if source1.GetSourceType() != "sdk" {
// 			t.Errorf("Expected source type 'sdk', got '%s'", source1.GetSourceType())
// 		}
// 	})

// 	t.Run("AbisDetails Source Sharing", func(t *testing.T) {
// 		// Get the shared source multiple times
// 		source1 := GetSharedAbisDetailsSource()
// 		source2 := GetSharedAbisDetailsSource()

// 		// Should be the same instance (singleton)
// 		if source1 != source2 {
// 			t.Error("GetSharedAbisDetailsSource should return the same instance (singleton)")
// 		}

// 		// Verify source type
// 		if source1.GetSourceType() != "sdk" {
// 			t.Errorf("Expected source type 'sdk', got '%s'", source1.GetSourceType())
// 		}
// 	})
// }

// // TestSourceBasedCollectionCompatibility verifies that source-based collections
// // maintain the same interface as traditional collections
// func TestSourceBasedCollectionCompatibility(t *testing.T) {
// 	oldCollection := NewAbisCollection()
// 	newCollection := NewAbisCollection()

// 	// Both should implement the same basic operations
// 	listKinds := []types.ListKind{AbisDownloaded, AbisKnown, AbisFunctions, AbisEvents}

// 	for _, kind := range listKinds {
// 		t.Run("NeedsUpdate-"+string(kind), func(t *testing.T) {
// 			// Both should report needing updates initially
// 			if !oldCollection.NeedsUpdate(kind) {
// 				t.Errorf("Old collection should need update for %s", kind)
// 			}
// 			if !newCollection.NeedsUpdate(kind) {
// 				t.Errorf("New collection should need update for %s", kind)
// 			}
// 		})

// 		t.Run("ClearCache-"+string(kind), func(t *testing.T) {
// 			// Both should support cache clearing
// 			oldCollection.Reset(kind)
// 			newCollection.Reset(kind)

// 			// Should still need updates after clearing
// 			if !oldCollection.NeedsUpdate(kind) {
// 				t.Errorf("Old collection should need update after clearing %s", kind)
// 			}
// 			if !newCollection.NeedsUpdate(kind) {
// 				t.Errorf("New collection should need update after clearing %s", kind)
// 			}
// 		})
// 	}
// }

// // TestConcurrentSourceAccess verifies thread safety of shared sources without triggering SDK calls
// func TestConcurrentSourceAccess(t *testing.T) {
// 	const numGoroutines = 10

// 	t.Run("Concurrent AbisList Source Access", func(t *testing.T) {
// 		var wg sync.WaitGroup
// 		sources := make([]interface{}, numGoroutines)

// 		for i := 0; i < numGoroutines; i++ {
// 			wg.Add(1)
// 			go func(index int) {
// 				defer wg.Done()
// 				sources[index] = GetSharedAbisListSource()
// 			}(i)
// 		}

// 		wg.Wait()

// 		// All should be the same instance (singleton pattern)
// 		firstSource := sources[0]
// 		for i, source := range sources {
// 			if source != firstSource {
// 				t.Errorf("Source %d differs from first source - singleton not working", i)
// 			}
// 		}
// 		t.Log("All goroutines got the same shared AbisList source instance")
// 	})

// 	t.Run("Concurrent AbisDetails Source Access", func(t *testing.T) {
// 		var wg sync.WaitGroup
// 		sources := make([]interface{}, numGoroutines)

// 		for i := 0; i < numGoroutines; i++ {
// 			wg.Add(1)
// 			go func(index int) {
// 				defer wg.Done()
// 				sources[index] = GetSharedAbisDetailsSource()
// 			}(i)
// 		}

// 		wg.Wait()

// 		// All should be the same instance (singleton pattern)
// 		firstSource := sources[0]
// 		for i, source := range sources {
// 			if source != firstSource {
// 				t.Errorf("Source %d differs from first source - singleton not working", i)
// 			}
// 		}
// 		t.Log("All goroutines got the same shared AbisDetails source instance")
// 	})
// }

// // TestLoadData tests the LoadData function from abis_load.go
// func TestLoadData(t *testing.T) {
// 	// TODO: Turn this back on: Removed t.Parallel() to prevent concurrent SDK access causing race conditions
// 	ac := NewAbisCollection()

// 	// Verify initial state - all facets should not be loaded
// 	if ac.downloadedFacet.IsLoaded() || ac.knownFacet.IsLoaded() || ac.functionsFacet.IsLoaded() || ac.eventsFacet.IsLoaded() {
// 		t.Error("NewAbisCollection should not have any loaded facets")
// 	}
// 	if ac.downloadedFacet.IsFetching() || ac.knownFacet.IsFetching() || ac.functionsFacet.IsFetching() || ac.eventsFacet.IsFetching() {
// 		t.Error("NewAbisCollection should not have any fetching facets")
// 	}

// 	// Call LoadData for known ABIs
// 	ac.LoadData(AbisKnown)

// 	// Give the goroutine a moment to start
// 	time.Sleep(10 * time.Millisecond)

// 	// After calling LoadData, the known facet should either be fetching or loaded
// 	// Note: Due to the asynchronous nature, we can't easily test completion without
// 	// more complex synchronization mechanisms
// 	isFetchingOrLoaded := ac.knownFacet.IsFetching() || ac.knownFacet.IsLoaded()

// 	if !isFetchingOrLoaded {
// 		t.Error("After LoadData, known facet should be fetching or loaded")
// 	}

// 	// Test that calling LoadData again doesn't cause issues when already fetching/loaded
// 	prevState := ac.knownFacet.IsFetching() || ac.knownFacet.IsLoaded()
// 	ac.LoadData(AbisKnown)

// 	newState := ac.knownFacet.IsFetching() || ac.knownFacet.IsLoaded()

// 	if !prevState || !newState {
// 		t.Error("LoadData should not change state when already fetching/loaded")
// 	}
// }

// func TestLoadAbis(t *testing.T) {
// 	ac := NewAbisCollection()

// 	// Verify initial state - all facets should not be loaded
// 	if ac.downloadedFacet.IsLoaded() || ac.knownFacet.IsLoaded() || ac.functionsFacet.IsLoaded() || ac.eventsFacet.IsLoaded() {
// 		t.Fatalf("NewAbisCollection should not have any loaded facets initially")
// 	}

// 	ac.LoadData(AbisKnown)
// 	// Note: LoadData is async. Testing its completion requires a different approach.
// 	// For now, we'll assume it kicks off the process.
// 	// We can check isFetching state.

// 	// The rest of this test needs significant rework due to the asynchronous nature
// 	// of loadInternal and the facet pattern replacing direct slice management.
// 	// We cannot directly check ac.loaded or counts immediately after LoadData.
// 	// We would need to:
// 	// 1. Wait for isLoaded to become true (with a timeout).
// 	// 2. Or, inspect events emitted.
// 	// For now, this part of the test is effectively disabled or needs to be redesigned.
// 	t.Skip("TestLoadAbis needs rework for asynchronous loading via LoadData and facet pattern.")
// }

// // TestReloadCancellation tests that Reload properly cancels ongoing operations
// func TestReloadCancellation(t *testing.T) {
// 	abisAddr := base.ZeroAddr.Hex()
// 	renderCtx := sources.RegisterContext(abisAddr)

// 	// Verify the context was registered
// 	if sources.CtxCount(abisAddr) != 1 {
// 		t.Errorf("Expected 1 registered context, got %d", sources.CtxCount(abisAddr))
// 	}

// 	// Verify the context is not nil
// 	if renderCtx == nil {
// 		t.Error("RegisterContext should return non-nil context")
// 	}

// 	// Simulate a reload operation by cancelling the context
// 	cancelled, found := sources.UnregisterContext(abisAddr)
// 	if !found {
// 		t.Error("Cancel should find the registered context")
// 	}
// 	if cancelled != 1 {
// 		t.Errorf("Expected 1 cancelled context, got %d", cancelled)
// 	}

// 	// Verify the context was cancelled and removed
// 	if sources.CtxCount(abisAddr) != 0 {
// 		t.Errorf("Expected 0 registered contexts after reload, got %d", sources.CtxCount(abisAddr))
// 	}

// 	// Note: We can't easily test if the context was actually cancelled since
// 	// the Cancel method removes it from the map, but the fact that it was
// 	// removed indicates it was processed correctly
// }

// // TestContextRegistration tests that contexts are properly registered and cleaned up
// func TestContextRegistration(t *testing.T) {
// 	addr1 := "0x1234567890123456789012345678901234567890"
// 	addr2 := "0x2234567890123456789012345678901234567890"

// 	ctx1 := sources.RegisterContext(addr1)
// 	ctx2 := sources.RegisterContext(addr2)

// 	cnt := sources.CtxCount(addr1) + sources.CtxCount(addr1)
// 	if cnt != 2 {
// 		t.Errorf("Expected 2 registered contexts, got %d", cnt)
// 	}

// 	if ctx1 == nil || ctx2 == nil {
// 		t.Error("RegisterContext should return non-nil contexts")
// 	}

// 	// Test Cancel for specific address
// 	cancelled, found := sources.UnregisterContext(addr1)
// 	if !found {
// 		t.Error("Cancel should find the registered context")
// 	}
// 	if cancelled != 1 {
// 		t.Errorf("Expected 1 cancelled context, got %d", cancelled)
// 	}
// 	cnt = sources.CtxCount(addr1) + sources.CtxCount(addr2)
// 	if cnt != 1 {
// 		t.Errorf("Expected 1 remaining context after cancel, got %d", cnt)
// 	}

// 	// Test Cancel for non-existent address
// 	nonExistentAddr := "0x9999999999999999999999999999999999999999"
// 	cancelled, found = sources.UnregisterContext(nonExistentAddr)
// 	if found {
// 		t.Error("Cancel should not find non-existent context")
// 	}
// 	if cancelled != 0 {
// 		t.Errorf("Expected 0 cancelled contexts for non-existent address, got %d", cancelled)
// 	}
// }

// func TestGetPage_BasicFunctionality(t *testing.T) {
// 	ac := NewAbisCollection()

// 	tests := []struct {
// 		name     string
// 		listKind types.ListKind
// 	}{
// 		{"Downloaded", AbisDownloaded},
// 		{"Known", AbisKnown},
// 		{"Functions", AbisFunctions},
// 		{"Events", AbisEvents},
// 	}

// 	for _, tc := range tests {
// 		t.Run(tc.name, func(t *testing.T) {
// 			page, err := ac.GetPage(tc.listKind, 0, 10, sorting.EmptySortSpec(), "")
// 			if err != nil {
// 				t.Fatalf("GetPage returned error: %v", err)
// 			}

// 			if page.Kind != tc.listKind {
// 				t.Errorf("Expected page kind %s, got %s", tc.listKind, page.Kind)
// 			}

// 			if page.TotalItems != 0 {
// 				t.Errorf("Expected 0 total items for empty facet, got %d", page.TotalItems)
// 			}

// 			if len(page.Abis) != 0 && len(page.Functions) != 0 {
// 				t.Errorf("Expected empty results for empty facet")
// 			}
// 		})
// 	}
// }

// func TestGetPage_InvalidListKind(t *testing.T) {
// 	ac := NewAbisCollection()

// 	_, err := ac.GetPage(types.ListKind("InvalidKind"), 0, 10, sorting.EmptySortSpec(), "")
// 	if err == nil {
// 		t.Error("Expected error for invalid list kind, got nil")
// 	}
// }
// ABIS_ROUTE
