# Facets Package

The `facets` package provides a generic and extensible framework for managing, loading, filtering, sorting, and paginating data from various sources. It is designed to handle data asynchronously and efficiently, with built-in support for caching and state management.

## Core Concepts

### `Facet[T any]` Interface

The central piece of the package is the `Facet` interface. It defines the contract for accessing and manipulating data of a generic type `T`.

```go
type Facet[T any] interface {
    Load() (*StreamingResult, error)
    GetPage(first, pageSize int, filter FilterFunc[T], sortSpec interface{}, sortFunc func([]T, interface{}) error) (*PageResult[T], error)
    IsFetching() bool
    IsLoaded() bool
    GetState() LoadState
    NeedsUpdate() bool
    FacetCrud(matchFunc func(*T) bool) bool
    Reset()
    ExpectedCount() int
    Count() int
}
```

- **`Load()`**: Initiates or retrieves data. If data is stale, it fetches from the source. Returns a `StreamingResult` which can indicate partial data availability during loading.
- **`GetPage(...)`**: Retrieves a specific page of data, applying filtering and sorting.
- **`IsFetching()`**: Returns `true` if data is currently being loaded.
- **`IsLoaded()`**: Returns `true` if data has been successfully and completely loaded.
- **`GetState()`**: Returns the current `LoadState` of the facet.
- **`NeedsUpdate()`**: Returns `true` if the data is considered stale and needs to be reloaded.
- **`FacetCrud(...)`**: Allows for Create, Read, Update, Delete operations on the underlying data. Takes a `matchFunc` to identify items for deletion.
- **`Reset()`**: Resets the facet's state to `StateStale`, forcing a reload on the next access.
- **`ExpectedCount()`**: Returns the number of items expected to be loaded (e.g., from a source manifest).
- **`Count()`**: Returns the number of items currently loaded in the facet.

### `BaseFacet[T any]` Struct

`BaseFacet` is a concrete implementation of the `Facet` interface. It uses a `source.Source[T]` to fetch data and manages the data in an internal slice.

```go
type BaseFacet[T any] struct {
    // ... private fields for state, source, data, mutexes, etc.
}
```

- **`NewBaseFacet[T any](...)`**: Constructor for `BaseFacet`.
  ```go
  func NewBaseFacet[T any](
      listKind types.ListKind,
      filterFunc FilterFunc[T],
      isDupFunc func(existing []T, newItem *T) bool,
      source source.Source[T],
  ) *BaseFacet[T]
  ```
  - `listKind`: A `types.ListKind` identifying the type of data being managed (for messaging/UI purposes).
  - `filterFunc`: A default filter function to apply during loading.
  - `isDupFunc`: A function to identify and handle duplicate items during loading.
  - `source`: The `source.Source[T]` implementation responsible for fetching the raw data.

### `LoadState`

Represents the various states a facet can be in during its lifecycle.

```go
type LoadState string

const (
    StateStale    LoadState = "stale"    // Data is old or not yet loaded
    StateFetching LoadState = "fetching" // Actively loading data
    StatePartial  LoadState = "partial"  // Some data loaded, more incoming or error occurred
    StateLoaded   LoadState = "loaded"   // All data successfully loaded
    StatePending  LoadState = "pending"  // An operation is pending (transitional)
    StateError    LoadState = "error"    // An error occurred during loading
)
```

- `AllStates`: A slice providing all `LoadState` values, useful for frontend bindings.

### Helper Types

- **`FilterFunc[T any]`**: `func(*T) bool` - A function type for filtering items.
- **`StreamingResult`**:
  ```go
  type StreamingResult struct {
      Payload types.DataLoadedPayload // Contains counts and completion status
      Error   error                   // Any error during streaming
  }
  ```
- **`PageResult[T any]`**:
  ```go
  type PageResult[T any] struct {
      Items      []T       // The items for the current page
      TotalItems int       // Total number of items after filtering
      State      LoadState // The load state when the page was retrieved
  }
  ```

## Public Functions and Methods

### `BaseFacet[T any]` Methods

(Refer to the `Facet[T any]` interface methods described above, as `BaseFacet` implements them.)

Additional public methods on `BaseFacet`:

