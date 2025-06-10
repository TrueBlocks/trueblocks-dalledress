package facets

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/store"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

// Test data structures
type TestItem struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Value int    `json:"value"`
}

// Implement the Modeler interface with the correct signature
func (t *TestItem) Model(chain, format string, verbose bool, extraOptions map[string]any) coreTypes.Model {
	// Return a proper Model struct with Data and Order fields
	return coreTypes.Model{
		Data: map[string]any{
			"id":    t.ID,
			"name":  t.Name,
			"value": t.Value,
		},
		Order: []string{"id", "name", "value"},
	}
}

type TestListKind types.ListKind

const (
	TestList TestListKind = "test-list"
)

// Test utilities
func setupTestEnvironment() *msgs.TestHelpers {
	return msgs.NewTestHelpers()
}

func createTestStore() *store.Store[TestItem] {
	return store.NewStore(
		"test-store",
		func(ctx *output.RenderCtx) error {
			// Simulate data fetching
			testData := []TestItem{
				{ID: 1, Name: "Item1", Value: 10},
				{ID: 2, Name: "Item2", Value: 20},
				{ID: 3, Name: "Item3", Value: 30},
				{ID: 4, Name: "Item4", Value: 40},
				{ID: 5, Name: "Item5", Value: 50},
			}

			for _, item := range testData {
				// Send items one by one, waiting for each to be received
				// This matches the typical streaming behavior
				ctx.ModelChan <- &item
			}

			// Close both channels to signal completion
			close(ctx.ModelChan)
			close(ctx.ErrorChan)
			return nil
		},
		func(itemIntf interface{}) *TestItem {
			if item, ok := itemIntf.(TestItem); ok {
				return &item
			}
			if item, ok := itemIntf.(*TestItem); ok {
				return item
			}
			return nil
		},
	)
}

func createTestFacet(store *store.Store[TestItem]) *Facet[TestItem] {
	return NewFacet(
		types.ListKind(TestList),
		nil, // No filter
		nil, // No duplicate check
		store,
	)
}

func createFilteredFacet(store *store.Store[TestItem], minValue int) *Facet[TestItem] {
	return NewFacet(
		types.ListKind(TestList),
		func(item *TestItem) bool {
			return item.Value >= minValue
		},
		nil, // No duplicate check
		store,
	)
}

func createDedupedFacet(store *store.Store[TestItem]) *Facet[TestItem] {
	return NewFacet(
		types.ListKind(TestList),
		nil, // No filter
		func(existing []*TestItem, newItem *TestItem) bool {
			for _, existingItem := range existing {
				if existingItem.ID == newItem.ID {
					return true // Is duplicate
				}
			}
			return false // Not a duplicate
		},
		store,
	)
}

// Tests
func TestNewFacet(t *testing.T) {
	helpers := setupTestEnvironment()
	defer helpers.Cleanup()

	testStore := createTestStore()
	facet := createTestFacet(testStore)

	// Test initial state
	if facet.GetState() != StateStale {
		t.Errorf("Expected initial state to be StateStale, got %v", facet.GetState())
	}

	if facet.Count() != 0 {
		t.Errorf("Expected initial count to be 0, got %d", facet.Count())
	}

	if !facet.NeedsUpdate() {
		t.Error("New facet should need update")
	}

	if facet.IsLoaded() {
		t.Error("New facet should not be loaded")
	}

	if facet.IsFetching() {
		t.Error("New facet should not be fetching initially")
	}

	// Test that store is properly set
	if facet.GetStore() != testStore {
		t.Error("Facet should reference the provided store")
	}
}

