package repository

// Remove deletes items matching the predicate from the repository data
func (r *BaseRepository[T]) Remove(predicate func(*T) bool) bool {
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
		r.expectedCount = len(r.data) // Update expected count after removal
	}

	return removed
}
