# TrueBlocks Dalledress - AI Development Instructions

## Project Architecture

This is a **Wails desktop application** with Go backend (`app/`) and React frontend (`frontend/src/`). The app interfaces with TrueBlocks Core SDK for blockchain data analysis and uses auto-generated TypeScript bindings from Go models.

### Core Components
- **Backend**: Go app in `app/` with API handlers (`api_*.go`) exposing TrueBlocks functionality
- **Frontend**: React/TypeScript with Mantine UI, organized by Views (`views/`) with Facets pattern
- **Models**: Auto-generated TypeScript types in `frontend/wailsjs/go/models.ts` from Go structs
- **Data Flow**: Frontend → Wails bindings → Go backend → TrueBlocks Core SDK

## ABSOLUTE REQUIREMENTS - ZERO TOLERANCE

### Package Management
- **YARN ONLY** - Never use `npm` or `npx`
- Run commands from repo root: `yarn start`, `yarn build`, `yarn test`
- Frontend commands run through root package.json, not `frontend/` directory

### Development Workflow
- **Always use `yarn start`** (Wails dev mode) - never browser mode or localhost:5173
- **Use `Log` from `@utils`** instead of console.log (console.log is invisible in Wails)
- **Read file contents first** before editing - files change between requests
- **No comments in production code** - only use comments for TODOs or explanations
- **Use `wails generate module`** to regenerate bindings after changes to the backend go code.
- **Always run the linter** from the root of the project with `yarn lint`.
- **Run linting, testing, and building** together with one command `yarn lint && yarn test && yarn start`. If any fail, stop and do not try to fix it. The developer will fix and provide guidance.
- **Removing files**: If you delete a file, use `rm -f` and ask for confirmation before proceeding. If you need to delete a folder, use `rm -R` (do not include the `-f` flag). Again, Ask for confirmation before proceeding.