func TestFacetStateTransitions(t *testing.T) {
	helpers := setupTestEnvironment()
	defer helpers.Cleanup()

	testStore := createTestStore()
	facet := createTestFacet(testStore)

	// Test StartFetching
	if !facet.StartFetching() {
		t.Error("StartFetching should return true on first call")
	}

	if facet.GetState() != StateFetching {
		t.Errorf("Expected state to be StateFetching, got %v", facet.GetState())
	}

	if facet.StartFetching() {
		t.Error("StartFetching should return false when already fetching")
	}

	// Test SetPartial
	facet.SetPartial()
	if facet.GetState() != StatePartial {
		t.Errorf("Expected state to be StatePartial, got %v", facet.GetState())
	}

	// Test Reset
	facet.Reset()
	if facet.GetState() != StateStale {
		t.Errorf("Expected state to be StateStale after reset, got %v", facet.GetState())
	}

	if facet.Count() != 0 {
		t.Errorf("Expected count to be 0 after reset, got %d", facet.Count())
	}
}

func TestFacetLoad(t *testing.T) {
	helpers := setupTestEnvironment()
	defer helpers.Cleanup()

	testStore := createTestStore()
	facet := createTestFacet(testStore)

	// Test initial load
	result, err := facet.Load()
	if err != nil {
		t.Fatalf("Load should not return error on first call: %v", err)
	}

	if result == nil {
		t.Fatal("Load should return a result")
	}

	if result.Payload.ListKind != types.ListKind(TestList) {
		t.Errorf("Expected ListKind to be %v, got %v", TestList, result.Payload.ListKind)
	}

	// Wait for loading to complete
	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatal("Timeout waiting for facet to load")
		case <-ticker.C:
			if facet.IsLoaded() {
				goto loaded
			}
		}
	}

loaded:
	// Test that data was loaded
	if facet.Count() != 5 {
		t.Errorf("Expected 5 items after load, got %d", facet.Count())
	}

	// Test second load returns cached result
	result2, err := facet.Load()
	if err != nil {
		t.Fatalf("Second load should not return error: %v", err)
	}

	if result2.Payload.CurrentCount != 5 {
		t.Errorf("Expected cached result to have 5 items, got %d", result2.Payload.CurrentCount)
	}

	// Test load while already loading
	facet.Reset()
	facet.StartFetching()
	_, err = facet.Load()
	if err != ErrAlreadyLoading {
		t.Errorf("Expected ErrAlreadyLoading, got %v", err)
	}
}

func TestFacetFiltering(t *testing.T) {
	helpers := setupTestEnvironment()
	defer helpers.Cleanup()

	testStore := createTestStore()
	facet := createFilteredFacet(testStore, 30) // Filter items with value >= 30

	// Load data
	_, err := facet.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Wait for loading to complete
	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatal("Timeout waiting for facet to load")
		case <-ticker.C:
			if facet.IsLoaded() {
				goto loaded
			}
		}
	}

loaded:
	// Should only have items with value >= 30 (items 3, 4, 5)
	if facet.Count() != 3 {
		t.Errorf("Expected 3 filtered items, got %d", facet.Count())
	}
}

func TestFacetDeduplication(t *testing.T) {
	helpers := setupTestEnvironment()
	defer helpers.Cleanup()

	testStore := createTestStore()
	facet := createDedupedFacet(testStore)

	// Manually add duplicate items to test deduplication
	testItem := &TestItem{ID: 1, Name: "Duplicate", Value: 100}

	// Add item twice
	facet.OnNewItem(testItem, 0)
	facet.OnNewItem(testItem, 1)

	// Should only have one item due to deduplication
	if facet.Count() != 1 {
		t.Errorf("Expected 1 item after deduplication, got %d", facet.Count())
	}
}

func TestFacetPagination(t *testing.T) {
	helpers := setupTestEnvironment()
	defer helpers.Cleanup()

	testStore := createTestStore()
	facet := createTestFacet(testStore)

	// Load data
	_, err := facet.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Wait for loading to complete
	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatal("Timeout waiting for facet to load")
		case <-ticker.C:
			if facet.IsLoaded() {
				goto loaded
			}
		}
	}