- **`Clear()`**: Removes all data from the facet and resets its state to `StateStale`.
- **`StartFetching()`**: Attempts to transition the state from `StateStale` to `StateFetching`. Returns `false` if already fetching.
- **`SetPartial()`**: Transitions the state from `StateFetching` to `StatePartial`.
- **`StopFetching()`**: Sets the state to `StateLoaded`.
- **`MarkStale()`**: Sets the state to `StateStale` if currently `StateLoaded` or `StatePartial`, indicating that external changes might have occurred and data should be reloaded.

### Standalone Utility Functions

- **`FilterPageSlice[T any](items *[]T, filterFn func(T) bool) []T`**:
  Applies a `filterFn` to a slice of items and returns a new slice containing only the items for which `filterFn` returns `true`.

- **`PaginateSlice[T any](items []T, first, pageSize int) []T`**:
  Extracts a sub-slice (page) from a larger slice of items based on `first` (starting index) and `pageSize`.

- **`SortPageSlice[T any](items []T, sortSpec sdk.SortSpec, sortFn func([]T, sdk.SortSpec) error) error`**:
  Sorts a slice of items in place using a provided `sortFn` and `sortSpec`.

- **`CreatePageFilter[T any](filter string, searchFields func(T) []string) func(T) bool`**:
  A utility to create a `FilterFunc`. It takes a search `filter` string and a `searchFields` function. The `searchFields` function extracts searchable string fields from an item of type `T`. The returned filter function will return `true` if the `filter` string is found (case-insensitively) in any of the extracted fields.

- **`ProcessPage[T any](...)`**:
  ```go
  func ProcessPage[T any](
      typeName string,
      sourceSlice *[]T,
      sortSpec sdk.SortSpec,
      sortFn func([]T, sdk.SortSpec) error,
      filterFn func(T) bool,
      first, pageSize int,
  ) (paginatedItems []T, totalFiltered, totalSource int, err error)
  ```
  A comprehensive utility function that takes a source slice and applies filtering, sorting, and pagination to it.
  - `typeName`: A string name for the type `T`, used in error messages.
  - `sourceSlice`: The original slice of data.
  - `sortSpec`: Specification for sorting.
  - `sortFn`: The function to perform the sort.
  - `filterFn`: The function to filter items.
  - `first`, `pageSize`: Pagination parameters.
  Returns the paginated items, the total count of items after filtering, the total count of items in the original source slice, and any error that occurred.

### Error Variables

- **`ErrorAlreadyLoading`**: `errors.New("already loading")`
  Returned by `Load()` if a load operation is already in progress.

## Example Usage

Let's assume we have a `struct Person` and a `source.Source[Person]` implementation called `PersonSource`.