### Code Patterns
- **Use existing patterns**: BaseTab, Table components, DataFacet enums, Collection/Store/Page architecture
- **Follow established architecture**: Import from `@models`, `@components`, `@utils`, `@hooks`
- **No React imports** (implicitly available)
- **No comments** in production code
- **Do not use `any` in TypeScript** - always use specific types (our linter won't allow it)
- **No custom strings for DataFacet** - always use `types.DataFacet.*` values
- **No custom ViewStateKey** - always use `{ viewName: string, tabName: types.DataFacet }`

## Header Actions Contract

Backend provides a single source of truth for header-level actions per facet, and the frontend must treat this as always-present (possibly empty).

### Backend
- Every facet config must define `HeaderActions []string` and it must never be nil. Use `[]` when there are no header actions for a facet.
- All data-table facets must include `export` in `HeaderActions`. A data-table facet is any facet with `isForm == false`.
- Action identifiers must align with the frontend `ActionType` values used by `useActions()` (e.g., `add`, `export`, `publish`, `pin`). Do not invent new strings.
- When you introduce a new action identifier, also expose metadata in the backend `Actions` map (title/label/icon) to keep UI consistent, then run `wails generate module`.
- Per-view notes (in addition to the universal `export` on all non-form facets):
  - Names: all facets include `add` (and `export` by the rule above). The `CUSTOM` facet additionally includes `publish` and `pin`.
  - Exports: all facets include `export` (already satisfied by the rule).
  - Monitors: include `export` (row actions like `delete`, `remove` remain as row-only).
  - Abis, Chunks, Status: include `export` on non-form facets (e.g., abis: all; chunks: stats/index/blooms; status: caches/chains). Form facets like `status` and `manifest` do not include `export`.
  - Contracts: table facets like `events` include `export`. Form facets like `dashboard` and `execute` do not include `export` and manage their own controls.

### Frontend
- Assume `config.headerActions` is always an array. Do not null-check it. Use length checks only, e.g., `if (!config.headerActions.length) return null`.
- Render header actions using the identifiers from `config.headerActions` and map them to handlers from `useActions()` (e.g., `handleAdd`, `handleExport`, `handlePublish`, `handlePin`).
- Do not hard-code action strings; rely on backend-provided identifiers and the `useActions()` wiring.
- If backend header actions change, regenerate models (`wails generate module`) and restart dev mode with `yarn start`.

### Reference: current facets by `isForm`

Data-table facets (must include `export`):
- abis: downloaded, known, functions, events
- names: all, custom, prefund, regular, baddress
- exports: statements, balances, transfers, transactions, withdrawals, assets, logs, traces, receipts
- monitors: monitors
- chunks: stats, index, blooms
- status: caches, chains
- contracts: events

Form facets (do not include `export`):
- status: status
- chunks: manifest
- contracts: dashboard, execute

## View Architecture Pattern

Each view follows this structure:
```
views/[viewname]/
├── [ViewName].tsx     # Main component with sections: Imports, Hooks, Handlers, Render
├── facets.ts          # DataFacet configurations and routing
├── columns.ts         # Table column definitions per facet
└── index.ts           # Exports
```

### Key Patterns
- **DataFacet enum**: Use `types.DataFacet.*` values, never custom strings
- **ViewStateKey**: `{ viewName: string, tabName: types.DataFacet }`
- **Page Data**: Auto-generated types like `monitors.MonitorsPage` with `facet`, `data[]`, `state`
- **BaseTab component**: Handles tables with `data`, `columns`, `viewStateKey`, `loading`, `error`

## Backend Integration

### API Endpoints
- Backend functions in `app/api_*.go` are auto-bound to frontend as `@app` imports
- Page data fetched via functions like `GetMonitorsPage(payload)`
- CRUD operations via `*Crud(action, data)` functions
- Error handling through `types.LoadState` enum

### Auto-Generated Models
- TypeScript types generated from Go structs in `frontend/wailsjs/go/models.ts`
- Import namespaced: `{ monitors, types, msgs } from '@models'`
- Never duplicate these types - always use generated ones

## Development Commands

```bash
# Development (from root)
yarn start              # Wails dev mode (NOT yarn dev)
yarn build             # Production build
yarn test              # Run all tests (Go + TypeScript + Dalle)
yarn lint              # Lint Go and TypeScript

# Testing specific components
yarn test-go           # Go backend tests
yarn test-tsx          # Frontend tests
yarn test-tsx <filename> # Run a single frontend test file
yarn test-dalle        # Dalle module tests
```

## Component Usage

### Tables
```tsx
import { BaseTab } from '@components';
import { getColumns } from './columns';

<BaseTab
  data={pageData?.monitors || []}
  columns={getColumns(getCurrentDataFacet())}
  viewStateKey={viewStateKey}
  loading={pageData?.isFetching || false}
  error={error}
  onSubmit={handleSubmit}
  onDelete={handleDelete}
/>
```

### Hooks
- `useActiveFacet()`: Manages view facets and routing
- `usePayload()`: Creates API request payloads
- `useActions()`: CRUD operations with error handling
- `useFiltering()`, `useSorting()`, `usePagination()`: Table state management

### Logging
```tsx
import { Log } from '@utils';
Log('Debug message here'); // Single string parameter only
```

### Address Handling
- **NEVER use manual address conversion** - use standardized utilities from `@utils`
- **Frontend ↔ Backend consistency**: Mirror Go patterns with TypeScript equivalents
- **Standard conversions**:
  ```tsx
  import { addressToHex, hexToAddress, getDisplayAddress, isValidAddress } from '@utils';
  
  // Convert address to hex string (equivalent to Go's addr.Hex())
  const hexString = addressToHex(address);
  
  // Convert hex string to address object (equivalent to Go's base.HexToAddress())
  const addressObj = hexToAddress('0x1234...');
  
  // Display truncated address (0x1234...abcd)
  const display = getDisplayAddress(address);
  
  // Validate non-zero address
  const valid = isValidAddress(address);
  ```
- **Forbidden patterns**:
  - `String(contract.address)` ❌
  - Manual byte array conversion ❌ 
  - `address as string` casting ❌
  - Custom address formatting functions ❌
  - Legacy function names `getAddressString`, `stringToAddress` ❌

## Error Handling Protocol

### Stop Conditions
- **Test failures**: Stop, report exact error, await instructions
- **Lint errors**: Stop, report issues, await fixes  
- **Build failures**: Stop, provide full output, await guidance
- **Unclear requirements**: Stop, ask specific questions

### Don't Guess
- Ask "Please clarify: [specific question]" instead of assuming
- Be honest about mock vs. real implementations
- Acknowledge when you don't understand something

### Anti-Bloat Principles
- **Use existing utilities** before creating new ones - check `@utils` first
- **Maintain consistency** with established patterns across the codebase
- **Remove redundant code** when standardizing - don't leave both old and new patterns
- **Check for existing solutions** before implementing custom logic
- **Follow the principle**: "If we've solved this problem once, reuse that solution"

## File Structure Context

```
app/                    # Go backend with TrueBlocks integration
├── api_*.go           # API handlers for each data collection
├── app.go             # Main Wails app struct and lifecycle
└── *.go               # Business logic and utilities

frontend/src/
├── components/        # Reusable UI components (BaseTab, Table, etc.)
├── views/            # Main application views with Facets pattern
├── hooks/            # Custom React hooks for data and state
├── contexts/         # React context providers
├── stores/           # Application state management
├── utils/            # Utilities including Log function
└── wailsjs/          # Auto-generated Wails bindings and types

pkg/                   # Go packages for backend functionality
dalle/                 # Separate Go module for AI/image generation
```

## Integration Points

- **TrueBlocks Core**: Backend integrates via SDK for blockchain data
- **Wails Bindings**: Auto-generated TypeScript/Go bridge
- **External APIs**: OpenAI integration in `dalle/` module
- **File System**: Project management and preferences via Go backend

Follow these patterns precisely. When in doubt, examine existing views like `monitors/` or `status/` for reference implementations.

## Wails Architecture

Go backend runs in the same process as the frontend
CGO bridge enables direct JavaScript ↔ Go function calls
Webview (like embedded Chromium) hosts the frontend
No separate processes - everything runs in a single process
Direct memory sharing between JS and Go (with serialization at the boundary) So it's even more impressive than IPC! The calls are:

- Synchronous from JS perspective (though Go functions can be async)
- Direct function invocation across the language boundary
- Minimal overhead compared to network calls or traditional IPC
- Shared memory space with serialization only at the JS/Go boundary
This makes the caching discussion even more interesting because:

- Calls are very fast - Direct CGO bridge, not IPC
- But serialization still exists - JSON-like marshaling between JS objects and Go structs
- Caching is still beneficial - Avoids repeated marshaling overhead
- Cache coherence is still the real problem - Backend transforms data without frontend awareness

## Information for future AI coders

- All main views follow a common structure: Imports → Hooks → Data fetching → Events → Actions (useActions) → Columns/forms → Render.
- The matrix below is the single source of truth for cross-view consistency. Treat it as the checklist when refactoring or adding views.
- Goal: a fully parameterized, generatable view scaffold driven by backend `ViewConfig`.
- Start by reading the matrix, then focus on the next non-green row (if any). Don’t revisit already-green rows unless new inconsistencies appear.
- Prefer the most recently refactored view as the reference implementation.
- Keep all changes lint- and test-compliant. Run `yarn lint && yarn test` regularly during refactors.

## Cross-view consistency matrix

Views in scope: Monitors, Names, Abis, Contracts, Status, Chunks, Exports

Legend (Key → Aspect):
- A: Form facet(s) present
- B: Form rendering abstracted (useFacetForm)
- C: Dynamic enabledActions per facet (backend-driven via ViewConfig)
- D: Confirm modal logic abstracted (centralized, per-operation silencing for remove/clean)
- E: Export format modal (unified, respects IsDialogSilenced)
- F: Inline facet-specific branching minimized inside perTabContent
- I: Header action mapping to handlers (config.headerActions → useActions handlers)

All Greens (fully consistent):

| Key | Monitors | Names | Abis | Contracts | Status | Chunks | Exports |
|-----|----------|-------|------|-----------|--------|--------|---------|
| A   | ✅        | ✅     | ✅    | ✅         | ✅      | ✅      | ✅       |
| B   | ✅        | ✅     | ✅    | ✅         | ✅      | ✅      | ✅       |
| C   | ✅        | ✅     | ✅    | ✅         | ✅      | ✅      | ✅       |
| D   | ✅        | ✅     | ✅    | ✅         | ✅      | ✅      | ✅       |
| E   | ✅        | ✅     | ✅    | ✅         | ✅      | ✅      | ✅       |
| F   | ✅        | ✅     | ✅    | ✅         | ✅      | ✅      | ✅       |
| I   | ✅        | ✅     | ✅    | ✅         | ✅      | ✅      | ✅       |

Notes:
- Header vs row actions are derived from backend `ViewConfig` per facet: `headerActions` vs `actions`, mapped through ACTION_DEFINITIONS and exposed by `useActions` as `config.headerActions` and `config.rowActions`.
- Confirm prompts only for remove and clean, using dialog keys `confirm.remove` and `confirm.clean` with per-operation silencing.
- Export uses ExportFormatModal unless `IsDialogSilenced('exportFormat')` is true.