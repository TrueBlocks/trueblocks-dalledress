// ADD_ROUTE
package types

import (
	"strings"
	"testing"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
)

// MockApp is a mock implementation of the App interface for testing.
type MockApp struct {
	RegisteredCtxs map[base.Address]*output.RenderCtx
	Events         []struct {
		EventType msgs.EventType
		Payload   interface{}
	}
}

func (m *MockApp) LogBackend(msg string) {
}

func (m *MockApp) RegisterCtx(addr base.Address) *output.RenderCtx {
	if m.RegisteredCtxs == nil {
		m.RegisteredCtxs = make(map[base.Address]*output.RenderCtx)
	}
	renderCtx := output.NewStreamingContext()
	m.RegisteredCtxs[addr] = renderCtx
	return renderCtx
}

func (m *MockApp) Cancel(addr base.Address) (int, bool) {
	if m.RegisteredCtxs == nil {
		return 0, false
	}
	if ctx, exists := m.RegisteredCtxs[addr]; exists {
		ctx.Cancel()
		delete(m.RegisteredCtxs, addr)
		return 1, true
	}
	return 0, false
}

func (m *MockApp) EmitEvent(eventType msgs.EventType, payload interface{}) {
	m.Events = append(m.Events, struct {
		EventType msgs.EventType
		Payload   interface{}
	}{EventType: eventType, Payload: payload})
}

func NewMockApp() *MockApp {
	return &MockApp{
		RegisteredCtxs: make(map[base.Address]*output.RenderCtx),
		Events: make([]struct {
			EventType msgs.EventType
			Payload   interface{}
		}, 0),
	}
}

func TestNewAbisCollection(t *testing.T) {
	mockApp := NewMockApp()
	ac := NewAbisCollection(mockApp)
	if ac.isLoaded { // Changed from ac.loaded to ac.isLoaded
		t.Error("NewAbisCollection should not be fully loaded")
	}
	if ac.isLoading { // Added check for isLoading
		t.Error("NewAbisCollection should not be loading")
	}
	// Check that slices are initialized (non-nil and empty)
	if ac.downloadedAbis == nil || len(ac.downloadedAbis) != 0 {
		t.Errorf("NewAbisCollection.downloadedAbis not initialized correctly: got %v, want empty non-nil slice", ac.downloadedAbis)
	}
	if ac.knownAbis == nil || len(ac.knownAbis) != 0 {
		t.Errorf("NewAbisCollection.knownAbis not initialized correctly: got %v, want empty non-nil slice", ac.knownAbis)
	}
	if ac.allFunctions == nil || len(ac.allFunctions) != 0 {
		t.Errorf("NewAbisCollection.allFunctions not initialized correctly: got %v", ac.allFunctions)
	}
	if ac.allEvents == nil || len(ac.allEvents) != 0 {
		t.Errorf("NewAbisCollection.allEvents not initialized correctly: got %v", ac.allEvents)
	}
}

