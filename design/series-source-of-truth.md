# DalleDress Series Source-of-Truth Design

## Goal

Establish `dalle/pkg/storage/series.tar.gz` as the canonical, embedded source of truth for the list of DalleDress series, using the same caching-aware replacement strategy already employed for `databases.tar.gz`.

## Current State (Before Implementation)

- The app loaded series from the runtime directory `~/.local/share/trueblocks/dalle/series/`.
- `dalle/pkg/storage/series.tar.gz` existed in the repo but was unused, truncated (29 bytes), and not referenced by any `//go:embed` directive.
- `dalle/pkg/storage/database.go` embedded `databases.tar.gz` and built an in-memory cache keyed by the archive hash.

## Final Design (Implemented)

1. **Embed the archive**: `dalle/pkg/storage/series.go` embeds `series.tar.gz`.
2. **In-memory cache**: On startup, `CacheManager` parses the embedded archive into `cache/series_v1.0.0.gob`, a `SeriesCache` of suffix -> raw JSON bytes.
3. **Hash-based invalidation**: The cache is rebuilt whenever the SHA256 hash of the embedded `series.tar.gz` differs from the stored `SourceHash`.
4. **User series isolation**: User-created series live in `~/.local/share/trueblocks/dalle/user-series/` and are merged with built-ins at query time.
5. **Immutable built-ins**: Built-in series cannot be edited, hidden, or deleted; the UI and backend both enforce this.
6. **Source of truth**: The embedded archive is the authoritative set of default series. Developer edits go through `dalle/pkg/storage/series/` followed by `make build-db`.

## Resolved Design Questions

- **Mutability**: Built-ins are immutable. Only the developer updates them by editing `dalle/pkg/storage/series/` and rebuilding the archive.
- **Overrides**: User-created series live in a separate `user-series/` directory. They shadow built-ins on suffix collision.
- **Replacement rule**: Hash-based, identical to `databases.tar.gz`.
- **Build process**: `make build-db` (existing target) regenerates `series.tar.gz` from `dalle/pkg/storage/series/`.
- **Migration**: None required; no user data existed in the runtime `series/` directory.
- **Windows**: Not supported.

## Related Files

- `dalle/pkg/storage/series.go` — `//go:embed series.tar.gz` and archive helpers.
- `dalle/pkg/storage/cache.go` — `CacheManager` series cache implementation.
- `dalle/pkg/storage/datadir.go` — `UserSeriesDir()` definition.
- `dalle/engine.go` — `ListSeries()`, `GetSeries()`, `SaveSeries()`, `SetSeriesHidden()`.
- `dalle/series.go` — `Series` struct with `Source` and `Version` fields.
- `dalle/series_crud.go` — user series merge helpers.
- `dalle/context.go` — `loadSeries()` using cache and user series.
- `dalledress/frontend/src/views/Series.tsx` — UI enforcement of built-in immutability.
- `dalle/pkg/storage/series/` — canonical source for built-in series JSON files.

## Acceptance Criteria

- [x] `series.tar.gz` is embedded in the binary.
- [x] A fresh install with no local user series receives the embedded set on first run.
- [x] Updating the app replaces stale built-in series while preserving user-created ones.
- [x] `ListSeries()` continues to return the same shape of data it does today, with the addition of a `source` field (`builtin` | `user`).
- [x] Built-in series are read-only in the UI.
- [x] Tests cover series archive parsing, cache build, cache reuse, invalidation, and hash-mismatch rebuild.

## See Also

- `dalle/design/embedded-series-data.md` for the detailed cache design.
- Existing embedded database handling in `dalle/pkg/storage/cache.go`.
