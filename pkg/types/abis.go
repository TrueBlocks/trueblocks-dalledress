// ADD_ROUTE
package types

import (
	"fmt"
	"sync"

	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
)

type AbisCollection struct {
	App       App
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

func NewAbisCollection(app App) AbisCollection {
	return AbisCollection{
		App:            app,
		deduper:        make(map[string]struct{}),
		downloadedAbis: make([]coreTypes.Abi, 0),
		knownAbis:      make([]coreTypes.Abi, 0),
		allFunctions:   make([]coreTypes.Function, 0),
		allEvents:      make([]coreTypes.Function, 0),
	}
}

func (ac *AbisCollection) ClearCache(listKind ListKind, loaded bool) {
	ac.mutex.Lock()
	defer ac.mutex.Unlock()

	switch listKind {
	case AbisDownloaded:
		ac.isDownloadedLoaded = loaded
	case AbisKnown:
		ac.isKnownLoaded = loaded
	case AbisFunctions:
		ac.isFuncsLoaded = loaded
	case AbisEvents:
		ac.isEventsLoaded = loaded
	default:
		ac.App.LogBackend(fmt.Sprintf("Unknown ListKind in ClearCache: %s", listKind))
	}
}

// ADD_ROUTE
