// MONITORS_ROUTE
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
	State         facets.LoadState    `json:"state"`
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

	switch listKind {
	case MonitorsList:
		go func() {
			if result, err := mc.monitorsFacet.Load(); err != nil {
				logging.LogError("LoadData.MonitorsList from store: %v", err, facets.ErrAlreadyLoading)
			} else {
				msgs.EmitLoaded("monitors", result.Payload)
			}
		}()
	}
}

func (mc *MonitorsCollection) GetPage(
	kind types.ListKind,
	first, pageSize int,
	sort sdk.SortSpec,
	filter string,
) (*MonitorsPage, error) {
	switch kind {
	case MonitorsList:
		// Create filter function for monitors
		var filterFunc func(*coreTypes.Monitor) bool
		if filter != "" {
			filterFunc = func(monitor *coreTypes.Monitor) bool {
				return mc.matchesFilter(monitor, filter)
			}
		}

		// Create sort function for monitors
		var sortFunc func([]coreTypes.Monitor, sdk.SortSpec) error
		sortFunc = func(items []coreTypes.Monitor, sort sdk.SortSpec) error {
			return sdk.SortMonitors(items, sort)
		}

		pageResult, err := mc.monitorsFacet.GetPage(
			first,
			pageSize,
			filterFunc,
			sort,
			sortFunc,
		)
		if err != nil {
			return nil, err
		}

		return &MonitorsPage{
			Kind:          kind,
			Monitors:      pageResult.Items,
			TotalItems:    pageResult.TotalItems,
			ExpectedTotal: mc.getExpectedTotal(),
			IsFetching:    mc.monitorsFacet.IsFetching(),
			State:         pageResult.State,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported list kind: %s", kind)
	}
}

func (mc *MonitorsCollection) Reset(listKind types.ListKind) {
	switch listKind {
	case MonitorsList:
		monitorsStore.Reset()
	}
}

func (mc *MonitorsCollection) NeedsUpdate(listKind types.ListKind) bool {
	switch listKind {
	case MonitorsList:
		return mc.monitorsFacet.NeedsUpdate()
	}
	return false
}

func (mc *MonitorsCollection) getExpectedTotal() int {
	// Try to get accurate count from MonitorsCount, fallback to facet count
	if count, err := GetMonitorsCount(); err == nil && count > 0 {
		return count
	}
	// Fallback to current facet expected count
	return mc.monitorsFacet.ExpectedCount()
}

// matchesFilter checks if a monitor matches the given filter string
func (mc *MonitorsCollection) matchesFilter(monitor *coreTypes.Monitor, filter string) bool {
	filterLower := strings.ToLower(filter)

	// Check address (with and without 0x prefix)
	addressHex := strings.ToLower(monitor.Address.Hex())
	addressNoPrefix := strings.TrimPrefix(addressHex, "0x")
	addressNoLeadingZeros := strings.TrimLeft(addressNoPrefix, "0")

	if strings.Contains(addressHex, filterLower) ||
		strings.Contains(addressNoPrefix, filterLower) ||
		strings.Contains(addressNoLeadingZeros, filterLower) {
		return true
	}

	// Check name if available
	if strings.Contains(strings.ToLower(monitor.Name), filterLower) {
		return true
	}

	// Check numeric fields as strings
	if strings.Contains(fmt.Sprintf("%d", monitor.NRecords), filterLower) ||
		strings.Contains(fmt.Sprintf("%d", monitor.FileSize), filterLower) ||
		strings.Contains(fmt.Sprintf("%d", monitor.LastScanned), filterLower) {
		return true
	}

	// Check boolean fields
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

// MONITORS_ROUTE
