// ADD_ROUTE
package abis

import (
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

// NeedsUpdate checks if the specified facet needs to be updated
func (ac *AbisCollection) NeedsUpdate(listKind types.ListKind) bool {
	var result bool

	switch listKind {
	case AbisDownloaded:
		result = ac.downloadedFacet.NeedsUpdate()
	case AbisKnown:
		result = ac.knownFacet.NeedsUpdate()
	case AbisFunctions:
		result = ac.functionsFacet.NeedsUpdate()
	case AbisEvents:
		result = ac.eventsFacet.NeedsUpdate()
	default:
		result = true
	}

	return result
}

func (a *AbisCollection) Reset(listKind types.ListKind) {
	switch listKind {
	case AbisDownloaded:
		a.downloadedFacet.Reset()
	case AbisKnown:
		a.knownFacet.Reset()
	case AbisFunctions:
		a.functionsFacet.Reset()
	case AbisEvents:
		a.eventsFacet.Reset()
	}
}

// ADD_ROUTE
