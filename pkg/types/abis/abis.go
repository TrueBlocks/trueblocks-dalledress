package abis

import (
	"fmt"
	"sync"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/facets"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/logging"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
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

type AbisCollection struct {
	downloadedFacet *facets.Facet[Abi]
	knownFacet      *facets.Facet[Abi]
	functionsFacet  *facets.Facet[Function]
	eventsFacet     *facets.Facet[Function]
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
		func(fn *Function) bool { return fn.FunctionType != "event" },
		isDupEncoding(),
		abisDetailStore,
		"abis",
		ac,
	)

	eventsFacet := facets.NewFacetWithSummary(
		AbisEvents,
		func(fn *Function) bool { return fn.FunctionType == "event" },
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

func isDownload(abi *Abi) bool {
	return !abi.IsKnown
}

func isKnown(abi *Abi) bool {
	return abi.IsKnown
}

func isDupEncoding() func(existing []*Function, newItem *Function) bool {
	seen := make(map[string]bool)
	lastExistingLen := 0

	return func(existing []*Function, newItem *Function) bool {
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

	var facetAbi *facets.Facet[Abi]
	var facetFunction *facets.Facet[Function]
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
	var facetAbi *facets.Facet[Abi]
	var facetFunction *facets.Facet[Function]

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
	case *Abi:
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

	case *Function:
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

func (ac *AbisCollection) GetSummary() types.Summary {
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
