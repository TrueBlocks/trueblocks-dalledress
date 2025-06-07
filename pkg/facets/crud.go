package facets

// Remove deletes items matching the predicate from the facet data
func (r *BaseFacet[T]) Remove(predicate func(*T) bool) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	var removed bool
	filteredData := make([]T, 0, len(r.data))
	for _, item := range r.data {
		if !predicate(&item) {
			filteredData = append(filteredData, item)
		} else {
			removed = true
		}
	}

	if removed {
		r.data = filteredData
		r.expectedCnt = len(r.data) // Update expected count after removal
	}

	return removed
}
