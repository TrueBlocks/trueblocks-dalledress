// ADD_ROUTE
package abis

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/logging"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

func (ac *AbisCollection) ClearCache(listKind types.ListKind) {
	switch listKind {
	case AbisDownloaded:
		ac.downloadedFacet.Clear()
	case AbisKnown:
		ac.knownFacet.Clear()
	case AbisFunctions:
		ac.functionsFacet.Clear()
	case AbisEvents:
		ac.eventsFacet.Clear()
	default:
		logging.LogBackend(fmt.Sprintf("Unknown ListKind in ClearCache: %s", listKind))
	}
}

// ADD_ROUTE
