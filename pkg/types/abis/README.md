# ABIs Package (`pkg/types/abis`)

The `abis` package is responsible for managing Application Binary Interface (ABI) data within the TrueBlocks ecosystem. It provides a structured way to handle different categories of ABIs (downloaded by the user, known by TrueBlocks) and their constituent parts (functions and events). This package makes heavy use of the `facets` package for efficient data loading, caching, filtering, sorting, and pagination.

## Core Concepts

### `AbisCollection`

This is the central struct that orchestrates various ABI-related data views (facets). It holds separate facets for:
- **Downloaded ABIs**: ABIs that the user has explicitly fetched or added.
- **Known ABIs**: A curated list of ABIs known to the TrueBlocks system.
- **Functions**: A collection of all unique function definitions extracted from all available ABIs.
- **Events**: A collection of all unique event definitions extracted from all available ABIs.

```go
type AbisCollection struct {
    downloadedFacet facets.Facet[coreTypes.Abi]
    knownFacet      facets.Facet[coreTypes.Abi]
    functionsFacet  facets.Facet[coreTypes.Function]
    eventsFacet     facets.Facet[coreTypes.Function]
}
```

### Shared Data Sources

A key design principle in this package is the use of shared `source.Source` implementations. This minimizes redundant calls to the TrueBlocks SDK:
- `sharedAbisListSource`: Fetches a list of all ABIs (both downloaded and known) once. The `downloadedFacet` and `knownFacet` then filter this shared data in memory based on the `IsKnown` property of `coreTypes.Abi`.
- `sharedAbisDetailsSource`: Fetches detailed information (including functions and events) for all ABIs once. The `functionsFacet` and `eventsFacet` then filter this shared data in memory based on the `FunctionType` property of `coreTypes.Function`.

This approach significantly improves performance by reducing the number of potentially expensive SDK queries from four to two for a full data refresh.

### `ListKind` Constants

The package defines specific `types.ListKind` constants to identify the different ABI data categories:

```go
const (
    AbisDownloaded types.ListKind = "Downloaded"
    AbisKnown      types.ListKind = "Known"
    AbisFunctions  types.ListKind = "Functions"
    AbisEvents     types.ListKind = "Events"
)
```
These are registered with the `types` package and used throughout the application to refer to these specific data sets.

## Public API

### `NewAbisCollection() AbisCollection`

Constructs and returns a new `AbisCollection`. It initializes the four facets, configuring them with their respective `ListKind`, filter functions (to differentiate between downloaded/known ABIs, and functions/events), and the shared data sources.

For `functionsFacet` and `eventsFacet`, it uses `IsDupFuncByEncoding()` to ensure that functions/events with the same signature encoding are treated as duplicates and only one instance is stored.

### `(ac *AbisCollection) LoadData(listKind types.ListKind)`

Initiates data loading for the specified `listKind`. It checks if an update is needed (via `NeedsUpdate`) before starting the load. Loading is performed asynchronously in a goroutine. Once loaded, it emits a `msgs.LoadedMsg` to notify other parts of the application.

### `(ac *AbisCollection) NeedsUpdate(listKind types.ListKind) bool`

Checks if the facet corresponding to the given `listKind` needs to reload its data (i.e., if its state is `StateStale`).

### `(ac *AbisCollection) Reset(listKind types.ListKind)`

Resets the state of the facet for the given `listKind` to `StateStale`. This will typically force a data reload on the next access or `LoadData` call.

### `(ac *AbisCollection) GetPage(...) (AbisPage, error)`

```go
func (ac *AbisCollection) GetPage(
    listKind types.ListKind,
    first, pageSize int,
    sortSpec sdk.SortSpec,
    filter string,
) (AbisPage, error)
```

Retrieves a paginated, filtered, and sorted view of ABI data. 
- It first ensures data is loaded by calling `ac.LoadData(listKind)`.
- It then constructs an `AbisPage` struct containing the items, total counts, and loading state.
- Filtering is done based on the `filter` string (case-insensitive search across relevant fields like ABI name, address, function name, signature, encoding).
- Sorting is performed using `sdk.SortAbis` for ABI lists and `sdk.SortFunctions` for function/event lists, based on the provided `sdk.SortSpec`.

The `AbisPage` struct is defined as:
```go
type AbisPage struct {
    Kind          types.ListKind       `json:"kind"`
    Abis          []coreTypes.Abi      `json:"abis,omitempty"`
    Functions     []coreTypes.Function `json:"functions,omitempty"`
    TotalItems    int                  `json:"totalItems"`
    ExpectedTotal int                  `json:"expectedTotal"`
    IsFetching    bool                 `json:"isFetching"`
    State         facets.LoadState     `json:"state"`
}
```

### `(ac *AbisCollection) AbisCrud(listKind types.ListKind, op crud.Operation, abi *coreTypes.Abi) error`

Performs Create, Read, Update, Delete (CRUD) operations. Currently, only `crud.Remove` is implemented for the `AbisDownloaded` list kind.
- For `crud.Remove`, it uses `sdk.AbisOptions{Decache: true}` to remove the ABI from the TrueBlocks cache via an SDK call.
- It then updates the internal `downloadedFacet` data by removing the matching ABI.
- Emits status messages about the operation's success or failure.

