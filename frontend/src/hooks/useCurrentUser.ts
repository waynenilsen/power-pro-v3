import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { queryKeys } from '../lib/query';
import { users, programs } from '../api';
import type { ListLiftMaxesParams, ListProgressionHistoryParams } from '../api/endpoints/users';
import type { EnrollRequest } from '../api/types';

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

export function useEnrollInProgram(userId: string | undefined) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (request: EnrollRequest) => {
      if (!userId) throw new Error('User ID is required');
      return programs.enrollInProgram(userId, request);
    },
    onSuccess: () => {
      if (userId) {
        queryClient.invalidateQueries({ queryKey: queryKeys.users.enrollment(userId) });
      }
    },
  });
}

export function useUnenroll(userId: string | undefined) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => {
      if (!userId) throw new Error('User ID is required');
      return users.unenroll(userId);
    },
    onSuccess: () => {
      if (userId) {
        queryClient.invalidateQueries({ queryKey: queryKeys.users.enrollment(userId) });
      }
    },
  });
}
