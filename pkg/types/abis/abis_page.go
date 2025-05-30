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

	ac.mutex.RLock()
	defer ac.mutex.RUnlock()

	filter = strings.ToLower(filter)
	sortSpec := sorting.ConvertToSortSpec(sortDef)
	var err error
	var page = AbisPage{
		Kind:      listKind,
		IsLoading: ac.isLoading == 1,
		IsLoaded:  ac.isKnownLoaded,
	}

	switch listKind {
	case AbisDownloaded:
		if page.Abis, page.TotalItems, page.ExpectedTotal, err = streaming.ProcessPage(
			"downloaded abis",
			&ac.downloadedAbis,
			sortSpec,
			getAbiSortFunc,
			streaming.CreatePageFilter(filter, getAbiSearchFields),
			first, pageSize,
		); err != nil {
			return AbisPage{}, err
		} else {
			return page, nil
		}

	case AbisKnown:
		if page.Abis, page.TotalItems, page.ExpectedTotal, err = streaming.ProcessPage(
			"known abis",
			&ac.knownAbis,
			sortSpec,
			getAbiSortFunc,
			streaming.CreatePageFilter(filter, getAbiSearchFields),
			first, pageSize,
		); err != nil {
			return AbisPage{}, err
		} else {
			return page, nil
		}

	case AbisFunctions:
		if page.Functions, page.TotalItems, page.ExpectedTotal, err = streaming.ProcessPage(
			"functions",
			&ac.allFunctions,
			sortSpec,
			getFunctionSortFunc,
			streaming.CreatePageFilter(filter, getFunctionSearchFields),
			first, pageSize,
		); err != nil {
			return AbisPage{}, err
		} else {
			return page, nil
		}

	case AbisEvents:
		if page.Functions, page.TotalItems, page.ExpectedTotal, err = streaming.ProcessPage(
			"events",
			&ac.allEvents,
			sortSpec,
			getFunctionSortFunc,
			streaming.CreatePageFilter(filter, getFunctionSearchFields),
			first, pageSize,
		); err != nil {
			return AbisPage{}, err
		} else {
			return page, nil
		}

	default:
		return AbisPage{}, fmt.Errorf("invalid list kind: %s", listKind)
	}
}

// ------------------------------------------------------------------------------
func getAbiSearchFields(item coreTypes.Abi) []string {
	return []string{item.Name, item.Address.Hex()}
}

// ------------------------------------------------------------------------------
func getAbiSortFunc(items []coreTypes.Abi, sortSpec interface{}) error {
	return streaming.SortPageSlice(items, sortSpec, sdk.SortAbis)
}

// ------------------------------------------------------------------------------
func getFunctionSearchFields(item coreTypes.Function) []string {
	return []string{item.Name, item.Signature, item.Encoding}
}

// ------------------------------------------------------------------------------
func getFunctionSortFunc(items []coreTypes.Function, sortSpec interface{}) error {
	return streaming.SortPageSlice(items, sortSpec, sdk.SortFunctions)
}

// ADD_ROUTE
