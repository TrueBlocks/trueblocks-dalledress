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
		assert.True(t, collection.matchesFilter(testMonitor, "1234567890123456"))
	})

	t.Run("NameMatch", func(t *testing.T) {
		assert.True(t, collection.matchesFilter(testMonitor, "test"))
		assert.True(t, collection.matchesFilter(testMonitor, "Monitor"))
		assert.True(t, collection.matchesFilter(testMonitor, "TEST")) // case insensitive
	})

	t.Run("NumericFieldsMatch", func(t *testing.T) {
		assert.True(t, collection.matchesFilter(testMonitor, "100"))   // NRecords
		assert.True(t, collection.matchesFilter(testMonitor, "1024"))  // FileSize
		assert.True(t, collection.matchesFilter(testMonitor, "12345")) // LastScanned
	})

	t.Run("BooleanFieldsMatch", func(t *testing.T) {
		assert.True(t, collection.matchesFilter(testMonitor, "staged"))
		emptyMonitor := &coreTypes.Monitor{IsEmpty: true}
		assert.True(t, collection.matchesFilter(emptyMonitor, "empty"))
		deletedMonitor := &coreTypes.Monitor{Deleted: true}
		assert.True(t, collection.matchesFilter(deletedMonitor, "deleted"))
	})

	t.Run("NoMatch", func(t *testing.T) {
		assert.False(t, collection.matchesFilter(testMonitor, "nonexistent"))
		assert.False(t, collection.matchesFilter(testMonitor, "xyz"))
	})

	t.Run("EmptyFilter", func(t *testing.T) {
		result := collection.matchesFilter(testMonitor, "")
		assert.True(t, result)
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

	t.Run("GetExpectedTotal", func(t *testing.T) {
		assert.NotPanics(t, func() {
			page, _ := collection.GetPage(MonitorsList, 0, 10, sdk.SortSpec{}, "")
			if page != nil {
				assert.GreaterOrEqual(t, page.ExpectedTotal, 0, "ExpectedTotal should be non-negative")
			}
		})
	})

	t.Run("FacetStateIntegration", func(t *testing.T) {
		needsUpdate := collection.NeedsUpdate(MonitorsList)
		assert.IsType(t, true, needsUpdate, "NeedsUpdate should return a boolean")

		page, err := collection.GetPage(MonitorsList, 0, 1, sdk.SortSpec{}, "")
		if err == nil && page != nil {
			assert.Equal(t, MonitorsList, page.Kind, "Page kind should match request")
			assert.GreaterOrEqual(t, page.TotalItems, 0, "TotalItems should be non-negative")
			assert.GreaterOrEqual(t, page.ExpectedTotal, 0, "ExpectedTotal should be non-negative")
			assert.IsType(t, true, page.IsFetching, "IsFetching should be a boolean")
		}
	})
}

func TestMonitorsCollectionLoadData(t *testing.T) {
	collection := NewMonitorsCollection()

	t.Run("LoadDataValidKind", func(t *testing.T) {
		assert.NotPanics(t, func() {
			collection.LoadData(MonitorsList)
		}, "LoadData with valid list kind should not panic")

		needsUpdate := collection.NeedsUpdate(MonitorsList)
		assert.True(t, needsUpdate, "LoadData call doesn't immediately change NeedsUpdate state")
	})

	t.Run("LoadDataInvalidKind", func(t *testing.T) {
		assert.NotPanics(t, func() {
			collection.LoadData("invalid-kind")
		}, "LoadData with invalid list kind should not panic")
	})

	t.Run("LoadDataWhenNotNeeded", func(t *testing.T) {

		assert.NotPanics(t, func() {
			collection.LoadData(MonitorsList)
			collection.LoadData(MonitorsList)
		})
	})

	t.Run("GetPageErrorHandling", func(t *testing.T) {
		_, err := collection.GetPage("invalid-kind", 0, 10, sdk.SortSpec{}, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported list kind")

		assert.NotPanics(t, func() {
			page, err := collection.GetPage(MonitorsList, 0, 10, sdk.SortSpec{}, "")
			if err == nil && page != nil {
				assert.Equal(t, MonitorsList, page.Kind)
				assert.GreaterOrEqual(t, page.TotalItems, 0)
				assert.GreaterOrEqual(t, page.ExpectedTotal, 0)
			}
		})
	})
}