### `IsDupFuncByEncoding() func(existing []coreTypes.Function, newItem *coreTypes.Function) bool`

Returns a deduplication function specifically for `coreTypes.Function` items (which are used for both functions and events). This function uses a map (`seen`) to keep track of function/event encodings that have already been added to the list. This provides an efficient O(1) lookup for duplicates, which is crucial when processing potentially large numbers of functions and events from many ABIs.
The map is reset if the `existing` slice becomes empty, indicating a full data reload.

### Shared Source Getters

- **`GetSharedAbisListSource() source.Source[coreTypes.Abi]`**: Returns the singleton instance of the shared data source for ABI lists. If not already initialized, it creates a new `source.SDKSource` configured to call `sdk.AbisOptions{}.AbisList()`.
- **`GetSharedAbisDetailsSource() source.Source[coreTypes.Function]`**: Returns the singleton instance of the shared data source for ABI details (functions/events). If not already initialized, it creates a new `source.SDKSource` configured to call `sdk.AbisOptions{}.AbisDetails()`.

Both source getters use `preferences.GetPreferredChainName()` to ensure data is fetched for the currently selected blockchain.

## Example Usage (Conceptual)

Below is a conceptual example. For a runnable `main.go` file, see the section [Runnable Example](#runnable-example) further down.

```go
// Initialize the collection
abiCollection := abis.NewAbisCollection()

// Load downloaded ABIs (triggers async load if stale)
abiCollection.LoadData(abis.AbisDownloaded)

// Later, retrieve a page of downloaded ABIs
// Sort by Name, ascending, no filter, first 10 items
page, err := abiCollection.GetPage(
    abis.AbisDownloaded,
    0, 10,
    sdk.SortSpec{Field: "Name", Direction: sdk.SortAsc},
    "", // no filter
)
if err != nil {
    // handle error
}

fmt.Printf("Displaying %d of %d downloaded ABIs (State: %s):
", len(page.Abis), page.TotalItems, page.State)
for _, abi := range page.Abis {
    fmt.Printf("- %s (%s)
", abi.Name, abi.Address.Hex())
}

// To remove an ABI:
addressToRemove := "0xSomeAddress..."
var abiToRemove coreTypes.Abi
// (Assume abiToRemove is populated, e.g., from the page.Abis)
for _, a := range page.Abis {
    if a.Address.Hex() == addressToRemove {
        abiToRemove = a
        break
    }
}

if abiToRemove.Address.Hex() != "" { // check if found
    err = abiCollection.AbisCrud(abis.AbisDownloaded, crud.Remove, &abiToRemove)
    if err != nil {
        // handle error
    }
    fmt.Println("ABI removal process initiated.")
    // Data will be updated in the facet, and next GetPage will reflect the change after reload.
}

// To refresh data if external changes are suspected:
abiCollection.Reset(abis.AbisDownloaded)
abiCollection.LoadData(abis.AbisDownloaded) // or wait for next GetPage
```

This package provides a comprehensive solution for managing ABI data, emphasizing efficiency and integration with the broader TrueBlocks application architecture.

## Runnable Example

Here is a full `main.go` file that you can use to run the example. Ensure you have a Go environment set up and the TrueBlocks dependencies available in your Go module.

```go
package main

import (
	"fmt"
	"time" // For brief sleeps to allow async operations in this example

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types" // For types.ListKind
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/abis"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"

	// In a full TrueBlocks application, you would also initialize:
	// "github.com/TrueBlocks/trueblocks-dalledress/pkg/logging"
	// "github.com/TrueBlocks/trueblocks-dalledress/pkg/preferences"
	// "github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	// For this example, we\'ll proceed without them, which means some
	// internal logging or preference-based behavior might use defaults or be silent.
)

// main function demonstrates the usage of the AbisCollection.
func main() {
	fmt.Println("--- Initializing ABIs Collection ---")
	// NewAbisCollection sets up facets with shared sources.
	// These sources might attempt to use preferences (e.g., for chain selection)
	// or logging, which are not explicitly initialized in this standalone example.
	abiCollection := abis.NewAbisCollection()

	// --- Initial Load and Display ---
	// We\'ll demonstrate with AbisDownloaded. Other kinds (AbisKnown, AbisFunctions, AbisEvents) work similarly.
	fmt.Println("\\n--- 1. Initial Load and Page Display (Downloaded ABIs) ---")
	// LoadData is asynchronous. GetPage will also trigger LoadData if the facet is stale.
	abiCollection.LoadData(abis.AbisDownloaded)
	// Crude delay for async operation; in a real app, use events or check facet state.
	fmt.Println("Called LoadData for AbisDownloaded. Waiting briefly for potential async processing...")
	time.Sleep(1 * time.Second)
	displayAbisPage(&abiCollection, abis.AbisDownloaded, "Initial Page")

	// --- CRUD Operation (Attempt to Remove an ABI) ---
	fmt.Println("\\n--- 2. CRUD Operation (Attempt to Remove First Downloaded ABI if any) ---")
	// Get the current page to find an ABI to remove.
	// We fetch just one item to see if there\'s anything to remove.
	pageDataForRemoval, err := abiCollection.GetPage(abis.AbisDownloaded, 0, 1, sdk.SortSpec{Field: "Name", Direction: sdk.SortAsc}, "")
	if err == nil && len(pageDataForRemoval.Abis) > 0 {
		abiToRemove := pageDataForRemoval.Abis[0] // Get the first ABI from the page
		fmt.Printf("Attempting to remove ABI: \'%s\' (Address: %s)\\n", abiToRemove.Name, abiToRemove.Address.Hex())

		// Perform the CRUD operation (Remove)
		err = abiCollection.AbisCrud(abis.AbisDownloaded, crud.Remove, &abiToRemove)
		if err != nil {
			fmt.Printf("Error during AbisCrud operation for \'%s\': %v\\n", abiToRemove.Name, err)
		} else {
			fmt.Printf("ABI removal process for \'%s\' initiated successfully via AbisCrud.\\n", abiToRemove.Name)
			// The underlying data source (chifra cache) should be updated.
			// The facet\'s internal data is also modified.
			// Wait a moment for changes to propagate or for cache updates.
			fmt.Println("Waiting briefly after removal attempt...")
			time.Sleep(2 * time.Second)
			fmt.Println("Displaying page after removal attempt to see changes...")
			displayAbisPage(&abiCollection, abis.AbisDownloaded, "Page After Removal Attempt")
		}
	} else if err != nil {
		fmt.Printf("Could not get an ABI to remove due to error: %v\\n", err)
	} else {
		fmt.Println("No downloaded ABIs found on the first page to attempt removal.")
	}

	// --- Reset and Reload Data ---
	fmt.Println("\\n--- 3. Reset and Reload Data (Downloaded ABIs) ---")
	abiCollection.Reset(abis.AbisDownloaded)
	fmt.Println("AbisDownloaded facet has been reset (marked stale).")
	// Calling LoadData again will trigger a fresh fetch from the source if it\'s stale.
	abiCollection.LoadData(abis.AbisDownloaded)
	fmt.Println("LoadData called for AbisDownloaded. Waiting briefly for potential async processing...")
	// Another crude delay.
	time.Sleep(2 * time.Second)

	fmt.Println("Displaying page after reset and reload attempt...")
	displayAbisPage(&abiCollection, abis.AbisDownloaded, "Page After Reset and Reload")

	fmt.Println("\\n--- Example Finished ---")
	fmt.Println("Important Notes:")
	fmt.Println("- This example uses time.Sleep() for simplicity. Real applications should use event-driven mechanisms or polling of facet states to manage asynchronous operations.")
	fmt.Println("- The behavior of data sources (especially caching and live fetches) depends on the TrueBlocks backend (chifra) configuration and state.")
	fmt.Println("- For this to fetch real data, your TrueBlocks environment must be correctly set up and running.")
}

// displayAbisPage is a helper function to fetch and print a page of ABI data.
func displayAbisPage(ac *abis.AbisCollection, listKind types.ListKind, contextLabel string) {
	fmt.Printf("\\n-- Displaying ABIs for ListKind: \'%s\' (Context: %s) --\\n", listKind, contextLabel)

	// Fetch a page (e.g., first 5 items, sorted by Name)
	page, err := ac.GetPage(listKind, 0, 5, sdk.SortSpec{Field: "Name", Direction: sdk.SortAsc}, "")
	if err != nil {
		fmt.Printf("Error getting page for \'%s\': %v\\n", listKind, err)
		return
	}

	fmt.Printf("Page Details - State: %s, IsFetching: %t, Total Items in Facet: %d, Expected Total from Source: %d\\n",
		page.State, page.IsFetching, page.TotalItems, page.ExpectedTotal)

	itemCount := 0
	if listKind == abis.AbisDownloaded || listKind == abis.AbisKnown {
		itemCount = len(page.Abis)
		if itemCount > 0 {
			fmt.Printf("Found %d ABIs on this page:\\n", itemCount)
			for i, abi := range page.Abis {
				fmt.Printf("  %d. Name: \'%s\', Address: %s, IsKnown: %t\\n",
					i+1, abi.Name, abi.Address.Hex(), abi.IsKnown)
			}
		}
	} else if listKind == abis.AbisFunctions || listKind == abis.AbisEvents {
		itemCount = len(page.Functions)
		if itemCount > 0 {
			fmt.Printf("Found %d Functions/Events on this page:\\n", itemCount)
			for i, fn := range page.Functions {
				fmt.Printf("  %d. Name: \'%s\', Type: %s, Encoding: %s, Signature: %s\\n",
					i+1, fn.Name, fn.FunctionType, fn.Encoding, fn.Signature)
			}
		}
	}

	if itemCount == 0 {
		fmt.Println("  No items found on this page for the given ListKind.")
	}
	fmt.Println("----------------------------------------------------")
}
```
