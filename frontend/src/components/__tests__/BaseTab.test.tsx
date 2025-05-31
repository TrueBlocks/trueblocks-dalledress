import React from 'react';

import { FormField } from '@components';
import { describe, expect, it } from 'vitest';

import { render } from '../../__tests__/mocks';
import { BaseTab } from '../BaseTab';

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

  it('renders table structure even when loading', () => {
    const { container } = render(
      <BaseTab
        data={[]}
        columns={mockColumns}
        loading={true}
        error={null}
        tableKey={mockTableKey}
      />,
    );

    expect(container.innerHTML).toContain('data-testid="mock-table-provider"');
    expect(container.innerHTML).toContain('data-testid="mock-table"');
  });

  it('renders table structure even when empty', () => {
    const { container } = render(
      <BaseTab
        data={[]}
        columns={mockColumns}
        loading={false}
        error={null}
        tableKey={mockTableKey}
      />,
    );

    expect(container.innerHTML).toContain('data-testid="mock-table-provider"');
    expect(container.innerHTML).toContain('data-testid="mock-table"');
  });
});
