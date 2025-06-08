package abis

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

// Crud performs CRUD operations on ABIs - handles all operation types
func (ac *AbisCollection) AbisCrud(listKind types.ListKind, op crud.Operation, abi *coreTypes.Abi) error {
	action := "unknown"
	nRemoved := 0

	switch op {
	case crud.Remove:
		action = "removed"
		actionFunc := func(item *coreTypes.Abi) (error, bool) {
			addr := item.Address.Hex()
			opts := sdk.AbisOptions{
				Addrs:   []string{addr},
				Globals: sdk.Globals{Decache: true},
			}
			if _, _, err := opts.Abis(); err != nil {
				return err, false /* don't remove */ // nolint:nilerr
			}
			return nil, true /* remove */
		}

		matchFunc := func(existing *coreTypes.Abi) bool {
			return existing.Address == abi.Address
		}

		nRemoved, _ = ac.downloadedFacet.ForEvery(actionFunc, matchFunc)
	}

	if nRemoved > 0 {
		msgs.EmitStatus(fmt.Sprintf("%s ABI for address in %s facet: %s", action, listKind, abi.Address))
	} else {
		msgs.EmitStatus(fmt.Sprintf("Address %s not found in %s facet", abi.Address, listKind))
	}

	return nil
}

