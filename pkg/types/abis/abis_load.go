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

// LoadData demonstrates using sources for data loading (Phase 2 pattern)
// This method shows how to use the enhanced source-based loading where multiple facets
// share the same data sources, eliminating redundant SDK queries.
func (ac *AbisCollection) LoadData(listKind types.ListKind) {
	if !ac.NeedsUpdate(listKind) {
		return
	}

	switch listKind {
	case AbisDownloaded:
		go func() {
			if result, err := ac.downloadedFacet.Load(facets.LoadOptions{}); err != nil {
				logging.LogError("LoadData.AbisDownloaded from source: %v", err, facets.ErrorAlreadyLoading)
			} else {
				msgs.EmitLoaded(result.Payload.Reason, result.Payload)
			}
		}()
	case AbisKnown:
		go func() {
			if result, err := ac.knownFacet.Load(facets.LoadOptions{}); err != nil {
				logging.LogError("LoadData.AbisKnown from source: %v", err, facets.ErrorAlreadyLoading)
			} else {
				msgs.EmitLoaded(result.Payload.Reason, result.Payload)
			}
		}()
	case AbisFunctions:
		go func() {
			if result, err := ac.functionsFacet.Load(facets.LoadOptions{}); err != nil {
				logging.LogError("LoadData.AbisFunction from source: %v", err, facets.ErrorAlreadyLoading)
			} else {
				msgs.EmitLoaded(result.Payload.Reason, result.Payload)
			}
		}()
	case AbisEvents:
		go func() {
			if result, err := ac.eventsFacet.Load(facets.LoadOptions{}); err != nil {
				logging.LogError("LoadData.AbisEvents from source: %v", err, facets.ErrorAlreadyLoading)
			} else {
				msgs.EmitLoaded(result.Payload.Reason, result.Payload)
			}
		}()
	default:
		logger.Error(fmt.Sprintf("LoadData: unexpected list kind: %v", listKind))
	}
}

// ADD_ROUTE
