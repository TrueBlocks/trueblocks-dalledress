import { Body, Column } from '@components';
import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

describe('Body', () => {
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
    {
      name: 'Bob',
      address: '0x2',
      tags: 'tag2',
      source: 'src2',
      deleted: true,
    },
    {
      name: 'Carol',
      address: '0x3',
      tags: 'tag3',
      source: 'src3',
      isCustom: true,
    },
    {
      name: 'Dave',
      address: '0x4',
      tags: 'tag4',
      source: 'src4',
      isPrefund: true,
    },
  ];
  const handleRowClick = vi.fn();

  it('renders all rows and columns', () => {
    render(
      <table>
        <tbody>
          <Body
            columns={columns}
            data={names}
            selectedRowIndex={1}
            handleRowClick={handleRowClick}
            expandedRowIndex={null}
            setExpandedRowIndex={() => {}}
          />
        </tbody>
      </table>,
    );
    expect(screen.getByText('Alice')).toBeInTheDocument();
    expect(screen.getByText('Bob')).toBeInTheDocument();
    expect(screen.getByText('Carol')).toBeInTheDocument();
    expect(screen.getByText('Dave')).toBeInTheDocument();
    expect(screen.getAllByRole('row')).toHaveLength(4);
  });

  it('calls handleRowClick when a row is clicked', () => {
    render(
      <table>
        <tbody>
          <Body
            columns={columns}
            data={names}
            selectedRowIndex={0}
            handleRowClick={handleRowClick}
            expandedRowIndex={null}
            setExpandedRowIndex={() => {}}
          />
        </tbody>
      </table>,
    );
    fireEvent.click(screen.getByText('Bob'));
    expect(handleRowClick).toHaveBeenCalledWith(1);
  });

  it('applies selected class to the selected row', () => {
    render(
      <table>
        <tbody>
          <Body
            columns={columns}
            data={names}
            selectedRowIndex={2}
            handleRowClick={handleRowClick}
            expandedRowIndex={null}
            setExpandedRowIndex={() => {}}
          />
        </tbody>
      </table>,
    );
    const rows = screen.getAllByRole('row');
    expect(rows[2]?.className).toContain('selected');
  });

  it('shows status text for deleted, custom, and prefund', () => {
    render(
      <table>
        <tbody>
          <Body
            columns={columns}
            data={names}
            selectedRowIndex={0}
            handleRowClick={handleRowClick}
            expandedRowIndex={null}
            setExpandedRowIndex={() => {}}
          />
        </tbody>
      </table>,
    );
    expect(screen.getByText(/Deleted/)).toBeInTheDocument();
    expect(screen.getByText(/Custom/)).toBeInTheDocument();
    expect(screen.getByText(/Prefund/)).toBeInTheDocument();
  });
});
