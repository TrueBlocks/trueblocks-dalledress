import { useCallback, useEffect, useMemo, useState } from 'react';
import { Group, Stack, Tabs, Text, Title } from '@mantine/core';
import { IconDatabase } from '@tabler/icons-react';
import { DetailHeader, usePersistedTab } from '@trueblocks/ui';
import { Column, DataTable } from '../components/DataTable';
import { isEditableElement } from '../utils/keyboard';
import { uniqueSortedValues } from '../utils/table';
import {
  GetTab,
  ListDatabaseArchives,
  ListDatabaseRecords,
  SetTab,
} from '../../wailsjs/go/app/App';
import { dalle } from '../../wailsjs/go/models';
import { storage } from '../../wailsjs/go/models';

type DatabaseFileRow = {
  archive: storage.DatabaseArchiveManifest;
  file: storage.DatabaseFileManifest;
};

type DatabaseRecordRow = storage.DatabaseRecord & {
  rowIndex: number;
};

function messageFromError(error: unknown): string {
  return error instanceof Error ? error.message : String(error);
}

export function Databases() {
  const [archives, setArchives] = useState<storage.DatabaseArchiveManifest[]>([]);
  const [selectedName, setSelectedName] = useState('');
  const [records, setRecords] = useState<dalle.DatabaseRecordsResult | null>(null);
  const [filteredFiles, setFilteredFiles] = useState<DatabaseFileRow[]>([]);
  const [error, setError] = useState('');
  const { activeTab, setActiveTab } = usePersistedTab({
    key: 'databases',
    defaultTab: 'list',
    loadTab: GetTab,
    saveTab: SetTab,
    tabs: ['list', 'detail'],
    cycleViewId: 'databases',
  });

  const files = archives.flatMap((archive) => archive.files.map((file) => ({ archive, file })));
  const selectedFileRow = files.find(({ file }) => file.name === selectedName);
  const selectedFile = selectedFileRow?.file;
  const databaseList = filteredFiles.length > 0 ? filteredFiles : files;
  const selectedIndex = databaseList.findIndex(({ file }) => file.name === selectedName);
  const hasPrevious = selectedIndex > 0;
  const hasNext = selectedIndex >= 0 && selectedIndex < databaseList.length - 1;
  const recordRows: DatabaseRecordRow[] = (records?.records ?? []).map((record, rowIndex) => ({
    ...record,
    rowIndex,
  }));

  const fileNames = useMemo(() => uniqueSortedValues(files.map(({ file }) => file.name)), [files]);
  const fileVersions = useMemo(
    () => uniqueSortedValues(files.map(({ archive }) => archive.version)),
    [files],
  );
  const fileColumnOptions = useMemo(
    () => uniqueSortedValues(files.flatMap(({ file }) => file.columns)),
    [files],
  );

  const databaseValueGetter = (row: DatabaseFileRow, column: string) => {
    switch (column) {
      case 'name':
        return row.file.name;
      case 'rows':
        return row.file.rows;
      case 'columns':
        return row.file.columns.join(', ');
      case 'version':
        return row.archive.version;
      default:
        return '';
    }
  };

  const recordValueGetter = (row: DatabaseRecordRow, column: string) => {
    if (column === 'key') return row.key;
    if (column.startsWith('value-')) {
      const valueIndex = Number.parseInt(column.replace('value-', ''), 10);
      return row.values[valueIndex] ?? '';
    }
    return '';
  };

  const fileColumns: Column<DatabaseFileRow>[] = useMemo(
    () => [
      {
        key: 'name',
        label: 'Name',
        width: '30%',
        render: ({ file }) => file.name,
        sortValue: ({ file }) => file.name.toLowerCase(),
        filterOptions: fileNames,
        scrollOnSelect: true,
      },
      {
        key: 'rows',
        label: 'Rows',
        width: '12%',
        render: ({ file }) => file.rows.toLocaleString(),
        sortValue: ({ file }) => file.rows,
        filterRange: true,
      },
      {
        key: 'columns',
        label: 'Columns',
        width: '43%',
        render: ({ file }) => file.columns.join(', '),
        sortValue: ({ file }) => file.columns.length,
        filterOptions: fileColumnOptions,
        filterFn: ({ file }, selected) => file.columns.some((column) => selected.has(column)),
      },
      {
        key: 'version',
        label: 'Version',
        width: '15%',
        render: ({ archive }) => archive.version,
        sortValue: ({ archive }) => archive.version,
        filterOptions: fileVersions,
      },
    ],
    [fileColumnOptions, fileNames, fileVersions],
  );

  const recordColumns: Column<DatabaseRecordRow>[] = useMemo(() => {
    const valueColumns = (selectedFile?.columns ?? []).map((column, index) => ({
      key: `value-${index}`,
      label: column,
      render: (record: DatabaseRecordRow) => record.values[index] ?? '',
      sortValue: (record: DatabaseRecordRow) => record.values[index] ?? '',
      filterOptions: uniqueSortedValues(recordRows.map((record) => record.values[index] ?? '')),
    }));
    return [
      {
        key: 'key',
        label: 'Key',
        width: '20%',
        render: (record: DatabaseRecordRow) => record.key,
        sortValue: (record: DatabaseRecordRow) => record.key,
        filterOptions: uniqueSortedValues(recordRows.map((record) => record.key)),
        scrollOnSelect: true,
      },
      ...valueColumns,
    ];
  }, [recordRows, selectedFile]);

  const loadRecords = useCallback((name: string) => {
    setSelectedName(name);
    setError('');
    ListDatabaseRecords(name, 200)
      .then(setRecords)
      .catch((err: unknown) => setError(messageFromError(err)));
  }, []);

  const load = useCallback(
    (preferred = '') => {
      setError('');
      ListDatabaseArchives()
        .then((items) => {
          const next = items ?? [];
          setArchives(next);
          const nextFiles = next.flatMap((archive) => archive.files.map((file) => file.name));
          const nextSelected = nextFiles.includes(preferred) ? preferred : nextFiles[0] || '';
          if (nextSelected) loadRecords(nextSelected);
        })
        .catch((err: unknown) => setError(messageFromError(err)));
    },
    [loadRecords],
  );

  useEffect(() => {
    load();
  }, [load]);

  useEffect(() => {
    const handleRefresh = (event: Event) => {
      if ((event as CustomEvent).detail !== 'databases') return;
      load(selectedName);
    };

    window.addEventListener('view:refresh', handleRefresh);
    return () => window.removeEventListener('view:refresh', handleRefresh);
  }, [load, selectedName]);

  const selectDatabase = useCallback(
    (row: DatabaseFileRow) => {
      loadRecords(row.file.name);
      setActiveTab('detail');
    },
    [loadRecords, setActiveTab],
  );

  const selectDatabaseByIndex = useCallback(
    (index: number) => {
      const next = databaseList[index];
      if (next) selectDatabase(next);
    },
    [databaseList, selectDatabase],
  );

  const returnToList = useCallback(() => {
    setActiveTab('list');
  }, [setActiveTab]);

  const selectPrevious = useCallback(() => {
    if (hasPrevious) selectDatabaseByIndex(selectedIndex - 1);
  }, [hasPrevious, selectDatabaseByIndex, selectedIndex]);

  const selectNext = useCallback(() => {
    if (hasNext) selectDatabaseByIndex(selectedIndex + 1);
  }, [hasNext, selectDatabaseByIndex, selectedIndex]);

  const searchFiles = (row: DatabaseFileRow, search: string) => {
    const query = search.toLowerCase();
    return [row.file.name, row.archive.version, row.file.path, ...row.file.columns].some((value) =>
      value.toLowerCase().includes(query),
    );
  };

  const searchRecords = (record: storage.DatabaseRecord, search: string) => {
    const query = search.toLowerCase();
    return [record.key, ...record.values].some((value) => value.toLowerCase().includes(query));
  };

  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      if (activeTab !== 'detail') return;
      if ((event.metaKey || event.ctrlKey) && event.shiftKey && event.key === 'ArrowLeft') {
        event.preventDefault();
        returnToList();
        return;
      }
      if (event.metaKey || event.ctrlKey || event.altKey || isEditableElement(event.target)) return;
      if (event.key === 'ArrowLeft') {
        event.preventDefault();
        selectPrevious();
      }
      if (event.key === 'ArrowRight') {
        event.preventDefault();
        selectNext();
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [activeTab, returnToList, selectNext, selectPrevious]);

  return (
    <Stack gap="md">
      <Group justify="space-between" align="end">
        <Stack gap={2}>
          <Title order={2}>Databases</Title>
          <Text c="dimmed" size="sm">
            {files.length} tables
          </Text>
        </Stack>
      </Group>
      {error && <Text c="red">{error}</Text>}
      {archives.map((archive) => (
        <Text key={archive.version} size="sm" c="dimmed">
          {archive.version} · {archive.archiveHash}
        </Text>
      ))}
      <Tabs value={activeTab} onChange={(value) => value && setActiveTab(value)}>
        <Tabs.List>
          <Tabs.Tab value="list">List</Tabs.Tab>
          <Tabs.Tab value="detail" disabled={!selectedName}>
            Detail
          </Tabs.Tab>
        </Tabs.List>
        <Tabs.Panel value="list" pt="md">
          <DataTable<DatabaseFileRow>
            tableName="dalle-databases"
            data={files}
            columns={fileColumns}
            getRowKey={({ archive, file }) => `${archive.version}-${file.path}`}
            onRowClick={selectDatabase}
            onFilteredSortedChange={setFilteredFiles}
            searchFn={searchFiles}
            valueGetter={databaseValueGetter}
          />
        </Tabs.Panel>
        <Tabs.Panel value="detail" pt="md">
          <Stack gap="md">
            <DetailHeader
              hasPrev={hasPrevious}
              hasNext={hasNext}
              onPrev={selectPrevious}
              onNext={selectNext}
              onBack={returnToList}
              currentIndex={selectedIndex >= 0 ? selectedIndex : undefined}
              totalCount={databaseList.length}
              icon={<IconDatabase size={24} />}
              title={<Text fw={700}>{selectedName || 'Records'}</Text>}
              subtitle={
                <Text size="sm" c="dimmed">
                  {selectedFileRow?.archive.version ?? ''} ·{' '}
                  {selectedFile?.rows.toLocaleString() ?? 0} rows · showing first{' '}
                  {recordRows.length.toLocaleString()}
                </Text>
              }
            />
            <DataTable<DatabaseRecordRow>
              tableName={`dalle-database-${selectedName || 'records'}`}
              data={recordRows}
              columns={recordColumns}
              getRowKey={(record) => `${record.key}-${record.rowIndex}`}
              searchFn={searchRecords}
              valueGetter={recordValueGetter}
            />
          </Stack>
        </Tabs.Panel>
      </Tabs>
    </Stack>
  );
}
