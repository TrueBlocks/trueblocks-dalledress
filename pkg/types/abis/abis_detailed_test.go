//go:build detailed

package abis

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
	"github.com/stretchr/testify/assert"
)

func TestAbisCollectionLoadDataAsync(t *testing.T) {
	t.Run("LoadDataDoesNotBlock", func(t *testing.T) {
		collection := NewAbisCollection()

		// Test that LoadData call completes without panicking
		assert.NotPanics(t, func() {
			collection.LoadData(AbisDownloaded)
		})
	})

	t.Run("LoadDataStartsAsyncOperation", func(t *testing.T) {
		collection := NewAbisCollection()

		collection.LoadData(AbisDownloaded)

		needsUpdate := collection.NeedsUpdate(AbisDownloaded)

		_ = needsUpdate

		time.Sleep(10 * time.Millisecond)

		_, err := collection.GetPage(AbisDownloaded, 0, 1, sdk.SortSpec{}, "")
		assert.NoError(t, err, "Should be able to call GetPage after LoadData")
	})

	t.Run("MultipleLoadDataCalls", func(t *testing.T) {
		collection := NewAbisCollection()

		assert.NotPanics(t, func() {
			collection.LoadData(AbisDownloaded)
			collection.LoadData(AbisKnown)
			collection.LoadData(AbisFunctions)
			collection.LoadData(AbisEvents)
		}, "Multiple LoadData calls should not panic")

		time.Sleep(50 * time.Millisecond)

		needsUpdate := collection.NeedsUpdate(AbisDownloaded)
		_ = needsUpdate
	})

	t.Run("ConcurrentLoadDataCalls", func(t *testing.T) {
		collection := NewAbisCollection()

		var wg sync.WaitGroup
		numGoroutines := 3

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				time.Sleep(time.Duration(index*10) * time.Millisecond)
				collection.LoadData([]types.DataFacet{AbisDownloaded, AbisKnown, AbisFunctions, AbisEvents}[index%4])
			}(i)
		}

		done := make(chan bool)
		go func() {
			wg.Wait()
			done <- true
		}()

		select {
		case <-done:
		case <-time.After(10 * time.Second):
			t.Fatal("Timed out waiting for concurrent LoadData calls to complete")
		}

		_, err := collection.GetPage(AbisDownloaded, 0, 1, sdk.SortSpec{}, "")
		assert.NoError(t, err, "Collection should be functional after concurrent operations")
	})
}

func TestAbisCollectionAdvancedAsync(t *testing.T) {
	t.Run("LoadDataWhileGettingPages", func(t *testing.T) {
		collection := NewAbisCollection()

		var wg sync.WaitGroup
		errors := make(chan error, 10)

		for i := 0; i < 2; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < 3; j++ {
					_, err := collection.GetPage(AbisDownloaded, 0, 5, sdk.SortSpec{}, "")
					if err != nil {
						errors <- err
						return
					}
					time.Sleep(20 * time.Millisecond)
				}
			}()
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 2; j++ {
				collection.LoadData(AbisDownloaded)
				time.Sleep(50 * time.Millisecond)
			}
		}()

		done := make(chan bool)
		go func() {
			wg.Wait()
			done <- true
		}()

		select {
		case <-done:
			close(errors)
			for err := range errors {
				assert.NoError(t, err, "Concurrent operations should not cause errors")
			}
		case <-time.After(15 * time.Second):
			t.Fatal("Timed out waiting for concurrent operations")
		}
	})

	t.Run("ResetDuringOperations", func(t *testing.T) {
		collection := NewAbisCollection()

		var wg sync.WaitGroup
		stopFlag := make(chan bool)

		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stopFlag:
					return
				default:
					collection.LoadData(AbisDownloaded)
					_, _ = collection.GetPage(AbisDownloaded, 0, 1, sdk.SortSpec{}, "")
					time.Sleep(10 * time.Millisecond)
				}
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 3; i++ {
				collection.Reset(AbisDownloaded)
				time.Sleep(20 * time.Millisecond)
			}
		}()

		time.Sleep(200 * time.Millisecond)
		close(stopFlag)

		done := make(chan bool)
		go func() {
			wg.Wait()
			done <- true
		}()

		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("Timed out waiting for reset operations")
		}

		assert.NotPanics(t, func() {
			collection.NeedsUpdate(AbisDownloaded)
		}, "Collection should be functional after concurrent reset operations")
	})
}

