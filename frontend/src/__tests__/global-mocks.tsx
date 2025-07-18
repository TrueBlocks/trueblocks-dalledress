// Copyright 2016, 2026 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */
import '@testing-library/jest-dom';
import { beforeEach, vi } from 'vitest';

beforeEach(() => {
  vi.clearAllMocks();
});

// Utility to generate icon mocks
const iconNames = [
  'Home',
  'Khedra',
  'Exports',
  'Monitors',
  'Abis',
  'Names',
  'Chunks',
  'Status',
  'DalleDress',
  'Settings',
  'Wizard',
  'Switch',
  'File',
  'Twitter',
  'Github',
  'Website',
  'Email',
  'Add',
  'Edit',
  'Delete',
  'Undelete',
  'Remove',
  'ChevronLeft',
  'ChevronRight',
  'ChevronUp',
  'ChevronDown',
  'Light',
  'Dark',
  'Autoname',
];

const iconMocks = Object.fromEntries(
  iconNames.map((name) => [
    name,
    (props: any) => (
      <div data-testid={`${name.toLowerCase()}-icon`} {...props}>
        {name} Icon
      </div>
    ),
  ]),
);

vi.mock('@utils', async (importOriginal) => {
  const original = await importOriginal();
  return {
    ...(original as object),
    Log: vi.fn(),
    checkAndNavigateToWizard: () => Promise.resolve(null),
    useEmitters: () => ({ emitStatus: vi.fn(), emitError: vi.fn() }),
  };
});

vi.mock('@hooks', async (importOriginal) => {
  const original = await importOriginal();
  return {
    ...(original as object),
    useIconSets: () => iconMocks,
  };
});

Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: (query: string) => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: () => {}, // Deprecated
    removeListener: () => {}, // Deprecated
    addEventListener: () => {},
    removeEventListener: () => {},
    dispatchEvent: () => false,
  }),
});
