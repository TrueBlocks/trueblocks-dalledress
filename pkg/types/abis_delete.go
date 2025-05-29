// ADD_ROUTE
package types

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
)

func (ac *AbisCollection) Delete(address string) error {
	ac.mutex.Lock()
	defer ac.mutex.Unlock()

	for i, abi := range ac.downloadedAbis {
		if abi.Address.Hex() == address {
			ac.downloadedAbis = append(ac.downloadedAbis[:i], ac.downloadedAbis[i+1:]...)
			msgs.EmitStatus(fmt.Sprintf("Deleted downloaded ABI for address: %s", address))
			return nil
		}
	}

	return fmt.Errorf("ABI with address %s not found", address)
}
// ADD_ROUTE
