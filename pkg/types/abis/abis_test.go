package abis

import (
	"testing"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/facets"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
)

func TestEnhancedAbisCollection(t *testing.T) {
	// Enable test mode for messaging
	msgs.SetTestMode(true)
	defer msgs.SetTestMode(false)

	// Create a new enhanced collection
	collection := NewAbisCollection()

	// Test initial state
	if state := collection.downloadedFacet.GetState(); state != facets.StateStale1 {
		t.Errorf("Expected initial state to be Stale, got %v", state)
	}

	// In a real test, we'd load data and verify it works properly
	// However, since that would require hitting the SDK, we'll just verify
	// that the collection was created successfully with all facets

	// Test that all facets were initialized
	if collection.downloadedFacet.GetState() != facets.StateStale1 {
		t.Error("downloadedFacet should be in StateStale")
	}
	if collection.knownFacet.GetState() != facets.StateStale1 {
		t.Error("knownFacet should be in StateStale")
	}
	if collection.functionsFacet.GetState() != facets.StateStale1 {
		t.Error("functionsFacet should be in StateStale")
	}
	if collection.eventsFacet.GetState() != facets.StateStale1 {
		t.Error("eventsFacet should be in StateStale")
	}

	// Test that NeedsUpdate works properly
	if !collection.NeedsUpdate(AbisDownloaded) {
		t.Error("downloadedFacet should need update")
	}
	if !collection.NeedsUpdate(AbisKnown) {
		t.Error("knownFacet should need update")
	}
	if !collection.NeedsUpdate(AbisFunctions) {
		t.Error("functionsFacet should need update")
	}
	if !collection.NeedsUpdate(AbisEvents) {
		t.Error("eventsFacet should need update")
	}

	// Test Reset function
	collection.Reset(AbisDownloaded)
	if !collection.NeedsUpdate(AbisDownloaded) {
		t.Error("downloadedFacet should need update after reset")
	}
}
