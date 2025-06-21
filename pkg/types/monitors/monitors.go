package monitors

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/facets"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/logging"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

const (
	MonitorsList types.DataFacet = "monitors"
)

func init() {
	types.RegisterDataFacet(MonitorsList)
}

type MonitorsCollection struct {
	monitorsFacet *facets.Facet[Monitor]
	summary       types.Summary
	summaryMutex  sync.RWMutex
}

func NewMonitorsCollection() *MonitorsCollection {
	c := &MonitorsCollection{
		summary: types.Summary{
			TotalCount:  0,
			FacetCounts: make(map[types.DataFacet]int),
			CustomData:  make(map[string]interface{}),
		},
	}
	c.initializeFacets()
	return c
}

func (c *MonitorsCollection) initializeFacets() {
	monitorsStore := GetMonitorsStore()

	monitorsFacet := facets.NewFacetWithSummary(
		MonitorsList,
		nil,
		nil,
		monitorsStore,
		"monitors",
		c,
	)

	c.monitorsFacet = monitorsFacet
}

func (mc *MonitorsCollection) LoadData(dataFacet types.DataFacet) {
	if !mc.NeedsUpdate(dataFacet) {
		return
	}

	var facet *facets.Facet[Monitor]
	var facetName string

	switch dataFacet {
	case MonitorsList:
		facet = mc.monitorsFacet
		facetName = "monitors"
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

func (c *MonitorsCollection) Reset(dataFacet types.DataFacet) {
	switch dataFacet {
	case MonitorsList:
		monitorsStore.Reset()
	default:
		return
	}
}

func (c *MonitorsCollection) NeedsUpdate(dataFacet types.DataFacet) bool {
	var facet *facets.Facet[Monitor]

	switch dataFacet {
	case MonitorsList:
		facet = c.monitorsFacet
	default:
		return false
	}

	return facet.NeedsUpdate()
}

func (mc *MonitorsCollection) getExpectedTotal(dataFacet types.DataFacet) int {
	_ = dataFacet
	if count, err := GetMonitorsCount(); err == nil && count > 0 {
		return count
	}
	return mc.monitorsFacet.ExpectedCount()
}

func (mc *MonitorsCollection) matchesFilter(monitor *Monitor, filter string) bool {
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

func (c *MonitorsCollection) GetMonitorsPage(
	dataFacet types.DataFacet,
	first, pageSize int,
	sortSpec sdk.SortSpec,
	filter string,
) (*MonitorsPage, error) {
	page, err := c.GetPage(dataFacet, first, pageSize, sortSpec, filter)
	if err != nil {
		return nil, err
	}

	monitorsPage, ok := page.(*MonitorsPage)
	if !ok {
		return nil, fmt.Errorf("internal error: GetPage returned unexpected type %T", page)
	}

	return monitorsPage, nil
}

func (c *MonitorsCollection) GetSupportedFacets() []types.DataFacet {
	return []types.DataFacet{
		MonitorsList,
	}
}

func (c *MonitorsCollection) GetStoreForFacet(dataFacet types.DataFacet) string {
	switch dataFacet {
	case MonitorsList:
		return "monitors"
	default:
		return ""
	}
}

func (c *MonitorsCollection) GetCollectionName() string {
	return "monitors"
}

func (c *MonitorsCollection) AccumulateItem(item interface{}, summary *types.Summary) {
	monitor, ok := item.(*Monitor)
	if !ok {
		return
	}

	c.summaryMutex.Lock()
	defer c.summaryMutex.Unlock()

	summary.TotalCount++

	if summary.FacetCounts == nil {
		summary.FacetCounts = make(map[types.DataFacet]int)
	}

	summary.FacetCounts[MonitorsList]++
	if summary.CustomData == nil {
		summary.CustomData = make(map[string]interface{})
	}

	emptyCount, _ := summary.CustomData["emptyCount"].(int)
	stagedCount, _ := summary.CustomData["stagedCount"].(int)
	deletedCount, _ := summary.CustomData["deletedCount"].(int)
	totalRecords, _ := summary.CustomData["totalRecords"].(int)
	totalFileSize, _ := summary.CustomData["totalFileSize"].(int64)

	if monitor.IsEmpty {
		emptyCount++
	}
	if monitor.IsStaged {
		stagedCount++
	}
	if monitor.Deleted {
		deletedCount++
	}

	totalRecords += int(monitor.NRecords)
	totalFileSize += int64(monitor.FileSize)

	summary.CustomData["emptyCount"] = emptyCount
	summary.CustomData["stagedCount"] = stagedCount
	summary.CustomData["deletedCount"] = deletedCount
	summary.CustomData["totalRecords"] = totalRecords
	summary.CustomData["totalFileSize"] = totalFileSize
}

func (c *MonitorsCollection) GetSummary() types.Summary {
	c.summaryMutex.RLock()
	defer c.summaryMutex.RUnlock()

	summary := c.summary
	summary.FacetCounts = make(map[types.DataFacet]int)
	for k, v := range c.summary.FacetCounts {
		summary.FacetCounts[k] = v
	}

	if c.summary.CustomData != nil {
		summary.CustomData = make(map[string]interface{})
		for k, v := range c.summary.CustomData {
			summary.CustomData[k] = v
		}
	}

	return summary
}

func (c *MonitorsCollection) ResetSummary() {
	c.summaryMutex.Lock()
	defer c.summaryMutex.Unlock()
	c.summary = types.Summary{
		TotalCount:  0,
		FacetCounts: make(map[types.DataFacet]int),
		CustomData:  make(map[string]interface{}),
		LastUpdated: time.Now().Unix(),
	}
}
