import React, { Fragment } from 'react';

import { FieldRenderer, FormField } from '@components';

import './Body.css';

interface BodyProps<T extends Record<string, unknown>> {
  columns: FormField<T>[];
  data: T[];
  selectedRowIndex: number;
  handleRowClick: (index: number) => void;
  noDataMessage?: string;
  onSubmit?: (data: Record<string, unknown>) => void;
}

export const Body = <T extends Record<string, unknown>>({
  columns,
  data,
  selectedRowIndex,
  handleRowClick,
  noDataMessage = 'No data found.',
  // eslint-disable-next-line unused-imports/no-unused-vars
  onSubmit,
}: BodyProps<T>) => {
  if (data.length === 0) {
    return (
      <tr className="selected">
        <td colSpan={columns.length} className="no-data-cell">
          {noDataMessage}
        </td>
      </tr>
    );
  }
  return (
    <>
      {data.map((row, rowIndex) => (
        <Fragment key={rowIndex}>
          <tr
            className={selectedRowIndex === rowIndex ? 'selected' : ''}
            onClick={() => handleRowClick(rowIndex)}
            aria-selected={selectedRowIndex === rowIndex}
          >
            {columns.map((col) => {
              const cellClassNames = [
                col.className || '',
                typeof col.width === 'string' && col.width.startsWith('col-')
                  ? col.width
                  : '',
              ]
                .filter(Boolean)
                .join(' ');

              return (
                <td
                  key={col.key}
                  className={cellClassNames}
                  style={
                    typeof col.width === 'string' &&
                    col.width.startsWith('col-')
                      ? { textAlign: col.textAlign || 'left' }
                      : {
                          ...(col.width ? { width: col.width } : undefined),
                          textAlign: col.textAlign || 'left',
                        }
                  }
                >
                  {col.render
                    ? col.render(row, rowIndex)
                    : col.accessor
                      ? (col.accessor(row) as React.ReactNode)
                      : col.key !== undefined
                        ? (() => {
                            const value = row[col.key as keyof T];
                            return (
                              <FieldRenderer
                                field={
                                  {
                                    type: col.type || 'text',
                                    value: value as string,
                                  } as FormField<Record<string, unknown>>
                                }
                                mode="display"
                                tableCell={true}
                              />
                            );
                          })()
                        : null}
                </td>
              );
            })}
          </tr>
        </Fragment>
      ))}
    </>
  );
};
