import React from 'react';

import { sorting } from '@models';

import type { Column } from '.';
import './Header.css';

export function Header<T extends Record<string, unknown>>({
  columns,
  sort,
  onSortChange,
}: {
  columns: Column<T>[];
  sort?: sorting.SortDef | null;
  onSortChange?: (sort: sorting.SortDef | null) => void;
}) {
  const handleClick = (col: Column<T>) => {
    if (!onSortChange) return;
    if (!col.sortable) return;
    if (!sort || sort.key !== col.key) {
      onSortChange({ key: col.key, direction: 'asc' });
    } else if (sort.direction === 'asc') {
      onSortChange({ key: col.key, direction: 'desc' });
    } else {
      onSortChange(null); // Remove sort
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
          return (
            <th
              key={col.key}
              style={col.width ? { width: col.width } : undefined}
              className={
                (col.className || '') +
                (col.sortable ? ' sortable' : '') +
                (isSorted ? ' sorted' : '')
              }
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
}