loaded:
	// Test first page
	page1, err := facet.GetPage(0, 2, nil, sdk.SortSpec{}, nil)
	if err != nil {
		t.Fatalf("GetPage failed: %v", err)
	}

	if len(page1.Items) != 2 {
		t.Errorf("Expected 2 items in first page, got %d", len(page1.Items))
	}

	if page1.TotalItems != 5 {
		t.Errorf("Expected total items to be 5, got %d", page1.TotalItems)
	}

	if page1.State != StateLoaded {
		t.Errorf("Expected page state to be StateLoaded, got %v", page1.State)
	}

	// Test second page
	page2, err := facet.GetPage(2, 2, nil, sdk.SortSpec{}, nil)
	if err != nil {
		t.Fatalf("GetPage failed: %v", err)
	}

	if len(page2.Items) != 2 {
		t.Errorf("Expected 2 items in second page, got %d", len(page2.Items))
	}

	// Test last page
	page3, err := facet.GetPage(4, 2, nil, sdk.SortSpec{}, nil)
	if err != nil {
		t.Fatalf("GetPage failed: %v", err)
	}

	if len(page3.Items) != 1 {
		t.Errorf("Expected 1 item in last page, got %d", len(page3.Items))
	}

	// Test out of bounds
	pageEmpty, err := facet.GetPage(10, 2, nil, sdk.SortSpec{}, nil)
	if err != nil {
		t.Fatalf("GetPage failed: %v", err)
	}

	if len(pageEmpty.Items) != 0 {
		t.Errorf("Expected 0 items for out of bounds page, got %d", len(pageEmpty.Items))
	}
}

func TestFacetPageWithFilter(t *testing.T) {
	helpers := setupTestEnvironment()
	defer helpers.Cleanup()

	testStore := createTestStore()
	facet := createTestFacet(testStore)

	// Load data
	_, err := facet.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Wait for loading to complete
	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatal("Timeout waiting for facet to load")
		case <-ticker.C:
			if facet.IsLoaded() {
				goto loaded
			}
		}
	}

loaded:
	// Test pagination with runtime filter
	filterFunc := func(item *TestItem) bool {
		return item.Value >= 30
	}

	page, err := facet.GetPage(0, 5, filterFunc, sdk.SortSpec{}, nil)
	if err != nil {
		t.Fatalf("GetPage with filter failed: %v", err)
	}

	// Should have 3 items (values 30, 40, 50)
	if len(page.Items) != 3 {
		t.Errorf("Expected 3 filtered items, got %d", len(page.Items))
	}

	if page.TotalItems != 3 {
		t.Errorf("Expected total filtered items to be 3, got %d", page.TotalItems)
	}
}

func TestFacetSorting(t *testing.T) {
	helpers := setupTestEnvironment()
	defer helpers.Cleanup()

	testStore := createTestStore()
	facet := createTestFacet(testStore)

	// Load data
	_, err := facet.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Wait for loading to complete
	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatal("Timeout waiting for facet to load")
		case <-ticker.C:
			if facet.IsLoaded() {
				goto loaded
			}
		}
	}

loaded:
	// Custom sort function (reverse by value)
	sortFunc := func(items []TestItem, spec sdk.SortSpec) error {
		n := len(items)
		for i := 0; i < n-1; i++ {
			for j := 0; j < n-i-1; j++ {
				if items[j].Value < items[j+1].Value {
					items[j], items[j+1] = items[j+1], items[j]
				}
			}
		}
		return nil
	}

	page, err := facet.GetPage(0, 5, nil, sdk.SortSpec{}, sortFunc)
	if err != nil {
		t.Fatalf("GetPage with sort failed: %v", err)
	}

	// Should be sorted in reverse order by value
	if len(page.Items) != 5 {
		t.Errorf("Expected 5 items, got %d", len(page.Items))
	}

	if page.Items[0].Value != 50 {
		t.Errorf("Expected first item value to be 50, got %d", page.Items[0].Value)
	}

	if page.Items[4].Value != 10 {
		t.Errorf("Expected last item value to be 10, got %d", page.Items[4].Value)
	}
}

func TestFacetSortingError(t *testing.T) {
	helpers := setupTestEnvironment()
	defer helpers.Cleanup()

	testStore := createTestStore()
	facet := createTestFacet(testStore)

	// Load data
	_, err := facet.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Wait for loading to complete
	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatal("Timeout waiting for facet to load")
		case <-ticker.C:
			if facet.IsLoaded() {
				goto loaded
			}
		}
	}

