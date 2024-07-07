import React, { useState, useEffect } from 'react';
import { Select } from '@mantine/core';
import { SetLastSeries } from '@gocode/app/App';

type SelectItem = {
  value: string;
  label: string;
};

async function GetExistingAddrs(): Promise<string[]> {
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve([
        'gitcoin.eth',
        'giveth.eth',
        'chase.wright.eth',
        'cnn.eth',
        'dawid.eth',
        'dragonstone.eth',
        'eats.eth',
        'ens.eth',
        'gameofthrones.eth',
        'jen.eth',
        'makingprogress.eth',
        'meriam.eth',
        'nate.eth',
        'poap.eth',
        'revenge.eth',
        'rotki.eth',
        'trueblocks.eth',
        'unchainedindex.eth',
        'vitalik.eth',
        'when.eth',
      ]);
    }, 1000);
  });
}

interface EditableSelectProps {
  value: string;
  onChange: (value: string) => void;
  label: string;
  placeholder: string;
}

const EditableSelect: React.FC<EditableSelectProps> = ({ value, onChange, label, placeholder }) => {
  const [options, setOptions] = useState<SelectItem[]>([]);
  const [inputValue, setInputValue] = useState<string>(value);

  useEffect(() => {
    async function fetchAddresses() {
      GetExistingAddrs().then((existingAddrs) => {
        const formattedOptions = existingAddrs.map((addr) => ({ value: addr, label: addr }));
        setOptions(formattedOptions);
        if (!existingAddrs.includes(value)) {
          setInputValue(existingAddrs.length > 0 ? existingAddrs[0] : '');
        }
      });
    }
    fetchAddresses();
  }, []);

  const handleKeyDown = (event: React.KeyboardEvent<HTMLInputElement>) => {
    if (event.key === 'Enter' && inputValue.trim() !== '') {
      const newOption: SelectItem = { value: inputValue, label: inputValue };
      setOptions((current) => {
        const existingIndex = current.findIndex((option) => option.value.toLowerCase() === inputValue.toLowerCase());
        if (existingIndex === -1) {
          return [...current, newOption];
        } else {
          return current.map((option, index) => (index === existingIndex ? newOption : option));
        }
      });
      onChange(inputValue);
      event.preventDefault();
      setTimeout(() => {
        const selectElement = document.querySelector('.mantine-Select-input') as HTMLElement;
        if (selectElement) {
          selectElement.blur();
        }
      }, 0);
    }
  };

  const handleChange = (value: string | null) => {
    if (value !== null) {
      setInputValue(value);
      onChange(value);
    }
  };

  return (
    <Select
      label={label}
      placeholder={placeholder}
      data={options}
      searchable
      value={inputValue}
      onSearchChange={setInputValue}
      onChange={handleChange}
      onKeyDown={handleKeyDown}
      style={{ width: '600px' }}
    />
  );
};

export default EditableSelect;
