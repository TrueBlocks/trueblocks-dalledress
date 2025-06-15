package monitors

import (
	"testing"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
	"github.com/stretchr/testify/assert"
)

func TestNewMonitorsCollection(t *testing.T) {
	assert.NotPanics(t, func() {
		collection := NewMonitorsCollection()
		assert.NotNil(t, collection)
	}, "NewMonitorsCollection should not panic")
}

func TestMonitorsMatchesFilter(t *testing.T) {
	collection := NewMonitorsCollection()
	testMonitor := &coreTypes.Monitor{
		Address:     base.HexToAddress("0x1234567890123456789012345678901234567890"),
		Name:        "Test Monitor",
		NRecords:    100,
		FileSize:    1024,
		LastScanned: 12345,
		IsEmpty:     false,
		IsStaged:    true,
		Deleted:     false,
	}

	t.Run("AddressMatch", func(t *testing.T) {
		assert.True(t, collection.matchesFilter(testMonitor, "1234"))
		assert.True(t, collection.matchesFilter(testMonitor, "0x1234"))
	})

	t.Run("NameMatch", func(t *testing.T) {
		assert.True(t, collection.matchesFilter(testMonitor, "test"))
		assert.True(t, collection.matchesFilter(testMonitor, "Monitor"))
	})

	t.Run("EmptyFilter", func(t *testing.T) {
		result := collection.matchesFilter(testMonitor, "")
		assert.True(t, result)
	})

	t.Run("NoMatch", func(t *testing.T) {
		assert.False(t, collection.matchesFilter(testMonitor, "nonexistent"))
	})
}

func TestMonitorsCollectionStateManagement(t *testing.T) {
	collection := NewMonitorsCollection()

	t.Run("NeedsUpdate", func(t *testing.T) {
		needsUpdate := collection.NeedsUpdate(MonitorsList)
		assert.True(t, needsUpdate, "New collection should need update")

		needsUpdate = collection.NeedsUpdate("invalid-kind")
		assert.False(t, needsUpdate, "Invalid list kind should return false")
	})

	t.Run("Reset", func(t *testing.T) {
		assert.NotPanics(t, func() {
			collection.Reset(MonitorsList)
		}, "Reset with valid list kind should not panic")

		needsUpdate := collection.NeedsUpdate(MonitorsList)
		assert.True(t, needsUpdate, "After reset, collection should need update")

		assert.NotPanics(t, func() {
			collection.Reset("invalid-kind")
		}, "Reset with invalid list kind should not panic")
	})
}

func TestMonitorsCollectionLoadData(t *testing.T) {
	collection := NewMonitorsCollection()

	t.Run("LoadDataValidKind", func(t *testing.T) {
		assert.NotPanics(t, func() {
			collection.LoadData(MonitorsList)
		}, "LoadData with valid list kind should not panic")
	})

	t.Run("LoadDataInvalidKind", func(t *testing.T) {
		assert.NotPanics(t, func() {
			collection.LoadData("invalid-kind")
		}, "LoadData with invalid list kind should not panic")
	})
}

func TestMonitorsCollectionGetPage(t *testing.T) {
	collection := NewMonitorsCollection()

	t.Run("BasicGetPage", func(t *testing.T) {
		page, err := collection.GetPage(MonitorsList, 0, 10, sdk.SortSpec{}, "")

		if err == nil && page != nil {
			assert.Equal(t, MonitorsList, page.Kind)
			assert.GreaterOrEqual(t, page.TotalItems, 0)
			assert.GreaterOrEqual(t, page.ExpectedTotal, 0)
			assert.LessOrEqual(t, len(page.Monitors), 10, "Returned items should not exceed page size")
		}
	})

	t.Run("GetPageWithFilter", func(t *testing.T) {
		page, err := collection.GetPage(MonitorsList, 0, 10, sdk.SortSpec{}, "test")

		if err == nil && page != nil {
			assert.Equal(t, MonitorsList, page.Kind)
			assert.GreaterOrEqual(t, page.TotalItems, 0)
		}
	})

	t.Run("InvalidListKind", func(t *testing.T) {
		_, err := collection.GetPage("invalid-kind", 0, 10, sdk.SortSpec{}, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported list kind")
	})

	t.Run("ZeroPageSize", func(t *testing.T) {
		page, err := collection.GetPage(MonitorsList, 0, 0, sdk.SortSpec{}, "")
		if err == nil && page != nil {
			assert.Equal(t, MonitorsList, page.Kind)
			assert.Len(t, page.Monitors, 0, "Zero page size should return no items")
		}
	})
}
