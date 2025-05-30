// ADD_ROUTE
package abis

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

func (ac *AbisCollection) Delete(address string) error {
	opts := sdk.AbisOptions{
		Addrs: []string{address},
		Globals: sdk.Globals{
			Decache: true,
		},
	}
	if _, _, err := opts.Abis(); err != nil {
		return err
	}

	ac.mutex.Lock()
	defer ac.mutex.Unlock()
	for i, abi := range ac.downloadedAbis {
		if abi.Address.Hex() == address {
			ac.downloadedAbis = append(ac.downloadedAbis[:i], ac.downloadedAbis[i+1:]...)
			msgs.EmitStatus(fmt.Sprintf("Deleted ABI for address: %s", address))
			ac.App.LogBackend(fmt.Sprintf("Deleted ABI for address: %s", address))
			return nil
		}
	}

	return fmt.Errorf("ABI with address %s not found", address)
}

// ADD_ROUTE
