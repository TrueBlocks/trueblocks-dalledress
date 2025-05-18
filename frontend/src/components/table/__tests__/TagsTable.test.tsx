import { MantineProvider } from '@mantine/core';
import { fireEvent, render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

import { TagsTable } from '../TagsTable';

describe('TagsTable', () => {
  const mockTags = ['Tag1', 'Tag2', 'Tag3'];
  const mockOnTagSelect = vi.fn();

  beforeEach(() => {
    mockOnTagSelect.mockClear();
  });

  it('renders tags list correctly', () => {
    render(
      <MantineProvider>
        <TagsTable
          tags={mockTags}
          onTagSelect={mockOnTagSelect}
          selectedTag={null}
          visible={true}
        />
      </MantineProvider>,
    );

    // Check that all tags are rendered
    expect(screen.getByText('Tag1')).toBeDefined();
    expect(screen.getByText('Tag2')).toBeDefined();
    expect(screen.getByText('Tag3')).toBeDefined();

    // Header should be visible
    expect(screen.getByText('Tags')).toBeDefined();
  });

  it('not visible when visible prop is false', () => {
    render(
      <MantineProvider>
        <TagsTable
          tags={mockTags}
          onTagSelect={mockOnTagSelect}
          selectedTag={null}
          visible={false}
        />{' '}
      </MantineProvider>,
    );

    // Container should be empty - nothing should be rendered
    expect(screen.queryByText('Tags')).not.toBeInTheDocument();
    expect(screen.queryByRole('list')).not.toBeInTheDocument();
  });

  it('shows selected tag correctly', () => {
    render(
      <MantineProvider>
        <TagsTable
          tags={mockTags}
          onTagSelect={mockOnTagSelect}
          selectedTag="Tag2"
          visible={true}
        />
      </MantineProvider>,
    );

    // Tag2 should have the selected class
    expect(screen.getByRole('button', { name: 'Tag2' })).toHaveClass(
      'selected',
    );
  });

  it('calls onTagSelect when tag is clicked, with focus=false', () => {
    render(
      <MantineProvider>
        <TagsTable
          tags={mockTags}
          onTagSelect={mockOnTagSelect}
          selectedTag={null}
          visible={true}
        />
      </MantineProvider>,
    );

    // Click the second tag
    fireEvent.click(screen.getByText('Tag2'));
    expect(mockOnTagSelect).toHaveBeenCalledWith('Tag2', false);
  });

  it('toggles tag selection when clicking already selected tag', () => {
    render(
      <MantineProvider>
        <TagsTable
          tags={mockTags}
          onTagSelect={mockOnTagSelect}
          selectedTag="Tag2"
          visible={true}
        />
      </MantineProvider>,
    );

    // Click the already selected tag should deselect it
    fireEvent.click(screen.getByText('Tag2'));
    expect(mockOnTagSelect).toHaveBeenCalledWith(null, false);
  });

  it('calls onTagSelect with focus=true when pressing Enter', () => {
    render(
      <MantineProvider>
        <TagsTable
          tags={mockTags}
          onTagSelect={mockOnTagSelect}
          selectedTag={null}
          visible={true}
        />
      </MantineProvider>,
    );

    // Get the tags table using the test id
    const tagsTable = screen.getByTestId('tags-table');

    // Press down to focus first tag
    fireEvent.keyDown(tagsTable, { key: 'ArrowDown' });
    // Press Enter to select and focus main table
    fireEvent.keyDown(tagsTable, { key: 'Enter' });
    expect(mockOnTagSelect).toHaveBeenCalledWith('Tag1', true);
  });

  it('shows no tags message when tags array is empty', () => {
    render(
      <MantineProvider>
        <TagsTable
          tags={[]}
          onTagSelect={mockOnTagSelect}
          selectedTag={null}
          visible={true}
        />
      </MantineProvider>,
    );

    expect(screen.getByText('No tags available')).toBeInTheDocument();
  });

  it('handles keyboard navigation with arrow keys', () => {
    render(
      <MantineProvider>
        <TagsTable
          tags={mockTags}
          onTagSelect={mockOnTagSelect}
          selectedTag={null}
          visible={true}
        />
      </MantineProvider>,
    );

    // Get the tags table using the test id
    const tagsTable = screen.getByTestId('tags-table');

    // Press down to move focus
    fireEvent.keyDown(tagsTable, { key: 'ArrowDown' });
    // Press down again to move to second tag
    fireEvent.keyDown(tagsTable, { key: 'ArrowDown' });
    // Press Enter on the second tag
    fireEvent.keyDown(tagsTable, { key: 'Enter' });
    expect(mockOnTagSelect).toHaveBeenCalledWith('Tag2', true);

    mockOnTagSelect.mockClear();

    // Press Escape to clear selection
    fireEvent.keyDown(tagsTable, { key: 'Escape' });
    expect(mockOnTagSelect).toHaveBeenCalledWith(null);
  });
});
