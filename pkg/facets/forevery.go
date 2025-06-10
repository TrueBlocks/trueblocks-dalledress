package facets

func (r *Facet[T]) ForEvery(actionFunc func(itemMatched *T) (error, bool), matchFunc func(item *T) bool) (int, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	var matchCount int = 0
	filteredData := make([]*T, 0, len(r.view))
	for _, itemPtr := range r.view {
		if !matchFunc(itemPtr) {
			filteredData = append(filteredData, itemPtr)
		} else {
			matchCount++
		}
	}

	if matchCount > 0 {
		r.view = filteredData
		r.expectedCnt = len(r.view)
	}

	return matchCount, nil
}