func TestGetPage(t *testing.T) {
	mockApp := NewMockApp()
	ac := NewAbisCollection(mockApp)
	ac.EnsureInitialLoad() // Changed from ac.Load()

	// Need to wait for loading to complete before testing GetPage,
	// or GetPage should correctly reflect the loading state.
	// The current GetPage implementation calls EnsureInitialLoad itself.

	// Let's simulate a loaded state for GetPage tests for now.
	// This bypasses the async loading for the purpose of testing GetPage's logic.
	// In a real scenario, you'd wait or use a more complex setup.
	ac.mutex.Lock()
	ac.isLoaded = true // Simulate loaded state
	// Populate with some dummy data for testing pagination and filtering
	ac.downloadedAbis = []coreTypes.Abi{
		{Name: "TestABI1", Address: base.HexToAddress("0x1")},   // Use base.HexToAddress
		{Name: "AnotherABI", Address: base.HexToAddress("0x2")}, // Use base.HexToAddress
	}
	ac.knownAbis = []coreTypes.Abi{
		{Name: "KnownABI1", Address: base.HexToAddress("0x3")}, // Use base.HexToAddress
	}
	ac.allFunctions = []coreTypes.Function{
		{Name: "func1", Signature: "func1()"},
		{Name: "func2", Signature: "func2(uint256)"},
	}
	ac.allEvents = []coreTypes.Function{
		{Name: "event1", Signature: "event1()"},
	}
	ac.mutex.Unlock()
	t.Run("DownloadedAbis", func(t *testing.T) {
		if len(ac.downloadedAbis) == 0 {
			// This skip might still be relevant if, after EnsureInitialLoad and waiting,
			// there truly are no downloaded ABIs.
			t.Skip("Skipping DownloadedAbis tests as ac.downloadedAbis is empty")
		}
		nItems := len(ac.downloadedAbis)
		var itemName, itemAddr string
		if nItems > 0 {
			itemName = ac.downloadedAbis[0].Name
			itemAddr = ac.downloadedAbis[0].Address.Hex()
		}

		countMatches := func(filter string) int {
			count := 0
			lFilter := strings.ToLower(filter)
			for _, abi := range ac.downloadedAbis {
				if filter == "" || strings.Contains(strings.ToLower(abi.Name), lFilter) || strings.Contains(strings.ToLower(abi.Address.Hex()), lFilter) {
					count++
				}
			}
			return count
		}

		minPageLen := func(pageSize, actualTotal, first int) int {
			if first >= actualTotal {
				return 0
			}
			if first+pageSize > actualTotal {
				return actualTotal - first
			}
			return pageSize
		}

		tests := []struct {
			name            string
			first           int
			pageSize        int
			filter          string
			expectedTotal   func(filter string) int
			expectedPageLen func(filter string, first, pageSize, actualTotal int) int
		}{
			{"Page1_Size1", 0, 1, "", countMatches, func(f string, first, ps, total int) int { return minPageLen(ps, countMatches(f), first) }},
			{"AllItems", 0, nItems, "", countMatches, func(f string, first, ps, total int) int { return minPageLen(ps, countMatches(f), first) }},
			{"FilterByName_MatchFirst", 0, nItems, itemName, countMatches, func(f string, first, ps, total int) int { return minPageLen(ps, countMatches(f), first) }},
			{"FilterByAddr_MatchFirst", 0, nItems, itemAddr, countMatches, func(f string, first, ps, total int) int { return minPageLen(ps, countMatches(f), first) }},
			{"Filter_NoMatch", 0, nItems, "THIS_SHOULD_NOT_MATCH_ANYTHING_EVER_12345", countMatches, func(f string, first, ps, total int) int { return 0 }},
			{"PageSizeLargerThanTotal", 0, nItems + 5, "", countMatches, func(f string, first, ps, total int) int { return minPageLen(ps, countMatches(f), first) }},
			{"OffsetBeyondTotal", nItems + 1, 5, "", countMatches, func(f string, first, ps, total int) int { return 0 }},
		}

		for _, tc := range tests {
			filter := tc.filter
			if nItems == 0 && (tc.filter == itemName || tc.filter == itemAddr) {
				filter = "FILTER_WHEN_NO_ITEMS_EXIST"
			}
			t.Run(tc.name, func(t *testing.T) {
				page, err := ac.GetPage("Downloaded", tc.first, tc.pageSize, nil, filter)
				if err != nil {
					t.Fatalf("GetPage returned error: %v", err)
				}
				actualNItems := page.TotalItems
				expectedNItems := tc.expectedTotal(filter)

				if actualNItems != expectedNItems {
					t.Errorf("Expected total %d, got %d from page.TotalAbis (filter: '%s')", expectedNItems, actualNItems, filter)
				}

				expectedLen := tc.expectedPageLen(filter, tc.first, tc.pageSize, actualNItems)

				if len(page.Abis) != expectedLen {
					t.Errorf("Expected page length %d, got %d (TotalAbis: %d, first: %d, pageSize: %d, filter: '%s')",
						expectedLen, len(page.Abis), actualNItems, tc.first, tc.pageSize, filter)
				}
			})
		}
	})

	t.Run("KnownAbis", func(t *testing.T) {
		if len(ac.knownAbis) == 0 {
			t.Skip("Skipping KnownAbis tests as ac.knownAbis is empty")
		}
		nItems := len(ac.knownAbis)
		var itemName, itemAddr string
		if nItems > 0 {
			itemName = ac.knownAbis[0].Name
			itemAddr = ac.knownAbis[0].Address.Hex()
		}

		countMatches := func(filter string) int {
			count := 0
			lFilter := strings.ToLower(filter)
			for _, abi := range ac.knownAbis {
				if filter == "" || strings.Contains(strings.ToLower(abi.Name), lFilter) || strings.Contains(strings.ToLower(abi.Address.Hex()), lFilter) {
					count++
				}
			}
			return count
		}

		minPageLen := func(pageSize, actualTotal, first int) int {
			if first >= actualTotal {
				return 0
			}
			if first+pageSize > actualTotal {
				return actualTotal - first
			}
			return pageSize
		}

		tests := []struct {
			name            string
			first           int
			pageSize        int
			filter          string
			expectedTotal   func(filter string) int
			expectedPageLen func(filter string, first, pageSize, actualTotal int) int
		}{
			{"Page1_Size1", 0, 1, "", countMatches, func(f string, first, ps, total int) int { return minPageLen(ps, countMatches(f), first) }},
			{"AllItems", 0, nItems, "", countMatches, func(f string, first, ps, total int) int { return minPageLen(ps, countMatches(f), first) }},
			{"FilterByName_MatchFirst", 0, nItems, itemName, countMatches, func(f string, first, ps, total int) int { return minPageLen(ps, countMatches(f), first) }},
			{"FilterByAddr_MatchFirst", 0, nItems, itemAddr, countMatches, func(f string, first, ps, total int) int { return minPageLen(ps, countMatches(f), first) }},
			{"Filter_NoMatch", 0, nItems, "THIS_SHOULD_NOT_MATCH_ANYTHING_EVER_12345", countMatches, func(f string, first, ps, total int) int { return 0 }},
			{"PageSizeLargerThanTotal", 0, nItems + 5, "", countMatches, func(f string, first, ps, total int) int { return minPageLen(ps, countMatches(f), first) }},
			{"OffsetBeyondTotal", nItems + 1, 5, "", countMatches, func(f string, first, ps, total int) int { return 0 }},
		}

		for _, tc := range tests {
			filter := tc.filter
			if nItems == 0 && (tc.filter == itemName || tc.filter == itemAddr) {
				filter = "FILTER_WHEN_NO_ITEMS_EXIST"
			}
			t.Run(tc.name, func(t *testing.T) {
				page, err := ac.GetPage("Known", tc.first, tc.pageSize, nil, filter)
				if err != nil {
					t.Fatalf("GetPage returned error: %v", err)
				}
				actualNItems := page.TotalItems
				expectedNItems := tc.expectedTotal(filter)

				if actualNItems != expectedNItems {
					t.Errorf("Expected total %d, got %d from page.TotalAbis (filter: '%s')", expectedNItems, actualNItems, filter)
				}

				expectedLen := tc.expectedPageLen(filter, tc.first, tc.pageSize, actualNItems)

				if len(page.Abis) != expectedLen {
					t.Errorf("Expected page length %d, got %d (TotalAbis: %d, first: %d, pageSize: %d, filter: '%s')",
						expectedLen, len(page.Abis), actualNItems, tc.first, tc.pageSize, filter)
				}
			})
		}
	})

	t.Run("Functions", func(t *testing.T) {
		if len(ac.allFunctions) == 0 {
			t.Skip("Skipping Functions tests: no functions found.")
			return
		}

		nItems := len(ac.allFunctions)
		var itemName, itemSig string
		if nItems > 0 {
			itemName = ac.allFunctions[0].Name
			itemSig = ac.allFunctions[0].Signature
		}

		countMatches := func(filter string) int {
			if filter == "" {
				return nItems
			}
			count := 0
			lFilter := strings.ToLower(filter)
			for _, f := range ac.allFunctions {
				if strings.Contains(strings.ToLower(f.Name), lFilter) ||
					strings.Contains(strings.ToLower(f.Signature), lFilter) {
					count++
				}
			}
			return count
		}

		minPageLen := func(pageSize, actualTotal, first int) int {
			if first >= actualTotal {
				return 0
			}
			return min(pageSize, actualTotal-first)
		}

		tests := []struct {
			name            string
			first           int
			pageSize        int
			filter          string
			expectedTotal   func(filter string) int
			expectedPageLen func(filter string, first, pageSize, actualTotal int) int
		}{
			{"Page1_Size1", 0, 1, "", countMatches, func(f string, first, ps, total int) int { return minPageLen(ps, countMatches(f), first) }},
			{"AllItems", 0, nItems, "", countMatches, func(f string, first, ps, total int) int { return minPageLen(ps, countMatches(f), first) }},
			{"FilterByName_MatchFirst", 0, nItems, itemName, countMatches, func(f string, first, ps, total int) int { return minPageLen(ps, countMatches(f), first) }},
			{"FilterBySig_MatchFirst", 0, nItems, itemSig, countMatches, func(f string, first, ps, total int) int { return minPageLen(ps, countMatches(f), first) }},
			{"Filter_NoMatch", 0, nItems, "THIS_SHOULD_NOT_MATCH_ANYTHING_EVER_12345", countMatches, func(f string, first, ps, total int) int { return 0 }},
		}

		for _, tc := range tests {
			filter := tc.filter
			if nItems == 0 && (tc.filter == itemName || tc.filter == itemSig) {
				filter = "FILTER_WHEN_NO_FUNCS_EXIST"
			}
			t.Run(tc.name, func(t *testing.T) {
				// Assuming sortDef is nil for these tests as in original
				page, err := ac.GetPage("Functions", tc.first, tc.pageSize, nil, filter)
				if err != nil {
					t.Fatalf("GetPage returned error: %v", err)
				}
				actualNItems := page.TotalItems
				expectedNItems := tc.expectedTotal(filter)
				if actualNItems != expectedNItems {
					t.Errorf("Expected total %d, got %d (filter: '%s')", expectedNItems, actualNItems, filter)
				}
				expectedLen := tc.expectedPageLen(filter, tc.first, tc.pageSize, actualNItems)
				if len(page.Functions) != expectedLen {
					t.Errorf("Expected page length %d, got %d (TotalFunctions: %d, first: %d, pageSize: %d, filter: '%s')",
						expectedLen, len(page.Functions), actualNItems, tc.first, tc.pageSize, filter)
				}
			})
		}
	})

	t.Run("Events", func(t *testing.T) {
		if len(ac.allEvents) == 0 {
			t.Skip("Skipping Events tests: no events found.")
			return
		}

		nItems := len(ac.allEvents)
		var itemName, itemSig string
		if nItems > 0 {
			itemName = ac.allEvents[0].Name
			itemSig = ac.allEvents[0].Signature
		}

		countMatches := func(filter string) int {
			if filter == "" {
				return nItems
			}
			count := 0
			lFilter := strings.ToLower(filter)
			for _, e := range ac.allEvents {
				if strings.Contains(strings.ToLower(e.Name), lFilter) || strings.Contains(strings.ToLower(e.Signature), lFilter) {
					count++
				}
			}
			return count
		}

		minPageLen := func(pageSize, actualTotal, first int) int {
			if first >= actualTotal {
				return 0
			}
			return min(pageSize, actualTotal-first)
		}

		tests := []struct {
			name            string
			first           int
			pageSize        int
			filter          string
			expectedTotal   func(filter string) int
			expectedPageLen func(filter string, first, pageSize, actualTotal int) int
		}{
			{"Page1_Size1", 0, 1, "", countMatches, func(f string, first, ps, total int) int { return minPageLen(ps, countMatches(f), first) }},
			{"AllItems", 0, nItems, "", countMatches, func(f string, first, ps, total int) int { return minPageLen(ps, countMatches(f), first) }},
			{"FilterByName_MatchFirst", 0, nItems, itemName, countMatches, func(f string, first, ps, total int) int { return minPageLen(ps, countMatches(f), first) }},
			{"FilterBySig_MatchFirst", 0, nItems, itemSig, countMatches, func(f string, first, ps, total int) int { return minPageLen(ps, countMatches(f), first) }},
			{"Filter_NoMatch", 0, nItems, "THIS_SHOULD_NOT_MATCH_ANYTHING_EVER_12345", countMatches, func(f string, first, ps, total int) int { return 0 }},
		}

		for _, tc := range tests {
			filter := tc.filter
			if nItems == 0 && (tc.filter == itemName || tc.filter == itemSig) {
				filter = "FILTER_WHEN_NO_EVENTS_EXIST"
			}
			t.Run(tc.name, func(t *testing.T) {
				page, err := ac.GetPage("Events", tc.first, tc.pageSize, nil, filter)
				if err != nil {
					t.Fatalf("GetPage returned error: %v", err)
				}
				actualNItems := page.TotalItems
				expectedNItems := tc.expectedTotal(filter)
				if actualNItems != expectedNItems {
					t.Errorf("Expected total %d, got %d (filter: '%s')", expectedNItems, actualNItems, filter)
				}
				expectedLen := tc.expectedPageLen(filter, tc.first, tc.pageSize, actualNItems)
				if len(page.Functions) != expectedLen {
					t.Errorf("Expected page length %d, got %d (TotalEvents: %d, first: %d, pageSize: %d, filter: '%s')",
						expectedLen, len(page.Functions), actualNItems, tc.first, tc.pageSize, filter)
				}
			})
		}
	})

	t.Run("InvalidKind", func(t *testing.T) {
		_, err := ac.GetPage("UnknownKind", 0, 1, nil, "")
		if err == nil {
			t.Error("Expected error for invalid kind, got nil")
		} else if !strings.Contains(strings.ToLower(err.Error()), "unknown abi page kind") {
			t.Errorf("Expected error message for 'unknown abi page kind', got: %v", err)
		}
	})
}

// ADD_ROUTE
