// Package abis defines types and logic for managing contract ABIs, including CRUD operations and update tracking.
//
// Usage:
//
//	Import this package to define, store, and update contract ABI definitions.
//	Use the provided types and functions to perform CRUD operations and track ABI changes.
//
// Example:
//
//	ac := abis.NewAbisCollection()
//	page, err := ac.GetPage(abis.AbisDownloaded, 0, 10, sorting.EmptySortSpec(), "")
//	if err != nil {
//	    // handle error
//	}
//	for _, abi := range page.Abis {
//	    fmt.Println(abi.Name, abi.Address)
//	}
package abis
