package names

func (n *NamesCollection) ClearCaches() NamesCollection {
	ret := NamesCollection{selectedTags: n.selectedTags}
	_ = ret.LoadData(nil)
	return ret
}
