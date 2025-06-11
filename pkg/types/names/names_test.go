// NAMES_ROUTE
package names

import (
	"testing"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/sorting"
)

func makeTestNames() NamesCollection {
	return NamesCollection{
		Map: map[base.Address]types.Name{
			base.HexToAddress("0x1"): {Address: base.HexToAddress("0x1"), Tags: "custom", Parts: types.Custom},
			base.HexToAddress("0x2"): {Address: base.HexToAddress("0x2"), Tags: "prefund", Parts: types.Prefund},
			base.HexToAddress("0x3"): {Address: base.HexToAddress("0x3"), Tags: "regular", Parts: types.Regular},
			base.HexToAddress("0x4"): {Address: base.HexToAddress("0x4"), Tags: "baddress", Parts: types.Baddress},
		},
		List: []*types.Name{
			{Name: "A", Address: base.HexToAddress("0x1"), Tags: "custom", Parts: types.Custom},
			{Name: "B", Address: base.HexToAddress("0x2"), Tags: "prefund", Parts: types.Prefund},
			{Name: "C", Address: base.HexToAddress("0x3"), Tags: "regular", Parts: types.Regular},
			{Name: "D", Address: base.HexToAddress("0x4"), Tags: "baddress", Parts: types.Baddress},
		},
		Custom:   []*types.Name{{Name: "A", Address: base.HexToAddress("0x1"), Tags: "custom", Parts: types.Custom}},
		Prefund:  []*types.Name{{Name: "B", Address: base.HexToAddress("0x2"), Tags: "prefund", Parts: types.Prefund}},
		Regular:  []*types.Name{{Name: "C", Address: base.HexToAddress("0x3"), Tags: "regular", Parts: types.Regular}},
		Baddress: []*types.Name{{Name: "D", Address: base.HexToAddress("0x4"), Tags: "baddress", Parts: types.Baddress}},
	}
}

var noSort = sorting.EmptySortSpec()

func TestNamesStructFields(t *testing.T) {
	names := NamesCollection{
		Map:  make(map[base.Address]types.Name),
		List: []*types.Name{},
	}
	addr := base.HexToAddress("0x123")
	name := types.Name{
		Address: addr,
		Tags:    "test",
		Parts:   types.Custom,
	}
	names.Map[addr] = name
	names.List = append(names.List, &name)

	if len(names.Map) != 1 {
		t.Errorf("expected Map length 1, got %d", len(names.Map))
	}
	if len(names.List) != 1 {
		t.Errorf("expected List length 1, got %d", len(names.List))
	}
	if names.List[0].Address != addr {
		t.Errorf("expected address %s, got %s", addr.Hex(), names.List[0].Address.Hex())
	}
}

func TestGetNamesPagination(t *testing.T) {
	names := NamesCollection{
		List: []*types.Name{
			{Tags: "a"}, {Tags: "b"}, {Tags: "c"}, {Tags: "d"}, {Tags: "e"},
		},
	}
	page, _ := names.GetPage("all", 1, 2, noSort, "")
	if page.Total != 5 {
		t.Errorf("expected total 5, got %d", page.Total)
	}
	if len(page.Names) != 2 {
		t.Errorf("expected 2 names, got %d", len(page.Names))
	}
	if page.Names[0].Tags != "b" || page.Names[1].Tags != "c" {
		t.Errorf("unexpected pagination result: %+v", page.Names)
	}
}

func TestLoadDataNoReloadIfCountsMatch(t *testing.T) {
	names := NamesCollection{
		List: []*types.Name{
			{Parts: types.Custom},
			{Parts: types.Custom},
		},
	}
	// Pass nil for WaitGroup to avoid negative counter panic.
	err := names.LoadData(nil)
	if err != nil && err.Error() != "no names found" {
		t.Errorf("expected nil or 'no names found' error, got %v", err)
	}
}

func TestGetNames_All(t *testing.T) {
	names := makeTestNames()
	page, _ := names.GetPage("all", 0, 10, noSort, "")
	if page.Total != 4 {
		t.Errorf("expected total 4, got %d", page.Total)
	}
	if len(page.Names) != 4 {
		t.Errorf("expected 4 results, got %d", len(page.Names))
	}
}

func TestGetNames_Custom(t *testing.T) {
	names := makeTestNames()
	page, _ := names.GetPage("custom", 0, 10, noSort, "")
	if page.Total != 1 {
		t.Errorf("expected total 1, got %d", page.Total)
	}
	if len(page.Names) != 1 || page.Names[0].Tags != "custom" {
		t.Errorf("expected custom name, got %+v", page.Names)
	}
}

func TestGetNames_Paging(t *testing.T) {
	names := makeTestNames()
	page, _ := names.GetPage("all", 2, 2, noSort, "")
	if page.Total != 4 {
		t.Errorf("expected total 4, got %d", page.Total)
	}
	if len(page.Names) != 2 || page.Names[0].Tags != "regular" || page.Names[1].Tags != "baddress" {
		t.Errorf("unexpected page results: %+v", page.Names)
	}
}

func TestGetNames_OutOfRange(t *testing.T) {
	names := makeTestNames()
	page, _ := names.GetPage("all", 10, 2, noSort, "")
	if page.Total != 4 {
		t.Errorf("expected total 4, got %d", page.Total)
	}
	if len(page.Names) != 0 {
		t.Errorf("expected empty page, got %+v", page.Names)
	}
}

func TestGetNames_Filtering(t *testing.T) {
	names := NamesCollection{
		List: []*types.Name{
			{Name: "Alice", Address: base.HexToAddress("0x1"), Tags: "custom", Source: "manual", Parts: types.Custom},
			{Name: "Bob", Address: base.HexToAddress("0x2"), Tags: "prefund", Source: "imported", Parts: types.Prefund},
			{Name: "Charlie", Address: base.HexToAddress("0x3"), Tags: "regular", Source: "auto", Parts: types.Regular},
			{Name: "Dave", Address: base.HexToAddress("0x4"), Tags: "baddress", Source: "manual", Parts: types.Baddress},
		},
	}

	// Filter by name
	page, _ := names.GetPage("all", 0, 10, noSort, "ali")
	if len(page.Names) != 1 || page.Names[0].Name != "Alice" {
		t.Errorf("expected Alice, got %+v", page.Names)
	}

	// Filter by address
	page, _ = names.GetPage("all", 0, 10, noSort, "0x2")
	if len(page.Names) != 1 || page.Names[0].Name != "Bob" {
		t.Errorf("expected Bob, got %+v", page.Names)
	}

	// Filter by tags
	page, _ = names.GetPage("all", 0, 10, noSort, "regular")
	if len(page.Names) != 1 || page.Names[0].Name != "Charlie" {
		t.Errorf("expected Charlie, got %+v", page.Names)
	}

	// Filter by source
	page, _ = names.GetPage("all", 0, 10, noSort, "manual")
	if len(page.Names) != 2 {
		t.Errorf("expected 2 manual, got %+v", page.Names)
	}

	// Case-insensitive
	page, _ = names.GetPage("all", 0, 10, noSort, "ALICE")
	if len(page.Names) != 1 || page.Names[0].Name != "Alice" {
		t.Errorf("expected Alice (case-insensitive), got %+v", page.Names)
	}

	// Empty filter returns all
	page, _ = names.GetPage("all", 0, 10, noSort, "")
	if len(page.Names) != 4 {
		t.Errorf("expected all names, got %+v", page.Names)
	}
}

// NAMES_ROUTE
