import {
  ChangeEvent,
  FormEvent,
  useEffect,
  useMemo,
  useRef,
  useState,
} from 'react';

import { useTableContext, useTableKeys } from '@components';
import { Form } from '@components';
import { TableKey } from '@contexts';
import { Modal } from '@mantine/core';
import { sorting } from '@models';
import { Log } from '@utils';

import {
  Body,
  Column,
  Header,
  Pagination,
  PerPage,
  Stats,
  usePagination,
} from '.';
import { SearchBox } from './SearchBox';
import './Table.css';

export interface TableProps<T extends Record<string, unknown>> {
  columns: Column<T>[];
  data: T[];
  loading: boolean;
  sort?: sorting.SortDef | null;
  onSortChange?: (sort: sorting.SortDef | null) => void;
  filter?: string;
  onFilterChange?: (filter: string) => void;
  tableKey: TableKey;
  onSubmit: (data: Record<string, unknown>) => void;
  validate?: Record<
    string,
    (value: unknown, values: Record<string, unknown>) => string | null
  >;
}

export const Table = <T extends Record<string, unknown>>({
  columns,
  data,
  loading,
  sort: controlledSort,
  onSortChange,
  filter: controlledFilter,
  onFilterChange,
  tableKey,
  onSubmit,
  validate,
}: TableProps<T>) => {
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

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [currentRowData, setCurrentRowData] = useState<T | null>(null);

  const isModalOpenRef = useRef(false);

  useEffect(() => {
    isModalOpenRef.current = isModalOpen;
  }, [isModalOpen]);

  const closeModal = () => {
    setIsModalOpen(false);
    focusTable();
  };

  const handleFormSubmit = (e: FormEvent) => {
    e.preventDefault();
    const formData = new FormData(e.target as HTMLFormElement);
    const data = Object.fromEntries(formData.entries());
    // Log(`Table::handleFormSubmit: ${data['name']} ${data['address']}`);
    setIsModalOpen(false);
    onSubmit(data as T);
  };

  const { handleKeyDown } = useTableKeys({
    itemCount: data.length,
    currentPage,
    totalPages,
    tableKey,
    onEnter: () => {
      if (selectedRowIndex >= 0 && selectedRowIndex < data.length) {
        const rowData = data[selectedRowIndex] || null;
        setCurrentRowData(rowData);
        setIsModalOpen(true);
      }
    },
    onEscape: () => {
      setIsModalOpen(false);
    },
  });

  useEffect(() => {
    const tableElement = tableRef.current;
    if (tableElement) {
      const nativeKeyDownHandler = (e: KeyboardEvent) => {
        handleKeyDown({
          key: e.key,
          preventDefault: () => e.preventDefault(),
          stopPropagation: () => e.stopPropagation(),
        } as React.KeyboardEvent);
      };

      tableElement.addEventListener('keydown', nativeKeyDownHandler);
      return () => {
        tableElement.removeEventListener('keydown', nativeKeyDownHandler);
      };
    }
  }, [handleKeyDown, tableRef]);

  const handleFieldChange = (e: ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setCurrentRowData((prev) => {
      if (!prev) return null;
      return { ...prev, [name]: value };
    });
  };

  useEffect(() => {
    if (data.length > 0 && !loading) {
      const safeFocusTable = () => {
        if (!isModalOpenRef.current) {
          focusTable();
        }
      };

      const timer = setTimeout(() => {
        safeFocusTable();
      }, 100);

      return () => clearTimeout(timer);
    }
  }, [data, loading, focusTable]);

  useEffect(() => {
    if (isModalOpen) {
      const firstInput = document.querySelector(
        '.mantine-Modal input',
      ) as HTMLInputElement | null;
      firstInput?.focus();
    }
  }, [isModalOpen]);

  const handleFormKeyDown = (e: React.KeyboardEvent) => {
    const navigationKeys = [
      'ArrowUp',
      'ArrowDown',
      'ArrowLeft',
      'ArrowRight',
      'PageUp',
      'PageDown',
      'Home',
      'End',
    ];
    if (navigationKeys.includes(e.key)) {
      e.stopPropagation();
    }
  };

  useEffect(() => {
    focusTable();
  }, [currentPage, focusTable]);

  useEffect(() => {
    if (selectedRowIndex === -1 || selectedRowIndex >= data.length) {
      setSelectedRowIndex(Math.max(0, data.length - 1));
    }
  }, [data, selectedRowIndex, setSelectedRowIndex]);

  const [internalSort, setInternalSort] = useState<sorting.SortDef | null>(
    null,
  );
  const sort = controlledSort !== undefined ? controlledSort : internalSort;
  const handleSortChange = (nextSort: sorting.SortDef | null) => {
    if (onSortChange) onSortChange(nextSort);
    else setInternalSort(nextSort);
  };

  const [internalFilter, setInternalFilter] = useState('');
  const filter =
    controlledFilter !== undefined ? controlledFilter : internalFilter;
  const handleFilterChange = (v: string) => {
    if (onFilterChange) onFilterChange(v);
    else setInternalFilter(v);
  };

  const sortedData = useMemo(() => {
    if (!sort || !sort.key) return data;
    const col = columns.find((c) => c.key === sort.key);
    if (!col) return data;
    const getValue = col.accessor
      ? (row: T) => (col.accessor ? col.accessor(row) : undefined)
      : (row: T) => {
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

  const handleRowClick = (index: number) => {
    setSelectedRowIndex(index);
    setCurrentRowData(data[index] || null);
  };

  if (loading) return <div>Loading...</div>;

  return (
    <div className="table-outer-container">
      <div className="top-pagination-container">
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

      <div className="table-header-container">
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

      <div className="table-body-scroll" onClick={focusTable}>
        <table
          className="data-table"
          ref={tableRef}
          tabIndex={0}
          aria-label="Table with keyboard navigation"
          style={{ width: '100%', tableLayout: 'fixed' }}
        >
          <colgroup>
            {columns.map((col) => (
              <col key={col.key} style={{ width: col.width || 'auto' }} />
            ))}
          </colgroup>
          <tbody>
            <Body
              columns={columns}
              data={sortedData}
              selectedRowIndex={selectedRowIndex}
              handleRowClick={handleRowClick}
              noDataMessage="No data found."
            />
          </tbody>
        </table>
      </div>

      <div className="table-footer">
        <Stats namesLength={data.length} tableKey={tableKey} />
      </div>

      <Modal
        opened={isModalOpen}
        onClose={closeModal}
        centered
        size="lg"
        closeOnClickOutside={true}
        closeOnEscape={true}
        styles={{
          header: { display: 'none' },
          overlay: {
            backgroundColor: 'rgba(0, 0, 0, 0.5)', // Semi-transparent gray overlay
            backdropFilter: 'blur(1px)', // Optional slight blur effect
          },
          inner: {
            padding: '20px',
          },
        }}
      >
        <div onKeyDown={handleFormKeyDown}>
          <Form
            title={`Edit ${tableKey.tabName.replace(/\b\w/g, (char) =>
              char.toUpperCase(),
            )} ${tableKey.viewName
              .replace(/^\//, '')
              .replace(/\b\w/g, (char) => char.toUpperCase())}`}
            fields={columns.map((col) => ({
              name: col.key,
              label: col.header || col.key,
              placeholder: `Enter ${col.header || col.key}`,
              sameLine: col.sameLine,
              value: currentRowData
                ? String(
                    (currentRowData as Record<string, unknown>)[col.key] || '',
                  )
                : '',
            }))}
            onSubmit={handleFormSubmit}
            onCancel={closeModal}
            onChange={handleFieldChange}
            initMode="edit"
            compact
            validate={validate}
          />
        </div>
      </Modal>
    </div>
  );
};
