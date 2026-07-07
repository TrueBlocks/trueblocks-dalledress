export const DASHBOARD_PREFS = {
  input: 'dashboard.input',
  series: 'dashboard.series',
  enhance: 'dashboard.enhance',
  generateImage: 'dashboard.generateImage',
  annotate: 'dashboard.annotate',
  imageModel: 'settings.imageModel',
};

export const IMAGE_MODELS = [
  { value: 'gpt-image-1', label: 'GPT Image 1 (current default)' },
] as const;

export const DEFAULT_IMAGE_MODEL = 'gpt-image-1';

export function booleanPref(value: string): boolean {
  return value === 'true';
}

export function serializeBooleanPref(value: boolean): string {
  return String(value);
}
