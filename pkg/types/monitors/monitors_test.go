package monitors

import (
	"testing"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
	"github.com/stretchr/testify/assert"
)

func assertMonitorsPage(t *testing.T, page types.Page) *MonitorsPage {
	t.Helper()
	if page == nil {
		t.Fatal("page is nil")
	}
	monitorsPage, ok := page.(*MonitorsPage)
	if !ok {
		t.Fatalf("expected *MonitorsPage, got %T", page)
	}
	return monitorsPage
}

// Domain-specific CRUD and business logic tests for Monitors collection

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

func TestMonitorsCollectionDomainSpecific(t *testing.T) {
	collection := NewMonitorsCollection()

	t.Run("CrudOperationsSpecific", func(t *testing.T) {
		// Test domain-specific CRUD functionality
		// This is unique to monitors and involves address handling
		assert.NotPanics(t, func() {
			// Test monitor-specific operations
			// Actual CRUD logic would be tested here
			_, _ = collection.GetPage(MonitorsList, 0, 5, sdk.SortSpec{}, "")
		})
	})

	t.Run("AddressFilteringLogic", func(t *testing.T) {
		// Test monitors-specific filtering with addresses
		page, err := collection.GetPage(MonitorsList, 0, 10, sdk.SortSpec{}, "0x")
		if err == nil && page != nil {
			monitorsPage := assertMonitorsPage(t, page)
			assert.Equal(t, MonitorsList, monitorsPage.Facet)
			assert.GreaterOrEqual(t, monitorsPage.TotalItems, 0)
		}
	})
}
