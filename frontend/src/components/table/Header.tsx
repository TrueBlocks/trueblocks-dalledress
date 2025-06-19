import React from 'react';

import { FormField } from '@components';
import { TableKey, useSorting } from '@contexts';

import { SortSpecManager } from '../../utils/sortSpec';
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
    if (!col.sortable || typeof col.key !== 'string') return;

    // Use SortSpecManager to handle multi-field sorting
    const newSort = SortSpecManager.handleFieldClick(
      sort || { fields: [], orders: [] },
      col.key as string,
    );
    setSorting(newSort.fields.length > 0 ? newSort : null);
  };

  return (
    <thead>
      <tr>
        {columns.map((col) => {
          // Get sort info from SortSpecManager for this column
          const sortInfo = SortSpecManager.getSortInfo(
            sort || { fields: [], orders: [] },
            col.key as string,
          );

          const isSorted = sortInfo.active;

          // Determine aria-sort value
          let ariaSort: 'ascending' | 'descending' | undefined;
          if (isSorted) {
            ariaSort =
              sortInfo.direction === 'asc' ? 'ascending' : 'descending';
          }

          // Determine sort indicator with priority
          let sortIndicator = '';
          if (col.sortable) {
            if (isSorted) {
              // Show direction arrow with priority subscript for active sorts
              const arrow = sortInfo.direction === 'asc' ? '↑' : '↓';
              const priority =
                sortInfo.priority > 1
                  ? String.fromCharCode(8320 + sortInfo.priority)
                  : '';
              sortIndicator = ` ${arrow}${priority}`;
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
              style={{
                ...(typeof col.width === 'string' &&
                col.width.startsWith('col-')
                  ? undefined
                  : col.width
                    ? { width: col.width }
                    : undefined),
                textAlign: 'center',
              }}
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
