import { useMutation, useQueryClient } from '@tanstack/react-query';
import { queryKeys } from '../lib/query';
import { users } from '../api';
import type { CreateLiftMaxRequest, UpdateLiftMaxRequest } from '../api/types';

export function useCreateLiftMax(userId: string | undefined) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (request: CreateLiftMaxRequest) => {
      if (!userId) throw new Error('User ID is required');
      return users.createLiftMax(userId, request);
    },
    onSuccess: () => {
      if (userId) {
        queryClient.invalidateQueries({ queryKey: queryKeys.users.liftMaxes(userId) });
      }
    },
  });
}

export function useUpdateLiftMax(userId: string | undefined) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ liftMaxId, request }: { liftMaxId: string; request: UpdateLiftMaxRequest }) => {
      if (!userId) throw new Error('User ID is required');
      return users.updateLiftMax(userId, liftMaxId, request);
    },
    onSuccess: () => {
      if (userId) {
        queryClient.invalidateQueries({ queryKey: queryKeys.users.liftMaxes(userId) });
      }
    },
  });
}

export function useDeleteLiftMax(userId: string | undefined) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (liftMaxId: string) => {
      if (!userId) throw new Error('User ID is required');
      return users.deleteLiftMax(userId, liftMaxId);
    },
    onSuccess: () => {
      if (userId) {
        queryClient.invalidateQueries({ queryKey: queryKeys.users.liftMaxes(userId) });
      }
    },
  });
}
