import { useCallback, useEffect, useState } from 'react';
import {
  Badge,
  Button,
  Checkbox,
  Group,
  Paper,
  SimpleGrid,
  Stack,
  Table,
  Text,
  Textarea,
  TextInput,
  Title,
} from '@mantine/core';
import { ListSeries, SaveSeries, SetSeriesHidden } from '../../wailsjs/go/app/App';
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
};

const emptyDraft: SeriesDraft = {
  suffix: '',
  purpose: '',
  last: '0',
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
  const [selectedSuffix, setSelectedSuffix] = useState('');
  const [draft, setDraft] = useState<SeriesDraft>(emptyDraft);
  const [error, setError] = useState('');
  const [message, setMessage] = useState('');
  const [saving, setSaving] = useState(false);

  const selected = items.find((item) => item.suffix === selectedSuffix);

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

  const selectSeries = (series: dalle.Series) => {
    setSelectedSuffix(series.suffix);
    setDraft(draftFromSeries(series));
    setMessage('');
    setError('');
  };

  const newSeries = () => {
    setSelectedSuffix('');
    setDraft(emptyDraft);
    setMessage('');
    setError('');
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

      <SimpleGrid cols={{ base: 1, lg: 2 }} spacing="md">
        <Paper withBorder p="md">
          <Table>
            <Table.Thead>
              <Table.Tr>
                <Table.Th>Suffix</Table.Th>
                <Table.Th>Purpose</Table.Th>
                <Table.Th>Last</Table.Th>
                <Table.Th>State</Table.Th>
              </Table.Tr>
            </Table.Thead>
            <Table.Tbody>
              {items.map((item) => (
                <Table.Tr
                  key={item.suffix}
                  onClick={() => selectSeries(item)}
                  style={{
                    cursor: 'pointer',
                    background:
                      item.suffix === selectedSuffix
                        ? 'var(--mantine-color-blue-light)'
                        : undefined,
                  }}
                >
                  <Table.Td>{item.suffix}</Table.Td>
                  <Table.Td>{item.purpose}</Table.Td>
                  <Table.Td>{item.last ?? 0}</Table.Td>
                  <Table.Td>
                    <Badge variant="light" color={item.deleted ? 'gray' : 'green'}>
                      {item.deleted ? 'hidden' : 'active'}
                    </Badge>
                  </Table.Td>
                </Table.Tr>
              ))}
            </Table.Tbody>
          </Table>
        </Paper>

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
      </SimpleGrid>
    </Stack>
  );
}
