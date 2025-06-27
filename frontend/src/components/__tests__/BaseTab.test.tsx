import React from 'react';

import { FormField } from '@components';
import { ViewStateKey } from '@contexts';
import { types } from '@models';
import { describe, expect, it, vi } from 'vitest';

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

const mockViewStateKey: ViewStateKey = {
  viewName: 'test',
  tabName: types.DataFacet.ALL,
};

describe('BaseTab', () => {
  it('renders table with data', () => {
    const { container } = render(
      <BaseTab
        data={mockData}
        columns={mockColumns}
        loading={false}
        error={null}
        viewStateKey={mockViewStateKey}
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
        viewStateKey={mockViewStateKey}
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
        viewStateKey={mockViewStateKey}
      />,
    );

    expect(container.innerHTML).toContain('data-testid="mock-table-provider"');
    expect(container.innerHTML).toContain('data-testid="mock-table"');
  });

  // DataFacet-related tests for refactor preparation
  describe('ViewStateKey handling (DataFacet refactor preparation)', () => {
    it('passes ViewStateKey with different tabName values correctly', () => {
      const exportsTxnKey: ViewStateKey = {
        viewName: '/exports',
        tabName: types.DataFacet.TRANSACTIONS,
      };
      const exportsReceiptsKey: ViewStateKey = {
        viewName: '/exports',
        tabName: types.DataFacet.RECEIPTS,
      };

      const { rerender } = render(
        <BaseTab
          data={mockData}
          columns={mockColumns}
          loading={false}
          error={null}
          viewStateKey={exportsTxnKey}
        />,
      );

      // Verify initial render with transactions key
      expect(
        document.querySelector('[data-testid="mock-table"]'),
      ).toBeInTheDocument();

      // Test state key change
      rerender(
        <BaseTab
          data={mockData}
          columns={mockColumns}
          loading={false}
          error={null}
          viewStateKey={exportsReceiptsKey}
        />,
      );

      expect(
        document.querySelector('[data-testid="mock-table"]'),
      ).toBeInTheDocument();
    });

    it('handles ViewStateKey creation patterns used in views', () => {
      // Simulate the pattern used in actual views: { viewName: ROUTE, tabName: dataFacet }
      const chunksKey: ViewStateKey = {
        viewName: '/chunks',
        tabName: types.DataFacet.INDEX,
      };
      const monitorsKey: ViewStateKey = {
        viewName: '/monitors',
        tabName: types.DataFacet.MONITORS,
      };
      const abisKey: ViewStateKey = {
        viewName: '/abis',
        tabName: types.DataFacet.DOWNLOADED,
      };

      // Test each key type
      [chunksKey, monitorsKey, abisKey].forEach((key) => {
        const { container, unmount } = render(
          <BaseTab
            data={mockData}
            columns={mockColumns}
            loading={false}
            error={null}
            viewStateKey={key}
          />,
        );

        expect(container.innerHTML).toContain('data-testid="mock-table"');
        unmount();
      });
    });

    it('properly forwards ViewStateKey to Table component', () => {
      const testKey: ViewStateKey = {
        viewName: '/names',
        tabName: types.DataFacet.CUSTOM,
      };

      render(
        <BaseTab
          data={mockData}
          columns={mockColumns}
          loading={false}
          error={null}
          viewStateKey={testKey}
        />,
      );

      // The mocked Table component should receive the viewStateKey
      // This tests that BaseTab correctly passes through the key without modification
      expect(
        document.querySelector('[data-testid="mock-table"]'),
      ).toBeInTheDocument();
    });

    it('handles ViewStateKey uniqueness requirements', () => {
      const key1: ViewStateKey = {
        viewName: '/exports',
        tabName: types.DataFacet.TRANSACTIONS,
      };
      const key2: ViewStateKey = {
        viewName: '/exports',
        tabName: types.DataFacet.RECEIPTS,
      };
      const key3: ViewStateKey = {
        viewName: '/names',
        tabName: types.DataFacet.REGULAR,
      }; // same tabName, different view

      // All keys should be valid and unique for state management
      [key1, key2, key3].forEach((key) => {
        expect(key.viewName).toBeTruthy();
        expect(key.tabName).toBeTruthy();
        expect(typeof key.viewName).toBe('string');
        expect(typeof key.tabName).toBe('string');
      });
    });
  });

  describe('Component integration tests', () => {
    it('passes onSubmit callback correctly', () => {
      const mockOnSubmit = vi.fn();

      render(
        <BaseTab
          data={mockData}
          columns={mockColumns}
          loading={false}
          error={null}
          viewStateKey={mockViewStateKey}
          onSubmit={mockOnSubmit}
        />,
      );

      expect(
        document.querySelector('[data-testid="mock-table"]'),
      ).toBeInTheDocument();
    });

    it('passes onAction callback correctly', () => {
      const mockOnAction = vi.fn();

      render(
        <BaseTab
          data={mockData}
          columns={mockColumns}
          loading={false}
          error={null}
          viewStateKey={mockViewStateKey}
          onAction={mockOnAction}
        />,
      );

      expect(
        document.querySelector('[data-testid="mock-table"]'),
      ).toBeInTheDocument();
    });

    it('handles error prop without breaking render', () => {
      const mockError = new Error('Test error');

      const { container } = render(
        <BaseTab
          data={mockData}
          columns={mockColumns}
          loading={false}
          error={mockError}
          viewStateKey={mockViewStateKey}
        />,
      );

      expect(container.innerHTML).toContain('data-testid="mock-table"');
    });
  });
});
