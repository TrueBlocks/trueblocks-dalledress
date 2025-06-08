// Copyright (c) 2023 TrueBlocks
// Licensed under the MIT License
//
// Package facets provides a generic way to manage, load, filter, sort, and
// paginate data from various sources. It is designed to be flexible and
// extensible, allowing for different data types and source implementations.
//
// The core component is the Facet interface, which defines the contract for
// data access. BaseFacet provides a concrete implementation of this interface,
// utilizing a Source to fetch data.
//
// Key features include:
//   - Asynchronous data loading with state management (stale, fetching, loaded, etc.).
//   - Streaming of data during load operations.
//   - Paging, filtering, and sorting of data.
//   - Cache management and automatic updates when data becomes stale.
//   - CRUD operations on the underlying data.
package facets
