// === SECTION 1: Imports & Dependencies ===
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { GetNamesPage, NamesCrud, Reload } from '@app';
import { Action, Chips, mapNameToChips } from '@components';
import { BaseTab, usePagination } from '@components';
import { ViewStateKey, useFiltering, useSorting } from '@contexts';
import {
  ActionData,
  useActionConfig,
  useActionMsgs,
  useCrudOperations,
} from '@hooks';
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
  const dataFacetRef = useRef(getCurrentDataFacet() as types.DataFacet);
  useEffect(() => {
    dataFacetRef.current = getCurrentDataFacet() as types.DataFacet;
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
      if (payload?.collection === ROUTE) {
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
    operations: ['delete', 'undelete', 'remove'],
  });

  // Use the new CRUD operations hook for handleRemove
  const { handleRemove } = useCrudOperations({
    collectionName: ROUTE,
    getCurrentDataFacet,
    pageData,
    setPageData,
    setTotalItems,
    crudFunction: NamesCrud,
    getPageFunction: GetNamesPage,
    dataFacetRef,
    actionConfig,
    PageClass: names.NamesPage,
    emptyItem: types.Name.createFrom({}),
  });

  const { emitSuccess, failure } = useActionMsgs('names');

  const handleDelete = useCallback(
    (address: string) => {
      clearError();
      actionConfig.startProcessing(address);

      try {
        const original = [...(pageData?.names || [])];
        const optimisticValues = original.map((name) => {
          const nameAddress = getAddressString(name.address);
          if (nameAddress === address) {
            return { ...name, deleted: true };
          }
          return name;
        });
        setPageData((prev) => {
          if (!prev) return null;
          return new names.NamesPage({
            ...prev,
            names: optimisticValues,
          });
        });
        NamesCrud(
          createPayload(dataFacetRef.current, address),
          crud.Operation.DELETE,
          {} as types.Name,
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
            emitSuccess('delete', address);
          })
          .catch((err) => {
            setPageData((prev) => {
              if (!prev) return null;
              return new names.NamesPage({
                ...prev,
                names: original,
              });
            });
            handleError(err, failure('delete', address, err.message));
          })
          .finally(() => {
            setTimeout(() => {
              actionConfig.stopProcessing(address);
            }, 100);
          });
      } catch (err: unknown) {
        handleError(err, `Failed to delete name ${address}`);
        actionConfig.stopProcessing(address);
      }
    },
    [
      clearError,
      actionConfig,
      pageData?.names,
      handleError,
      pagination.currentPage,
      pagination.pageSize,
      sort,
      filter,
      setTotalItems,
      emitSuccess,
      failure,
      createPayload,
    ],
  );

  const handleUndelete = useCallback(
    (address: string) => {
      clearError();
      actionConfig.startProcessing(address);

      try {
        const original = [...(pageData?.names || [])];
        const optimisticValues = original.map((name) => {
          const nameAddress = getAddressString(name.address);
          if (nameAddress === address) {
            return { ...name, deleted: false };
          }
          return name;
        });
        setPageData((prev) => {
          if (!prev) return null;
          return new names.NamesPage({
            ...prev,
            names: optimisticValues,
          });
        });
        NamesCrud(
          createPayload(dataFacetRef.current, address),
          crud.Operation.UNDELETE,
          {} as types.Name,
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
            emitSuccess('undelete', address);
          })
          .catch((err) => {
            setPageData((prev) => {
              if (!prev) return null;
              return new names.NamesPage({
                ...prev,
                names: original,
              });
            });
            handleError(err, failure('undelete', address, err.message));
          })
          .finally(() => {
            setTimeout(() => {
              actionConfig.stopProcessing(address);
            }, 100);
          });
      } catch (err: unknown) {
        handleError(err, `Failed to undelete name ${address}`);
        actionConfig.stopProcessing(address);
      }
    },
    [
      clearError,
      actionConfig,
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

  const handleAutoname = useCallback(
    (address: string) => {
      clearError();
      actionConfig.startProcessing(address);

      try {
        const original = [...(pageData?.names || [])];
        const optimisticValues = original.map((name) => {
          const nameAddress = getAddressString(name.address);
          if (nameAddress === address) {
            return { ...name, name: 'Generating...' };
          }
          return name;
        });
        setPageData((prev) => {
          if (!prev) return null;
          return new names.NamesPage({
            ...prev,
            names: optimisticValues,
          });
        });
        NamesCrud(
          createPayload(dataFacetRef.current, address),
          crud.Operation.AUTONAME,
          {} as types.Name,
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
            emitSuccess('autoname', address);
          })
          .catch((err) => {
            setPageData((prev) => {
              if (!prev) return null;
              return new names.NamesPage({
                ...prev,
                names: original,
              });
            });
            handleError(err, failure('autoname', address, err.message));
          })
          .finally(() => {
            setTimeout(() => {
              actionConfig.stopProcessing(address);
            }, 100);
          });
      } catch (err: unknown) {
        handleError(err, `Failed to autoname address ${address}`);
        actionConfig.stopProcessing(address);
      }
    },
    [
      clearError,
      actionConfig,
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
  // EXISTING_CODE
  // === END SECTION 6 ===

  // === SECTION 7: Form & UI Handlers ===
  // EXISTING_CODE
  type IndexableName = types.Name & Record<string, unknown>;

  function removeUndefinedProps(
    obj: Record<string, unknown>,
  ): Record<string, unknown> {
    const result: Record<string, unknown> = {};
    for (const [key, value] of Object.entries(obj)) {
      if (value !== undefined) {
        result[key] = value;
      }
    }
    return result;
  }

  const handleSubmit = useCallback(
    (data: Record<string, unknown>) => {
      const submittedName = data as IndexableName;
      const addressStr = getAddressString(submittedName.address);

      // Set default source if not provided
      if (!submittedName.source || submittedName.source === '') {
        submittedName.source = 'TrueBlocks';
      }

      const originalNames = [...(pageData?.names || [])];

      // Optimistic UI Update
      let optimisticNames: IndexableName[];
      const existingNameIndex = originalNames.findIndex(
        (n) => getAddressString(n.address) === addressStr,
      );

      if (existingNameIndex !== -1) {
        // Update existing name
        optimisticNames = originalNames.map((n, index) =>
          index === existingNameIndex
            ? ({
                ...n,
                ...removeUndefinedProps(submittedName),
              } as IndexableName)
            : (n as IndexableName),
        );
      } else {
        // Add new name
        optimisticNames = [
          submittedName as IndexableName,
          ...(originalNames as IndexableName[]),
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

    const baseColumns = getColumns(
      getCurrentDataFacet() as types.DataFacet,
    ).map((col) =>
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
