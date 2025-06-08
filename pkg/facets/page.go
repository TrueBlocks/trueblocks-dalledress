package facets

import (
	"fmt"
	"strings"

	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

/*
Performance Note:
When filtering is applied, the current implementation creates a new slice
containing copies of all matching items, which can be expensive for large datasets.
An alternative approach would be to use indices or iterators instead:

1. Create a slice of indices that point to matching items in the original data
2. Use these indices to access the original data when returning the page
3. Only copy the specific items needed for the returned page

This would avoid the full copy during filtering while still protecting the original data.
For example:

	matchingIndices := []int{}
	for i, item := range data {
	    if filter(&item) {
	        matchingIndices = append(matchingIndices, i)
	    }
	}

	Then use matchingIndices[first:end] to access only the needed items
	and make copies just of those for the returned page

This approach trades memory efficiency for slightly more complex access patterns.
*/
// GetPage returns a filtered, sorted, and paginated page of data
func (r *BaseFacet[T]) GetPage(first, pageSize int, filter FilterFunc[T], sortSpec interface{}, sortFunc func([]T, interface{}) error) (*PageResult[T], error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	data := r.data
	if filter != nil {
		filtered := make([]T, 0, len(data))
		for _, item := range data {
			if filter(&item) {
				filtered = append(filtered, item)
			}
		}
		data = filtered
	}
	if sortFunc != nil {
		dataCopy := make([]T, len(data))
		copy(dataCopy, data)
		if err := sortFunc(dataCopy, sortSpec); err != nil {
			return nil, fmt.Errorf("error sorting data: %w", err)
		}
		data = dataCopy
	}
	if first < 0 || pageSize <= 0 {
		return nil, fmt.Errorf("invalid pagination parameters")
	}
	end := first + pageSize
	if end > len(data) {
		end = len(data)
	}
	if first >= len(data) {
		return &PageResult[T]{
			Items:      []T{},
			TotalItems: len(data),
			State:      r.GetState(),
		}, nil
	}
	return &PageResult[T]{
		Items:      data[first:end],
		TotalItems: len(data),
		State:      r.GetState(),
	}, nil
}

// FilterPageSlice returns a filtered slice using the provided filter function
func FilterPageSlice[T any](items *[]T, filterFn func(T) bool) []T {
	filtered := make([]T, 0)
	for _, item := range *items {
		if filterFn(item) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

// PaginateSlice returns a paginated slice from items
func PaginateSlice[T any](items []T, first, pageSize int) []T {
	totalItems := len(items)
	start := first
	end := first + pageSize

	if start >= totalItems {
		return []T{}
	}
	if end > totalItems {
		end = totalItems
	}
	return items[start:end]
}

// SortPageSlice sorts items using the provided sort function
func SortPageSlice[T any](items []T, sortSpec sdk.SortSpec, sortFn func([]T, sdk.SortSpec) error) error {
	return sortFn(items, sortSpec)
}

// CreatePageFilter returns a filter function for a given search string and field extractor
func CreatePageFilter[T any](filter string, searchFields func(T) []string) func(T) bool {
	return func(item T) bool {
		if filter == "" {
			return true
		}
		fields := searchFields(item)
		for _, field := range fields {
			if strings.Contains(strings.ToLower(field), filter) {
				return true
			}
		}
		return false
	}
}

// ProcessPage processes the store slice: filters, sorts, and paginates the data
func ProcessPage[T any](
	typeName string,
	storeSlice *[]T,
	sortSpec sdk.SortSpec,
	sortFn func([]T, sdk.SortSpec) error,
	filterFn func(T) bool,
	first, pageSize int,
) (paginatedItems []T, totalFiltered, totalSource int, err error) {
	filteredItems := FilterPageSlice(storeSlice, filterFn)

	if err := sortFn(filteredItems, sortSpec); err != nil {
		return nil, 0, 0, fmt.Errorf("error sorting %s: %w", typeName, err)
	}

	paginatedItems = PaginateSlice(filteredItems, first, pageSize)

	return paginatedItems, len(filteredItems), len(*storeSlice), nil
}
