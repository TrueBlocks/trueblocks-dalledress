// ADD_ROUTE
package abis

import "github.com/TrueBlocks/trueblocks-dalledress/pkg/types"

// NeedsUpdate checks if the specified facet needs to be updated
func (ac *AbisCollection) NeedsUpdate(listKind types.ListKind) bool {
	switch listKind {
	case AbisDownloaded:
		return ac.downloadedFacet.NeedsUpdate()
	case AbisKnown:
		return ac.knownFacet.NeedsUpdate()
	case AbisFunctions:
		return ac.functionsFacet.NeedsUpdate()
	case AbisEvents:
		return ac.eventsFacet.NeedsUpdate()
	default:
		return true
	}
}

// ADD_ROUTE
