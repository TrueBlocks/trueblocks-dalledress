// NAMES_ROUTE
package names

func (n *NamesCollection) ClearCache() NamesCollection {
	ret := NamesCollection{}
	_ = ret.LoadData(nil)
	return ret
}

// NAMES_ROUTE
