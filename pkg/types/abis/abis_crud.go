// ABIS_ROUTE
package abis

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/logging"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

func (ac *AbisCollection) Crud(
	listKind types.ListKind,
	op crud.Operation,
	abi *coreTypes.Abi,
) error {
	switch op {
	case crud.Remove:
		opts := sdk.AbisOptions{
			Addrs:   []string{abi.Address.Hex()},
			Globals: sdk.Globals{Decache: true},
		}
		if _, _, err := opts.Abis(); err != nil {
			return err
		}

		actionFunc := func(itemMatched *coreTypes.Abi) (error, bool) {
			return nil, true
		}
		matchFunc := func(existing *coreTypes.Abi) bool {
			return existing.Address == abi.Address
		}
		removedCount, err := ac.downloadedFacet.ForEvery(actionFunc, matchFunc)

		if err != nil {
			return fmt.Errorf("failed to remove ABI from facet: %w", err)
		}

		if removedCount > 0 {
			msgs.EmitStatus(fmt.Sprintf("deleted ABI for address: %s", abi.Address))
		} else {
			msgs.EmitStatus(fmt.Sprintf("ABI for address %s was not found in cache", abi.Address))
		}
		logging.LogBackend(fmt.Sprintf("Deleted ABI for address: %s", abi.Address))
		return nil

	default:
		logging.LogBackend(fmt.Sprintf("ABI operation %s not implemented for address: %s", op, abi.Address))
		return fmt.Errorf("operation %s not yet implemented for ABIs", op)
	}
}

// TODO: Consider adding batch operations for Abis, similar to MonitorsCollection.Clean (e.g., batch remove).
// ABIS_ROUTE
