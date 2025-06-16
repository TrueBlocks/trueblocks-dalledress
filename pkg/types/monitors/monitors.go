package monitors

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
	MonitorsList types.ListKind = "Monitors"
)

func init() {
	types.RegisterKind(MonitorsList)
}

type MonitorsPage struct {
	Kind          types.ListKind      `json:"kind"`
	Monitors      []coreTypes.Monitor `json:"monitors,omitempty"`
	TotalItems    int                 `json:"totalItems"`
	ExpectedTotal int                 `json:"expectedTotal"`
	IsFetching    bool                `json:"isFetching"`
	State         types.LoadState     `json:"state"`
}

func (mp *MonitorsPage) GetKind() types.ListKind {
	return mp.Kind
}

func (mp *MonitorsPage) GetTotalItems() int {
	return mp.TotalItems
}

func (mp *MonitorsPage) GetExpectedTotal() int {
	return mp.ExpectedTotal
}

func (mp *MonitorsPage) GetIsFetching() bool {
	return mp.IsFetching
}

func (mp *MonitorsPage) GetState() types.LoadState {
	return mp.State
}

type MonitorsCollection struct {
	monitorsFacet *facets.Facet[coreTypes.Monitor]
}

func NewMonitorsCollection() *MonitorsCollection {
	monitorsStore := GetMonitorsStore()

	monitorsFacet := facets.NewFacet(
		MonitorsList,
		nil,
		nil,
		monitorsStore,
	)

	return &MonitorsCollection{
		monitorsFacet: monitorsFacet,
	}
}

func (mc *MonitorsCollection) LoadData(listKind types.ListKind) {
	if !mc.NeedsUpdate(listKind) {
		return
	}

	var facet *facets.Facet[coreTypes.Monitor]
	var facetName string

	switch listKind {
	case MonitorsList:
		facet = mc.monitorsFacet
		facetName = "monitors"
	default:
		logging.LogError("LoadData: unexpected list kind: %v", fmt.Errorf("invalid list kind: %s", listKind), nil)
		return
	}

	go func() {
		if result, err := facet.Load(); err != nil {
			logging.LogError(fmt.Sprintf("LoadData.%s from store: %%v", facetName), err, facets.ErrAlreadyLoading)
		} else {
			msgs.EmitLoaded(facetName, result.Payload)
		}
	}()
}

func (mc *MonitorsCollection) GetPage(
	listKind types.ListKind,
	first, pageSize int,
	sortSpec sdk.SortSpec,
	filter string,
) (types.Page, error) {
	switch listKind {
	case MonitorsList:
		var filterFunc func(*coreTypes.Monitor) bool
		if filter != "" {
			filterFunc = func(monitor *coreTypes.Monitor) bool {
				return mc.matchesFilter(monitor, filter)
			}
		}

		var sortFunc func([]coreTypes.Monitor, sdk.SortSpec) error
		sortFunc = func(items []coreTypes.Monitor, sort sdk.SortSpec) error {
			return sdk.SortMonitors(items, sort)
		}

		pageResult, err := mc.monitorsFacet.GetPage(
			first,
			pageSize,
			filterFunc,
			sortSpec,
			sortFunc,
		)
		if err != nil {
			return nil, err
		}

		return &MonitorsPage{
			Kind:          listKind,
			Monitors:      pageResult.Items,
			TotalItems:    pageResult.TotalItems,
			ExpectedTotal: mc.getExpectedTotal(listKind),
			IsFetching:    mc.monitorsFacet.IsFetching(),
			State:         pageResult.State,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported list kind: %s", listKind)
	}
}

func (mc *MonitorsCollection) Reset(listKind types.ListKind) {
	switch listKind {
	case MonitorsList:
		monitorsStore.Reset()
	default:
		return
	}
}

func (mc *MonitorsCollection) NeedsUpdate(listKind types.ListKind) bool {
	var facet *facets.Facet[coreTypes.Monitor]

	switch listKind {
	case MonitorsList:
		facet = mc.monitorsFacet
	default:
		return false
	}

	return facet.NeedsUpdate()
}

func (mc *MonitorsCollection) getExpectedTotal(listKind types.ListKind) int {
	_ = listKind
	if count, err := GetMonitorsCount(); err == nil && count > 0 {
		return count
	}
	return mc.monitorsFacet.ExpectedCount()
}

func (mc *MonitorsCollection) matchesFilter(monitor *coreTypes.Monitor, filter string) bool {
	if filter == "" {
		return true
	}

	filterLower := strings.ToLower(filter)

	addressHex := strings.ToLower(monitor.Address.Hex())
	addressNoPrefix := strings.TrimPrefix(addressHex, "0x")
	addressNoLeadingZeros := strings.TrimLeft(addressNoPrefix, "0")

	if strings.Contains(addressHex, filterLower) ||
		strings.Contains(addressNoPrefix, filterLower) ||
		strings.Contains(addressNoLeadingZeros, filterLower) {
		return true
	}

	if strings.Contains(strings.ToLower(monitor.Name), filterLower) {
		return true
	}

	if strings.Contains(fmt.Sprintf("%d", monitor.NRecords), filterLower) ||
		strings.Contains(fmt.Sprintf("%d", monitor.FileSize), filterLower) ||
		strings.Contains(fmt.Sprintf("%d", monitor.LastScanned), filterLower) {
		return true
	}

	if monitor.IsEmpty && strings.Contains("empty", filterLower) {
		return true
	}
	if monitor.IsStaged && strings.Contains("staged", filterLower) {
		return true
	}
	if monitor.Deleted && strings.Contains("deleted", filterLower) {
		return true
	}

	return false
}

func (mc *MonitorsCollection) GetMonitorsPage(
	listKind types.ListKind,
	first, pageSize int,
	sortSpec sdk.SortSpec,
	filter string,
) (*MonitorsPage, error) {
	page, err := mc.GetPage(listKind, first, pageSize, sortSpec, filter)
	if err != nil {
		return nil, err
	}

	monitorsPage, ok := page.(*MonitorsPage)
	if !ok {
		return nil, fmt.Errorf("internal error: GetPage returned unexpected type %T", page)
	}

	return monitorsPage, nil
}
