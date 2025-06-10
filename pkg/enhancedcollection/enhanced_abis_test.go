package enhancedcollection_test

import (
	"testing"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/enhancedcollection"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/enhancedfacet"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
)

func TestEnhancedAbisCollection(t *testing.T) {
	// Enable test mode for messaging
	msgs.SetTestMode(true)
	defer msgs.SetTestMode(false)

	// Create a new enhanced collection
	collection := enhancedcollection.NewEnhancedAbisCollection()

	// Test initial state
	if state := collection.GetDownloadedState(); state != enhancedfacet.StateStale {
		t.Errorf("Expected initial state to be Stale, got %v", state)
	}

	// In a real test, we'd load data and verify it works properly
	// However, since that would require hitting the SDK, we'll just verify
	// that the collection was created successfully with all facets

	// Test that all facets were initialized
	if collection.GetDownloadedState() != enhancedfacet.StateStale {
		t.Error("downloadedFacet should be in StateStale")
	}
	if collection.GetKnownState() != enhancedfacet.StateStale {
		t.Error("knownFacet should be in StateStale")
	}
	if collection.GetFunctionsState() != enhancedfacet.StateStale {
		t.Error("functionsFacet should be in StateStale")
	}
	if collection.GetEventsState() != enhancedfacet.StateStale {
		t.Error("eventsFacet should be in StateStale")
	}

	// Test that NeedsUpdate works properly
	if !collection.NeedsUpdate(enhancedcollection.AbisDownloaded) {
		t.Error("downloadedFacet should need update")
	}
	if !collection.NeedsUpdate(enhancedcollection.AbisKnown) {
		t.Error("knownFacet should need update")
	}
	if !collection.NeedsUpdate(enhancedcollection.AbisFunctions) {
		t.Error("functionsFacet should need update")
	}
	if !collection.NeedsUpdate(enhancedcollection.AbisEvents) {
		t.Error("eventsFacet should need update")
	}

	// Test Reset function
	collection.Reset(enhancedcollection.AbisDownloaded)
	if !collection.NeedsUpdate(enhancedcollection.AbisDownloaded) {
		t.Error("downloadedFacet should need update after reset")
	}
}
