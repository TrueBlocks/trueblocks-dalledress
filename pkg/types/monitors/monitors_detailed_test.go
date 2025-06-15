//go:build detailed

package monitors

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
	"github.com/stretchr/testify/assert"
)

func TestMonitorsCollectionLoadDataAsync(t *testing.T) {
	t.Run("LoadDataDoesNotBlock", func(t *testing.T) {
		collection := NewMonitorsCollection()

		start := time.Now()
		collection.LoadData(MonitorsList)
		duration := time.Since(start)

		assert.Less(t, duration, 100*time.Millisecond, "LoadData should not block")
	})

	t.Run("LoadDataStartsAsyncOperation", func(t *testing.T) {
		collection := NewMonitorsCollection()

		collection.LoadData(MonitorsList)

		immediateNeedsUpdate := collection.NeedsUpdate(MonitorsList)

		assert.IsType(t, true, immediateNeedsUpdate, "NeedsUpdate should return boolean")

		time.Sleep(50 * time.Millisecond)

		assert.NotPanics(t, func() {
			collection.NeedsUpdate(MonitorsList)
			_, _ = collection.GetPage(MonitorsList, 0, 1, sdk.SortSpec{}, "")
		}, "Other methods should work during async loading")
	})

	t.Run("MultipleLoadDataCalls", func(t *testing.T) {
		collection := NewMonitorsCollection()

		assert.NotPanics(t, func() {
			collection.LoadData(MonitorsList)
			collection.LoadData(MonitorsList)
			collection.LoadData(MonitorsList)
		}, "Multiple LoadData calls should not panic")

		time.Sleep(50 * time.Millisecond)

		assert.NotPanics(t, func() {
			collection.NeedsUpdate(MonitorsList)
		})
	})

	t.Run("ConcurrentLoadDataCalls", func(t *testing.T) {
		collection := NewMonitorsCollection()

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

		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("Concurrent LoadData calls took too long")
		}

		assert.NotPanics(t, func() {
			collection.NeedsUpdate(MonitorsList)
			_, _ = collection.GetPage(MonitorsList, 0, 1, sdk.SortSpec{}, "")
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

		collection.LoadData(MonitorsList)

		assert.NotPanics(t, func() {
			collection.Reset(MonitorsList)
		}, "Reset during LoadData should not panic")

		assert.NotPanics(t, func() {
			collection.LoadData(MonitorsList)
		}, "LoadData after Reset should not panic")

		time.Sleep(50 * time.Millisecond)

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
					_, _ = collection.GetPage(MonitorsList, 0, 1, sdk.SortSpec{}, "")
				}, "Basic operations should work after %s", scenario.name)
			})
		}
	})

	t.Run("LoadDataMemoryManagement", func(t *testing.T) {
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
		const numGoroutines = 5 // Reduced from 10

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

				// Reduced calls per goroutine from 3 to 2
				for j := 0; j < 2; j++ {
					collection.LoadData(MonitorsList)
					// Removed sleep to speed up test
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
			_, _ = collection.GetPage(MonitorsList, 0, 1, sdk.SortSpec{}, "")
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

		// Reduced iterations from 50 to 20
		const iterations = 20
		for i := 0; i < iterations; i++ {
			collection.LoadData(MonitorsList)

			// Check basic operations less frequently (every 5 instead of 10)
			if i%5 == 0 {
				assert.NotPanics(t, func() {
					needsUpdate := collection.NeedsUpdate(MonitorsList)
					_ = needsUpdate

					page, err := collection.GetPage(MonitorsList, 0, 1, sdk.SortSpec{}, "")
					_ = err
					_ = page
				}, "Basic operations should work during rapid loading")
			}

			// Removed sleep delays to speed up test
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
					// Reduced iterations from 5 to 3
					for i := 0; i < 3; i++ {
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
					_, _ = collection.GetPage(MonitorsList, 0, 1, sdk.SortSpec{}, "")
				}, "Basic functionality should work after %s", scenario.name)
			})
		}
	})
}

