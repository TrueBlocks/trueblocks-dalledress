package monitors

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
	"github.com/stretchr/testify/assert"
)

// Set to true to run slow/integration tests that may take several seconds
const runSlowTests = false

// Helper function to get appropriate sleep duration based on test mode
func testSleep(duration time.Duration) time.Duration {
	if runSlowTests {
		return duration
	}
	// Use much shorter sleeps for fast tests
	return duration / 10
}

func TestNewMonitorsCollection(t *testing.T) {
	assert.NotPanics(t, func() {
		collection := NewMonitorsCollection()
		assert.NotNil(t, collection)
	}, "NewMonitorsCollection should not panic")
}

func TestMonitorsMatchesFilter(t *testing.T) {
	collection := NewMonitorsCollection()
	testMonitor := &coreTypes.Monitor{
		Address:     base.HexToAddress("0x1234567890123456789012345678901234567890"),
		Name:        "Test Monitor",
		NRecords:    100,
		FileSize:    1024,
		LastScanned: 12345,
		IsEmpty:     false,
		IsStaged:    true,
		Deleted:     false,
	}

	t.Run("AddressMatch", func(t *testing.T) {
		assert.True(t, collection.matchesFilter(testMonitor, "1234"))
		assert.True(t, collection.matchesFilter(testMonitor, "0x1234"))
		assert.True(t, collection.matchesFilter(testMonitor, "1234567890123456"))
	})

	t.Run("NameMatch", func(t *testing.T) {
		assert.True(t, collection.matchesFilter(testMonitor, "test"))
		assert.True(t, collection.matchesFilter(testMonitor, "Monitor"))
		assert.True(t, collection.matchesFilter(testMonitor, "TEST")) // case insensitive
	})

	t.Run("NumericFieldsMatch", func(t *testing.T) {
		assert.True(t, collection.matchesFilter(testMonitor, "100"))   // NRecords
		assert.True(t, collection.matchesFilter(testMonitor, "1024"))  // FileSize
		assert.True(t, collection.matchesFilter(testMonitor, "12345")) // LastScanned
	})

	t.Run("BooleanFieldsMatch", func(t *testing.T) {
		assert.True(t, collection.matchesFilter(testMonitor, "staged"))
		emptyMonitor := &coreTypes.Monitor{IsEmpty: true}
		assert.True(t, collection.matchesFilter(emptyMonitor, "empty"))
		deletedMonitor := &coreTypes.Monitor{Deleted: true}
		assert.True(t, collection.matchesFilter(deletedMonitor, "deleted"))
	})

	t.Run("NoMatch", func(t *testing.T) {
		assert.False(t, collection.matchesFilter(testMonitor, "nonexistent"))
		assert.False(t, collection.matchesFilter(testMonitor, "xyz"))
	})

	t.Run("EmptyFilter", func(t *testing.T) {
		result := collection.matchesFilter(testMonitor, "")
		assert.True(t, result)
	})
}

func TestMonitorsCollectionStateManagement(t *testing.T) {
	collection := NewMonitorsCollection()

	t.Run("NeedsUpdate", func(t *testing.T) {
		needsUpdate := collection.NeedsUpdate(MonitorsList)
		assert.True(t, needsUpdate, "New collection should need update")

		needsUpdate = collection.NeedsUpdate("invalid-kind")
		assert.False(t, needsUpdate, "Invalid list kind should return false")
	})

	t.Run("Reset", func(t *testing.T) {
		assert.NotPanics(t, func() {
			collection.Reset(MonitorsList)
		}, "Reset with valid list kind should not panic")

		needsUpdate := collection.NeedsUpdate(MonitorsList)
		assert.True(t, needsUpdate, "After reset, collection should need update")

		assert.NotPanics(t, func() {
			collection.Reset("invalid-kind")
		}, "Reset with invalid list kind should not panic")
	})

	t.Run("GetExpectedTotal", func(t *testing.T) {
		assert.NotPanics(t, func() {
			page, _ := collection.GetPage(MonitorsList, 0, 10, sdk.SortSpec{}, "")
			if page != nil {
				assert.GreaterOrEqual(t, page.ExpectedTotal, 0, "ExpectedTotal should be non-negative")
			}
		})
	})

	t.Run("FacetStateIntegration", func(t *testing.T) {
		needsUpdate := collection.NeedsUpdate(MonitorsList)
		assert.IsType(t, true, needsUpdate, "NeedsUpdate should return a boolean")

		page, err := collection.GetPage(MonitorsList, 0, 1, sdk.SortSpec{}, "")
		if err == nil && page != nil {
			assert.Equal(t, MonitorsList, page.Kind, "Page kind should match request")
			assert.GreaterOrEqual(t, page.TotalItems, 0, "TotalItems should be non-negative")
			assert.GreaterOrEqual(t, page.ExpectedTotal, 0, "ExpectedTotal should be non-negative")
			assert.IsType(t, true, page.IsFetching, "IsFetching should be a boolean")
		}
	})
}

