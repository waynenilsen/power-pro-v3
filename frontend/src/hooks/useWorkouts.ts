import { useQuery } from '@tanstack/react-query';
import { queryKeys } from '../lib/query';
import { workouts } from '../api';
import type { ListWorkoutSessionsParams } from '../api/endpoints/workouts';

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
