import React from 'react';

import { FormField } from '@components';
import { TableKey, useSorting } from '@contexts';

import './Header.css';

export const Header = <T extends Record<string, unknown>>({
  columns,
  tableKey,
}: {
  columns: FormField<T>[];
  tableKey: TableKey;
}) => {
  const { sort, setSorting } = useSorting(tableKey);
  const handleClick = (col: FormField<T>) => {
    if (!col.sortable) return;
    if (!sort || sort.key !== col.key) {
      if (typeof col.key === 'string') {
        setSorting({ key: col.key, direction: 'asc' });
      }
    } else if (sort.direction === 'asc') {
      setSorting({ key: col.key, direction: 'desc' });
    } else {
      setSorting(null); // Remove sort
    }
  };

  return (
    <thead>
      <tr>
        {columns.map((col) => {
          const isSorted = sort && sort.key === col.key;
          // Determine aria-sort value safely
          let ariaSort: 'ascending' | 'descending' | undefined;
          if (isSorted) {
            ariaSort =
              sort && sort.direction === 'asc' ? 'ascending' : 'descending';
          }
          // Determine sort indicator safely
          let sortIndicator = '';
          if (col.sortable) {
            if (isSorted) {
              sortIndicator = sort && sort.direction === 'asc' ? ' ▲' : ' ▼';
            } else {
              sortIndicator = ' ⇅';
            }
          }
          const classNames = [
            col.className || '', // Custom class from column definition
            col.sortable ? 'sortable' : '',
            isSorted ? 'sorted' : '',
            typeof col.width === 'string' && col.width.startsWith('col-')
              ? col.width
              : '',
          ]
            .filter(Boolean)
            .join(' ');

          return (
            <th
              key={col.key}
              style={
                typeof col.width === 'string' && col.width.startsWith('col-')
                  ? undefined
                  : col.width
                    ? { width: col.width }
                    : undefined
              }
              className={classNames}
              onClick={col.sortable ? () => handleClick(col) : undefined}
              tabIndex={col.sortable ? 0 : undefined}
              aria-sort={ariaSort}
              role={col.sortable ? 'button' : undefined}
              onKeyDown={
                col.sortable
                  ? (e) => {
                      if (e.key === 'Enter' || e.key === ' ') {
                        e.preventDefault();
                        handleClick(col);
                      }
                    }
                  : undefined
              }
            >
              <span className="header-cell">
                <span className="header-label">{col.header}</span>
                {col.sortable && (
                  <span className="sort-indicator">{sortIndicator}</span>
                )}
              </span>
            </th>
          );
        })}
      </tr>
    </thead>
  );
};
