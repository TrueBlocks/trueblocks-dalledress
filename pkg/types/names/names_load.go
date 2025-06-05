package names

import (
	"fmt"
	"path/filepath"
	"sort"
	"sync"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/config"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/names"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
)

var namesMutex sync.Mutex

// LoadData loads the names database into the NamesCollection struct and populates all category lists.
func (n *NamesCollection) LoadData(wg *sync.WaitGroup) error {
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
