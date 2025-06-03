// ADD_ROUTE
package abis

import (
	"fmt"
	"strings"

	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/sorting"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/streaming"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

// ------------------------------------------------------------------------------
type AbisPage struct {
	Kind          types.ListKind       `json:"kind"`
	Abis          []coreTypes.Abi      `json:"abis,omitempty"`
	Functions     []coreTypes.Function `json:"functions,omitempty"`
	TotalItems    int                  `json:"totalItems"`
	ExpectedTotal int                  `json:"expectedTotal"`
	IsLoading     bool                 `json:"isLoading"`
	IsLoaded      bool                 `json:"isLoaded"`
}

// ------------------------------------------------------------------------------
func (ac *AbisCollection) GetPage(listKind types.ListKind, first, pageSize int, sortDef *sorting.SortDef, filter string) (AbisPage, error) {
	ac.LoadData(listKind)

	filter = strings.ToLower(filter)
	sortSpec := sorting.ConvertToSortSpec(sortDef)
	var page = AbisPage{
		Kind: listKind,
	}

	switch listKind {
	case AbisDownloaded:
		// Use Repository pattern for downloaded ABIs
		page.IsLoading = ac.downloadedRepo.IsLoading()
		page.IsLoaded = ac.downloadedRepo.IsLoaded()

		// Create filter function from the string filter
		filterFunc := func(abi *coreTypes.Abi) bool {
			if filter == "" {
				return true
			}
			searchFields := []string{abi.Name, abi.Address.Hex()}
			for _, field := range searchFields {
				if strings.Contains(strings.ToLower(field), filter) {
					return true
				}
			}
			return false
		}

		// Create wrapper sort function to match Repository interface
		sortFunc := func(items []coreTypes.Abi, sortSpecInterface interface{}) error {
			if sortSpec, ok := sortSpecInterface.(sdk.SortSpec); ok {
				return streaming.SortPageSlice(items, sortSpec, sdk.SortAbis)
			}
			return nil // No sorting if type doesn't match
		}

		// Get paginated data from repository with sorting
		result, err := ac.downloadedRepo.GetPage(first, pageSize, filterFunc, sortSpec, sortFunc)
		if err != nil {
			return AbisPage{}, err
		}

		page.Abis = result.Items
		page.TotalItems = result.TotalItems
		page.ExpectedTotal = ac.downloadedRepo.ExpectedCount()
		return page, nil

	case AbisKnown:
		// Use Repository pattern for known ABIs
		page.IsLoading = ac.knownRepo.IsLoading()
		page.IsLoaded = ac.knownRepo.IsLoaded()

		// Create filter function from the string filter
		filterFunc := func(abi *coreTypes.Abi) bool {
			if filter == "" {
				return true
			}
			searchFields := []string{abi.Name, abi.Address.Hex()}
			for _, field := range searchFields {
				if strings.Contains(strings.ToLower(field), filter) {
					return true
				}
			}
			return false
		}

		// Create wrapper sort function to match Repository interface
		sortFunc := func(items []coreTypes.Abi, sortSpecInterface interface{}) error {
			if sortSpec, ok := sortSpecInterface.(sdk.SortSpec); ok {
				return streaming.SortPageSlice(items, sortSpec, sdk.SortAbis)
			}
			return nil // No sorting if type doesn't match
		}

		// Get paginated data from repository with sorting
		result, err := ac.knownRepo.GetPage(first, pageSize, filterFunc, sortSpec, sortFunc)
		if err != nil {
			return AbisPage{}, err
		}

		page.Abis = result.Items
		page.TotalItems = result.TotalItems
		page.ExpectedTotal = ac.knownRepo.ExpectedCount()
		return page, nil

	case AbisFunctions:
		// Use Repository pattern for functions
		page.IsLoading = ac.functionsRepo.IsLoading()
		page.IsLoaded = ac.functionsRepo.IsLoaded()

		// Create filter function from the string filter
		filterFunc := func(fn *coreTypes.Function) bool {
			if filter == "" {
				return true
			}
			searchFields := []string{fn.Name, fn.Signature, fn.Encoding}
			for _, field := range searchFields {
				if strings.Contains(strings.ToLower(field), filter) {
					return true
				}
			}
			return false
		}

		// Create wrapper sort function to match Repository interface
		sortFunc := func(items []coreTypes.Function, sortSpecInterface interface{}) error {
			if sortSpec, ok := sortSpecInterface.(sdk.SortSpec); ok {
				return streaming.SortPageSlice(items, sortSpec, sdk.SortFunctions)
			}
			return nil // No sorting if type doesn't match
		}

		// Get paginated data from repository with sorting
		result, err := ac.functionsRepo.GetPage(first, pageSize, filterFunc, sortSpec, sortFunc)
		if err != nil {
			return AbisPage{}, err
		}

		page.Functions = result.Items
		page.TotalItems = result.TotalItems
		page.ExpectedTotal = ac.functionsRepo.ExpectedCount()
		return page, nil

	case AbisEvents:
		// Use Repository pattern for events
		page.IsLoading = ac.eventsRepo.IsLoading()
		page.IsLoaded = ac.eventsRepo.IsLoaded()

		// Create filter function from the string filter
		filterFunc := func(fn *coreTypes.Function) bool {
			if filter == "" {
				return true
			}
			searchFields := []string{fn.Name, fn.Signature, fn.Encoding}
			for _, field := range searchFields {
				if strings.Contains(strings.ToLower(field), filter) {
					return true
				}
			}
			return false
		}

		// Create wrapper sort function to match Repository interface
		sortFunc := func(items []coreTypes.Function, sortSpecInterface interface{}) error {
			if sortSpec, ok := sortSpecInterface.(sdk.SortSpec); ok {
				return streaming.SortPageSlice(items, sortSpec, sdk.SortFunctions)
			}
			return nil // No sorting if type doesn't match
		}

		// Get paginated data from repository with sorting
		result, err := ac.eventsRepo.GetPage(first, pageSize, filterFunc, sortSpec, sortFunc)
		if err != nil {
			return AbisPage{}, err
		}

		page.Functions = result.Items
		page.TotalItems = result.TotalItems
		page.ExpectedTotal = ac.eventsRepo.ExpectedCount()
		return page, nil

	default:
		return AbisPage{}, fmt.Errorf("invalid list kind: %s", listKind)
	}
}

// ADD_ROUTE
