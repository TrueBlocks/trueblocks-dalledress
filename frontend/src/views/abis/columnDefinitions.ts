import { FormField } from '@components';
import { types } from '@models';

// ABI columns for Downloaded and Known tabs
export const ABI_COLUMNS: FormField<types.Abi>[] = [
  { key: 'address', header: 'Address', sortable: true, type: 'text' },
  { key: 'name', header: 'Name', sortable: true, type: 'text' },
  { key: 'fileSize', header: 'File Size', sortable: true, type: 'number' },
  { key: 'nFunctions', header: 'Functions', sortable: true, type: 'number' },
  { key: 'nEvents', header: 'Events', sortable: true, type: 'number' },
  { key: 'nErrors', header: 'Errors', sortable: true, type: 'number' },
];

// Function/event columns for Functions and Events tabs
export const FUNCTION_COLUMNS: FormField<types.Function>[] = [
  { key: 'encoding', header: 'Encoding', sortable: true, type: 'text' },
  { key: 'name', header: 'Name', sortable: true, type: 'text' },
  { key: 'type', header: 'Type', sortable: true, type: 'text' },
  { key: 'signature', header: 'Signature', sortable: true, type: 'text' },
];

// Known tab omits address column
export const KNOWN_ABI_COLUMNS: FormField<types.Abi>[] = ABI_COLUMNS.slice(1);
