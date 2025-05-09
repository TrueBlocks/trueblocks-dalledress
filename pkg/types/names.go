package types

import (
	"fmt"
	"log"
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
)

type Names struct {
	Map      map[base.Address]types.Name `json:"map"`
	List     []*types.Name               `json:"list"`
	Custom   []*types.Name               `json:"custom"`
	Prefund  []*types.Name               `json:"prefund"`
	Regular  []*types.Name               `json:"regular"`
	Baddress []*types.Name               `json:"baddress"`
}

type NamesPage struct {
	Names []*types.Name `json:"names"`
	Total int           `json:"total"`
}

// GetNamesPage returns a page of names for the given list type and the total count.
func (n *Names) GetNamesPage(listType string, first, pageSize int, sortKey sorting.SortDef, filter string) NamesPage {
	if len(n.List) == 0 {
		if err := n.LoadNames(nil); err != nil {
			return NamesPage{Names: nil, Total: 0}
		}
	}

	namesMutex.Lock()
	defer namesMutex.Unlock()

	var list []*types.Name
	switch listType {
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

	if filter != "" {
		filtered := make([]*types.Name, 0, len(list))
		f := strings.ToLower(filter)
		for _, name := range list {
			addrHex := strings.ToLower(name.Address.Hex())
			addrNoPrefix := strings.TrimPrefix(addrHex, "0x")
			addrNoLeadingZeros := strings.TrimLeft(addrNoPrefix, "0")
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
			if match {
				filtered = append(filtered, name)
			}
		}
		list = filtered
	}

	total := len(list)
	if total == 0 || first >= total {
		return NamesPage{Names: nil, Total: total}
	}

	// Sorting
	if sortKey.Key != "" {
		sort.SliceStable(list, func(i, j int) bool {
			var vi, vj string
			switch sortKey.Key {
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
			if sortKey.Direction == "desc" {
				return vi > vj
			}
			return vi < vj
		})
	}

	last := min(total, first+pageSize)
	return NamesPage{Names: list[first:last], Total: total}
}

var namesMutex sync.Mutex

// LoadNames loads the names database into the Names struct and populates all category lists.
func (n *Names) LoadNames(wg *sync.WaitGroup) error {
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
		log.Printf("Loaded %d names", len(namesMap))
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

func (n *Names) ReloadNames() {
	*n = Names{}
	_ = n.LoadNames(nil)
}
