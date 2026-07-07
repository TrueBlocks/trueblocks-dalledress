import { useCallback, useEffect, useMemo, useState } from 'react';
import {
  Badge,
  Button,
  Checkbox,
  Group,
  Paper,
  Select,
  SimpleGrid,
  Stack,
  Tabs,
  Text,
  Textarea,
  TextInput,
  Title,
} from '@mantine/core';
import { IconStack2 } from '@tabler/icons-react';
import { DetailHeader, usePersistedTab } from '@trueblocks/ui';
import { Column, DataTable } from '../components/DataTable';
import { isEditableElement } from '../utils/keyboard';
import { uniqueSortedValues } from '../utils/table';
import { GetTab, ListSeries, SaveSeries, SetSeriesHidden, SetTab } from '../../wailsjs/go/app/App';
import { dalle } from '../../wailsjs/go/models';

const FILTER_FIELDS = [
  'adverbs',
  'adjectives',
  'nouns',
  'emotions',
  'occupations',
  'actions',
  'artstyles',
  'litstyles',
  'colors',
  'viewpoints',
  'gazes',
  'backstyles',
  'compositions',
] as const;

type FilterField = (typeof FILTER_FIELDS)[number];
type SeriesDraft = Record<FilterField, string> & {
  suffix: string;
  purpose: string;
  last: string;
  colorLimit: string;
};

const emptyDraft: SeriesDraft = {
  suffix: '',
  purpose: '',
  last: '0',
  colorLimit: '',
  adverbs: '',
  adjectives: '',
  nouns: '',
  emotions: '',
  occupations: '',
  actions: '',
  artstyles: '',
  litstyles: '',
  colors: '',
  viewpoints: '',
  gazes: '',
  backstyles: '',
  compositions: '',
};

function messageFromError(error: unknown): string {
  return error instanceof Error ? error.message : String(error);
}

function valuesToText(values?: string[]): string {
  return (values ?? []).join('\n');
}

function textToValues(value: string): string[] {
  return value
    .split(/[\n,]/)
    .map((item) => item.trim())
    .filter(Boolean);
}

function draftFromSeries(series: dalle.Series): SeriesDraft {
  return {
    suffix: series.suffix ?? '',
    purpose: series.purpose ?? '',
    last: String(series.last ?? 0),
    colorLimit: series.colorLimit ?? '',
    adverbs: valuesToText(series.adverbs),
    adjectives: valuesToText(series.adjectives),
    nouns: valuesToText(series.nouns),
    emotions: valuesToText(series.emotions),
    occupations: valuesToText(series.occupations),
    actions: valuesToText(series.actions),
    artstyles: valuesToText(series.artstyles),
    litstyles: valuesToText(series.litstyles),
    colors: valuesToText(series.colors),
    viewpoints: valuesToText(series.viewpoints),
    gazes: valuesToText(series.gazes),
    backstyles: valuesToText(series.backstyles),
    compositions: valuesToText(series.compositions),
  };
}

function seriesFromDraft(draft: SeriesDraft): dalle.Series {
  return dalle.Series.createFrom({
    suffix: draft.suffix,
    purpose: draft.purpose,
    last: Number.parseInt(draft.last || '0', 10) || 0,
    colorLimit: draft.colorLimit || undefined,
    adverbs: textToValues(draft.adverbs),
    adjectives: textToValues(draft.adjectives),
    nouns: textToValues(draft.nouns),
    emotions: textToValues(draft.emotions),
    occupations: textToValues(draft.occupations),
    actions: textToValues(draft.actions),
    artstyles: textToValues(draft.artstyles),
    litstyles: textToValues(draft.litstyles),
    colors: textToValues(draft.colors),
    viewpoints: textToValues(draft.viewpoints),
    gazes: textToValues(draft.gazes),
    backstyles: textToValues(draft.backstyles),
    compositions: textToValues(draft.compositions),
  });
}

