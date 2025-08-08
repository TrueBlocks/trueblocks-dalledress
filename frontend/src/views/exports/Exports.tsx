// Copyright 2016, 2026 The Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */
// === SECTION 1: Imports & Dependencies ===
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { GetExportsPage, Reload } from '@app';
import { BaseTab, usePagination } from '@components';
import { Action, ExportFormatModal } from '@components';
import { useFiltering, useSorting } from '@contexts';
import {
  ActionType,
  DataFacetConfig,
  toPageDataProp,
  useActions,
  useActiveFacet,
  useColumns,
  useEvent,
  usePayload,
} from '@hooks';
import { TabView } from '@layout';
import { Group } from '@mantine/core';
import { useHotkeys } from '@mantine/hooks';
import { exports } from '@models';
import { msgs, project, types } from '@models';
import { Debugger, getDisplayAddress, useErrorHandler } from '@utils';

import { getColumns } from './columns';
import { exportsFacets } from './facets';

export const ROUTE = 'exports' as const;
export const Exports = () => {
  // === SECTION 2: Hook Initialization ===
  const renderCnt = useRef(0);
  const createPayload = usePayload();
  const activeFacetHook = useActiveFacet({
    facets: exportsFacets,
    viewRoute: ROUTE,
  });
  const { availableFacets, getCurrentDataFacet } = activeFacetHook;

  const [pageData, setPageData] = useState<exports.ExportsPage | null>(null);
  const viewStateKey = useMemo(
    (): project.ViewStateKey => ({
      viewName: ROUTE,
      facetName: getCurrentDataFacet(),
    }),
    [getCurrentDataFacet],
  );

  const { error, handleError, clearError } = useErrorHandler();
  const { pagination, setTotalItems } = usePagination(viewStateKey);
  const { sort } = useSorting(viewStateKey);
  const { filter } = useFiltering(viewStateKey);

  // Custom bank statement-style detail panel for Statement rows
  const renderStatementDetailPanel = useCallback(
    (rowData: Record<string, unknown> | null) => {
      if (!rowData) return null;
      const statement = rowData as unknown as types.Statement;

      // Helper function to format values as ether
      const formatValue = (value: unknown) => {
        if (!value || value === '0') return '-';

        // Convert wei to ether using the statement's decimals
        const valueStr = typeof value === 'string' ? value : String(value);
        const decimals = statement.decimals || 18; // Default to 18 for ETH

        try {
          // Simple conversion: divide by 10^decimals
          const divisor = Math.pow(10, decimals);
          const etherValue = parseFloat(valueStr) / divisor;

          // Format to 3 decimal places
          return etherValue.toFixed(3);
        } catch {
          return valueStr; // Fallback to original value if conversion fails
        }
      };

      // Check if reconciled (no correcting reasons means it's reconciled)
      const isReconciled =
        !statement.correctingReasons || statement.correctingReasons === '';

      return (
        <div
          style={{
            fontFamily: 'monospace',
            fontSize: '14px',
            lineHeight: '1.6',
          }}
        >
          <div
            style={{
              textAlign: 'center',
              marginBottom: '20px',
              borderBottom: '2px solid #333',
              paddingBottom: '10px',
            }}
          >
            <h3>Jay&apos;s Bank Statement</h3>
            <div style={{ fontSize: '12px', color: '#666' }}>
              Block: {statement.blockNumber} | Tx: {statement.transactionIndex}|
              Date: {new Date(statement.timestamp * 1000).toLocaleDateString()}
            </div>
          </div>

          <div style={{ marginBottom: '15px' }}>
            <strong>Transaction Identification:</strong>
            <br />
            <div style={{ marginLeft: '10px' }}>
              Account: {getDisplayAddress(statement.accountedFor)}
              <br />
              Asset: {statement.symbol} ({getDisplayAddress(statement.asset)})
              <br />
              From: {getDisplayAddress(statement.sender)}
              <br />
              To: {getDisplayAddress(statement.recipient)}
            </div>
          </div>

          <table
            style={{
              width: '100%',
              borderCollapse: 'collapse',
              marginBottom: '15px',
            }}
          >
            <thead>
              <tr style={{ borderBottom: '1px solid #333' }}>
                <th style={{ textAlign: 'left', padding: '5px' }}>Type</th>
                <th style={{ textAlign: 'right', padding: '5px' }}>In</th>
                <th style={{ textAlign: 'right', padding: '5px' }}>Out</th>
                <th style={{ textAlign: 'right', padding: '5px' }}>Gas</th>
                <th style={{ textAlign: 'right', padding: '5px' }}>Balance</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td style={{ padding: '5px' }}>Beginning Balance</td>
                <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
                <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
                <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
                <td style={{ textAlign: 'right', padding: '5px' }}>
                  {formatValue(statement.begBal)}
                </td>
              </tr>
              <tr>
                <td style={{ padding: '5px' }}>Regular Amount</td>
                <td style={{ textAlign: 'right', padding: '5px' }}>
                  {formatValue(statement.amountIn)}
                </td>
                <td style={{ textAlign: 'right', padding: '5px' }}>
                  {formatValue(statement.amountOut)}
                </td>
                <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
                <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
              </tr>
              <tr>
                <td style={{ padding: '5px' }}>Internal Transfers</td>
                <td style={{ textAlign: 'right', padding: '5px' }}>
                  {formatValue(statement.internalIn)}
                </td>
                <td style={{ textAlign: 'right', padding: '5px' }}>
                  {formatValue(statement.internalOut)}
                </td>
                <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
                <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
              </tr>
              <tr>
                <td style={{ padding: '5px' }}>Self Destruct</td>
                <td style={{ textAlign: 'right', padding: '5px' }}>
                  {formatValue(statement.selfDestructIn)}
                </td>
                <td style={{ textAlign: 'right', padding: '5px' }}>
                  {formatValue(statement.selfDestructOut)}
                </td>
                <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
                <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
              </tr>
              <tr>
                <td style={{ padding: '5px' }}>Gas Costs</td>
                <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
                <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
                <td style={{ textAlign: 'right', padding: '5px' }}>
                  {formatValue(statement.gasOut)}
                </td>
                <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
              </tr>
              <tr>
                <td style={{ padding: '5px' }}>Miner Rewards</td>
                <td style={{ textAlign: 'right', padding: '5px' }}>
                  {formatValue(statement.minerBaseRewardIn)}
                </td>
                <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
                <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
                <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
              </tr>
              <tr style={{ borderTop: '1px solid #333', fontWeight: 'bold' }}>
                <td style={{ padding: '5px' }}>Ending Balance</td>
                <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
                <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
                <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
                <td style={{ textAlign: 'right', padding: '5px' }}>
                  {formatValue(statement.endBal)}
                </td>
              </tr>
            </tbody>
          </table>

          <div
            style={{
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              marginTop: '10px',
            }}
          >
            <span style={{ marginRight: '10px' }}>Reconciled:</span>
            {isReconciled ? (
              <span style={{ color: 'green', fontSize: '18px' }}>✓</span>
            ) : (
              <span style={{ color: 'red', fontSize: '18px' }}>✗</span>
            )}
            {!isReconciled && statement.correctingReasons && (
              <span
                style={{ marginLeft: '10px', fontSize: '12px', color: '#666' }}
              >
                ({statement.correctingReasons})
              </span>
            )}
          </div>
        </div>
      );
    },
    [],
  );

  // === SECTION 3: Data Fetching ===
  const fetchData = useCallback(async () => {
    clearError();
    try {
      const result = await GetExportsPage(
        createPayload(getCurrentDataFacet()),
        pagination.currentPage * pagination.pageSize,
        pagination.pageSize,
        sort,
        filter,
      );
      setPageData(result);
      setTotalItems(result.totalItems || 0);
    } catch (err: unknown) {
      handleError(err, `Failed to fetch ${getCurrentDataFacet()}`);
    }
  }, [
    clearError,
    createPayload,
    getCurrentDataFacet,
    pagination.currentPage,
    pagination.pageSize,
    sort,
    filter,
    setTotalItems,
    handleError,
  ]);

  const currentData = useMemo(() => {
    if (!pageData) return [];
    const facet = getCurrentDataFacet();
    switch (facet) {
      case types.DataFacet.STATEMENTS:
        return pageData.statements || [];
      case types.DataFacet.BALANCES:
        return pageData.balances || [];
      case types.DataFacet.TRANSFERS:
        return pageData.transfers || [];
      case types.DataFacet.TRANSACTIONS:
        return pageData.transactions || [];
      case types.DataFacet.WITHDRAWALS:
        return pageData.withdrawals || [];
      case types.DataFacet.ASSETS:
        return pageData.assets || [];
      case types.DataFacet.LOGS:
        return pageData.logs || [];
      case types.DataFacet.TRACES:
        return pageData.traces || [];
      case types.DataFacet.RECEIPTS:
        return pageData.receipts || [];
      default:
        return [];
    }
  }, [pageData, getCurrentDataFacet]);

  // === SECTION 4: Event Handling ===
  useEvent(
    msgs.EventType.DATA_LOADED,
    (_message: string, payload?: Record<string, unknown>) => {
      if (payload?.collection === 'exports') {
        const eventDataFacet = payload.dataFacet;
        if (eventDataFacet === getCurrentDataFacet()) {
          fetchData();
        }
      }
    },
  );

  // Listen for active address/chain/period changes to refresh data
  useEvent(msgs.EventType.ADDRESS_CHANGED, fetchData);
  useEvent(msgs.EventType.CHAIN_CHANGED, fetchData);
  useEvent(msgs.EventType.PERIOD_CHANGED, fetchData);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const handleReload = useCallback(async () => {
    clearError();
    try {
      Reload(createPayload(getCurrentDataFacet())).then(() => {
        // The data will reload when the DataLoaded event is fired.
      });
    } catch (err: unknown) {
      handleError(err, `Failed to reload ${getCurrentDataFacet()}`);
    }
  }, [clearError, getCurrentDataFacet, createPayload, handleError]);

  useHotkeys([['mod+r', handleReload]]);

  // === SECTION 5: Actions Configuration ===
  const enabledActions = useMemo(() => {
    // Auto-enable Export for all DataTable views that are not forms
    return ['export'] as ActionType[];
  }, []);

  const { handlers, config, exportFormatModal } = useActions({
    collection: 'exports',
    viewStateKey,
    pagination,
    goToPage: () => {}, // Exports typically don't have pagination navigation
    sort,
    filter,
    enabledActions,
    pageData,
    setPageData,
    setTotalItems,
    crudFunc: () => Promise.resolve(), // Exports don't have CRUD operations
    pageFunc: GetExportsPage,
    pageClass: exports.ExportsPage,
    updateItem: undefined,
    createPayload,
    getCurrentDataFacet,
  });

  const headerActions = useMemo(() => {
    if (!config.headerActions?.length) return null;
    return (
      <Group gap="xs" style={{ flexShrink: 0 }}>
        {config.headerActions.map((action) => {
          const handlerKey =
            `handle${action.type.charAt(0).toUpperCase() + action.type.slice(1)}` as keyof typeof handlers;
          const handler = handlers[handlerKey] as () => void;
          return (
            <Action
              key={action.type}
              icon={
                action.icon as keyof ReturnType<
                  typeof import('@hooks').useIconSets
                >
              }
              onClick={handler}
              title={
                action.requiresWallet && !config.isWalletConnected
                  ? `${action.title} (requires wallet connection)`
                  : action.title
              }
              hotkey={action.type === 'export' ? 'mod+x' : undefined}
              size="sm"
            />
          );
        })}
      </Group>
    );
  }, [config.headerActions, config.isWalletConnected, handlers]);

  // === SECTION 6: UI Configuration ===
  const currentColumns = useColumns(
    getColumns(getCurrentDataFacet()),
    {
      showActions: false,
      actions: [],
      getCanRemove: useCallback((_row: unknown) => false, []),
    },
    {},
    toPageDataProp(pageData),
    { rowActions: [] },
  );

  const perTabContent = useMemo(() => {
    const isStatementsView =
      getCurrentDataFacet() === types.DataFacet.STATEMENTS;
    return (
      <BaseTab<Record<string, unknown>>
        data={currentData as unknown as Record<string, unknown>[]}
        columns={currentColumns}
        loading={!!pageData?.isFetching}
        error={error}
        viewStateKey={viewStateKey}
        headerActions={headerActions}
        detailPanel={isStatementsView ? renderStatementDetailPanel : undefined}
      />
    );
  }, [
    currentData,
    currentColumns,
    pageData?.isFetching,
    error,
    viewStateKey,
    headerActions,
    getCurrentDataFacet,
    renderStatementDetailPanel,
  ]);

  const tabs = useMemo(
    () =>
      availableFacets.map((facetConfig: DataFacetConfig) => ({
        key: facetConfig.id,
        label: facetConfig.label,
        value: facetConfig.id,
        content: perTabContent,
        dividerBefore: facetConfig.dividerBefore,
      })),
    [availableFacets, perTabContent],
  );

  // === SECTION 7: Render ===
  return (
    <div className="mainView">
      <TabView tabs={tabs} route={ROUTE} />
      {error && (
        <div>
          <h3>{`Error fetching ${getCurrentDataFacet()}`}</h3>
          <p>{error.message}</p>
        </div>
      )}
      <Debugger
        rowActions={[]}
        headerActions={[]}
        count={++renderCnt.current}
      />
      <ExportFormatModal
        opened={exportFormatModal.opened}
        onClose={exportFormatModal.onClose}
        onFormatSelected={exportFormatModal.onFormatSelected}
      />
    </div>
  );
};

// EXISTING_CODE
