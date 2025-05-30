// ADD_ROUTE
package abis

import (
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

// ADD_ROUTE
