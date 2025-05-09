import { useCallback, useEffect, useState } from 'react';

import { useTableContext } from '@components';
import { useFormHotkeys } from '@hooks';

interface SearchBoxProps {
  value: string;
  onChange: (v: string) => void;
}

export function SearchBox({ value, onChange }: SearchBoxProps) {
  useFormHotkeys({ keys: ['mod+a'] });

  const { focusControls } = useTableContext();
  const [inputValue, setInputValue] = useState(value);

  useEffect(() => {
    setInputValue(value);
  }, [value]);

  const handleFocus = useCallback(() => {
    focusControls();
  }, [focusControls]);

  const commitChange = useCallback(() => {
    if (inputValue !== value) {
      onChange(inputValue);
    }
  }, [inputValue, value, onChange]);

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      commitChange();
    }
  };

  const handleBlur = () => {
    setInputValue('');
  };

  return (
    <input
      type="text"
      placeholder="Search by name, address, tags, or source..."
      value={inputValue}
      onChange={(e) => setInputValue(e.target.value)}
      onFocus={handleFocus}
      onKeyDown={handleKeyDown}
      onBlur={handleBlur}
      style={{ width: 220, marginRight: 8, padding: 4 }}
      aria-label="Search table"
    />
  );
}
