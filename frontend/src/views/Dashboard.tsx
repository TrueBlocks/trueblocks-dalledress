import { useEffect, useState } from 'react';
import {
  Button,
  Checkbox,
  Group,
  Paper,
  Select,
  Stack,
  Text,
  Textarea,
  Title,
} from '@mantine/core';
import {
  DASHBOARD_PREFS,
  DEFAULT_IMAGE_MODEL,
  IMAGE_MODELS,
  booleanPref,
  serializeBooleanPref,
} from '../dallePrefs';
import { StatusLevel } from '../components/StatusBar';
import { playCompletionBeep } from '../App';
import {
  Generate,
  GetImage,
  GetImageModel,
  GetPref,
  ListBackstyles,
  ListSeries,
  NormalizeSeed,
  Preview,
  SetImageModel,
  SetPref,
} from '../../wailsjs/go/app/App';
import { dalle } from '../../wailsjs/go/models';

type DashboardProps = {
  onGeneratedImage: (imageId: string) => void;
  currentImage: dalle.ImageMetadataRecord | null;
  onCurrentImageChange: (record: dalle.ImageMetadataRecord | null) => void;
  onStatusChange: (status: {
    visible: boolean;
    level: StatusLevel;
    message: string;
    meta?: string;
    percent?: number;
  }) => void;
  onProgressStart: (series: string, seed: string) => void;
};

function messageFromError(error: unknown): string {
  return error instanceof Error ? error.message : String(error);
}

