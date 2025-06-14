// NAMES_ROUTE
package names

import (
	"testing"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/sorting"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

var noSort = sorting.EmptySortSpec()

func TestNewNamesCollection(t *testing.T) {
	names := NewNamesCollection()
	if names == nil {
		t.Errorf("NewNamesCollection() returned nil")
		return
	}

	if names.allFacet == nil {
		t.Errorf("allFacet not initialized")
	}
	if names.customFacet == nil {
		t.Errorf("customFacet not initialized")
	}
	if names.prefundFacet == nil {
		t.Errorf("prefundFacet not initialized")
	}
	if names.regularFacet == nil {
		t.Errorf("regularFacet not initialized")
	}
	if names.baddressFacet == nil {
		t.Errorf("baddressFacet not initialized")
	}
}

func TestNamesCollectionMethods(t *testing.T) {
	names := NewNamesCollection()

	if !names.NeedsUpdate(NamesAll) {
		t.Errorf("Expected NeedsUpdate to be true initially")
	}

	if expected := names.getExpectedTotal(NamesAll); expected != 0 {
		t.Errorf("Expected getExpectedTotal to be 0 initially, got %d", expected)
	}
}

func TestNamesPageStructure(t *testing.T) {
	names := NewNamesCollection()
	if names == nil {
		t.Errorf("NewNamesCollection() returned nil")
		return
	}

	page, err := names.GetPage(NamesAll, 0, 10, noSort, "")
	if err != nil {
		t.Errorf("GetPage should not error with empty data: %v", err)
	} else if page == nil {
		t.Errorf("GetPage should not return nil page")
		return
	}
	if page.Kind != NamesAll {
		t.Errorf("Expected page kind to be %s, got %s", NamesAll, page.Kind)
	}
	if page.Names == nil {
		t.Errorf("Expected Names array to be initialized")
	}
	if len(page.Names) != 0 {
		t.Errorf("Expected empty Names array, got %d items", len(page.Names))
	}
}

func TestFindNameByAddress(t *testing.T) {
	names := NewNamesCollection()

	addr := base.HexToAddress("0x123")
	name, found := names.FindNameByAddress(addr)
	if found {
		t.Errorf("Expected FindNameByAddress to return false for non-existent address")
	}
	if name != nil {
		t.Errorf("Expected FindNameByAddress to return nil name for non-existent address")
	}
}

func TestGetPageWithDifferentListKinds(t *testing.T) {
	names := NewNamesCollection()

	listKinds := []types.ListKind{
		NamesAll,
		NamesCustom,
		NamesPrefund,
		NamesRegular,
		NamesBaddress,
	}

	for _, kind := range listKinds {
		t.Run("ListKind_"+string(kind), func(t *testing.T) {
			page, err := names.GetPage(NamesAll, 0, 10, noSort, "")
			if err != nil {
				t.Errorf("GetPage should not error for kind %s: %v", kind, err)
			}
			if page == nil {
				t.Errorf("GetPage should not return nil page for kind %s", kind)
			}
		})
	}
}

func TestGetPageWithFilter(t *testing.T) {
	names := NewNamesCollection()

	filters := []string{"", "test", "0x123", "alice"}

	for _, filter := range filters {
		t.Run("Filter_"+filter, func(t *testing.T) {
			page, err := names.GetPage(NamesAll, 0, 10, noSort, filter)
			if err != nil {
				t.Errorf("GetPage should not error with filter '%s': %v", filter, err)
			}
			if page == nil {
				t.Errorf("GetPage should not return nil page with filter '%s'", filter)
			}
		})
	}
}

func TestGetPagePagination(t *testing.T) {
	names := NewNamesCollection()

	testCases := []struct {
		first    int
		pageSize int
	}{
		{0, 10},
		{0, 5},
		{5, 10},
		{10, 5},
	}

	for _, tc := range testCases {
		t.Run("Pagination", func(t *testing.T) {
			page, err := names.GetPage(NamesAll, tc.first, tc.pageSize, noSort, "")
			if err != nil {
				t.Errorf("GetPage should not error with first=%d, pageSize=%d: %v", tc.first, tc.pageSize, err)
			}
			if page == nil {
				t.Errorf("GetPage should not return nil page with first=%d, pageSize=%d", tc.first, tc.pageSize)
			}
		})
	}
}

// NAMES_ROUTE
