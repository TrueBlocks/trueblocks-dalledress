// ADD_ROUTE
package abis

import (
	"fmt"
	"strings"

	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/streaming"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

// AbisPage represents a paginated view of ABIs or functions
// for a given ListKind (Downloaded, Known, Functions, Events)
type AbisPage struct {
	Kind          types.ListKind       `json:"kind"`
	Abis          []coreTypes.Abi      `json:"abis,omitempty"`
	Functions     []coreTypes.Function `json:"functions,omitempty"`
	TotalItems    int                  `json:"totalItems"`
	ExpectedTotal int                  `json:"expectedTotal"`
	IsLoading     bool                 `json:"isLoading"`
	IsLoaded      bool                 `json:"isLoaded"`
}

// GetPage returns a filtered, sorted, and paginated AbisPage for the given ListKind
func (ac *AbisCollection) GetPage(
	listKind types.ListKind,
	first, pageSize int,
	sortSpec sdk.SortSpec,
	filter string,
) (AbisPage, error) {
	ac.LoadData(listKind)

	filter = strings.ToLower(filter)
	var page = AbisPage{Kind: listKind}

	switch listKind {
	case AbisDownloaded:
		// Downloaded ABIs
		page.IsLoading = ac.downloadedRepo.IsLoading()
		page.IsLoaded = ac.downloadedRepo.IsLoaded()

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

		sortFunc := func(items []coreTypes.Abi, sortSpecInterface interface{}) error {
			if sortSpec, ok := sortSpecInterface.(sdk.SortSpec); ok {
				return streaming.SortPageSlice(items, sortSpec, sdk.SortAbis)
			}
			return nil
		}

		result, err := ac.downloadedRepo.GetPage(first, pageSize, filterFunc, sortSpec, sortFunc)
		if err != nil {
			return AbisPage{}, err
		}
		page.Abis = result.Items
		page.TotalItems = result.TotalItems
		page.ExpectedTotal = ac.downloadedRepo.ExpectedCount()
		return page, nil

	case AbisKnown:
		// Known ABIs
		page.IsLoading = ac.knownRepo.IsLoading()
		page.IsLoaded = ac.knownRepo.IsLoaded()

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

		sortFunc := func(items []coreTypes.Abi, sortSpecInterface interface{}) error {
			if sortSpec, ok := sortSpecInterface.(sdk.SortSpec); ok {
				return streaming.SortPageSlice(items, sortSpec, sdk.SortAbis)
			}
			return nil
		}

		result, err := ac.knownRepo.GetPage(first, pageSize, filterFunc, sortSpec, sortFunc)
		if err != nil {
			return AbisPage{}, err
		}
		page.Abis = result.Items
		page.TotalItems = result.TotalItems
		page.ExpectedTotal = ac.knownRepo.ExpectedCount()
		return page, nil

	case AbisFunctions, AbisEvents:
		// Functions or Events
		repo := ac.functionsRepo
		if listKind == AbisEvents {
			repo = ac.eventsRepo
		}
		page.IsLoading = repo.IsLoading()
		page.IsLoaded = repo.IsLoaded()

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

		sortFunc := func(items []coreTypes.Function, sortSpecInterface interface{}) error {
			if sortSpec, ok := sortSpecInterface.(sdk.SortSpec); ok {
				return streaming.SortPageSlice(items, sortSpec, sdk.SortFunctions)
			}
			return nil
		}

		result, err := repo.GetPage(first, pageSize, filterFunc, sortSpec, sortFunc)
		if err != nil {
			return AbisPage{}, err
		}
		page.Functions = result.Items
		page.TotalItems = result.TotalItems
		page.ExpectedTotal = repo.ExpectedCount()
		return page, nil
	}
	return AbisPage{}, fmt.Errorf("unknown list kind: %s", listKind)
}

// ADD_ROUTE
