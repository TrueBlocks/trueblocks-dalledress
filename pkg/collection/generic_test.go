package collection

import (
	"testing"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
	"github.com/stretchr/testify/assert"
)

func TestNewGenericCollection(t *testing.T) {
	gc := NewGenericCollection("test-collection")

	assert.NotNil(t, gc)
	assert.Equal(t, "test-collection", gc.GetCollectionName())
	assert.Empty(t, gc.GetSupportedKinds())
}

func TestGenericCollectionBasicInterface(t *testing.T) {
	gc := NewGenericCollection("test-collection")

	// Test interface compliance
	var _ types.Collection = gc

	// Test with non-existent kind
	_, err := gc.GetPage("non-existent", 0, 10, sdk.SortSpec{}, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported list kind")

	// Test CRUD with non-existent kind
	err = gc.Crud("non-existent", crud.Create, "test-item")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported list kind")

	// Test other methods with non-existent kind
	assert.False(t, gc.NeedsUpdate("non-existent"))
	assert.Empty(t, gc.GetStoreForKind("non-existent"))

	// LoadData and Reset should not panic with non-existent kinds
	assert.NotPanics(t, func() {
		gc.LoadData("non-existent")
		gc.Reset("non-existent")
	})
}

func TestGenericCollectionThreadSafety(t *testing.T) {
	gc := NewGenericCollection("test-collection")

	// Test that concurrent access doesn't panic
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()

			// Concurrent reads
			gc.GetSupportedKinds()
			gc.GetCollectionName()
			gc.NeedsUpdate("test-kind")
			gc.GetStoreForKind("test-kind")

			// These operations should be safe even with non-existent kinds
			gc.LoadData("test-kind")
			gc.Reset("test-kind")
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

// Test basic registration workflow (without actual facets for now)
func TestGenericCollectionRegistration(t *testing.T) {
	gc := NewGenericCollection("test-collection")

	// Initially empty
	assert.Empty(t, gc.GetSupportedKinds())

	// TODO: Add tests with actual handler registration once we have
	// concrete implementations working with the existing collections
}
