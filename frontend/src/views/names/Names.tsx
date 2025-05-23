import { useEffect, useMemo, useRef, useState } from 'react';

import {
  AutonameName,
  DeleteName,
  GetNamesPage,
  RemoveName,
  UndeleteName,
  UpdateName,
} from '@app';
import {
  Action,
  FormField,
  Table,
  TableProvider,
  TagsTable,
  TagsTableHandle,
  usePagination,
} from '@components';
import { TableKey, useAppContext } from '@contexts';
import { TabView } from '@layout';
import { msgs, sorting, types } from '@models';
import { ClearSelectedTag, GetSelectedTag, SetSelectedTag } from '@names';
import { EventsOn } from '@runtime';
import { Log, useEmitters } from '@utils';

import './Names.css';

type IndexableName = types.Name & { [key: string]: unknown };
export const FocusSider = 'focus-tags-table';

// Helper function to remove undefined properties from an object
function removeUndefinedProps(
  obj: Record<string, unknown>,
): Record<string, unknown> {
  const result: Record<string, unknown> = {};
  for (const key in obj) {
    if (
      Object.prototype.hasOwnProperty.call(obj, key) &&
      obj[key] !== undefined
    ) {
      result[key] = obj[key];
    }
  }
  return result;
}

export const Names = () => {
  const { lastTab } = useAppContext();

  const [sort, setSort] = useState<sorting.SortDef | null>(null);
  const [filter, setFilter] = useState('');
  const [names, setNames] = useState<IndexableName[]>([]);
  const [loading, setLoading] = useState(false);
  const [processingAddresses, setProcessingAddresses] = useState<Set<string>>(
    new Set(),
  );
  const [tags, setTags] = useState<string[]>([]);
  const [selectedTag, setSelectedTag] = useState<string | null>(null);
  const [listType, setListType] = useState<ListType>(
    getListTypeFromLabel(lastTab['/names'] || ''),
  );
  const [showTagsView, setShowTagsView] = useState<Record<string, boolean>>({
    All: false,
    Custom: false,
    Prefund: false,
    Regular: false,
    Baddress: false,
  });
  const { emitStatus, emitError } = useEmitters();

  // References for each tab's TagsTable to support targeting specific instances
  const tagsTableRefs = useRef<Record<string, TagsTableHandle | null>>({});
  const mainTableRef = useRef<HTMLDivElement>(null);

  // Single ref for backward compatibility
  const tagsTableRef = useRef<TagsTableHandle>(null);

  const tableKey = useMemo<TableKey>(
    () => ({
      viewName: '/names',
      tabName: listType,
    }),
    [listType],
  );

  const { pagination, setTotalItems, goToPage } = usePagination(tableKey);

  // Load the selected tag when the list type changes
  useEffect(() => {
    const fetchSelectedTag = async () => {
      try {
        GetSelectedTag(listType).then((tag) => {
          setSelectedTag(tag || null);
        });
      } catch (err) {
        Log(`Error fetching selected tag: ${err}`);
        setSelectedTag(null);
      }
    };
    fetchSelectedTag();
  }, [listType]);

  useEffect(() => {
    const loadNames = async () => {
      setLoading(true);
      // Log(`Names:loadNames ${listType}`);
      GetNamesPage(
        listType,
        pagination.currentPage * pagination.pageSize,
        pagination.pageSize,
        sort as sorting.SortDef,
        filter ?? '',
      )
        .then((result) => {
          setNames((result.names || []) as IndexableName[]);
          setTotalItems(result.total || 0);
          setTags(result.tags || []);
        })
        .catch((err) => {
          emitError(err);
          setNames([]);
          setTags([]);
        })
        .finally(() => {
          setLoading(false);
        });
    };
    loadNames();

    // Set up event listener for REFRESH events from backend
    const unsubscribe = EventsOn(msgs.EventType.REFRESH, () => {
      loadNames();
    });
    return unsubscribe;
  }, [
    listType,
    pagination.currentPage,
    pagination.pageSize,
    sort,
    filter,
    selectedTag,
    setTotalItems,
    emitError,
  ]);

  useEffect(() => {
    const currentTabLabel = lastTab['/names'] || '';
    const newListType = getListTypeFromLabel(currentTabLabel);
    setListType(newListType);

    // Fetch the selected tag for the new list type when changing tabs
    const fetchSelectedTag = async () => {
      try {
        GetSelectedTag(newListType).then((tag) => {
          setSelectedTag(tag || null);
        });
      } catch (err) {
        Log(`Error fetching selected tag: ${err}`);
        setSelectedTag(null);
      }
    };
    fetchSelectedTag();
  }, [lastTab]);

  // Handle cmd+shift+e to open/focus tags table
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      // Cmd+Shift+E (Mac) or Ctrl+Shift+E (Win/Linux)
      const isCmdShiftE =
        (e.metaKey || e.ctrlKey) && e.shiftKey && e.key.toLowerCase() === 'e';
      if (isCmdShiftE) {
        e.preventDefault();
        const currentTabLabel = lastTab['/names'] || 'Custom';
        if (!showTagsView[currentTabLabel]) {
          setShowTagsView((prev) => ({
            ...prev,
            [currentTabLabel]: true,
          }));
        } else {
          // Focus the tags table for the current tab
          if (tagsTableRefs.current[currentTabLabel]) {
            tagsTableRefs.current[currentTabLabel].focus();
          } else if (tagsTableRef.current) {
            tagsTableRef.current.focus();
          }
        }
      }
    };
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [lastTab, showTagsView]);

  // Handle Escape to close tags table if focused
  useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        const currentTabLabel = lastTab['/names'] || 'Custom';
        if (showTagsView[currentTabLabel]) {
          // Check if the tags table for this tab is focused using DOM API
          const ref =
            tagsTableRefs.current[currentTabLabel] || tagsTableRef.current;
          // Try to get the root element of the tags table
          let tagsTableEl: HTMLElement | null = null;
          // Use a safer check for 'root' property without using 'any'
          const maybeRoot =
            ref && typeof ref === 'object' && ref !== null && 'root' in ref
              ? (ref as { root?: unknown }).root
              : undefined;
          if (
            maybeRoot &&
            typeof maybeRoot === 'object' &&
            maybeRoot !== null &&
            'contains' in maybeRoot
          ) {
            tagsTableEl = maybeRoot as HTMLElement;
          } else if (ref && 'focus' in ref && typeof ref.focus === 'function') {
            // fallback: try to find by class
            tagsTableEl = document.querySelector(
              '.tags-table',
            ) as HTMLElement | null;
          }
          if (
            tagsTableEl &&
            document.activeElement &&
            tagsTableEl.contains(document.activeElement)
          ) {
            setShowTagsView((prev) => ({
              ...prev,
              [currentTabLabel]: false,
            }));
            e.preventDefault();
          }
        }
      }
    };
    window.addEventListener('keydown', handleEscape);
    return () => window.removeEventListener('keydown', handleEscape);
  }, [lastTab, showTagsView]);

  // Setup event listener for focusing the tags table
  useEffect(() => {
    const focusTagsTable = (_data: { activeTab?: string } = {}) => {
      // First try to derive tab from current location
      const currentTabLabel = lastTab['/names'] || 'Custom';

      // Try to focus the TagsTable for the current tab first
      if (currentTabLabel && tagsTableRefs.current[currentTabLabel]) {
        tagsTableRefs.current[currentTabLabel].focus();
        return;
      }

      // If we couldn't focus the current tab's TagsTable, iterate through all available refs
      const availableTabs = Object.keys(tagsTableRefs.current);
      if (availableTabs.length > 0) {
        // Try to find any tab that has a valid ref
        for (const tabName of availableTabs) {
          if (tagsTableRefs.current[tabName]) {
            tagsTableRefs.current[tabName].focus();
            return;
          }
        }
      }

      // Fall back to legacy ref as a last resort
      if (tagsTableRef.current) {
        tagsTableRef.current.focus();
      }
    };

    const unsubscribe = EventsOn(FocusSider, focusTagsTable);
    return unsubscribe;
  }, [lastTab]);

  // Handle tag selection
  const handleTagSelect = async (
    tag: string | null,
    focusMainTable: boolean = false,
  ) => {
    try {
      if (tag) {
        try {
          SetSelectedTag(listType, tag).then(() => {});
        } catch (e) {
          Log(`Error in SetSelectedTag: ${e}`);
        }
      } else {
        try {
          ClearSelectedTag(listType).then(() => {});
        } catch (e) {
          Log(`Error in ClearSelectedTag: ${e}`);
        }
      }
      setSelectedTag(tag);
    } catch (err) {
      console.error(`Error setting selected tag: ${err}`);
    }

    // Reset to first page when selecting a tag
    if (pagination.currentPage !== 0) {
      goToPage(0);
    }

    // Focus the main table only if explicitly requested (via Enter key)
    if (focusMainTable && mainTableRef.current) {
      setTimeout(() => {
        const tableElement = mainTableRef.current?.querySelector('.data-table');
        if (tableElement) {
          (tableElement as HTMLElement).focus();
        }
      }, 0);
    }
  };

  const handleFormSubmit = (data: Record<string, unknown>) => {
    if (!data.source || data.source === '') {
      data.source = 'TrueBlocks';
    }

    // Cast the form data to IndexableName. This assumes that the data from the form,
    // combined with any defaults, is intended to be treated as an IndexableName for
    // the optimistic update and for the UpdateName API call.
    const submittedName = data as IndexableName;

    const originalNames = [...names]; // 1. Save current names state

    // 2. Optimistic UI Update
    let optimisticNames: IndexableName[];
    const existingNameIndex = originalNames.findIndex(
      // Compare addresses. Assumes 'address' property is comparable and of the same type.
      (n) => n.address === submittedName.address,
    );

    if (existingNameIndex !== -1) {
      // Update existing name
      optimisticNames = originalNames.map((n, index) =>
        index === existingNameIndex
          ? ({ ...n, ...removeUndefinedProps(submittedName) } as IndexableName) // Use helper
          : n,
      );
    } else {
      // Add new name: Prepend submittedName to the list for immediate visibility.
      // The subsequent GetNamesPage call will provide the correctly sorted/paginated list.
      optimisticNames = [submittedName, ...originalNames];
    }
    setNames(optimisticNames); // Update UI optimistically

    // 3. Call UpdateName
    UpdateName(submittedName)
      .then(() => {
        // 4. If UpdateName is successful, call GetNamesPage
        return GetNamesPage(
          listType,
          pagination.currentPage * pagination.pageSize,
          pagination.pageSize,
          sort as sorting.SortDef,
          filter ?? '',
        );
      })
      .then((result) => {
        // 5. Update UI with definitive result from GetNamesPage
        setNames((result.names || []) as IndexableName[]);
        setTotalItems(result.total || 0);
        setTags(result.tags || []);

        // Ensure address is stringified for the status message if it's not already a string.
        const addressStr =
          typeof submittedName.address === 'string'
            ? submittedName.address
            : String(submittedName.address); // Fallback to String() conversion
        emitStatus(
          `Name updated successfully: ${submittedName.name} ${addressStr}`,
        );
      })
      .catch((err) => {
        // 6. If there's an error anywhere in the chain, undo the optimistic setNames
        setNames(originalNames); // Revert to the original names
        emitError(err);
      });
  };

  // Action types for name operations
  type NameActionType = 'delete' | 'edit' | 'remove' | 'autoname';

  const handleAction = (
    address: string,
    isDeleted: boolean,
    actionType: NameActionType = 'delete',
  ) => {
    const addressStr = address;

    // Add address to processing set
    setProcessingAddresses((prev) => new Set(prev).add(addressStr));

    try {
      // Handle different action types
      if (actionType === 'delete') {
        // Existing delete/undelete functionality
        const originalNames = [...names]; // Save current names state

        // Optimistic UI Update - toggle deleted state
        const optimisticNames = originalNames.map((name) => {
          // Convert address to string for comparison
          const nameAddress =
            typeof name.address === 'string'
              ? name.address
              : String(name.address);

          return nameAddress === addressStr
            ? ({ ...name, deleted: !isDeleted } as IndexableName)
            : name;
        });
        setNames(optimisticNames);

        // Determine which API to call based on current state
        const apiCall = isDeleted
          ? UndeleteName(addressStr)
          : DeleteName(addressStr);

        apiCall
          .then(() => {
            // If API call is successful, refresh the data to get the definitive state
            return GetNamesPage(
              listType,
              pagination.currentPage * pagination.pageSize,
              pagination.pageSize,
              sort as sorting.SortDef,
              filter ?? '',
            );
          })
          .then((result) => {
            // Update UI with definitive result from GetNamesPage
            setNames((result.names || []) as IndexableName[]);
            setTotalItems(result.total || 0);
            setTags(result.tags || []);

            const action = isDeleted ? 'undeleted' : 'deleted';
            emitStatus(`Address ${addressStr} was ${action} successfully`);
          })
          .catch((err) => {
            // If there's an error, revert the optimistic update
            setNames(originalNames);
            emitError(err);
            Log(`Error in handleAction: ${err}`);
          });
      }
      // Add implementations for future action types
      else if (actionType === 'edit') {
        // For edit, we could open a form or implement inline editing
        // For now, we'll log that this functionality needs UI implementation
        emitStatus(
          `Edit functionality for ${addressStr} needs UI implementation`,
        );
      } else if (actionType === 'remove') {
        // Remove can only be performed on deleted items
        if (!isDeleted) {
          emitError(
            'Cannot remove a name that is not deleted. Delete it first.',
          );
          return;
        }

        const originalNames = [...names]; // Save current names state

        // Optimistic UI Update - remove the item from the list
        const optimisticNames = originalNames.filter((name) => {
          const nameAddress =
            typeof name.address === 'string'
              ? name.address
              : String(name.address);
          return nameAddress !== addressStr;
        });
        setNames(optimisticNames);

        RemoveName(addressStr)
          .then(() => {
            // If API call is successful, refresh the data to get the definitive state
            return GetNamesPage(
              listType,
              pagination.currentPage * pagination.pageSize,
              pagination.pageSize,
              sort as sorting.SortDef,
              filter ?? '',
            );
          })
          .then((result) => {
            // Update UI with definitive result from GetNamesPage
            setNames((result.names || []) as IndexableName[]);
            setTotalItems(result.total || 0);
            setTags(result.tags || []);

            emitStatus(`Address ${addressStr} was removed successfully`);
          })
          .catch((err) => {
            // If there's an error, revert the optimistic update
            setNames(originalNames);
            emitError(err);
            Log(`Error in handleAction (remove): ${err}`);
          });
      } else if (actionType === 'autoname') {
        const originalNames = [...names]; // Save current names state

        // Optimistic UI Update - update the name to show it's being processed
        const optimisticNames = originalNames.map((name) => {
          const nameAddress =
            typeof name.address === 'string'
              ? name.address
              : String(name.address);

          return nameAddress === addressStr
            ? ({ ...name, name: 'Generating...' } as IndexableName)
            : name;
        });
        setNames(optimisticNames);

        AutonameName(addressStr)
          .then(() => {
            // If API call is successful, refresh the data to get the definitive state
            return GetNamesPage(
              listType,
              pagination.currentPage * pagination.pageSize,
              pagination.pageSize,
              sort as sorting.SortDef,
              filter ?? '',
            );
          })
          .then((result) => {
            // Update UI with definitive result from GetNamesPage
            setNames((result.names || []) as IndexableName[]);
            setTotalItems(result.total || 0);
            setTags(result.tags || []);

            emitStatus(`Address ${addressStr} was auto-named successfully`);
          })
          .catch((err) => {
            // If there's an error, revert the optimistic update
            setNames(originalNames);
            emitError(err);
            Log(`Error in handleAction (autoname): ${err}`);
          });
      }
    } finally {
      // Always clean up the processing state
      setProcessingAddresses((prev) => {
        const newSet = new Set(prev);
        newSet.delete(addressStr);
        return newSet;
      });
    }
  };

  const formValidation = {
    name: (value: unknown) => {
      if (!value || String(value).trim() === '') {
        return 'Name is required';
      }
      return null;
    },
    address: (value: unknown) => {
      if (!value || String(value).trim() === '') {
        return 'Address is required';
      }

      const addressRegex = /^0x[a-fA-F0-9]{40}$/;
      if (!addressRegex.test(String(value))) {
        return 'Invalid Ethereum address format';
      }

      return null;
    },
  };

  const tagsTable = (tabLabel: string) => {
    return (
      <TagsTable
        key={`tags-table-${tabLabel}`}
        tags={tags}
        onTagSelect={handleTagSelect}
        selectedTag={selectedTag}
        visible={true}
        ref={(ref) => {
          if (ref) {
            tagsTableRefs.current[tabLabel] = ref;
          }
          if (tabLabel === 'Custom' && ref) {
            tagsTableRef.current = ref;
          }
        }}
      />
    );
  };

  const dataTable = (tabLabel: string) => {
    // Column override configurations
    const autonameOverride: Partial<FormField<IndexableName>> = {
      width: '40px',
      sortable: false,
      editable: false,
      visible: true,
      render: (row: IndexableName) => {
        const isDeleted = Boolean(row.deleted);
        const addressStr =
          typeof row.address === 'string' ? row.address : String(row.address);
        const isProcessing = processingAddresses.has(addressStr);

        return (
          <div className="action-buttons-container">
            <Action
              icon="Autoname"
              onClick={() => handleAction(addressStr, isDeleted, 'autoname')}
              disabled={isProcessing}
              title="Auto-Update"
              size="sm"
            />
          </div>
        );
      },
    };

    const actionsOverride: Partial<FormField<IndexableName>> = {
      sortable: false,
      editable: false,
      visible: true,
      render: (row: IndexableName) => {
        const isDeleted = Boolean(row.deleted);
        const addressStr =
          typeof row.address === 'string' ? row.address : String(row.address);
        const isProcessing = processingAddresses.has(addressStr);
        return (
          <div className="action-buttons-container">
            <Action
              icon="Edit"
              onClick={() => handleAction(addressStr, isDeleted, 'edit')}
              disabled={isProcessing}
              title="Edit"
              size="sm"
            />
            <Action
              icon={isDeleted ? 'Undelete' : 'Delete'}
              onClick={() => handleAction(addressStr, isDeleted, 'delete')}
              disabled={isProcessing}
              title={isDeleted ? 'Undelete' : 'Delete'}
              size="sm"
            />
            <Action
              icon="Remove"
              onClick={() => handleAction(addressStr, isDeleted, 'remove')}
              disabled={isProcessing || !isDeleted}
              title="Remove"
              size="sm"
            />
          </div>
        );
      },
    };

    const nameColumns: FormField<IndexableName>[] = [
      createColumn('name'),
      createColumn('', autonameOverride),
      createColumn('address', { readOnly: true, width: '350px' }),
      createColumn('tags'),
      createColumn('source', { sameLine: true }),
      createColumn('actions', actionsOverride),
    ];

    return (
      <>
        <Table
          key={`data-table-${tabLabel}`}
          columns={nameColumns}
          data={names}
          loading={loading}
          sort={sort}
          onSortChange={setSort}
          filter={filter}
          onFilterChange={setFilter}
          tableKey={tableKey}
          onSubmit={handleFormSubmit}
          validate={formValidation}
        />
      </>
    );
  };

  const createTableContent = (tabLabel: string) => {
    const tagsVisible = !!showTagsView[tabLabel];
    return (
      <TableProvider key={`table-provider-${tabLabel}`}>
        {tagsVisible ? (
          <div className="dual-table-layout">
            {tagsTable(tabLabel)}
            <div className="table-wrapper" ref={mainTableRef}>
              {dataTable(tabLabel)}
            </div>
          </div>
        ) : (
          dataTable(tabLabel)
        )}
      </TableProvider>
    );
  };

  const tabs = [
    { label: 'All', content: createTableContent('All') },
    { label: 'Custom', content: createTableContent('Custom') },
    { label: 'Prefund', content: createTableContent('Prefund') },
    { label: 'Regular', content: createTableContent('Regular') },
    { label: 'Baddress', content: createTableContent('Baddress') },
  ];

  return (
    <div className="table-tab-view">
      <TabView tabs={tabs} route="/names" />
    </div>
  );
};

const createColumn = (
  baseName: string,
  overrides: Partial<FormField<IndexableName>> = {},
): FormField<IndexableName> => {
  const capitalizedBase = baseName.charAt(0).toUpperCase() + baseName.slice(1);
  return {
    name: baseName.toLowerCase(),
    key: baseName.toLowerCase(),
    header: capitalizedBase,
    label: capitalizedBase,
    type: 'text',
    sortable: true,
    ...overrides,
  };
};

const getListTypeFromLabel = (label: string): ListType => {
  const tabToListType: Record<string, ListType> = {
    All: 'all',
    Custom: 'custom',
    Prefund: 'prefund',
    Regular: 'regular',
    Baddress: 'baddress',
  };
  return tabToListType[label] || 'all';
};

// Define the ListType directly here as it's not exposed from the Wails bindings
export type ListType = 'all' | 'custom' | 'prefund' | 'regular' | 'baddress';
