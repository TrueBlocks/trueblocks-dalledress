// ADD_ROUTE
package abis

import (
	"fmt"
	"strings"

	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/sorting"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

type AbisPage struct {
	Kind          types.ListKind       `json:"kind"`
	Abis          []coreTypes.Abi      `json:"abis,omitempty"`
	Functions     []coreTypes.Function `json:"functions,omitempty"`
	TotalItems    int                  `json:"totalItems"`
	ExpectedTotal int                  `json:"expectedTotal"`
	IsLoading     bool                 `json:"isLoading"`
	IsLoaded      bool                 `json:"isLoaded"`
}

func (ac *AbisCollection) GetPage(listKind types.ListKind, first, pageSize int, sortDef *sorting.SortDef, filter string) (AbisPage, error) {
	ac.LoadData(listKind)

	ac.mutex.RLock()
	defer ac.mutex.RUnlock()

	// Get direct reference to slice - no copying!
	var sourceAbis []coreTypes.Abi
	var sourceFunctions []coreTypes.Function
	var listLen int

	switch listKind {
	case AbisDownloaded:
		sourceAbis = ac.downloadedAbis
		listLen = len(ac.downloadedAbis)
	case AbisKnown:
		sourceAbis = ac.knownAbis
		listLen = len(ac.knownAbis)
	case AbisFunctions:
		sourceFunctions = ac.allFunctions
		listLen = len(ac.allFunctions)
	case AbisEvents:
		sourceFunctions = ac.allEvents
		listLen = len(ac.allEvents)
	default:
		return AbisPage{}, fmt.Errorf("unknown ABI page kind: %s", listKind)
	}

	var isLoaded bool
	switch listKind {
	case AbisDownloaded:
		isLoaded = ac.isDownloadedLoaded
	case AbisKnown:
		isLoaded = ac.isKnownLoaded
	case AbisFunctions:
		isLoaded = ac.isFuncsLoaded
	case AbisEvents:
		isLoaded = ac.isEventsLoaded
	}

	page := AbisPage{
		Kind:          listKind,
		IsLoading:     ac.isLoading == 1,
		IsLoaded:      isLoaded,
		ExpectedTotal: listLen,
	}
	filter = strings.ToLower(filter)
	sortSpec := sorting.ConvertToSortSpec(sortDef)

	switch listKind {
	case AbisDownloaded, AbisKnown:
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

	case AbisFunctions, AbisEvents:
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
