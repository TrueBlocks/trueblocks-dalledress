import React from "react";
import { Modal, Text } from "@mantine/core";

interface ResultDialogProps {
  opened: boolean;
  onClose: () => void;
  success: boolean;
}

export const ResultDialog: React.FC<ResultDialogProps> = ({ opened, onClose, success }) => {
  return (
    <Modal opened={opened} onClose={onClose} title="Result">
      <Text>{success ? "Success!" : "Failed!"}</Text>
    </Modal>
  );
};
