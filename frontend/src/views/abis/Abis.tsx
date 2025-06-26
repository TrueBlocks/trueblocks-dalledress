// === SECTION 1: Imports & Dependencies ===
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { AbisCrud, GetAbisPage, Reload } from '@app';
import { Action } from '@components';
import { BaseTab, usePagination } from '@components';
import { ViewStateKey, useFiltering, useSorting } from '@contexts';
import { ActionData, useActionConfig, useActionMsgs } from '@hooks';
import { DataFacetConfig, useActiveFacet, useEvent, usePayload } from '@hooks';
import { TabView } from '@layout';
import { useHotkeys } from '@mantine/hooks';
import { abis, crud, msgs, types } from '@models';
import { getAddressString, useErrorHandler } from '@utils';

import { getColumns } from './columns';
import { DEFAULT_FACET, ROUTE, abisFacets } from './facets';

// === END SECTION 1 ===

export const Abis = () => {
  // === SECTION 2: Hook Initialization ===
  const createPayload = usePayload();

  const activeFacetHook = useActiveFacet({
    facets: abisFacets,
    defaultFacet: DEFAULT_FACET,
    viewRoute: ROUTE,
  });
  const { availableFacets, getCurrentDataFacet } = activeFacetHook;

  const [pageData, setPageData] = useState<abis.AbisPage | null>(null);
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
  const dataFacetRef = useRef(getCurrentDataFacet() as types.DataFacet);
  useEffect(() => {
    dataFacetRef.current = getCurrentDataFacet() as types.DataFacet;
  }, [getCurrentDataFacet]);
  // === END SECTION 3 ===

  // === SECTION 4: Data Fetching Logic ===
  const fetchData = useCallback(async () => {
    clearError();
    try {
      const result = await GetAbisPage(
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
      case types.DataFacet.DOWNLOADED:
        return pageData.abis || [];
      case types.DataFacet.KNOWN:
        return pageData.abis || [];
      case types.DataFacet.FUNCTIONS:
        return pageData.functions || [];
      case types.DataFacet.EVENTS:
        return pageData.functions || [];
      default:
        return [];
    }
  }, [pageData, getCurrentDataFacet]);
  // === END SECTION 4 ===

  // === SECTION 5: Event Handling ===
  useEvent(
    msgs.EventType.DATA_LOADED,
    (_message: string, payload?: Record<string, unknown>) => {
      if (payload?.collection === 'abis') {
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
        fetchData();
      });
    } catch (err: unknown) {
      handleError(err, `Failed to reload ${getCurrentDataFacet()}`);
    }
  }, [getCurrentDataFacet, createPayload, fetchData, handleError]);

  useHotkeys([['mod+r', handleReload]]);
  // === END SECTION 5 ===

  // === SECTION 6: CRUD Operations ===
  // EXISTING_CODE
  const actionConfig = useActionConfig({
    operations: ['remove'],
  });

  const { emitSuccess, failure } = useActionMsgs('abis');
  const { goToPage } = usePagination(viewStateKey);

  const handleRemove = useCallback(
    (address: string) => {
      clearError();
      actionConfig.startProcessing(address);

      try {
        const original = [...(pageData?.abis || [])];
        const isOnlyRowOnPage = original.length === 1;
        const optimisticValues = original.filter((abi) => {
          const abiAddress = getAddressString(abi.address);
          return abiAddress !== address;
        });

        setPageData((prev) => {
          if (!prev) return null;
          return new abis.AbisPage({
            ...prev,
            abis: optimisticValues,
          });
        });
        const currentTotal = pageData?.totalItems || 0;
        setTotalItems(Math.max(0, currentTotal - 1));

        AbisCrud(
          createPayload(dataFacetRef.current, address),
          crud.Operation.REMOVE,
          {} as types.Abi,
        )
          .then(async () => {
            const result = await GetAbisPage(
              createPayload(dataFacetRef.current),
              pagination.currentPage * pagination.pageSize,
              pagination.pageSize,
              sort,
              filter,
            );
            setPageData(result);
            setTotalItems(result.totalItems || 0);

            if (isOnlyRowOnPage && result.totalItems > 0) {
              const newTotalPages = Math.ceil(
                result.totalItems / pagination.pageSize,
              );
              const lastPageIndex = Math.max(0, newTotalPages - 1);

              if (lastPageIndex !== pagination.currentPage) {
                goToPage(lastPageIndex);
              }
            }
            emitSuccess('remove', address);
          })
          .catch((err) => {
            setPageData((prev) => {
              if (!prev) return null;
              return new abis.AbisPage({
                ...prev,
                abis: original,
              });
            });
            setTotalItems(pageData?.totalItems || 0);
            handleError(err, failure('remove', address, err.message));
          })
          .finally(() => {
            setTimeout(() => {
              actionConfig.stopProcessing(address);
            }, 100);
          });
      } catch (err: unknown) {
        handleError(err, `Failed to remove abi ${address}`);
        actionConfig.stopProcessing(address);
      }
    },
    [
      clearError,
      actionConfig,
      pageData?.abis,
      pageData?.totalItems,
      setTotalItems,
      pagination.currentPage,
      pagination.pageSize,
      sort,
      filter,
      emitSuccess,
      failure,
      goToPage,
      handleError,
      createPayload,
    ],
  );
  // EXISTING_CODE
  // === END SECTION 6 ===

  // === SECTION 7: Form & UI Handlers ===
  // EXISTING_CODE
  const handleSubmit = useCallback(
    (_formData: Record<string, unknown>) => {},
    [],
  );

  const currentColumns = useMemo(() => {
    const baseColumns = getColumns(
      pageData?.facet || types.DataFacet.DOWNLOADED,
    );

    const shouldShowActions =
      (pageData?.facet || types.DataFacet.DOWNLOADED) ===
        types.DataFacet.DOWNLOADED ||
      (pageData?.facet || types.DataFacet.DOWNLOADED) === types.DataFacet.KNOWN;

    if (!shouldShowActions) {
      return actionConfig.injectActionColumn(baseColumns, () => null);
    }

    const renderActions = (actionData: ActionData) => {
      const canRemove = pageData?.facet === types.DataFacet.DOWNLOADED;

      return (
        <div className="action-buttons-container">
          <Action
            icon="Remove"
            onClick={() => handleRemove(actionData.addressStr)}
            disabled={actionData.isProcessing || !canRemove}
            title="Remove"
            size="sm"
          />
        </div>
      );
    };

    const getCanRemove = (_row: Record<string, unknown>) => {
      return pageData?.facet === types.DataFacet.DOWNLOADED;
    };

    return actionConfig.injectActionColumn(
      baseColumns,
      renderActions,
      getCanRemove,
    );
  }, [pageData?.facet, handleRemove, actionConfig]);
  // EXISTING_CODE
  // === END SECTION 7 ===

  // === SECTION 8: Tab Configuration ===
  const perTabTable = useMemo(
    () => (
      <BaseTab
        data={currentData as unknown as Record<string, unknown>[]}
        columns={currentColumns}
        loading={!!pageData?.isFetching}
        error={error}
        onSubmit={handleSubmit}
        viewStateKey={viewStateKey}
      />
    ),
    [
      currentData,
      currentColumns,
      pageData?.isFetching,
      error,
      handleSubmit,
      viewStateKey,
    ],
  );

  const tabs = useMemo(
    () =>
      availableFacets.map((facetConfig: DataFacetConfig) => ({
        label: facetConfig.label,
        value: facetConfig.id,
        content: perTabTable,
      })),
    [availableFacets, perTabTable],
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
