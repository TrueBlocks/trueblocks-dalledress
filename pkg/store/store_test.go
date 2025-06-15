package store

import (
	"fmt"
	"testing"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
)

func TestStoreDirectly(t *testing.T) {
	testStore := NewStore(
		"direct-test-store",
		func(ctx *output.RenderCtx) error {
			fmt.Println("Direct test: Query function starting...")

			testItem := &TestItem{ID: 1, Name: "Direct", Value: 100}

			fmt.Println("Direct test: Sending item...")
			ctx.ModelChan <- testItem

			fmt.Println("Direct test: Closing ModelChan...")
			close(ctx.ModelChan)

			fmt.Println("Direct test: Closing ErrorChan...")
			close(ctx.ErrorChan)

			fmt.Println("Direct test: Query function completed")
			return nil
		},
		func(itemIntf interface{}) *TestItem {
			fmt.Println("Direct test: Processing item...")
			if item, ok := itemIntf.(*TestItem); ok {
				fmt.Println("Direct test: Item processed successfully")
				return item
			}
			fmt.Println("Direct test: Failed to process item")
			return nil
		},
		nil,
	)

	fmt.Println("Direct test: Initial store state:", testStore.GetState())

	fmt.Println("Direct test: Starting fetch...")
	err := testStore.Fetch()
	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}

	fmt.Println("Direct test: Fetch completed")
	fmt.Println("Direct test: Final store state:", testStore.GetState())
	fmt.Println("Direct test: Items count:", len(testStore.GetItems()))

	if testStore.GetState() != StateLoaded {
		t.Errorf("Expected store state to be StateLoaded (%d), got %d", StateLoaded, testStore.GetState())
	}

	if len(testStore.GetItems()) != 1 {
		t.Errorf("Expected 1 item, got %d", len(testStore.GetItems()))
	}
}

type TestItem struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func (t *TestItem) Model(chain, format string, verbose bool, extraOptions map[string]any) types.Model {
	return types.Model{
		Data: map[string]any{
			"id":    t.ID,
			"name":  t.Name,
			"value": t.Value,
		},
		Order: []string{"id", "name", "value"},
	}
}
