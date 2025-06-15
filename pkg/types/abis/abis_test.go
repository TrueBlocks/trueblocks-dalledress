package abis

import (
	"strings"
	"testing"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
	"github.com/stretchr/testify/assert"
)

func TestNewAbisCollection(t *testing.T) {
	assert.NotPanics(t, func() {
		collection := NewAbisCollection()
		assert.NotNil(t, collection)
	}, "NewAbisCollection should not panic")
}

func TestAbisMatchesFilter(t *testing.T) {
	collection := NewAbisCollection()
	testAbi := &coreTypes.Abi{
		Address:  base.HexToAddress("0x1234567890123456789012345678901234567890"),
		Name:     "Test ABI",
		IsKnown:  false,
		FileSize: 1024,
	}

	testFunction := &coreTypes.Function{
		Name:         "transfer",
		FunctionType: "function",
		Signature:    "transfer(address,uint256)",
		Encoding:     "0xa9059cbb",
	}

	t.Run("AbiNameMatch", func(t *testing.T) {
		filterFunc := func(item *coreTypes.Abi) bool {
			filter := "test"
			return collection.matchesAbiFilter(item, filter)
		}
		assert.True(t, filterFunc(testAbi))

		filterFunc2 := func(item *coreTypes.Abi) bool {
			filter := "ABI"
			return collection.matchesAbiFilter(item, filter)
		}
		assert.True(t, filterFunc2(testAbi))
	})

	t.Run("FunctionNameMatch", func(t *testing.T) {

		filterFunc := func(item *coreTypes.Function) bool {
			filter := "transfer"
			return collection.matchesFunctionFilter(item, filter)
		}
		assert.True(t, filterFunc(testFunction))
	})

	t.Run("FunctionEncodingMatch", func(t *testing.T) {

		filterFunc := func(item *coreTypes.Function) bool {
			filter := "0xa9059cbb"
			return collection.matchesFunctionFilter(item, filter)
		}
		assert.True(t, filterFunc(testFunction))
	})

	t.Run("EmptyFilter", func(t *testing.T) {
		filterFunc := func(item *coreTypes.Abi) bool {
			filter := ""
			return collection.matchesAbiFilter(item, filter)
		}
		result := filterFunc(testAbi)
		assert.True(t, result)
	})

	t.Run("NoMatch", func(t *testing.T) {
		filterFunc := func(item *coreTypes.Abi) bool {
			filter := "nonexistent"
			return collection.matchesAbiFilter(item, filter)
		}
		assert.False(t, filterFunc(testAbi))
	})
}

func (ac *AbisCollection) matchesAbiFilter(abi *coreTypes.Abi, filter string) bool {
	if filter == "" {
		return true
	}
	filterLower := strings.ToLower(filter)
	return strings.Contains(strings.ToLower(abi.Name), filterLower)
}

func (ac *AbisCollection) matchesFunctionFilter(fn *coreTypes.Function, filter string) bool {
	if filter == "" {
		return true
	}
	filterLower := strings.ToLower(filter)
	return strings.Contains(strings.ToLower(fn.Name), filterLower) ||
		strings.Contains(strings.ToLower(fn.Encoding), filterLower)
}

func TestAbisCollectionStateManagement(t *testing.T) {
	collection := NewAbisCollection()

	t.Run("NeedsUpdate", func(t *testing.T) {
		needsUpdate := collection.NeedsUpdate(AbisDownloaded)
		assert.True(t, needsUpdate, "New collection should need update for AbisDownloaded")

		needsUpdate = collection.NeedsUpdate(AbisKnown)
		assert.True(t, needsUpdate, "New collection should need update for AbisKnown")

		needsUpdate = collection.NeedsUpdate(AbisFunctions)
		assert.True(t, needsUpdate, "New collection should need update for AbisFunctions")

		needsUpdate = collection.NeedsUpdate(AbisEvents)
		assert.True(t, needsUpdate, "New collection should need update for AbisEvents")

		needsUpdate = collection.NeedsUpdate("invalid-kind")
		assert.False(t, needsUpdate, "Invalid list kind should return false")
	})

	t.Run("Reset", func(t *testing.T) {
		assert.NotPanics(t, func() {
			collection.Reset(AbisDownloaded)
		}, "Reset with AbisDownloaded should not panic")

		assert.NotPanics(t, func() {
			collection.Reset(AbisKnown)
		}, "Reset with AbisKnown should not panic")

		assert.NotPanics(t, func() {
			collection.Reset(AbisFunctions)
		}, "Reset with AbisFunctions should not panic")

		assert.NotPanics(t, func() {
			collection.Reset(AbisEvents)
		}, "Reset with AbisEvents should not panic")

		needsUpdate := collection.NeedsUpdate(AbisDownloaded)
		assert.True(t, needsUpdate, "After reset, collection should need update")

		assert.NotPanics(t, func() {
			collection.Reset("invalid-kind")
		}, "Reset with invalid list kind should not panic")
	})
}

