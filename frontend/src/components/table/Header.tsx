import React from 'react';

import { FormField } from '@components';
import { TableKey, useSorting } from '@contexts';

import {
  createSingleFieldSortSpec,
  getSortDirection,
  getSortField,
} from '../../utils/sortSpec';
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
    const field = getSortField(sort);
    if (!col.sortable) return;
    if (field !== col.key) {
      if (typeof col.key === 'string') {
        setSorting(createSingleFieldSortSpec(col.key, 'asc'));
      }
    } else if (getSortDirection(sort) === 'asc') {
      setSorting(createSingleFieldSortSpec(col.key as string, 'desc'));
    } else {
      setSorting(null);
    }
  };

  return (
    <thead>
      <tr>
        {columns.map((col) => {
          const field = getSortField(sort);
          const isSorted = field === col.key;
          // Determine aria-sort value
          let ariaSort: 'ascending' | 'descending' | undefined;
          if (isSorted) {
            ariaSort =
              getSortDirection(sort) === 'asc' ? 'ascending' : 'descending';
          }
          // Determine sort indicator
          let sortIndicator = '';
          if (col.sortable) {
            if (isSorted) {
              sortIndicator = getSortDirection(sort) === 'asc' ? ' ▲' : ' ▼';
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
