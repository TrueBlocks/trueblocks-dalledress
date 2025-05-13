import React, { useState } from 'react';

import { Column } from '@components';
import { TextInput } from '@mantine/core';

interface BodyProps<T extends Record<string, unknown>> {
  columns: Column<T>[];
  data: T[];
  selectedRowIndex: number;
  handleRowClick: (index: number) => void;
  noDataMessage?: string;
  expandedRowIndex: number | null;
  setExpandedRowIndex: (index: number | null) => void;
  onSaveRow?: (row: T, updated: Partial<T>) => void; // optional save handler
  onCancelRow?: () => void; // optional cancel handler
}

export function Body<T extends Record<string, unknown>>({
  columns,
  data,
  selectedRowIndex,
  handleRowClick,
  noDataMessage = 'No data found.',
  expandedRowIndex,
  setExpandedRowIndex,
  onSaveRow,
  onCancelRow,
}: BodyProps<T>) {
  if (data.length === 0) {
    return (
      <tr className="selected">
        <td colSpan={columns.length} style={{ textAlign: 'center' }}>
          {noDataMessage}
        </td>
      </tr>
    );
  }
  return (
    <>
      {data.map((row, rowIndex) => (
        <React.Fragment key={rowIndex}>
          <tr
            className={selectedRowIndex === rowIndex ? 'selected' : ''}
            onClick={() => handleRowClick(rowIndex)}
            aria-selected={selectedRowIndex === rowIndex}
          >
            {columns.map((col) => (
              <td key={col.key} className={col.className}>
                {col.render
                  ? col.render(row, rowIndex)
                  : col.accessor
                    ? (col.accessor(row) as React.ReactNode)
                    : ((row as Record<string, unknown>)[
                        col.key
                      ] as React.ReactNode)}
              </td>
            ))}
          </tr>
          {expandedRowIndex === rowIndex && (
            <tr>
              <td colSpan={columns.length}>
                <div
                  style={{
                    padding: '1rem',
                  }}
                >
                  <EditableRowForm
                    row={row}
                    columns={columns}
                    onClose={() => {
                      setExpandedRowIndex(null);
                      if (onCancelRow) onCancelRow();
                    }}
                    onSave={(row, updated) => {
                      if (onSaveRow) onSaveRow(row, updated);
                      setExpandedRowIndex(null);
                    }}
                  />
                </div>
              </td>
            </tr>
          )}
        </React.Fragment>
      ))}
    </>
  );
}

// Editable form for a row
const EditableRowForm = <T extends Record<string, unknown>>({
  row,
  columns,
  onClose,
  onSave,
}: {
  row: T;
  columns: Column<T>[];
  onClose: () => void;
  onSave?: (row: T, updated: Partial<T>) => void;
}) => {
  // Re-initialize formData whenever the row changes
  const [formData, setFormData] = useState<Partial<T>>(() => {
    const initial: Partial<T> = {};
    columns.forEach((col) => {
      if (col.accessor) {
        try {
          initial[col.key as keyof T] = col.accessor(row) as T[keyof T];
        } catch {
          initial[col.key as keyof T] = row[col.key as keyof T];
        }
      } else {
        initial[col.key as keyof T] = row[col.key as keyof T];
      }
    });
    return initial;
  });

  // Update formData if row changes
  React.useEffect(() => {
    const initial: Partial<T> = {};
    columns.forEach((col) => {
      if (col.accessor) {
        try {
          initial[col.key as keyof T] = col.accessor(row) as T[keyof T];
        } catch {
          initial[col.key as keyof T] = row[col.key as keyof T];
        }
      } else {
        initial[col.key as keyof T] = row[col.key as keyof T];
      }
    });
    setFormData(initial);
  }, [row, columns]);

  const handleChange = (key: string, value: unknown) => {
    setFormData((prev) => ({ ...prev, [key]: value }));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (onSave) onSave(row, formData);
    onClose();
  };

  return (
    <div
      style={{
        maxWidth: 700,
        margin: '2rem auto',
        backgroundColor: '#23272f',
        border: '1px solid #333',
        borderRadius: 12,
        boxShadow: '0 2px 16px rgba(0,0,0,0.25)',
        padding: '2.5rem 2rem',
        color: '#e0e0e0',
      }}
    >
      <form onSubmit={handleSubmit}>
        {columns.map((col) => (
          <div
            key={col.key}
            style={{
              marginBottom: '1.5rem',
              display: 'flex',
              alignItems: 'center',
              gap: '1rem',
            }}
          >
            <label
              htmlFor={col.key}
              style={{
                minWidth: 120,
                fontWeight: 500,
                color: '#e0e0e0',
                marginRight: 8,
              }}
            >
              {col.header}
            </label>
            <TextInput
              id={col.key}
              value={String(formData[col.key as keyof T] ?? '')}
              onChange={(e) => handleChange(col.key, e.target.value)}
              name={col.key}
              style={{ flex: 1, marginBottom: 0 }}
            />
          </div>
        ))}
        <div style={{ display: 'flex', justifyContent: 'flex-end', gap: 12 }}>
          <button
            type="button"
            onClick={onClose}
            style={{
              minWidth: 100,
              background: 'transparent',
              color: '#7fffa9',
              border: '1px solid #7fffa9',
              borderRadius: 6,
              padding: '0.5rem 1.5rem',
              fontWeight: 500,
              cursor: 'pointer',
              transition: 'background 0.2s, color 0.2s',
            }}
          >
            Cancel
          </button>
          <button
            type="submit"
            style={{
              minWidth: 100,
              background: '#7fffa9',
              color: '#23272f',
              border: 'none',
              borderRadius: 6,
              padding: '0.5rem 1.5rem',
              fontWeight: 500,
              cursor: 'pointer',
              transition: 'background 0.2s, color 0.2s',
            }}
          >
            Save
          </button>
        </div>
      </form>
    </div>
  );
};