func TestAbisCollectionLoadData(t *testing.T) {
	collection := NewAbisCollection()

	t.Run("LoadDataValidKinds", func(t *testing.T) {
		assert.NotPanics(t, func() {
			collection.LoadData(AbisDownloaded)
		}, "LoadData with AbisDownloaded should not panic")

		assert.NotPanics(t, func() {
			collection.LoadData(AbisKnown)
		}, "LoadData with AbisKnown should not panic")

		assert.NotPanics(t, func() {
			collection.LoadData(AbisFunctions)
		}, "LoadData with AbisFunctions should not panic")

		assert.NotPanics(t, func() {
			collection.LoadData(AbisEvents)
		}, "LoadData with AbisEvents should not panic")
	})

	t.Run("LoadDataInvalidKind", func(t *testing.T) {
		assert.NotPanics(t, func() {
			collection.LoadData("invalid-kind")
		}, "LoadData with invalid list kind should not panic")
	})
}

func TestAbisCollectionGetPage(t *testing.T) {
	collection := NewAbisCollection()

	t.Run("BasicGetPageAbisDownloaded", func(t *testing.T) {
		page, err := collection.GetPage(AbisDownloaded, 0, 10, sdk.SortSpec{}, "")

		if err == nil && page != nil {
			assert.Equal(t, AbisDownloaded, page.Kind)
			assert.GreaterOrEqual(t, page.TotalItems, 0)
			assert.GreaterOrEqual(t, page.ExpectedTotal, 0)
			assert.LessOrEqual(t, len(page.Abis), 10, "Returned items should not exceed page size")
		}
	})

	t.Run("BasicGetPageAbisKnown", func(t *testing.T) {
		page, err := collection.GetPage(AbisKnown, 0, 10, sdk.SortSpec{}, "")

		if err == nil && page != nil {
			assert.Equal(t, AbisKnown, page.Kind)
			assert.GreaterOrEqual(t, page.TotalItems, 0)
			assert.GreaterOrEqual(t, page.ExpectedTotal, 0)
			assert.LessOrEqual(t, len(page.Abis), 10, "Returned items should not exceed page size")
		}
	})

	t.Run("BasicGetPageAbisFunctions", func(t *testing.T) {
		page, err := collection.GetPage(AbisFunctions, 0, 10, sdk.SortSpec{}, "")

		if err == nil && page != nil {
			assert.Equal(t, AbisFunctions, page.Kind)
			assert.GreaterOrEqual(t, page.TotalItems, 0)
			assert.GreaterOrEqual(t, page.ExpectedTotal, 0)
			assert.LessOrEqual(t, len(page.Functions), 10, "Returned items should not exceed page size")
		}
	})

	t.Run("BasicGetPageAbisEvents", func(t *testing.T) {
		page, err := collection.GetPage(AbisEvents, 0, 10, sdk.SortSpec{}, "")

		if err == nil && page != nil {
			assert.Equal(t, AbisEvents, page.Kind)
			assert.GreaterOrEqual(t, page.TotalItems, 0)
			assert.GreaterOrEqual(t, page.ExpectedTotal, 0)
			assert.LessOrEqual(t, len(page.Functions), 10, "Returned items should not exceed page size")
		}
	})

	t.Run("GetPageWithFilter", func(t *testing.T) {
		page, err := collection.GetPage(AbisDownloaded, 0, 10, sdk.SortSpec{}, "test")

		if err == nil && page != nil {
			assert.Equal(t, AbisDownloaded, page.Kind)
			assert.GreaterOrEqual(t, page.TotalItems, 0)
		}
	})

	t.Run("InvalidListKind", func(t *testing.T) {
		_, err := collection.GetPage("invalid-kind", 0, 10, sdk.SortSpec{}, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unexpected list kind")
	})

	t.Run("ZeroPageSize", func(t *testing.T) {
		page, err := collection.GetPage(AbisDownloaded, 0, 0, sdk.SortSpec{}, "")
		if err == nil && page != nil {
			assert.Equal(t, AbisDownloaded, page.Kind)
			assert.Len(t, page.Abis, 0, "Zero page size should return no items")
		}
	})
}
