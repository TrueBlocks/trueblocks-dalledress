import React from 'react';

import { render, screen } from '@testing-library/react';
import { Stats } from 'src/components/table/Stats';

describe('Stats', () => {
  it('renders the correct range and total', () => {
    render(
      <Stats currentPage={1} pageSize={10} namesLength={10} totalItems={35} />,
    );
    expect(
      screen.getByText(/Showing 11 to 20 of 35 entries/),
    ).toBeInTheDocument();
  });

  it('shows 1 to totalItems when on first page and less than pageSize', () => {
    render(
      <Stats currentPage={0} pageSize={10} namesLength={5} totalItems={5} />,
    );
    expect(screen.getByText(/Showing 1 to 5 of 5 entries/)).toBeInTheDocument();
  });

  it('shows correct range when on last page with partial page', () => {
    render(
      <Stats currentPage={3} pageSize={10} namesLength={5} totalItems={35} />,
    );
    expect(
      screen.getByText(/Showing 31 to 35 of 35 entries/),
    ).toBeInTheDocument();
  });

  it('shows 0 results correctly', () => {
    render(
      <Stats currentPage={0} pageSize={10} namesLength={0} totalItems={0} />,
    );
    expect(screen.getByText(/Showing 0 to 0 of 0 entries/)).toBeInTheDocument();
  });
});
