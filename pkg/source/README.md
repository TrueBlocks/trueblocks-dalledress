# Package `source`

Provides interfaces and implementations for fetching and streaming data from various sources, primarily the TrueBlocks SDK. It includes mechanisms for managing rendering contexts, processing data streams with features like filtering and deduplication, and reporting progress.

## Overview

The `source` package is designed to abstract the complexities of data acquisition and provide a streamlined way to consume data streams within an application. It is particularly tailored for scenarios where data is fetched from external systems like the TrueBlocks SDK and needs to be processed, displayed, or managed with considerations for user experience, such as progress reporting and the ability to cancel operations.

## Core Components

-   `Source[T]`: This interface defines the contract for data sources.
    -   `Fetch(ctx *output.RenderCtx, processor func(item *T) bool) error`: Implementations of this method are responsible for retrieving data, processing each raw item into the type `*T` (if necessary), and then passing it to the `processor` callback. The `processor` callback (typically provided by `ProcessStream`) handles further logic like filtering, deduplication, and storage. If the `processor` returns `false`, the `Fetch` operation should stop. The method must also respect context cancellation via `ctx.Ctx.Done()`.
    -   `GetSourceType() string`: Returns a string identifier for the type of the source (e.g., "sdk").

-   `SDKSource[T]`: A concrete implementation of `Source[T]` designed for the TrueBlocks SDK.
    -   It is initialized with `NewSDKSource[T](queryFunc func(*output.RenderCtx) error, processFunc func(interface{}) *T)`.
        -   `queryFunc`: A function that performs the actual SDK query, sending raw data items to `ctx.ModelChan` and errors to `ctx.ErrorChan`.
        -   `processFunc`: A function that converts a raw data item (of type `interface{}`) received from `ctx.ModelChan` into the typed item `*T`.
    -   Its `Fetch` method orchestrates running the `queryFunc`, uses `processFunc` to convert items, and then passes them to the `processor` callback provided to `Fetch`.
    -   Its `GetSourceType()` method returns `"sdk"`.

-   `ContextManager`: This exported struct is responsible for managing `output.RenderCtx` instances. `RenderCtx` objects are crucial for controlling the lifecycle of data fetching operations, primarily enabling cancellation. The `ContextManager` is a singleton accessed via the exported function `GetContextManager()`. It provides the following exported functions (which are called on the instance returned by `GetContextManager()`):
    -   `RegisterContext(key string) *output.RenderCtx`: Registers a new context for a key, cancelling any existing context for the same key.
    -   `UnregisterContext(key string) (int, bool)`: Removes and cancels the context for a given key. Returns `1, true` if a context was found and removed, `0, false` otherwise.
    -   `CancelFetch(contextKey string)`: Cancels and removes the context associated with `contextKey`.
    -   `CtxCount(key string) int`: Returns 1 if a context for the key exists, 0 otherwise.

-   `ProcessStream[T]`: This high-level function orchestrates the data streaming pipeline.
    -   Signature: `ProcessStream[T](contextKey string, source Source[T], filterFunc func(item *T) bool, isDupFunc func(existing []T, newItem *T) bool, targetSlice *[]T, expectedCnt *int, listKind types.ListKind, m sync.Locker, onFirstData func()) (types.DataLoadedPayload, error)`
        (Note: `sync.Locker` is an interface satisfied by `*sync.Mutex`).
    -   It takes a `Source[T]`, a `contextKey`, callbacks for filtering (`filterFunc`) and deduplication (`isDupFunc` - should return `true` if `newItem` is a duplicate within `existing`), a pointer to the target slice, an optional pointer to the expected total count, the `types.ListKind`, a mutex (e.g., `*sync.Mutex`), and an `onFirstData` callback.
    -   It manages the `RenderCtx` lifecycle using the `ContextManager` (via `RegisterContext` and `UnregisterContext` called on the singleton instance).
    -   It calls the `source.Fetch` method, providing a processor callback that incorporates the `filterFunc`, `isDupFunc`, appends valid items to `targetSlice` (thread-safely using the mutex), and utilizes an internal (unexported) progress reporter to send updates.
    -   The internal progress reporter handles updates when the first chunk of data is received, at dynamic item count intervals, and provides "heartbeat" updates if no new data has arrived for a certain period.
    -   Returns a `types.DataLoadedPayload` containing the final counts and the `listKind`. An error is returned if the streaming failed *and* no data was retrieved. If partial data was retrieved before an error occurred, that data will be in `targetSlice`, the payload will reflect its count, and `nil` will be returned as the error by `ProcessStream`.

## Workflow and Usage

A typical workflow using this package involves:

