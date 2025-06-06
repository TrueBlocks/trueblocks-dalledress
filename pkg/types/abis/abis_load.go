// ADD_ROUTE
package abis

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/facets"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/logging"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

func (ac *AbisCollection) LoadData(listKind types.ListKind) {
	if !ac.NeedsUpdate(listKind) {
		return
	}

	switch listKind {
	case AbisDownloaded:
		go ac.loadDownloaded()
	case AbisKnown:
		go ac.loadKnown()
	case AbisFunctions:
		go ac.loadFunctions()
	case AbisEvents:
		go ac.loadEvents()
	default:
		logger.Error(fmt.Sprintf("LoadData: unexpected list kind: %v", listKind))
	}
}

func (ac *AbisCollection) loadDownloaded() {
	if result, err := ac.downloadedFacet.Load(facets.LoadOptions{}); err != nil {
		logging.LogError("loadDownloaded: %v", err, facets.ErrorAlreadyLoading)
	} else {
		msgs.EmitLoaded(result.Payload.Reason, result.Payload)
	}
}

func (ac *AbisCollection) loadKnown() {
	if result, err := ac.knownFacet.Load(facets.LoadOptions{}); err != nil {
		logging.LogError("loadKnown: %v", err, facets.ErrorAlreadyLoading)
	} else {
		msgs.EmitLoaded(result.Payload.Reason, result.Payload)
	}
}

func (ac *AbisCollection) loadFunctions() {
	if result, err := ac.functionsFacet.Load(facets.LoadOptions{}); err != nil {
		logging.LogError("loadFunctions: %v", err, facets.ErrorAlreadyLoading)
	} else {
		msgs.EmitLoaded(result.Payload.Reason, result.Payload)
	}
}

func (ac *AbisCollection) loadEvents() {
	if result, err := ac.eventsFacet.Load(facets.LoadOptions{}); err != nil {
		logging.LogError("loadEvents: %v", err, facets.ErrorAlreadyLoading)
	} else {
		msgs.EmitLoaded(result.Payload.Reason, result.Payload)
	}
}

// ADD_ROUTE
