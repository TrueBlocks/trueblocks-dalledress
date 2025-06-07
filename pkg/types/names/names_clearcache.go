package names

func (n *NamesCollection) ClearCache() NamesCollection {
	ret := NamesCollection{selectedTags: n.selectedTags}
	_ = ret.LoadData(nil)
	return ret
}
