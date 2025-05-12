import React from 'react';

import { TableKey } from '@contexts';
import { usePagination } from 'src/components/table/usePagination';

import './Stats.css';

// StatsProps defines the props for the Stats component.
interface StatsProps {
  namesLength: number;
  tableKey: TableKey;
}

// Stats displays a summary of the current entries being shown in the table.
export const Stats = ({ namesLength, tableKey }: StatsProps) => {
  const { pagination } = usePagination(tableKey);
  const { currentPage, pageSize, totalItems } = pagination;
  return (
    <div className="showing-entries">
      Showing {currentPage * pageSize + (namesLength > 0 ? 1 : 0)} to{' '}
      {Math.min((currentPage + 1) * pageSize, totalItems)} of {totalItems}{' '}
      entries
    </div>
  );
};
