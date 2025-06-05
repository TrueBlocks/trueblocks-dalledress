package repository

import "fmt"

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
func (r *BaseRepository[T]) GetPage(first, pageSize int, filter FilterFunc[T], sortSpec interface{}, sortFunc func([]T, interface{}) error) (*PageResult[T], error) {
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
			HasMore:    false,
			IsLoaded:   r.state.IsLoaded(),
		}, nil
	}
	return &PageResult[T]{
		Items:      data[first:end],
		TotalItems: len(data),
		HasMore:    end < len(data),
		IsLoaded:   r.state.IsLoaded(),
	}, nil
}
