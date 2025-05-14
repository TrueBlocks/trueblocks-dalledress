import { useEffect, useMemo, useRef, useState } from 'react';

import { useTableContext } from '@components';
import { TableKey } from '@contexts';
import { sorting } from '@models';

import {
  Body,
  Column,
  Header,
  Pagination,
  PerPage,
  Stats,
  usePagination,
  useTableKeys,
} from '.';
import { SearchBox } from './SearchBox';
import './Table.css';

export interface TableProps<T extends Record<string, unknown>> {
  columns: Column<T>[];
  data: T[];
  loading: boolean;
  error: string | null;
  sort?: sorting.SortDef | null;
  onSortChange?: (sort: sorting.SortDef | null) => void;
  filter?: string;
  onFilterChange?: (filter: string) => void;
  tableKey: TableKey;
  onSaveRow?: (row: T, updated: Partial<T>) => void;
  onCancelRow?: () => void;
}

export function Table<T extends Record<string, unknown>>({
  columns,
  data,
  loading,
  error,
  sort: controlledSort,
  onSortChange,
  filter: controlledFilter,
  onFilterChange,
  tableKey,
  onSaveRow,
  onCancelRow,
}: TableProps<T>) {
  const { pagination } = usePagination(tableKey);
  const { currentPage, pageSize, totalItems } = pagination;
  const totalPages = Math.ceil(totalItems / pageSize);

  const {
    tableRef,
    selectedRowIndex,
    setSelectedRowIndex,
    focusTable,
    focusControls,
  } = useTableContext();

  const [expandedRowIndex, setExpandedRowIndex] = useState<number | null>(null);
  const [detailClosed, setDetailClosed] = useState(true);

  // Ref to track if the next Enter key should re-open the detail view
  const shouldOpenDetailOnEnter = useRef(false);

  const { handleKeyDown } = useTableKeys({
    itemCount: data.length,
    currentPage,
    totalPages,
    tableKey,
    expandedRowIndex,
    setExpandedRowIndex: (idx) => {
      if (idx !== null && shouldOpenDetailOnEnter.current) {
        setDetailClosed(false);
        shouldOpenDetailOnEnter.current = false;
      }
      if (idx === null) setDetailClosed(true);
      setExpandedRowIndex(idx);
    },
    onEnter: () => {
      shouldOpenDetailOnEnter.current = true;
    },
    onEscape: () => {
      setDetailClosed(true);
      setExpandedRowIndex(null);
    },
  });

  useEffect(() => {
    if (data.length > 0 && !loading) {
      const timer = setTimeout(() => {
        focusTable();
      }, 100);

      return () => clearTimeout(timer);
    }
  }, [data, loading, focusTable]);

  useEffect(() => {
    focusTable();
    // If the selected row is still present on the new page, keep it expanded
    if (
      !detailClosed &&
      selectedRowIndex >= 0 &&
      selectedRowIndex < data.length
    ) {
      setExpandedRowIndex(selectedRowIndex);
    } else {
      setExpandedRowIndex(null);
    }
  }, [currentPage, focusTable, selectedRowIndex, data.length, detailClosed]);

  // Only open the expanded row dialog when the selected row changes and detailClosed is false
  useEffect(() => {
    if (
      !detailClosed &&
      selectedRowIndex >= 0 &&
      selectedRowIndex < data.length
    ) {
      setExpandedRowIndex(selectedRowIndex);
    }
  }, [selectedRowIndex, data.length, detailClosed]);

  const handleRowClick = (index: number) => {
    setSelectedRowIndex(index);
    focusTable();
  };

  useEffect(() => {
    if (selectedRowIndex === -1 || selectedRowIndex >= data.length) {
      setSelectedRowIndex(Math.max(0, data.length - 1));
    }
  }, [data, selectedRowIndex, setSelectedRowIndex]);

  // Sorting state (uncontrolled if not provided)
  const [internalSort, setInternalSort] = useState<sorting.SortDef | null>(
    null,
  );
  const sort = controlledSort !== undefined ? controlledSort : internalSort;
  const handleSortChange = (nextSort: sorting.SortDef | null) => {
    if (onSortChange) onSortChange(nextSort);
    else setInternalSort(nextSort);
  };

  // Filter state (uncontrolled if not provided)
  const [internalFilter, setInternalFilter] = useState('');
  const filter =
    controlledFilter !== undefined ? controlledFilter : internalFilter;
  const handleFilterChange = (v: string) => {
    if (onFilterChange) onFilterChange(v);
    else setInternalFilter(v);
  };

  // Memoized sorted data (client-side for now)
  const sortedData = useMemo(() => {
    if (!sort || !sort.key) return data;
    const col = columns.find((c) => c.key === sort.key);
    if (!col) return data;
    // Avoid non-null assertion and explicit any
    const getValue = col.accessor
      ? (row: T) => (col.accessor ? col.accessor(row) : undefined)
      : (row: T) => {
          // Try to get the property by key, fallback to undefined
          return Object.prototype.hasOwnProperty.call(row, col.key)
            ? (row as Record<string, unknown>)[col.key]
            : undefined;
        };
    return [...data].sort((a, b) => {
      const va = getValue(a);
      const vb = getValue(b);
      if (va == null && vb == null) return 0;
      if (va == null) return sort.direction === 'asc' ? -1 : 1;
      if (vb == null) return sort.direction === 'asc' ? 1 : -1;
      if (typeof va === 'number' && typeof vb === 'number')
        return sort.direction === 'asc' ? va - vb : vb - va;
      return sort.direction === 'asc'
        ? String(va).localeCompare(String(vb))
        : String(vb).localeCompare(String(va));
    });
  }, [data, sort, columns]);

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <div
      className="table-outer-container"
      style={{
        display: 'flex',
        flexDirection: 'column',
        height: '100%',
        minHeight: 0,
        position: 'relative',
      }}
    >
      {/* Top: Pagination and Search */}
      <div
        className="top-pagination-container"
        style={{
          display: 'flex',
          alignItems: 'center',
          gap: 8,
          flex: '0 0 auto',
          zIndex: 2,
          background: '#23272f',
          borderBottom: '1px solid #333',
        }}
      >
        <SearchBox value={filter} onChange={handleFilterChange} />
        <PerPage
          tableKey={tableKey}
          pageSize={pageSize}
          focusTable={focusTable}
          focusControls={focusControls}
        />
        <Pagination
          totalPages={totalPages}
          currentPage={currentPage}
          tableKey={tableKey}
          focusControls={focusControls}
        />
      </div>

      {/* Table header (sticky) */}
      <div
        className="table-header-container"
        style={{
          flex: '0 0 auto',
          zIndex: 2,
          background: '#23272f',
          borderBottom: '1px solid #333',
        }}
      >
        <table
          className="data-table"
          style={{ width: '100%', tableLayout: 'fixed' }}
        >
          <Header
            columns={columns}
            sort={sort}
            onSortChange={handleSortChange}
          />
        </table>
      </div>

      {/* Scrollable body */}
      <div
        className="table-body-scroll"
        style={{
          flex: '1 1 auto',
          overflowY: 'auto',
          minHeight: 0,
          background: '#181818',
        }}
        onClick={focusTable}
      >
        <table
          className="data-table"
          ref={tableRef}
          tabIndex={0}
          onKeyDown={handleKeyDown}
          aria-label="Table with keyboard navigation"
          style={{ width: '100%', tableLayout: 'fixed' }}
        >
          {/* Empty header row for column widths */}
          <colgroup>
            {columns.map((col) => (
              <col key={col.key} style={{ width: col.width || 'auto' }} />
            ))}
          </colgroup>
          <tbody>
            {/* The Body component renders <tbody>, so render its content directly here */}
            <Body
              columns={columns}
              data={sortedData}
              selectedRowIndex={selectedRowIndex}
              handleRowClick={handleRowClick}
              noDataMessage="No data found."
              expandedRowIndex={expandedRowIndex}
              setExpandedRowIndex={(idx) => {
                if (idx === null) setDetailClosed(true);
                setExpandedRowIndex(idx);
              }}
              onSaveRow={onSaveRow || (() => setDetailClosed(true))}
              onCancelRow={onCancelRow || (() => setDetailClosed(true))}
            />
          </tbody>
        </table>
      </div>

      {/* Bottom: Stats */}
      <div
        className="table-footer"
        style={{
          flex: '0 0 auto',
          zIndex: 2,
          background: '#23272f',
          borderTop: '1px solid #333',
        }}
      >
        <Stats namesLength={data.length} tableKey={tableKey} />
      </div>
    </div>
  );
}
