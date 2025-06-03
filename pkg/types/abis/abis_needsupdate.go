// ADD_ROUTE
package abis

import "github.com/TrueBlocks/trueblocks-dalledress/pkg/types"

func (ac *AbisCollection) NeedsUpdate(listKind types.ListKind) bool {
	switch listKind {
	case AbisDownloaded:
		return ac.downloadedRepo.NeedsUpdate()
	case AbisKnown:
		return ac.knownRepo.NeedsUpdate()
	case AbisFunctions:
		return ac.functionsRepo.NeedsUpdate()
	case AbisEvents:
		return ac.eventsRepo.NeedsUpdate()
	default:
		return true
	}
}

// ADD_ROUTE
