import { TableKey } from '@contexts';
import { render, screen } from '@testing-library/react';
import { Stats } from 'src/components/table/Stats';
import { describe, expect, it, vi } from 'vitest';

// Mock the usePagination hook
vi.mock('../usePagination', () => ({
  usePagination: vi.fn().mockImplementation(({ viewName, tabName }) => {
    // Different mock implementations for different test cases based on viewName/tabName
    if (viewName === 'test-view' && tabName === 'page1') {
      return {
        pagination: {
          currentPage: 1,
          pageSize: 10,
          totalItems: 35,
        },
      };
    } else if (viewName === 'test-view' && tabName === 'first-page') {
      return {
        pagination: {
          currentPage: 0,
          pageSize: 10,
          totalItems: 5,
        },
      };
    } else if (viewName === 'test-view' && tabName === 'last-page') {
      return {
        pagination: {
          currentPage: 3,
          pageSize: 10,
          totalItems: 35,
        },
      };
    } else {
      return {
        pagination: {
          currentPage: 0,
          pageSize: 10,
          totalItems: 0,
        },
      };
    }
  }),
}));

describe('Stats', () => {
  it('renders the correct range and total', () => {
    const tableKey: TableKey = { viewName: 'test-view', tabName: 'page1' };
    render(<Stats namesLength={10} tableKey={tableKey} />);

    expect(
      screen.getByText(/Showing 11 to 20 of 35 entries/),
    ).toBeInTheDocument();
  });

  it('shows 1 to totalItems when on first page and less than pageSize', () => {
    const tableKey: TableKey = { viewName: 'test-view', tabName: 'first-page' };
    render(<Stats namesLength={5} tableKey={tableKey} />);

    expect(screen.getByText(/Showing 1 to 5 of 5 entries/)).toBeInTheDocument();
  });

  it('shows correct range when on last page with partial page', () => {
    const tableKey: TableKey = { viewName: 'test-view', tabName: 'last-page' };
    render(<Stats namesLength={5} tableKey={tableKey} />);

    expect(
      screen.getByText(/Showing 31 to 35 of 35 entries/),
    ).toBeInTheDocument();
  });

  it('shows 0 results correctly', () => {
    const tableKey: TableKey = { viewName: 'test-view', tabName: 'empty' };
    render(<Stats namesLength={0} tableKey={tableKey} />);

    expect(screen.getByText(/Showing 0 to 0 of 0 entries/)).toBeInTheDocument();
  });
});
