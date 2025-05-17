import React, { useState } from 'react';

import { Logger } from '@app';
import { Column, Form, FormField } from '@components';

interface BodyProps<T extends Record<string, unknown>> {
  columns: Column<T>[];
  data: T[];
  selectedRowIndex: number;
  handleRowClick: (index: number) => void;
  noDataMessage?: string;
  expandedRowIndex: number | null;
  setExpandedRowIndex: (index: number | null) => void;
  onSubmit?: (data: Record<string, unknown>) => void;
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
  onSubmit,
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
                    onSave={(updatedData) => {
                      Logger(
                        'DEBUGGING: onSubmit in body' +
                          JSON.stringify(updatedData),
                      );
                      if (onSubmit) onSubmit(updatedData);
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

const EditableRowForm = <T extends Record<string, unknown>>({
  row,
  columns,
  onClose,
  onSave,
}: {
  row: T;
  columns: Column<T>[];
  onClose: () => void;
  onSave?: (updatedData: Record<string, unknown>) => void;
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

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (onSave) onSave(formData);
    onClose();
  };

  const formFields: FormField<T>[] = columns
    .filter((col) => col.editable !== false)
    .map((col) => ({
      name: col.key,
      label: col.header,
      value: String(formData[col.key as keyof T] ?? ''),
      type: 'text',
      readOnly: col.readOnly,
    }));

  return (
    <div
      style={{
        maxWidth: 700,
        margin: '1rem auto',
        backgroundColor: '#23272f',
        border: '1px solid #333',
        borderRadius: 12,
        boxShadow: '0 2px 16px rgba(0,0,0,0.25)',
        padding: '1rem 2rem',
        color: '#e0e0e0',
      }}
    >
      <Form<T>
        fields={formFields}
        initMode="edit"
        onSubmit={handleSubmit}
        onCancel={onClose}
        onChange={handleChange}
        submitText="Save"
        compact={true}
      />
    </div>
  );
};
