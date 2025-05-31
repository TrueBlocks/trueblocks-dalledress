package streaming

import (
	"errors"
	"reflect"
	"testing"

	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

// Test data structures
type TestItem struct {
	ID   int
	Name string
	Age  int
}

func TestFilterPageSlice(t *testing.T) {
	items := []TestItem{
		{ID: 1, Name: "Alice", Age: 30},
		{ID: 2, Name: "Bob", Age: 25},
		{ID: 3, Name: "Charlie", Age: 35},
		{ID: 4, Name: "David", Age: 20},
	}

	tests := []struct {
		name     string
		items    *[]TestItem
		filterFn func(TestItem) bool
		expected []TestItem
	}{
		{
			name:  "filter by age > 25",
			items: &items,
			filterFn: func(item TestItem) bool {
				return item.Age > 25
			},
			expected: []TestItem{
				{ID: 1, Name: "Alice", Age: 30},
				{ID: 3, Name: "Charlie", Age: 35},
			},
		},
		{
			name:  "filter by name starts with 'B'",
			items: &items,
			filterFn: func(item TestItem) bool {
				return item.Name[0] == 'B'
			},
			expected: []TestItem{
				{ID: 2, Name: "Bob", Age: 25},
			},
		},
		{
			name:  "no filter (all pass)",
			items: &items,
			filterFn: func(item TestItem) bool {
				return true
			},
			expected: items,
		},
		{
			name:  "filter none",
			items: &items,
			filterFn: func(item TestItem) bool {
				return false
			},
			expected: []TestItem{},
		},
		{
			name:  "empty slice",
			items: &[]TestItem{},
			filterFn: func(item TestItem) bool {
				return true
			},
			expected: []TestItem{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterPageSlice(tt.items, tt.filterFn)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("FilterPageSlice() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPaginateSlice(t *testing.T) {
	items := []TestItem{
		{ID: 1, Name: "Alice", Age: 30},
		{ID: 2, Name: "Bob", Age: 25},
		{ID: 3, Name: "Charlie", Age: 35},
		{ID: 4, Name: "David", Age: 20},
		{ID: 5, Name: "Eve", Age: 28},
	}

	tests := []struct {
		name     string
		items    []TestItem
		first    int
		pageSize int
		expected []TestItem
	}{
		{
			name:     "first page",
			items:    items,
			first:    0,
			pageSize: 2,
			expected: []TestItem{
				{ID: 1, Name: "Alice", Age: 30},
				{ID: 2, Name: "Bob", Age: 25},
			},
		},
		{
			name:     "second page",
			items:    items,
			first:    2,
			pageSize: 2,
			expected: []TestItem{
				{ID: 3, Name: "Charlie", Age: 35},
				{ID: 4, Name: "David", Age: 20},
			},
		},
		{
			name:     "last page (partial)",
			items:    items,
			first:    4,
			pageSize: 2,
			expected: []TestItem{
				{ID: 5, Name: "Eve", Age: 28},
			},
		},
		{
			name:     "start beyond items",
			items:    items,
			first:    10,
			pageSize: 2,
			expected: []TestItem{},
		},
		{
			name:     "page size larger than remaining items",
			items:    items,
			first:    3,
			pageSize: 10,
			expected: []TestItem{
				{ID: 4, Name: "David", Age: 20},
				{ID: 5, Name: "Eve", Age: 28},
			},
		},
		{
			name:     "empty slice",
			items:    []TestItem{},
			first:    0,
			pageSize: 2,
			expected: []TestItem{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PaginateSlice(tt.items, tt.first, tt.pageSize)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("PaginateSlice() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSortPageSlice(t *testing.T) {
	items := []TestItem{
		{ID: 3, Name: "Charlie", Age: 35},
		{ID: 1, Name: "Alice", Age: 30},
		{ID: 2, Name: "Bob", Age: 25},
	}

	t.Run("successful sort", func(t *testing.T) {
		sortFn := func(items []TestItem, sortSpec sdk.SortSpec) error {
			// Simple sort by ID for testing
			for i := 0; i < len(items)-1; i++ {
				for j := i + 1; j < len(items); j++ {
					if items[i].ID > items[j].ID {
						items[i], items[j] = items[j], items[i]
					}
				}
			}
			return nil
		}

		sortSpec := sdk.SortSpec{} // Mock sort spec
		err := SortPageSlice(items, sortSpec, sortFn)
		if err != nil {
			t.Errorf("SortPageSlice() error = %v, want nil", err)
		}

		// Check if items are sorted by ID
		expected := []TestItem{
			{ID: 1, Name: "Alice", Age: 30},
			{ID: 2, Name: "Bob", Age: 25},
			{ID: 3, Name: "Charlie", Age: 35},
		}

		if !reflect.DeepEqual(items, expected) {
			t.Errorf("SortPageSlice() sorted items = %v, want %v", items, expected)
		}
	})

	t.Run("sort function error", func(t *testing.T) {
		sortFn := func(items []TestItem, sortSpec sdk.SortSpec) error {
			return errors.New("sort error")
		}

		sortSpec := sdk.SortSpec{}
		err := SortPageSlice(items, sortSpec, sortFn)
		if err == nil {
			t.Error("SortPageSlice() error = nil, want error")
		}
	})
}

func TestCreatePageFilter(t *testing.T) {
	searchFields := func(item TestItem) []string {
		return []string{item.Name}
	}

	tests := []struct {
		name     string
		filter   string
		item     TestItem
		expected bool
	}{
		{
			name:     "empty filter returns true",
			filter:   "",
			item:     TestItem{ID: 1, Name: "Alice", Age: 30},
			expected: true,
		},
		{
			name:     "filter matches name (case insensitive)",
			filter:   "alice",
			item:     TestItem{ID: 1, Name: "Alice", Age: 30},
			expected: true,
		},
		{
			name:     "filter matches partial name",
			filter:   "lic",
			item:     TestItem{ID: 1, Name: "Alice", Age: 30},
			expected: true,
		},
		{
			name:     "filter does not match",
			filter:   "xyz",
			item:     TestItem{ID: 1, Name: "Alice", Age: 30},
			expected: false,
		},
		{
			name:     "filter matches different case",
			filter:   "alice",
			item:     TestItem{ID: 1, Name: "Alice", Age: 30},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filterFn := CreatePageFilter(tt.filter, searchFields)
			result := filterFn(tt.item)
			if result != tt.expected {
				t.Errorf("CreatePageFilter()(%v) = %v, want %v", tt.item, result, tt.expected)
			}
		})
	}
}

func TestProcessPage(t *testing.T) {
	sourceItems := []TestItem{
		{ID: 1, Name: "Alice", Age: 30},
		{ID: 2, Name: "Bob", Age: 25},
		{ID: 3, Name: "Charlie", Age: 35},
		{ID: 4, Name: "David", Age: 20},
		{ID: 5, Name: "Eve", Age: 28},
	}

	t.Run("successful processing", func(t *testing.T) {
		filterFn := func(item TestItem) bool {
			return item.Age >= 25 // Filter out David (age 20)
		}

		sortFn := func(items []TestItem, sortSpec sdk.SortSpec) error {
			// Sort by ID
			for i := 0; i < len(items)-1; i++ {
				for j := i + 1; j < len(items); j++ {
					if items[i].ID > items[j].ID {
						items[i], items[j] = items[j], items[i]
					}
				}
			}
			return nil
		}

		sortSpec := sdk.SortSpec{}
		paginatedItems, totalFiltered, totalSource, err := ProcessPage(
			"TestItem",
			&sourceItems,
			sortSpec,
			sortFn,
			filterFn,
			0, // first
			2, // pageSize
		)

		if err != nil {
			t.Errorf("ProcessPage() error = %v, want nil", err)
		}

		expectedPaginated := []TestItem{
			{ID: 1, Name: "Alice", Age: 30},
			{ID: 2, Name: "Bob", Age: 25},
		}

		if !reflect.DeepEqual(paginatedItems, expectedPaginated) {
			t.Errorf("ProcessPage() paginatedItems = %v, want %v", paginatedItems, expectedPaginated)
		}

		if totalFiltered != 4 {
			t.Errorf("ProcessPage() totalFiltered = %d, want 4", totalFiltered)
		}

		if totalSource != 5 {
			t.Errorf("ProcessPage() totalSource = %d, want 5", totalSource)
		}
	})

	t.Run("sort function error", func(t *testing.T) {
		filterFn := func(item TestItem) bool {
			return true
		}

		sortFn := func(items []TestItem, sortSpec sdk.SortSpec) error {
			return errors.New("sort failed")
		}

		sortSpec := sdk.SortSpec{}
		_, _, _, err := ProcessPage(
			"TestItem",
			&sourceItems,
			sortSpec,
			sortFn,
			filterFn,
			0,
			2,
		)

		if err == nil {
			t.Error("ProcessPage() error = nil, want error")
		}

		expectedError := "error sorting TestItem: sort failed"
		if err.Error() != expectedError {
			t.Errorf("ProcessPage() error = %v, want %v", err.Error(), expectedError)
		}
	})

	t.Run("empty source slice", func(t *testing.T) {
		emptyItems := []TestItem{}
		filterFn := func(item TestItem) bool {
			return true
		}

		sortFn := func(items []TestItem, sortSpec sdk.SortSpec) error {
			return nil
		}

		sortSpec := sdk.SortSpec{}
		paginatedItems, totalFiltered, totalSource, err := ProcessPage(
			"TestItem",
			&emptyItems,
			sortSpec,
			sortFn,
			filterFn,
			0,
			2,
		)

		if err != nil {
			t.Errorf("ProcessPage() error = %v, want nil", err)
		}

		if len(paginatedItems) != 0 {
			t.Errorf("ProcessPage() paginatedItems length = %d, want 0", len(paginatedItems))
		}

		if totalFiltered != 0 {
			t.Errorf("ProcessPage() totalFiltered = %d, want 0", totalFiltered)
		}

		if totalSource != 0 {
			t.Errorf("ProcessPage() totalSource = %d, want 0", totalSource)
		}
	})

	t.Run("filter excludes all items", func(t *testing.T) {
		filterFn := func(item TestItem) bool {
			return false // Exclude all items
		}

		sortFn := func(items []TestItem, sortSpec sdk.SortSpec) error {
			return nil
		}

		sortSpec := sdk.SortSpec{}
		paginatedItems, totalFiltered, totalSource, err := ProcessPage(
			"TestItem",
			&sourceItems,
			sortSpec,
			sortFn,
			filterFn,
			0,
			2,
		)

		if err != nil {
			t.Errorf("ProcessPage() error = %v, want nil", err)
		}

		if len(paginatedItems) != 0 {
			t.Errorf("ProcessPage() paginatedItems length = %d, want 0", len(paginatedItems))
		}

		if totalFiltered != 0 {
			t.Errorf("ProcessPage() totalFiltered = %d, want 0", totalFiltered)
		}

		if totalSource != 5 {
			t.Errorf("ProcessPage() totalSource = %d, want 5", totalSource)
		}
	})
}
