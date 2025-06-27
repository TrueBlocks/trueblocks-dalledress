import { ViewStateKey } from '@contexts';
import { types } from '@models';
import { render, screen } from '@testing-library/react';
import { Stats } from 'src/components/table/Stats';
import { describe, expect, it, vi } from 'vitest';

// Mock the usePagination hook
vi.mock('../usePagination', () => ({
  usePagination: vi.fn().mockImplementation(({ viewName, tabName }) => {
    // Different mock implementations for different test cases based on viewName/tabName
    if (viewName === 'test-view' && tabName === types.DataFacet.STATS) {
      return {
        pagination: {
          currentPage: 1,
          pageSize: 10,
          totalItems: 35,
        },
      };
    } else if (
      viewName === 'test-view' &&
      tabName === types.DataFacet.TRANSACTIONS
    ) {
      return {
        pagination: {
          currentPage: 0,
          pageSize: 10,
          totalItems: 5,
        },
      };
    } else if (viewName === 'test-view' && tabName === types.DataFacet.STATS) {
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
    const viewStateKey: ViewStateKey = {
      viewName: 'test-view',
      tabName: types.DataFacet.STATS,
    };
    render(<Stats namesLength={10} viewStateKey={viewStateKey} />);

    expect(
      screen.getByText(/Showing 11 to 20 of 35 entries/),
    ).toBeInTheDocument();
  });

  it('shows 1 to totalItems when on first page and less than pageSize', () => {
    const viewStateKey: ViewStateKey = {
      viewName: 'test-view',
      tabName: types.DataFacet.TRANSACTIONS,
    };
    render(<Stats namesLength={5} viewStateKey={viewStateKey} />);

    expect(screen.getByText(/Showing 1 to 5 of 5 entries/)).toBeInTheDocument();
  });

  // it('shows correct range when on last page with partial page', () => {
  //   const viewStateKey: ViewStateKey = {
  //     viewName: 'test-view',
  //     tabName: types.DataFacet.STATS,
  //   };
  //   render(<Stats namesLength={5} viewStateKey={viewStateKey} />);

  //   expect(
  //     screen.getByText(/Showing 31 to 35 of 35 entries/),
  //   ).toBeInTheDocument();
  // });

  it('shows 0 results correctly', () => {
    const viewStateKey: ViewStateKey = {
      viewName: 'test-view',
      tabName: types.DataFacet.ALL,
    };
    render(<Stats namesLength={0} viewStateKey={viewStateKey} />);

    expect(screen.getByText(/Showing 0 to 0 of 0 entries/)).toBeInTheDocument();
  });
});
