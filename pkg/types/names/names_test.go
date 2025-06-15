package names

import (
	"testing"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
	"github.com/stretchr/testify/assert"
)

func TestNewNamesCollection(t *testing.T) {
	assert.NotPanics(t, func() {
		collection := NewNamesCollection()
		assert.NotNil(t, collection)
	}, "NewNamesCollection should not panic")
}

func TestNamesMatchesFilter(t *testing.T) {
	collection := NewNamesCollection()
	testName := &coreTypes.Name{
		Address:    base.HexToAddress("0x1234567890123456789012345678901234567890"),
		Name:       "Test Name",
		Tags:       "testing,example",
		Source:     "test",
		Symbol:     "TEST",
		Decimals:   18,
		IsCustom:   true,
		IsPrefund:  false,
		IsContract: true,
		IsErc20:    true,
		IsErc721:   false,
		Parts:      coreTypes.Custom,
	}

	t.Run("AddressMatch", func(t *testing.T) {
		assert.True(t, collection.matchesFilter(testName, "1234"))
		assert.True(t, collection.matchesFilter(testName, "0x1234"))
	})

	t.Run("NameMatch", func(t *testing.T) {
		assert.True(t, collection.matchesFilter(testName, "test"))
		assert.True(t, collection.matchesFilter(testName, "Name"))
	})

	t.Run("TagsMatch", func(t *testing.T) {
		assert.True(t, collection.matchesFilter(testName, "testing"))
		assert.True(t, collection.matchesFilter(testName, "example"))
	})

	t.Run("SourceMatch", func(t *testing.T) {
		assert.True(t, collection.matchesFilter(testName, "test"))
	})

	t.Run("EmptyFilter", func(t *testing.T) {
		result := collection.matchesFilter(testName, "")
		assert.True(t, result)
	})

	t.Run("NoMatch", func(t *testing.T) {
		assert.False(t, collection.matchesFilter(testName, "nonexistent"))
	})
}

func TestNamesCollectionStateManagement(t *testing.T) {
	collection := NewNamesCollection()

	t.Run("NeedsUpdate", func(t *testing.T) {
		needsUpdate := collection.NeedsUpdate(NamesAll)
		assert.True(t, needsUpdate, "New collection should need update for NamesAll")

		needsUpdate = collection.NeedsUpdate(NamesCustom)
		assert.True(t, needsUpdate, "New collection should need update for NamesCustom")

		needsUpdate = collection.NeedsUpdate(NamesPrefund)
		assert.True(t, needsUpdate, "New collection should need update for NamesPrefund")

		needsUpdate = collection.NeedsUpdate(NamesRegular)
		assert.True(t, needsUpdate, "New collection should need update for NamesRegular")

		needsUpdate = collection.NeedsUpdate(NamesBaddress)
		assert.True(t, needsUpdate, "New collection should need update for NamesBaddress")

		needsUpdate = collection.NeedsUpdate("invalid-kind")
		assert.False(t, needsUpdate, "Invalid list kind should return false")
	})

	t.Run("Reset", func(t *testing.T) {
		assert.NotPanics(t, func() {
			collection.Reset(NamesAll)
		}, "Reset with NamesAll should not panic")

		assert.NotPanics(t, func() {
			collection.Reset(NamesCustom)
		}, "Reset with NamesCustom should not panic")

		assert.NotPanics(t, func() {
			collection.Reset(NamesPrefund)
		}, "Reset with NamesPrefund should not panic")

		assert.NotPanics(t, func() {
			collection.Reset(NamesRegular)
		}, "Reset with NamesRegular should not panic")

		assert.NotPanics(t, func() {
			collection.Reset(NamesBaddress)
		}, "Reset with NamesBaddress should not panic")

		needsUpdate := collection.NeedsUpdate(NamesAll)
		assert.True(t, needsUpdate, "After reset, collection should need update")

		assert.NotPanics(t, func() {
			collection.Reset("invalid-kind")
		}, "Reset with invalid list kind should not panic")
	})
}

