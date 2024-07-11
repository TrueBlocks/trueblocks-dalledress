import React, { useState } from "react";
import { Table, Checkbox } from "@mantine/core";

export interface DataItem {
  id: number;
  value: string;
}

interface StringTableProps {
  data: DataItem[];
}

const StringTable: React.FC<StringTableProps> = ({ data }) => {
  const [selectedIds, setSelectedIds] = useState<number[]>([]);

  const handleCheckboxChange = (id: number) => {
    setSelectedIds((prevSelectedIds) =>
      prevSelectedIds.includes(id)
        ? prevSelectedIds.filter((selectedId) => selectedId !== id)
        : [...prevSelectedIds, id]
    );
  };

  const allChecked = data.length === selectedIds.length;
  const indeterminate = selectedIds.length > 0 && selectedIds.length < data.length;

  const handleAllCheckboxChange = () => {
    setSelectedIds(allChecked ? [] : data.map((item) => item.id));
  };

  const rows = data.map((item) => (
    <Table.Tr key={item.id}>
      <Table.Td>
        <Checkbox checked={selectedIds.includes(item.id)} onChange={() => handleCheckboxChange(item.id)} />
      </Table.Td>
      <Table.Td>{item.value}</Table.Td>
    </Table.Tr>
  ));

  return (
    <Table>
      <Table.Thead>
        <Table.Tr>
          <Table.Th>
            <Checkbox checked={allChecked} indeterminate={indeterminate} onChange={handleAllCheckboxChange} />
          </Table.Th>
          <Table.Th>Value</Table.Th>
        </Table.Tr>
      </Table.Thead>
      <Table.Tbody>{rows}</Table.Tbody>
    </Table>
  );
};

export default StringTable;
