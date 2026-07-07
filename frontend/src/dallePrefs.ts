export const DASHBOARD_PREFS = {
  input: 'dashboard.input',
  series: 'dashboard.series',
  enhance: 'dashboard.enhance',
  generateImage: 'dashboard.generateImage',
  annotate: 'dashboard.annotate',
  imageModel: 'settings.imageModel',
  currentImageId: 'images.currentImageId',
};

export const IMAGE_MODELS = [
  { value: 'gpt-image-2', label: 'GPT Image 2 (latest, best quality)' },
  { value: 'gpt-image-1.5', label: 'GPT Image 1.5' },
  { value: 'gpt-image-1', label: 'GPT Image 1' },
  { value: 'gpt-image-1-mini', label: 'GPT Image 1 Mini (fast, cheap)' },
  { value: 'dall-e-3', label: 'DALL-E 3 (may require org verification)' },
  { value: 'dall-e-2', label: 'DALL-E 2 (legacy)' },
] as const;

export const DEFAULT_IMAGE_MODEL = 'gpt-image-2';

export function booleanPref(value: string): boolean {
  return value === 'true';
}

export function serializeBooleanPref(value: boolean): string {
  return String(value);
}
