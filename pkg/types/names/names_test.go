// NAMES_ROUTE
package names

// func setupTest(t *testing.T) *NamesCollection {
// 	t.Helper()

// 	th := msgs.NewTestHelpers()
// 	defer th.Cleanup()

// 	collection := NewNamesCollection()
// 	assert.NotNil(t, collection, "NewNamesCollection should return a non-nil collection")

// 	loadedCh := make(chan bool)
// 	msgs.On(msgs.EventDataLoaded, func(optionalData ...interface{}) {
// 		var payload types.DataLoadedPayload
// 		isRelevantEvent := false
// 		if len(optionalData) >= 1 {
// 			if p, ok := optionalData[len(optionalData)-1].(types.DataLoadedPayload); ok {
// 				payload = p
// 				if payload.ListKind == NamesAll {
// 					isRelevantEvent = true
// 				}
// 			}
// 		}

// 		if isRelevantEvent {
// 			if collection.allFacet.GetState() == facets.StateLoaded && payload.CurrentCount == payload.ExpectedTotal && payload.ExpectedTotal > 0 {
// 				select {
// 				case <-loadedCh:
// 					// Already closed
// 				default:
// 					close(loadedCh)
// 				}
// 			}
// 		}
// 	})

// 	collection.LoadData(NamesAll)

// 	select {
// 	case <-loadedCh:
// 		assert.Equal(t, facets.StateLoaded, collection.allFacet.GetState(), "allFacet should be in StateLoaded after event")
// 		assert.Greater(t, collection.allFacet.Count(), 0, "allFacet should have items after load")
// 		assert.True(t, collection.allFacet.IsLoaded(), "allFacet.IsLoaded() should be true")
// 	case <-time.After(15 * time.Second): // Increased timeout for CI or slower environments
// 		t.Fatalf("Timeout waiting for NamesCollection (%s) to load. Current state of allFacet: %s. Facet count: %d. Expected total in facet (from progress): %d. IsLoaded: %v",
// 			NamesAll, collection.allFacet.GetState(), collection.allFacet.Count(), collection.allFacet.ExpectedCount(), collection.allFacet.IsLoaded())
// 	}

// 	return collection
// }

// func TestNamesCollection_SimpleAccessors(t *testing.T) {
// 	collection := setupTest(t)

// 	t.Run("TestAllFacet_InitialLoad", func(t *testing.T) {
// 		// Access the NamesAll facet directly from the collection
// 		allFacet := collection.allFacet // Assuming 'allFacet' is the correct field name for the 'All' names facet
// 		assert.NotNil(t, allFacet, "allFacet should not be nil")

// 		// Check if the facet is loaded (assuming setupTest ensures this)
// 		assert.True(t, allFacet.IsLoaded(), "allFacet should be loaded")

// 		// Check that the count of items is non-negative
// 		// In a real scenario, after initial load, this might be expected to be > 0
// 		// if there's default data.
// 		count := allFacet.Count()
// 		assert.True(t, count >= 0, "allFacet count should be non-negative")

// 		// To get actual items, you would use GetPage.
// 		// For this initial test, checking IsLoaded and Count should suffice to confirm initialization.
// 		// Example of getting items (though not strictly needed for this specific subtest's goal):
// 		// pageResult, err := allFacet.GetPage(0, 10, nil, sdk.SortSpec{}, nil) // Get first 10 items
// 		// assert.NoError(t, err)
// 		// assert.NotNil(t, pageResult)
// 		// assert.NotNil(t, pageResult.Items)
// 	})
// }

// Test cases 1:
// - Ensure that the names collection is properly initialized and can handle CRUD operations
// - Ensure that the names collection can handle pagination and sorting correctly
// - Ensure that the names collection can handle custom names, prefund names, and regular names correctly
// - Ensure that the names collection can handle baddress names correctly
// - Ensure that the names collection can handle all names correctly
// - Ensure that the names collection can handle names with different properties (e.g., custom, prefund, regular, baddress)
// - Ensure that the names collection can handle names with different tags and sources
// - Ensure that the names collection can handle names with different decimals
// - Ensure that the names collection can handle names with different pre-fund amounts
// - Ensure that the names collection can handle names with different parts
// - Ensure that the names collection can handle names with different isCustom flags
// - Ensure that the names collection can handle names with different deleted flags
// - Ensure that the names collection can handle names with different isContract flags
// - Ensure that the names collection can handle names with different isErc20 flags
// - Ensure that the names collection can handle names with different isErc721 flags
// - Ensure that the names collection can handle names with different isPrefund flags
// - Ensure that the names collection can handle names with different addresses
// - Ensure that the names collection can handle names with different name strings
// - Ensure that the names collection can handle names with different tags
// - Ensure that the names collection can handle names with different sources

// Test cases 2:
// - Add a new name
// - Autoname an existing name
// - Edit an existing name
// - Delete an existing name
// - Undelete a deleted name
// - Remove a name (should fail if not deleted)
// - Delete a name and then remove it
// - Remove a name that has been deleted
// - Ensure removed names are not found in any search method
// - Ensure names can be added, edited, deleted, and removed correctly
// - Ensure that the CRUD operations work as expected with the new facet-based architecture
