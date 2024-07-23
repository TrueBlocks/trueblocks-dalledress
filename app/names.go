package app

import (
	"sort"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/names"
)

func (a *App) GetNamesPage(first, pageSize int) []NameEx {
	if len(a.names) == 0 {
		return a.names
	}

	first = base.Max(0, base.Min(first, len(a.names)-1))
	last := base.Min(len(a.names), first+pageSize)
	return a.names[first:last]
}

func (a *App) GetNamesCnt() int {
	return len(a.names)
}

func (a *App) loadNames() error {
	types := []names.Parts{names.Regular, names.Custom, names.Prefund, names.Baddress}
	for _, t := range types {
		if m, err := names.LoadNamesMap("mainnet", t, nil); err != nil {
			return err
		} else {
			for addr, name := range m {
				namex := NameEx{
					Name: name,
					Type: t,
				}
				vv := a.namesMap[addr]
				namex.Type |= vv.Type
				a.namesMap[addr] = namex
			}
		}
	}
	for _, name := range a.namesMap {
		a.names = append(a.names, name)
	}
	sort.Slice(a.names, func(i, j int) bool {
		ti := a.names[i].Type
		tj := a.names[j].Type
		if ti == tj {
			if a.names[i].Tags == a.names[j].Tags {
				return a.names[i].Address.Hex() < a.names[j].Address.Hex()
			}
			return a.names[i].Tags < a.names[j].Tags
		}
		if ti == 4 || ti == 18 {
			return true
		}
		if tj == 4 || tj == 18 {
			return false
		}
		return ti < tj
	})
	return nil
}

func (a *App) GetNameTypes() []string {
	return []string{"Regular", "Custom", "Prefund", "Baddress"}
}
