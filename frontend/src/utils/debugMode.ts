export const isDebugMode = (): boolean => {
  return import.meta.env.VITE_TB_DEBUG === 'true';
};