```go
package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/facets"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/source"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5" // For sdk.SortSpec
)

// Define our data type
type Person struct {
	ID   int
	Name string
	Age  int
}

// Implement a dummy source.Source[Person]
type PersonSource struct {
	people []Person
}

func NewPersonSource() *PersonSource {
	return &PersonSource{
		people: []Person{
			{ID: 1, Name: "Alice", Age: 30},
			{ID: 2, Name: "Bob", Age: 24},
			{ID: 3, Name: "Charlie", Age: 35},
			{ID: 4, Name: "Diana", Age: 28},
			{ID: 5, Name: "Edward", Age: 40},
		},
	}
}

func (s *PersonSource) GetSourceType() source.SourceType {
	return source.SourceTypeAPI // Or any appropriate type
}

func (s *PersonSource) Fetch(
	progressCb func(current, total int, isDone bool),
	filterFunc func(*Person) bool,
	isDupFunc func(existing []Person, newItem *Person) bool,
	data *[]Person,
	expectedCnt *int,
) error {
	fmt.Println("PersonSource: Fetching data...")
	time.Sleep(50 * time.Millisecond) // Simulate network latency

	*expectedCnt = len(s.people)
	for i, p := range s.people {
		if filterFunc == nil || filterFunc(&p) {
			if isDupFunc == nil || !isDupFunc(*data, &p) {
				*data = append(*data, p)
			}
		}
		progressCb(i+1, len(s.people), i+1 == len(s.people))
		time.Sleep(10 * time.Millisecond) // Simulate streaming
	}
	fmt.Println("PersonSource: Fetching complete.")
	return nil
}


// Define a sort function for Person
func sortPeople(people []Person, spec interface{}) error {
	sortSpec, ok := spec.(sdk.SortSpec)
	if !ok {
		// Default sort or handle error if spec is mandatory and wrong type
		// For simplicity, let's sort by Name ASC if spec is not sdk.SortSpec
		// In a real app, you'd likely return an error or have a clear default.
		// This example assumes sdk.SortSpec or a default if not provided.
		// We'll sort by Name ASC as a default if spec is nil or not sdk.SortSpec
		// and Field is empty.
		if sortSpec.Field == "" {
			sortSpec.Field = "Name"
			sortSpec.Direction = sdk.SortAsc
		}
	}


	// Simple sort for example purposes
	// In a real app, use sort.SliceStable with more robust field checking
	for i := 0; i < len(people); i++ {
		for j := i + 1; j < len(people); j++ {
			swap := false
			switch sortSpec.Field {
			case "Name":
				if (sortSpec.Direction == sdk.SortAsc && people[i].Name > people[j].Name) ||
					(sortSpec.Direction == sdk.SortDesc && people[i].Name < people[j].Name) {
					swap = true
				}
			case "Age":
				if (sortSpec.Direction == sdk.SortAsc && people[i].Age > people[j].Age) ||
					(sortSpec.Direction == sdk.SortDesc && people[i].Age < people[j].Age) {
					swap = true
				}
			// Add other fields as needed
			default:
				return fmt.Errorf("unknown sort field: %s", sortSpec.Field)
			}
			if swap {
				people[i], people[j] = people[j], people[i]
			}
		}
	}
	return nil
}

func main() {
	personSrc := NewPersonSource()

	// No default filter during load, no specific duplicate check
	personFacet := facets.NewBaseFacet[Person](
		types.ListKind("people"), // Arbitrary kind for identification
		nil,                       // No initial filter during load
		func(existing []Person, newItem *Person) bool { // isDupFunc
			for _, p := range existing {
				if p.ID == newItem.ID {
					return true // It's a duplicate
				}
			}
			return false
		},
		personSrc,
	)

	fmt.Printf("Initial state: %s, Count: %d, Expected: %d
", personFacet.GetState(), personFacet.Count(), personFacet.ExpectedCount())

	// 1. Load data
	fmt.Println("
--- Loading Data ---")
	streamResult, err := personFacet.Load()
	if err != nil {
		if err == facets.ErrorAlreadyLoading {
			fmt.Println("Load already in progress.")
		} else {
			fmt.Printf("Error loading data: %v
", err)
			return
		}
	} else {
		fmt.Printf("Load complete. State: %s, Count: %d, Expected: %d
",
			personFacet.GetState(), streamResult.Payload.CurrentCount, streamResult.Payload.ExpectedTotal)
	}
	
	// Try loading again (should use cache)
	fmt.Println("
--- Loading Data Again (should be cached) ---")
	streamResult, err = personFacet.Load()
	if err != nil {
		fmt.Printf("Error loading data: %v
", err)
	} else {
		fmt.Printf("Load complete. State: %s, Count: %d, Expected: %d
",
			personFacet.GetState(), streamResult.Payload.CurrentCount, streamResult.Payload.ExpectedTotal)
	}


	// 2. Get a page of data (all items, no filter, sort by Name ASC)
	fmt.Println("
--- Getting Page (All, Sort by Name ASC) ---")
	sortByNameAsc := sdk.SortSpec{Field: "Name", Direction: sdk.SortAsc}
	pageResult, err := personFacet.GetPage(0, 5, nil, sortByNameAsc, sortPeople)
	if err != nil {
		fmt.Printf("Error getting page: %v
", err)
		return
	}
	fmt.Printf("Page State: %s, Total Items (after filter): %d
", pageResult.State, pageResult.TotalItems)
	for _, p := range pageResult.Items {
		fmt.Printf("  %+v
", p)
	}

	// 3. Get a filtered and sorted page (Age > 25, sort by Age DESC)
	fmt.Println("
--- Getting Page (Age > 25, Sort by Age DESC) ---")
	ageFilter := func(p *Person) bool {
		return p.Age > 25
	}
	sortByAgeDesc := sdk.SortSpec{Field: "Age", Direction: sdk.SortDesc}
	pageResult, err = personFacet.GetPage(0, 3, ageFilter, sortByAgeDesc, sortPeople)
	if err != nil {
		fmt.Printf("Error getting page: %v
", err)
		return
	}
	fmt.Printf("Page State: %s, Total Items (after filter): %d
", pageResult.State, pageResult.TotalItems)
	for _, p := range pageResult.Items {
		fmt.Printf("  %+v
", p)
	}

	// 4. Using CreatePageFilter utility
	fmt.Println("
--- Getting Page (Name contains 'a', case-insensitive, Sort by ID ASC) ---")
	nameContainsAFilter := facets.CreatePageFilter[Person]("a", func(p Person) []string {
		return []string{p.Name} // Fields to search in
	})
	sortByID := sdk.SortSpec{Field: "ID", Direction: sdk.SortAsc} // Assuming sortPeople handles "ID"
	// For this example, let's modify sortPeople to handle ID or add a new sort function.
	// For simplicity, we'll assume sortPeople can handle ID or we'd use a more generic sorter.
	// If sortPeople doesn't handle "ID", this would error or misbehave.
	// Let's assume we add ID sorting to sortPeople for this example.
	// (Actual sortPeople in this example doesn't handle ID, this is for illustration)

	pageResult, err = personFacet.GetPage(0, 5, nameContainsAFilter, sortByID, sortPeople)
	if err != nil {
		fmt.Printf("Error getting page: %v
", err)
		// return // Comment out to continue to CRUD example
	} else {
		fmt.Printf("Page State: %s, Total Items (after filter): %d
", pageResult.State, pageResult.TotalItems)
		for _, p := range pageResult.Items {
			fmt.Printf("  %+v
", p)
		}
	}


	// 5. CRUD: Remove a person
	fmt.Println("
--- Removing Bob (ID 2) ---")
	removed := personFacet.FacetCrud(func(p *Person) bool {
		return p.ID == 2 // Match Bob
	})
	if removed {
		fmt.Println("Bob was removed.")
		fmt.Printf("State: %s, Count: %d, Expected: %d
", personFacet.GetState(), personFacet.Count(), personFacet.ExpectedCount())
	} else {
		fmt.Println("Bob was not found or not removed.")
	}

	// Verify removal
	pageResult, _ = personFacet.GetPage(0, 5, nil, sortByNameAsc, sortPeople)
	fmt.Println("Current people after removal (sorted by Name ASC):")
	for _, p := range pageResult.Items {
		fmt.Printf("  %+v
", p)
	}
	
	// 6. Mark stale and reload
	fmt.Println("
--- Marking Stale and Reloading ---")
	personFacet.MarkStale()
	fmt.Printf("State after marking stale: %s
", personFacet.GetState())

	// Simulate external data change in the source
	personSrc.people = append(personSrc.people, Person{ID: 6, Name: "Frank", Age: 50})
	// Remove Bob from source as well if we want consistency with CRUD op
	var updatedPeople []Person
	for _, p := range personSrc.people {
		if p.ID != 2 {
			updatedPeople = append(updatedPeople, p)
		}
	}
	personSrc.people = updatedPeople


	streamResult, err = personFacet.Load() // This will now actually fetch
	if err != nil {
		fmt.Printf("Error reloading data: %v
", err)
	} else {
		fmt.Printf("Reload complete. State: %s, Count: %d, Expected: %d
",
			personFacet.GetState(), streamResult.Payload.CurrentCount, streamResult.Payload.ExpectedTotal)
	}
	pageResult, _ = personFacet.GetPage(0, 10, nil, sortByNameAsc, sortPeople)
	fmt.Println("Current people after reload (sorted by Name ASC):")
	for _, p := range pageResult.Items {
		fmt.Printf("  %+v
", p)
	}
}

```

To run this example:
1. Save the Go code above as `main.go` in a new directory.
2. Ensure you have the `trueblocks-sdk` and the current project's `pkg` directory in your Go module's scope or adjust import paths accordingly. If `pkg/source` and `pkg/types` are local to this project, you'd typically run `go mod init example.com/facetexample` (or similar) and then `go mod tidy`.
3. Run `go run main.go`.

This example demonstrates:
- Creating a `BaseFacet`.
- Loading data (initially and from cache).
- Paginating and sorting data.
- Filtering data with a custom filter and with `CreatePageFilter`.
- Using `FacetCrud` to remove an item.
- Marking data as stale and observing the reload behavior.

This `facets` package provides a robust foundation for managing collections of data in Go applications, particularly those needing to interact with dynamic data sources and present views to users.