loaded:
	// Sort function that returns an error
	sortFunc := func(items []TestItem, spec sdk.SortSpec) error {
		return errors.New("sort error")
	}

	_, err = facet.GetPage(0, 5, nil, sdk.SortSpec{}, sortFunc)
	if err == nil {
		t.Error("Expected GetPage to return error when sort function fails")
	}

	if err.Error() != "error sorting data: sort error" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestFacetObserverInterface(t *testing.T) {
	helpers := setupTestEnvironment()
	defer helpers.Cleanup()

	testStore := createTestStore()
	facet := createTestFacet(testStore)

	// Test OnStateChanged
	facet.OnStateChanged(store.StateFetching, "Test fetching")
	if facet.GetState() != StateFetching {
		t.Errorf("Expected state to be StateFetching, got %v", facet.GetState())
	}

	facet.OnStateChanged(store.StateLoaded, "Test loaded")
	if facet.GetState() != StateLoaded {
		t.Errorf("Expected state to be StateLoaded, got %v", facet.GetState())
	}

	facet.OnStateChanged(store.StateError, "Test error")
	if facet.GetState() != StateError {
		t.Errorf("Expected state to be StateError, got %v", facet.GetState())
	}

	// Test partial state with existing data
	facet.OnNewItem(&TestItem{ID: 1, Name: "Test", Value: 100}, 0)
	facet.OnStateChanged(store.StateError, "Test partial error")
	if facet.GetState() != StatePartial {
		t.Errorf("Expected state to be StatePartial when error occurs with existing data, got %v", facet.GetState())
	}

	facet.OnStateChanged(store.StateCanceled, "Test canceled")
	if facet.GetState() != StateStale {
		t.Errorf("Expected state to be StateStale after cancel, got %v", facet.GetState())
	}
}

func TestFacetOnNewItem(t *testing.T) {
	helpers := setupTestEnvironment()
	defer helpers.Cleanup()

	testStore := createTestStore()
	facet := createTestFacet(testStore)

	// Test adding items
	item1 := &TestItem{ID: 1, Name: "Test1", Value: 100}
	item2 := &TestItem{ID: 2, Name: "Test2", Value: 200}

	facet.OnNewItem(item1, 0)
	if facet.Count() != 1 {
		t.Errorf("Expected count to be 1, got %d", facet.Count())
	}

	facet.OnNewItem(item2, 1)
	if facet.Count() != 2 {
		t.Errorf("Expected count to be 2, got %d", facet.Count())
	}
}

func TestFacetForEvery(t *testing.T) {
	helpers := setupTestEnvironment()
	defer helpers.Cleanup()

	testStore := createTestStore()
	facet := createTestFacet(testStore)

	// Add test data
	items := []*TestItem{
		{ID: 1, Name: "Test1", Value: 100},
		{ID: 2, Name: "Test2", Value: 200},
		{ID: 3, Name: "Test3", Value: 300},
	}

	for i, item := range items {
		facet.OnNewItem(item, i)
	}

	// Test ForEvery with match function
	matchCount, err := facet.ForEvery(
		func(item *TestItem) (error, bool) {
			return nil, true // Action function
		},
		func(item *TestItem) bool {
			return item.Value >= 200 // Match items with value >= 200
		},
	)

	if err != nil {
		t.Fatalf("ForEvery failed: %v", err)
	}

	if matchCount != 2 {
		t.Errorf("Expected 2 matches, got %d", matchCount)
	}

	// Should have 1 item left (the one with Value < 200)
	if facet.Count() != 1 {
		t.Errorf("Expected 1 item remaining, got %d", facet.Count())
	}
}

func TestFacetSyncWithStore(t *testing.T) {
	helpers := setupTestEnvironment()
	defer helpers.Cleanup()

	testStore := createTestStore()
	facet := createTestFacet(testStore)

	// Manually add items to store
	testStore.AddItem(TestItem{ID: 1, Name: "Direct1", Value: 111}, 0)
	testStore.AddItem(TestItem{ID: 2, Name: "Direct2", Value: 222}, 1)

	// Sync facet with store
	facet.SyncWithStore()

	if facet.Count() != 2 {
		t.Errorf("Expected 2 items after sync, got %d", facet.Count())
	}
}