1.  Creating an instance of a `Source[T]`, usually an `SDKSource[T]`, by providing it with a function that performs the actual SDK query and another function that processes each item returned by the SDK.
2.  Getting the singleton `ContextManager` instance via `source.GetContextManager()`.
3.  Calling `ProcessStream[T]` with the created source, a unique `contextKey` for the operation, and various callback functions and parameters to control how data is filtered, stored, and how progress is reported. `ProcessStream` will use the `ContextManager`'s methods (like `RegisterContext` and `UnregisterContext`) internally.
4.  `ProcessStream[T]` internally registers a new `RenderCtx` using `RegisterContext` (called on the `ContextManager` instance).
5.  The `Source`\'s `Fetch` method is called by `ProcessStream` with this `RenderCtx` and a processor callback. For `SDKSource`, data items are fetched by its `queryFunc` and sent through channels in the `RenderCtx`, then converted by its `processFunc`, and finally passed to the processor callback.
6.  The processor callback (within `ProcessStream`) applies the `filterFunc` and `isDupFunc`, appends valid items to the `targetSlice` (using the provided mutex for safety), and calls methods of an internal progress reporter to potentially send updates.
7.  If the operation needs to be cancelled (e.g., user navigates away), the `CancelFetch` function of the `ContextManager` can be called with the `contextKey` (on the singleton instance), which will cancel the underlying `context.Context` in the `RenderCtx`, signaling the `Fetch` operation to stop.
8.  Once the `Fetch` operation completes (either normally, due to an error, or cancellation), `ProcessStream[T]` unregisters the `RenderCtx` using `UnregisterContext` (on the `ContextManager` instance).

## Key Features

-   **Cancellability**: Operations can be cancelled cleanly through the `ContextManager` and `RenderCtx`, preventing resource leaks and unresponsive applications.
-   **Generic Data Handling**: The use of generics (`[T any]`) allows the package to work with any data type.
-   **Filtering and Deduplication**: Callers can provide custom logic to filter out unwanted items and prevent duplicates from being added to the result set.
-   **Progress Reporting**: `ProcessStream[T]` provides detailed feedback on the data loading process through an internal mechanism, enhancing user experience for potentially long-running operations. This includes initial load, incremental updates, and heartbeat messages.
-   **Thread Safety**: `ProcessStream[T]` requires a `sync.Locker` (like `*sync.Mutex`) to be passed in to ensure that modifications to the shared `targetSlice` are thread-safe, which is important as data fetching often happens in separate goroutines.

## Example

(Conceptual example - actual SDK and type details would vary)

```go
import (
    "fmt"
    "log"
    "sync"
    // Assuming output and types are correctly pathed
    "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
    "github.com/TrueBlocks/trueblocks-dalledress/pkg/source"
    "github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

// Assume coreTypes.Transaction is a defined type like:
// type Transaction struct { ID string; /* other fields */ }

func main() {
    // Assume coreTypes.Transaction is a defined type
    var transactions []types.Transaction // Using a placeholder type
    var transactionsMutex sync.Mutex
    expectedCount := 1000

    // 1. Define the SDK query function
    queryFn := func(ctx *output.RenderCtx) error {
        // Hypothetical SDK call
        // for i := 0; i < expectedCount; i++ {
        //     // Simulate receiving a transaction
        //     tx := &types.Transaction{ID: fmt.Sprintf("tx-%d", i)}
        //     select {
        //     case <-ctx.Ctx.Done():
        //         log.Println("Query function cancelled")
        //         return ctx.Ctx.Err()
        //     case ctx.ModelChan <- tx: // Send the raw item (can be any interface{})
        //     }
        // }
        // log.Println("Query function finished sending items")
        // Ensure ModelChan is closed by SDKSource.Fetch's defer, or close it here if appropriate
        return nil
    }

    // 2. Define the item processing function (interface{} to *Transaction)
    processFn := func(itemIntf interface{}) *types.Transaction {
        if tx, ok := itemIntf.(*types.Transaction); ok {
            return tx
        }
        // log.Printf("ProcessFn: could not convert item: %T", itemIntf)
        return nil
    }

    // 3. Create the SDKSource
    txSource := source.NewSDKSource(queryFn, processFn)

    // 4. Define filter and duplicate check functions (optional)
    filterFn := func(item *types.Transaction) bool {
        return true // Keep all items for this example
    }
    isDupFn := func(existing []types.Transaction, newItem *types.Transaction) bool {
        for _, ex := range existing {
            if ex.ID == newItem.ID {
                return true // It's a duplicate
            }
        }
        return false // Not a duplicate
    }

    // 5. Define a callback for when the first data arrives
    onFirstData := func() {
        fmt.Println("First batch of transactions received!")
    }

    // 6. Call ProcessStream
    contextKey := "myTransactionsFetch"
    // Get the context manager instance (though ProcessStream handles its usage internally)
    // ctxManager := source.GetContextManager() 

    payload, err := source.ProcessStream(
        contextKey,
        txSource,
        filterFn,
        isDupFn,
        &transactions,
        &expectedCount,
        types.ListKind("transactions"), // types.ListKind
        &transactionsMutex,
        onFirstData,
    )

    if err != nil {
        log.Printf("Error streaming transactions: %v (Note: partial data might be in transactions slice)", err)
    }
    log.Printf("Stream processing finished. Payload: %+v", payload)
    log.Printf("Successfully streamed %d transactions. Expected: %d", len(transactions), expectedCount)


    // To cancel (e.g., from another goroutine or UI event):
    // ctxManager := source.GetContextManager()
    // source.CancelFetch(contextKey) // CancelFetch is a top-level function now
    // log.Println("Cancellation requested for", contextKey)
}

// Placeholder for types.Transaction if not defined elsewhere for the example
// package types
// type Transaction struct {
//  ID string
// }
// type ListKind string
// type DataLoadedPayload struct {
//  CurrentCount  int
//  ExpectedTotal int
//  ListKind      ListKind
// }
```
