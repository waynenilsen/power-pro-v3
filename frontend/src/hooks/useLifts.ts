import { useQuery } from '@tanstack/react-query';
import { lifts } from '../api';
import type { ListLiftsParams } from '../api/endpoints/lifts';

export const liftQueryKeys = {
  all: ['lifts'] as const,
  lists: () => [...liftQueryKeys.all, 'list'] as const,
  list: <T extends object>(filters?: T) => [...liftQueryKeys.lists(), filters] as const,
};

export function useLifts(params?: ListLiftsParams) {
  return useQuery({
    queryKey: liftQueryKeys.list(params),
    queryFn: () => lifts.listLifts(params),
  });
}
