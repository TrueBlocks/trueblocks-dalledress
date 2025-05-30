package abis

import (
	"fmt"

	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

func (ac *AbisCollection) ClearCache(listKind types.ListKind) {
	ac.mutex.Lock()
	defer ac.mutex.Unlock()

	switch listKind {
	case AbisDownloaded:
		ac.expectedDownloaded = len(ac.downloadedAbis)
		ac.isDownloadedLoaded = false
		ac.downloadedAbis = make([]coreTypes.Abi, 0)
	case AbisKnown:
		ac.expectedKnown = len(ac.knownAbis)
		ac.isKnownLoaded = false
		ac.knownAbis = make([]coreTypes.Abi, 0)
	case AbisFunctions:
		ac.expectedFunctions = len(ac.allFunctions)
		ac.isFuncsLoaded = false
		ac.allFunctions = make([]coreTypes.Function, 0)
	case AbisEvents:
		ac.expectedEvents = len(ac.allEvents)
		ac.isEventsLoaded = false
		ac.allEvents = make([]coreTypes.Function, 0)
	default:
		ac.App.LogBackend(fmt.Sprintf("Unknown ListKind in ClearCache: %s", listKind))
	}
	ac.deduper = make(map[string]struct{})
}