func TestAbisCollectionIntegration(t *testing.T) {
	t.Run("FullWorkflowAllDataFacets", func(t *testing.T) {
		collection := NewAbisCollection()
		dataFacets := []types.DataFacet{AbisDownloaded, AbisKnown, AbisFunctions, AbisEvents}

		for _, dataFacet := range dataFacets {
			t.Run(fmt.Sprintf("Workflow_%s", dataFacet), func(t *testing.T) {
				needsUpdate := collection.NeedsUpdate(dataFacet)
				assert.True(t, needsUpdate, "Should need update initially")

				collection.LoadData(dataFacet)

				page1, err := collection.GetPage(dataFacet, 0, 5, sdk.SortSpec{}, "")
				if err == nil && page1 != nil {
					// Cast to concrete type to access fields
					abisPage1, ok := page1.(*AbisPage)
					assert.True(t, ok, "Expected *AbisPage type")
					assert.Equal(t, dataFacet, abisPage1.Facet)
					assert.GreaterOrEqual(t, abisPage1.TotalItems, 0)

					if abisPage1.TotalItems > 5 {
						page2, err := collection.GetPage(dataFacet, 5, 5, sdk.SortSpec{}, "")
						if err == nil && page2 != nil {
							// Cast to concrete type to access fields
							abisPage2, ok := page2.(*AbisPage)
							assert.True(t, ok, "Expected *AbisPage type")
							assert.Equal(t, dataFacet, abisPage2.Facet)
						}
					}

					if abisPage1.TotalItems > 0 {
						filteredPage, err := collection.GetPage(dataFacet, 0, 5, sdk.SortSpec{}, "test")
						if err == nil && filteredPage != nil {
							// Cast to concrete type to access fields
							abisFilteredPage, ok := filteredPage.(*AbisPage)
							assert.True(t, ok, "Expected *AbisPage type")
							assert.Equal(t, dataFacet, abisFilteredPage.Facet)
							assert.LessOrEqual(t, abisFilteredPage.TotalItems, abisPage1.TotalItems)
						}
					}
				}

				collection.Reset(dataFacet)
				needsUpdate = collection.NeedsUpdate(dataFacet)
				assert.True(t, needsUpdate, "Should need update after reset")
			})
		}
	})

	t.Run("MixedOperationsAllFacets", func(t *testing.T) {
		collection := NewAbisCollection()

		collection.LoadData(AbisDownloaded)
		collection.LoadData(AbisKnown)
		collection.LoadData(AbisFunctions)
		collection.LoadData(AbisEvents)

		time.Sleep(50 * time.Millisecond)

		pages := make([]*AbisPage, 4)
		dataFacets := []types.DataFacet{AbisDownloaded, AbisKnown, AbisFunctions, AbisEvents}

		for i, dataFacet := range dataFacets {
			page, err := collection.GetPage(dataFacet, 0, 2, sdk.SortSpec{}, "")
			if err == nil {
				// Cast to concrete type to access fields
				abisPage, ok := page.(*AbisPage)
				assert.True(t, ok, "Expected *AbisPage type")
				pages[i] = abisPage
				assert.Equal(t, dataFacet, abisPage.Facet)
			}
		}

		collection.Reset(AbisDownloaded)

		needsUpdate := collection.NeedsUpdate(AbisDownloaded)
		assert.True(t, needsUpdate, "AbisDownloaded should need update after reset")

		_, err := collection.GetPage(AbisFunctions, 0, 1, sdk.SortSpec{}, "")
		assert.NoError(t, err, "AbisFunctions should still be accessible")
	})
}

func TestAbisCollectionBoundaryConditions(t *testing.T) {
	t.Run("AllDataFacetsCoverage", func(t *testing.T) {
		collection := NewAbisCollection()
		validFacets := []types.DataFacet{AbisDownloaded, AbisKnown, AbisFunctions, AbisEvents}

		for _, facet := range validFacets {
			t.Run(fmt.Sprintf("Facet_%s", facet), func(t *testing.T) {
				assert.NotPanics(t, func() {
					collection.LoadData(facet)
				}, "LoadData should not panic with valid facet")

				assert.NotPanics(t, func() {
					collection.Reset(facet)
				}, "Reset should not panic with valid facet")

				needsUpdate := collection.NeedsUpdate(facet)
				assert.True(t, needsUpdate, "Should need update for valid facet")

				_, err := collection.GetPage(facet, 0, 5, sdk.SortSpec{}, "")
				_ = err
			})
		}
	})

	t.Run("InvalidDataFacets", func(t *testing.T) {
		collection := NewAbisCollection()
		invalidFacets := []string{"", "invalid", "Monitors", "Unknown"}

		for _, facet := range invalidFacets {
			t.Run(fmt.Sprintf("InvalidFacet_%s", facet), func(t *testing.T) {
				assert.NotPanics(t, func() {
					collection.LoadData(types.DataFacet(facet))
				}, "LoadData should not panic with invalid facet")

				assert.NotPanics(t, func() {
					collection.Reset(types.DataFacet(facet))
				}, "Reset should not panic with invalid facet")

				needsUpdate := collection.NeedsUpdate(types.DataFacet(facet))
				assert.False(t, needsUpdate, "Should not need update for invalid facet")

				_, err := collection.GetPage(types.DataFacet(facet), 0, 5, sdk.SortSpec{}, "")
				assert.Error(t, err, "GetPage should return error for invalid facet")
				assert.Contains(t, err.Error(), "unexpected dataFacet", "Error should mention unexpected dataFacet")
			})
		}
	})

	t.Run("RapidResetOperations", func(t *testing.T) {
		collection := NewAbisCollection()

		assert.NotPanics(t, func() {
			for i := 0; i < 20; i++ {
				collection.Reset(AbisDownloaded)
				collection.Reset(AbisFunctions)
				time.Sleep(1 * time.Millisecond)
			}
		}, "Rapid reset operations should not panic")

		needsUpdate := collection.NeedsUpdate(AbisDownloaded)
		assert.True(t, needsUpdate, "Collection should be functional after rapid resets")
	})

	t.Run("EmptyStringFilters", func(t *testing.T) {
		collection := NewAbisCollection()

		specialFilters := []string{"", " ", "\t", "\n", "\r\n", "   "}

		for i, filter := range specialFilters {
			t.Run(fmt.Sprintf("SpecialFilter_%d", i), func(t *testing.T) {
				page, err := collection.GetPage(AbisDownloaded, 0, 5, sdk.SortSpec{}, filter)
				if err == nil && page != nil {
					// Cast to concrete type to access fields
					abisPage, ok := page.(*AbisPage)
					assert.True(t, ok, "Expected *AbisPage type")
					assert.Equal(t, AbisDownloaded, abisPage.Facet)
					assert.GreaterOrEqual(t, abisPage.TotalItems, 0)
				}
			})
		}
	})
}
