package facets_test

import (
	"testing"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/facets"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

// mockSource is a test implementation of Source[string]
type mockSource struct {
	testData   []string
	sourceType string
}

func (m *mockSource) Fetch(ctx *output.RenderCtx, processor func(item *string) bool) error {
	for _, data := range m.testData {
		item := data
		if !processor(&item) {
			break // Stop if processor returns false
		}
	}
	return nil
}

func (m *mockSource) GetSourceType() string {
	return m.sourceType
}

func TestFacetPattern(t *testing.T) {
	testData := []string{"ab", "abcd", "abcde", "xyz", "hello"}

	filterFunc := func(item *string) bool {
		return len(*item) > 3
	}

	isDupFunc := func(existing []string, newItem *string) bool { return false }

	mockSrc := &mockSource{
		testData:   testData,
		sourceType: "mock",
	}

	facet := facets.NewBaseFacet(
		types.ListKind("test"),
		filterFunc,
		isDupFunc,
		mockSrc,
	)

	if facet.IsLoaded() {
		t.Error("Facet should not be loaded initially")
	}

	if !facet.NeedsUpdate() {
		t.Error("Facet should need update initially")
	}

	if facet.Count() != 0 {
		t.Errorf("Expected count 0, got %d", facet.Count())
	}

	result, err := facet.Load()
	if err != nil {
		t.Errorf("Unexpected error with mock implementation: %v", err)
	}
	if result == nil {
		t.Error("Expected non-nil result from Load")
	}

	// Verify state after loading
	if !facet.IsLoaded() {
		t.Error("Facet should be loaded after Load() call")
	}

	// Check filtered results - should have 3 items: "abcd", "abcde", "hello"
	expectedCnt := 3
	if facet.Count() != expectedCnt {
		t.Errorf("Expected count %d after filtering, got %d", expectedCnt, facet.Count())
	}

	// Test paging
	page, err := facet.GetPage(0, 10, nil, nil, nil)
	if err != nil {
		t.Errorf("Unexpected error getting page: %v", err)
	}
	if page == nil {
		t.Error("Expected non-nil page result")
	} else {
		if len(page.Items) != expectedCnt {
			t.Errorf("Expected %d items in page, got %d", expectedCnt, len(page.Items))
		}
	}

	t.Log("Facet pattern test completed successfully")
}
