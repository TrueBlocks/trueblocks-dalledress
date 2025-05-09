import React from 'react';

import { Column } from '@components';

// BodyProps defines the props for the TableBody component.
interface BodyProps<T> {
  columns: Column<T>[];
  data: T[];
  selectedRowIndex: number;
  handleRowClick: (index: number) => void;
  noDataMessage?: string;
}

// TableBody renders the table body rows for the given data and manages row selection.
export function Body<T>({
  columns,
  data,
  selectedRowIndex,
  handleRowClick,
  noDataMessage = 'No data found.',
}: BodyProps<T>) {
  return (
    <tbody>
      {data.length === 0 ? (
        <tr className="selected">
          <td colSpan={columns.length} style={{ textAlign: 'center' }}>
            {noDataMessage}
          </td>
        </tr>
      ) : (
        data.map((row, rowIndex) => (
          <tr
            key={rowIndex}
            className={selectedRowIndex === rowIndex ? 'selected' : ''}
            onClick={() => handleRowClick(rowIndex)}
            aria-selected={selectedRowIndex === rowIndex}
          >
            {columns.map((col) => (
              <td key={col.key} className={col.className}>
                {col.render
                  ? col.render(row, rowIndex)
                  : col.accessor
                    ? col.accessor(row)
                    : // fallback: try to display property by key if exists
                      // eslint-disable-next-line @typescript-eslint/no-explicit-any
                      (row as any)[col.key]}
              </td>
            ))}
          </tr>
        ))
      )}
    </tbody>
  );
}
