import { TableKey } from '@contexts';
import { render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

import { Column } from '../Column';
// Import after mocks are defined
import { Table } from '../Table';

// Mock the useTableKeys hook (must be before any imports that use it)
vi.mock('../useTableKeys', () => ({
  useTableKeys: () => ({
    handleKeyDown: vi.fn(),
    requestFocus: vi.fn(),
  }),
}));

// Mock the usePagination hook
vi.mock('../usePagination', () => ({
  usePagination: () => ({
    pagination: {
      currentPage: 0,
      pageSize: 10,
      totalItems: 3,
    },
    goToPage: vi.fn(),
    changePageSize: vi.fn(),
    setTotalItems: vi.fn(),
  }),
}));

// Mock the useTableContext hook (must be before any imports that use it)
vi.mock('@components', () => ({
  useTableContext: () => ({
    focusState: 'table',
    selectedRowIndex: -1,
    setSelectedRowIndex: vi.fn(),
    focusTable: vi.fn(),
    focusControls: vi.fn(),
    tableRef: { current: null },
  }),
}));

describe('Table', () => {
  type TestRow = {
    id: number;
    name: string;
    description: string;
    status?: string;
    deleted?: boolean;
  };

  const mockColumns: Column<TestRow>[] = [
    { key: 'id', header: 'ID', sortable: true },
    { key: 'name', header: 'Name', sortable: true },
    { key: 'description', header: 'Description', sortable: false },
    {
      key: 'status',
      header: 'Status',
      render: (row: TestRow) => `${row.deleted ? 'Deleted' : 'Active'}`,
      sortable: false,
    },
  ];

  const mockData = [
    { id: 1, name: 'Item 1', description: 'First item' },
    { id: 2, name: 'Item 2', description: 'Second item', deleted: true },
    { id: 3, name: 'Item 3', description: 'Third item' },
  ];

  const tableKey: TableKey = { viewName: 'test-view', tabName: 'test-tab' };

  const defaultProps = {
    columns: mockColumns,
    data: mockData,
    loading: false,
    error: null,
    tableKey,
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  // Group 1: Basic rendering tests
  describe('Rendering', () => {
    it('renders column headers', () => {
      render(<Table {...defaultProps} />);

      expect(screen.getByText('ID')).toBeInTheDocument();
      expect(screen.getByText('Name')).toBeInTheDocument();
      expect(screen.getByText('Description')).toBeInTheDocument();
      expect(screen.getByText('Status')).toBeInTheDocument();
    });

    it('renders data rows', () => {
      render(<Table {...defaultProps} />);

      expect(screen.getByText('Item 1')).toBeInTheDocument();
      expect(screen.getByText('Item 2')).toBeInTheDocument();
      expect(screen.getByText('Item 3')).toBeInTheDocument();
      expect(screen.getByText('First item')).toBeInTheDocument();
      expect(screen.getByText('Second item')).toBeInTheDocument();
      expect(screen.getByText('Third item')).toBeInTheDocument();
    });

    it('renders custom cell content using render function', () => {
      render(<Table {...defaultProps} />);

      const activeElements = screen.getAllByText('Active');
      const deletedElements = screen.getAllByText('Deleted');

      expect(activeElements.length).toBe(2);
      expect(deletedElements.length).toBe(1);
    });
  });

  // Group 2: State handling tests
  describe('State handling', () => {
    it('shows loading state', () => {
      render(<Table {...defaultProps} loading={true} />);

      expect(screen.getByText('Loading...')).toBeInTheDocument();
    });

    it('shows error state', () => {
      const errorMessage = 'Failed to load data';
      render(<Table {...defaultProps} error={errorMessage} />);

      expect(screen.getByText(`Error: ${errorMessage}`)).toBeInTheDocument();
    });

    it('shows no data message when data is empty', () => {
      render(<Table {...defaultProps} data={[]} />);

      expect(screen.getByText('No data found.')).toBeInTheDocument();
    });

    it('applies sorting when provided', () => {
      const sort = { key: 'id', direction: 'desc' as const };
      render(<Table {...defaultProps} sort={sort} />);

      const rows = screen.getAllByRole('row').slice(1);
      expect(rows[0]).toHaveTextContent('Item 3');
      expect(rows[1]).toHaveTextContent('Item 2');
      expect(rows[2]).toHaveTextContent('Item 1');
    });
  });

  // Group 4: TableKey integration tests
  describe('TableKey integration', () => {
    it('includes tableKey prop in rendered table', () => {
      expect(() => render(<Table {...defaultProps} />)).not.toThrow();
    });
  });

  // Group 5: Edge cases and special scenarios
  describe('Edge cases', () => {
    it('handles no sortable columns', () => {
      const nonSortableColumns = mockColumns.map((col) => ({
        ...col,
        sortable: false,
      }));
      render(<Table {...defaultProps} columns={nonSortableColumns} />);

      expect(screen.getByText('ID')).toBeInTheDocument();
    });

    it('handles large datasets', () => {
      const largeDataset = Array.from({ length: 100 }, (_, i) => ({
        id: i + 1,
        name: `Item ${i + 1}`,
        description: `Description ${i + 1}`,
      }));

      // Mock usePagination to return large dataset pagination
      vi.mock('../usePagination', () => ({
        usePagination: () => ({
          pagination: {
            currentPage: 0,
            pageSize: 10,
            totalItems: 100,
          },
          goToPage: vi.fn(),
          changePageSize: vi.fn(),
          setTotalItems: vi.fn(),
        }),
      }));

      render(
        <Table
          {...defaultProps}
          data={largeDataset.slice(0, 10)}
          tableKey={tableKey}
        />,
      );

      // Count rows using getByRole instead of container.querySelectorAll
      const rows = screen.getAllByRole('row');
      expect(rows.length).toBe(11);
    });

    it('handles column with custom render function', () => {
      const columnsWithCustomRenderer: Column<TestRow>[] = [
        ...mockColumns,
        {
          key: 'custom',
          header: 'Custom Column',
          render: () => 'Static content',
        },
      ];

      render(<Table {...defaultProps} columns={columnsWithCustomRenderer} />);

      expect(screen.getByText('Custom Column')).toBeInTheDocument();
      expect(screen.getAllByText('Static content').length).toBe(3);
    });
  });
});
