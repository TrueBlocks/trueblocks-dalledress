package facets

func (r *BaseFacet[T]) ForEvery(actionFunc func(itemMatched *T) (error, bool), matchFunc func(item *T) bool) (int, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	var matchCount int = 0
	filteredData := make([]T, 0, len(r.data))
	for _, item := range r.data {
		if !matchFunc(&item) {
			filteredData = append(filteredData, item)
		} else {
			matchCount++
		}
	}

	if matchCount > 0 {
		r.data = filteredData
		r.expectedCnt = len(r.data)
	}

	return matchCount, nil
}
