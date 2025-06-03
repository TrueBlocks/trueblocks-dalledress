// SortSpec utility functions for frontend
import { sdk } from '@models';

// Create a SortSpec for single field sorting (current UI behavior)
export const createSingleFieldSortSpec = (
  field: string,
  direction: 'asc' | 'desc',
): sdk.SortSpec => ({
  fields: [field],
  orders: [direction === 'asc'],
});

// Get the first field from a SortSpec (for single-field sorting UI)
export const getSortField = (sortSpec: sdk.SortSpec): string => {
  return sortSpec.fields && sortSpec.fields.length > 0
    ? sortSpec.fields[0] || ''
    : '';
};

// Get the first direction from a SortSpec as string (for single-field sorting UI)
export const getSortDirection = (sortSpec: sdk.SortSpec): 'asc' | 'desc' => {
  return sortSpec.orders && sortSpec.orders.length > 0 && sortSpec.orders[0]
    ? 'asc'
    : 'desc';
};

// Check if SortSpec is empty (no fields)
export const isEmptySort = (sortSpec: sdk.SortSpec): boolean => {
  return !sortSpec.fields || sortSpec.fields.length === 0;
};

// Create an empty SortSpec
export const EMPTY_SORT_SPEC: sdk.SortSpec = {
  fields: [],
  orders: [],
};

// Create an empty SortSpec (function version for when you need a new instance)
export const createEmptySortSpec = (): sdk.SortSpec => ({
  fields: [],
  orders: [],
});
