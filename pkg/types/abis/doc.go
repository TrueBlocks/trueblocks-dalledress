// Copyright (c) 2023 TrueBlocks
// Licensed under the MIT License
//
// Package abis provides types and functions for managing ABI (Application Binary Interface)
// data within the TrueBlocks system. It leverages the `facets` package for efficient
// data handling, including loading, filtering, sorting, and pagination.
//
// The core of this package is the `AbisCollection`, which orchestrates multiple
// facets related to ABIs:
//   - Downloaded ABIs: ABIs fetched by the user.
//   - Known ABIs: ABIs that are part of the TrueBlocks curated list.
//   - Functions: ABI function definitions extracted from all ABIs.
//   - Events: ABI event definitions extracted from all ABIs.
//
// A key feature is the use of shared data sources (`source.Source`) to minimize
// redundant calls to the underlying TrueBlocks SDK. For example, both Downloaded
// and Known ABIs share a single source that lists all ABIs, and they are then
// filtered in memory. Similarly, Functions and Events share a source that provides
// detailed ABI information.
//
// This package defines list kinds (`AbisDownloaded`, `AbisKnown`, etc.) for use
// with the facet system and provides methods for loading data, retrieving paginated
// views, performing CRUD operations, and checking if data needs an update.
package abis
