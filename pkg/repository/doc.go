// Package repository provides data persistence, retrieval, and caching mechanisms for application models.
//
// Usage:
//
//	Import this package to interact with persistent storage, perform CRUD operations, and manage in-memory caches for models.
//	Typical usage involves calling repository functions to load, save, or update data entities.
//
// Example:
//
//	repo := repository.NewBaseRepository(
//	    types.ListKind("example"),
//	    filterFunc, processFunc, queryFunc, dedupeFunc,
//	)
//	result, err := repo.Load(repository.LoadOptions{})
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println(result.Status)
package repository
