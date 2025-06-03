package repository_test

import (
	"testing"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/mocks"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/repository"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

func TestRepositoryPattern(t *testing.T) {
	app := mocks.NewMockApp()

	// Test creating a simple repository
	filterFunc := func(item *string) bool {
		return len(*item) > 3
	}

	processFunc := func(itemIntf interface{}) *string {
		if str, ok := itemIntf.(*string); ok {
			return str
		}
		return nil
	}

	queryFunc := func(renderCtx *output.RenderCtx) {
		// Mock query function that would normally call SDK
		// In real usage this would populate renderCtx.ModelChan
	}

	// Simple deduplication function
	dedupeFunc := func(existing []string, newItem *string) bool {
		if newItem == nil {
			return false
		}
		for _, item := range existing {
			if item == *newItem {
				return true // Already exists
			}
		}
		return false
	}

	repo := repository.NewBaseRepository(
		app,
		types.ListKind("test"),
		filterFunc,
		processFunc,
		queryFunc,
		dedupeFunc,
	)

	// Test initial state
	if repo.IsLoaded() {
		t.Error("Repository should not be loaded initially")
	}

	if !repo.NeedsUpdate() {
		t.Error("Repository should need update initially")
	}

	if repo.Count() != 0 {
		t.Errorf("Expected count 0, got %d", repo.Count())
	}

	// Test loading (this will succeed in mock but with no data)
	opts := repository.LoadOptions{}
	result, err := repo.Load(opts)
	if err != nil {
		t.Errorf("Unexpected error with mock implementation: %v", err)
	}
	if result == nil {
		t.Error("Expected non-nil result from Load")
	}
	// if result != nil && len(result.Data) != 0 {
	// 	t.Errorf("Expected 0 items in mock result, got %d", len(result.Data))
	// }

	t.Log("Repository pattern test completed successfully")
}
