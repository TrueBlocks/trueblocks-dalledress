// ADD_ROUTE
package abis

import (
	"fmt"
	"strings"

	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/facets"
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
		page.IsLoading = ac.downloadedFacet.IsLoading()
		page.IsLoaded = ac.downloadedFacet.IsLoaded()

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
				return facets.SortPageSlice(items, sortSpec, sdk.SortAbis)
			}
			return nil
		}

		result, err := ac.downloadedFacet.GetPage(first, pageSize, filterFunc, sortSpec, sortFunc)
		if err != nil {
			return AbisPage{}, err
		}
		page.Abis = result.Items
		page.TotalItems = result.TotalItems
		page.ExpectedTotal = ac.downloadedFacet.ExpectedCount()
		return page, nil

	case AbisKnown:
		page.IsLoading = ac.knownFacet.IsLoading()
		page.IsLoaded = ac.knownFacet.IsLoaded()

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
				return facets.SortPageSlice(items, sortSpec, sdk.SortAbis)
			}
			return nil
		}

		result, err := ac.knownFacet.GetPage(first, pageSize, filterFunc, sortSpec, sortFunc)
		if err != nil {
			return AbisPage{}, err
		}
		page.Abis = result.Items
		page.TotalItems = result.TotalItems
		page.ExpectedTotal = ac.knownFacet.ExpectedCount()
		return page, nil

	case AbisFunctions, AbisEvents:
		facet := ac.functionsFacet
		if listKind == AbisEvents {
			facet = ac.eventsFacet
		}
		page.IsLoading = facet.IsLoading()
		page.IsLoaded = facet.IsLoaded()

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
				return facets.SortPageSlice(items, sortSpec, sdk.SortFunctions)
			}
			return nil
		}

		result, err := facet.GetPage(first, pageSize, filterFunc, sortSpec, sortFunc)
		if err != nil {
			return AbisPage{}, err
		}
		page.Functions = result.Items
		page.TotalItems = result.TotalItems
		page.ExpectedTotal = facet.ExpectedCount()
		return page, nil
	}
	return AbisPage{}, fmt.Errorf("unknown list kind: %s", listKind)
}

// ADD_ROUTE
