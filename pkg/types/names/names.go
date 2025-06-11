// NAMES_ROUTE
package names

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
)

type namesMap map[base.Address]types.Name
type NamesCollection struct {
	Map      namesMap      `json:"map"`
	List     []*types.Name `json:"list"`
	Custom   []*types.Name `json:"custom"`
	Prefund  []*types.Name `json:"prefund"`
	Regular  []*types.Name `json:"regular"`
	Baddress []*types.Name `json:"baddress"`
}

func NewNamesCollection() NamesCollection {
	return NamesCollection{
		Map:      make(namesMap),
		List:     make([]*types.Name, 0),
		Custom:   make([]*types.Name, 0),
		Prefund:  make([]*types.Name, 0),
		Regular:  make([]*types.Name, 0),
		Baddress: make([]*types.Name, 0),
	}
}

type NamesPage struct {
	Names []*types.Name `json:"names"`
	Total int           `json:"total"`
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

// NAMES_ROUTE
