import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import {
  AutonameName,
  CleanNames,
  DeleteName,
  Logger,
  PublishNames,
  Reload,
  RemoveName,
  SaveName,
  UndeleteName,
} from '@app';
import {
  ModifyIcons,
  Table,
  TableProvider,
  TagsTable,
  TagsTableHandle,
  usePagination,
} from '@components';
import { TableKey, useAppContext } from '@contexts';
import { TabView } from '@layout';
import { msgs, sorting, types } from '@models';
import { base } from '@models';
import {
  ClearSelectedTag,
  GetNamesPage,
  GetSelectedTag,
  SetSelectedTag,
} from '@names';
import { EventsOn } from '@runtime';

import { createEnhancedName } from '../utils/NameTypeUtils';
import './Names.css';

export const FocusSider = 'focus-tags-table';

export const Names = () => {
  const { lastTab } = useAppContext();

  const [sort, setSort] = useState<sorting.SortDef | null>(null);
  const [filter, setFilter] = useState('');
  const [names, setNames] = useState<types.Name[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
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
  const [forceRenderKey, setForceRenderKey] = useState(0);

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

  // Add a ref to hold the names array so we can log its length without adding names as a dependency
  const namesRef = useRef(names);
  useEffect(() => {
    // Update the ref whenever names change
    namesRef.current = names;
  }, [names]);

  // Load the selected tag when the list type changes
  useEffect(() => {
    const fetchSelectedTag = async () => {
      try {
        const tag = await GetSelectedTag(listType);
        setSelectedTag(tag || null);
      } catch (err) {
        console.error('Error fetching selected tag:', err);
        setSelectedTag(null);
      }
    };
    fetchSelectedTag();
  }, [listType]);

  // Fix the loadNames function inside the useEffect
  useEffect(() => {
    const loadNames = async () => {
      setLoading(true);
      setError('');
      Logger(
        `Loading names for listType: ${listType}, page: ${pagination.currentPage}, filter: '${filter}'`,
      );

      try {
        const result = await GetNamesPage(
          listType,
          pagination.currentPage,
          pagination.pageSize,
          sort as sorting.SortDef,
          filter,
        );

        Logger(
          `DEBUG_DELETED_STATUS1: GetNamesPage result: ` +
            JSON.stringify(
              result.names.map((n) => ({
                address: n.address,
                deleted: n.deleted,
                name: n.name,
              })),
            ),
        );

        const processedNames = result.names.map((name) => {
          const newName = types.Name.createFrom(name);
          newName.deleted = Boolean(newName.deleted); // Ensure boolean

          if (newName.deleted === true) {
            Logger(
              `DEBUG_DELETED_STATUS2: Processing incoming deleted record: Address: ${newName.address}, Name: ${newName.name}, Deleted: ${newName.deleted}, Type: ${typeof newName.deleted}`,
            );
          }
          return newName;
        });

        Logger(
          `DEBUG_DELETED_STATUS3: processedNames before setNames: ` +
            JSON.stringify(
              processedNames.map((n) => ({
                address: n.address,
                deleted: n.deleted,
                name: n.name,
              })),
            ),
        );

        setNames(processedNames);

        setTimeout(() => {
          Logger(
            `DEBUG_DELETED_STATUS4: names state after setNames (via namesRef.current): ` +
              JSON.stringify(
                namesRef.current.map((n) => ({
                  address: n.address,
                  deleted: n.deleted,
                  name: n.name,
                })),
              ),
          );
        }, 0);

        setTotalItems(result.total);
      } catch (err) {
        setError('Failed to load names');
        Logger('Error loading names:' + err);
      } finally {
        setLoading(false);
        setForceRenderKey((prevKey) => prevKey + 1);
      }
    };

    loadNames();

    var unsubscribe = EventsOn(msgs.EventType.REFRESH, () => {
      Logger('Got refresh event - loading names');
      Logger(`Current names count before refresh: ${namesRef.current.length}`);
      Logger(`Current selected tag: ${selectedTag || 'none'}`);
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
  ]);

  const [editRowData, setEditRowData] = useState<types.Name | null>(null);

  // Reset editRowData when modal is closed
  useEffect(() => {
    const handleModalClose = () => {
      // Add a small delay to ensure we don't reset while the modal is still potentially processing
      setTimeout(() => {
        setEditRowData(null);
      }, 100);
    };

    // Listen for Escape key to detect modal closing
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && editRowData) {
        handleModalClose();
      }
    };

    // Modal could also close via save button or cancel button
    const handleClick = (e: MouseEvent) => {
      // Check if we clicked a button in the modal that might close it
      const target = e.target as HTMLElement;
      if (
        editRowData &&
        (target.matches('.mantine-Modal button') ||
          target.closest('.mantine-Modal button'))
      ) {
        handleModalClose();
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    document.addEventListener('click', handleClick, true);

    return () => {
      window.removeEventListener('keydown', handleKeyDown);
      document.removeEventListener('click', handleClick, true);
    };
  }, [editRowData]);

  useEffect(() => {
    const currentTabLabel = lastTab['/names'] || '';
    const newListType = getListTypeFromLabel(currentTabLabel);
    setListType(newListType);

    // Fetch the selected tag for the new list type when changing tabs
    const fetchSelectedTag = async () => {
      try {
        const tag = await GetSelectedTag(newListType);
        setSelectedTag(tag || null);
      } catch (err) {
        console.error('Error fetching selected tag:', err);
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
          await SetSelectedTag(listType, tag);
        } catch (e) {
          console.error(`Error in SetSelectedTag: ${e}`);
        }
      } else {
        try {
          await ClearSelectedTag(listType);
        } catch (e) {
          console.error(`Error in ClearSelectedTag: ${e}`);
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

  const handleSaveRow = (
    row: Record<string, unknown>,
    updated: Partial<Record<string, unknown>>,
  ) => {
    console.log('Saving row:', row, 'with updates:', updated);
    // TODO: Implement actual save to backend later

    // Reset editRowData after saving
    setEditRowData(null);
  };

  const handleSubmit = (data: Record<string, unknown>) => {
    // Convert the data to a proper Name object using our enhanced utility
    const nameData = createEnhancedName(data);
    SaveName(nameData).then(() => {
      Logger('Front end got returned from SaveName');
      // Reset editRowData after saving
      setEditRowData(null);
    });
  };

  // Handle the cancel action for the edit modal
  const handleCancelRow = useCallback(() => {
    // Reset editRowData when modal is closed
    setEditRowData(null);
  }, [setEditRowData]);

  const handleAction = useCallback(
    (type: string, row?: types.Name, tabNameParam?: string) => {
      const currentTab = tabNameParam || listType;
      let addressString = '';
      if (row && row.address) {
        // Ensure address is consistently handled as a string for operations/logging
        if (typeof row.address === 'object' && 'address' in row.address) {
          const ba = row.address as base.Address;
          if (ba.address) {
            addressString =
              '0x' +
              ba.address.map((b) => b.toString(16).padStart(2, '0')).join('');
          }
        } else {
          addressString = String(row.address);
        }
      }

      if (row) {
        Logger(`Action ${type} clicked for ${addressString}`);
      } else {
        Logger(`Action ${type} clicked for tab ${currentTab}`);
      }

      switch (type) {
        case 'add':
          Logger('Add name action clicked');
          // Create a new empty name object using the createFrom method
          const newName = types.Name.createFrom({
            address: base.Address.createFrom({ address: [] }), // Ensure address is base.Address
            name: '',
            decimals: 0,
            source: 'TrueBlocks Browse',
            symbol: '',
            tags: 'User-Defined',
            isCustom: true,
          });
          setEditRowData(newName);
          break;
        case 'edit':
          if (row) {
            Logger(`Edit action clicked for: ${addressString}`);
            // Directly set the row to edit rather than simulating key press
            setEditRowData(row);
          }
          break;
        case 'delete':
          Logger(`Action ${type} clicked for ${addressString}`);
          if (addressString) {
            DeleteName(addressString)
              .then(() => {
                // Optimistic update
                setNames((prev) => {
                  return prev.map((nameEntry) => {
                    let nameEntryAddressString = '';
                    if (
                      nameEntry.address &&
                      typeof nameEntry.address === 'object' &&
                      'address' in nameEntry.address
                    ) {
                      const ba = nameEntry.address as base.Address;
                      if (ba.address) {
                        nameEntryAddressString =
                          '0x' +
                          ba.address
                            .map((b) => b.toString(16).padStart(2, '0'))
                            .join('');
                      }
                    } else {
                      nameEntryAddressString = String(nameEntry.address);
                    }

                    if (
                      nameEntryAddressString.toLowerCase() ===
                      addressString.toLowerCase()
                    ) {
                      const updatedName = types.Name.createFrom(nameEntry);
                      updatedName.deleted = true;
                      return updatedName;
                    }
                    return nameEntry;
                  });
                });
                // Backend will emit REFRESH, no need to call frontend Reload() here
              })
              .catch((err) => {
                Logger(`Error deleting name ${addressString}: ${err}`);
                // Optionally, revert optimistic update or show error to user
              });
          }
          break;
        case 'undelete':
          Logger(`Action ${type} clicked for ${addressString}`);
          if (addressString) {
            UndeleteName(addressString)
              .then(() => {
                // Optimistic update
                setNames((prev) => {
                  return prev.map((nameEntry) => {
                    let nameEntryAddressString = '';
                    if (
                      nameEntry.address &&
                      typeof nameEntry.address === 'object' &&
                      'address' in nameEntry.address
                    ) {
                      const ba = nameEntry.address as base.Address;
                      if (ba.address) {
                        nameEntryAddressString =
                          '0x' +
                          ba.address
                            .map((b) => b.toString(16).padStart(2, '0'))
                            .join('');
                      }
                    } else {
                      nameEntryAddressString = String(nameEntry.address);
                    }

                    if (
                      nameEntryAddressString.toLowerCase() ===
                      addressString.toLowerCase()
                    ) {
                      const updatedName = types.Name.createFrom(nameEntry);
                      updatedName.deleted = false;
                      return updatedName;
                    }
                    return nameEntry;
                  });
                });
                // Backend will emit REFRESH
              })
              .catch((err) => {
                Logger(`Error undeleting name ${addressString}: ${err}`);
              });
          }
          break;
        case 'remove':
          if (row && addressString) {
            // Show confirmation dialog
            if (
              confirm(
                `Are you sure you want to permanently remove the name entry for ${addressString}? This action cannot be undone.`,
              )
            ) {
              RemoveName(addressString)
                .then(() => {
                  Logger(`Successfully removed name: ${addressString}`);
                  // Backend will emit REFRESH
                })
                .catch((err) => {
                  Logger(`Error removing name ${addressString}: ${err}`);
                });
            }
          }
          break;
        case 'autoname':
          if (row && addressString) {
            AutonameName(addressString)
              .then(() => {
                Logger(`Successfully auto-named address: ${addressString}`);
                // Backend will emit REFRESH
              })
              .catch((err) => {
                Logger(`Error auto-naming address ${addressString}: ${err}`);
              });
          }
          break;
        case 'clean':
          CleanNames(currentTab)
            .then(() => {
              Logger(`Successfully cleaned names for tab: ${currentTab}`);
              Reload();
            })
            .catch((err) => {
              Logger(`Error cleaning names: ${err}`);
            });
          break;
        case 'publish':
          PublishNames(currentTab)
            .then(() => {
              Logger(`Successfully published names for tab: ${currentTab}`);
              Reload();
            })
            .catch((err) => {
              Logger(`Error publishing names: ${err}`);
            });
          break;
        default:
          Logger(`Unknown action: ${type}`);
      }
    },
    [listType, setNames], // Removed Reload from dependencies
  );

  const nameColumns = useMemo(
    () => [
      { key: 'name', header: 'Name', sortable: true },
      { key: 'address', header: 'Address', sortable: true, readOnly: true },
      { key: 'tags', header: 'Tags', sortable: true },
      { key: 'source', header: 'Source', sortable: true },
      {
        key: 'actions',
        header: 'Actions',
        render: (row: Record<string, unknown>) => {
          Logger('RENDER ROW: ' + JSON.stringify(row));
          const typedRow = row as unknown as types.Name;

          let addressForKey = '';
          if (typedRow.address) {
            if (
              typeof typedRow.address === 'object' &&
              'address' in typedRow.address
            ) {
              const ba = typedRow.address as base.Address;
              if (ba.address) {
                addressForKey =
                  '0x' +
                  ba.address
                    .map((b) => b.toString(16).padStart(2, '0'))
                    .join('');
              }
            } else {
              addressForKey = String(typedRow.address);
            }
          }

          Logger(
            `DEBUG_DELETED_STATUS5: nameColumns.render for row: ${addressForKey}, Name: ${typedRow.name}, Deleted: ${typedRow.deleted}, Typeof Deleted: ${typeof typedRow.deleted}`,
          );

          const isDeleted = typedRow.deleted === true;

          return (
            <div style={{ display: 'flex', alignItems: 'center' }}>
              {isDeleted && <span className="deleted-label">DELETED</span>}
              <ModifyIcons
                key={`${addressForKey}-${isDeleted ? 'deleted' : 'active'}`}
                row={typedRow}
                onAction={(type, rowData, tabName) =>
                  handleAction(type, rowData, tabName)
                }
                tabName={listType}
              />
            </div>
          );
        },
        sortable: false,
        editable: false,
      },
    ],
    [listType, handleAction],
  );

  const createTableContent = (tabLabel: string) => {
    const tagsVisible = !!showTagsView[tabLabel];
    Logger('TABLE DATA: ' + JSON.stringify(names));
    return (
      <TableProvider key={`table-provider-${tabLabel}`}>
        {tagsVisible ? (
          <div className="dual-table-layout">
            <TagsTable
              key={`tags-table-${tabLabel}`}
              tags={[]}
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
            <div className="table-wrapper" ref={mainTableRef}>
              <Table
                key={`data-table-${tabLabel}-${forceRenderKey}`}
                columns={nameColumns}
                data={names as unknown as Record<string, unknown>[]}
                loading={loading}
                error={error}
                sort={sort}
                onSortChange={setSort}
                filter={filter}
                onFilterChange={setFilter}
                tableKey={tableKey}
                onSaveRow={handleSaveRow}
                onCancelRow={handleCancelRow}
                onSubmit={handleSubmit}
                editRow={editRowData as unknown as Record<string, unknown>}
                onAddRow={() => handleAction('add')}
              />
            </div>
          </div>
        ) : (
          <Table
            key={`data-table-${tabLabel}-${forceRenderKey}`}
            columns={nameColumns}
            data={names as unknown as Record<string, unknown>[]}
            loading={loading}
            error={error}
            sort={sort}
            onSortChange={setSort}
            filter={filter}
            onFilterChange={setFilter}
            tableKey={tableKey}
            onSaveRow={handleSaveRow}
            onCancelRow={handleCancelRow}
            onSubmit={handleSubmit}
            editRow={editRowData as unknown as Record<string, unknown>}
            onAddRow={() => handleAction('add')}
          />
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
