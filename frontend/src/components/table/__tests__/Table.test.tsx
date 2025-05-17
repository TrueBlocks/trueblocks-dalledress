import { TableKey } from '@contexts';
import { render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { MantineProvider } from '@mantine/core';

import { Column, Table, TableProps, TableProvider } from '@components';

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
  }),
}));

// Mock the useTableContext hook (must be before any imports that use it)
vi.mock('@components', async (importOriginal) => {
  const actual = await importOriginal() as any;
  return {
    ...actual, // Preserve other exports from @components if any
    useTableContext: () => ({
      focusState: 'table',
      selectedRowIndex: -1,
      setSelectedRowIndex: vi.fn(),
      focusTable: vi.fn(),
      focusControls: vi.fn(),
      tableRef: { current: null },
    }),
    // Add the missing useTableKeys mock here
    useTableKeys: () => ({
      handleKeyDown: vi.fn(),
      requestFocus: vi.fn(),
    }),
  };
});

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

  const defaultProps: TableProps<TestRow> = { // Explicitly type defaultProps
    columns: mockColumns,
    data: mockData,
    loading: false,
    error: null,
    tableKey,
  };

  // Helper function to setup the test environment
  const setupTest = (props: Partial<TableProps<TestRow>> = {}) => {
    const testProps = { ...defaultProps, ...props };
    return render(
      <MantineProvider>
        <Table {...testProps} />
      </MantineProvider>
    );
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  // Group 1: Basic rendering tests
  describe('Rendering', () => {
    it('renders column headers', () => {
      setupTest();
      expect(screen.getByText('ID')).toBeInTheDocument();
      expect(screen.getByText('Name')).toBeInTheDocument();
      expect(screen.getByText('Description')).toBeInTheDocument();
      expect(screen.getByText('Status')).toBeInTheDocument();
    });

    it('renders data rows', () => {
      setupTest();
      expect(screen.getByText('Item 1')).toBeInTheDocument();
      expect(screen.getByText('Item 2')).toBeInTheDocument();
      expect(screen.getByText('Item 3')).toBeInTheDocument();
      expect(screen.getByText('First item')).toBeInTheDocument();
      expect(screen.getByText('Second item')).toBeInTheDocument();
      expect(screen.getByText('Third item')).toBeInTheDocument();
    });

    it('renders custom cell content using render function', () => {
      setupTest();
      const activeElements = screen.getAllByText('Active');
      const deletedElements = screen.getAllByText('Deleted');

      expect(activeElements.length).toBe(2);
      expect(deletedElements.length).toBe(1);
    });
  });

  // Group 2: State handling tests
  describe('State handling', () => {
    it('shows loading state', () => {
      setupTest({ loading: true });
      expect(screen.getByText('Loading...')).toBeInTheDocument();
    });

    it('shows error state', () => {
      const errorMessage = 'Failed to load data';
      setupTest({ error: errorMessage });
      expect(screen.getByText(`Error: ${errorMessage}`)).toBeInTheDocument();
    });

    it('shows no data message when data is empty', () => {
      setupTest({ data: [] });
      expect(screen.getByText('No data found.')).toBeInTheDocument();
    });

    it('applies sorting when provided', () => {
      const sort = { key: 'id', direction: 'desc' as const };
      setupTest({ sort });
      const rows = screen.getAllByRole('row').slice(1); // Exclude header row

      expect(rows[0]).toHaveTextContent('3'); // ID of Item 3
      expect(rows[1]).toHaveTextContent('2'); // ID of Item 2
      expect(rows[2]).toHaveTextContent('1'); // ID of Item 1
    });
  });

  // Group 4: TableKey integration tests
  describe('TableKey integration', () => {
    it('includes tableKey prop in rendered table', () => {
      expect(() => setupTest()).not.toThrow();
    });
  });

  // // Group 5: Edge cases and special scenarios
  describe('Edge cases', () => {
    it('handles no sortable columns', () => {
      const nonSortableColumns = mockColumns.map((col) => ({
        ...col,
        sortable: false,
      }));
      render(<MantineProvider><Table {...defaultProps} columns={nonSortableColumns} /></MantineProvider>);
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
        }),
      }));

      render(
        <MantineProvider>
          <Table
            {...defaultProps}
            data={largeDataset.slice(0, 10)}
            tableKey={tableKey}
          />
        </MantineProvider>
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

      render(<MantineProvider><Table {...defaultProps} columns={columnsWithCustomRenderer} /></MantineProvider>);

      expect(screen.getByText('Custom Column')).toBeInTheDocument();
      expect(screen.getAllByText('Static content').length).toBe(3);
    });
  });
});
