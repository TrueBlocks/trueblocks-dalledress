import { useEffect, useState } from 'react';

import { Logger } from '@app';
import { Column, Table, TableProvider, usePagination } from '@components';
import { useAppContext } from '@contexts';
import { TabView } from '@layout';
import { msgs, sorting, types } from '@models';
import { GetNamesPage } from '@names';
import { EventsOn } from '@runtime';

export const Names = () => {
  const { lastTab } = useAppContext();

  const [sort, setSort] = useState<sorting.SortDef | null>(null);
  const [filter, setFilter] = useState('');
  const [names, setNames] = useState<types.Name[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [listType, setListType] = useState<ListType>(
    getListTypeFromLabel(lastTab['/names'] || ''),
  );

  const { pagination, goToPage, changePageSize, setTotalItems } = usePagination(
    '/names',
    listType,
  );

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
      } catch (err) {
        console.error('Error loading names:', err);
        setError(err instanceof Error ? err.message : 'Failed to load names');
        setNames([]);
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
    setTotalItems,
  ]);

  useEffect(() => {
    const currentTabLabel = lastTab['/names'] || '';
    setListType(getListTypeFromLabel(currentTabLabel));
  }, [lastTab]);

  // Each tab gets its own TableProvider instance to ensure state isolation
  const createTableContent = () => (
    <TableProvider>
      <Table<types.Name>
        columns={nameColumns}
        data={names}
        loading={loading}
        error={error}
        pagination={pagination}
        onPageChange={goToPage}
        onPageSizeChange={changePageSize}
        sort={sort}
        onSortChange={setSort}
        filter={filter}
        onFilterChange={setFilter}
      />
    </TableProvider>
  );

  const tabs = [
    { label: 'All', content: createTableContent() },
    { label: 'Custom', content: createTableContent() },
    { label: 'Prefund', content: createTableContent() },
    { label: 'Regular', content: createTableContent() },
    { label: 'Baddress', content: createTableContent() },
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
