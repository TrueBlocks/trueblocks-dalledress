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
import { Generate, ListSeries, Preview } from '../../wailsjs/go/app/App';
import { dalle } from '../../wailsjs/go/models';

function messageFromError(error: unknown): string {
  return error instanceof Error ? error.message : String(error);
}

export function Dashboard() {
  const [series, setSeries] = useState<dalle.Series[]>([]);
  const [selectedSeries, setSelectedSeries] = useState<string>('');
  const [input, setInput] = useState('Person Tour Coordinates');
  const [generateImage, setGenerateImage] = useState(false);
  const [annotate, setAnnotate] = useState(false);
  const [result, setResult] = useState<dalle.GenerateResult | null>(null);
  const [error, setError] = useState('');

  useEffect(() => {
    ListSeries(false, false)
      .then((items) => {
        const next = items ?? [];
        setSeries(next);
        setSelectedSeries((current) => current || next[0]?.suffix || '');
      })
      .catch((err: unknown) => setError(messageFromError(err)));
  }, []);

  const request = () =>
    dalle.GenerateRequest.createFrom({
      input,
      series: selectedSeries || undefined,
      image: generateImage,
      annotate: generateImage && annotate,
    });

  const runPreview = () => {
    setError('');
    Preview(request())
      .then(setResult)
      .catch((err: unknown) => setError(messageFromError(err)));
  };

  const runGenerate = () => {
    setError('');
    Generate(request())
      .then(setResult)
      .catch((err: unknown) => setError(messageFromError(err)));
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
        <Button onClick={runPreview}>Preview</Button>
        <Button onClick={runGenerate}>Generate</Button>
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