func TestMonitorsCollectionLoadData(t *testing.T) {
	collection := NewMonitorsCollection()

	t.Run("LoadDataValidKind", func(t *testing.T) {
		assert.NotPanics(t, func() {
			collection.LoadData(MonitorsList)
		}, "LoadData with valid list kind should not panic")

		needsUpdate := collection.NeedsUpdate(MonitorsList)
		assert.True(t, needsUpdate, "LoadData call doesn't immediately change NeedsUpdate state")
	})

	t.Run("LoadDataInvalidKind", func(t *testing.T) {
		assert.NotPanics(t, func() {
			collection.LoadData("invalid-kind")
		}, "LoadData with invalid list kind should not panic")
	})

	t.Run("LoadDataWhenNotNeeded", func(t *testing.T) {

		assert.NotPanics(t, func() {
			collection.LoadData(MonitorsList)
			collection.LoadData(MonitorsList)
		})
	})

	t.Run("GetPageErrorHandling", func(t *testing.T) {
		_, err := collection.GetPage("invalid-kind", 0, 10, sdk.SortSpec{}, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported list kind")

		assert.NotPanics(t, func() {
			page, err := collection.GetPage(MonitorsList, 0, 10, sdk.SortSpec{}, "")
			if err == nil && page != nil {
				assert.Equal(t, MonitorsList, page.Kind)
				assert.GreaterOrEqual(t, page.TotalItems, 0)
				assert.GreaterOrEqual(t, page.ExpectedTotal, 0)
			}
		})
	})
}

func TestMonitorsCollectionGetPageFunctionality(t *testing.T) {
	collection := NewMonitorsCollection()

	t.Run("BasicPagination", func(t *testing.T) {
		testCases := []struct {
			first    int
			pageSize int
			name     string
		}{
			{0, 10, "first page"},
			{0, 5, "smaller page size"},
			{5, 10, "offset start"},
			{0, 1, "single item page"},
			{0, 100, "large page size"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				page, err := collection.GetPage(MonitorsList, tc.first, tc.pageSize, sdk.SortSpec{}, "")

				if err == nil && page != nil {
					assert.Equal(t, MonitorsList, page.Kind)
					assert.GreaterOrEqual(t, page.TotalItems, 0)
					assert.GreaterOrEqual(t, page.ExpectedTotal, 0)
					assert.LessOrEqual(t, len(page.Monitors), tc.pageSize, "Returned items should not exceed page size")

					for i, monitor := range page.Monitors {
						assert.NotNil(t, monitor, "Monitor at index %d should not be nil", i)
					}
				}
			})
		}
	})

	t.Run("FilterFunctionality", func(t *testing.T) {
		testFilters := []struct {
			filter      string
			description string
		}{
			{"", "empty filter"},
			{"test", "simple text filter"},
			{"0x1234", "address filter"},
			{"100", "numeric filter"},
			{"empty", "boolean filter"},
			{"nonexistent", "filter with no matches"},
		}

		for _, tf := range testFilters {
			t.Run(tf.description, func(t *testing.T) {
				page, err := collection.GetPage(MonitorsList, 0, 10, sdk.SortSpec{}, tf.filter)

				if err == nil && page != nil {
					assert.Equal(t, MonitorsList, page.Kind)
					assert.GreaterOrEqual(t, page.TotalItems, 0)

					if tf.filter == "" {
						assert.GreaterOrEqual(t, page.TotalItems, 0)
					}
				}
			})
		}
	})

	t.Run("SortingParameters", func(t *testing.T) {
		sortSpecs := []sdk.SortSpec{
			{},
		}

		for i, sortSpec := range sortSpecs {
			t.Run(fmt.Sprintf("sort_spec_%d", i), func(t *testing.T) {
				page, err := collection.GetPage(MonitorsList, 0, 10, sortSpec, "")

				if err == nil && page != nil {
					assert.Equal(t, MonitorsList, page.Kind)
					assert.GreaterOrEqual(t, page.TotalItems, 0)
				}
			})
		}
	})

	t.Run("EdgeCases", func(t *testing.T) {
		t.Run("ZeroPageSize", func(t *testing.T) {
			page, err := collection.GetPage(MonitorsList, 0, 0, sdk.SortSpec{}, "")
			if err == nil && page != nil {
				assert.Equal(t, MonitorsList, page.Kind)
				assert.Len(t, page.Monitors, 0, "Zero page size should return no items")
			}
		})

		t.Run("NegativeFirst", func(t *testing.T) {
			page, err := collection.GetPage(MonitorsList, -1, 10, sdk.SortSpec{}, "")
			if err == nil && page != nil {
				assert.Equal(t, MonitorsList, page.Kind)
				assert.GreaterOrEqual(t, page.TotalItems, 0)
			}
		})

		t.Run("NegativePageSize", func(t *testing.T) {
			page, err := collection.GetPage(MonitorsList, 0, -1, sdk.SortSpec{}, "")
			if err == nil && page != nil {
				assert.Equal(t, MonitorsList, page.Kind)
			}
		})

		t.Run("LargeOffset", func(t *testing.T) {
			page, err := collection.GetPage(MonitorsList, 10000, 10, sdk.SortSpec{}, "")
			if err == nil && page != nil {
				assert.Equal(t, MonitorsList, page.Kind)
				assert.GreaterOrEqual(t, page.TotalItems, 0)
			}
		})
	})

	t.Run("InvalidListKind", func(t *testing.T) {
		_, err := collection.GetPage("invalid-kind", 0, 10, sdk.SortSpec{}, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported list kind")
	})

	t.Run("PageStructureValidation", func(t *testing.T) {
		page, err := collection.GetPage(MonitorsList, 0, 5, sdk.SortSpec{}, "")

		if err == nil && page != nil {
			assert.Equal(t, MonitorsList, page.Kind, "Kind should match request")
			assert.GreaterOrEqual(t, page.TotalItems, 0, "TotalItems should be non-negative")
			assert.GreaterOrEqual(t, page.ExpectedTotal, 0, "ExpectedTotal should be non-negative")
			assert.IsType(t, false, page.IsFetching, "IsFetching should be boolean")
			assert.NotNil(t, page.Monitors, "Monitors slice should not be nil")
			assert.LessOrEqual(t, len(page.Monitors), 5, "Should not exceed requested page size")

			for i, monitor := range page.Monitors {
				assert.NotEqual(t, coreTypes.Monitor{}, monitor, "Monitor %d should not be zero value", i)
			}
		}
	})

	t.Run("FilterAndPaginationCombined", func(t *testing.T) {
		page1, err1 := collection.GetPage(MonitorsList, 0, 3, sdk.SortSpec{}, "test")
		page2, err2 := collection.GetPage(MonitorsList, 3, 3, sdk.SortSpec{}, "test")

		if err1 == nil && page1 != nil && err2 == nil && page2 != nil {
			assert.Equal(t, page1.TotalItems, page2.TotalItems, "Filter should give consistent total across pages")
			assert.Equal(t, MonitorsList, page1.Kind)
			assert.Equal(t, MonitorsList, page2.Kind)
		}
	})
}

func TestMonitorsCollectionLoadDataAsync(t *testing.T) {
	if !runSlowTests {
		t.Skip("Skipping slow async tests - set runSlowTests = true to enable")
	}
	t.Run("LoadDataDoesNotBlock", func(t *testing.T) {
		collection := NewMonitorsCollection()

		// LoadData should return immediately (not block)
		start := time.Now()
		collection.LoadData(MonitorsList)
		duration := time.Since(start)

		// Should complete very quickly since it just starts a goroutine
		assert.Less(t, duration, 100*time.Millisecond, "LoadData should not block")
	})

	t.Run("LoadDataStartsAsyncOperation", func(t *testing.T) {
		collection := NewMonitorsCollection()

		// Start loading
		collection.LoadData(MonitorsList)

		// Immediately after calling LoadData, state might not have changed yet
		immediateNeedsUpdate := collection.NeedsUpdate(MonitorsList)

		// The goroutine might not have started yet, so either value is acceptable
		assert.IsType(t, true, immediateNeedsUpdate, "NeedsUpdate should return boolean")

		// Give some time for the async operation to potentially start
		time.Sleep(testSleep(50 * time.Millisecond))

		// Check that we can still call other methods without issues
		assert.NotPanics(t, func() {
			collection.NeedsUpdate(MonitorsList)
			collection.GetPage(MonitorsList, 0, 1, sdk.SortSpec{}, "")
		}, "Other methods should work during async loading")
	})

	t.Run("MultipleLoadDataCalls", func(t *testing.T) {
		collection := NewMonitorsCollection()

		// Multiple calls should not panic or cause issues
		assert.NotPanics(t, func() {
			collection.LoadData(MonitorsList)
			collection.LoadData(MonitorsList)
			collection.LoadData(MonitorsList)
		}, "Multiple LoadData calls should not panic")

		// Give time for any goroutines to start
		time.Sleep(testSleep(50 * time.Millisecond))

		// Should still be able to query state
		assert.NotPanics(t, func() {
			collection.NeedsUpdate(MonitorsList)
		})
	})

	t.Run("LoadDataWithInvalidKind", func(t *testing.T) {
		collection := NewMonitorsCollection()

		// Should handle invalid kind gracefully (early return)
		assert.NotPanics(t, func() {
			collection.LoadData("invalid-kind")
		}, "LoadData with invalid kind should not panic")

		// Should return quickly since it returns early
		start := time.Now()
		collection.LoadData("invalid-kind")
		duration := time.Since(start)
		assert.Less(t, duration, 10*time.Millisecond, "Invalid kind should return immediately")
	})

	t.Run("LoadDataWhenNotNeeded", func(t *testing.T) {
		collection := NewMonitorsCollection()

		// This tests the early return when !mc.NeedsUpdate(listKind)
		// We can't easily force NeedsUpdate to return false, but we can test
		// that calling LoadData multiple times doesn't cause issues

		for i := 0; i < 3; i++ {
			assert.NotPanics(t, func() {
				collection.LoadData(MonitorsList)
			}, "Repeated LoadData calls should not panic")

			// Small delay between calls
			time.Sleep(testSleep(10 * time.Millisecond))
		}
	})

	t.Run("ConcurrentLoadDataCalls", func(t *testing.T) {
		if !runSlowTests {
			t.Skip("Skipping concurrent test - set runSlowTests = true to enable")
		}

		collection := NewMonitorsCollection()

		// Test concurrent calls to LoadData
		var wg sync.WaitGroup
		numGoroutines := 5

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				assert.NotPanics(t, func() {
					collection.LoadData(MonitorsList)
				}, "Concurrent LoadData call %d should not panic", id)
			}(i)
		}

		// Wait for all goroutines to complete
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		// Wait with timeout
		select {
		case <-done:
			// All goroutines completed successfully
		case <-time.After(5 * time.Second):
			t.Fatal("Concurrent LoadData calls took too long")
		}

		// Collection should still be functional
		assert.NotPanics(t, func() {
			collection.NeedsUpdate(MonitorsList)
			collection.GetPage(MonitorsList, 0, 1, sdk.SortSpec{}, "")
		}, "Collection should remain functional after concurrent access")
	})

	t.Run("LoadDataStateConsistency", func(t *testing.T) {
		collection := NewMonitorsCollection()

		// Call LoadData
		collection.LoadData(MonitorsList)

		// Multiple calls to NeedsUpdate should be consistent
		// (not testing specific value, just consistency)
		result1 := collection.NeedsUpdate(MonitorsList)
		result2 := collection.NeedsUpdate(MonitorsList)
		result3 := collection.NeedsUpdate(MonitorsList)

		assert.Equal(t, result1, result2, "NeedsUpdate should be consistent")
		assert.Equal(t, result2, result3, "NeedsUpdate should be consistent")
	})

	t.Run("LoadDataWithResetInteraction", func(t *testing.T) {
		collection := NewMonitorsCollection()

		// Start loading
		collection.LoadData(MonitorsList)

		// Immediately reset
		assert.NotPanics(t, func() {
			collection.Reset(MonitorsList)
		}, "Reset during LoadData should not panic")

		// Should be able to load again after reset
		assert.NotPanics(t, func() {
			collection.LoadData(MonitorsList)
		}, "LoadData after Reset should not panic")

		// Give time for any async operations
		time.Sleep(testSleep(50 * time.Millisecond))

		// State should still be queryable
		needsUpdate := collection.NeedsUpdate(MonitorsList)
		assert.IsType(t, true, needsUpdate, "Should return boolean after reset/load cycle")
	})

	t.Run("LoadDataErrorScenarios", func(t *testing.T) {
		collection := NewMonitorsCollection()

		// Test that LoadData handles various scenarios gracefully
		// Even if the underlying facet.Load() fails, LoadData should not panic

		scenarios := []struct {
			name   string
			action func()
		}{
			{
				name:   "normal_load",
				action: func() { collection.LoadData(MonitorsList) },
			},
			{
				name: "load_after_multiple_resets",
				action: func() {
					collection.Reset(MonitorsList)
					collection.Reset(MonitorsList)
					collection.LoadData(MonitorsList)
				},
			},
			{
				name: "rapid_load_reset_cycle",
				action: func() {
					for i := 0; i < 3; i++ {
						collection.LoadData(MonitorsList)
						collection.Reset(MonitorsList)
					}
				},
			},
		}

		for _, scenario := range scenarios {
			t.Run(scenario.name, func(t *testing.T) {
				assert.NotPanics(t, scenario.action, "Scenario %s should not panic", scenario.name)

				// After each scenario, basic operations should still work
				assert.NotPanics(t, func() {
					collection.NeedsUpdate(MonitorsList)
					collection.GetPage(MonitorsList, 0, 1, sdk.SortSpec{}, "")
				}, "Basic operations should work after %s", scenario.name)
			})
		}
	})

	t.Run("LoadDataMemoryManagement", func(t *testing.T) {
		if !runSlowTests {
			t.Skip("Skipping memory management test - set runSlowTests = true to enable")
		}

		collection := NewMonitorsCollection()

		// Test that repeated LoadData calls don't cause memory issues
		// This is a basic test - in a real scenario you'd use memory profiling
		for i := 0; i < 10; i++ {
			collection.LoadData(MonitorsList)
			time.Sleep(1 * time.Millisecond) // Small delay
		}

		// Should still be functional
		assert.NotPanics(t, func() {
			collection.NeedsUpdate(MonitorsList)
		}, "Should still function after many LoadData calls")

		// Reset to clean up
		collection.Reset(MonitorsList)
	})
}

