// ABIS_ROUTE
package abis

import (
	"fmt"
	"strings"

	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/facets"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/logging"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

const (
	AbisDownloaded types.ListKind = "Downloaded"
	AbisKnown      types.ListKind = "Known"
	AbisFunctions  types.ListKind = "Functions"
	AbisEvents     types.ListKind = "Events"
)

func init() {
	types.RegisterKind(AbisDownloaded)
	types.RegisterKind(AbisKnown)
	types.RegisterKind(AbisFunctions)
	types.RegisterKind(AbisEvents)
}

type AbisPage struct {
	Kind          types.ListKind       `json:"kind"`
	Abis          []coreTypes.Abi      `json:"abis,omitempty"`
	Functions     []coreTypes.Function `json:"functions,omitempty"`
	TotalItems    int                  `json:"totalItems"`
	ExpectedTotal int                  `json:"expectedTotal"`
	IsFetching    bool                 `json:"isFetching"`
	State         facets.LoadState     `json:"state"`
}

type AbisCollection struct {
	downloadedFacet *facets.Facet[coreTypes.Abi]
	knownFacet      *facets.Facet[coreTypes.Abi]
	functionsFacet  *facets.Facet[coreTypes.Function]
	eventsFacet     *facets.Facet[coreTypes.Function]
}

func NewAbisCollection() *AbisCollection {
	abisListStore := GetAbisListStore()
	abisDetailStore := GetAbisDetailStore()

	downloadedFacet := facets.NewFacet(
		AbisDownloaded,
		isDownload,
		nil,
		abisListStore,
	)

	knownFacet := facets.NewFacet(
		AbisKnown,
		isKnown,
		nil,
		abisListStore,
	)

	functionsFacet := facets.NewFacet(
		AbisFunctions,
		func(fn *coreTypes.Function) bool { return fn.FunctionType != "event" },
		isDupEncoding(),
		abisDetailStore,
	)

	eventsFacet := facets.NewFacet(
		AbisEvents,
		func(fn *coreTypes.Function) bool { return fn.FunctionType == "event" },
		isDupEncoding(),
		abisDetailStore,
	)

	return &AbisCollection{
		downloadedFacet: downloadedFacet,
		knownFacet:      knownFacet,
		functionsFacet:  functionsFacet,
		eventsFacet:     eventsFacet,
	}
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

func (ac *AbisCollection) LoadData(listKind types.ListKind) {
	if !ac.NeedsUpdate(listKind) {
		return
	}

	switch listKind {
	case AbisDownloaded:
		go func() {
			if result, err := ac.downloadedFacet.Load(); err != nil {
				logging.LogError("LoadData.AbisDownloaded from store: %v", err, facets.ErrAlreadyLoading)
			} else {
				msgs.EmitLoaded("downloaded", result.Payload)
			}
		}()
	case AbisKnown:
		go func() {
			if result, err := ac.knownFacet.Load(); err != nil {
				logging.LogError("LoadData.AbisKnown from store: %v", err, facets.ErrAlreadyLoading)
			} else {
				msgs.EmitLoaded("known", result.Payload)
			}
		}()
	case AbisFunctions:
		go func() {
			if result, err := ac.functionsFacet.Load(); err != nil {
				logging.LogError("LoadData.AbisFunctions from store: %v", err, facets.ErrAlreadyLoading)
			} else {
				msgs.EmitLoaded("functions", result.Payload)
			}
		}()
	case AbisEvents:
		go func() {
			if result, err := ac.eventsFacet.Load(); err != nil {
				logging.LogError("LoadData.AbisEvents from store: %v", err, facets.ErrAlreadyLoading)
			} else {
				msgs.EmitLoaded("events", result.Payload)
			}
		}()
	default:
		logging.LogError("LoadData: unexpected list kind: %v", fmt.Errorf("invalid list kind: %s", listKind), nil)
	}
}

func (ac *AbisCollection) Reset(listKind types.ListKind) {
	switch listKind {
	case AbisDownloaded, AbisKnown:
		abisListStore.Reset()
	case AbisFunctions, AbisEvents:
		abisDetailStore.Reset()
	}
}

func (ac *AbisCollection) NeedsUpdate(listKind types.ListKind) bool {
	switch listKind {
	case AbisDownloaded:
		return ac.downloadedFacet.NeedsUpdate()
	case AbisKnown:
		return ac.knownFacet.NeedsUpdate()
	case AbisFunctions:
		return ac.functionsFacet.NeedsUpdate()
	case AbisEvents:
		return ac.eventsFacet.NeedsUpdate()
	}
	return false
}

func (ac *AbisCollection) GetPage(
	listKind types.ListKind,
	first, pageSize int,
	sortSpec sdk.SortSpec,
	filter string,
) (*AbisPage, error) {
	page := &AbisPage{
		Kind: listKind,
	}
	filter = strings.ToLower(filter)

	var listFacet *facets.Facet[coreTypes.Abi]
	var detailFacet *facets.Facet[coreTypes.Function]

	switch listKind {
	case AbisDownloaded:
		listFacet = ac.downloadedFacet
	case AbisKnown:
		listFacet = ac.knownFacet
	case AbisFunctions:
		detailFacet = ac.functionsFacet
	case AbisEvents:
		detailFacet = ac.eventsFacet
	default:
		return nil, fmt.Errorf("GetPage: unexpected list kind: %v", listKind)
	}

	if listFacet != nil {
		var listFilterFunc = func(item *coreTypes.Abi) bool {
			return strings.Contains(strings.ToLower(item.Name), filter)
		}
		var listSortFunc = func(items []coreTypes.Abi, sort sdk.SortSpec) error {
			return sdk.SortAbis(items, sort)
		}
		if result, err := listFacet.GetPage(first, pageSize, listFilterFunc, sortSpec, listSortFunc); err != nil {
			page.Abis, page.TotalItems, page.State = nil, 0, facets.StateError
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
			page.Functions, page.TotalItems, page.State = nil, 0, facets.StateError
		} else {
			page.Functions, page.TotalItems, page.State = result.Items, result.TotalItems, result.State
		}
		page.IsFetching = detailFacet.IsFetching()
		page.ExpectedTotal = detailFacet.ExpectedCount()

	} else {
		return nil, fmt.Errorf("GetPage: no facet found for list kind: %v", listKind)
	}

	return page, nil
}

// ABIS_ROUTE
