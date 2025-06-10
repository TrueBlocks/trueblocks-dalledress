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
	filterString string,
) (*AbisPage, error) {
	page := &AbisPage{
		Kind: listKind,
	}
	filterString = strings.ToLower(filterString)

	switch listKind {
	case AbisDownloaded:
		items, total, state := ac.GetDownloadedPage(first, pageSize, sortSpec)
		page.Abis = items
		page.TotalItems = total
		page.State = state
		page.IsFetching = ac.downloadedFacet.IsFetching()
		page.ExpectedTotal = ac.downloadedFacet.ExpectedCount()

	case AbisKnown:
		items, total, state := ac.GetKnownPage(first, pageSize, sortSpec)
		page.Abis = items
		page.TotalItems = total
		page.State = state
		page.IsFetching = ac.knownFacet.IsFetching()
		page.ExpectedTotal = ac.knownFacet.ExpectedCount()

	case AbisFunctions:
		items, total, state := ac.GetFunctionsPage(first, pageSize, sortSpec)
		page.Functions = items
		page.TotalItems = total
		page.State = state
		page.IsFetching = ac.functionsFacet.IsFetching()
		page.ExpectedTotal = ac.functionsFacet.ExpectedCount()

	case AbisEvents:
		items, total, state := ac.GetEventsPage(first, pageSize, sortSpec)
		page.Functions = items
		page.TotalItems = total
		page.State = state
		page.IsFetching = ac.functionsFacet.IsFetching()
		page.ExpectedTotal = ac.functionsFacet.ExpectedCount()

	default:
		return nil, fmt.Errorf("GetPage: unexpected list kind: %v", listKind)
	}

	// TODO: FILTERING IS NOT WORKING
	if filterString != "" {
		logging.LogBackend(fmt.Sprintf("GetPage: filterString is present but not yet implemented: %s", filterString))
	}

	// TODO: SORTING IS NOT WORKING
	if equals(sortSpec, sdk.SortSpec{}) {
		logging.LogBackend("GetPage: sortSpec is present but not yet implemented")
	}

	return page, nil
}

func (ac *AbisCollection) GetDownloadedPage(first, pageSize int, sortSpec sdk.SortSpec) ([]coreTypes.Abi, int, facets.LoadState) {
	if result, err := ac.downloadedFacet.GetPage(first, pageSize, nil, sortSpec, nil); err != nil {
		return nil, 0, facets.StateError
	} else {
		return result.Items, result.TotalItems, result.State
	}
}

func (ac *AbisCollection) GetKnownPage(first, pageSize int, sortSpec sdk.SortSpec) ([]coreTypes.Abi, int, facets.LoadState) {
	if result, err := ac.knownFacet.GetPage(first, pageSize, nil, sortSpec, nil); err != nil {
		return nil, 0, facets.StateError
	} else {
		return result.Items, result.TotalItems, result.State
	}
}

func (ac *AbisCollection) GetFunctionsPage(first, pageSize int, sortSpec sdk.SortSpec) ([]coreTypes.Function, int, facets.LoadState) {
	if result, err := ac.functionsFacet.GetPage(first, pageSize, nil, sortSpec, nil); err != nil {
		return nil, 0, facets.StateError
	} else {
		return result.Items, result.TotalItems, result.State
	}
}

func (ac *AbisCollection) GetEventsPage(first, pageSize int, sortSpec sdk.SortSpec) ([]coreTypes.Function, int, facets.LoadState) {
	if result, err := ac.eventsFacet.GetPage(first, pageSize, nil, sortSpec, nil); err != nil {
		return nil, 0, facets.StateError
	} else {
		return result.Items, result.TotalItems, result.State
	}
}

// equals compares two sdk.SortSpecs and returns true if they are equal
func equals(s sdk.SortSpec, other sdk.SortSpec) bool {
	if len(s.Fields) != len(other.Fields) || len(s.Order) != len(other.Order) {
		return false
	}
	for i := range s.Fields {
		if s.Fields[i] != other.Fields[i] || s.Order[i] != other.Order[i] {
			return false
		}
	}
	return true
}