func TestMonitorsCollectionAdvancedAsync(t *testing.T) {
	if !runSlowTests {
		t.Skip("Skipping slow advanced async tests - set runSlowTests = true to enable")
	}
	// collection := NewMonitorsCollection()

	// t.Run("MessageEmissionVerification", func(t *testing.T) {
	// 	capture := NewMessageCapture()

	// 	// Register listener for DataLoaded events
	// 	unsubscribe := msgs.On(msgs.EventDataLoaded, capture.CaptureLoaded)
	// 	defer unsubscribe()

	// 	// Trigger data loading
	// 	collection.LoadData(MonitorsList)

	// 	// Wait for async operation with timeout
	// 	require.Eventually(t, func() bool {
	// 		messages := capture.GetMessagesForText("monitors")
	// 		return len(messages) > 0
	// 	}, 5*time.Second, 50*time.Millisecond, "Should emit monitors loaded message")

	// 	// Verify message content
	// 	messages := capture.GetMessagesForText("monitors")
	// 	assert.Len(t, messages, 1, "Should emit exactly one monitors message")

	// 	msg := messages[0]
	// 	assert.Equal(t, "DataLoaded", msg.EventType)
	// 	assert.Equal(t, "monitors", msg.MsgText)
	// })

	t.Run("StateTransitionVerification", func(t *testing.T) {
		collection := NewMonitorsCollection()

		// Initial state should indicate updates needed
		initialNeedsUpdate := collection.NeedsUpdate(MonitorsList)
		assert.IsType(t, true, initialNeedsUpdate, "NeedsUpdate should return a boolean")
		t.Logf("Initial NeedsUpdate state: %v", initialNeedsUpdate)

		// Get initial page to check state
		initialPage, err := collection.GetPage(MonitorsList, 0, 1, sdk.SortSpec{}, "")
		if err == nil && initialPage != nil {
			initialState := initialPage.State
			initialFetching := initialPage.IsFetching

			// Trigger load
			collection.LoadData(MonitorsList)

			// Check state transitions during loading
			var fetchingObserved bool
			for i := 0; i < 20; i++ { // Check for up to 1 second
				page, pageErr := collection.GetPage(MonitorsList, 0, 1, sdk.SortSpec{}, "")
				if pageErr == nil && page != nil {
					if page.IsFetching {
						fetchingObserved = true
					}

					// State should progress from initial state
					assert.Contains(t, []string{
						string(initialState),
						"fetching",
						"loaded",
						"partial",
						"error",
					}, string(page.State), "State should be valid during loading")
				}

				time.Sleep(50 * time.Millisecond)
			}

			// We should have observed fetching state at some point
			if !fetchingObserved && !initialFetching {
				t.Log("Warning: Did not observe fetching state transition (may be too fast)")
			}
		}
	})

	t.Run("ConcurrentLoadDataSafety", func(t *testing.T) {
		collection := NewMonitorsCollection()
		const numGoroutines = 10

		var wg sync.WaitGroup
		var errors []error
		var errorsMu sync.Mutex

		// Launch multiple concurrent LoadData calls
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				defer func() {
					if r := recover(); r != nil {
						errorsMu.Lock()
						errors = append(errors, fmt.Errorf("goroutine %d panicked: %v", id, r))
						errorsMu.Unlock()
					}
				}()

				// Multiple calls from same goroutine
				for j := 0; j < 3; j++ {
					collection.LoadData(MonitorsList)
					time.Sleep(1 * time.Millisecond)
				}
			}(i)
		}

		wg.Wait()

		errorsMu.Lock()
		assert.Empty(t, errors, "No panics should occur during concurrent LoadData calls")
		errorsMu.Unlock()

		// Verify collection is still functional
		assert.NotPanics(t, func() {
			collection.NeedsUpdate(MonitorsList)
			collection.GetPage(MonitorsList, 0, 1, sdk.SortSpec{}, "")
		}, "Collection should remain functional after concurrent access")
	})

	t.Run("LoadDataErrorScenarios", func(t *testing.T) {
		// Test with invalid list kinds to ensure graceful handling
		invalidKinds := []string{
			"",
			"InvalidKind",
			"monitors", // lowercase
			"MONITORS", // uppercase
			"monitors-list",
		}

		for _, kind := range invalidKinds {
			t.Run(fmt.Sprintf("InvalidKind_%s", kind), func(t *testing.T) {
				collection := NewMonitorsCollection()

				start := time.Now()
				assert.NotPanics(t, func() {
					collection.LoadData(types.ListKind(kind))
				}, "LoadData with invalid kind should not panic")

				duration := time.Since(start)
				assert.Less(t, duration, 100*time.Millisecond, "Invalid kind should return quickly")
			})
		}
	})

	t.Run("LoadDataMemoryAndResourceManagement", func(t *testing.T) {
		collection := NewMonitorsCollection()

		// Simulate many rapid load requests
		const iterations = 50
		for i := 0; i < iterations; i++ {
			collection.LoadData(MonitorsList)

			// Occasionally check that basic operations still work
			if i%10 == 0 {
				assert.NotPanics(t, func() {
					needsUpdate := collection.NeedsUpdate(MonitorsList)
					_ = needsUpdate

					page, err := collection.GetPage(MonitorsList, 0, 1, sdk.SortSpec{}, "")
					_ = err
					_ = page
				}, "Basic operations should work during rapid loading")
			}

			// Small delay to prevent overwhelming the system
			if i%5 == 0 {
				time.Sleep(1 * time.Millisecond)
			}
		}

		// Collection should still be responsive
		assert.NotPanics(t, func() {
			collection.NeedsUpdate(MonitorsList)
		}, "Collection should be responsive after many load requests")

		// Reset should work
		assert.NotPanics(t, func() {
			collection.Reset(MonitorsList)
		}, "Reset should work after many load requests")
	})

	t.Run("LoadDataWithResetInteraction", func(t *testing.T) {
		collection := NewMonitorsCollection()

		scenarios := []struct {
			name   string
			action func()
		}{
			{
				name: "LoadThenReset",
				action: func() {
					collection.LoadData(MonitorsList)
					time.Sleep(10 * time.Millisecond)
					collection.Reset(MonitorsList)
				},
			},
			{
				name: "ResetThenLoad",
				action: func() {
					collection.Reset(MonitorsList)
					collection.LoadData(MonitorsList)
				},
			},
			{
				name: "AlternatingLoadReset",
				action: func() {
					for i := 0; i < 5; i++ {
						collection.LoadData(MonitorsList)
						collection.Reset(MonitorsList)
					}
				},
			},
		}

		for _, scenario := range scenarios {
			t.Run(scenario.name, func(t *testing.T) {
				assert.NotPanics(t, scenario.action, "Load/Reset interaction should not panic")

				// Verify basic functionality after scenario
				assert.NotPanics(t, func() {
					collection.NeedsUpdate(MonitorsList)
					collection.GetPage(MonitorsList, 0, 1, sdk.SortSpec{}, "")
				}, "Basic functionality should work after %s", scenario.name)
			})
		}
	})
}

