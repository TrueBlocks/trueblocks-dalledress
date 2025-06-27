// === SECTION 1: Imports & Dependencies ===
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { GetNamesPage, NamesCrud, Reload } from '@app';
import { Action, Chips, mapNameToChips } from '@components';
import { BaseTab, usePagination } from '@components';
import { ViewStateKey, useFiltering, useSorting } from '@contexts';
// prettier-ignore
import { ActionData, useActionConfig, useActionMsgs, useCrudOperations } from '@hooks';
import { DataFacetConfig, useActiveFacet, useEvent, usePayload } from '@hooks';
import { TabView } from '@layout';
import { useHotkeys } from '@mantine/hooks';
import { crud, msgs, names, types } from '@models';
import { getAddressString, useErrorHandler } from '@utils';

import { getColumns } from './columns';
import { DEFAULT_FACET, ROUTE, namesFacets } from './facets';

// === END SECTION 1 ===

export const Names = () => {
  // === SECTION 2: Hook Initialization ===
  const createPayload = usePayload();

  const activeFacetHook = useActiveFacet({
    facets: namesFacets,
    defaultFacet: DEFAULT_FACET,
    viewRoute: ROUTE,
  });
  const { availableFacets, getCurrentDataFacet } = activeFacetHook;

  const [pageData, setPageData] = useState<names.NamesPage | null>(null);
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
      const result = await GetNamesPage(
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
      case types.DataFacet.ALL:
        return pageData.names || [];
      case types.DataFacet.CUSTOM:
        return pageData.names || [];
      case types.DataFacet.PREFUND:
        return pageData.names || [];
      case types.DataFacet.REGULAR:
        return pageData.names || [];
      case types.DataFacet.BADDRESS:
        return pageData.names || [];
      default:
        return [];
    }
  }, [pageData, getCurrentDataFacet]);
  // === END SECTION 4 ===

  // === SECTION 5: Event Handling ===
  useEvent(
    msgs.EventType.DATA_LOADED,
    (_message: string, payload?: Record<string, unknown>) => {
      if (payload?.collection === 'names') {
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
  const actionConfig = useActionConfig({
    operations: ['delete', 'undelete', 'remove', 'autoname'],
  });

  // prettier-ignore
  const { handleDelete, handleUndelete, handleRemove, handleAutoname } = useCrudOperations({
    collectionName: 'names',
    crudFunc: NamesCrud,
    pageFunc: GetNamesPage,
    pageClass: names.NamesPage,
    emptyItem: types.Name.createFrom({}),
    getCurrentDataFacet,
    pageData,
    setPageData,
    setTotalItems,
    dataFacetRef,
    actionConfig,
  });
  // === END SECTION 6 ===

  // === SECTION 7: Form & UI Handlers ===
  // EXISTING_CODE
  const { emitSuccess, failure } = useActionMsgs('names');
  const handleSubmit = useCallback(
    (data: Record<string, unknown>) => {
      const submittedName = data as unknown as types.Name;
      const addressStr = getAddressString(submittedName.address);

      // Set default source if not provided
      if (!submittedName.source || submittedName.source === '') {
        submittedName.source = 'TrueBlocks';
      }

      const originalNames = [...(pageData?.names || [])];

      // Optimistic UI Update
      let optimisticNames: types.Name[];
      const existingNameIndex = originalNames.findIndex(
        (n) => getAddressString(n.address) === addressStr,
      );

      if (existingNameIndex !== -1) {
        // Update existing name
        optimisticNames = originalNames.map((n, index) =>
          index === existingNameIndex
            ? ({
                ...n,
                ...submittedName,
              } as types.Name)
            : (n as types.Name),
        );
      } else {
        // Add new name
        optimisticNames = [
          submittedName as types.Name,
          ...(originalNames as types.Name[]),
        ];
      }
      setPageData((prev) => {
        if (!prev) return null;
        return new names.NamesPage({
          ...prev,
          names: optimisticNames,
        });
      });

      NamesCrud(
        createPayload(dataFacetRef.current, ''), // TODO: This may not work for AddName...
        crud.Operation.UPDATE,
        submittedName as types.Name,
      )
        .then(async () => {
          const result = await GetNamesPage(
            createPayload(dataFacetRef.current),
            pagination.currentPage * pagination.pageSize,
            pagination.pageSize,
            sort,
            filter ?? '',
          );
          setPageData(result);
          setTotalItems(result.totalItems || 0);

          const displayName = submittedName.name || addressStr;
          emitSuccess('update', displayName);
        })
        .catch((err) => {
          setPageData((prev) => {
            if (!prev) return null;
            return new names.NamesPage({
              ...prev,
              names: originalNames,
            });
          });
          handleError(err, failure('update', '', err.message));
        });
    },
    [
      pageData?.names,
      pagination.currentPage,
      pagination.pageSize,
      sort,
      filter,
      setTotalItems,
      emitSuccess,
      handleError,
      failure,
      createPayload,
    ],
  );

  const { setFiltering } = useFiltering(viewStateKey);

  const currentColumns = useMemo(() => {
    const handleChipClick = (chip: string) => {
      setFiltering(chip);
      Reload(createPayload(dataFacetRef.current)).then(() => {
        fetchData();
      });
    };

    const baseColumns = getColumns(getCurrentDataFacet()).map((col) =>
      col.key === 'chips'
        ? {
            ...col,
            render: (row: Record<string, unknown>) => {
              const nameObject = row as unknown as types.Name;
              const chipItems = mapNameToChips(nameObject);
              return <Chips items={chipItems} onChipClick={handleChipClick} />;
            },
          }
        : col,
    );

    const renderActions = (actionData: ActionData) => {
      const isDeleted = actionData.isDeleted;

      return (
        <div className="action-buttons-container">
          <Action
            icon={isDeleted ? 'Undelete' : 'Delete'}
            onClick={() => {
              if (isDeleted) {
                handleUndelete(actionData.addressStr);
              } else {
                handleDelete(actionData.addressStr);
              }
            }}
            disabled={actionData.isProcessing}
            title={isDeleted ? 'Undelete' : 'Delete'}
            size="sm"
          />
          <Action
            icon="Remove"
            onClick={() => handleRemove(actionData.addressStr)}
            disabled={actionData.isProcessing || !isDeleted}
            title="Remove"
            size="sm"
          />
          <Action
            icon="Autoname"
            onClick={() => handleAutoname(actionData.addressStr)}
            disabled={actionData.isProcessing}
            title="Auto-generate name"
            size="sm"
          />
        </div>
      );
    };

    const getCanRemove = (row: Record<string, unknown>) => {
      const name = row as unknown as types.Name;
      return Boolean(name.deleted);
    };

    return actionConfig.injectActionColumn(
      baseColumns,
      renderActions,
      getCanRemove,
    );
  }, [
    setFiltering,
    fetchData,
    getCurrentDataFacet,
    actionConfig,
    handleDelete,
    handleUndelete,
    handleRemove,
    handleAutoname,
    createPayload,
  ]);
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
    // prettier-ignore
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
