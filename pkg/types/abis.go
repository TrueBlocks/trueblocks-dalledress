// ADD_ROUTE
package types

import (
	"sync"

	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
)

type AbisCollection struct {
	App App

	mutex sync.RWMutex

	isLoading          bool
	isLoaded           bool
	expectedFunctions  int
	expectedEvents     int
	expectedDownloaded int
	expectedKnown      int

	downloadedAbis []coreTypes.Abi
	knownAbis      []coreTypes.Abi
	allFunctions   []coreTypes.Function
	allEvents      []coreTypes.Function
	seenEncodings  map[string]bool
}

func NewAbisCollection(app App) AbisCollection {
	return AbisCollection{
		App:            app,
		downloadedAbis: make([]coreTypes.Abi, 0),
		knownAbis:      make([]coreTypes.Abi, 0),
		allFunctions:   make([]coreTypes.Function, 0),
		allEvents:      make([]coreTypes.Function, 0),
		seenEncodings:  make(map[string]bool),
	}
}

// ADD_ROUTE
