// Package streaming manages real-time and batched data streaming, rendering contexts, and page state for the application.
//
// Usage:
//
//	Import this package to handle live or paginated data flows, manage rendering contexts, and coordinate updates to UI or consumers.
//	Use streaming functions to subscribe to data, render pages, or update state in response to events.
//
// Example:
//
//	status, payload, err := streaming.StreamData(
//	    "context-key", queryFunc, filterFunc, processFunc, dedupeFunc,
//	    &targetSlice, &expectedCount, &loadedFlag, listKind, &mutex,
//	)
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println(status, payload.CurrentCount)
package streaming
