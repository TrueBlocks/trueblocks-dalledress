import { useEffect, useMemo, useRef, useState } from 'react';

import { Logger } from '@app';
import {
  Column,
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

  // Reference to the tags table for focusing
  const tagsTableRef = useRef<TagsTableHandle>(null);
  const mainTableRef = useRef<HTMLDivElement>(null);

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
      Logger('Refreshing names in the frontend');
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

  // Setup event listener for focusing the tags table
  useEffect(() => {
    const focusTagsTable = () => {
      if (tagsTableRef.current) {
        tagsTableRef.current.focus();
      }
    };

    const unsubscribe = EventsOn(FocusSider, focusTagsTable);
    return unsubscribe;
  }, []);

  // Handle tag selection
  const handleTagSelect = async (
    tag: string | null,
    focusMainTable: boolean = false,
  ) => {
    try {
      Logger(`handleTagSelect called with tag: ${tag}, listType: ${listType}`);
      if (tag) {
        Logger(`Calling SetSelectedTag(${listType}, ${tag})`);
        try {
          const result = await SetSelectedTag(listType, tag);
          Logger(`SetSelectedTag result: ${JSON.stringify(result)}`);
        } catch (e) {
          Logger(`Error in SetSelectedTag: ${e}`);
        }
      } else {
        Logger(`Calling ClearSelectedTag(${listType})`);
        try {
          const result = await ClearSelectedTag(listType);
          Logger(`ClearSelectedTag result: ${JSON.stringify(result)}`);
        } catch (e) {
          Logger(`Error in ClearSelectedTag: ${e}`);
        }
      }
      setSelectedTag(tag);
    } catch (err) {
      Logger(`Error setting selected tag: ${err}`);
    }

    // Reset to first page when selecting a tag
    if (pagination.currentPage !== 0) {
      Logger(`Resetting to first page`);
      goToPage(0);
    }

    // Focus the main table only if explicitly requested (via Enter key)
    if (focusMainTable && mainTableRef.current) {
      Logger(`Focusing main table`);
      setTimeout(() => {
        const tableElement = mainTableRef.current?.querySelector('.data-table');
        if (tableElement) {
          (tableElement as HTMLElement).focus();
        }
      }, 0);
    }
  };

  // Each tab gets its own TableProvider instance to ensure state isolation
  const createTableContent = (showTagsTable: boolean = false) => (
    <TableProvider>
      {showTagsTable ? (
        <div className="dual-table-layout">
          <TagsTable
            tags={tags}
            onTagSelect={handleTagSelect}
            selectedTag={selectedTag}
            visible={true}
            ref={tagsTableRef}
          />
          <div className="table-wrapper" ref={mainTableRef}>
            <Table<types.Name>
              columns={nameColumns}
              data={names}
              loading={loading}
              error={error}
              sort={sort}
              onSortChange={setSort}
              filter={filter}
              onFilterChange={setFilter}
              tableKey={tableKey}
            />
          </div>
        </div>
      ) : (
        <Table<types.Name>
          columns={nameColumns}
          data={names}
          loading={loading}
          error={error}
          sort={sort}
          onSortChange={setSort}
          filter={filter}
          onFilterChange={setFilter}
          tableKey={tableKey}
        />
      )}
    </TableProvider>
  );

  const tabs = [
    { label: 'All', content: createTableContent(false) },
    { label: 'Custom', content: createTableContent(true) },
    { label: 'Prefund', content: createTableContent(false) },
    { label: 'Regular', content: createTableContent(false) },
    { label: 'Baddress', content: createTableContent(false) },
  ];

  return (
    <div className="table-tab-view">
      <TabView tabs={tabs} route="/names" />
    </div>
  );
};

const nameColumns: Column<types.Name>[] = [
  { key: 'name', header: 'Name', sortable: true },
  { key: 'address', header: 'Address', sortable: true },
  { key: 'tags', header: 'Tags', sortable: true },
  { key: 'source', header: 'Source', sortable: true },
  {
    key: 'actions',
    header: 'Actions',
    render: (row: types.Name) => `${row.deleted ? 'Deleted ' : ''}`,
    sortable: false,
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