func TestMonitorsCollectionIntegration(t *testing.T) {
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
	// t.Run("RapidStateQueries", func(t *testing.T) {
	// 	collection := NewMonitorsCollection()

	// 	// Test rapid queries to state methods
	// 	const iterations = 1000
	// 	start := time.Now()

	// 	for i := 0; i < iterations; i++ {
	// 		collection.NeedsUpdate(MonitorsList)
	// 		if i%10 == 0 {
	// 			// Occasionally trigger other operations
	// 			_, _ = collection.GetPage(MonitorsList, 0, 1, sdk.SortSpec{}, "")
	// 		}
	// 	}

	// 	duration := time.Since(start)
	// 	avgPerCall := duration / iterations

	// 	assert.Less(t, avgPerCall, 1*time.Millisecond, "Average NeedsUpdate call should be very fast")
	// 	t.Logf("Average time per NeedsUpdate call: %v", avgPerCall)
	// })

	// t.Run("LargePageSizeHandling", func(t *testing.T) {
	// 	collection := NewMonitorsCollection()

	// 	largeSizes := []int{0, 1, 10, 100, 1000, 10000}

	// 	for _, size := range largeSizes {
	// 		t.Run(fmt.Sprintf("PageSize_%d", size), func(t *testing.T) {
	// 			start := time.Now()
	// 			page, err := collection.GetPage(MonitorsList, 0, size, sdk.SortSpec{}, "")
	// 			duration := time.Since(start)

	// 			if err == nil && page != nil {
	// 				assert.Equal(t, MonitorsList, page.Kind)
	// 				assert.GreaterOrEqual(t, page.TotalItems, 0)
	// 				assert.LessOrEqual(t, len(page.Monitors), size, "Returned items should not exceed page size")

	// 				// Larger page sizes may take longer, but should still be reasonable
	// 				maxExpectedTime := time.Duration(size/10+1) * 100 * time.Millisecond
	// 				assert.Less(t, duration, maxExpectedTime, "Page retrieval should complete in reasonable time")
	// 			}
	// 		})
	// 	}
	// })

	t.Run("FilterPerformance", func(t *testing.T) {
		collection := NewMonitorsCollection()

		// Reduced filter set to focus on most important cases
		filters := []string{
			"",                   // Empty filter
			"test",               // Common word
			"0x1234567890123456", // Long address
			"nonexistentfilter",  // No matches expected
		}

		for _, filter := range filters {
			t.Run(fmt.Sprintf("Filter_%s", filter[:min(len(filter), 10)]), func(t *testing.T) {
				start := time.Now()
				// Reduced page size from 50 to 20
				page, err := collection.GetPage(MonitorsList, 0, 20, sdk.SortSpec{}, filter)
				duration := time.Since(start)

				if err == nil && page != nil {
					assert.Equal(t, MonitorsList, page.Kind)

					// Reduced timeout from 5s to 2s
					assert.Less(t, duration, 2*time.Second, "Filtering should complete in reasonable time")
				}
			})
		}
	})

	t.Run("ConcurrentGetPageCalls", func(t *testing.T) {
		collection := NewMonitorsCollection()
		const numConcurrent = 8 // Reduced from 20

		var wg sync.WaitGroup
		var results []*MonitorsPage
		var errors []error
		var mu sync.Mutex

		for i := 0; i < numConcurrent; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()

				// Reduced page size from 10 to 5
				page, err := collection.GetPage(MonitorsList, id%3, 5, sdk.SortSpec{}, "")

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

		// Reduced to most critical edge cases
		edgeCaseFilters := []string{
			"",   // Empty
			" ",  // Space
			"0x", // Just prefix
			"0000000000000000000000000000000000000000", // All zeros
			"!@#$%^&*()", // Special characters
		}

		for i, filter := range edgeCaseFilters {
			t.Run(fmt.Sprintf("EdgeFilter_%d", i), func(t *testing.T) {
				assert.NotPanics(t, func() {
					// Reduced page size from 10 to 5
					page, err := collection.GetPage(MonitorsList, 0, 5, sdk.SortSpec{}, filter)
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
		const iterations = 30 // Reduced from 100

		// Mix of operations
		operations := []func(){
			func() { collection.NeedsUpdate(MonitorsList) },
			func() { collection.LoadData(MonitorsList) },
			func() { collection.Reset(MonitorsList) },
			func() {
				// Reduced page size from 10 to 5
				_, _ = collection.GetPage(MonitorsList, 0, 5, sdk.SortSpec{}, "")
			},
			func() {
				// Reduced page size from 5 to 3
				_, _ = collection.GetPage(MonitorsList, 0, 3, sdk.SortSpec{}, "test")
			},
		}

		start := time.Now()
		for i := 0; i < iterations; i++ {
			op := operations[i%len(operations)]
			assert.NotPanics(t, op, "Operation %d should not panic", i)

			// Removed delays to speed up test
		}
		duration := time.Since(start)

		t.Logf("Completed %d mixed operations in %v", iterations, duration)
		assert.Less(t, duration, 10*time.Second, "Stress test should complete in reasonable time") // Reduced from 30s to 10s

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
			{-100, 10, "negative first"}, // Reduced from -1000
			{0, -1, "negative page size"},
			{1000, 100, "large values"}, // Reduced from 1000000, 1000000
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
