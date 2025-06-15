//go:build detailed

package names

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
	"github.com/stretchr/testify/assert"
)

func TestNamesCollectionLoadDataAsync(t *testing.T) {
	t.Run("LoadDataDoesNotBlock", func(t *testing.T) {
		collection := NewNamesCollection()

		start := time.Now()
		collection.LoadData(NamesAll)
		duration := time.Since(start)

		assert.Less(t, duration, 100*time.Millisecond, "LoadData should not block")
	})

	t.Run("LoadDataStartsAsyncOperation", func(t *testing.T) {
		collection := NewNamesCollection()

		collection.LoadData(NamesAll)

		needsUpdate := collection.NeedsUpdate(NamesAll)

		_ = needsUpdate

		time.Sleep(5 * time.Millisecond)

		collection.NeedsUpdate(NamesCustom)
	})

	t.Run("MultipleLoadDataCalls", func(t *testing.T) {
		collection := NewNamesCollection()

		for i := 0; i < 5; i++ {
			assert.NotPanics(t, func() {
				collection.LoadData(NamesAll)
			})
		}

		time.Sleep(10 * time.Millisecond)

		collection.NeedsUpdate(NamesAll)
	})

	t.Run("ConcurrentLoadDataCalls", func(t *testing.T) {
		collection := NewNamesCollection()
		var wg sync.WaitGroup
		numGoroutines := 3

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				listKind := []types.ListKind{NamesAll, NamesCustom, NamesPrefund}[index%3]
				collection.LoadData(listKind)
			}(i)
		}

		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Fatal("Concurrent LoadData calls took too long")
		}

		assert.NotPanics(t, func() {
			collection.NeedsUpdate(NamesAll)
		})
	})
}

func TestNamesCollectionAdvancedAsync(t *testing.T) {
	t.Run("LoadDataAllListKinds", func(t *testing.T) {
		collection := NewNamesCollection()

		listKinds := []types.ListKind{
			NamesAll, NamesCustom, NamesPrefund, NamesRegular, NamesBaddress,
		}

		for _, listKind := range listKinds {
			t.Run(fmt.Sprintf("LoadData_%s", listKind), func(t *testing.T) {
				start := time.Now()
				collection.LoadData(listKind)
				duration := time.Since(start)

				assert.Less(t, duration, 100*time.Millisecond, "LoadData should be non-blocking")
			})
		}

		time.Sleep(25 * time.Millisecond)

		for _, listKind := range listKinds {
			assert.NotPanics(t, func() {
				collection.NeedsUpdate(listKind)
			})
		}
	})

	t.Run("ResetDuringAsyncLoad", func(t *testing.T) {
		collection := NewNamesCollection()

		collection.LoadData(NamesAll)

		assert.NotPanics(t, func() {
			collection.Reset(NamesAll)
		})

		assert.NotPanics(t, func() {
			collection.LoadData(NamesAll)
		})
	})

	t.Run("GetPageDuringAsyncLoad", func(t *testing.T) {
		collection := NewNamesCollection()

		collection.LoadData(NamesAll)

		assert.NotPanics(t, func() {
			_, _ = collection.GetPage(NamesAll, 0, 5, sdk.SortSpec{}, "")
		})
	})
}

