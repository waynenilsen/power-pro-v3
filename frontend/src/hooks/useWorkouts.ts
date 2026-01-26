import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { queryKeys } from '../lib/query';
import { workouts } from '../api';
import type { ListWorkoutSessionsParams } from '../api/endpoints/workouts';
import type { AdvanceStateRequest } from '../api/types';

export function useCurrentWorkout(userId: string | undefined) {
  return useQuery({
    queryKey: queryKeys.workouts.current(userId!),
    queryFn: () => workouts.getCurrentWorkout(userId!),
    enabled: !!userId,
  });
}

export function useWorkoutForDay(
  userId: string | undefined,
  week: number,
  dayIndex: number
) {
  return useQuery({
    queryKey: queryKeys.workouts.forDay(userId!, week, dayIndex),
    queryFn: () => workouts.getWorkoutForDay(userId!, week, dayIndex),
    enabled: !!userId,
  });
}

export function useWorkoutSessions(
  userId: string | undefined,
  params?: ListWorkoutSessionsParams
) {
  return useQuery({
    queryKey: queryKeys.workouts.sessionList(userId!, params),
    queryFn: () => workouts.listWorkoutSessions(userId!, params),
    enabled: !!userId,
  });
}

export function useWorkoutSession(
  userId: string | undefined,
  sessionId: string | undefined
) {
  return useQuery({
    queryKey: queryKeys.workouts.sessionDetail(userId!, sessionId!),
    queryFn: () => workouts.getWorkoutSession(userId!, sessionId!),
    enabled: !!userId && !!sessionId,
  });
}

export function useCurrentWorkoutSession(userId: string | undefined) {
  return useQuery({
    queryKey: queryKeys.workouts.sessionCurrent(userId!),
    queryFn: () => workouts.getCurrentWorkoutSession(userId!),
    enabled: !!userId,
  });
}

export function useStartWorkoutSession(userId: string | undefined) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => {
      if (!userId) throw new Error('User ID is required');
      return workouts.startWorkoutSession(userId);
    },
    onSuccess: () => {
      if (userId) {
        queryClient.invalidateQueries({ queryKey: queryKeys.workouts.sessionCurrent(userId) });
        queryClient.invalidateQueries({ queryKey: queryKeys.workouts.sessions(userId) });
        queryClient.invalidateQueries({ queryKey: queryKeys.users.enrollment(userId) });
      }
    },
  });
}

export function useFinishWorkoutSession(userId: string | undefined) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => {
      if (!userId) throw new Error('User ID is required');
      return workouts.finishWorkoutSession(userId);
    },
    onSuccess: () => {
      if (userId) {
        queryClient.invalidateQueries({ queryKey: queryKeys.workouts.sessionCurrent(userId) });
        queryClient.invalidateQueries({ queryKey: queryKeys.workouts.sessions(userId) });
        queryClient.invalidateQueries({ queryKey: queryKeys.users.enrollment(userId) });
      }
    },
  });
}

export function useAbandonWorkoutSession(userId: string | undefined) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => {
      if (!userId) throw new Error('User ID is required');
      return workouts.abandonWorkoutSession(userId);
    },
    onSuccess: () => {
      if (userId) {
        queryClient.invalidateQueries({ queryKey: queryKeys.workouts.sessionCurrent(userId) });
        queryClient.invalidateQueries({ queryKey: queryKeys.workouts.sessions(userId) });
        queryClient.invalidateQueries({ queryKey: queryKeys.users.enrollment(userId) });
      }
    },
  });
}

export function useAdvanceState(userId: string | undefined) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (request: AdvanceStateRequest) => {
      if (!userId) throw new Error('User ID is required');
      return workouts.advanceState(userId, request);
    },
    onSuccess: () => {
      if (userId) {
        queryClient.invalidateQueries({ queryKey: queryKeys.users.enrollment(userId) });
        queryClient.invalidateQueries({ queryKey: queryKeys.workouts.current(userId) });
      }
    },
  });
}