func TestFacetProgressIntegration(t *testing.T) {
	helpers := setupTestEnvironment()
	defer helpers.Cleanup()

	testStore := createTestStore()
	facet := createTestFacet(testStore)

	// Test that progress is initialized
	if facet.progress == nil {
		t.Error("Progress should be initialized in facet")
	}

	// Test progress updates during OnNewItem
	item := &TestItem{ID: 1, Name: "Progress Test", Value: 100}
	facet.OnNewItem(item, 0)

	// Progress should be called internally, but we can't easily test the specific calls
	// without more complex mocking. The important thing is that it doesn't panic.
}

func TestFacetConcurrency(t *testing.T) {
	helpers := setupTestEnvironment()
	defer helpers.Cleanup()

	testStore := createTestStore()
	facet := createTestFacet(testStore)

	// Test concurrent access to facet
	var wg sync.WaitGroup
	numGoroutines := 10
	itemsPerGoroutine := 5

	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(routineID int) {
			defer wg.Done()
			for j := 0; j < itemsPerGoroutine; j++ {
				item := &TestItem{
					ID:    routineID*itemsPerGoroutine + j,
					Name:  fmt.Sprintf("Concurrent-%d-%d", routineID, j),
					Value: (routineID*itemsPerGoroutine + j) * 10,
				}
				facet.OnNewItem(item, routineID*itemsPerGoroutine+j)
			}
		}(i)
	}

	wg.Wait()

	expectedCount := numGoroutines * itemsPerGoroutine
	if facet.Count() != expectedCount {
		t.Errorf("Expected %d items after concurrent access, got %d", expectedCount, facet.Count())
	}
}

func TestFacetEdgeCases(t *testing.T) {
	helpers := setupTestEnvironment()
	defer helpers.Cleanup()

	testStore := createTestStore()
	facet := createTestFacet(testStore)

	// Test GetPage with empty data
	page, err := facet.GetPage(0, 10, nil, sdk.SortSpec{}, nil)
	if err != nil {
		t.Fatalf("GetPage should not fail with empty data: %v", err)
	}

	if len(page.Items) != 0 {
		t.Errorf("Expected 0 items with empty data, got %d", len(page.Items))
	}

	if page.TotalItems != 0 {
		t.Errorf("Expected total items to be 0, got %d", page.TotalItems)
	}

	// Test negative pagination values
	_, err = facet.GetPage(-1, 5, nil, sdk.SortSpec{}, nil)
	if err != nil {
		t.Fatalf("GetPage should handle negative start: %v", err)
	}

	// Test zero page size
	page, err = facet.GetPage(0, 0, nil, sdk.SortSpec{}, nil)
	if err != nil {
		t.Fatalf("GetPage should handle zero page size: %v", err)
	}

	if len(page.Items) != 0 {
		t.Errorf("Expected 0 items with zero page size, got %d", len(page.Items))
	}
}

// Benchmark tests
func BenchmarkFacetOnNewItem(b *testing.B) {
	helpers := setupTestEnvironment()
	defer helpers.Cleanup()

	testStore := createTestStore()
	facet := createTestFacet(testStore)

	item := &TestItem{ID: 1, Name: "Benchmark", Value: 100}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		facet.OnNewItem(item, i)
	}
}

