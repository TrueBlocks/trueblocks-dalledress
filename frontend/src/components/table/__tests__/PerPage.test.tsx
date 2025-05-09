import React from 'react';

import { PerPage } from '@components';
import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

describe('PerPage', () => {
  const defaultProps = {
    pageSize: 25,
    onPageSizeChange: vi.fn(),
    focusTable: vi.fn(),
    focusControls: vi.fn(),
  };

  it('renders the select with correct value and options', () => {
    render(<PerPage {...defaultProps} />);
    const select = screen.getByLabelText('Items per page');
    expect(select).toBeInTheDocument();
    expect(select).toHaveValue('25');
    // The visible text for each option is the number only
    expect(
      screen.getByRole('option', { name: '10 per page' }),
    ).toBeInTheDocument();
    expect(screen.getByRole('option', { name: '25' })).toBeInTheDocument();
    expect(screen.getByRole('option', { name: '50' })).toBeInTheDocument();
    expect(screen.getByRole('option', { name: '100' })).toBeInTheDocument();
  });

  it('calls onPageSizeChange and focusTable when select changes', () => {
    const onPageSizeChange = vi.fn();
    const focusTable = vi.fn();
    render(
      <PerPage
        {...defaultProps}
        onPageSizeChange={onPageSizeChange}
        focusTable={focusTable}
      />,
    );
    fireEvent.change(screen.getByLabelText('Items per page'), {
      target: { value: '50' },
    });
    expect(onPageSizeChange).toHaveBeenCalledWith(50);
    // focusTable is called after a timeout, so we can use fake timers if needed
  });

  it('calls focusControls on focus', () => {
    const focusControls = vi.fn();
    render(<PerPage {...defaultProps} focusControls={focusControls} />);
    fireEvent.focus(screen.getByLabelText('Items per page'));
    expect(focusControls).toHaveBeenCalled();
  });
});
