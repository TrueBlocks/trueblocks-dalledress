// ADD_ROUTE
package abis

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

func (ac *AbisCollection) ClearCache(listKind types.ListKind) {
	switch listKind {
	case AbisDownloaded:
		ac.downloadedRepo.Clear()
	case AbisKnown:
		ac.knownRepo.Clear()
	case AbisFunctions:
		ac.functionsRepo.Clear()
	case AbisEvents:
		ac.eventsRepo.Clear()
	default:
		ac.App.LogBackend(fmt.Sprintf("Unknown ListKind in ClearCache: %s", listKind))
	}
}

// ADD_ROUTE
