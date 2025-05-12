import { useEffect, useMemo, useState } from 'react';

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

export interface TableProps<T> {
  columns: Column<T>[];
  data: T[];
  loading: boolean;
  error: string | null;
  sort?: sorting.SortDef | null;
  onSortChange?: (sort: sorting.SortDef | null) => void;
  filter?: string;
  onFilterChange?: (filter: string) => void;
  tableKey: TableKey;
}

export function Table<T>({
  columns,
  data,
  loading,
  error,
  sort: controlledSort,
  onSortChange,
  filter: controlledFilter,
  onFilterChange,
  tableKey,
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

  const { handleKeyDown } = useTableKeys({
    itemCount: data.length,
    currentPage,
    totalPages,
    tableKey,
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
  }, [currentPage, focusTable]);

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
    <>
      <div
        className="top-pagination-container"
        style={{ display: 'flex', alignItems: 'center', gap: 8 }}
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

      <div className="table-container" onClick={focusTable}>
        <table
          className="data-table"
          ref={tableRef}
          tabIndex={0}
          onKeyDown={handleKeyDown}
          aria-label="Table with keyboard navigation"
        >
          <Header
            columns={columns}
            sort={sort}
            onSortChange={handleSortChange}
          />
          <Body
            columns={columns}
            data={sortedData}
            selectedRowIndex={selectedRowIndex}
            handleRowClick={handleRowClick}
            noDataMessage="No data found."
          />
        </table>
      </div>

      <div className="table-footer">
        <Stats namesLength={data.length} tableKey={tableKey} />
      </div>
    </>
  );
}