func TestNamesCollectionIntegration(t *testing.T) {
	t.Run("MultipleListKindsWorkflow", func(t *testing.T) {
		collection := NewNamesCollection()

		listKinds := []types.ListKind{
			NamesAll, NamesCustom, NamesPrefund, NamesRegular, NamesBaddress,
		}

		for _, listKind := range listKinds {
			collection.LoadData(listKind)
		}

		time.Sleep(15 * time.Millisecond)

		for _, listKind := range listKinds {
			t.Run(fmt.Sprintf("GetPage_%s", listKind), func(t *testing.T) {
				page, err := collection.GetPage(listKind, 0, 5, sdk.SortSpec{}, "")

				if err == nil && page != nil {
					assert.Equal(t, listKind, page.Kind)
					assert.GreaterOrEqual(t, page.TotalItems, 0)
					assert.LessOrEqual(t, len(page.Names), 5)
				}
			})
		}
	})

	t.Run("FilteringAcrossListKinds", func(t *testing.T) {
		collection := NewNamesCollection()

		listKinds := []types.ListKind{
			NamesAll, NamesCustom, NamesPrefund, NamesRegular, NamesBaddress,
		}

		filters := []string{"", "test", "0x", "custom"}

		for _, listKind := range listKinds {
			for _, filter := range filters {
				t.Run(fmt.Sprintf("Filter_%s_%s", listKind, filter), func(t *testing.T) {
					page, err := collection.GetPage(listKind, 0, 3, sdk.SortSpec{}, filter)

					if err == nil && page != nil {
						assert.Equal(t, listKind, page.Kind)
						assert.GreaterOrEqual(t, page.TotalItems, 0)
					}
				})
			}
		}
	})

	t.Run("SortingAndPagination", func(t *testing.T) {
		collection := NewNamesCollection()

		listKinds := []types.ListKind{NamesAll, NamesCustom}
		sortSpecs := []sdk.SortSpec{
			{},
		}
		pageSizes := []int{1, 5, 10}

		for _, listKind := range listKinds {
			for _, sortSpec := range sortSpecs {
				for _, pageSize := range pageSizes {
					t.Run(fmt.Sprintf("Sort_%s_%d", listKind, pageSize), func(t *testing.T) {
						page, err := collection.GetPage(listKind, 0, pageSize, sortSpec, "")

						if err == nil && page != nil {
							assert.Equal(t, listKind, page.Kind)
							assert.LessOrEqual(t, len(page.Names), pageSize)
						}
					})
				}
			}
		}
	})
}

func TestNamesCollectionPerformanceAndEdgeCases(t *testing.T) {
	t.Run("RapidStateQueries", func(t *testing.T) {
		collection := NewNamesCollection()

		const iterations = 50
		start := time.Now()

		for i := 0; i < iterations; i++ {
			listKind := []types.ListKind{NamesAll, NamesCustom, NamesPrefund, NamesRegular, NamesBaddress}[i%5]
			collection.NeedsUpdate(listKind)
			if i%10 == 0 {
				_, _ = collection.GetPage(NamesAll, 0, 1, sdk.SortSpec{}, "")
			}
		}

		duration := time.Since(start)
		avgPerCall := duration / iterations

		assert.Less(t, avgPerCall, 10*time.Millisecond)
		t.Logf("Average time per NeedsUpdate call: %v", avgPerCall)
	})

	t.Run("LargePageSizeHandling", func(t *testing.T) {
		collection := NewNamesCollection()

		largeSizes := []int{0, 1, 10, 25}

		for _, size := range largeSizes {
			t.Run(fmt.Sprintf("PageSize_%d", size), func(t *testing.T) {
				start := time.Now()
				page, err := collection.GetPage(NamesAll, 0, size, sdk.SortSpec{}, "")
				duration := time.Since(start)

				if err == nil && page != nil {
					assert.Equal(t, NamesAll, page.Kind)
					assert.GreaterOrEqual(t, page.TotalItems, 0)
					assert.LessOrEqual(t, len(page.Names), size, "Returned items should not exceed page size")

					maxExpectedTime := time.Duration(size/5+1) * 200 * time.Millisecond
					assert.Less(t, duration, maxExpectedTime, "Page retrieval should complete in reasonable time")
				}
			})
		}
	})

	t.Run("ConcurrentGetPageCalls", func(t *testing.T) {
		collection := NewNamesCollection()
		var wg sync.WaitGroup
		numGoroutines := 3

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				listKind := []types.ListKind{NamesAll, NamesCustom, NamesPrefund}[index%3]
				_, _ = collection.GetPage(listKind, 0, 5, sdk.SortSpec{}, "")
			}(i)
		}

		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Fatal("Concurrent GetPage calls took too long")
		}
	})
}

