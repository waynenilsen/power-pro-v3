import { useQuery } from '@tanstack/react-query';
import { queryKeys } from '../lib/query';
import { users } from '../api';
import type { ListLiftMaxesParams, ListProgressionHistoryParams } from '../api/endpoints/users';

export function useEnrollment(userId: string | undefined) {
  return useQuery({
    queryKey: queryKeys.users.enrollment(userId!),
    queryFn: () => users.getEnrollment(userId!),
    enabled: !!userId,
  });
}

export function useLiftMaxes(
  userId: string | undefined,
  params?: ListLiftMaxesParams
) {
  return useQuery({
    queryKey: queryKeys.users.liftMaxes(userId!, params),
    queryFn: () => users.listLiftMaxes(userId!, params),
    enabled: !!userId,
  });
}

export function useProgressionHistory(
  userId: string | undefined,
  params?: ListProgressionHistoryParams
) {
  return useQuery({
    queryKey: queryKeys.users.progressionHistory(userId!, params),
    queryFn: () => users.listProgressionHistory(userId!, params),
    enabled: !!userId,
  });
}
