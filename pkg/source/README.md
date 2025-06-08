# Package `sources`

Provides interfaces and implementations for fetching and streaming data from various sources, primarily the TrueBlocks SDK. It includes mechanisms for managing rendering contexts, processing data streams with features like filtering and deduplication, and reporting progress.

## Overview

The `sources` package is designed to abstract the complexities of data acquisition and provide a streamlined way to consume data streams within an application. It is particularly tailored for scenarios where data is fetched from external systems like the TrueBlocks SDK and needs to be processed, displayed, or managed with considerations for user experience, such as progress reporting and the ability to cancel operations.

## Core Components

-   `Source[T]`: This interface is the fundamental contract for all data sources. It defines a generic `Fetch` method that concrete source types must implement to retrieve data. The type parameter `T` allows sources to be strongly typed for the data they provide.

-   `SDKSource[T]`: A concrete implementation of the `Source[T]` interface. It is specifically designed to work with the TrueBlocks SDK. It takes a query function (which interacts with the SDK) and a processing function (to convert raw SDK output into the desired type `T`).

-   `ContextManager`: This component is responsible for managing `output.RenderCtx` instances. `RenderCtx` objects are crucial for controlling the lifecycle of data fetching operations. They encapsulate a Go `context.Context`, enabling features like cancellation, timeouts, and passing request-scoped values. The `ContextManager` provides functions to register, unregister, and retrieve these contexts, ensuring that operations can be cleanly started and stopped.

-   `ProcessStream[T]`: This is a high-level function that orchestrates the entire data streaming pipeline. It takes a `Source[T]`, a `contextKey` (to identify the operation), callback functions for filtering data (`filterFunc`) and detecting duplicates (`isDupFunc`), a target slice to store the results, an optional pointer to the expected total count of items, the kind of list being processed, a mutex for thread-safe access to the target slice, and a callback for when the first piece of data arrives. `ProcessStream` handles the interaction with the `Source`, manages the `RenderCtx` lifecycle via the `ContextManager`, and integrates with the `Progress` reporter.

-   `Progress`: This struct manages the logic for sending progress updates during data streaming. It helps in providing feedback to the user or system about the status of an ongoing data fetch operation. It can emit updates when the first chunk of data is received, at regular item count intervals, and also provide "heartbeat" updates if no new data has arrived for a certain period, assuring the user that the process is still active.

## Workflow and Usage

A typical workflow using this package involves:

1.  Creating an instance of a `Source[T]`, usually an `SDKSource[T]`, by providing it with a function that performs the actual SDK query and another function that processes each item returned by the SDK.
2.  Calling `ProcessStream[T]` with the created source, a unique `contextKey` for the operation, and various callback functions and parameters to control how data is filtered, stored, and how progress is reported.
3.  `ProcessStream` internally registers a new `RenderCtx` using the `ContextManager`.
4.  The `Source`'s `Fetch` method is called with this `RenderCtx`. Data items are streamed back through channels in the `RenderCtx`.
5.  As items arrive, `ProcessStream` applies the `filterFunc` and `isDupFunc`, appends valid items to the `targetSlice` (using the provided mutex for safety), and calls the `Progress` reporter's methods to potentially send updates.
6.  If the operation needs to be cancelled (e.g., user navigates away), the `CancelFetch` function of the `ContextManager` can be called with the `contextKey`, which will cancel the underlying `context.Context` in the `RenderCtx`, signaling the `Fetch` operation to stop.
7.  Once the `Fetch` operation completes (either normally, due to an error, or cancellation), `ProcessStream` unregisters the `RenderCtx`.

## Key Features

-   **Cancellability**: Operations can be cancelled cleanly through the `ContextManager` and `RenderCtx`, preventing resource leaks and unresponsive applications.
-   **Generic Data Handling**: The use of generics (`[T any]`) allows the package to work with any data type.
-   **Filtering and Deduplication**: Callers can provide custom logic to filter out unwanted items and prevent duplicates from being added to the result set.
-   **Progress Reporting**: Provides detailed feedback on the data loading process, enhancing user experience for potentially long-running operations. This includes initial load, incremental updates, and heartbeat messages.
-   **Thread Safety**: `ProcessStream` requires a mutex to be passed in to ensure that modifications to the shared `targetSlice` are thread-safe, which is important as data fetching often happens in separate goroutines.

## Example

(Conceptual example - actual SDK and type details would vary)

```go
// Assume coreTypes.Transaction is a defined type
var transactions []coreTypes.Transaction
var transactionsMutex sync.Mutex
expectedCount := 1000

// 1. Define the SDK query function
queryFn := func(ctx *output.RenderCtx) error {
    // Hypothetical SDK call
    // for i := 0; i < expectedCount; i++ {
    //     select {
    //     case <-ctx.Ctx.Done():
    //         return ctx.Ctx.Err()
    //     case ctx.ModelChan <- &coreTypes.Transaction{ID: fmt.Sprintf("tx-%d", i)}:
    //     }
    // }
    // close(ctx.ModelChan)
    return nil
}

// 2. Define the item processing function
processFn := func(itemIntf interface{}) *coreTypes.Transaction {
    if tx, ok := itemIntf.(*coreTypes.Transaction); ok {
        return tx
    }
    return nil
}

// 3. Create the SDKSource
txSource := sources.NewSDKSource(queryFn, processFn)

// 4. Define filter and duplicate check functions (optional)
filterFn := func(item *coreTypes.Transaction) bool {
    return true // Keep all items for this example
}
isDupFn := func(existing []coreTypes.Transaction, newItem *coreTypes.Transaction) bool {
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
_, err := sources.ProcessStream(
    contextKey,
    txSource,
    filterFn,
    isDupFn,
    &transactions,
    &expectedCount,
    "transactions", // types.ListKind
    &transactionsMutex,
    onFirstData,
)

if err != nil {
    log.Printf("Error streaming transactions: %v", err)
} else {
    log.Printf("Successfully streamed %d transactions", len(transactions))
}

// To cancel (e.g., from another goroutine or UI event):
// sources.CancelFetch(contextKey)
```