func TestNamesCollectionBoundaryConditions(t *testing.T) {
	t.Run("InvalidListKindOperations", func(t *testing.T) {
		collection := NewNamesCollection()

		invalidKinds := []types.ListKind{"", "invalid", "NotAListKind", "names"}

		for _, invalidKind := range invalidKinds {
			t.Run(fmt.Sprintf("Invalid_%s", invalidKind), func(t *testing.T) {
				assert.NotPanics(t, func() {
					collection.LoadData(invalidKind)
				})

				assert.NotPanics(t, func() {
					collection.Reset(invalidKind)
				})

				needsUpdate := collection.NeedsUpdate(invalidKind)
				assert.False(t, needsUpdate)

				_, err := collection.GetPage(invalidKind, 0, 10, sdk.SortSpec{}, "")
				assert.Error(t, err)
			})
		}
	})

	t.Run("ExtremePageParameters", func(t *testing.T) {
		collection := NewNamesCollection()

		extremeParams := []struct {
			first    int
			pageSize int
			desc     string
		}{
			{-1, 10, "negative_first"},
			{0, -1, "negative_pageSize"},
			{1000000, 1, "very_large_first"},
			{0, 1000000, "very_large_pageSize"},
			{-100, -100, "both_negative"},
		}

		for _, param := range extremeParams {
			t.Run(param.desc, func(t *testing.T) {
				assert.NotPanics(t, func() {
					_, _ = collection.GetPage(NamesAll, param.first, param.pageSize, sdk.SortSpec{}, "")
				})
			})
		}
	})

	t.Run("SpecialFilterStrings", func(t *testing.T) {
		collection := NewNamesCollection()

		specialFilters := []string{
			"",
			" ",
			"0x0000000000000000000000000000000000000000",
			"0xffffffffffffffffffffffffffffffffffffffffff",
			"@#$%^&*()",
			"verylongfilterstringshouldnotcauseproblems",
			"🚀🌟💫",
		}

		for _, filter := range specialFilters {
			t.Run(fmt.Sprintf("Filter_%s", filter), func(t *testing.T) {
				assert.NotPanics(t, func() {
					_, _ = collection.GetPage(NamesAll, 0, 5, sdk.SortSpec{}, filter)
				})
			})
		}
	})

	t.Run("ResetAfterMultipleLoads", func(t *testing.T) {
		collection := NewNamesCollection()

		listKinds := []types.ListKind{NamesAll, NamesCustom, NamesPrefund}
		for _, listKind := range listKinds {
			collection.LoadData(listKind)
		}

		time.Sleep(10 * time.Millisecond)

		collection.Reset(NamesAll)

		for _, listKind := range listKinds {
			assert.NotPanics(t, func() {
				collection.NeedsUpdate(listKind)
			})
		}
	})

	t.Run("StressTestMixedOperations", func(t *testing.T) {
		collection := NewNamesCollection()
		var wg sync.WaitGroup
		numGoroutines := 3

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				listKind := []types.ListKind{NamesAll, NamesCustom, NamesPrefund}[index%3]

				for j := 0; j < 3; j++ {
					switch j % 4 {
					case 0:
						collection.LoadData(listKind)
					case 1:
						collection.NeedsUpdate(listKind)
					case 2:
						_, _ = collection.GetPage(listKind, 0, 2, sdk.SortSpec{}, "")
					case 3:
						collection.Reset(listKind)
					}
				}
			}(i)
		}

		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
		case <-time.After(3 * time.Second):
			t.Fatal("Stress test took too long")
		}

		assert.NotPanics(t, func() {
			collection.NeedsUpdate(NamesAll)
		})
	})
}
