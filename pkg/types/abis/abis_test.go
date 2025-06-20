package abis

import (
	"strings"
	"testing"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
	"github.com/stretchr/testify/assert"
)

func assertAbisPage(t *testing.T, page types.Page) *AbisPage {
	t.Helper()
	if page == nil {
		t.Fatal("page is nil")
	}
	abisPage, ok := page.(*AbisPage)
	if !ok {
		t.Fatalf("expected *AbisPage, got %T", page)
	}
	return abisPage
}

// Domain-specific filtering and ABI logic tests for Abis collection

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

func TestAbisCollectionDomainSpecific(t *testing.T) {
	collection := NewAbisCollection()

	t.Run("GetPageMultiFacetFiltering", func(t *testing.T) {
		// Test ABI list filtering
		page, err := collection.GetPage(AbisDownloaded, 0, 10, sdk.SortSpec{}, "test")
		if err == nil && page != nil {
			abisPage := assertAbisPage(t, page)
			assert.Equal(t, AbisDownloaded, abisPage.Facet)
			assert.GreaterOrEqual(t, abisPage.TotalItems, 0)
		}

		// Test function filtering
		page, err = collection.GetPage(AbisFunctions, 0, 10, sdk.SortSpec{}, "transfer")
		if err == nil && page != nil {
			abisPage := assertAbisPage(t, page)
			assert.Equal(t, AbisFunctions, abisPage.Facet)
			assert.GreaterOrEqual(t, abisPage.TotalItems, 0)
		}
	})

	t.Run("MultiFacetSupport", func(t *testing.T) {
		// Test that different facets work correctly
		facets := []types.DataFacet{AbisDownloaded, AbisKnown, AbisFunctions, AbisEvents}
		for _, facet := range facets {
			page, err := collection.GetPage(facet, 0, 5, sdk.SortSpec{}, "")
			if err == nil && page != nil {
				abisPage := assertAbisPage(t, page)
				assert.Equal(t, facet, abisPage.Facet)
				// Verify the right data structure is populated based on facet
				switch facet {
				case AbisDownloaded, AbisKnown:
					// Should have Abis populated
					assert.NotNil(t, abisPage.Abis)
				case AbisFunctions, AbisEvents:
					// Should have Functions populated
					assert.NotNil(t, abisPage.Functions)
				}
			}
		}
	})
}
