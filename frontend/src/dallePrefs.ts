export const DASHBOARD_PREFS = {
  input: 'dashboard.input',
  series: 'dashboard.series',
  backstyle: 'dashboard.backstyle',
  enhance: 'dashboard.enhance',
  generateImage: 'dashboard.generateImage',
  annotate: 'dashboard.annotate',
  imageModel: 'settings.imageModel',
  currentImageId: 'images.currentImageId',
};

export const IMAGE_MODELS = [
  { value: 'gemini-3-pro-image', label: 'Gemini 3 Pro Image (pro tier)' },
  { value: 'gemini-3.1-flash-image', label: 'Gemini 3.1 Flash Image (cheap tier)' },
  { value: 'gpt-image-2', label: 'GPT Image 2 (latest OpenAI)' },
  { value: 'gpt-image-1.5', label: 'GPT Image 1.5' },
  { value: 'gpt-image-1', label: 'GPT Image 1' },
  { value: 'gpt-image-1-mini', label: 'GPT Image 1 Mini (fast, cheap)' },
  { value: 'dall-e-3', label: 'DALL-E 3 (may require org verification)' },
  { value: 'dall-e-2', label: 'DALL-E 2 (legacy)' },
] as const;

// The engine's own default (the registry's pro image model) is preferred;
// this is only the picker's value when the engine reports none.
export const DEFAULT_IMAGE_MODEL = 'gemini-3-pro-image';

export function booleanPref(value: string): boolean {
  return value === 'true';
}

export function serializeBooleanPref(value: boolean): string {
  return String(value);
}
