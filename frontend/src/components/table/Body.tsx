import React, { Fragment } from 'react';

import { Column } from '@components';

import './Body.css';

interface BodyProps<T extends Record<string, unknown>> {
  columns: Column<T>[];
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
            {columns.map((col) => (
              <td key={col.key} className={col.className}>
                {col.render
                  ? col.render(row, rowIndex)
                  : col.accessor
                    ? (col.accessor(row) as React.ReactNode)
                    : (row[col.key] as React.ReactNode)}
              </td>
            ))}
          </tr>
        </Fragment>
      ))}
    </>
  );
};
