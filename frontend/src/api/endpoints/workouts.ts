import { get, post, buildPaginationParams } from '../client';
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
  return post<EnrollmentResponse, AdvanceStateRequest>(`/users/${userId}/state/advance`, request);
}

// Workout Sessions

export async function startWorkoutSession(userId: string): Promise<WorkoutSessionDataResponse> {
  return post<WorkoutSessionDataResponse>(`/users/${userId}/workout-sessions/start`);
}

export async function finishWorkoutSession(userId: string): Promise<WorkoutSessionDataResponse> {
  return post<WorkoutSessionDataResponse>(`/users/${userId}/workout-sessions/finish`);
}

export async function abandonWorkoutSession(userId: string): Promise<WorkoutSessionDataResponse> {
  return post<WorkoutSessionDataResponse>(`/users/${userId}/workout-sessions/abandon`);
}

export async function getCurrentWorkoutSession(userId: string): Promise<WorkoutSessionDataResponse> {
  return get<WorkoutSessionDataResponse>(`/users/${userId}/workout-sessions/current`);
}

export interface ListWorkoutSessionsParams extends PaginationParams {
  status?: 'IN_PROGRESS' | 'COMPLETED' | 'ABANDONED';
}

export async function listWorkoutSessions(
  userId: string,
  params?: ListWorkoutSessionsParams
): Promise<WorkoutSessionListResponse> {
  return get<WorkoutSessionListResponse>(`/users/${userId}/workout-sessions`, {
    params: {
      ...buildPaginationParams(params),
      status: params?.status,
    },
  });
}

export async function getWorkoutSession(
  userId: string,
  sessionId: string
): Promise<WorkoutSessionDataResponse> {
  return get<WorkoutSessionDataResponse>(`/users/${userId}/workout-sessions/${sessionId}`);
}
