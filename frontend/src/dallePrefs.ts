export const DASHBOARD_PREFS = {
  input: 'dashboard.input',
  series: 'dashboard.series',
  enhance: 'dashboard.enhance',
  generateImage: 'dashboard.generateImage',
  annotate: 'dashboard.annotate',
};

export function booleanPref(value: string): boolean {
  return value === 'true';
}

export function serializeBooleanPref(value: boolean): string {
  return String(value);
}
