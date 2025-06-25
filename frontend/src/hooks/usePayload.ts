import { useCallback } from 'react';

import { types } from '@models';

import { useActiveProject } from './useActiveProject';

export const usePayload = () => {
  const { effectiveAddress, effectiveChain } = useActiveProject();
  return useCallback(
    (dataFacet: types.DataFacet, address?: string) => {
      return types.Payload.createFrom({
        dataFacet,
        chain: effectiveChain,
        address: address || effectiveAddress,
      });
    },
    [effectiveChain, effectiveAddress],
  );
};