export function Dashboard({
  onGeneratedImage,
  currentImage,
  onCurrentImageChange,
  onStatusChange,
  onProgressStart,
}: DashboardProps) {
  const [series, setSeries] = useState<dalle.Series[]>([]);
  const [selectedSeries, setSelectedSeries] = useState<string>('');
  const [input, setInput] = useState('Person Tour Coordinates');
  const [enhance, setEnhance] = useState(false);
  const [generateImage, setGenerateImage] = useState(false);
  const [annotate, setAnnotate] = useState(false);
  const [result, setResult] = useState<dalle.GenerateResult | null>(null);
  const [working, setWorking] = useState(false);
  const [prefsLoaded, setPrefsLoaded] = useState(false);
  const [imageModel, setImageModelState] = useState(DEFAULT_IMAGE_MODEL);
  const [backstyles, setBackstyles] = useState<string[]>([]);
  const [selectedBackstyle, setSelectedBackstyle] = useState<string>('');
  const [error, setError] = useState('');

  useEffect(() => {
    Promise.all([
      ListSeries(false, false),
      ListBackstyles(200),
      GetPref(DASHBOARD_PREFS.input),
      GetPref(DASHBOARD_PREFS.series),
      GetPref(DASHBOARD_PREFS.backstyle),
      GetPref(DASHBOARD_PREFS.enhance),
      GetPref(DASHBOARD_PREFS.generateImage),
      GetPref(DASHBOARD_PREFS.annotate),
      GetPref(DASHBOARD_PREFS.imageModel),
      GetImageModel(),
    ])
      .then(
        ([
          items,
          loadedBackstyles,
          savedInput,
          savedSeries,
          savedBackstyle,
          savedEnhance,
          savedGenerateImage,
          savedAnnotate,
          savedModel,
          engineModel,
        ]) => {
          const next = items ?? [];
          setSeries(next);
          setBackstyles(loadedBackstyles);
          setInput(savedInput || 'Person Tour Coordinates');
          setSelectedBackstyle(
            loadedBackstyles.includes(savedBackstyle) ? savedBackstyle : loadedBackstyles[0] || '',
          );
          setEnhance(booleanPref(savedEnhance));
          setGenerateImage(booleanPref(savedGenerateImage));
          setAnnotate(booleanPref(savedAnnotate));
          setSelectedSeries(
            next.some((item) => item.suffix === savedSeries) ? savedSeries : next[0]?.suffix || '',
          );
          const model = savedModel || engineModel || DEFAULT_IMAGE_MODEL;
          setImageModelState(model);
          SetImageModel(model);
          setPrefsLoaded(true);
        },
      )
      .catch((err: unknown) => setError(messageFromError(err)));
  }, []);

  useEffect(() => {
    if (!prefsLoaded) return;
    SetPref(DASHBOARD_PREFS.input, input);
    NormalizeSeed(input, '')
      .then((hash) => GetImage(hash))
      .then((record) => onCurrentImageChange(record))
      .catch(() => onCurrentImageChange(null));
  }, [input, onCurrentImageChange, prefsLoaded]);

  useEffect(() => {
    if (!prefsLoaded) return;
    SetPref(DASHBOARD_PREFS.series, selectedSeries);
  }, [prefsLoaded, selectedSeries]);

  useEffect(() => {
    if (!prefsLoaded) return;
    SetPref(DASHBOARD_PREFS.backstyle, selectedBackstyle);
  }, [prefsLoaded, selectedBackstyle]);

  useEffect(() => {
    if (!prefsLoaded) return;
    SetPref(DASHBOARD_PREFS.enhance, serializeBooleanPref(enhance));
  }, [enhance, prefsLoaded]);

  useEffect(() => {
    if (!prefsLoaded) return;
    SetPref(DASHBOARD_PREFS.generateImage, serializeBooleanPref(generateImage));
  }, [generateImage, prefsLoaded]);

  useEffect(() => {
    if (!prefsLoaded) return;
    SetPref(DASHBOARD_PREFS.annotate, serializeBooleanPref(annotate));
  }, [annotate, prefsLoaded]);

  useEffect(() => {
    if (!prefsLoaded || !currentImage) return;
    const meta = currentImage.metadata;
    if (meta.input) setInput(meta.input);
    if (meta.series?.name) setSelectedSeries(meta.series.name);
  }, [currentImage, prefsLoaded]);

  const handleModelChange = (value: string | null) => {
    if (!value) return;
    setImageModelState(value);
    SetPref(DASHBOARD_PREFS.imageModel, value);
    SetImageModel(value);
  };

  const loadFromCurrentImage = () => {
    if (!currentImage) return;
    const meta = currentImage.metadata;
    if (meta.input) setInput(meta.input);
    if (meta.series?.name) setSelectedSeries(meta.series.name);
  };

  const request = () =>
    dalle.GenerateRequest.createFrom({
      input,
      series: selectedSeries || undefined,
      backstyle: selectedBackstyle || undefined,
      enhance,
      image: generateImage,
      annotate: generateImage && annotate,
    });

  const runPreview = () => {
    setError('');
    onStatusChange({ visible: true, level: 'progress', message: 'Building preview prompt' });
    setWorking(true);
    Preview(request())
      .then((next) => {
        setResult(next);
        onStatusChange({ visible: true, level: 'success', message: 'Preview ready' });
      })
      .catch((err: unknown) => {
        const message = messageFromError(err);
        setError(message);
        onStatusChange({ visible: true, level: 'error', message });
      })
      .finally(() => setWorking(false));
  };

  const runGenerate = () => {
    const nextRequest = request();
    setError('');
    onStatusChange({ visible: true, level: 'progress', message: 'Preparing generation' });
    setWorking(true);
    Preview(nextRequest)
      .then((preview) => {
        setResult(preview);
        onProgressStart(preview.series, preview.seed);
        onStatusChange({
          visible: true,
          level: 'progress',
          message: generateImage ? 'Starting image generation' : 'Writing prompt metadata',
          meta: preview.series,
        });
        return Generate(nextRequest);
      })
      .then((next) => {
        setResult(next);
        const imageId = next.metadata?.imageId || next.seed;
        playCompletionBeep();
        onStatusChange({
          visible: true,
          level: 'success',
          message: generateImage ? 'Image generation complete' : 'Generation metadata complete',
          meta: next.metadata?.prompts?.titlePrompt || next.seed,
        });
        setWorking(false);
        if (generateImage && imageId) onGeneratedImage(imageId);
      })
      .catch((err: unknown) => {
        const message = messageFromError(err);
        setError(message);
        onStatusChange({ visible: true, level: 'error', message });
        setWorking(false);
      });
  };

  return (
    <Stack>
      <Title order={2}>Dashboard</Title>
      <Textarea
        label="Seed"
        value={input}
        minRows={3}
        onChange={(event) => setInput(event.currentTarget.value)}
      />
      <Group align="end">
        <Select
          label="Series"
          value={selectedSeries}
          data={series.map((item) => ({ value: item.suffix, label: item.suffix }))}
          onChange={(value) => setSelectedSeries(value ?? '')}
        />
        <Select
          label="Model"
          value={imageModel}
          data={IMAGE_MODELS.map((m) => ({ value: m.value, label: m.label }))}
          onChange={handleModelChange}
          allowDeselect={false}
        />
        <Select
          label="Background"
          value={selectedBackstyle}
          data={backstyles.map((style) => ({ value: style, label: style }))}
          onChange={(value) => setSelectedBackstyle(value ?? '')}
          disabled={backstyles.length === 0}
          allowDeselect={false}
        />
        <Checkbox
          label="Enhance prompt"
          checked={enhance}
          onChange={(event) => setEnhance(event.currentTarget.checked)}
        />
        <Checkbox
          label="Generate image"
          checked={generateImage}
          onChange={(event) => setGenerateImage(event.currentTarget.checked)}
        />
        <Checkbox
          label="Annotate"
          checked={annotate}
          disabled={!generateImage}
          onChange={(event) => setAnnotate(event.currentTarget.checked)}
        />
        <Button onClick={runPreview} loading={working}>
          Preview
        </Button>
        <Button onClick={runGenerate} loading={working}>
          Generate
        </Button>
        {currentImage && (
          <Button variant="light" onClick={loadFromCurrentImage}>
            Load from current image
          </Button>
        )}
      </Group>
      {error && <Text c="red">{error}</Text>}
      {result && (
        <Paper withBorder p="md">
          <Stack gap="xs">
            <Text fw={700}>{result.metadata?.prompts?.titlePrompt || result.seed}</Text>
            <Text size="sm" c="dimmed">
              {result.series} · {result.recipe} · {result.seed}
            </Text>
            <Text>{result.metadata?.prompts?.prompt}</Text>
            {result.generatedPath && <Text size="sm">Generated: {result.generatedPath}</Text>}
            {result.annotatedPath && <Text size="sm">Annotated: {result.annotatedPath}</Text>}
          </Stack>
        </Paper>
      )}
    </Stack>
  );
}
