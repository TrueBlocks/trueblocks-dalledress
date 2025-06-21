package names

import (
	"fmt"
	"strings"
	"sync"
	"time"

	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/facets"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/logging"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/preferences"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

const (
	NamesAll      types.DataFacet = "all"
	NamesCustom   types.DataFacet = "custom"
	NamesPrefund  types.DataFacet = "prefund"
	NamesRegular  types.DataFacet = "regular"
	NamesBaddress types.DataFacet = "baddress"
)

func init() {
	types.RegisterDataFacet(NamesAll)
	types.RegisterDataFacet(NamesCustom)
	types.RegisterDataFacet(NamesPrefund)
	types.RegisterDataFacet(NamesRegular)
	types.RegisterDataFacet(NamesBaddress)
}

type NamesCollection struct {
	allFacet      *facets.Facet[Name]
	customFacet   *facets.Facet[Name]
	prefundFacet  *facets.Facet[Name]
	regularFacet  *facets.Facet[Name]
	baddressFacet *facets.Facet[Name]
	summary       types.Summary
	summaryMutex  sync.RWMutex
}

func NewNamesCollection() *NamesCollection {
	nc := &NamesCollection{}
	nc.ResetSummary()
	nc.initializeFacets()
	return nc
}

func (nc *NamesCollection) initializeFacets() {
	namesStore := GetNamesStore()

	allFacet := facets.NewFacetWithSummary(
		NamesAll,
		nil,
		nil,
		namesStore,
		"names",
		nc,
	)

	customFacet := facets.NewFacetWithSummary(
		NamesCustom,
		func(name *Name) bool { return name.Parts&coreTypes.Custom != 0 },
		nil,
		namesStore,
		"names",
		nc,
	)

	prefundFacet := facets.NewFacetWithSummary(
		NamesPrefund,
		func(name *Name) bool { return name.Parts&coreTypes.Prefund != 0 },
		nil,
		namesStore,
		"names",
		nc,
	)

	regularFacet := facets.NewFacetWithSummary(
		NamesRegular,
		func(name *Name) bool { return name.Parts&coreTypes.Regular != 0 },
		nil,
		namesStore,
		"names",
		nc,
	)

	baddressFacet := facets.NewFacetWithSummary(
		NamesBaddress,
		func(name *Name) bool { return name.Parts&coreTypes.Baddress != 0 },
		nil,
		namesStore,
		"names",
		nc,
	)

	nc.allFacet = allFacet
	nc.customFacet = customFacet
	nc.prefundFacet = prefundFacet
	nc.regularFacet = regularFacet
	nc.baddressFacet = baddressFacet
}

func (nc *NamesCollection) LoadData(dataFacet types.DataFacet) {
	if !nc.NeedsUpdate(dataFacet) {
		return
	}

	var facet *facets.Facet[Name]
	var facetName types.DataFacet

	switch dataFacet {
	case NamesAll:
		facet = nc.allFacet
		facetName = NamesAll
	case NamesCustom:
		facet = nc.customFacet
		facetName = NamesCustom
	case NamesPrefund:
		facet = nc.prefundFacet
		facetName = NamesPrefund
	case NamesRegular:
		facet = nc.regularFacet
		facetName = NamesRegular
	case NamesBaddress:
		facet = nc.baddressFacet
		facetName = NamesBaddress
	default:
		logging.LogError("LoadData: unexpected dataFacet: %v", fmt.Errorf("invalid dataFacet: %s", dataFacet), nil)
		return
	}

	go func() {
		if err := facet.Load(); err != nil {
			logging.LogError(fmt.Sprintf("LoadData.%s from store: %%v", facetName), err, facets.ErrAlreadyLoading)
		}
	}()
}

func (nc *NamesCollection) Reset(dataFacet types.DataFacet) {
	switch dataFacet {
	case NamesAll, NamesCustom, NamesPrefund, NamesRegular, NamesBaddress:
		namesStore.Reset()
	default:
		return
	}
}

func (nc *NamesCollection) NeedsUpdate(dataFacet types.DataFacet) bool {
	var facet *facets.Facet[Name]

	switch dataFacet {
	case NamesAll:
		facet = nc.allFacet
	case NamesCustom:
		facet = nc.customFacet
	case NamesPrefund:
		facet = nc.prefundFacet
	case NamesRegular:
		facet = nc.regularFacet
	case NamesBaddress:
		facet = nc.baddressFacet
	default:
		return false
	}

	return facet.NeedsUpdate()
}

func (nc *NamesCollection) getExpectedTotal(dataFacet types.DataFacet) int {
	_ = dataFacet
	if count, err := GetNamesCount(); err == nil && count > 0 {
		return count
	}
	return nc.allFacet.ExpectedCount()
}

func (nc *NamesCollection) matchesFilter(name *Name, filter string) bool {
	filterLower := strings.ToLower(filter)

	addressHex := strings.ToLower(name.Address.Hex())
	addressNoPrefix := strings.TrimPrefix(addressHex, "0x")
	addressNoLeadingZeros := strings.TrimLeft(addressNoPrefix, "0")

	if strings.Contains(addressHex, filterLower) ||
		strings.Contains(addressNoPrefix, filterLower) ||
		strings.Contains(addressNoLeadingZeros, filterLower) {
		return true
	}

	if strings.Contains(strings.ToLower(name.Name), filterLower) {
		return true
	}

	if strings.Contains(strings.ToLower(name.Tags), filterLower) {
		return true
	}

	if strings.Contains(strings.ToLower(name.Source), filterLower) {
		return true
	}

	if strings.HasPrefix(filterLower, "0x") {
		fNoPrefix := strings.TrimPrefix(filterLower, "0x")
		if strings.Contains(addressNoPrefix, fNoPrefix) || strings.Contains(addressNoLeadingZeros, fNoPrefix) {
			return true
		}
	}

	return false
}

func GetNamesCount() (int, error) {
	chainName := preferences.GetChain()
	countOpts := sdk.NamesOptions{
		Globals: sdk.Globals{Cache: true, Chain: chainName},
	}
	if countResult, _, err := countOpts.NamesCount(); err != nil {
		return 0, fmt.Errorf("NamesCount query error: %v", err)
	} else if len(countResult) > 0 {
		return int(countResult[0].Count), nil
	}
	return 0, nil
}

func (nc *NamesCollection) GetSupportedFacets() []types.DataFacet {
	return []types.DataFacet{
		NamesAll,
		NamesCustom,
		NamesPrefund,
		NamesRegular,
		NamesBaddress,
	}
}

func (nc *NamesCollection) GetStoreForFacet(dataFacet types.DataFacet) string {
	switch dataFacet {
	case NamesAll, NamesCustom, NamesPrefund, NamesRegular, NamesBaddress:
		return "names"
	default:
		return ""
	}
}

func (nc *NamesCollection) GetCollectionName() string {
	return "names"
}

func (nc *NamesCollection) GetSummary() types.Summary {
	nc.summaryMutex.RLock()
	defer nc.summaryMutex.RUnlock()
	return nc.summary
}

func (nc *NamesCollection) AccumulateItem(item interface{}, summary *types.Summary) {
	if name, ok := item.(*Name); ok {
		nc.summaryMutex.Lock()
		defer nc.summaryMutex.Unlock()

		nc.summary.TotalCount++
		if nc.summary.FacetCounts == nil {
			nc.summary.FacetCounts = make(map[types.DataFacet]int)
		}
		if name.Parts&coreTypes.Custom != 0 {
			nc.summary.FacetCounts[NamesCustom]++
		}
		if name.Parts&coreTypes.Prefund != 0 {
			nc.summary.FacetCounts[NamesPrefund]++
		}
		if name.Parts&coreTypes.Regular != 0 {
			nc.summary.FacetCounts[NamesRegular]++
		}
		if name.Parts&coreTypes.Baddress != 0 {
			nc.summary.FacetCounts[NamesBaddress]++
		}
		nc.summary.LastUpdated = time.Now().Unix()
	}
}

func (nc *NamesCollection) ResetSummary() {
	nc.summaryMutex.Lock()
	defer nc.summaryMutex.Unlock()
	nc.summary = types.Summary{
		TotalCount:  0,
		FacetCounts: make(map[types.DataFacet]int),
		CustomData:  make(map[string]interface{}),
		LastUpdated: time.Now().Unix(),
	}
}
