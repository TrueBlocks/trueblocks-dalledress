import React from 'react';

import { Column, Table } from '@components';
import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

type NameRow = {
  name: string;
  address: string;
  tags: string;
  source: string;
  deleted?: boolean;
  isCustom?: boolean;
  isPrefund?: boolean;
};

const columns: Column<NameRow>[] = [
  { key: 'name', header: 'Name' },
  { key: 'address', header: 'Address' },
  { key: 'tags', header: 'Tags' },
  { key: 'source', header: 'Source' },
  {
    key: 'status',
    header: 'Status',
    render: (row) =>
      `${row.deleted ? 'Deleted ' : ''}${row.isCustom ? 'Custom ' : ''}${row.isPrefund ? 'Prefund' : ''}`,
  },
];

const names: NameRow[] = [
  { name: 'Alice', address: '0x1', tags: 'tag1', source: 'src1' },
  { name: 'Bob', address: '0x2', tags: 'tag2', source: 'src2' },
  { name: 'Carol', address: '0x3', tags: 'tag3', source: 'src3' },
];

describe('Table', () => {
  const defaultProps = {
    columns,
    data: names,
    loading: false,
    error: null,
    pagination: { currentPage: 0, pageSize: 2, totalItems: 3 },
    onPageChange: vi.fn(),
    onPageSizeChange: vi.fn(),
  };

  it('renders table headers and rows', () => {
    render(<Table {...defaultProps} />);
    expect(screen.getByText('Name')).toBeInTheDocument();
    expect(screen.getByText('Alice')).toBeInTheDocument();
    expect(screen.getByText('Bob')).toBeInTheDocument();
  });

  it('shows loading message', () => {
    render(<Table {...defaultProps} loading={true} />);
    expect(screen.getByText(/Loading/)).toBeInTheDocument();
  });

  it('shows error message', () => {
    render(<Table {...defaultProps} error="Something went wrong" />);
    expect(screen.getByText(/Something went wrong/)).toBeInTheDocument();
  });

  it('shows no data found', () => {
    render(<Table {...defaultProps} data={[]} />);
    expect(screen.getByText(/No data found/)).toBeInTheDocument();
  });

  it('calls onPageChange when pagination is used', () => {
    const onPageChange = vi.fn();
    render(<Table {...defaultProps} onPageChange={onPageChange} />);
    fireEvent.click(screen.getByTitle('Next Page'));
    expect(onPageChange).toHaveBeenCalled();
  });

  it('calls onPageSizeChange when per page is changed', () => {
    const onPageSizeChange = vi.fn();
    render(<Table {...defaultProps} onPageSizeChange={onPageSizeChange} />);
    fireEvent.change(screen.getByLabelText('Items per page'), {
      target: { value: '10' },
    });
    expect(onPageSizeChange).toHaveBeenCalledWith(10);
  });
});