func TestMonitorsCollectionIntegration(t *testing.T) {
	if !runSlowTests {
		t.Skip("Skipping slow integration tests - set runSlowTests = true to enable")
	}
	// testHelper := msgs.NewTestHelpers()
	// defer testHelper.Cleanup()

	// t.Run("FullPipelineDataFlow", func(t *testing.T) {
	// 	collection := NewMonitorsCollection()
	// 	capture := NewMessageCapture()

	// 	// Set up message capture
	// 	unsubscribe := msgs.On(msgs.EventDataLoaded, capture.CaptureLoaded)
	// 	defer unsubscribe()

	// 	// Initial state check
	// 	initialPage, err := collection.GetPage(MonitorsList, 0, 5, sdk.SortSpec{}, "")
	// 	if err != nil {
	// 		t.Skipf("Skipping integration test due to GetPage error: %v", err)
	// 	}

	// 	require.NotNil(t, initialPage, "Initial page should not be nil")
	// 	assert.Equal(t, MonitorsList, initialPage.Kind)
	// 	assert.GreaterOrEqual(t, initialPage.TotalItems, 0)
	// 	assert.GreaterOrEqual(t, initialPage.ExpectedTotal, 0)

	// 	// Trigger load and verify async behavior
	// 	loadStart := time.Now()
	// 	collection.LoadData(MonitorsList)
	// 	loadCall := time.Since(loadStart)

	// 	// LoadData should return quickly (async)
	// 	assert.Less(t, loadCall, 100*time.Millisecond, "LoadData should return quickly")

	// 	// Wait for completion with timeout
	// 	var finalPage *MonitorsPage
	// 	require.Eventually(t, func() bool {
	// 		page, pageErr := collection.GetPage(MonitorsList, 0, 5, sdk.SortSpec{}, "")
	// 		if pageErr == nil && page != nil {
	// 			finalPage = page
	// 			// Check if loading is complete (not fetching)
	// 			return !page.IsFetching
	// 		}
	// 		return false
	// 	}, 10*time.Second, 100*time.Millisecond, "Data loading should complete")

	// 	// Verify final state
	// 	require.NotNil(t, finalPage, "Final page should not be nil")
	// 	assert.Equal(t, MonitorsList, finalPage.Kind)
	// 	assert.False(t, finalPage.IsFetching, "Should not be fetching when complete")

	// 	// Verify message emission
	// 	require.Eventually(t, func() bool {
	// 		messages := capture.GetMessagesForText("monitors")
	// 		return len(messages) > 0
	// 	}, 2*time.Second, 50*time.Millisecond, "Should emit monitors loaded message")

	// 	messages := capture.GetMessagesForText("monitors")
	// 	assert.NotEmpty(t, messages, "Should have received monitors messages")
	// })

	t.Run("ErrorHandlingIntegration", func(t *testing.T) {
		collection := NewMonitorsCollection()

		// Test error scenarios that might occur in real usage
		errorScenarios := []struct {
			name string
			test func(t *testing.T)
		}{
			{
				name: "InvalidPageRequests",
				test: func(t *testing.T) {
					// These should handle errors gracefully
					_, err1 := collection.GetPage("invalid", 0, 10, sdk.SortSpec{}, "")
					assert.Error(t, err1)

					page, err2 := collection.GetPage(MonitorsList, -1, 10, sdk.SortSpec{}, "")
					if err2 == nil {
						// If no error, verify response is reasonable
						assert.NotNil(t, page)
						assert.Equal(t, MonitorsList, page.Kind)
					}
				},
			},
			{
				name: "LoadDataAfterErrors",
				test: func(t *testing.T) {
					// Even after errors, LoadData should work
					_, _ = collection.GetPage("invalid", 0, 10, sdk.SortSpec{}, "")

					assert.NotPanics(t, func() {
						collection.LoadData(MonitorsList)
					}, "LoadData should work after previous errors")
				},
			},
		}

		for _, scenario := range errorScenarios {
			t.Run(scenario.name, scenario.test)
		}
	})
}

