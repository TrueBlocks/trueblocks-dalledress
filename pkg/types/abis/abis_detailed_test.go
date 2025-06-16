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
				collection.LoadData([]types.ListKind{AbisDownloaded, AbisKnown, AbisFunctions, AbisEvents}[index%4])
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
	t.Run("FullWorkflowAllListKinds", func(t *testing.T) {
		collection := NewAbisCollection()
		listKinds := []types.ListKind{AbisDownloaded, AbisKnown, AbisFunctions, AbisEvents}

		for _, listKind := range listKinds {
			t.Run(fmt.Sprintf("Workflow_%s", listKind), func(t *testing.T) {
				needsUpdate := collection.NeedsUpdate(listKind)
				assert.True(t, needsUpdate, "Should need update initially")

				collection.LoadData(listKind)

				page1, err := collection.GetPage(listKind, 0, 5, sdk.SortSpec{}, "")
				if err == nil && page1 != nil {
					// Cast to concrete type to access fields
					abisPage1, ok := page1.(*AbisPage)
					assert.True(t, ok, "Expected *AbisPage type")
					assert.Equal(t, listKind, abisPage1.Kind)
					assert.GreaterOrEqual(t, abisPage1.TotalItems, 0)

					if abisPage1.TotalItems > 5 {
						page2, err := collection.GetPage(listKind, 5, 5, sdk.SortSpec{}, "")
						if err == nil && page2 != nil {
							// Cast to concrete type to access fields
							abisPage2, ok := page2.(*AbisPage)
							assert.True(t, ok, "Expected *AbisPage type")
							assert.Equal(t, listKind, abisPage2.Kind)
						}
					}

					if abisPage1.TotalItems > 0 {
						filteredPage, err := collection.GetPage(listKind, 0, 5, sdk.SortSpec{}, "test")
						if err == nil && filteredPage != nil {
							// Cast to concrete type to access fields
							abisFilteredPage, ok := filteredPage.(*AbisPage)
							assert.True(t, ok, "Expected *AbisPage type")
							assert.Equal(t, listKind, abisFilteredPage.Kind)
							assert.LessOrEqual(t, abisFilteredPage.TotalItems, abisPage1.TotalItems)
						}
					}
				}

				collection.Reset(listKind)
				needsUpdate = collection.NeedsUpdate(listKind)
				assert.True(t, needsUpdate, "Should need update after reset")
			})
		}
	})

	t.Run("MixedOperationsAllKinds", func(t *testing.T) {
		collection := NewAbisCollection()

		collection.LoadData(AbisDownloaded)
		collection.LoadData(AbisKnown)
		collection.LoadData(AbisFunctions)
		collection.LoadData(AbisEvents)

		time.Sleep(50 * time.Millisecond)

		pages := make([]*AbisPage, 4)
		listKinds := []types.ListKind{AbisDownloaded, AbisKnown, AbisFunctions, AbisEvents}

		for i, listKind := range listKinds {
			page, err := collection.GetPage(listKind, 0, 2, sdk.SortSpec{}, "")
			if err == nil {
				// Cast to concrete type to access fields
				abisPage, ok := page.(*AbisPage)
				assert.True(t, ok, "Expected *AbisPage type")
				pages[i] = abisPage
				assert.Equal(t, listKind, abisPage.Kind)
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
	t.Run("AllListKindsCoverage", func(t *testing.T) {
		collection := NewAbisCollection()
		validKinds := []types.ListKind{AbisDownloaded, AbisKnown, AbisFunctions, AbisEvents}

		for _, kind := range validKinds {
			t.Run(fmt.Sprintf("Kind_%s", kind), func(t *testing.T) {
				assert.NotPanics(t, func() {
					collection.LoadData(kind)
				}, "LoadData should not panic with valid kind")

				assert.NotPanics(t, func() {
					collection.Reset(kind)
				}, "Reset should not panic with valid kind")

				needsUpdate := collection.NeedsUpdate(kind)
				assert.True(t, needsUpdate, "Should need update for valid kind")

				_, err := collection.GetPage(kind, 0, 5, sdk.SortSpec{}, "")
				_ = err
			})
		}
	})

	t.Run("InvalidListKinds", func(t *testing.T) {
		collection := NewAbisCollection()
		invalidKinds := []string{"", "invalid", "Monitors", "Unknown"}

		for _, kind := range invalidKinds {
			t.Run(fmt.Sprintf("InvalidKind_%s", kind), func(t *testing.T) {
				assert.NotPanics(t, func() {
					collection.LoadData(types.ListKind(kind))
				}, "LoadData should not panic with invalid kind")

				assert.NotPanics(t, func() {
					collection.Reset(types.ListKind(kind))
				}, "Reset should not panic with invalid kind")

				needsUpdate := collection.NeedsUpdate(types.ListKind(kind))
				assert.False(t, needsUpdate, "Should not need update for invalid kind")

				_, err := collection.GetPage(types.ListKind(kind), 0, 5, sdk.SortSpec{}, "")
				assert.Error(t, err, "GetPage should return error for invalid kind")
				assert.Contains(t, err.Error(), "unexpected list kind", "Error should mention unexpected list kind")
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
					assert.Equal(t, AbisDownloaded, abisPage.Kind)
					assert.GreaterOrEqual(t, abisPage.TotalItems, 0)
				}
			})
		}
	})
}
