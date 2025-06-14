// NAMES_ROUTE
package names

import (
	"fmt"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/facets"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/logging"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
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
	State         facets.LoadState  `json:"state"`
}

type NamesCollection struct {
	allFacet      *facets.Facet[coreTypes.Name]
	customFacet   *facets.Facet[coreTypes.Name]
	prefundFacet  *facets.Facet[coreTypes.Name]
	regularFacet  *facets.Facet[coreTypes.Name]
	baddressFacet *facets.Facet[coreTypes.Name]
}

func NewNamesCollection() *NamesCollection {
	namesStore := GetNamesStore()

	allFacet := facets.NewFacet(
		NamesAll,
		nil,
		nil,
		namesStore,
	)

	customFacet := facets.NewFacet(
		NamesCustom,
		func(name *coreTypes.Name) bool { return name.Parts&coreTypes.Custom != 0 },
		nil,
		namesStore,
	)

	prefundFacet := facets.NewFacet(
		NamesPrefund,
		func(name *coreTypes.Name) bool { return name.Parts&coreTypes.Prefund != 0 },
		nil,
		namesStore,
	)

	regularFacet := facets.NewFacet(
		NamesRegular,
		func(name *coreTypes.Name) bool { return name.Parts&coreTypes.Regular != 0 },
		nil,
		namesStore,
	)

	baddressFacet := facets.NewFacet(
		NamesBaddress,
		func(name *coreTypes.Name) bool { return name.Parts&coreTypes.Baddress != 0 },
		nil,
		namesStore,
	)

	return &NamesCollection{
		allFacet:      allFacet,
		customFacet:   customFacet,
		prefundFacet:  prefundFacet,
		regularFacet:  regularFacet,
		baddressFacet: baddressFacet,
	}
}

func (nc *NamesCollection) LoadData(listKind types.ListKind) {
	if !nc.NeedsUpdate(listKind) {
		return
	}

	var facet *facets.Facet[coreTypes.Name]
	var facetName string

	switch listKind {
	case NamesAll:
		facet = nc.allFacet
		facetName = "all"
	case NamesCustom:
		facet = nc.customFacet
		facetName = "custom"
	case NamesPrefund:
		facet = nc.prefundFacet
		facetName = "prefund"
	case NamesRegular:
		facet = nc.regularFacet
		facetName = "regular"
	case NamesBaddress:
		facet = nc.baddressFacet
		facetName = "baddress"
	default:
		logging.LogError("LoadData: unexpected list kind: %v", fmt.Errorf("invalid list kind: %s", listKind), nil)
		return
	}

	// Single goroutine implementation for all facets
	go func() {
		if result, err := facet.Load(); err != nil {
			logging.LogError(fmt.Sprintf("LoadData.%s from store: %%v", facetName), err, facets.ErrAlreadyLoading)
		} else {
			msgs.EmitLoaded(facetName, result.Payload)
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
	_ = listKind // delinter
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
) (*NamesPage, error) {
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
		return nil, fmt.Errorf("GetPage: unexpected list kind: %v", listKind)
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
		return nil, err
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

func (nc *NamesCollection) FindNameByAddress(addr base.Address) (*coreTypes.Name, bool) {
	if nc.allFacet != nil {
		found := false
		var result *coreTypes.Name

		_, _ = nc.allFacet.ForEvery(
			func(name *coreTypes.Name) (error, bool) {
				result = name
				found = true
				return nil, true
			},
			func(name *coreTypes.Name) bool {
				return name.Address == addr
			},
		)

		if found {
			return result, true
		}
	}

	return nil, false
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

// NAMES_ROUTE
