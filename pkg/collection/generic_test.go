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
	supportedFacets   []types.DataFacet
	invalidFacets     []types.DataFacet
	collectionName    string
	expectedStoreName string
}

// CreateCollectionTestSuite creates a test suite for any Collection implementation
func CreateCollectionTestSuite(
	collection types.Collection,
	supportedFacets []types.DataFacet,
	invalidFacets []types.DataFacet,
	collectionName string,
	expectedStoreName string,
) *CollectionTestSuite {
	return &CollectionTestSuite{
		collection:        collection,
		supportedFacets:   supportedFacets,
		invalidFacets:     invalidFacets,
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

	t.Run("GetSupportedFacets", func(t *testing.T) {
		facets := suite.collection.GetSupportedFacets()
		assert.NotEmpty(t, facets, "Supported facets should not be empty")
		assert.ElementsMatch(t, suite.supportedFacets, facets, "Supported facets should match expected")
	})

	t.Run("GetStoreForFacet_ValidFacets", func(t *testing.T) {
		for _, facet := range suite.supportedFacets {
			storeName := suite.collection.GetStoreForFacet(facet)
			assert.NotEmpty(t, storeName, "Store name should not be empty for valid facet %s", facet)
		}
	})

	t.Run("GetStoreForFacet_invalidFacets", func(t *testing.T) {
		for _, facet := range suite.invalidFacets {
			storeName := suite.collection.GetStoreForFacet(facet)
			assert.Empty(t, storeName, "Store name should be empty for invalid facet %s", facet)
		}
	})
}

// TestCollectionStateManagement tests NeedsUpdate/Reset/LoadData behavior
func (suite *CollectionTestSuite) TestCollectionStateManagement(t *testing.T) {
	t.Run("NeedsUpdate_ValidFacets", func(t *testing.T) {
		for _, facet := range suite.supportedFacets {
			needsUpdate := suite.collection.NeedsUpdate(facet)
			// We can't assert true/false here as it depends on actual state
			// But the call should not panic
			_ = needsUpdate
		}
	})

	t.Run("NeedsUpdate_InvalidFacets", func(t *testing.T) {
		for _, facet := range suite.invalidFacets {
			needsUpdate := suite.collection.NeedsUpdate(facet)
			assert.False(t, needsUpdate, "Invalid facet %s should return false for NeedsUpdate", facet)
		}
	})

	t.Run("Reset_ValidFacets", func(t *testing.T) {
		for _, facet := range suite.supportedFacets {
			assert.NotPanics(t, func() {
				suite.collection.Reset(facet)
			}, "Reset should not panic for valid facet %s", facet)
		}
	})

	t.Run("Reset_InvalidFacets", func(t *testing.T) {
		for _, facet := range suite.invalidFacets {
			assert.NotPanics(t, func() {
				suite.collection.Reset(facet)
			}, "Reset should not panic for invalid facet %s", facet)
		}
	})

	t.Run("LoadData_ValidFacets", func(t *testing.T) {
		for _, facet := range suite.supportedFacets {
			assert.NotPanics(t, func() {
				suite.collection.LoadData(facet)
			}, "LoadData should not panic for valid facet %s", facet)
		}
	})

	t.Run("LoadData_InvalidFacets", func(t *testing.T) {
		for _, facet := range suite.invalidFacets {
			assert.NotPanics(t, func() {
				suite.collection.LoadData(facet)
			}, "LoadData should not panic for invalid facet %s", facet)
		}
	})
}

// TestCollectionPageInterface tests Page interface compliance
func (suite *CollectionTestSuite) TestCollectionPageInterface(t *testing.T) {
	t.Run("GetPage_ValidFacets", func(t *testing.T) {
		for _, facet := range suite.supportedFacets {
			page, err := suite.collection.GetPage(facet, 0, 10, sdk.SortSpec{}, "")

			if err == nil && page != nil {
				// Test Page interface compliance
				assert.Equal(t, facet, page.GetFacet(), "Page facet should match requested facet")
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

	t.Run("GetPage_InvalidFacets", func(t *testing.T) {
		for _, facet := range suite.invalidFacets {
			_, err := suite.collection.GetPage(facet, 0, 10, sdk.SortSpec{}, "")
			assert.Error(t, err, "GetPage should return error for invalid facet %s", facet)
		}
	})

	t.Run("GetPage_ZeroPageSize", func(t *testing.T) {
		for _, facet := range suite.supportedFacets {
			page, err := suite.collection.GetPage(facet, 0, 0, sdk.SortSpec{}, "")
			if err == nil && page != nil {
				assert.Equal(t, facet, page.GetFacet(), "Page facet should match even with zero page size")
			}
		}
	})

	t.Run("GetPage_WithFilter", func(t *testing.T) {
		for _, facet := range suite.supportedFacets {
			page, err := suite.collection.GetPage(facet, 0, 5, sdk.SortSpec{}, "test")
			if err == nil && page != nil {
				assert.Equal(t, facet, page.GetFacet(), "Page facet should match with filter")
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

		for _, facet := range suite.supportedFacets {
			for _, param := range extremeParams {
				t.Run(param.desc+"_"+string(facet), func(t *testing.T) {
					assert.NotPanics(t, func() {
						_, _ = suite.collection.GetPage(facet, param.first, param.pageSize, sdk.SortSpec{}, "")
					}, "Extreme parameters should not cause panic")
				})
			}
		}
	})

	t.Run("CrossFacetIsolation", func(t *testing.T) {
		// Test that operations on one facet don't affect others
		if len(suite.supportedFacets) >= 2 {
			facet1 := suite.supportedFacets[0]
			facet2 := suite.supportedFacets[1]

			// Reset one facet
			suite.collection.Reset(facet1)

			// Operations on other facets should still work
			assert.NotPanics(t, func() {
				suite.collection.NeedsUpdate(facet2)
			}, "Operations on other facets should work after reset")

			assert.NotPanics(t, func() {
				_, _ = suite.collection.GetPage(facet2, 0, 5, sdk.SortSpec{}, "")
			}, "GetPage on other facets should work after reset")
		}
	})
}

// TestCollectionCrudInterface tests CRUD method signature compliance
func (suite *CollectionTestSuite) TestCollectionCrudInterface(t *testing.T) {
	t.Run("Crud_MethodExists", func(t *testing.T) {
		// Test that Crud method exists and can be called without panicking
		// We don't test actual CRUD functionality here as it's domain-specific
		for _, facet := range suite.supportedFacets {
			assert.NotPanics(t, func() {
				// Using nil item - actual CRUD logic is tested in domain-specific tests
				_ = suite.collection.Crud(facet, "create", nil)
			}, "Crud method should exist and not panic for valid facet %s", facet)
		}
	})
}

// Generic test functions that create the test suites for each collection type

func TestNamesCollectionGeneric(t *testing.T) {
	collection := names.NewNamesCollection()
	supportedFacets := []types.DataFacet{
		names.NamesAll, names.NamesCustom, names.NamesPrefund,
		names.NamesRegular, names.NamesBaddress,
	}
	invalidFacets := []types.DataFacet{"", "invalid", "NotADataFacet", "monitors"}

	suite := CreateCollectionTestSuite(collection, supportedFacets, invalidFacets, "names", "names")

	t.Run("Interface", suite.TestCollectionInterface)
	t.Run("StateManagement", suite.TestCollectionStateManagement)
	t.Run("PageInterface", suite.TestCollectionPageInterface)
	t.Run("BoundaryConditions", suite.TestCollectionBoundaryConditions)
	t.Run("CrudInterface", suite.TestCollectionCrudInterface)
}

func TestAbisCollectionGeneric(t *testing.T) {
	collection := abis.NewAbisCollection()
	supportedFacets := []types.DataFacet{
		abis.AbisDownloaded, abis.AbisKnown, abis.AbisFunctions, abis.AbisEvents,
	}
	invalidFacets := []types.DataFacet{"", "invalid", "NotADataFacet", "names"}

	suite := CreateCollectionTestSuite(collection, supportedFacets, invalidFacets, "abis", "abis")

	t.Run("Interface", suite.TestCollectionInterface)
	t.Run("StateManagement", suite.TestCollectionStateManagement)
	t.Run("PageInterface", suite.TestCollectionPageInterface)
	t.Run("BoundaryConditions", suite.TestCollectionBoundaryConditions)
	t.Run("CrudInterface", suite.TestCollectionCrudInterface)
}

func TestMonitorsCollectionGeneric(t *testing.T) {
	collection := monitors.NewMonitorsCollection()
	supportedFacets := []types.DataFacet{monitors.MonitorsList}
	invalidFacets := []types.DataFacet{"", "invalid", "NotADataFacet", "names"}

	suite := CreateCollectionTestSuite(collection, supportedFacets, invalidFacets, "monitors", "monitors")

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
