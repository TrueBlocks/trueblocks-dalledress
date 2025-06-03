import { useCallback, useState } from 'react';

import { Log, useEmitters } from '.';

export const useErrorHandler = () => {
  const [error, setError] = useState<Error | null>(null);
  const { emitError } = useEmitters();

  const handleError = useCallback(
    (err: unknown, context: string) => {
      const e = err instanceof Error ? err : new Error(String(err));
      setError(e);
      emitError(`${e.message} ${context}`);
      Log(`Error in ${context}: ${e}`);
    },
    [emitError],
  );

  const clearError = useCallback(() => setError(null), []);

  return { error, handleError, clearError };
};
