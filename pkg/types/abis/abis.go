package abis

import (
	"fmt"
	"strings"
	"sync"

	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/facets"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/logging"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

const (
	AbisDownloaded types.DataFacet = "downloaded"
	AbisKnown      types.DataFacet = "known"
	AbisFunctions  types.DataFacet = "functions"
	AbisEvents     types.DataFacet = "events"
)

func init() {
	types.RegisterDataFacet(AbisDownloaded)
	types.RegisterDataFacet(AbisKnown)
	types.RegisterDataFacet(AbisFunctions)
	types.RegisterDataFacet(AbisEvents)
}

type AbisPage struct {
	Facet         types.DataFacet      `json:"facet"`
	Abis          []coreTypes.Abi      `json:"abis,omitempty"`
	Functions     []coreTypes.Function `json:"functions,omitempty"`
	TotalItems    int                  `json:"totalItems"`
	ExpectedTotal int                  `json:"expectedTotal"`
	IsFetching    bool                 `json:"isFetching"`
	State         types.LoadState      `json:"state"`
}

func (ap *AbisPage) GetFacet() types.DataFacet {
	return ap.Facet
}

func (ap *AbisPage) GetTotalItems() int {
	return ap.TotalItems
}

func (ap *AbisPage) GetExpectedTotal() int {
	return ap.ExpectedTotal
}

func (ap *AbisPage) GetIsFetching() bool {
	return ap.IsFetching
}

func (ap *AbisPage) GetState() types.LoadState {
	return ap.State
}

type AbisCollection struct {
	downloadedFacet *facets.Facet[coreTypes.Abi]
	knownFacet      *facets.Facet[coreTypes.Abi]
	functionsFacet  *facets.Facet[coreTypes.Function]
	eventsFacet     *facets.Facet[coreTypes.Function]
	summary         types.Summary
	summaryMutex    sync.RWMutex
}

func NewAbisCollection() *AbisCollection {
	ac := &AbisCollection{
		summary: types.Summary{
			TotalCount:  0,
			FacetCounts: make(map[types.DataFacet]int),
			CustomData:  make(map[string]interface{}),
		},
	}

	ac.initializeFacets()
	return ac
}

func (ac *AbisCollection) initializeFacets() {
	abisListStore := GetAbisListStore()
	abisDetailStore := GetAbisDetailStore()

	downloadedFacet := facets.NewFacetWithSummary(
		AbisDownloaded,
		isDownload,
		nil,
		abisListStore,
		"abis",
		ac,
	)

	knownFacet := facets.NewFacetWithSummary(
		AbisKnown,
		isKnown,
		nil,
		abisListStore,
		"abis",
		ac,
	)

	functionsFacet := facets.NewFacetWithSummary(
		AbisFunctions,
		func(fn *coreTypes.Function) bool { return fn.FunctionType != "event" },
		isDupEncoding(),
		abisDetailStore,
		"abis",
		ac,
	)

	eventsFacet := facets.NewFacetWithSummary(
		AbisEvents,
		func(fn *coreTypes.Function) bool { return fn.FunctionType == "event" },
		isDupEncoding(),
		abisDetailStore,
		"abis",
		ac,
	)

	ac.downloadedFacet = downloadedFacet
	ac.knownFacet = knownFacet
	ac.functionsFacet = functionsFacet
	ac.eventsFacet = eventsFacet
}

func isDownload(abi *coreTypes.Abi) bool {
	return !abi.IsKnown
}

func isKnown(abi *coreTypes.Abi) bool {
	return abi.IsKnown
}

func isDupEncoding() func(existing []*coreTypes.Function, newItem *coreTypes.Function) bool {
	seen := make(map[string]bool)
	lastExistingLen := 0

	return func(existing []*coreTypes.Function, newItem *coreTypes.Function) bool {
		if newItem == nil {
			return false
		}

		if len(existing) == 0 && lastExistingLen > 0 {
			seen = make(map[string]bool)
		}
		lastExistingLen = len(existing)

		if seen[newItem.Encoding] {
			return true
		}
		seen[newItem.Encoding] = true
		return false
	}
}

func (ac *AbisCollection) LoadData(dataFacet types.DataFacet) {
	if !ac.NeedsUpdate(dataFacet) {
		return
	}

	var facetAbi *facets.Facet[coreTypes.Abi]
	var facetFunction *facets.Facet[coreTypes.Function]
	var facetName string

	switch dataFacet {
	case AbisDownloaded:
		facetAbi = ac.downloadedFacet
		facetName = "downloaded"
	case AbisKnown:
		facetAbi = ac.knownFacet
		facetName = "known"
	case AbisFunctions:
		facetFunction = ac.functionsFacet
		facetName = "functions"
	case AbisEvents:
		facetFunction = ac.eventsFacet
		facetName = "events"
	default:
		logging.LogError("LoadData: unexpected dataFacet: %v", fmt.Errorf("invalid dataFacet: %s", dataFacet), nil)
		return
	}

	go func() {
		if facetAbi != nil {
			if err := facetAbi.Load(); err != nil {
				logging.LogError(fmt.Sprintf("LoadData.%s from store: %%v", facetName), err, facets.ErrAlreadyLoading)
			}
		} else if facetFunction != nil {
			if err := facetFunction.Load(); err != nil {
				logging.LogError(fmt.Sprintf("LoadData.%s from store: %%v", facetName), err, facets.ErrAlreadyLoading)
			}
		}
	}()
}

func (ac *AbisCollection) Reset(dataFacet types.DataFacet) {
	switch dataFacet {
	case AbisDownloaded, AbisKnown:
		abisListStore.Reset()
	case AbisFunctions, AbisEvents:
		abisDetailStore.Reset()
	default:
		return
	}
}

func (ac *AbisCollection) NeedsUpdate(dataFacet types.DataFacet) bool {
	var facetAbi *facets.Facet[coreTypes.Abi]
	var facetFunction *facets.Facet[coreTypes.Function]

	switch dataFacet {
	case AbisDownloaded:
		facetAbi = ac.downloadedFacet
	case AbisKnown:
		facetAbi = ac.knownFacet
	case AbisFunctions:
		facetFunction = ac.functionsFacet
	case AbisEvents:
		facetFunction = ac.eventsFacet
	default:
		return false
	}

	if facetAbi != nil {
		return facetAbi.NeedsUpdate()
	} else if facetFunction != nil {
		return facetFunction.NeedsUpdate()
	}

	return false
}

func (ac *AbisCollection) GetPage(
	dataFacet types.DataFacet,
	first, pageSize int,
	sortSpec sdk.SortSpec,
	filter string,
) (types.Page, error) {
	page := &AbisPage{
		Facet: dataFacet,
	}
	filter = strings.ToLower(filter)

	var listFacet *facets.Facet[coreTypes.Abi]
	var detailFacet *facets.Facet[coreTypes.Function]

	switch dataFacet {
	case AbisDownloaded:
		listFacet = ac.downloadedFacet
	case AbisKnown:
		listFacet = ac.knownFacet
	case AbisFunctions:
		detailFacet = ac.functionsFacet
	case AbisEvents:
		detailFacet = ac.eventsFacet
	default:
		// This is truly a validation error - invalid DataFacet for this collection
		return nil, types.NewValidationError("abis", dataFacet, "GetPage",
			fmt.Errorf("unsupported dataFacet: %v", dataFacet))
	}

	if listFacet != nil {
		var listFilterFunc = func(item *coreTypes.Abi) bool {
			return strings.Contains(strings.ToLower(item.Name), filter)
		}
		var listSortFunc = func(items []coreTypes.Abi, sort sdk.SortSpec) error {
			return sdk.SortAbis(items, sort)
		}
		if result, err := listFacet.GetPage(first, pageSize, listFilterFunc, sortSpec, listSortFunc); err != nil {
			// This is likely an SDK or store error, not a validation error
			return nil, types.NewStoreError("abis", dataFacet, "GetPage", err)
		} else {
			page.Abis, page.TotalItems, page.State = result.Items, result.TotalItems, result.State
		}
		page.IsFetching = listFacet.IsFetching()
		page.ExpectedTotal = listFacet.ExpectedCount()

	} else if detailFacet != nil {
		var detailFilter = func(item *coreTypes.Function) bool {
			return strings.Contains(strings.ToLower(item.Name), filter) ||
				strings.Contains(strings.ToLower(item.Encoding), filter)
		}
		var detailSortFunc = func(items []coreTypes.Function, sort sdk.SortSpec) error {
			return sdk.SortFunctions(items, sort)
		}
		if result, err := detailFacet.GetPage(first, pageSize, detailFilter, sortSpec, detailSortFunc); err != nil {
			// This is likely an SDK or store error, not a validation error
			return nil, types.NewStoreError("abis", dataFacet, "GetPage", err)
		} else {
			page.Functions, page.TotalItems, page.State = result.Items, result.TotalItems, result.State
		}
		page.IsFetching = detailFacet.IsFetching()
		page.ExpectedTotal = detailFacet.ExpectedCount()

	} else {
		// This should not happen since we validated dataFacet above
		return nil, types.NewValidationError("abis", dataFacet, "GetPage",
			fmt.Errorf("no facet found for dataFacet: %v", dataFacet))
	}

	return page, nil
}

func (ac *AbisCollection) GetAbisPage(
	dataFacet types.DataFacet,
	first, pageSize int,
	sortSpec sdk.SortSpec,
	filter string,
) (*AbisPage, error) {
	page, err := ac.GetPage(dataFacet, first, pageSize, sortSpec, filter)
	if err != nil {
		return nil, err
	}

	abisPage, ok := page.(*AbisPage)
	if !ok {
		return nil, fmt.Errorf("internal error: GetPage returned unexpected type %T", page)
	}

	return abisPage, nil
}

func (ac *AbisCollection) GetSupportedFacets() []types.DataFacet {
	return []types.DataFacet{
		AbisDownloaded,
		AbisKnown,
		AbisFunctions,
		AbisEvents,
	}
}

func (ac *AbisCollection) GetStoreForFacet(dataFacet types.DataFacet) string {
	switch dataFacet {
	case AbisDownloaded, AbisKnown:
		return "abis-list"
	case AbisFunctions, AbisEvents:
		return "abis-detail"
	default:
		return ""
	}
}

func (ac *AbisCollection) GetCollectionName() string {
	return "abis"
}

func (ac *AbisCollection) AccumulateItem(item interface{}, summary *types.Summary) {
	ac.summaryMutex.Lock()
	defer ac.summaryMutex.Unlock()

	if summary.FacetCounts == nil {
		summary.FacetCounts = make(map[types.DataFacet]int)
	}

	switch v := item.(type) {
	case *coreTypes.Abi:
		summary.TotalCount++

		if v.IsKnown {
			summary.FacetCounts[AbisKnown]++
		} else {
			summary.FacetCounts[AbisDownloaded]++
		}

		if summary.CustomData == nil {
			summary.CustomData = make(map[string]interface{})
		}

		knownCount, _ := summary.CustomData["knownCount"].(int)
		downloadedCount, _ := summary.CustomData["downloadedCount"].(int)

		if v.IsKnown {
			knownCount++
		} else {
			downloadedCount++
		}

		summary.CustomData["knownCount"] = knownCount
		summary.CustomData["downloadedCount"] = downloadedCount

	case *coreTypes.Function:
		summary.TotalCount++

		if v.FunctionType == "event" {
			summary.FacetCounts[AbisEvents]++
		} else {
			summary.FacetCounts[AbisFunctions]++
		}

		if summary.CustomData == nil {
			summary.CustomData = make(map[string]interface{})
		}

		functionsCount, _ := summary.CustomData["functionsCount"].(int)
		eventsCount, _ := summary.CustomData["eventsCount"].(int)

		if v.FunctionType == "event" {
			eventsCount++
		} else {
			functionsCount++
		}

		summary.CustomData["functionsCount"] = functionsCount
		summary.CustomData["eventsCount"] = eventsCount
	}
}

func (ac *AbisCollection) GetCurrentSummary() types.Summary {
	ac.summaryMutex.RLock()
	defer ac.summaryMutex.RUnlock()

	summary := ac.summary
	summary.FacetCounts = make(map[types.DataFacet]int)
	for k, v := range ac.summary.FacetCounts {
		summary.FacetCounts[k] = v
	}

	if ac.summary.CustomData != nil {
		summary.CustomData = make(map[string]interface{})
		for k, v := range ac.summary.CustomData {
			summary.CustomData[k] = v
		}
	}

	return summary
}

func (ac *AbisCollection) ResetSummary() {
	ac.summaryMutex.Lock()
	defer ac.summaryMutex.Unlock()
	ac.summary = types.Summary{
		TotalCount:  0,
		FacetCounts: make(map[types.DataFacet]int),
		CustomData:  make(map[string]interface{}),
		LastUpdated: 0,
	}
}

func (ac *AbisCollection) GetSummary() types.Summary {
	return ac.GetCurrentSummary()
}
