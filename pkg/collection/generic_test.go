package collection

import (
	"testing"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/abis"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/monitors"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/names"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
	"github.com/stretchr/testify/assert"
)

// CollectionTestSuite represents a generic test suite that can be applied to any Collection implementation
type CollectionTestSuite struct {
	collection        types.Collection
	supportedKinds    []types.ListKind
	invalidKinds      []types.ListKind
	collectionName    string
	expectedStoreName string
}

// CreateCollectionTestSuite creates a test suite for any Collection implementation
func CreateCollectionTestSuite(
	collection types.Collection,
	supportedKinds []types.ListKind,
	invalidKinds []types.ListKind,
	collectionName string,
	expectedStoreName string,
) *CollectionTestSuite {
	return &CollectionTestSuite{
		collection:        collection,
		supportedKinds:    supportedKinds,
		invalidKinds:      invalidKinds,
		collectionName:    collectionName,
		expectedStoreName: expectedStoreName,
	}
}

// TestCollectionInterface tests basic Collection interface compliance
func (suite *CollectionTestSuite) TestCollectionInterface(t *testing.T) {
	t.Run("GetCollectionName", func(t *testing.T) {
		name := suite.collection.GetCollectionName()
		assert.Equal(t, suite.collectionName, name, "Collection name should match expected")
	})

	t.Run("GetSupportedKinds", func(t *testing.T) {
		kinds := suite.collection.GetSupportedKinds()
		assert.NotEmpty(t, kinds, "Supported kinds should not be empty")
		assert.ElementsMatch(t, suite.supportedKinds, kinds, "Supported kinds should match expected")
	})

	t.Run("GetStoreForKind_ValidKinds", func(t *testing.T) {
		for _, kind := range suite.supportedKinds {
			storeName := suite.collection.GetStoreForKind(kind)
			assert.NotEmpty(t, storeName, "Store name should not be empty for valid kind %s", kind)
		}
	})

	t.Run("GetStoreForKind_InvalidKinds", func(t *testing.T) {
		for _, kind := range suite.invalidKinds {
			storeName := suite.collection.GetStoreForKind(kind)
			assert.Empty(t, storeName, "Store name should be empty for invalid kind %s", kind)
		}
	})
}

// TestCollectionStateManagement tests NeedsUpdate/Reset/LoadData behavior
func (suite *CollectionTestSuite) TestCollectionStateManagement(t *testing.T) {
	t.Run("NeedsUpdate_ValidKinds", func(t *testing.T) {
		for _, kind := range suite.supportedKinds {
			needsUpdate := suite.collection.NeedsUpdate(kind)
			// We can't assert true/false here as it depends on actual state
			// But the call should not panic
			_ = needsUpdate
		}
	})

	t.Run("NeedsUpdate_InvalidKinds", func(t *testing.T) {
		for _, kind := range suite.invalidKinds {
			needsUpdate := suite.collection.NeedsUpdate(kind)
			assert.False(t, needsUpdate, "Invalid kind %s should return false for NeedsUpdate", kind)
		}
	})

	t.Run("Reset_ValidKinds", func(t *testing.T) {
		for _, kind := range suite.supportedKinds {
			assert.NotPanics(t, func() {
				suite.collection.Reset(kind)
			}, "Reset should not panic for valid kind %s", kind)
		}
	})

	t.Run("Reset_InvalidKinds", func(t *testing.T) {
		for _, kind := range suite.invalidKinds {
			assert.NotPanics(t, func() {
				suite.collection.Reset(kind)
			}, "Reset should not panic for invalid kind %s", kind)
		}
	})

	t.Run("LoadData_ValidKinds", func(t *testing.T) {
		for _, kind := range suite.supportedKinds {
			assert.NotPanics(t, func() {
				suite.collection.LoadData(kind)
			}, "LoadData should not panic for valid kind %s", kind)
		}
	})

	t.Run("LoadData_InvalidKinds", func(t *testing.T) {
		for _, kind := range suite.invalidKinds {
			assert.NotPanics(t, func() {
				suite.collection.LoadData(kind)
			}, "LoadData should not panic for invalid kind %s", kind)
		}
	})
}

