package streaming

import (
	"fmt"
	"strings"

	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

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

// ProcessPage processes the source slice: filters, sorts, and paginates the data
func ProcessPage[T any](
	typeName string,
	sourceSlice *[]T,
	sortSpec sdk.SortSpec,
	sortFn func([]T, sdk.SortSpec) error,
	filterFn func(T) bool,
	first, pageSize int,
) (paginatedItems []T, totalFiltered, totalSource int, err error) {
	filteredItems := FilterPageSlice(sourceSlice, filterFn)

	if err := sortFn(filteredItems, sortSpec); err != nil {
		return nil, 0, 0, fmt.Errorf("error sorting %s: %w", typeName, err)
	}

	paginatedItems = PaginateSlice(filteredItems, first, pageSize)

	return paginatedItems, len(filteredItems), len(*sourceSlice), nil
}
