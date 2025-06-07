// ADD_ROUTE
package abis

import (
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/facets"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

// ListKind constants for ABI types
const (
	AbisDownloaded types.ListKind = "Downloaded"
	AbisKnown      types.ListKind = "Known"
	AbisFunctions  types.ListKind = "Functions"
	AbisEvents     types.ListKind = "Events"
)

func init() {
	types.RegisterKind(AbisDownloaded)
	types.RegisterKind(AbisKnown)
	types.RegisterKind(AbisFunctions)
	types.RegisterKind(AbisEvents)
}

// AbisCollection orchestrates all ABI facets
// (Downloaded, Known, Functions, Events)
type AbisCollection struct {
	downloadedFacet facets.Facet[coreTypes.Abi]
	knownFacet      facets.Facet[coreTypes.Abi]
	functionsFacet  facets.Facet[coreTypes.Function]
	eventsFacet     facets.Facet[coreTypes.Function]
}

// NewAbisCollection demonstrates the new Source → Facet pattern
// This creates an AbisCollection where multiple facets share the same data sources,
// eliminating redundant SDK queries:
// - Downloaded + Known facets share ONE AbisList source
// - Functions + Events facets share ONE AbisDetails source
// Total: 2 SDK queries instead of 4!
func NewAbisCollection() AbisCollection {
	sharedAbisListSource := GetSharedAbisListSource()
	sharedAbisDetailsSource := GetSharedAbisDetailsSource()

	downloadedFacet := facets.NewBaseFacet(
		AbisDownloaded,
		func(item *coreTypes.Abi) bool { return !item.IsKnown },
		nil,
		sharedAbisListSource,
	)

	knownFacet := facets.NewBaseFacet(
		AbisKnown,
		func(item *coreTypes.Abi) bool { return item.IsKnown },
		nil,
		sharedAbisListSource,
	)

	functionsFacet := facets.NewBaseFacet(
		AbisFunctions,
		func(item *coreTypes.Function) bool { return item.FunctionType != "event" },
		IsDupFuncByEncoding(),
		sharedAbisDetailsSource,
	)

	eventsFacet := facets.NewBaseFacet(
		AbisEvents,
		func(item *coreTypes.Function) bool { return item.FunctionType == "event" },
		IsDupFuncByEncoding(),
		sharedAbisDetailsSource,
	)

	return AbisCollection{
		downloadedFacet: downloadedFacet,
		knownFacet:      knownFacet,
		functionsFacet:  functionsFacet,
		eventsFacet:     eventsFacet,
	}
}

// ADD_ROUTE
