import { useCallback, useEffect, useState } from 'react';

import { GetProjectViewState, SetProjectViewState } from '@app';
import { Group, NumberInput, Select } from '@mantine/core';
import { project } from '@models';
import { LogError } from '@utils';

export type GallerySortMode = 'series' | 'address';

export interface GalleryControlsState {
  sortMode: GallerySortMode;
  columns: number;
}

interface GalleryControlsProps {
  viewStateKey: project.ViewStateKey;
  value: GalleryControlsState;
  onChange: (state: GalleryControlsState) => void;
}

async function loadPersisted(
  viewStateKey: project.ViewStateKey,
): Promise<Partial<GalleryControlsState>> {
  try {
    const viewStates = await GetProjectViewState(viewStateKey.viewName);
    const facetState = viewStates[viewStateKey.facetName];
    const other = facetState?.other || {};
    const gallery = (other as Record<string, unknown>).gallery as
      | Record<string, unknown>
      | undefined;
    return {
      sortMode: gallery?.sortMode as GallerySortMode | undefined,
      columns:
        typeof gallery?.columns === 'number'
          ? (gallery.columns as number)
          : undefined,
    };
  } catch (e) {
    LogError('gallery:loadState:' + String(e));
    return {};
  }
}

async function persist(
  viewStateKey: project.ViewStateKey,
  state: GalleryControlsState,
) {
  try {
    const viewStates = await GetProjectViewState(viewStateKey.viewName);
    const facetName = viewStateKey.facetName;
    const existing = viewStates[facetName] || {
      sorting: {},
      filtering: {},
      other: {},
    };
    const updated = {
      ...existing,
      other: {
        ...(existing.other || {}),
        gallery: { sortMode: state.sortMode, columns: state.columns },
      },
    };
    await SetProjectViewState(viewStateKey.viewName, {
      ...viewStates,
      [facetName]: updated,
    });
  } catch (e) {
    LogError('gallery:persistState:' + String(e));
  }
}

export const GalleryControls = ({
  viewStateKey,
  value,
  onChange,
}: GalleryControlsProps) => {
  const [hydrated, setHydrated] = useState(false);

  useEffect(() => {
    if (hydrated) return;
    loadPersisted(viewStateKey).then((loaded) => {
      const next: GalleryControlsState = {
        sortMode: loaded.sortMode || value.sortMode,
        columns: loaded.columns || value.columns,
      };
      onChange(next);
      setHydrated(true);
    });
  }, [hydrated, viewStateKey, value.columns, value.sortMode, onChange]);

  const updateSort = useCallback(
    (mode: GallerySortMode) => {
      const next = { ...value, sortMode: mode };
      onChange(next);
      persist(viewStateKey, next);
    },
    [value, onChange, viewStateKey],
  );

  const updateColumns = useCallback(
    (cols: number) => {
      const capped = Math.min(12, Math.max(1, cols));
      const next = { ...value, columns: capped };
      onChange(next);
      persist(viewStateKey, next);
    },
    [value, onChange, viewStateKey],
  );

  return (
    <Group gap="md" mb="sm" wrap="nowrap" align="flex-end">
      <Select
        label="Sort"
        placeholder="Sort"
        value={value.sortMode}
        onChange={(v) => v && updateSort(v as GallerySortMode)}
        data={[
          { value: 'series', label: 'Series' },
          { value: 'address', label: 'Address' },
        ]}
        w={160}
        size="xs"
      />
      <NumberInput
        label="Columns"
        size="xs"
        value={value.columns}
        min={1}
        max={12}
        onChange={(val) => typeof val === 'number' && updateColumns(val)}
        w={120}
      />
    </Group>
  );
};
