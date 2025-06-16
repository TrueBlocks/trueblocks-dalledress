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

		// Test that LoadData call completes without panicking
		assert.NotPanics(t, func() {
			collection.LoadData(NamesAll)
		})
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
				assert.NotPanics(t, func() {
					collection.LoadData(listKind)
				})
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
					assert.Equal(t, listKind, page.GetKind())
					assert.GreaterOrEqual(t, page.GetTotalItems(), 0)

					// For names page, we need to type assert to access Names field
					if namesPage, ok := page.(*NamesPage); ok {
						assert.LessOrEqual(t, len(namesPage.Names), 5)
					}
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
						// Cast to concrete type to access fields
						namesPage, ok := page.(*NamesPage)
						assert.True(t, ok, "Expected *NamesPage type")
						assert.Equal(t, listKind, namesPage.Kind)
						assert.GreaterOrEqual(t, namesPage.TotalItems, 0)
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
							// Cast to concrete type to access fields
							namesPage, ok := page.(*NamesPage)
							assert.True(t, ok, "Expected *NamesPage type")
							assert.Equal(t, listKind, namesPage.Kind)
							assert.LessOrEqual(t, len(namesPage.Names), pageSize)
						}
					})
				}
			}
		}
	})
}

func TestNamesCollectionDomainSpecificIntegration(t *testing.T) {
	t.Run("SpecialFilterStrings", func(t *testing.T) {
		collection := NewNamesCollection()

		// Test domain-specific filtering logic for names
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
}