func BenchmarkFacetGetPage(b *testing.B) {
	helpers := setupTestEnvironment()
	defer helpers.Cleanup()

	testStore := createTestStore()
	facet := createTestFacet(testStore)

	// Add test data
	for i := 0; i < 1000; i++ {
		item := &TestItem{ID: i, Name: fmt.Sprintf("Item%d", i), Value: i * 10}
		facet.OnNewItem(item, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = facet.GetPage(0, 20, nil, sdk.SortSpec{}, nil)
	}
}

func TestStoreDirectly(t *testing.T) {
	testStore := store.NewStore(
		"direct-test-store",
		func(ctx *output.RenderCtx) error {
			fmt.Println("Direct test: Query function starting...")

			// Send one item
			testItem := &TestItem{ID: 1, Name: "Direct", Value: 100}

			fmt.Println("Direct test: Sending item...")
			ctx.ModelChan <- testItem

			fmt.Println("Direct test: Closing ModelChan...")
			close(ctx.ModelChan)

			fmt.Println("Direct test: Closing ErrorChan...")
			close(ctx.ErrorChan)

			fmt.Println("Direct test: Query function completed")
			return nil
		},
		func(itemIntf interface{}) *TestItem {
			fmt.Println("Direct test: Processing item...")
			if item, ok := itemIntf.(*TestItem); ok {
				fmt.Println("Direct test: Item processed successfully")
				return item
			}
			fmt.Println("Direct test: Failed to process item")
			return nil
		},
	)

	fmt.Println("Direct test: Initial store state:", testStore.GetState())

	// Fetch data
	fmt.Println("Direct test: Starting fetch...")
	err := testStore.Fetch()
	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}

	fmt.Println("Direct test: Fetch completed")
	fmt.Println("Direct test: Final store state:", testStore.GetState())
	fmt.Println("Direct test: Items count:", len(testStore.GetItems()))

	// Check if state is loaded
	if testStore.GetState() != store.StateLoaded {
		t.Errorf("Expected store state to be StateLoaded (%d), got %d", store.StateLoaded, testStore.GetState())
	}

	if len(testStore.GetItems()) != 1 {
		t.Errorf("Expected 1 item, got %d", len(testStore.GetItems()))
	}
}

func TestSimpleStoreLoad(t *testing.T) {
	helpers := setupTestEnvironment()
	defer helpers.Cleanup()

	testStore := store.NewStore(
		"debug-store",
		func(ctx *output.RenderCtx) error {
			fmt.Println("Query function starting...")

			// Simple test item
			testItem := &TestItem{ID: 1, Name: "Debug", Value: 100}

			fmt.Println("Sending item to ModelChan...")
			select {
			case ctx.ModelChan <- testItem:
				fmt.Println("Item sent successfully")
			default:
				fmt.Println("Failed to send item")
				return fmt.Errorf("failed to send item")
			}

			fmt.Println("Closing ModelChan...")
			close(ctx.ModelChan)

			fmt.Println("Closing ErrorChan...")
			close(ctx.ErrorChan)

			fmt.Println("Query function completed successfully")
			return nil
		},
		func(itemIntf interface{}) *TestItem {
			fmt.Println("Process function called...")
			if item, ok := itemIntf.(TestItem); ok {
				fmt.Println("Processed TestItem (value)")
				return &item
			}
			if item, ok := itemIntf.(*TestItem); ok {
				fmt.Println("Processed TestItem (pointer)")
				return item
			}
			fmt.Println("Process function failed to convert item")
			return nil
		},
	)

	facet := NewFacet(
		types.ListKind(TestList),
		nil, // No filter
		nil, // No duplicate check
		testStore,
	)

	fmt.Println("Initial state:", facet.GetState())

	// Load the facet
	result, err := facet.Load()
	if err != nil {
		t.Fatalf("Load should not return error: %v", err)
	}

	if result == nil {
		t.Fatal("Load should return a result")
	}

	fmt.Println("Load called, state:", facet.GetState())

	// Wait for loading to complete
	timeout := time.After(10 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			fmt.Println("Final state:", facet.GetState())
			fmt.Println("Final count:", facet.Count())
			fmt.Println("Store state:", testStore.GetState())
			t.Fatal("Timeout waiting for facet to load")
		case <-ticker.C:
			currentState := facet.GetState()
			fmt.Println("Current state:", currentState, "Count:", facet.Count())
			if facet.IsLoaded() {
				fmt.Println("Facet is loaded!")
				goto loaded
			}
		}
	}

loaded:
	if facet.Count() != 1 {
		t.Errorf("Expected 1 item after load, got %d", facet.Count())
	}
}
