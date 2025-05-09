import React from 'react';

import { Pagination } from '@components';
import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

describe('Pagination', () => {
  const defaultProps = {
    totalPages: 10,
    currentPage: 4,
    handlePageChange: vi.fn(),
    focusControls: vi.fn(),
  };

  it('renders navigation and correct page buttons', () => {
    render(<Pagination {...defaultProps} />);
    expect(screen.getByTitle('First Page')).toBeInTheDocument();
    expect(screen.getByTitle('Previous Page')).toBeInTheDocument();
    expect(screen.getByTitle('Next Page')).toBeInTheDocument();
    expect(screen.getByTitle('Last Page')).toBeInTheDocument();
    // Should show 5 page buttons: 3, 4, 5, 6, 7 (for currentPage=4)
    [3, 4, 5, 6, 7].forEach((n) => {
      expect(screen.getByText(n.toString())).toBeInTheDocument();
    });
  });

  it('disables prev/first on first page', () => {
    render(<Pagination {...defaultProps} currentPage={0} />);
    expect(screen.getByTitle('First Page')).toBeDisabled();
    expect(screen.getByTitle('Previous Page')).toBeDisabled();
  });

  it('disables next/last on last page', () => {
    render(<Pagination {...defaultProps} currentPage={9} />);
    expect(screen.getByTitle('Next Page')).toBeDisabled();
    expect(screen.getByTitle('Last Page')).toBeDisabled();
  });

  it('calls handlePageChange when a page button is clicked', () => {
    const handlePageChange = vi.fn();
    render(
      <Pagination {...defaultProps} handlePageChange={handlePageChange} />,
    );
    fireEvent.click(screen.getByText('6'));
    expect(handlePageChange).toHaveBeenCalledWith(5);
  });

  it('returns five disabled page buttons if only one page', () => {
    render(<Pagination {...defaultProps} totalPages={1} />);
    const buttons = screen.getAllByRole('button');
    expect(buttons).toHaveLength(4 + 5); // 5 nav + 5 page buttons
    // Navigation buttons
    expect(buttons[0]).toBeDisabled(); // First
    expect(buttons[1]).toBeDisabled(); // Prev
    expect(buttons[7]).toBeDisabled(); // Next
    expect(buttons[8]).toBeDisabled(); // Last
    // Page buttons
    for (let i = 2; i <= 6; i++) {
      expect(buttons[i]).toBeDisabled();
    }
    expect(buttons[2]).toHaveTextContent('1'); // The third page button is '1'
  });
});
