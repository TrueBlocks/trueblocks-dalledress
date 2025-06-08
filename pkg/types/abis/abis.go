package abis

import (
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/facets"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

const (
	AbisDownloaded types.ListKind = "Downloaded"
	AbisKnown      types.ListKind = "Known"
	AbisFunctions  types.ListKind = "Functions"
	AbisEvents     types.ListKind = "Events"
)

// Register list kinds
func init() {
	types.RegisterKind(AbisDownloaded)
	types.RegisterKind(AbisKnown)
	types.RegisterKind(AbisFunctions)
	types.RegisterKind(AbisEvents)
}

type AbisCollection struct {
	downloadedFacet facets.Facet[coreTypes.Abi]
	knownFacet      facets.Facet[coreTypes.Abi]
	functionsFacet  facets.Facet[coreTypes.Function]
	eventsFacet     facets.Facet[coreTypes.Function]
}

func NewAbisCollection() AbisCollection {
	abisListStore := GetListStore()
	abisDetailStore := GetDetailStore()

	downloadedFacet := facets.NewBaseFacet(
		AbisDownloaded,
		func(item *coreTypes.Abi) bool { return !item.IsKnown },
		nil,
		abisListStore,
	)

	knownFacet := facets.NewBaseFacet(
		AbisKnown,
		func(item *coreTypes.Abi) bool { return item.IsKnown },
		nil,
		abisListStore,
	)

	functionsFacet := facets.NewBaseFacet(
		AbisFunctions,
		func(item *coreTypes.Function) bool { return item.FunctionType != "event" },
		isEncodingDup(),
		abisDetailStore,
	)

	eventsFacet := facets.NewBaseFacet(
		AbisEvents,
		func(item *coreTypes.Function) bool { return item.FunctionType == "event" },
		isEncodingDup(),
		abisDetailStore,
	)

	return AbisCollection{
		downloadedFacet: downloadedFacet,
		knownFacet:      knownFacet,
		functionsFacet:  functionsFacet,
		eventsFacet:     eventsFacet,
	}
}
