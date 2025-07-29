import { useCallback } from 'react';

import { types } from '@models';

import { useActiveProject } from './useActiveProject';

export const usePayload = () => {
  const { activeAddress, activeChain } = useActiveProject();
  return useCallback(
    (dataFacet: types.DataFacet, address?: string) => {
      return types.Payload.createFrom({
        dataFacet,
        chain: activeChain,
        address: address || activeAddress,
      });
    },
    [activeChain, activeAddress],
  );
};
