// ADD_ROUTE
package abis

import (
	"fmt"

	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/logging"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

// Remove deletes an ABI by address from the downloadedRepo and emits status
func (ac *AbisCollection) Remove(address string) error {
	opts := sdk.AbisOptions{
		Addrs:   []string{address},
		Globals: sdk.Globals{Decache: true},
	}
	if _, _, err := opts.Abis(); err != nil {
		return err
	}

	removed := ac.downloadedRepo.Remove(func(abi *coreTypes.Abi) bool {
		return abi.Address.Hex() == address
	})

	if removed {
		msgs.EmitStatus(fmt.Sprintf("deleted ABI for address: %s", address))
	} else {
		msgs.EmitStatus(fmt.Sprintf("ABI for address %s was not found in cache", address))
	}
	logging.LogBackend(fmt.Sprintf("Deleted ABI for address: %s", address))

	return nil
}

// ADD_ROUTE
