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
	isLoading int32

	isDownloadedLoaded bool
	expectedDownloaded int
	downloaded         []coreTypes.Abi
	downloadedMutex sync.RWMutex

	isKnownLoaded bool
	expectedKnown int
	known         []coreTypes.Abi
	knownMutex      sync.RWMutex

	isFuncsLoaded     bool
	expectedFunctions int
	functions         []coreTypes.Function
	functionsMutex  sync.RWMutex

	isEventsLoaded bool
	expectedEvents int
	events         []coreTypes.Function
	eventsMutex     sync.RWMutex
}

func NewAbisCollection(app types.App) AbisCollection {
	return AbisCollection{
		App:        app,
		downloaded: make([]coreTypes.Abi, 0),
		known:      make([]coreTypes.Abi, 0),
		functions:  make([]coreTypes.Function, 0),
		events:     make([]coreTypes.Function, 0),
	}
}

func (ac *AbisCollection) ClearCache(listKind types.ListKind) {
	switch listKind {
	case AbisDownloaded:
		ac.downloadedMutex.Lock()
		defer ac.downloadedMutex.Unlock()
		ac.expectedDownloaded = len(ac.downloaded)
		ac.isDownloadedLoaded = false
		ac.downloaded = make([]coreTypes.Abi, 0)
	case AbisKnown:
		ac.knownMutex.Lock()
		defer ac.knownMutex.Unlock()
		ac.expectedKnown = len(ac.known)
		ac.isKnownLoaded = false
		ac.known = make([]coreTypes.Abi, 0)
	case AbisFunctions:
		ac.functionsMutex.Lock()
		defer ac.functionsMutex.Unlock()
		ac.expectedFunctions = len(ac.functions)
		ac.isFuncsLoaded = false
		ac.functions = make([]coreTypes.Function, 0)
	case AbisEvents:
		ac.eventsMutex.Lock()
		defer ac.eventsMutex.Unlock()
		ac.expectedEvents = len(ac.events)
		ac.isEventsLoaded = false
		ac.events = make([]coreTypes.Function, 0)
	default:
		ac.App.LogBackend(fmt.Sprintf("Unknown ListKind in ClearCache: %s", listKind))
	}
}

// NeedsUpdate checks if the AbisCollection needs to be updated.
func (ac *AbisCollection) NeedsUpdate(listKind types.ListKind) bool {
	switch listKind {
	case AbisDownloaded:
		ac.downloadedMutex.RLock()
		defer ac.downloadedMutex.RUnlock()
		return !ac.isDownloadedLoaded
	case AbisKnown:
		ac.knownMutex.RLock()
		defer ac.knownMutex.RUnlock()
		return !ac.isKnownLoaded
	case AbisFunctions:
		ac.functionsMutex.RLock()
		defer ac.functionsMutex.RUnlock()
		return !ac.isFuncsLoaded
	case AbisEvents:
		ac.eventsMutex.RLock()
		defer ac.eventsMutex.RUnlock()
		return !ac.isEventsLoaded
	default:
		return true
	}
}

// ADD_ROUTE
