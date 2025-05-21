# Table Component Design (Mermaid Diagrams & Narrative)

This document explains the structure and flow of the table system in the codebase, focusing on how data, events, and rendering are managed. It is intended for contributors and maintainers who want to understand or extend the table components. (For form components, see the README in the form folder.)

---

## 1. File/Component Relationship

This chart shows which files depend on or use which others. It helps you see the modularity and reusability of the table system.

```mermaid
graph TD
  Table --> Body
  Table --> Header
  Table --> Pagination
  Table --> PerPage
  Table --> Stats
  Table --> SearchBox
  Table --> TableContext
  Table --> usePagination
  Table --> useTableKeys
  Table --> Column
  Body --> Column
  Header --> Column
  Pagination --> usePagination
  PerPage --> usePagination
  Stats --> usePagination
  TableContext --> useTableContext
  SearchBox --> useFormHotkeys
  SearchBox --> useTableContext
  useTableKeys --> useTableContext
  useTableKeys --> usePagination
  index.tsx --> Table
  index.tsx --> Body
  index.tsx --> Header
  index.tsx --> Pagination
  index.tsx --> PerPage
  index.tsx --> Stats
  index.tsx --> TableContext
  index.tsx --> Column
  index.tsx --> usePagination
  index.tsx --> useTableKeys
```

---

## 2. Data & Event Lifecycle: How a Table Works

This chart and narrative walk you through what happens from the moment you provide `columns` and `data` to the moment a user interacts with the table.

```mermaid
sequenceDiagram
  participant Parent as ParentComponent
  participant Table as Table
  participant Body as Body
  participant Header as Header
  participant Pagination as Pagination
  participant PerPage as PerPage
  participant Stats as Stats
  Parent->>Table: Pass columns, data, handlers
  Table->>Header: Render column headers (with sort)
  Table->>Body: Render rows (with columns)
  Table->>Pagination: Render page controls
  Table->>PerPage: Render page size selector
  Table->>Stats: Render summary
  loop For each row
    Body->>Column: Render cell (accessor/render)
  end
  User->>Header: Click to sort
  Header->>Table: onSortChange
  User->>Pagination: Change page
  Pagination->>Table: goToPage
  User->>PerPage: Change page size
  PerPage->>Table: changePageSize
  User->>Body: Click/select row
  Body->>TableContext: setSelectedRowIndex
```

**Narrative:**

- The parent component provides `columns`, `data`, and handlers to `Table`.
- `Table` renders the header, body, pagination, per-page selector, and stats.
- Each row and cell is rendered using the `Column` definition (accessor/render).
- User actions (sort, paginate, select) flow through the appropriate subcomponents and update the table context or pagination state.

---

## 3. Column Structure (What is a Column?)

A `Column` is a flexible description of a single table column, including how to render and access its data.

```mermaid
classDiagram
  class Column {
    string key
    string header
    string label
    string type
    any value
    bool required
    string error
    function onChange
    function onBlur
    Column[] fields
    ReactNode customRender
    bool sortable
    function accessor
    function render
    %% ... (other props)
  }
```

---

## 4. Rendering Logic: Table Composition

This chart shows how the system supports both simple and complex tables, with sorting, pagination, and selection.

```mermaid
flowchart TD
  Table -->|"columns, data"| Header
  Table -->|"columns, data"| Body
  Table -->|"pagination"| Pagination
  Table -->|"pageSize"| PerPage
  Table -->|"summary"| Stats
  Header -->|"onSortChange"| Table
  Pagination -->|"goToPage"| Table
  PerPage -->|"changePageSize"| Table
  Body -->|"setSelectedRowIndex"| TableContext
```

---

## 5. Event Flow: User Interaction

This chart focuses on how user actions propagate through the system.

```mermaid
sequenceDiagram
  participant User
  participant Header
  participant Body
  participant Pagination
  participant PerPage
  participant TableContext
  User->>Header: Click to sort
  Header->>Table: onSortChange
  User->>Pagination: Change page
  Pagination->>Table: goToPage
  User->>PerPage: Change page size
  PerPage->>Table: changePageSize
  User->>Body: Click/select row
  Body->>TableContext: setSelectedRowIndex
```

---

## 6. Keyboard Navigation & Focus

The table system supports keyboard navigation and focus management for accessibility and speed.

```mermaid
flowchart TD
  User -->|"arrow keys, enter, esc"| useTableKeys
  useTableKeys --> TableContext
  useTableKeys --> usePagination
  TableContext --> Table
```

---

## Summary

- **Table** is the orchestrator: it manages state, rendering, and delegates to subcomponents.
- **Header, Body, Pagination, PerPage, Stats** are the visualizers: they know how to display their part of the table.
- **Column** is the schema: it describes what each column is and how it should behave.
- **TableContext** is the state manager: it tracks selection and focus.
- **usePagination, useTableKeys** are the hooks: they manage pagination and keyboard navigation.

This system is designed for flexibility, composability, and a clear separation of concerns.

---

## See Also

- For form components and their design, see the README in the `form` folder.
- This table system is independent of any wizard-related files/components.
