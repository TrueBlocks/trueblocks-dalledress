import { useEffect, useState } from 'react';

import { GetFormat, SetFormat, SilenceDialog } from '@app';
import {
  Button,
  Checkbox,
  Group,
  Modal,
  Radio,
  Stack,
  Text,
} from '@mantine/core';
import { Log } from '@utils';

export interface ExportFormatModalProps {
  opened: boolean;
  onClose: () => void;
  onFormatSelected: (format: string) => void;
}

const formatOptions = [
  { value: 'csv', label: 'CSV - Comma separated values (.csv)' },
  { value: 'txt', label: 'TXT - Tab separated values (.txt)' },
  { value: 'json', label: 'JSON - JavaScript Object Notation (.json)' },
];

export const ExportFormatModal = ({
  opened,
  onClose,
  onFormatSelected,
}: ExportFormatModalProps) => {
  const [selectedFormat, setSelectedFormat] = useState<string>('csv');
  const [dontShowAgain, setDontShowAgain] = useState(false);
  const [loading, setLoading] = useState(false);

  // Load the last format preference when modal opens
  useEffect(() => {
    if (opened) {
      setLoading(true);
      GetFormat()
        .then((lastFormat: string) => {
          Log(`[EXPORT FORMAT MODAL] Loaded last format: ${lastFormat}`);
          setSelectedFormat(lastFormat || 'csv');
        })
        .catch((error: Error) => {
          Log(`[EXPORT FORMAT MODAL] Error loading format: ${error}`);
          setSelectedFormat('csv');
        })
        .finally(() => {
          setLoading(false);
        });
    }
  }, [opened]);

  const handleFormatSelect = async (format: string) => {
    Log(
      `[EXPORT FORMAT MODAL] Format selected: ${format}, dontShowAgain: ${dontShowAgain}`,
    );

    try {
      // Save the selected format preference
      await SetFormat(format);
      Log(`[EXPORT FORMAT MODAL] Format preference saved: ${format}`);

      // If user chose "don't show again", silence the dialog
      if (dontShowAgain) {
        await SilenceDialog('exportFormat');
        Log('[EXPORT FORMAT MODAL] Export format dialog silenced');
      }

      // Close modal and proceed with export
      onClose();
      onFormatSelected(format);
    } catch (error) {
      Log(`[EXPORT FORMAT MODAL] Error saving preferences: ${error}`);
      // Still proceed with export even if preference saving fails
      onClose();
      onFormatSelected(format);
    }
  };

  const handleCancel = () => {
    Log('[EXPORT FORMAT MODAL] Export cancelled by user');
    onClose();
  };

  return (
    <Modal
      opened={opened}
      onClose={handleCancel}
      title="Select Export Format"
      size="md"
      centered
      withCloseButton={false}
    >
      <Stack gap="md">
        <Text size="sm" c="dimmed">
          Choose the format for your exported data:
        </Text>

        <Radio.Group
          value={selectedFormat}
          onChange={setSelectedFormat}
          name="exportFormat"
        >
          <Stack gap="xs">
            {formatOptions.map((option) => (
              <Radio
                key={option.value}
                value={option.value}
                label={option.label}
                disabled={loading}
              />
            ))}
          </Stack>
        </Radio.Group>

        <Checkbox
          checked={dontShowAgain}
          onChange={(event) => setDontShowAgain(event.currentTarget.checked)}
          label="Don't show this dialog again"
          disabled={loading}
        />

        <Group justify="flex-end" gap="sm">
          <Button variant="subtle" onClick={handleCancel} disabled={loading}>
            Cancel
          </Button>
          <Button
            onClick={() => handleFormatSelect(selectedFormat)}
            loading={loading}
          >
            Export
          </Button>
        </Group>
      </Stack>
    </Modal>
  );
};
