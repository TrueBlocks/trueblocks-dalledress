// ABIS_CODE
package types

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/sorting"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

var refreshRate = 91

type AbisCollection struct {
	App App

	mutex sync.RWMutex

	isLoading          bool
	isFullyLoaded      bool
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

// NewAbisCollection creates a new AbisCollection.
// Make sure to pass a valid App instance when creating AbisCollection.
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

// loadInternal is the private method that handles the actual data loading.
// It's designed to be run in a goroutine.
func (ac *AbisCollection) loadInternal() {
	defer func() {
		ac.mutex.Lock()
		ac.isLoading = false
		ac.mutex.Unlock()
	}()

	// At the start, reset flags and slices under a write lock
	ac.mutex.Lock()
	ac.downloadedAbis = make([]coreTypes.Abi, 0)
	ac.knownAbis = make([]coreTypes.Abi, 0)
	ac.allFunctions = make([]coreTypes.Function, 0)
	ac.allEvents = make([]coreTypes.Function, 0)
	ac.seenEncodings = make(map[string]bool) // Reset deduplication map
	ac.isFullyLoaded = false
	ac.expectedFunctions = 0
	ac.expectedEvents = 0
	ac.expectedDownloaded = 0
	ac.expectedKnown = 0

	ac.mutex.Unlock()

	listOpts := sdk.AbisOptions{
		Globals: sdk.Globals{Cache: true, Verbose: true},
	}
	sdkAbis, _, err := listOpts.AbisList()
	if err != nil {
		msgs.EmitError("AbisCollection.loadInternal: error fetching ABI list", err)
		return
	}

	if sdkAbis == nil {
		sdkAbis = make([]coreTypes.Abi, 0)
	}

	addrMap := make(map[base.Address]bool)
	ac.mutex.Lock()
	for _, abi := range sdkAbis {
		if abi.IsKnown {
			ac.knownAbis = append(ac.knownAbis, abi)
		} else {
			if _, exists := addrMap[abi.Address]; !exists {
				ac.downloadedAbis = append(ac.downloadedAbis, abi)
				ac.expectedFunctions += int(abi.NFunctions)
				ac.expectedEvents += int(abi.NEvents)
			}
			addrMap[abi.Address] = true
		}
	}
	ac.expectedDownloaded = len(ac.downloadedAbis)
	ac.expectedKnown = len(ac.knownAbis)
	currentCount := 0
	msgs.EmitStatus(fmt.Sprintf("Loading ABIs: %d downloaded, %d known. Fetching details...", len(ac.downloadedAbis), len(ac.knownAbis)))
	ac.mutex.Unlock()

	if len(addrMap) == 0 {
		logger.Info(fmt.Sprintln("AbisCollection.loadInternal: no non-known ABI addresses to fetch detailed functions/events."))
		ac.mutex.Lock()
		ac.isFullyLoaded = true
		statusMsg := "ABI details loaded: No new items to fetch."
		msgs.EmitStatus(statusMsg)
		ac.App.EmitEvent(msgs.EventDataLoaded, DataLoadedPayload{
			DataType:      "functions-events",
			CurrentCount:  currentCount,
			ExpectedTotal: ac.expectedFunctions + ac.expectedEvents,
			IsFullyLoaded: true,
			Category:      "abis",
		})
		ac.mutex.Unlock()
		return
	}

	abisAddr := base.ZeroAddr
	renderCtx := ac.App.RegisterCtx(abisAddr)

	// Ensure context is cancelled/cleaned up when we're done
	defer func() {
		ac.App.Cancel(abisAddr)
	}()

	terms := make([]string, 0, len(addrMap))
	for addr := range addrMap {
		terms = append(terms, addr.Hex())
	}
	detailOpts := sdk.AbisOptions{
		Globals:   sdk.Globals{Cache: true},
		Addrs:     terms,
		RenderCtx: renderCtx,
	}

	go func() {
		defer func() {
			if renderCtx.ModelChan != nil {
				close(renderCtx.ModelChan)
			}
			if renderCtx.ErrorChan != nil {
				close(renderCtx.ErrorChan)
			}
		}()

		_, _, streamInitiationErr := detailOpts.AbisListItems()
		if streamInitiationErr != nil {
			logger.Info(fmt.Sprintf("AbisCollection.loadInternal: error initiating stream: %v", streamInitiationErr))
		}
	}()

ProcessingLoop:
	for {
		select {
		case itemIntf, ok := <-renderCtx.ModelChan:
			if !ok {
				renderCtx.ModelChan = nil
				if renderCtx.ErrorChan == nil {
					break ProcessingLoop
				}
				continue
			}

			itemPtr, okAssert := itemIntf.(*coreTypes.Function)
			if !okAssert {
				logger.Info(fmt.Sprintf("AbisCollection.loadInternal: unexpected item type: %T", itemIntf))
				continue
			}

			ac.mutex.Lock()
			if ac.seenEncodings[itemPtr.Encoding] {
				ac.mutex.Unlock()
				continue
			}

			ac.seenEncodings[itemPtr.Encoding] = true
			if itemPtr.FunctionType == "event" {
				ac.allEvents = append(ac.allEvents, *itemPtr)
			} else {
				ac.allFunctions = append(ac.allFunctions, *itemPtr)
			}
			currentCount = len(ac.allFunctions) + len(ac.allEvents)
			ac.mutex.Unlock()

			if currentCount%refreshRate == 0 {
				ac.mutex.RLock()
				payload := DataLoadedPayload{
					DataType:      "functions-events",
					CurrentCount:  currentCount,
					ExpectedTotal: ac.expectedFunctions + ac.expectedEvents,
					IsFullyLoaded: false,
					Category:      "abis",
				}
				ac.App.EmitEvent(msgs.EventDataLoaded, payload)
				statusMsg := fmt.Sprintf("Loading ABI details: %d processed.", currentCount)
				if (ac.expectedFunctions + ac.expectedEvents) > 0 {
					statusMsg = fmt.Sprintf("Loading ABI details: %d of %d processed.", currentCount, ac.expectedFunctions+ac.expectedEvents)
				}
				msgs.EmitStatus(statusMsg)
				ac.mutex.RUnlock()
			}

		case streamErr, ok := <-renderCtx.ErrorChan:
			if !ok {
				renderCtx.ErrorChan = nil
				if renderCtx.ModelChan == nil {
					break ProcessingLoop
				}
				continue
			}
			msgs.EmitError("AbisCollection.loadInternal: streaming error", streamErr)
		}
	}

	ac.mutex.Lock()
	ac.isFullyLoaded = true
	finalCount := len(ac.allFunctions) + len(ac.allEvents)
	finalStatus := fmt.Sprintf("ABI details fully loaded: %d functions/events.", finalCount)
	msgs.EmitStatus(finalStatus)
	ac.App.EmitEvent(msgs.EventDataLoaded, DataLoadedPayload{
		DataType:      "functions-events",
		CurrentCount:  finalCount,
		ExpectedTotal: ac.expectedFunctions + ac.expectedEvents,
		IsFullyLoaded: true,
		Category:      "abis",
	})
	ac.mutex.Unlock()
}

func (ac *AbisCollection) EnsureInitialLoad() {
	ac.mutex.Lock()
	if !ac.isFullyLoaded && !ac.isLoading {
		ac.isLoading = true
		go ac.loadInternal()
	}
	ac.mutex.Unlock()
}

type AbisPage struct {
	Kind          string               `json:"Type"`
	Abis          []coreTypes.Abi      `json:"Abis,omitempty"`
	Functions     []coreTypes.Function `json:"Functions,omitempty"`
	TotalItems    int                  `json:"TotalItems"`
	IsLoading     bool                 `json:"IsLoading"`
	IsFullyLoaded bool                 `json:"IsFullyLoaded"`
	ExpectedTotal int                  `json:"ExpectedTotal"`
}

func (ac *AbisCollection) GetPage(kind string, first, pageSize int, sortDef *sorting.SortDef, filter string) (AbisPage, error) {
	ac.EnsureInitialLoad()

	ac.mutex.RLock()
	isLoadingSnapshot := ac.isLoading
	isFullyLoadedSnapshot := ac.isFullyLoaded
	expectedFunctionsSnapshot := ac.expectedFunctions
	expectedEventsSnapshot := ac.expectedEvents
	expectedDownloadedSnapshot := ac.expectedDownloaded
	expectedKnownSnapshot := ac.expectedKnown

	var currentAbis []coreTypes.Abi
	var currentFunctions []coreTypes.Function

	switch kind {
	case "Downloaded":
		currentAbis = make([]coreTypes.Abi, len(ac.downloadedAbis))
		copy(currentAbis, ac.downloadedAbis)
	case "Known":
		currentAbis = make([]coreTypes.Abi, len(ac.knownAbis))
		copy(currentAbis, ac.knownAbis)
	case "Functions":
		currentFunctions = make([]coreTypes.Function, len(ac.allFunctions))
		copy(currentFunctions, ac.allFunctions)
	case "Events":
		currentFunctions = make([]coreTypes.Function, len(ac.allEvents))
		copy(currentFunctions, ac.allEvents)
	default:
		ac.mutex.RUnlock()
		return AbisPage{}, fmt.Errorf("unknown ABI page kind: %s", kind)
	}
	ac.mutex.RUnlock()

	var expectedTotal int
	switch kind {
	case "Downloaded":
		expectedTotal = expectedDownloadedSnapshot
	case "Known":
		expectedTotal = expectedKnownSnapshot
	case "Functions":
		expectedTotal = expectedFunctionsSnapshot
	case "Events":
		expectedTotal = expectedEventsSnapshot
	}

	page := AbisPage{
		Kind:          kind,
		IsLoading:     isLoadingSnapshot,
		IsFullyLoaded: isFullyLoadedSnapshot,
		ExpectedTotal: expectedTotal,
	}
	filter = strings.ToLower(filter)

	switch kind {
	case "Downloaded", "Known":
		filteredAbis := make([]coreTypes.Abi, 0)
		for _, item := range currentAbis { // Use the copied slice
			if filter == "" || strings.Contains(strings.ToLower(item.Name), filter) || strings.Contains(strings.ToLower(item.Address.Hex()), filter) {
				filteredAbis = append(filteredAbis, item)
			}
		}
		if sortDef != nil && sortDef.Key != "" {
			sort.SliceStable(filteredAbis, func(i, j int) bool {
				valI, valJ := "", ""
				switch sortDef.Key {
				case "address":
					valI, valJ = filteredAbis[i].Address.Hex(), filteredAbis[j].Address.Hex()
				case "name":
					valI, valJ = filteredAbis[i].Name, filteredAbis[j].Name
				case "nFunctions":
					valI, valJ = fmt.Sprintf("%09d", filteredAbis[i].NFunctions), fmt.Sprintf("%09d", filteredAbis[j].NFunctions)
				case "nEvents":
					valI, valJ = fmt.Sprintf("%09d", filteredAbis[i].NEvents), fmt.Sprintf("%09d", filteredAbis[j].NEvents)
				case "fileSize":
					valI, valJ = fmt.Sprintf("%09d", filteredAbis[i].FileSize), fmt.Sprintf("%09d", filteredAbis[j].FileSize)
				case "lastModDate":
					valI, valJ = filteredAbis[i].LastModDate, filteredAbis[j].LastModDate
				default:
					return false // Should not happen if keys are validated
				}
				if sortDef.Direction == "desc" {
					return strings.ToLower(valI) > strings.ToLower(valJ)
				}
				return strings.ToLower(valI) < strings.ToLower(valJ)
			})
		}
		start := first
		end := first + pageSize
		if start < len(filteredAbis) {
			if end > len(filteredAbis) {
				end = len(filteredAbis)
			}
			page.Abis = filteredAbis[start:end]
		}
		page.TotalItems = len(filteredAbis) // Total items matching filter, not just current page

	case "Functions", "Events":
		filteredFunctions := make([]coreTypes.Function, 0)
		for _, item := range currentFunctions { // Use the copied slice
			name := strings.ToLower(item.Name)
			signature := strings.ToLower(item.Signature)
			encoding := strings.ToLower(item.Encoding)
			if filter == "" || strings.Contains(name, filter) || strings.Contains(signature, filter) || strings.Contains(encoding, filter) {
				filteredFunctions = append(filteredFunctions, item)
			}
		}
		if sortDef != nil && sortDef.Key != "" {
			sort.SliceStable(filteredFunctions, func(i, j int) bool {
				valI, valJ := "", ""
				switch sortDef.Key {
				case "name":
					valI, valJ = filteredFunctions[i].Name, filteredFunctions[j].Name
				case "signature":
					valI, valJ = filteredFunctions[i].Signature, filteredFunctions[j].Signature
				case "encoding":
					valI, valJ = filteredFunctions[i].Encoding, filteredFunctions[j].Encoding
				default:
					return false // Should not happen
				}
				if sortDef.Direction == "desc" {
					return strings.ToLower(valI) > strings.ToLower(valJ)
				}
				return strings.ToLower(valI) < strings.ToLower(valJ)
			})
		}
		start := first
		end := first + pageSize
		if start < len(filteredFunctions) {
			if end > len(filteredFunctions) {
				end = len(filteredFunctions)
			}
			page.Functions = filteredFunctions[start:end]
		}
		page.TotalItems = len(filteredFunctions) // Total items matching filter, not just current page
	}

	return page, nil
}

func (ac *AbisCollection) Reload() {
	abisAddr := base.ZeroAddr
	ac.App.Cancel(abisAddr)

	ac.mutex.Lock()
	defer ac.mutex.Unlock()

	ac.isLoading = false
	ac.isFullyLoaded = false
	ac.expectedFunctions = 0
	ac.expectedEvents = 0
	ac.expectedDownloaded = 0
	ac.expectedKnown = 0
	ac.downloadedAbis = make([]coreTypes.Abi, 0)
	ac.knownAbis = make([]coreTypes.Abi, 0)
	ac.allFunctions = make([]coreTypes.Function, 0)
	ac.allEvents = make([]coreTypes.Function, 0)
	ac.seenEncodings = make(map[string]bool)

	msgs.EmitStatus("ABI data reload requested. Cleared existing data and cancelled ongoing operations.")
}

func (ac *AbisCollection) Delete(address string) error {
	ac.mutex.Lock()
	defer ac.mutex.Unlock()

	for i, abi := range ac.downloadedAbis {
		if abi.Address.Hex() == address {
			ac.downloadedAbis = append(ac.downloadedAbis[:i], ac.downloadedAbis[i+1:]...)
			msgs.EmitStatus(fmt.Sprintf("Deleted downloaded ABI for address: %s", address))
			return nil
		}
	}

	return fmt.Errorf("ABI with address %s not found", address)
}
