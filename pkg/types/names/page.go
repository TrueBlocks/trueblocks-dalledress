package names

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/facets"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

type NamesPage struct {
	Facet         types.DataFacet `json:"facet"`
	Names         []*Name         `json:"names"`
	TotalItems    int             `json:"totalItems"`
	ExpectedTotal int             `json:"expectedTotal"`
	IsFetching    bool            `json:"isFetching"`
	State         types.LoadState `json:"state"`
}

func (p *NamesPage) GetFacet() types.DataFacet {
	return p.Facet
}

func (p *NamesPage) GetTotalItems() int {
	return p.TotalItems
}

func (p *NamesPage) GetExpectedTotal() int {
	return p.ExpectedTotal
}

func (p *NamesPage) GetIsFetching() bool {
	return p.IsFetching
}

func (p *NamesPage) GetState() types.LoadState {
	return p.State
}

func (c *NamesCollection) GetPage(
	dataFacet types.DataFacet,
	first, pageSize int,
	sortSpec sdk.SortSpec,
	filter string,
) (types.Page, error) {
	var facet *facets.Facet[Name]

	switch dataFacet {
	case NamesAll:
		facet = c.allFacet
	case NamesCustom:
		facet = c.customFacet
	case NamesPrefund:
		facet = c.prefundFacet
	case NamesRegular:
		facet = c.regularFacet
	case NamesBaddress:
		facet = c.baddressFacet
	default:
		// This is truly a validation error - invalid DataFacet for this collection
		return nil, types.NewValidationError("names", dataFacet, "GetPage",
			fmt.Errorf("unsupported dataFacet: %v", dataFacet))
	}

	var filterFunc func(*Name) bool
	if filter != "" {
		filterFunc = func(name *Name) bool {
			return c.matchesFilter(name, filter)
		}
	}

	sortFunc := func(items []Name, sort sdk.SortSpec) error {
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
		return nil, types.NewStoreError("names", dataFacet, "GetPage", err)
	}

	names := make([]*Name, 0, len(pageResult.Items))
	for i := range pageResult.Items {
		names = append(names, &pageResult.Items[i])
	}

	return &NamesPage{
		Facet:         dataFacet,
		Names:         names,
		TotalItems:    pageResult.TotalItems,
		ExpectedTotal: c.getExpectedTotal(dataFacet),
		IsFetching:    facet.IsFetching(),
		State:         pageResult.State,
	}, nil
}

func (c *NamesCollection) GetNamesPage(
	dataFacet types.DataFacet,
	first, pageSize int,
	sortSpec sdk.SortSpec,
	filter string,
) (*NamesPage, error) {
	page, err := c.GetPage(dataFacet, first, pageSize, sortSpec, filter)
	if err != nil {
		return nil, err
	}

	namesPage, ok := page.(*NamesPage)
	if !ok {
		return nil, fmt.Errorf("internal error: GetPage returned unexpected type %T", page)
	}

	return namesPage, nil
}

