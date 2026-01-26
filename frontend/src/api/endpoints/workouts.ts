import { get, getRaw, post, buildPaginationParams } from '../client';
import type {
  Workout,
  WorkoutSessionDataResponse,
  WorkoutSessionListResponse,
  PaginationParams,
  AdvanceStateRequest,
  EnrollmentResponse,
} from '../types';

export async function getCurrentWorkout(userId: string): Promise<Workout> {
  return get<Workout>(`/users/${userId}/workout`);
}

export async function getWorkoutForDay(
  userId: string,
  week: number,
  dayIndex: number
): Promise<Workout> {
  return get<Workout>(`/users/${userId}/workout`, {
    params: { week, dayIndex },
  });
}

export async function advanceState(
  userId: string,
  request: AdvanceStateRequest
): Promise<EnrollmentResponse> {
  return post<EnrollmentResponse, AdvanceStateRequest>(`/users/${userId}/program-state/advance`, request);
}

// Workout Sessions

export async function startWorkoutSession(userId: string): Promise<WorkoutSessionDataResponse> {
  return post<WorkoutSessionDataResponse>(`/workouts/start`);
}

export async function finishWorkoutSession(sessionId: string): Promise<WorkoutSessionDataResponse> {
  return post<WorkoutSessionDataResponse>(`/workouts/${sessionId}/finish`);
}

export async function abandonWorkoutSession(sessionId: string): Promise<WorkoutSessionDataResponse> {
  return post<WorkoutSessionDataResponse>(`/workouts/${sessionId}/abandon`);
}

export async function getCurrentWorkoutSession(userId: string): Promise<WorkoutSessionDataResponse> {
  return get<WorkoutSessionDataResponse>(`/users/${userId}/workouts/current`);
}

export interface ListWorkoutSessionsParams extends PaginationParams {
  status?: 'IN_PROGRESS' | 'COMPLETED' | 'ABANDONED';
}

export async function listWorkoutSessions(
  userId: string,
  params?: ListWorkoutSessionsParams
): Promise<WorkoutSessionListResponse> {
  return getRaw<WorkoutSessionListResponse>(`/users/${userId}/workouts`, {
    params: {
      ...buildPaginationParams(params),
      status: params?.status,
    },
  });
}

export async function getWorkoutSession(
  sessionId: string
): Promise<WorkoutSessionDataResponse> {
  return get<WorkoutSessionDataResponse>(`/workouts/${sessionId}`);
}
