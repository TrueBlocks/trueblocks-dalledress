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
	NamesAll      types.ListKind = "All"
	NamesCustom   types.ListKind = "Custom"
	NamesPrefund  types.ListKind = "Prefund"
	NamesRegular  types.ListKind = "Regular"
	NamesBaddress types.ListKind = "Baddress"
)

func init() {
	types.RegisterKind(NamesAll)
	types.RegisterKind(NamesCustom)
	types.RegisterKind(NamesPrefund)
	types.RegisterKind(NamesRegular)
	types.RegisterKind(NamesBaddress)
}

type NamesPage struct {
	Kind          types.ListKind    `json:"kind"`
	Names         []*coreTypes.Name `json:"names"`
	TotalItems    int               `json:"totalItems"`
	ExpectedTotal int               `json:"expectedTotal"`
	IsFetching    bool              `json:"isFetching"`
	State         types.LoadState   `json:"state"`
}

func (np *NamesPage) GetKind() types.ListKind {
	return np.Kind
}

func (np *NamesPage) GetTotalItems() int {
	return np.TotalItems
}

func (np *NamesPage) GetExpectedTotal() int {
	return np.ExpectedTotal
}

func (np *NamesPage) GetIsFetching() bool {
	return np.IsFetching
}

func (np *NamesPage) GetState() types.LoadState {
	return np.State
}

type NamesCollection struct {
	allFacet      *facets.Facet[coreTypes.Name]
	customFacet   *facets.Facet[coreTypes.Name]
	prefundFacet  *facets.Facet[coreTypes.Name]
	regularFacet  *facets.Facet[coreTypes.Name]
	baddressFacet *facets.Facet[coreTypes.Name]
	summary       types.Summary
	summaryMutex  sync.RWMutex
}

func NewNamesCollection() *NamesCollection {
	namesStore := GetNamesStore()

	nc := &NamesCollection{}
	nc.ResetSummary()

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
		func(name *coreTypes.Name) bool { return name.Parts&coreTypes.Custom != 0 },
		nil,
		namesStore,
		"names",
		nc,
	)

	prefundFacet := facets.NewFacetWithSummary(
		NamesPrefund,
		func(name *coreTypes.Name) bool { return name.Parts&coreTypes.Prefund != 0 },
		nil,
		namesStore,
		"names",
		nc,
	)

	regularFacet := facets.NewFacetWithSummary(
		NamesRegular,
		func(name *coreTypes.Name) bool { return name.Parts&coreTypes.Regular != 0 },
		nil,
		namesStore,
		"names",
		nc,
	)

	baddressFacet := facets.NewFacetWithSummary(
		NamesBaddress,
		func(name *coreTypes.Name) bool { return name.Parts&coreTypes.Baddress != 0 },
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

	return nc
}

func (nc *NamesCollection) LoadData(listKind types.ListKind) {
	if !nc.NeedsUpdate(listKind) {
		return
	}

	var facet *facets.Facet[coreTypes.Name]
	var facetName types.ListKind

	switch listKind {
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
		logging.LogError("LoadData: unexpected list kind: %v", fmt.Errorf("invalid list kind: %s", listKind), nil)
		return
	}

	go func() {
		if err := facet.Load(); err != nil {
			logging.LogError(fmt.Sprintf("LoadData.%s from store: %%v", facetName), err, facets.ErrAlreadyLoading)
		}
	}()
}

func (nc *NamesCollection) Reset(listKind types.ListKind) {
	switch listKind {
	case NamesAll, NamesCustom, NamesPrefund, NamesRegular, NamesBaddress:
		namesStore.Reset()
	default:
		return
	}
}

func (nc *NamesCollection) NeedsUpdate(listKind types.ListKind) bool {
	var facet *facets.Facet[coreTypes.Name]

	switch listKind {
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

func (nc *NamesCollection) getExpectedTotal(listKind types.ListKind) int {
	_ = listKind
	if count, err := GetNamesCount(); err == nil && count > 0 {
		return count
	}
	return nc.allFacet.ExpectedCount()
}

func (nc *NamesCollection) GetPage(
	listKind types.ListKind,
	first, pageSize int,
	sortSpec sdk.SortSpec,
	filter string,
) (types.Page, error) {
	var facet *facets.Facet[coreTypes.Name]

	switch listKind {
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
		// This is truly a validation error - invalid ListKind for this collection
		return nil, types.NewValidationError("names", listKind, "GetPage",
			fmt.Errorf("unsupported list kind: %v", listKind))
	}

	var filterFunc func(*coreTypes.Name) bool
	if filter != "" {
		filterFunc = func(name *coreTypes.Name) bool {
			return nc.matchesFilter(name, filter)
		}
	}

	sortFunc := func(items []coreTypes.Name, sort sdk.SortSpec) error {
		return sdk.SortNames(items, sort)
	}

	pageResult, err := facet.GetPage(
		first,
		pageSize,
		filterFunc,
		sortSpec,
		sortFunc,
	)
	if err != nil {
		// This is likely an SDK or store error, not a validation error
		return nil, types.NewStoreError("names", listKind, "GetPage", err)
	}

	names := make([]*coreTypes.Name, 0, len(pageResult.Items))
	for i := range pageResult.Items {
		names = append(names, &pageResult.Items[i])
	}

	return &NamesPage{
		Kind:          listKind,
		Names:         names,
		TotalItems:    pageResult.TotalItems,
		ExpectedTotal: nc.getExpectedTotal(listKind),
		IsFetching:    facet.IsFetching(),
		State:         pageResult.State,
	}, nil
}

func (nc *NamesCollection) matchesFilter(name *coreTypes.Name, filter string) bool {
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

func (nc *NamesCollection) GetNamesPage(
	listKind types.ListKind,
	first, pageSize int,
	sortSpec sdk.SortSpec,
	filter string,
) (*NamesPage, error) {
	page, err := nc.GetPage(listKind, first, pageSize, sortSpec, filter)
	if err != nil {
		return nil, err
	}

	namesPage, ok := page.(*NamesPage)
	if !ok {
		return nil, fmt.Errorf("internal error: GetPage returned unexpected type %T", page)
	}

	return namesPage, nil
}

func (nc *NamesCollection) GetSupportedKinds() []types.ListKind {
	return []types.ListKind{
		NamesAll,
		NamesCustom,
		NamesPrefund,
		NamesRegular,
		NamesBaddress,
	}
}

func (nc *NamesCollection) GetStoreForKind(kind types.ListKind) string {
	switch kind {
	case NamesAll, NamesCustom, NamesPrefund, NamesRegular, NamesBaddress:
		return "names"
	default:
		return ""
	}
}

func (nc *NamesCollection) GetCollectionName() string {
	return "names"
}

func (nc *NamesCollection) GetCurrentSummary() types.Summary {
	nc.summaryMutex.RLock()
	defer nc.summaryMutex.RUnlock()
	return nc.summary
}

func (nc *NamesCollection) AccumulateItem(item interface{}, summary *types.Summary) {
	if name, ok := item.(*coreTypes.Name); ok {
		nc.summaryMutex.Lock()
		defer nc.summaryMutex.Unlock()

		nc.summary.TotalCount++
		if nc.summary.FacetCounts == nil {
			nc.summary.FacetCounts = make(map[types.ListKind]int)
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
		FacetCounts: make(map[types.ListKind]int),
		CustomData:  make(map[string]interface{}),
		LastUpdated: time.Now().Unix(),
	}
}

func (nc *NamesCollection) GetSummary() types.Summary {
	return nc.GetCurrentSummary()
}
