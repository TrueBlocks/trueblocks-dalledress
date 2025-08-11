// Copyright 2016, 2026 The Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */
// === SECTION 1: Imports & Dependencies ===
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { GetContractsPage, Reload } from '@app';
import { BaseTab, usePagination } from '@components';
import { useFiltering, useSorting } from '@contexts';
import {
  DataFacetConfig,
  toPageDataProp,
  useActiveFacet,
  useColumns,
  useEvent,
  usePayload,
  useViewConfig,
} from '@hooks';
import { TabView } from '@layout';
import { Alert, Container, Stack, Text, Title } from '@mantine/core';
import { useHotkeys } from '@mantine/hooks';
import { contracts } from '@models';
import { msgs, project, types } from '@models';
import { Debugger, useErrorHandler } from '@utils';
import { toProperCase } from 'src/utils/toProper';

import { createDetailPanelFromViewConfig } from '../utils/detailPanel';
import { ContractDashboard, ContractExecute } from './components';

const ROUTE = 'contracts';

// Tiny component for consistent facet titles
interface FacetTitleProps {
  facet: types.DataFacet;
  contractName?: string;
  contractAddress?: string;
}

const FacetTitle: React.FC<FacetTitleProps> = ({
  facet,
  contractName,
  contractAddress,
}) => {
  const getTitleForFacet = (facet: types.DataFacet): string => {
    switch (facet) {
      case types.DataFacet.DASHBOARD:
        return 'Contract Dashboard';
      case types.DataFacet.EXECUTE:
        return 'Contract Interactions';
      case types.DataFacet.EVENTS:
        return 'Contract Events';
      default:
        return 'Contract';
    }
  };

  return (
    <div>
      <Title order={3}>{getTitleForFacet(facet)}</Title>
      {contractName && (
        <Text size="sm" color="dimmed">
          {contractName} ({contractAddress})
        </Text>
      )}
    </div>
  );
};

export const Contracts = () => {
  // === SECTION 2: Hook Initialization ===
  const renderCnt = useRef(0);
  const createPayload = usePayload();

  // === SECTION 2.5: Initial ViewConfig Load ===
  const { config: viewConfig } = useViewConfig({
    viewName: 'contracts',
  });

  // Generate facets from ViewConfig
  const contractsFacets: DataFacetConfig[] = useMemo(() => {
    if (!viewConfig?.facets) {
      // Fallback to default facets if ViewConfig not loaded yet
      return [
        {
          id: types.DataFacet.DASHBOARD,
          label: toProperCase(types.DataFacet.DASHBOARD),
        },
        {
          id: types.DataFacet.EXECUTE,
          label: toProperCase(types.DataFacet.EXECUTE),
        },
        {
          id: types.DataFacet.EVENTS,
          label: toProperCase(types.DataFacet.EVENTS),
        },
      ];
    }
    // Maintain specific order for contracts facets
    const orderedKeys = ['dashboard', 'execute', 'events'];
    const availableKeys = Object.keys(viewConfig.facets);
    const sortedKeys = orderedKeys.filter((key) => availableKeys.includes(key));

    return sortedKeys.map((facetKey) => ({
      id: facetKey as types.DataFacet,
      label: viewConfig.facets[facetKey]?.name || toProperCase(facetKey),
    }));
  }, [viewConfig]);

  const activeFacetHook = useActiveFacet({
    facets: contractsFacets,
    viewRoute: ROUTE,
  });
  const { availableFacets, getCurrentDataFacet } = activeFacetHook;

  const [pageData, setPageData] = useState<contracts.ContractsPage | null>(
    null,
  );
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

  // === SECTION 3: Data Fetching ===
  const fetchData = useCallback(async () => {
    clearError();
    try {
      const result = await GetContractsPage(
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
      case types.DataFacet.DASHBOARD:
        return pageData.contracts || [];
      case types.DataFacet.EXECUTE:
        return pageData.contracts || [];
      case types.DataFacet.EVENTS:
        return pageData.logs || [];
      default:
        return [];
    }
  }, [pageData, getCurrentDataFacet]);

  // === SECTION 4: Event Handling ===
  useEvent(
    msgs.EventType.DATA_LOADED,
    (_message: string, payload?: Record<string, unknown>) => {
      if (payload?.collection === 'contracts') {
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

  // === SECTION 5: CRUD Operations ===
  // TODO: Add CRUD operations for Contracts view when needed

  // === SECTION 6: UI Configuration ===
  const currentColumns = useColumns(
    viewConfig?.facets?.[getCurrentDataFacet()]?.columns || [],
    {
      showActions: false,
      actions: [],
      getCanRemove: useCallback((_row: unknown) => false, []),
    },
    {},
    toPageDataProp(pageData),
    { rowActions: [] },
  );

  // Create detail panel from ViewConfig
  const detailPanel = useMemo(() => {
    return createDetailPanelFromViewConfig(
      viewConfig,
      getCurrentDataFacet,
      'Contract Details',
    );
  }, [viewConfig, getCurrentDataFacet]);

  const perTabContent = useMemo(() => {
    const currentFacet = getCurrentDataFacet();
    const facetConfig = viewConfig?.facets?.[currentFacet];

    // Check if this facet should be displayed as a form using ViewConfig
    if (facetConfig?.isForm) {
      if (!pageData?.contracts || pageData.contracts.length === 0) {
        return (
          <Container size="lg" py="xl">
            <Alert color="yellow" title="No contracts available">
              No contracts available for {facetConfig.name?.toLowerCase()} view
            </Alert>
          </Container>
        );
      }

      // Use the first contract or selected contract
      const contractState = pageData.contracts[0];
      if (!contractState) {
        return (
          <Container size="lg" py="xl">
            <Alert color="yellow" title="No contract data">
              No contract data available
            </Alert>
          </Container>
        );
      }

      return (
        <Container size="lg" py="xl">
          <Stack gap="md">
            <FacetTitle
              facet={currentFacet}
              contractName={contractState.name}
              contractAddress={contractState.address?.toString()}
            />
            {currentFacet === types.DataFacet.DASHBOARD ? (
              <ContractDashboard
                contractState={contractState}
                onRefresh={() => fetchData()}
              />
            ) : currentFacet === types.DataFacet.EXECUTE ? (
              <ContractExecute
                contractState={contractState}
                functionName="all"
              />
            ) : null}
          </Stack>
        </Container>
      );
    }

    // For Events and other facets, use regular table view
    return (
      <BaseTab<Record<string, unknown>>
        data={currentData as unknown as Record<string, unknown>[]}
        columns={currentColumns}
        loading={!!pageData?.isFetching}
        error={error}
        viewStateKey={viewStateKey}
        headerActions={[]}
        detailPanel={detailPanel}
      />
    );
  }, [
    getCurrentDataFacet,
    viewConfig?.facets,
    pageData,
    currentData,
    currentColumns,
    error,
    viewStateKey,
    detailPanel,
    fetchData,
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
    </div>
  );
};

// EXISTING_CODE

// EXISTING_CODE
