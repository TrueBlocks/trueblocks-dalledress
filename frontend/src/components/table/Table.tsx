import { ChangeEvent, useEffect, useRef, useState } from 'react';

import { useTableContext, useTableKeys } from '@components';
import { Form, FormField } from '@components';
import { TableKey } from '@contexts';
import { Modal } from '@mantine/core';
import { sorting } from '@models';

import { Body, Header, Pagination, PerPage, Stats, usePagination } from '.';
import { SearchBox } from './SearchBox';
import './Table.css';
import './TableSizing.css';

export interface TableProps<T extends Record<string, unknown>> {
  columns: FormField<T>[];
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
  collectionIsLoading?: boolean;
  collectionIsLoaded?: boolean;
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
  collectionIsLoading,
  collectionIsLoaded,
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

  const handleModalFormSubmit = (values: T) => {
    setIsModalOpen(false);
    onSubmit(values);
    focusTable();
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

  const handleRowClick = (index: number) => {
    setSelectedRowIndex(index);
    setCurrentRowData(data[index] || null);
  };

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
              <col
                key={col.key}
                className={
                  typeof col.width === 'string' && col.width.startsWith('col-')
                    ? col.width
                    : undefined
                }
                style={
                  typeof col.width === 'string' && col.width.startsWith('col-')
                    ? undefined
                    : { width: col.width || 'auto' }
                }
              />
            ))}
          </colgroup>
          <tbody>
            {loading && data.length === 0 ? (
              <tr>
                <td
                  colSpan={columns.length}
                  style={{
                    textAlign: 'left',
                    padding: '20px',
                    color: 'red',
                  }}
                >
                  Loading...
                </td>
              </tr>
            ) : (
              <>
                <Body
                  columns={columns}
                  data={data}
                  selectedRowIndex={selectedRowIndex}
                  handleRowClick={handleRowClick}
                  noDataMessage="No data found."
                />
              </>
            )}
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
          content: {
            // Target the modal content itself
            boxShadow: '0 4px 12px rgba(0, 0, 0, 0.15)', // Add a subtle box shadow
            border: '1px solid var(--mantine-color-gray-3)', // Add a light border, theme-aware
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
            fields={columns.map((col) => {
              const fieldDefinition: FormField<T> = {
                ...col,
                name: col.name || col.key,
                label: col.label || col.header || col.name || col.key,
                placeholder:
                  col.placeholder ||
                  `Enter ${col.label || col.header || col.name || col.key}`,
                value:
                  currentRowData && (col.name || col.key) !== undefined
                    ? String(
                        (currentRowData as Record<string, unknown>)[
                          (col.name || col.key) as string
                        ] ?? '',
                      )
                    : '',
              };
              return fieldDefinition;
            })}
            onSubmit={handleModalFormSubmit}
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