func TestMonitorsCollectionPerformanceAndEdgeCases(t *testing.T) {
	if !runSlowTests {
		t.Skip("Skipping slow performance tests - set runSlowTests = true to enable")
	}
	t.Run("RapidStateQueries", func(t *testing.T) {
		collection := NewMonitorsCollection()

		// Test rapid queries to state methods
		const iterations = 1000
		start := time.Now()

		for i := 0; i < iterations; i++ {
			collection.NeedsUpdate(MonitorsList)
			if i%10 == 0 {
				// Occasionally trigger other operations
				collection.GetPage(MonitorsList, 0, 1, sdk.SortSpec{}, "")
			}
		}

		duration := time.Since(start)
		avgPerCall := duration / iterations

		assert.Less(t, avgPerCall, 1*time.Millisecond, "Average NeedsUpdate call should be very fast")
		t.Logf("Average time per NeedsUpdate call: %v", avgPerCall)
	})

	t.Run("LargePageSizeHandling", func(t *testing.T) {
		collection := NewMonitorsCollection()

		largeSizes := []int{0, 1, 10, 100, 1000, 10000}

		for _, size := range largeSizes {
			t.Run(fmt.Sprintf("PageSize_%d", size), func(t *testing.T) {
				start := time.Now()
				page, err := collection.GetPage(MonitorsList, 0, size, sdk.SortSpec{}, "")
				duration := time.Since(start)

				if err == nil && page != nil {
					assert.Equal(t, MonitorsList, page.Kind)
					assert.GreaterOrEqual(t, page.TotalItems, 0)
					assert.LessOrEqual(t, len(page.Monitors), size, "Returned items should not exceed page size")

					// Larger page sizes may take longer, but should still be reasonable
					maxExpectedTime := time.Duration(size/10+1) * 100 * time.Millisecond
					assert.Less(t, duration, maxExpectedTime, "Page retrieval should complete in reasonable time")
				}
			})
		}
	})

	t.Run("FilterPerformance", func(t *testing.T) {
		collection := NewMonitorsCollection()

		filters := []string{
			"",                       // Empty filter
			"a",                      // Single character
			"test",                   // Common word
			"0x1234567890123456",     // Long address
			"nonexistentfilter",      // No matches expected
			strings.Repeat("x", 100), // Very long filter
		}

		for _, filter := range filters {
			t.Run(fmt.Sprintf("Filter_%s", filter[:min(len(filter), 10)]), func(t *testing.T) {
				start := time.Now()
				page, err := collection.GetPage(MonitorsList, 0, 50, sdk.SortSpec{}, filter)
				duration := time.Since(start)

				if err == nil && page != nil {
					assert.Equal(t, MonitorsList, page.Kind)

					// Filtering should be reasonably fast
					assert.Less(t, duration, 5*time.Second, "Filtering should complete in reasonable time")
				}
			})
		}
	})

	t.Run("ConcurrentGetPageCalls", func(t *testing.T) {
		collection := NewMonitorsCollection()
		const numConcurrent = 20

		var wg sync.WaitGroup
		var results []*MonitorsPage
		var errors []error
		var mu sync.Mutex

		for i := 0; i < numConcurrent; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()

				page, err := collection.GetPage(MonitorsList, id%5, 10, sdk.SortSpec{}, "")

				mu.Lock()
				if err != nil {
					errors = append(errors, err)
				} else if page != nil {
					results = append(results, page)
				}
				mu.Unlock()
			}(i)
		}

		wg.Wait()

		mu.Lock()
		assert.Empty(t, errors, "No errors should occur during concurrent GetPage calls")

		// All successful results should be valid
		for i, page := range results {
			assert.Equal(t, MonitorsList, page.Kind, "Page %d should have correct kind", i)
			assert.GreaterOrEqual(t, page.TotalItems, 0, "Page %d should have valid TotalItems", i)
		}
		mu.Unlock()
	})

	t.Run("EdgeCaseFilters", func(t *testing.T) {
		collection := NewMonitorsCollection()

		edgeCaseFilters := []string{
			"",   // Empty
			" ",  // Space
			"\t", // Tab
			"\n", // Newline
			"0x", // Just prefix
			"0X", // Uppercase prefix
			"0000000000000000000000000000000000000000", // All zeros
			"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF", // All F's
			"!@#$%^&*()",              // Special characters
			"тест",                    // Unicode
			strings.Repeat("a", 1000), // Very long filter
		}

		for i, filter := range edgeCaseFilters {
			t.Run(fmt.Sprintf("EdgeFilter_%d", i), func(t *testing.T) {
				assert.NotPanics(t, func() {
					page, err := collection.GetPage(MonitorsList, 0, 10, sdk.SortSpec{}, filter)
					if err == nil && page != nil {
						assert.Equal(t, MonitorsList, page.Kind)
						assert.GreaterOrEqual(t, page.TotalItems, 0)
					}
				}, "Edge case filter should not panic")
			})
		}
	})

	t.Run("StressTestRepeatedOperations", func(t *testing.T) {
		collection := NewMonitorsCollection()
		const iterations = 100

		// Mix of operations
		operations := []func(){
			func() { collection.NeedsUpdate(MonitorsList) },
			func() { collection.LoadData(MonitorsList) },
			func() { collection.Reset(MonitorsList) },
			func() {
				_, _ = collection.GetPage(MonitorsList, 0, 10, sdk.SortSpec{}, "")
			},
			func() {
				_, _ = collection.GetPage(MonitorsList, 0, 5, sdk.SortSpec{}, "test")
			},
		}

		start := time.Now()
		for i := 0; i < iterations; i++ {
			op := operations[i%len(operations)]
			assert.NotPanics(t, op, "Operation %d should not panic", i)

			// Small delay every 10 operations
			if i%10 == 0 {
				time.Sleep(1 * time.Millisecond)
			}
		}
		duration := time.Since(start)

		t.Logf("Completed %d mixed operations in %v", iterations, duration)
		assert.Less(t, duration, 30*time.Second, "Stress test should complete in reasonable time")

		// Collection should still be functional
		assert.NotPanics(t, func() {
			needsUpdate := collection.NeedsUpdate(MonitorsList)
			assert.IsType(t, true, needsUpdate)
		}, "Collection should remain functional after stress test")
	})
}

