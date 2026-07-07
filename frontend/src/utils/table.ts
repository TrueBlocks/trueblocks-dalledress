export function uniqueSortedValues(values: (string | undefined | null)[]): string[] {
  return Array.from(new Set(values.filter((value): value is string => Boolean(value)))).sort(
    (left, right) => left.localeCompare(right, undefined, { sensitivity: 'base' }),
  );
}
