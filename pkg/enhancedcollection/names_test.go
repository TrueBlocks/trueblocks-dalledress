package enhancedcollection_test

// func TestNamesCollection(t *testing.T) {
// 	// Enable test mode for messaging and ensure cleanup
// 	helper := msgs.NewTestHelpers()
// 	defer helper.Cleanup()

// 	// Create a new collection
// 	// names := enhancedcollection.NewNamesCollection()

// 	// Test initial state
// 	if state := names.GetState(); state != enhancedfacet.StateStale {
// 		t.Errorf("Expected initial state to be Stale, got %v", state)
// 	}

// 	if count := len(names.GetItems()); count != 0 {
// 		t.Errorf("Expected 0 items initially, got %d", count)
// 	}

// 	// Create a channel to wait for loading completion
// 	loadedCh := msgs.WaitForEvent(msgs.EventDataLoaded)

// 	// Test loading
// 	err := names.Load()
// 	if err != nil {
// 		t.Fatalf("Error loading: %v", err)
// 	}

// 	// Load test data
// 	names.LoadTestData()

// 	// Wait for the loaded event instead of arbitrary sleep
// 	select {
// 	case <-loadedCh:
// 		// Event received, continue with test
// 	case <-time.After(time.Second):
// 		t.Fatal("Timed out waiting for loaded event")
// 	}

// 	// Test state after loading
// 	if state := names.GetState(); state != enhancedfacet.StateLoaded {
// 		t.Errorf("Expected state after loading to be Loaded, got %v", state)
// 	}

// 	// Print debug info before checking items
// 	names.DebugInfo()

// 	// Test item count
// 	items := names.GetItems()
// 	if count := len(items); count != 5 {
// 		t.Errorf("Expected 5 items after loading, got %d", count)
// 	}

// 	// Create a channel to wait for cancellation
// 	canceledCh := msgs.WaitForEvent(msgs.EventStatus)

// 	// Test cancellation
// 	names.SimulateCancellation()

// 	// Wait for cancellation notification instead of sleep
// 	select {
// 	case <-canceledCh:
// 		// Event received, continue with test
// 	case <-time.After(time.Second):
// 		t.Fatal("Timed out waiting for cancellation event")
// 	}

// 	// After cancellation, state should be Stale
// 	if state := names.GetState(); state != enhancedfacet.StateStale {
// 		t.Errorf("Expected state after cancellation to be Stale, got %v", state)
// 	}

// 	// Create another channel to wait for reloading
// 	reloadedCh := msgs.WaitForEvent(msgs.EventDataLoaded)

// 	// Test marking as stale
// 	names.LoadTestData() // reload data

// 	// Wait for reload completion
// 	select {
// 	case <-reloadedCh:
// 		// Event received, continue with test
// 	case <-time.After(time.Second):
// 		t.Fatal("Timed out waiting for reload event")
// 	}

// 	names.MarkStale()

// 	if state := names.GetState(); state != enhancedfacet.StateStale {
// 		t.Errorf("Expected state after marking stale to be Stale, got %v", state)
// 	}
// }

// func TestFilterBehavior(t *testing.T) {
// 	// Enable test mode for messaging and ensure cleanup
// 	helper := msgs.NewTestHelpers()
// 	defer helper.Cleanup()

// 	// Create a new collection
// 	// names := enhancedcollection.NewNamesCollection()

// 	// Create a channel to wait for loading completion
// 	loadedCh := msgs.WaitForEvent(msgs.EventDataLoaded)

// 	// Load normal data
// 	names.LoadTestData()

// 	// Wait for loaded event
// 	select {
// 	case <-loadedCh:
// 		// Event received, continue with test
// 	case <-time.After(time.Second):
// 		t.Fatal("Timed out waiting for loaded event")
// 	}

// 	// Count before adding filtered items
// 	countBefore := len(names.GetItems())

// 	// Create a channel to wait for filtered items
// 	filteredCh := msgs.WaitForEvent(msgs.EventDataLoaded)

// 	// Add items that should be filtered
// 	names.AddFilteredItems()

// 	// Wait for filtered items
// 	select {
// 	case <-filteredCh:
// 		// Event received, continue with test
// 	case <-time.After(time.Second):
// 		t.Fatal("Timed out waiting for filtered items event")
// 	}

// 	// Count after - should be the same if filtering works
// 	countAfter := len(names.GetItems())

// 	if countBefore != countAfter {
// 		t.Errorf("Filtering failed: count before %d, count after %d",
// 			countBefore, countAfter)
// 	}
// }
