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
import { StatusBar, StatusLevel } from '../components/StatusBar';
import {
  Generate,
  GetGenerationProgress,
  GetImageModel,
  GetPref,
  ListSeries,
  Preview,
  SetImageModel,
  SetPref,
} from '../../wailsjs/go/app/App';
import { app } from '../../wailsjs/go/models';
import { dalle } from '../../wailsjs/go/models';

type DashboardProps = {
  onGeneratedImage: (imageId: string) => void;
  currentImage: dalle.ImageMetadataRecord | null;
};

function messageFromError(error: unknown): string {
  return error instanceof Error ? error.message : String(error);
}

type ProgressTarget = {
  series: string;
  seed: string;
};

type StatusState = {
  visible: boolean;
  level: StatusLevel;
  message: string;
  meta?: string;
  percent?: number;
};

const PHASE_LABELS: Record<string, string> = {
  setup: 'Preparing generation run',
  base_prompts: 'Selecting records and building prompts',
  enhance_prompt: 'Enhancing prompt',
  image_prep: 'Preparing image request',
  image_wait: 'Waiting for image provider',
  image_download: 'Receiving image artifact',
  annotate: 'Annotating generated image',
  failed: 'Generation failed',
  completed: 'Generation complete',
};

function statusForProgress(progress: app.GenerationProgress): StatusState {
  if (progress.error) {
    return { visible: true, level: 'error', message: progress.error };
  }
  const message = progress.cacheHit
    ? 'Using cached image artifacts'
    : PHASE_LABELS[progress.phase] || 'Working';
  return {
    visible: true,
    level: progress.done ? 'success' : 'progress',
    message,
    percent: progress.percent > 0 ? progress.percent : undefined,
  };
}

export function Dashboard({ onGeneratedImage, currentImage }: DashboardProps) {
  const [series, setSeries] = useState<dalle.Series[]>([]);
  const [selectedSeries, setSelectedSeries] = useState<string>('');
  const [input, setInput] = useState('Person Tour Coordinates');
  const [enhance, setEnhance] = useState(false);
  const [generateImage, setGenerateImage] = useState(false);
  const [annotate, setAnnotate] = useState(false);
  const [result, setResult] = useState<dalle.GenerateResult | null>(null);
  const [working, setWorking] = useState(false);
  const [progressTarget, setProgressTarget] = useState<ProgressTarget | null>(null);
  const [prefsLoaded, setPrefsLoaded] = useState(false);
  const [imageModel, setImageModelState] = useState(DEFAULT_IMAGE_MODEL);
  const [status, setStatus] = useState<StatusState>({
    visible: false,
    level: 'progress',
    message: '',
  });
  const [error, setError] = useState('');

  useEffect(() => {
    Promise.all([
      ListSeries(false, false),
      GetPref(DASHBOARD_PREFS.input),
      GetPref(DASHBOARD_PREFS.series),
      GetPref(DASHBOARD_PREFS.enhance),
      GetPref(DASHBOARD_PREFS.generateImage),
      GetPref(DASHBOARD_PREFS.annotate),
      GetPref(DASHBOARD_PREFS.imageModel),
      GetImageModel(),
    ])
      .then(
        ([
          items,
          savedInput,
          savedSeries,
          savedEnhance,
          savedGenerateImage,
          savedAnnotate,
          savedModel,
          engineModel,
        ]) => {
          const next = items ?? [];
          setSeries(next);
          setInput(savedInput || 'Person Tour Coordinates');
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
  }, [input, prefsLoaded]);

  useEffect(() => {
    if (!prefsLoaded) return;
    SetPref(DASHBOARD_PREFS.series, selectedSeries);
  }, [prefsLoaded, selectedSeries]);

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
      enhance,
      image: generateImage,
      annotate: generateImage && annotate,
    });

  const runPreview = () => {
    setError('');
    setStatus({ visible: true, level: 'progress', message: 'Building preview prompt' });
    setWorking(true);
    Preview(request())
      .then((next) => {
        setResult(next);
        setStatus({ visible: true, level: 'success', message: 'Preview ready' });
      })
      .catch((err: unknown) => {
        const message = messageFromError(err);
        setError(message);
        setStatus({ visible: true, level: 'error', message });
      })
      .finally(() => setWorking(false));
  };

  const runGenerate = () => {
    const nextRequest = request();
    setError('');
    setProgressTarget(null);
    setStatus({ visible: true, level: 'progress', message: 'Preparing generation' });
    setWorking(true);
    Preview(nextRequest)
      .then((preview) => {
        setResult(preview);
        setProgressTarget({ series: preview.series, seed: preview.seed });
        setStatus({
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
        setProgressTarget(null);
        setStatus({
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
        setStatus({ visible: true, level: 'error', message });
        setProgressTarget(null);
        setWorking(false);
      });
  };

  useEffect(() => {
    if (!working || !progressTarget) return;

    const poll = () => {
      GetGenerationProgress(progressTarget.series, progressTarget.seed)
        .then((progress) => {
          if (progress.active) setStatus(statusForProgress(progress));
        })
        .catch((err: unknown) => {
          setStatus({ visible: true, level: 'error', message: messageFromError(err) });
        });
    };

    poll();
    const interval = window.setInterval(poll, 750);
    return () => window.clearInterval(interval);
  }, [progressTarget, working]);

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
      <StatusBar
        visible={status.visible}
        level={status.level}
        message={status.message}
        meta={status.meta}
        percent={status.percent}
      />
    </Stack>
  );
}