export function Series() {
  const [includeHidden, setIncludeHidden] = useState(false);
  const [items, setItems] = useState<dalle.Series[]>([]);
  const [filteredItems, setFilteredItems] = useState<dalle.Series[]>([]);
  const [selectedSuffix, setSelectedSuffix] = useState('');
  const [draft, setDraft] = useState<SeriesDraft>(emptyDraft);
  const [error, setError] = useState('');
  const [message, setMessage] = useState('');
  const [saving, setSaving] = useState(false);
  const { activeTab, setActiveTab } = usePersistedTab({
    key: 'series',
    defaultTab: 'list',
    loadTab: GetTab,
    saveTab: SetTab,
    tabs: ['list', 'detail'],
    cycleViewId: 'series',
  });

  const selected = items.find((item) => item.suffix === selectedSuffix);
  const seriesList = filteredItems.length > 0 ? filteredItems : items;
  const selectedIndex = seriesList.findIndex((item) => item.suffix === selectedSuffix);
  const hasPrevious = selectedIndex > 0;
  const hasNext = selectedIndex >= 0 && selectedIndex < seriesList.length - 1;
  const suffixOptions = useMemo(
    () => uniqueSortedValues(items.map((item) => item.suffix)),
    [items],
  );
  const purposeOptions = useMemo(
    () => uniqueSortedValues(items.map((item) => item.purpose)),
    [items],
  );
  const modifiedOptions = useMemo(
    () => uniqueSortedValues(items.map((item) => item.modifiedAt)),
    [items],
  );

  const seriesValueGetter = (series: dalle.Series, column: string) => {
    switch (column) {
      case 'suffix':
        return series.suffix;
      case 'purpose':
        return series.purpose || '';
      case 'last':
        return series.last ?? 0;
      case 'state':
        return series.deleted ? 'hidden' : 'active';
      case 'modifiedAt':
        return series.modifiedAt || '';
      default:
        return '';
    }
  };

  const columns: Column<dalle.Series>[] = useMemo(
    () => [
      {
        key: 'suffix',
        label: 'Suffix',
        width: '30%',
        render: (item) => item.suffix,
        sortValue: (item) => item.suffix.toLowerCase(),
        filterOptions: suffixOptions,
        scrollOnSelect: true,
      },
      {
        key: 'purpose',
        label: 'Purpose',
        width: '35%',
        render: (item) => item.purpose || '',
        sortValue: (item) => item.purpose || '',
        filterOptions: purposeOptions,
      },
      {
        key: 'last',
        label: 'Last',
        width: '10%',
        render: (item) => item.last ?? 0,
        sortValue: (item) => item.last ?? 0,
        filterRange: true,
      },
      {
        key: 'state',
        label: 'State',
        width: '15%',
        render: (item) => (
          <Badge variant="light" color={item.deleted ? 'gray' : 'green'}>
            {item.deleted ? 'hidden' : 'active'}
          </Badge>
        ),
        sortValue: (item) => (item.deleted ? 'hidden' : 'active'),
        filterOptions: ['active', 'hidden'],
      },
      {
        key: 'modifiedAt',
        label: 'Modified',
        width: '10%',
        render: (item) => item.modifiedAt || '',
        sortValue: (item) => item.modifiedAt || '',
        filterOptions: modifiedOptions,
      },
    ],
    [modifiedOptions, purposeOptions, suffixOptions],
  );

  const load = useCallback(
    (preferred = '') => {
      setError('');
      return ListSeries(includeHidden, false)
        .then((result) => {
          const next = result ?? [];
          setItems(next);
          const nextSelected = next.some((item) => item.suffix === preferred)
            ? preferred
            : next[0]?.suffix || '';
          setSelectedSuffix(nextSelected);
          const nextSeries = next.find((item) => item.suffix === nextSelected);
          setDraft(nextSeries ? draftFromSeries(nextSeries) : emptyDraft);
        })
        .catch((err: unknown) => setError(messageFromError(err)));
    },
    [includeHidden],
  );

  useEffect(() => {
    load();
  }, [load]);

  useEffect(() => {
    const handleRefresh = (event: Event) => {
      if ((event as CustomEvent).detail !== 'series') return;
      load(selectedSuffix);
    };

    window.addEventListener('view:refresh', handleRefresh);
    return () => window.removeEventListener('view:refresh', handleRefresh);
  }, [load, selectedSuffix]);

  const selectSeries = useCallback(
    (series: dalle.Series) => {
      setSelectedSuffix(series.suffix);
      setDraft(draftFromSeries(series));
      setMessage('');
      setError('');
      setActiveTab('detail');
    },
    [setActiveTab],
  );

  const selectSeriesByIndex = useCallback(
    (index: number) => {
      const next = seriesList[index];
      if (next) selectSeries(next);
    },
    [selectSeries, seriesList],
  );

  const returnToList = useCallback(() => {
    setActiveTab('list');
  }, [setActiveTab]);

  const selectPrevious = useCallback(() => {
    if (hasPrevious) selectSeriesByIndex(selectedIndex - 1);
  }, [hasPrevious, selectSeriesByIndex, selectedIndex]);

  const selectNext = useCallback(() => {
    if (hasNext) selectSeriesByIndex(selectedIndex + 1);
  }, [hasNext, selectSeriesByIndex, selectedIndex]);

  const newSeries = () => {
    setSelectedSuffix('');
    setDraft(emptyDraft);
    setMessage('');
    setError('');
    setActiveTab('detail');
  };

  const save = () => {
    setSaving(true);
    setError('');
    setMessage('');
    SaveSeries(seriesFromDraft(draft))
      .then((saved) => {
        setMessage(`Saved ${saved.suffix}.`);
        return load(saved.suffix);
      })
      .catch((err: unknown) => setError(messageFromError(err)))
      .finally(() => setSaving(false));
  };

  const setHidden = (hidden: boolean) => {
    if (!selectedSuffix) return;
    setSaving(true);
    setError('');
    setMessage('');
    SetSeriesHidden(selectedSuffix, hidden)
      .then((updated) => {
        setMessage(hidden ? `Hid ${updated.suffix}.` : `Restored ${updated.suffix}.`);
        return load(updated.suffix);
      })
      .catch((err: unknown) => setError(messageFromError(err)))
      .finally(() => setSaving(false));
  };

  const updateDraft = (key: keyof SeriesDraft, value: string) => {
    setDraft((current) => ({ ...current, [key]: value }));
  };

  const searchSeries = (series: dalle.Series, search: string) => {
    const query = search.toLowerCase();
    return [
      series.suffix,
      series.purpose || '',
      ...(series.adverbs ?? []),
      ...(series.adjectives ?? []),
      ...(series.nouns ?? []),
      ...(series.emotions ?? []),
      ...(series.occupations ?? []),
      ...(series.actions ?? []),
      ...(series.artstyles ?? []),
      ...(series.litstyles ?? []),
      ...(series.colors ?? []),
      ...(series.viewpoints ?? []),
      ...(series.gazes ?? []),
      ...(series.backstyles ?? []),
      ...(series.compositions ?? []),
    ].some((value) => value.toLowerCase().includes(query));
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
          <Title order={2}>Series</Title>
          <Text c="dimmed" size="sm">
            {items.length} definitions
          </Text>
        </Stack>
        <Group>
          <Checkbox
            label="Include hidden"
            checked={includeHidden}
            onChange={(event) => setIncludeHidden(event.currentTarget.checked)}
          />
          <Button onClick={newSeries}>New Series</Button>
        </Group>
      </Group>
      {error && <Text c="red">{error}</Text>}
      {message && <Text c="dimmed">{message}</Text>}

      <Tabs value={activeTab} onChange={(value) => value && setActiveTab(value)}>
        <Tabs.List>
          <Tabs.Tab value="list">List</Tabs.Tab>
          <Tabs.Tab value="detail">Detail</Tabs.Tab>
        </Tabs.List>
        <Tabs.Panel value="list" pt="md">
          <DataTable<dalle.Series>
            tableName="dalle-series"
            data={items}
            columns={columns}
            getRowKey={(item) => item.suffix}
            onRowClick={selectSeries}
            onFilteredSortedChange={setFilteredItems}
            searchFn={searchSeries}
            valueGetter={seriesValueGetter}
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
              totalCount={seriesList.length}
              icon={<IconStack2 size={24} />}
              title={<Text fw={700}>{selected?.suffix || 'New Series'}</Text>}
              subtitle={
                <Text size="sm" c="dimmed">
                  {selected?.purpose || draft.purpose || 'Series definition'}
                </Text>
              }
            />
            <Paper withBorder p="md">
              <Stack gap="sm">
                <Group grow>
                  <TextInput
                    label="Suffix"
                    value={draft.suffix}
                    onChange={(event) => updateDraft('suffix', event.currentTarget.value)}
                  />
                  <TextInput
                    label="Last index"
                    type="number"
                    value={draft.last}
                    onChange={(event) => updateDraft('last', event.currentTarget.value)}
                  />
                </Group>
                <TextInput
                  label="Purpose"
                  value={draft.purpose}
                  onChange={(event) => updateDraft('purpose', event.currentTarget.value)}
                />
                <Select
                  label="Color limit"
                  description="Empty = colors are palette suggestions. N-tone = strict color constraint."
                  value={draft.colorLimit || null}
                  data={[
                    { value: 'two-tone', label: 'Two-tone' },
                    { value: 'three-tone', label: 'Three-tone' },
                    { value: 'four-tone', label: 'Four-tone' },
                    { value: 'five-tone', label: 'Five-tone' },
                    { value: 'six-tone', label: 'Six-tone' },
                    { value: 'seven-tone', label: 'Seven-tone' },
                  ]}
                  onChange={(value) => updateDraft('colorLimit', value ?? '')}
                  clearable
                  placeholder="No limit (full palette)"
                />
                <SimpleGrid cols={{ base: 1, md: 2 }} spacing="sm">
                  {FILTER_FIELDS.map((field) => (
                    <Textarea
                      key={field}
                      label={field}
                      value={draft[field]}
                      minRows={3}
                      autosize
                      onChange={(event) => updateDraft(field, event.currentTarget.value)}
                    />
                  ))}
                </SimpleGrid>
                <Group justify="space-between">
                  <Group>
                    <Button onClick={save} loading={saving} disabled={!draft.suffix.trim()}>
                      Save
                    </Button>
                    {selected && !selected.deleted && (
                      <Button variant="light" color="gray" onClick={() => setHidden(true)}>
                        Hide
                      </Button>
                    )}
                    {selected?.deleted && (
                      <Button variant="light" onClick={() => setHidden(false)}>
                        Restore
                      </Button>
                    )}
                  </Group>
                  {selected?.modifiedAt && (
                    <Text size="xs" c="dimmed">
                      Modified {selected.modifiedAt}
                    </Text>
                  )}
                </Group>
              </Stack>
            </Paper>
          </Stack>
        </Tabs.Panel>
      </Tabs>
    </Stack>
  );
}
