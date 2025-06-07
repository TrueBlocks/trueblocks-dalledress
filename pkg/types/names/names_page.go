package names

import (
	"sort"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/sorting"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

// GetPage returns a page of names for the given list type and the total count.
func (n *NamesCollection) GetPage(
	listKind string,
	first, pageSize int,
	sortSpec sdk.SortSpec,
	filter string,
) (NamesPage, error) {
	if len(n.List) == 0 {
		if err := n.LoadData(nil); err != nil {
			return NamesPage{Names: nil, Total: 0, Tags: []string{}}, err
		}
	}

	namesMutex.Lock()
	defer namesMutex.Unlock()

	var list []*types.Name
	switch listKind {
	case "custom":
		list = n.Custom
	case "prefund":
		list = n.Prefund
	case "regular":
		list = n.Regular
	case "baddress":
		list = n.Baddress
	default:
		list = n.List
	}

	var filters []string
	if filter != "" {
		filters = append(filters, filter)
	}

	if selectedTag := n.GetSelectedTag(listKind); selectedTag != "" {
		filters = append(filters, selectedTag)
	}

	if len(filters) > 0 {
		filtered := make([]*types.Name, 0, len(list))
		for _, name := range list {
			addrHex := strings.ToLower(name.Address.Hex())
			addrNoPrefix := strings.TrimPrefix(addrHex, "0x")
			addrNoLeadingZeros := strings.TrimLeft(addrNoPrefix, "0")

			matchesAllFilters := true
			for _, filter := range filters {
				f := strings.ToLower(filter)
				match := strings.Contains(strings.ToLower(name.Name), f) ||
					strings.Contains(addrHex, f) ||
					strings.Contains(addrNoPrefix, f) ||
					strings.Contains(addrNoLeadingZeros, f) ||
					strings.Contains(strings.ToLower(name.Tags), f) ||
					strings.Contains(strings.ToLower(name.Source), f)

				// Extra: if filter starts with 0x, try matching without 0x and leading zeros
				if !match && strings.HasPrefix(f, "0x") {
					fNoPrefix := strings.TrimPrefix(f, "0x")
					if strings.Contains(addrNoPrefix, fNoPrefix) || strings.Contains(addrNoLeadingZeros, fNoPrefix) {
						match = true
					}
				}

				if !match {
					matchesAllFilters = false
					break
				}
			}

			if matchesAllFilters {
				filtered = append(filtered, name)
			}
		}
		list = filtered
	}

	total := len(list)
	if total == 0 || first >= total {
		var tags []string
		switch listKind {
		case "custom":
			tags = n.CustomTags
		case "prefund":
			tags = n.PrefundTags
		case "regular":
			tags = n.RegularTags
		case "baddress":
			tags = n.BaddressTags
		default:
			tags = n.ListTags
		}
		return NamesPage{Names: nil, Total: total, Tags: tags}, nil
	}

	// Sorting
	if !sorting.IsEmptySort(sortSpec) {
		sort.SliceStable(list, func(i, j int) bool {
			var vi, vj string
			switch sorting.GetSortField(sortSpec) {
			case "name":
				vi, vj = list[i].Name, list[j].Name
			case "address":
				vi, vj = list[i].Address.Hex(), list[j].Address.Hex()
			case "tags":
				vi, vj = list[i].Tags, list[j].Tags
			case "source":
				vi, vj = list[i].Source, list[j].Source
			default:
				vi, vj = list[i].Name, list[j].Name
			}
			if sorting.GetSortDirection(sortSpec) == "desc" {
				return vi > vj
			}
			return vi < vj
		})
	}

	var tags []string
	switch listKind {
	case "custom":
		tags = n.CustomTags
	case "prefund":
		tags = n.PrefundTags
	case "regular":
		tags = n.RegularTags
	case "baddress":
		tags = n.BaddressTags
	default:
		tags = n.ListTags
	}

	last := min(total, first+pageSize)
	return NamesPage{
		Names: list[first:last],
		Total: total,
		Tags:  tags,
	}, nil
}