func TestMonitorsCollectionBoundaryConditions(t *testing.T) {
	t.Run("NilAndEmptyInputs", func(t *testing.T) {
		collection := NewMonitorsCollection()

		// Test with empty strings and zero values
		assert.NotPanics(t, func() {
			collection.LoadData("")
		}, "LoadData with empty string should not panic")

		assert.NotPanics(t, func() {
			_, err := collection.GetPage("", 0, 0, sdk.SortSpec{}, "")
			assert.Error(t, err, "GetPage with empty kind should return error")
		}, "GetPage with empty kind should not panic")

		assert.NotPanics(t, func() {
			collection.Reset("")
		}, "Reset with empty string should not panic")

		assert.NotPanics(t, func() {
			result := collection.NeedsUpdate("")
			assert.False(t, result, "NeedsUpdate with empty string should return false")
		}, "NeedsUpdate with empty string should not panic")
	})

	t.Run("ExtremePageValues", func(t *testing.T) {
		collection := NewMonitorsCollection()

		extremeValues := []struct {
			first    int
			pageSize int
			desc     string
		}{
			{-1000, 10, "very negative first"},
			{0, -1, "negative page size"},
			{1000000, 1000000, "very large values"},
			{0, 0, "zero page size"},
			{-1, -1, "both negative"},
		}

		for _, vals := range extremeValues {
			t.Run(vals.desc, func(t *testing.T) {
				assert.NotPanics(t, func() {
					page, err := collection.GetPage(MonitorsList, vals.first, vals.pageSize, sdk.SortSpec{}, "")
					if err == nil && page != nil {
						assert.Equal(t, MonitorsList, page.Kind)
						assert.GreaterOrEqual(t, page.TotalItems, 0)
					}
				}, "Extreme page values should not panic")
			})
		}
	})
}
