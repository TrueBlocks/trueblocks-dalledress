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
	var err error
	var page = AbisPage{
		Kind:      listKind,
		IsLoading: ac.isLoading == 1,
	}

	switch listKind {
	case AbisDownloaded:
		ac.downloadedMutex.RLock()
		defer ac.downloadedMutex.RUnlock()
		page.IsLoaded = ac.isDownloadedLoaded
		if page.Abis, page.TotalItems, page.ExpectedTotal, err = streaming.ProcessPage(
			"downloaded abis",
			&ac.downloaded,
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
		ac.knownMutex.RLock()
		defer ac.knownMutex.RUnlock()
		page.IsLoaded = ac.isKnownLoaded
		if page.Abis, page.TotalItems, page.ExpectedTotal, err = streaming.ProcessPage(
			"known abis",
			&ac.known,
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
		ac.functionsMutex.RLock()
		defer ac.functionsMutex.RUnlock()
		page.IsLoaded = ac.isFuncsLoaded
		if page.Functions, page.TotalItems, page.ExpectedTotal, err = streaming.ProcessPage(
			"functions",
			&ac.functions,
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
		ac.eventsMutex.RLock()
		defer ac.eventsMutex.RUnlock()
		page.IsLoaded = ac.isEventsLoaded
		if page.Functions, page.TotalItems, page.ExpectedTotal, err = streaming.ProcessPage(
			"events",
			&ac.events,
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
func getAbiSortFunc(items []coreTypes.Abi, sortSpec sdk.SortSpec) error {
	return streaming.SortPageSlice(items, sortSpec, sdk.SortAbis)
}

// ------------------------------------------------------------------------------
func getFunctionSearchFields(item coreTypes.Function) []string {
	return []string{item.Name, item.Signature, item.Encoding}
}

// ------------------------------------------------------------------------------
func getFunctionSortFunc(items []coreTypes.Function, sortSpec sdk.SortSpec) error {
	return streaming.SortPageSlice(items, sortSpec, sdk.SortFunctions)
}

// ADD_ROUTE
