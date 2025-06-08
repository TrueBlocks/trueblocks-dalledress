package abis

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/facets"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/logging"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

// LoadData demonstrates using stores for data loading (Phase 2 pattern)
// This method shows how to use the enhanced store-based loading where multiple facets
// share the same data stores, eliminating redundant SDK queries.
func (ac *AbisCollection) LoadData(listKind types.ListKind) {
	if !ac.NeedsUpdate(listKind) {
		return
	}

	switch listKind {
	case AbisDownloaded:
		go func() {
			if result, err := ac.downloadedFacet.Load(); err != nil {
				logging.LogError("LoadData.AbisDownloaded from store: %v", err, facets.ErrorAlreadyLoading)
			} else {
				msgs.EmitLoaded("downloaded", result.Payload)
			}
		}()
	case AbisKnown:
		go func() {
			if result, err := ac.knownFacet.Load(); err != nil {
				logging.LogError("LoadData.AbisKnown from store: %v", err, facets.ErrorAlreadyLoading)
			} else {
				msgs.EmitLoaded("known", result.Payload)
			}
		}()
	case AbisFunctions:
		go func() {
			if result, err := ac.functionsFacet.Load(); err != nil {
				logging.LogError("LoadData.AbisFunction from store: %v", err, facets.ErrorAlreadyLoading)
			} else {
				msgs.EmitLoaded("functions", result.Payload)
			}
		}()
	case AbisEvents:
		go func() {
			if result, err := ac.eventsFacet.Load(); err != nil {
				logging.LogError("LoadData.AbisEvents from store: %v", err, facets.ErrorAlreadyLoading)
			} else {
				msgs.EmitLoaded("events", result.Payload)
			}
		}()
	default:
		logger.Error(fmt.Sprintf("LoadData: unexpected list kind: %v", listKind))
	}
}
