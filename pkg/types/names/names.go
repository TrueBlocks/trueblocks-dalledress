// ADD_ROUTE
package names

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/config"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/names"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/sorting"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

type namesMap map[base.Address]types.Name
type NamesCollection struct {
	Map          namesMap      `json:"map"`
	List         []*types.Name `json:"list"`
	Custom       []*types.Name `json:"custom"`
	Prefund      []*types.Name `json:"prefund"`
	Regular      []*types.Name `json:"regular"`
	Baddress     []*types.Name `json:"baddress"`
	ListTags     []string      `json:"listTags"`
	CustomTags   []string      `json:"customTags"`
	PrefundTags  []string      `json:"prefundTags"`
	RegularTags  []string      `json:"regularTags"`
	BaddressTags []string      `json:"baddressTags"`
	selectedTags map[string]string
}

func NewNamesCollection() NamesCollection {
	return NamesCollection{
		Map:          make(namesMap),
		List:         make([]*types.Name, 0),
		Custom:       make([]*types.Name, 0),
		Prefund:      make([]*types.Name, 0),
		Regular:      make([]*types.Name, 0),
		Baddress:     make([]*types.Name, 0),
		ListTags:     make([]string, 0),
		CustomTags:   make([]string, 0),
		PrefundTags:  make([]string, 0),
		RegularTags:  make([]string, 0),
		BaddressTags: make([]string, 0),
	}
}

type NamesPage struct {
	Names []*types.Name `json:"names"`
	Total int           `json:"total"`
	Tags  []string      `json:"tags"`
}

// GetSelectedTag returns the currently selected tag for a list type
func (n *NamesCollection) GetSelectedTag(key string) string {
	if n.selectedTags == nil {
		n.selectedTags = make(map[string]string)
	}

	return n.selectedTags[key]
}

// SetSelectedTag sets the selected tag for a specific list type
func (n *NamesCollection) SetSelectedTag(key string, tag string) {
	if n.selectedTags == nil {
		n.selectedTags = make(map[string]string)
	}

	if tag == "" {
		delete(n.selectedTags, key)
	} else {
		n.selectedTags[key] = tag
	}
}

// ClearSelectedTag clears the selected tag for a specific list type
func (n *NamesCollection) ClearSelectedTag(key string) {
	if n.selectedTags != nil {
		delete(n.selectedTags, key)
	}
}

// GetPage returns a page of names for the given list type and the total count.
func (n *NamesCollection) GetPage(listKind string, first, pageSize int, sortSpec sdk.SortSpec, filter string) NamesPage {
	if len(n.List) == 0 {
		if err := n.LoadNames(nil); err != nil {
			return NamesPage{Names: nil, Total: 0, Tags: []string{}}
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
		return NamesPage{Names: nil, Total: total, Tags: tags}
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
	}
}

var namesMutex sync.Mutex

// LoadNames loads the names database into the NamesCollection struct and populates all category lists.
func (n *NamesCollection) LoadNames(wg *sync.WaitGroup) error {
	namesMutex.Lock()
	defer func() {
		if wg != nil {
			wg.Done()
		}
		namesMutex.Unlock()
	}()

	chain := "mainnet"
	filePath := filepath.Join(config.MustGetPathToChainConfig(chain), string(names.DatabaseCustom))
	lineCount, _ := file.WordCount(filePath, true)

	customCount := 0
	for _, name := range n.List {
		if name.Parts&types.Custom != 0 {
			customCount++
		} else {
			break
		}
	}

	if lineCount == customCount {
		return nil
	}

	names.ClearCustomNames()

	parts := types.Regular | types.Custom | types.Prefund | types.Baddress
	if namesMap, err := names.LoadNamesMap(chain, parts, nil); err != nil {
		return err
	} else if (namesMap == nil) || (len(namesMap) == 0) {
		return fmt.Errorf("no names found")
	} else {
		n.Map = namesMap
		n.List = make([]*types.Name, 0, len(namesMap))
		n.Custom = n.Custom[:0]
		n.Prefund = n.Prefund[:0]
		n.Regular = n.Regular[:0]
		n.Baddress = n.Baddress[:0]
		for _, name := range n.Map {
			n.List = append(n.List, &name)
		}
		sort.Slice(n.List, func(i, j int) bool {
			return compare(*n.List[i], *n.List[j])
		})
		for _, name := range n.List {
			switch {
			case name.Parts&types.Custom != 0:
				n.Custom = append(n.Custom, name)
			case name.Parts&types.Prefund != 0:
				n.Prefund = append(n.Prefund, name)
			case name.Parts&types.Baddress != 0:
				n.Baddress = append(n.Baddress, name)
			case name.Parts&types.Regular != 0:
				n.Regular = append(n.Regular, name)
			}
		}

		n.ListTags = extractTagsFromNames(n.List)
		n.CustomTags = extractTagsFromNames(n.Custom)
		n.PrefundTags = extractTagsFromNames(n.Prefund)
		n.RegularTags = extractTagsFromNames(n.Regular)
		n.BaddressTags = extractTagsFromNames(n.Baddress)

		return nil
	}
}

func compare(nameI, nameJ types.Name) bool {
	ti := nameI.Parts
	if ti == types.Regular {
		ti = 7
	}
	tj := nameJ.Parts
	if tj == types.Regular {
		tj = 7
	}
	if ti == tj {
		if nameI.Tags == nameJ.Tags {
			return nameI.Address.Hex() < nameJ.Address.Hex()
		}
		return nameI.Tags < nameJ.Tags
	}
	return ti < tj
}

func (n *NamesCollection) ReloadNames() NamesCollection {
	ret := NamesCollection{selectedTags: n.selectedTags}
	_ = ret.LoadNames(nil)
	return ret
}

// extractTagsFromNames extracts unique tags from a list of names
func extractTagsFromNames(namesList []*types.Name) []string {
	tagsMap := make(map[string]bool)
	for _, name := range namesList {
		if name.Tags != "" {
			// Split the tags string by commas and extract individual tags
			tagsList := strings.Split(name.Tags, ",")
			for _, tag := range tagsList {
				tag = strings.TrimSpace(tag)
				if tag != "" {
					tagsMap[tag] = true
				}
			}
		}
	}

	uniqueTags := make([]string, 0, len(tagsMap))
	for tag := range tagsMap {
		uniqueTags = append(uniqueTags, tag)
	}
	sort.Strings(uniqueTags)

	return uniqueTags
}

// ADD_ROUTE
