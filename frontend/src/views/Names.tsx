import { useEffect, useMemo, useRef, useState } from 'react';

import {
  Table,
  TableProvider,
  TagsTable,
  TagsTableHandle,
  usePagination,
} from '@components';
import { TableKey, useAppContext } from '@contexts';
import { TabView } from '@layout';
import { msgs, sorting, types } from '@models';
import {
  ClearSelectedTag,
  GetNamesPage,
  GetSelectedTag,
  SetSelectedTag,
} from '@names';
import { EventsOn } from '@runtime';

import './Names.css';
import { Logger } from '@app';

export const FocusSider = 'focus-tags-table';

export const Names = () => {
  const { lastTab } = useAppContext();

  const [sort, setSort] = useState<sorting.SortDef | null>(null);
  const [filter, setFilter] = useState('');
  const [names, setNames] = useState<types.Name[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
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
        const tag = await GetSelectedTag(listType);
        setSelectedTag(tag || null);
      } catch (err) {
        console.error('Error fetching selected tag:', err);
        setSelectedTag(null);
      }
    };
    fetchSelectedTag();
  }, [listType]);

  useEffect(() => {
    const loadNames = async () => {
      setLoading(true);
      setError(null);

      try {
        const result = await GetNamesPage(
          listType,
          pagination.currentPage * pagination.pageSize,
          pagination.pageSize,
          sort as sorting.SortDef,
          filter ?? '',
        );
        setNames(result.names || []);
        setTotalItems(result.total || 0);
        setTags(result.tags || []);
      } catch (err) {
        console.error('Error loading names:', err);
        setError(err instanceof Error ? err.message : 'Failed to load names');
        setNames([]);
        setTags([]);
      } finally {
        setLoading(false);
      }
    };
    loadNames();

    var unsubscribe = EventsOn(msgs.EventType.REFRESH, () => {
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
  };

  const handleSubmit = (data: Record<string, unknown>) => {
    Logger('Front end got returned from SaveName');
  };

  // Each tab gets its own TableProvider instance to ensure state isolation
  const createTableContent = (tabLabel: string) => {
    const tagsVisible = !!showTagsView[tabLabel];
    return (
      <TableProvider key={`table-provider-${tabLabel}`}>
        {tagsVisible ? (
          <div className="dual-table-layout">
            <TagsTable
              key={`tags-table-${tabLabel}`}
              tags={tags}
              onTagSelect={handleTagSelect}
              selectedTag={selectedTag}
              visible={true}
              ref={(ref) => {
                // Store the ref by tab name for targeted focusing
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
                key={`data-table-${tabLabel}`}
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
                onSubmit={handleSubmit}
              />
            </div>
          </div>
        ) : (
          <Table
            key={`data-table-${tabLabel}`}
            columns={nameColumns}
            data={names as unknown as Record<string, unknown>[]}
            loading={loading}
            error={error}
            sort={sort}
            onSortChange={setSort}
            filter={filter}
            onFilterChange={setFilter}
            tableKey={tableKey}
            onSubmit={handleSubmit}
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

const nameColumns = [
  { key: 'name', header: 'Name', sortable: true },
  { key: 'address', header: 'Address', sortable: true, readOnly: true },
  { key: 'tags', header: 'Tags', sortable: true },
  { key: 'source', header: 'Source', sortable: true },
  {
    key: 'actions',
    header: 'Actions',
    render: (row: Record<string, unknown>) =>
      `${(row as { deleted?: boolean }).deleted ? 'Deleted ' : ''}`,
    sortable: false,
    editable: false,
  },
];

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
