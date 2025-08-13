package types

import "sort"

// NormalizeOrders sorts columns and detail fields by their explicit order values.
// It does not assign defaults; ordering must be provided in config.
func NormalizeOrders(vc *ViewConfig) {
	if vc == nil || vc.Facets == nil {
		return
	}
	for key, facet := range vc.Facets {
		// Columns: sort by Order when both orders are positive; otherwise keep input order
		sort.SliceStable(facet.Columns, func(i, j int) bool {
			oi, oj := facet.Columns[i].Order, facet.Columns[j].Order
			if oi > 0 && oj > 0 {
				return oi < oj
			}
			return false
		})

		// Detail panels: within each panel, sort by DetailOrder when both orders are positive
		for p := range facet.DetailPanels {
			sort.SliceStable(facet.DetailPanels[p].Fields, func(i, j int) bool {
				oi, oj := facet.DetailPanels[p].Fields[i].DetailOrder, facet.DetailPanels[p].Fields[j].DetailOrder
				if oi > 0 && oj > 0 {
					return oi < oj
				}
				return false
			})
		}

		vc.Facets[key] = facet
	}
}
