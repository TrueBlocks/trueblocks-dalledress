package enhancedcollection_test

// func TestABIsCollection(t *testing.T) {
// 	// Enable test mode for messaging and ensure cleanup
// 	helper := msgs.NewTestHelpers()
// 	defer helper.Cleanup()

// 	// Create a new collection
// 	abis := enhancedcollection.NewABIsCollection()

// 	// Test initial state
// 	if state := abis.GetState(); state != enhancedfacet.StateStale {
// 		t.Errorf("Expected initial state to be Stale, got %v", state)
// 	}

// 	if count := len(abis.GetItems()); count != 0 {
// 		t.Errorf("Expected 0 items initially, got %d", count)
// 	}

// 	// Create a channel to wait for loading completion
// 	loadedCh := msgs.WaitForEvent(msgs.EventDataLoaded)

// 	// Test loading
// 	err := abis.Load()
// 	if err != nil {
// 		t.Fatalf("Error loading: %v", err)
// 	}

// 	// Load test data
// 	abis.LoadTestData()

// 	// Wait for the loaded event instead of arbitrary sleep
// 	select {
// 	case <-loadedCh:
// 		// Event received, continue with test
// 	case <-time.After(time.Second):
// 		t.Fatal("Timed out waiting for loaded event")
// 	}

// 	// Test state after loading
// 	if state := abis.GetState(); state != enhancedfacet.StateLoaded {
// 		t.Errorf("Expected state after loading to be Loaded, got %v", state)
// 	}

// 	// Test item count
// 	items := abis.GetItems()
// 	if count := len(items); count != 3 {
// 		t.Errorf("Expected 3 items after loading, got %d", count)
// 	}
// }
