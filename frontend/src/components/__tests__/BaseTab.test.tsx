import React from 'react';

import { BaseTab, FormField } from '@components';
import { describe, expect, it, vi } from 'vitest';

import { render } from '../../__tests__/mocks';

vi.mock('../Table', () => ({
  Table: ({ data }: { data: any[] }) => (
    <div data-testid="mock-table">Table with {data.length} items</div>
  ),
}));

const mockColumns: FormField<{ id: string; name: string }>[] = [
  { key: 'id', header: 'ID' },
  { key: 'name', header: 'Name' },
];

const mockData = [
  { id: '1', name: 'Item 1' },
  { id: '2', name: 'Item 2' },
];

const mockTableKey = { viewName: 'test', tabName: 'test' };

describe('BaseTab', () => {
  it('shows loading state', () => {
    const { container } = render(
      <BaseTab
        data={[]}
        columns={mockColumns}
        loading={true}
        error={null}
        tableKey={mockTableKey}
      />,
    );

    expect(container.textContent).toContain('Loading...');
  });

  it('shows custom loading message', () => {
    const { container } = render(
      <BaseTab
        data={[]}
        columns={mockColumns}
        loading={true}
        error={null}
        tableKey={mockTableKey}
        loadingMessage="Custom loading..."
      />,
    );

    expect(container.textContent).toContain('Custom loading...');
  });

  it('shows error state', () => {
    const error = new Error('Test error');

    const { container } = render(
      <BaseTab
        data={[]}
        columns={mockColumns}
        loading={false}
        error={error}
        tableKey={mockTableKey}
      />,
    );

    expect(container.textContent).toContain('Error loading data');
    expect(container.textContent).toContain('Test error');
  });

  it('shows empty state', () => {
    const { container } = render(
      <BaseTab
        data={[]}
        columns={mockColumns}
        loading={false}
        error={null}
        tableKey={mockTableKey}
      />,
    );

    expect(container.textContent).toContain('No data available');
  });

  it('shows custom empty message', () => {
    const { container } = render(
      <BaseTab
        data={[]}
        columns={mockColumns}
        loading={false}
        error={null}
        tableKey={mockTableKey}
        emptyMessage="Custom empty message"
      />,
    );

    expect(container.textContent).toContain('Custom empty message');
  });

  it('renders table with data', () => {
    const { container } = render(
      <BaseTab
        data={mockData}
        columns={mockColumns}
        loading={false}
        error={null}
        tableKey={mockTableKey}
      />,
    );

    expect(container.textContent).toContain('Table with 2 items');
    expect(container.innerHTML).toContain('data-testid="mock-table"');
  });
});
