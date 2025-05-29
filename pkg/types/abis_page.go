// ADD_ROUTE
package types

import (
	"fmt"
	"strings"

	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/sorting"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

type AbisPage struct {
	Kind          string               `json:"Type"`
	Abis          []coreTypes.Abi      `json:"Abis,omitempty"`
	Functions     []coreTypes.Function `json:"Functions,omitempty"`
	TotalItems    int                  `json:"TotalItems"`
	ExpectedTotal int                  `json:"ExpectedTotal"`
	IsLoading     bool                 `json:"IsLoading"`
	IsLoaded      bool                 `json:"IsLoaded"`
}

func (ac *AbisCollection) GetPage(kind string, first, pageSize int, sortDef *sorting.SortDef, filter string) (AbisPage, error) {
	ac.EnsureInitialLoad()

	ac.mutex.RLock()
	defer ac.mutex.RUnlock()

	// Get direct reference to slice - no copying!
	var sourceAbis []coreTypes.Abi
	var sourceFunctions []coreTypes.Function
	var expectedTotal int

	switch kind {
	case "Downloaded":
		sourceAbis = ac.downloadedAbis
		expectedTotal = ac.expectedDownloaded
	case "Known":
		sourceAbis = ac.knownAbis
		expectedTotal = ac.expectedKnown
	case "Functions":
		sourceFunctions = ac.allFunctions
		expectedTotal = ac.expectedFunctions
	case "Events":
		sourceFunctions = ac.allEvents
		expectedTotal = ac.expectedEvents
	default:
		return AbisPage{}, fmt.Errorf("unknown ABI page kind: %s", kind)
	}

	page := AbisPage{
		Kind:          kind,
		IsLoading:     ac.isLoading,
		IsLoaded:      ac.isLoaded,
		ExpectedTotal: expectedTotal,
	}
	filter = strings.ToLower(filter)
	sortSpec := sorting.ConvertToSortSpec(sortDef)

	switch kind {
	case "Downloaded", "Known":
		filteredAbis := make([]coreTypes.Abi, 0)
		for _, item := range sourceAbis {
			if filter == "" || strings.Contains(strings.ToLower(item.Name), filter) || strings.Contains(strings.ToLower(item.Address.Hex()), filter) {
				filteredAbis = append(filteredAbis, item)
			}
		}
		if len(sortSpec.Fields) > 0 {
			if err := sdk.SortAbis(filteredAbis, sortSpec); err != nil {
				return AbisPage{}, fmt.Errorf("failed to sort ABIs: %w", err)
			}
		}
		start := first
		end := first + pageSize
		if start < len(filteredAbis) {
			if end > len(filteredAbis) {
				end = len(filteredAbis)
			}
			page.Abis = filteredAbis[start:end]
		}
		page.TotalItems = len(filteredAbis)

	case "Functions", "Events":
		filteredFunctions := make([]coreTypes.Function, 0)
		for _, item := range sourceFunctions {
			name := strings.ToLower(item.Name)
			signature := strings.ToLower(item.Signature)
			encoding := strings.ToLower(item.Encoding)
			if filter == "" || strings.Contains(name, filter) || strings.Contains(signature, filter) || strings.Contains(encoding, filter) {
				filteredFunctions = append(filteredFunctions, item)
			}
		}
		if len(sortSpec.Fields) > 0 {
			if err := sdk.SortFunctions(filteredFunctions, sortSpec); err != nil {
				return AbisPage{}, fmt.Errorf("failed to sort Functions: %w", err)
			}
		}
		start := first
		end := first + pageSize
		if start < len(filteredFunctions) {
			if end > len(filteredFunctions) {
				end = len(filteredFunctions)
			}
			page.Functions = filteredFunctions[start:end]
		}
		page.TotalItems = len(filteredFunctions)
	}

	return page, nil
}

// ADD_ROUTE
