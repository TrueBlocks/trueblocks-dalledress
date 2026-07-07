import { useEffect, useState } from 'react';
import { Button, Checkbox, Group, Stack, Table, Text, Title } from '@mantine/core';
import { ListSeries } from '../../wailsjs/go/app/App';
import { dalle } from '../../wailsjs/go/models';

function messageFromError(error: unknown): string {
  return error instanceof Error ? error.message : String(error);
}

export function Series() {
  const [includeHidden, setIncludeHidden] = useState(false);
  const [items, setItems] = useState<dalle.Series[]>([]);
  const [error, setError] = useState('');

  const load = () => {
    setError('');
    ListSeries(includeHidden, false)
      .then((result) => setItems(result ?? []))
      .catch((err: unknown) => setError(messageFromError(err)));
  };

  useEffect(() => {
    setError('');
    ListSeries(includeHidden, false)
      .then((result) => setItems(result ?? []))
      .catch((err: unknown) => setError(messageFromError(err)));
  }, [includeHidden]);

  return (
    <Stack>
      <Title order={2}>Series</Title>
      <Group>
        <Checkbox
          label="Include hidden"
          checked={includeHidden}
          onChange={(event) => setIncludeHidden(event.currentTarget.checked)}
        />
        <Button onClick={load}>Refresh</Button>
      </Group>
      {error && <Text c="red">{error}</Text>}
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
            <Table.Tr key={item.suffix}>
              <Table.Td>{item.suffix}</Table.Td>
              <Table.Td>{item.purpose}</Table.Td>
              <Table.Td>{item.last ?? 0}</Table.Td>
              <Table.Td>{item.deleted ? 'hidden' : 'active'}</Table.Td>
            </Table.Tr>
          ))}
        </Table.Tbody>
      </Table>
    </Stack>
  );
}