// TestCollectionPageInterface tests Page interface compliance
func (suite *CollectionTestSuite) TestCollectionPageInterface(t *testing.T) {
	t.Run("GetPage_ValidKinds", func(t *testing.T) {
		for _, kind := range suite.supportedKinds {
			page, err := suite.collection.GetPage(kind, 0, 10, sdk.SortSpec{}, "")

			if err == nil && page != nil {
				// Test Page interface compliance
				assert.Equal(t, kind, page.GetKind(), "Page kind should match requested kind")
				assert.GreaterOrEqual(t, page.GetTotalItems(), 0, "Total items should be non-negative")
				assert.GreaterOrEqual(t, page.GetExpectedTotal(), 0, "Expected total should be non-negative")

				// State should be one of the valid states
				state := page.GetState()
				validStates := []types.LoadState{
					types.StateStale, types.StateFetching, types.StatePartial,
					types.StateLoaded, types.StatePending, types.StateError,
				}
				assert.Contains(t, validStates, state, "Page state should be valid")
			}
		}
	})

	t.Run("GetPage_InvalidKinds", func(t *testing.T) {
		for _, kind := range suite.invalidKinds {
			_, err := suite.collection.GetPage(kind, 0, 10, sdk.SortSpec{}, "")
			assert.Error(t, err, "GetPage should return error for invalid kind %s", kind)
		}
	})

	t.Run("GetPage_ZeroPageSize", func(t *testing.T) {
		for _, kind := range suite.supportedKinds {
			page, err := suite.collection.GetPage(kind, 0, 0, sdk.SortSpec{}, "")
			if err == nil && page != nil {
				assert.Equal(t, kind, page.GetKind(), "Page kind should match even with zero page size")
			}
		}
	})

	t.Run("GetPage_WithFilter", func(t *testing.T) {
		for _, kind := range suite.supportedKinds {
			page, err := suite.collection.GetPage(kind, 0, 5, sdk.SortSpec{}, "test")
			if err == nil && page != nil {
				assert.Equal(t, kind, page.GetKind(), "Page kind should match with filter")
				assert.GreaterOrEqual(t, page.GetTotalItems(), 0, "Filtered results should be non-negative")
			}
		}
	})
}

// TestCollectionBoundaryConditions tests edge cases and boundary conditions
func (suite *CollectionTestSuite) TestCollectionBoundaryConditions(t *testing.T) {
	t.Run("ExtremePageParameters", func(t *testing.T) {
		extremeParams := []struct {
			first    int
			pageSize int
			desc     string
		}{
			{-1, 10, "negative_first"},
			{0, -1, "negative_pageSize"},
			{1000000, 1, "very_large_first"},
			{0, 1000000, "very_large_pageSize"},
			{-10, -5, "both_negative"},
		}

		for _, kind := range suite.supportedKinds {
			for _, param := range extremeParams {
				t.Run(param.desc+"_"+string(kind), func(t *testing.T) {
					assert.NotPanics(t, func() {
						_, _ = suite.collection.GetPage(kind, param.first, param.pageSize, sdk.SortSpec{}, "")
					}, "Extreme parameters should not cause panic")
				})
			}
		}
	})

	t.Run("CrossKindIsolation", func(t *testing.T) {
		// Test that operations on one kind don't affect others
		if len(suite.supportedKinds) >= 2 {
			kind1 := suite.supportedKinds[0]
			kind2 := suite.supportedKinds[1]

			// Reset one kind
			suite.collection.Reset(kind1)

			// Operations on other kinds should still work
			assert.NotPanics(t, func() {
				suite.collection.NeedsUpdate(kind2)
			}, "Operations on other kinds should work after reset")

			assert.NotPanics(t, func() {
				_, _ = suite.collection.GetPage(kind2, 0, 5, sdk.SortSpec{}, "")
			}, "GetPage on other kinds should work after reset")
		}
	})
}

