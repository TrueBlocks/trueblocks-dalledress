// ADD_ROUTE
package abis

import (
	"testing"

	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/mocks"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

func TestNewAbisCollection(t *testing.T) {
	mockApp := mocks.NewMockApp()
	ac := NewAbisCollection(mockApp)

	// Check that all repositories are properly initialized and not loaded
	if ac.downloadedRepo == nil {
		t.Error("NewAbisCollection should initialize downloadedRepo")
	}
	if ac.downloadedRepo.IsLoaded() {
		t.Error("NewAbisCollection downloadedRepo should not be loaded initially")
	}
	if ac.downloadedRepo.IsLoading() {
		t.Error("NewAbisCollection downloadedRepo should not be loading initially")
	}

	if ac.knownRepo == nil {
		t.Error("NewAbisCollection should initialize knownRepo")
	}
	if ac.knownRepo.IsLoaded() {
		t.Error("NewAbisCollection knownRepo should not be loaded initially")
	}
	if ac.knownRepo.IsLoading() {
		t.Error("NewAbisCollection knownRepo should not be loading initially")
	}

	if ac.functionsRepo == nil {
		t.Error("NewAbisCollection should initialize functionsRepo")
	}
	if ac.functionsRepo.IsLoaded() {
		t.Error("NewAbisCollection functionsRepo should not be loaded initially")
	}
	if ac.functionsRepo.IsLoading() {
		t.Error("NewAbisCollection functionsRepo should not be loading initially")
	}

	if ac.eventsRepo == nil {
		t.Error("NewAbisCollection should initialize eventsRepo")
	}
	if ac.eventsRepo.IsLoaded() {
		t.Error("NewAbisCollection eventsRepo should not be loaded initially")
	}
	if ac.eventsRepo.IsLoading() {
		t.Error("NewAbisCollection eventsRepo should not be loading initially")
	}
}

func TestAbisCollection_ClearCache(t *testing.T) {
	mockApp := mocks.NewMockApp()
	ac := NewAbisCollection(mockApp)

	// For all data types now using Repository pattern
	// We can't easily simulate initial data without loading, so we test the clear operation
	ac.ClearCache(AbisDownloaded)
	if ac.downloadedRepo.IsLoaded() {
		t.Error("ClearCache(AbisDownloaded) should clear the repository loaded state")
	}

	ac.ClearCache(AbisKnown)
	if ac.knownRepo.IsLoaded() {
		t.Error("ClearCache(AbisKnown) should clear the repository loaded state")
	}

	ac.ClearCache(AbisFunctions)
	if ac.functionsRepo.IsLoaded() {
		t.Error("ClearCache(AbisFunctions) should clear the repository loaded state")
	}

	ac.ClearCache(AbisEvents)
	if ac.eventsRepo.IsLoaded() {
		t.Error("ClearCache(AbisEvents) should clear the repository loaded state")
	}
}

func TestAbisCollection_ClearCache_UnknownListKind(t *testing.T) {
	mockApp := mocks.NewMockApp()
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
	mockApp := mocks.NewMockApp()
	ac := NewAbisCollection(mockApp)

	// Initially all should need update (repositories not loaded)
	if !ac.NeedsUpdate(AbisDownloaded) {
		t.Error("NeedsUpdate(AbisDownloaded) should return true when repository is not loaded")
	}
	if !ac.NeedsUpdate(AbisKnown) {
		t.Error("NeedsUpdate(AbisKnown) should return true when repository is not loaded")
	}
	if !ac.NeedsUpdate(AbisFunctions) {
		t.Error("NeedsUpdate(AbisFunctions) should return true when repository is not loaded")
	}
	if !ac.NeedsUpdate(AbisEvents) {
		t.Error("NeedsUpdate(AbisEvents) should return true when repository is not loaded")
	}

	// Note: We can't easily simulate loaded repositories in tests without actual data loading
	// so this test primarily verifies the initial state behavior
}

func TestAbisCollection_NeedsUpdate_UnknownListKind(t *testing.T) {
	mockApp := mocks.NewMockApp()
	ac := NewAbisCollection(mockApp)

	// Test with unknown ListKind - should return true
	unknownKind := types.ListKind("Unknown")
	if !ac.NeedsUpdate(unknownKind) {
		t.Error("NeedsUpdate with unknown ListKind should return true")
	}
}

func TestAbisCollection_ThreadSafety(t *testing.T) {
	mockApp := mocks.NewMockApp()
	ac := NewAbisCollection(mockApp)

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

func TestAbisRepository(t *testing.T) {
	app := mocks.NewMockApp()

	// Test creating Abis repositories
	downloadedRepo := NewAbisRepository(app, "Downloaded", func(abi *coreTypes.Abi) bool {
		return !abi.IsKnown
	})
	knownRepo := NewAbisRepository(app, "Known", func(abi *coreTypes.Abi) bool {
		return abi.IsKnown
	})
	functionsRepo := NewFunctionsRepository(app, "Functions", func(item *coreTypes.Function) bool {
		return item.FunctionType != "event"
	})
	eventsRepo := NewFunctionsRepository(app, "Events", func(item *coreTypes.Function) bool {
		return item.FunctionType == "event"
	})

	if downloadedRepo == nil {
		t.Error("Downloaded repository should not be nil")
	}

	if knownRepo == nil {
		t.Error("Known repository should not be nil")
	}

	if functionsRepo == nil {
		t.Error("Functions repository should not be nil")
	}

	if eventsRepo == nil {
		t.Error("Events repository should not be nil")
	}

	// Test initial states
	repos := []interface {
		IsLoaded() bool
		NeedsUpdate() bool
		Count() int
	}{downloadedRepo, knownRepo, functionsRepo, eventsRepo}

	for i, repo := range repos {
		if repo.IsLoaded() {
			t.Errorf("Repository %d should not be loaded initially", i)
		}

		if !repo.NeedsUpdate() {
			t.Errorf("Repository %d should need update initially", i)
		}

		if repo.Count() != 0 {
			t.Errorf("Repository %d: expected count 0, got %d", i, repo.Count())
		}
	}

	t.Log("Abis repository test completed successfully")
}

// ADD_ROUTE
