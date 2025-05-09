import React from 'react';

import './Stats.css';

// StatsProps defines the props for the Stats component.
interface StatsProps {
  currentPage: number;
  pageSize: number;
  namesLength: number;
  totalItems: number;
}

// Stats displays a summary of the current entries being shown in the table.
export const Stats = ({
  currentPage,
  pageSize,
  namesLength,
  totalItems,
}: StatsProps) => (
  <div className="showing-entries">
    Showing {currentPage * pageSize + (namesLength > 0 ? 1 : 0)} to{' '}
    {Math.min((currentPage + 1) * pageSize, totalItems)} of {totalItems} entries
  </div>
);
