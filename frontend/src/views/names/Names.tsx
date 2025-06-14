// NAMES_ROUTE
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { GetNamesPage, NamesCrud, Reload } from '@app';
import {
  Action,
  BaseTab,
  Chips,
  FormField,
  mapNameToChips,
  usePagination,
} from '@components';
import { TableKey, useAppContext, useFiltering, useSorting } from '@contexts';
import { useActionMsgs, useEvent } from '@hooks';
import { TabView } from '@layout';
import { useHotkeys } from '@mantine/hooks';
import { crud, facets, msgs, names, types } from '@models';
import { getAddressString, useErrorHandler } from '@utils';

import { Address } from '../../types/address';
import { getColumns } from './columns';

type IndexableName = types.Name & Record<string, unknown>;

// Helper function to remove undefined properties from an object
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

export const Names = () => {
  const { lastTab } = useAppContext();
  const [pageData, setPageData] = useState<names.NamesPage | null>(null);
  const [state, setState] = useState<facets.LoadState>();
  const [processingAddresses, setProcessingAddresses] = useState<Set<string>>(
    new Set(),
  );

  const [listKind, setListKind] = useState<types.ListKind>(
    lastTab[NAMES_ROUTE] || NAMES_DEFAULT_LIST,
  );
  const tableKey = useMemo(
    (): TableKey => ({ viewName: NAMES_ROUTE, tabName: listKind }),
    [listKind],
  );

  const { error, handleError, clearError } = useErrorHandler();
  const { pagination, setTotalItems } = usePagination(tableKey);
  const { sort } = useSorting(tableKey);
  const { filter, setFiltering } = useFiltering(tableKey);
  const { emitSuccess, failure } = useActionMsgs('names');

  const listKindRef = useRef(listKind);
  const renderCnt = useRef(0);
  // renderCnt.current++;

  useEffect(() => {
    listKindRef.current = listKind;
  }, [listKind]);

  const fetchData = useCallback(async () => {
    clearError();
    try {
      const result = await GetNamesPage(
        listKindRef.current,
        pagination.currentPage * pagination.pageSize,
        pagination.pageSize,
        sort,
        filter ?? '',
      );
      setState(result.state);
      setPageData(result);
      setTotalItems(result.totalItems || 0);
    } catch (err: unknown) {
      handleError(err, `Failed to fetch ${listKindRef.current}`);
    }
  }, [
    clearError,
    pagination.currentPage,
    pagination.pageSize,
    sort,
    filter,
    setTotalItems,
    handleError,
  ]);

  const currentData = useMemo(() => {
    return pageData?.names || [];
  }, [pageData?.names]);

  useEffect(() => {
    const currentTab = lastTab[NAMES_ROUTE];
    if (currentTab && currentTab !== listKind) {
      setListKind(currentTab);
    }
  }, [lastTab, listKind]);

  useEvent(
    msgs.EventType.DATA_LOADED,
    (_message: string, payload?: types.DataLoadedPayload) => {
      if (payload?.listKind === listKindRef.current) {
        fetchData();
      }
    },
  );

  useEffect(() => {
    fetchData();
  }, [fetchData, listKind]);

  useHotkeys([
    [
      'mod+r',
      () => {
        Reload(listKind).then(() => {
          fetchData();
        });
      },
    ],
  ]);

  // Handle CRUD actions for names
  const handleDelete = useCallback(
    (address: Address) => {
      clearError();
      try {
        const original = [...(pageData?.names || [])];
        const optimisticValues = original.map((name) => {
          const nameAddress = getAddressString(name.address);
          if (nameAddress === address) {
            return { ...name, deleted: true };
          }
          return name;
        });
        setState(facets.LoadState.PENDING);
        setPageData((prev) => {
          if (!prev) return null;
          return new names.NamesPage({
            ...prev,
            names: optimisticValues,
          });
        });
        NamesCrud(
          listKindRef.current,
          crud.Operation.DELETE,
          {} as types.Name,
          address,
        )
          .then(async () => {
            const result = await GetNamesPage(
              listKindRef.current,
              pagination.currentPage * pagination.pageSize,
              pagination.pageSize,
              sort,
              filter ?? '',
            );
            setState(result.state);
            setPageData(result);
            setTotalItems(result.totalItems || 0);
            emitSuccess('delete', address);
          })
          .catch((err) => {
            setState(facets.LoadState.ERROR);
            setPageData((prev) => {
              if (!prev) return null;
              return new names.NamesPage({
                ...prev,
                names: original,
              });
            });
            handleError(err, failure('delete', address, err.message));
          });
      } catch (err: unknown) {
        handleError(err, `Failed to delete name ${address}`);
      }
    },
    [
      clearError,
      pageData?.names,
      handleError,
      pagination.currentPage,
      pagination.pageSize,
      sort,
      filter,
      setTotalItems,
      emitSuccess,
      failure,
    ],
  );

  const handleUndelete = useCallback(
    (address: Address) => {
      clearError();
      try {
        const original = [...(pageData?.names || [])];
        const optimisticValues = original.map((name) => {
          const nameAddress = getAddressString(name.address);
          if (nameAddress === address) {
            return { ...name, deleted: false };
          }
          return name;
        });
        setState(facets.LoadState.PENDING);
        setPageData((prev) => {
          if (!prev) return null;
          return new names.NamesPage({
            ...prev,
            names: optimisticValues,
          });
        });
        NamesCrud(
          listKindRef.current,
          crud.Operation.UNDELETE,
          {} as types.Name,
          address,
        )
          .then(async () => {
            const result = await GetNamesPage(
              listKindRef.current,
              pagination.currentPage * pagination.pageSize,
              pagination.pageSize,
              sort,
              filter ?? '',
            );
            setState(result.state);
            setPageData(result);
            setTotalItems(result.totalItems || 0);
            emitSuccess('undelete', address);
          })
          .catch((err) => {
            setState(facets.LoadState.ERROR);
            setPageData((prev) => {
              if (!prev) return null;
              return new names.NamesPage({
                ...prev,
                names: original,
              });
            });
            handleError(err, failure('undelete', address, err.message));
          });
      } catch (err: unknown) {
        handleError(err, `Failed to undelete name ${address}`);
      }
    },
    [
      clearError,
      pageData?.names,
      pagination.currentPage,
      pagination.pageSize,
      sort,
      filter,
      setTotalItems,
      emitSuccess,
      handleError,
      failure,
    ],
  );

  const handleRemove = useCallback(
    (address: Address) => {
      clearError();
      try {
        const original = [...(pageData?.names || [])];
        const optimisticValues = original.filter((name) => {
          const nameAddress = getAddressString(name.address);
          return nameAddress !== address;
        });
        setState(facets.LoadState.PENDING);
        setPageData((prev) => {
          if (!prev) return null;
          return new names.NamesPage({
            ...prev,
            names: optimisticValues,
          });
        });
        NamesCrud(
          listKindRef.current,
          crud.Operation.REMOVE,
          {} as types.Name,
          address,
        )
          .then(async () => {
            const result = await GetNamesPage(
              listKindRef.current,
              pagination.currentPage * pagination.pageSize,
              pagination.pageSize,
              sort,
              filter ?? '',
            );
            setState(result.state);
            setPageData(result);
            setTotalItems(result.totalItems || 0);
            emitSuccess('remove', address);
          })
          .catch((err) => {
            setState(facets.LoadState.ERROR);
            setPageData((prev) => {
              if (!prev) return null;
              return new names.NamesPage({
                ...prev,
                names: original,
              });
            });
            handleError(err, failure('remove', address, err.message));
          });
      } catch (err: unknown) {
        handleError(err, `Failed to remove name ${address}`);
      }
    },
    [
      clearError,
      pageData?.names,
      pagination.currentPage,
      pagination.pageSize,
      sort,
      filter,
      setTotalItems,
      emitSuccess,
      handleError,
      failure,
    ],
  );

  const handleAutoname = useCallback(
    (address: Address) => {
      clearError();
      try {
        const original = [...(pageData?.names || [])];
        const optimisticValues = original.map((name) => {
          const nameAddress = getAddressString(name.address);
          if (nameAddress === address) {
            return { ...name, name: 'Generating...' };
          }
          return name;
        });
        setState(facets.LoadState.PENDING);
        setPageData((prev) => {
          if (!prev) return null;
          return new names.NamesPage({
            ...prev,
            names: optimisticValues,
          });
        });
        NamesCrud(
          listKindRef.current,
          crud.Operation.AUTONAME,
          {} as types.Name,
          address,
        )
          .then(async () => {
            const result = await GetNamesPage(
              listKindRef.current,
              pagination.currentPage * pagination.pageSize,
              pagination.pageSize,
              sort,
              filter ?? '',
            );
            setState(result.state);
            setPageData(result);
            setTotalItems(result.totalItems || 0);
            emitSuccess('autoname', address);
          })
          .catch((err) => {
            setState(facets.LoadState.ERROR);
            setPageData((prev) => {
              if (!prev) return null;
              return new names.NamesPage({
                ...prev,
                names: original,
              });
            });
            handleError(err, failure('autoname', address, err.message));
          });
      } catch (err: unknown) {
        handleError(err, `Failed to autoname address ${address}`);
      }
    },
    [
      clearError,
      pageData?.names,
      pagination.currentPage,
      pagination.pageSize,
      sort,
      filter,
      setTotalItems,
      emitSuccess,
      handleError,
      failure,
    ],
  );

  const handleNameAction = useCallback(
    (
      address: Address,
      isDeleted: boolean,
      actionType: 'delete' | 'undelete' | 'remove' | 'autoname',
    ) => {
      // Add address to processing set
      setProcessingAddresses((prev) => new Set(prev).add(address));

      try {
        switch (actionType) {
          case 'delete':
            handleDelete(address);
            break;
          case 'undelete':
            handleUndelete(address);
            break;
          case 'remove':
            handleRemove(address);
            break;
          case 'autoname':
            handleAutoname(address);
            break;
        }
      } finally {
        setTimeout(() => {
          setProcessingAddresses((prev) => {
            const newSet = new Set(prev);
            newSet.delete(address);
            return newSet;
          });
        }, 100);
      }
    },
    [handleDelete, handleUndelete, handleRemove, handleAutoname],
  );

  const currentColumns = useMemo(() => {
    const handleChipClick = (chip: string) => {
      setFiltering(chip);
      Reload(listKind).then(() => {
        fetchData();
      });
    };

    const baseColumns = getColumns(listKind).map((col) =>
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

    // Add action buttons render function to the actions column
    const actionsOverride: Partial<FormField> = {
      sortable: false,
      editable: false,
      visible: true,
      render: (row: Record<string, unknown>) => {
        const name = row as unknown as types.Name;
        const addressStr = getAddressString(name.address);
        const isProcessing = processingAddresses.has(addressStr);
        const isDeleted = Boolean(name.deleted);

        return (
          <div className="action-buttons-container">
            <Action
              icon={isDeleted ? 'Undelete' : 'Delete'}
              onClick={() =>
                handleNameAction(
                  addressStr,
                  isDeleted,
                  isDeleted ? 'undelete' : 'delete',
                )
              }
              disabled={isProcessing}
              title={isDeleted ? 'Undelete' : 'Delete'}
              size="sm"
            />
            <Action
              icon="Remove"
              onClick={() => handleNameAction(addressStr, isDeleted, 'remove')}
              disabled={isProcessing || !isDeleted}
              title="Remove"
              size="sm"
            />
            <Action
              icon="Autoname"
              onClick={() =>
                handleNameAction(addressStr, isDeleted, 'autoname')
              }
              disabled={isProcessing}
              title="Auto-generate name"
              size="sm"
            />
          </div>
        );
      },
    };

    return baseColumns.map((col) =>
      col.key === 'actions' ? { ...col, ...actionsOverride } : col,
    );
  }, [
    listKind,
    setFiltering,
    fetchData,
    processingAddresses,
    handleNameAction,
  ]);

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
        listKindRef.current,
        crud.Operation.UPDATE,
        submittedName as types.Name,
        '',
      )
        .then(async () => {
          const result = await GetNamesPage(
            listKindRef.current,
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
          // Revert optimistic update on error
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
    ],
  );

  const perTabTable = useMemo(
    () => (
      <BaseTab
        data={currentData as unknown as Record<string, unknown>[]}
        columns={currentColumns}
        loading={!!pageData?.isFetching}
        error={error}
        onSubmit={handleSubmit}
        tableKey={tableKey}
      />
    ),
    [
      currentData,
      currentColumns,
      pageData?.isFetching,
      error,
      handleSubmit,
      tableKey,
    ],
  );

  const tabs = useMemo(
    () => [
      {
        label: types.ListKind.ALL,
        value: types.ListKind.ALL,
        content: perTabTable,
      },
      {
        label: types.ListKind.CUSTOM,
        value: types.ListKind.CUSTOM,
        content: perTabTable,
      },
      {
        label: types.ListKind.PREFUND,
        value: types.ListKind.PREFUND,
        content: perTabTable,
      },
      {
        label: types.ListKind.REGULAR,
        value: types.ListKind.REGULAR,
        content: perTabTable,
      },
      {
        label: types.ListKind.BADDRESS,
        value: types.ListKind.BADDRESS,
        content: perTabTable,
      },
    ],
    [perTabTable],
  );

  return (
    <div className="mainView">
      {(state as string) === '' && <div>{`state: ${state}`}</div>}
      <TabView tabs={tabs} route={NAMES_ROUTE} />
      {error && (
        <div>
          <h3>{`Error fetching ${listKind}`}</h3>
          <p>{error.message}</p>
        </div>
      )}
      {renderCnt.current > 0 && <div>{`renderCnt: ${renderCnt.current}`}</div>}
    </div>
  );
};

const NAMES_DEFAULT_LIST = types.ListKind.ALL;
const NAMES_ROUTE = '/names';

// NAMES_ROUTE
