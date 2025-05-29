// ADD_ROUTE
package types

import (
	"fmt"
	"sort"
	"strings"

	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/sorting"
)

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
// ADD_ROUTE