// TestCollectionCrudInterface tests CRUD method signature compliance
func (suite *CollectionTestSuite) TestCollectionCrudInterface(t *testing.T) {
	t.Run("Crud_MethodExists", func(t *testing.T) {
		// Test that Crud method exists and can be called without panicking
		// We don't test actual CRUD functionality here as it's domain-specific
		for _, kind := range suite.supportedKinds {
			assert.NotPanics(t, func() {
				// Using nil item - actual CRUD logic is tested in domain-specific tests
				_ = suite.collection.Crud(kind, "create", nil)
			}, "Crud method should exist and not panic for valid kind %s", kind)
		}
	})
}

// Generic test functions that create the test suites for each collection type

func TestNamesCollectionGeneric(t *testing.T) {
	collection := names.NewNamesCollection()
	supportedKinds := []types.ListKind{
		names.NamesAll, names.NamesCustom, names.NamesPrefund,
		names.NamesRegular, names.NamesBaddress,
	}
	invalidKinds := []types.ListKind{"", "invalid", "NotAListKind", "monitors"}

	suite := CreateCollectionTestSuite(collection, supportedKinds, invalidKinds, "names", "names")

	t.Run("Interface", suite.TestCollectionInterface)
	t.Run("StateManagement", suite.TestCollectionStateManagement)
	t.Run("PageInterface", suite.TestCollectionPageInterface)
	t.Run("BoundaryConditions", suite.TestCollectionBoundaryConditions)
	t.Run("CrudInterface", suite.TestCollectionCrudInterface)
}

func TestAbisCollectionGeneric(t *testing.T) {
	collection := abis.NewAbisCollection()
	supportedKinds := []types.ListKind{
		abis.AbisDownloaded, abis.AbisKnown, abis.AbisFunctions, abis.AbisEvents,
	}
	invalidKinds := []types.ListKind{"", "invalid", "NotAListKind", "names"}

	suite := CreateCollectionTestSuite(collection, supportedKinds, invalidKinds, "abis", "abis")

	t.Run("Interface", suite.TestCollectionInterface)
	t.Run("StateManagement", suite.TestCollectionStateManagement)
	t.Run("PageInterface", suite.TestCollectionPageInterface)
	t.Run("BoundaryConditions", suite.TestCollectionBoundaryConditions)
	t.Run("CrudInterface", suite.TestCollectionCrudInterface)
}

func TestMonitorsCollectionGeneric(t *testing.T) {
	collection := monitors.NewMonitorsCollection()
	supportedKinds := []types.ListKind{monitors.MonitorsList}
	invalidKinds := []types.ListKind{"", "invalid", "NotAListKind", "names"}

	suite := CreateCollectionTestSuite(collection, supportedKinds, invalidKinds, "monitors", "monitors")

	t.Run("Interface", suite.TestCollectionInterface)
	t.Run("StateManagement", suite.TestCollectionStateManagement)
	t.Run("PageInterface", suite.TestCollectionPageInterface)
	t.Run("BoundaryConditions", suite.TestCollectionBoundaryConditions)
	t.Run("CrudInterface", suite.TestCollectionCrudInterface)
}

// NewCollectionTests tests creating new collection instances
func TestNewCollectionTests(t *testing.T) {
	t.Run("NewNamesCollection", func(t *testing.T) {
		assert.NotPanics(t, func() {
			collection := names.NewNamesCollection()
			assert.NotNil(t, collection, "NewNamesCollection should return non-nil")
		})
	})

	t.Run("NewAbisCollection", func(t *testing.T) {
		assert.NotPanics(t, func() {
			collection := abis.NewAbisCollection()
			assert.NotNil(t, collection, "NewAbisCollection should return non-nil")
		})
	})

	t.Run("NewMonitorsCollection", func(t *testing.T) {
		assert.NotPanics(t, func() {
			collection := monitors.NewMonitorsCollection()
			assert.NotNil(t, collection, "NewMonitorsCollection should return non-nil")
		})
	})
}
