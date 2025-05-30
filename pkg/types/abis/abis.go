// ADD_ROUTE
package abis

import (
	"fmt"
	"sync"

	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

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

type AbisCollection struct {
	App       types.App
	mutex     sync.RWMutex
	deduper   map[string]struct{}
	isLoading int32

	isDownloadedLoaded bool
	expectedDownloaded int
	downloadedAbis     []coreTypes.Abi

	isKnownLoaded bool
	expectedKnown int
	knownAbis     []coreTypes.Abi

	isFuncsLoaded     bool
	expectedFunctions int
	allFunctions      []coreTypes.Function

	isEventsLoaded bool
	expectedEvents int
	allEvents      []coreTypes.Function
}

func NewAbisCollection(app types.App) AbisCollection {
	return AbisCollection{
		App:            app,
		deduper:        make(map[string]struct{}),
		downloadedAbis: make([]coreTypes.Abi, 0),
		knownAbis:      make([]coreTypes.Abi, 0),
		allFunctions:   make([]coreTypes.Function, 0),
		allEvents:      make([]coreTypes.Function, 0),
	}
}

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

// NeedsUpdate checks if the AbisCollection needs to be updated.
func (ac *AbisCollection) NeedsUpdate(listKind types.ListKind) bool {
	ac.mutex.RLock()
	defer ac.mutex.RUnlock()

	switch listKind {
	case AbisDownloaded:
		return !ac.isDownloadedLoaded
	case AbisKnown:
		return !ac.isKnownLoaded
	case AbisFunctions:
		return !ac.isFuncsLoaded
	case AbisEvents:
		return !ac.isEventsLoaded
	default:
		return true
	}
}

// ADD_ROUTE
