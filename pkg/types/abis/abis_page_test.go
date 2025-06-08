package abis

import (
	"testing"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/sorting"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

func TestGetPage_BasicFunctionality(t *testing.T) {
	ac := NewAbisCollection()

	tests := []struct {
		name     string
		listKind types.ListKind
	}{
		{"Downloaded", AbisDownloaded},
		{"Known", AbisKnown},
		{"Functions", AbisFunctions},
		{"Events", AbisEvents},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			page, err := ac.GetPage(tc.listKind, 0, 10, sorting.EmptySortSpec(), "")
			if err != nil {
				t.Fatalf("GetPage returned error: %v", err)
			}

			if page.Kind != tc.listKind {
				t.Errorf("Expected page kind %s, got %s", tc.listKind, page.Kind)
			}

			if page.TotalItems != 0 {
				t.Errorf("Expected 0 total items for empty facet, got %d", page.TotalItems)
			}

			if len(page.Abis) != 0 && len(page.Functions) != 0 {
				t.Errorf("Expected empty results for empty facet")
			}
		})
	}
}

func TestGetPage_InvalidListKind(t *testing.T) {
	ac := NewAbisCollection()

	_, err := ac.GetPage(types.ListKind("InvalidKind"), 0, 10, sorting.EmptySortSpec(), "")
	if err == nil {
		t.Error("Expected error for invalid list kind, got nil")
	}
}