func TestNamesCollectionLoadData(t *testing.T) {
	collection := NewNamesCollection()

	t.Run("LoadDataValidKinds", func(t *testing.T) {
		assert.NotPanics(t, func() {
			collection.LoadData(NamesAll)
		}, "LoadData with NamesAll should not panic")

		assert.NotPanics(t, func() {
			collection.LoadData(NamesCustom)
		}, "LoadData with NamesCustom should not panic")

		assert.NotPanics(t, func() {
			collection.LoadData(NamesPrefund)
		}, "LoadData with NamesPrefund should not panic")

		assert.NotPanics(t, func() {
			collection.LoadData(NamesRegular)
		}, "LoadData with NamesRegular should not panic")

		assert.NotPanics(t, func() {
			collection.LoadData(NamesBaddress)
		}, "LoadData with NamesBaddress should not panic")
	})

	t.Run("LoadDataInvalidKind", func(t *testing.T) {
		assert.NotPanics(t, func() {
			collection.LoadData("invalid-kind")
		}, "LoadData with invalid list kind should not panic")
	})
}

func TestNamesCollectionGetPage(t *testing.T) {
	collection := NewNamesCollection()

	t.Run("BasicGetPageNamesAll", func(t *testing.T) {
		page, err := collection.GetPage(NamesAll, 0, 10, sdk.SortSpec{}, "")

		if err == nil && page != nil {
			assert.Equal(t, NamesAll, page.Kind)
			assert.GreaterOrEqual(t, page.TotalItems, 0)
			assert.GreaterOrEqual(t, page.ExpectedTotal, 0)
			assert.LessOrEqual(t, len(page.Names), 10, "Returned items should not exceed page size")
		}
	})

	t.Run("BasicGetPageNamesCustom", func(t *testing.T) {
		page, err := collection.GetPage(NamesCustom, 0, 10, sdk.SortSpec{}, "")

		if err == nil && page != nil {
			assert.Equal(t, NamesCustom, page.Kind)
			assert.GreaterOrEqual(t, page.TotalItems, 0)
			assert.GreaterOrEqual(t, page.ExpectedTotal, 0)
			assert.LessOrEqual(t, len(page.Names), 10, "Returned items should not exceed page size")
		}
	})

	t.Run("BasicGetPageNamesPrefund", func(t *testing.T) {
		page, err := collection.GetPage(NamesPrefund, 0, 10, sdk.SortSpec{}, "")

		if err == nil && page != nil {
			assert.Equal(t, NamesPrefund, page.Kind)
			assert.GreaterOrEqual(t, page.TotalItems, 0)
			assert.GreaterOrEqual(t, page.ExpectedTotal, 0)
			assert.LessOrEqual(t, len(page.Names), 10, "Returned items should not exceed page size")
		}
	})

	t.Run("BasicGetPageNamesRegular", func(t *testing.T) {
		page, err := collection.GetPage(NamesRegular, 0, 10, sdk.SortSpec{}, "")

		if err == nil && page != nil {
			assert.Equal(t, NamesRegular, page.Kind)
			assert.GreaterOrEqual(t, page.TotalItems, 0)
			assert.GreaterOrEqual(t, page.ExpectedTotal, 0)
			assert.LessOrEqual(t, len(page.Names), 10, "Returned items should not exceed page size")
		}
	})

	t.Run("BasicGetPageNamesBaddress", func(t *testing.T) {
		page, err := collection.GetPage(NamesBaddress, 0, 10, sdk.SortSpec{}, "")

		if err == nil && page != nil {
			assert.Equal(t, NamesBaddress, page.Kind)
			assert.GreaterOrEqual(t, page.TotalItems, 0)
			assert.GreaterOrEqual(t, page.ExpectedTotal, 0)
			assert.LessOrEqual(t, len(page.Names), 10, "Returned items should not exceed page size")
		}
	})

	t.Run("GetPageWithFilter", func(t *testing.T) {
		page, err := collection.GetPage(NamesAll, 0, 10, sdk.SortSpec{}, "test")

		if err == nil && page != nil {
			assert.Equal(t, NamesAll, page.Kind)
			assert.GreaterOrEqual(t, page.TotalItems, 0)
		}
	})

	t.Run("InvalidListKind", func(t *testing.T) {
		_, err := collection.GetPage("invalid-kind", 0, 10, sdk.SortSpec{}, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unexpected list kind")
	})

	t.Run("ZeroPageSize", func(t *testing.T) {
		page, err := collection.GetPage(NamesAll, 0, 0, sdk.SortSpec{}, "")
		if err == nil && page != nil {
			assert.Equal(t, NamesAll, page.Kind)
			assert.Len(t, page.Names, 0, "Zero page size should return no items")
		}
	})
}
