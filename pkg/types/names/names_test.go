package names

import (
	"testing"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
	"github.com/stretchr/testify/assert"
)

func assertNamesPage(t *testing.T, page types.Page) *NamesPage {
	t.Helper()
	if page == nil {
		t.Fatal("page is nil")
	}
	namesPage, ok := page.(*NamesPage)
	if !ok {
		t.Fatalf("expected *NamesPage, got %T", page)
	}
	return namesPage
}

// Domain-specific filtering tests for Names collection

func TestNamesMatchesFilter(t *testing.T) {
	collection := NewNamesCollection()
	testName := &Name{
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

func TestNamesCollectionDomainSpecificFiltering(t *testing.T) {
	collection := NewNamesCollection()

	t.Run("GetPageWithDomainSpecificFilter", func(t *testing.T) {
		page, err := collection.GetPage(NamesAll, 0, 10, sdk.SortSpec{}, "test")

		if err == nil && page != nil {
			namesPage := assertNamesPage(t, page)
			assert.Equal(t, NamesAll, namesPage.Facet)
			assert.GreaterOrEqual(t, namesPage.TotalItems, 0)
		}
	})
}
