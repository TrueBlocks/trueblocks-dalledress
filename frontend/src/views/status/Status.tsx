// Copyright 2016, 2026 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */
// === SECTION 1: Imports & Dependencies ===
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { GetStatusPage, Reload } from '@app';
import { BaseTab, usePagination } from '@components';
import { ViewStateKey, useFiltering, useSorting } from '@contexts';
import { toPageDataProp, useColumns } from '@hooks';
// prettier-ignore
import { useActionConfig } from '@hooks';
import { DataFacetConfig, useActiveFacet, useEvent, usePayload } from '@hooks';
import { FormView, TabView } from '@layout';
import { useHotkeys } from '@mantine/hooks';
import { msgs, status, types } from '@models';
import { useErrorHandler } from '@utils';

import { getColumns } from './columns';
import { DEFAULT_FACET, ROUTE, statusFacets } from './facets';

// === END SECTION 1 ===

export const Status = () => {
  // === SECTION 2: Hook Initialization ===
  const createPayload = usePayload();

  const activeFacetHook = useActiveFacet({
    facets: statusFacets,
    defaultFacet: DEFAULT_FACET,
    viewRoute: ROUTE,
  });
  const { availableFacets, getCurrentDataFacet } = activeFacetHook;

  const [pageData, setPageData] = useState<status.StatusPage | null>(null);
  const viewStateKey = useMemo(
    (): ViewStateKey => ({
      viewName: ROUTE,
      tabName: getCurrentDataFacet(),
    }),
    [getCurrentDataFacet],
  );

  const { error, handleError, clearError } = useErrorHandler();
  const { pagination, setTotalItems } = usePagination(viewStateKey);
  const { sort } = useSorting(viewStateKey);
  const { filter } = useFiltering(viewStateKey);
  // === END SECTION 2 ===

  // === SECTION 3: Refs & Effects Setup ===
  const dataFacetRef = useRef(getCurrentDataFacet());
  useEffect(() => {
    dataFacetRef.current = getCurrentDataFacet();
  }, [getCurrentDataFacet]);
  // === END SECTION 3 ===

  // === SECTION 4: Data Fetching Logic ===
  const fetchData = useCallback(async () => {
    clearError();
    try {
      const result = await GetStatusPage(
        createPayload(dataFacetRef.current),
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
    pagination.currentPage,
    pagination.pageSize,
    sort,
    filter,
    setTotalItems,
    handleError,
    getCurrentDataFacet,
  ]);

  const currentData = useMemo(() => {
    if (!pageData) return [];

    const facet = getCurrentDataFacet();
    switch (facet) {
      case types.DataFacet.STATUS:
        return pageData.status || [];
      case types.DataFacet.CACHES:
        return pageData.caches || [];
      case types.DataFacet.CHAINS:
        return pageData.chains || [];
      default:
        return [];
    }
  }, [pageData, getCurrentDataFacet]);
  // === END SECTION 4 ===

  // === SECTION 5: Event Handling ===
  useEvent(
    msgs.EventType.DATA_LOADED,
    (_message: string, payload?: Record<string, unknown>) => {
      if (payload?.collection === 'status') {
        const eventDataFacet = payload.dataFacet;
        if (eventDataFacet === dataFacetRef.current) {
          fetchData();
        }
      }
    },
  );

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const handleReload = useCallback(async () => {
    try {
      Reload(createPayload(dataFacetRef.current)).then(() => {
        // The data will reload when the DataLoaded event is fired.
      });
    } catch (err: unknown) {
      handleError(err, `Failed to reload ${getCurrentDataFacet()}`);
    }
  }, [getCurrentDataFacet, createPayload, handleError]);

  useHotkeys([['mod+r', handleReload]]);
  // === END SECTION 5 ===

  // === SECTION 6: CRUD Operations ===
  const actionConfig = useActionConfig({
    operations: [],
  });
  // === END SECTION 6 ===

  // === SECTION 7: Form & UI Handlers ===
  const showActions = false;
  const getCanRemove = (_row: unknown): boolean => {
    return true;
  };

  const currentColumns = useColumns(
    getColumns(getCurrentDataFacet()),
    {
      showActions,
      actions: [],
      getCanRemove,
    },
    {},
    toPageDataProp(pageData),
    actionConfig,
    true /* perRowCrud */,
  );
  // === END SECTION 7 ===

  // === SECTION 8: Tab Configuration ===
  const isForm = useCallback((facet: types.DataFacet) => {
    switch (facet) {
      case types.DataFacet.STATUS:
        return true;
      default:
        return false;
    }
  }, []);
  const perTabContent = useMemo(() => {
    const facet = getCurrentDataFacet();
    if (isForm(facet)) {
      const statusData = currentData[0] as unknown as Record<string, unknown>;
      if (!statusData) {
        return <div>No status data available</div>;
      }
      const fieldsWithValues = getColumns(getCurrentDataFacet()).map(
        (field) => ({
          ...field,
          value:
            (statusData?.[field.name as string] as
              | string
              | number
              | boolean
              | undefined) || field.value,
          readOnly: true,
        }),
      );
      return (
        <FormView
          title="Status Information"
          formFields={fieldsWithValues}
          onSubmit={() => {}}
        />
      );
    } else {
      return (
        <BaseTab
          data={currentData as unknown as Record<string, unknown>[]}
          columns={currentColumns}
          loading={!!pageData?.isFetching}
          error={error}
          viewStateKey={viewStateKey}
        />
      );
    }
  }, [
    currentData,
    currentColumns,
    pageData?.isFetching,
    error,
    viewStateKey,
    isForm,
    getCurrentDataFacet,
  ]);

  const tabs = useMemo(
    () =>
      availableFacets.map((facetConfig: DataFacetConfig) => ({
        label: facetConfig.label,
        value: facetConfig.id,
        content: perTabContent,
      })),
    [availableFacets, perTabContent],
  );
  // === END SECTION 8 ===

  // === SECTION 9: Render/JSX ===
  const renderCnt = useRef(0);
  // renderCnt.current++;
  return (
    <div className="mainView">
      <TabView tabs={tabs} route={ROUTE} />
      {error && (
        <div>
          <h3>{`Error fetching ${getCurrentDataFacet()}`}</h3>
          <p>{error.message}</p>
        </div>
      )}
      {renderCnt.current > 0 && <div>{`renderCnt: ${renderCnt.current}`}</div>}
    </div>
  );
  // === END SECTION 9 ===
};

// EXISTING_CODE
// EXISTING_CODE
