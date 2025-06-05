// ADD_ROUTE
package names

import (
	"sort"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
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
