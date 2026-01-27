import { get, getRaw, post, put, del, buildPaginationParams } from '../client';
import type {
  EnrollmentResponse,
  LiftMax,
  LiftMaxWithWarnings,
  CreateLiftMaxRequest,
  UpdateLiftMaxRequest,
  PaginatedResponse,
  PaginationParams,
  MaxType,
  ProgressionHistoryEntry,
  ManualTriggerRequest,
  ManualTriggerResponse,
} from '../types';

// Enrollment

export async function getEnrollment(userId: string): Promise<EnrollmentResponse> {
  return get<EnrollmentResponse>(`/users/${userId}/program`);
}

export async function unenroll(userId: string): Promise<void> {
  return del<void>(`/users/${userId}/program`);
}

// Lift Maxes

export interface ListLiftMaxesParams extends PaginationParams {
  liftId?: string;
  type?: MaxType;
}

export async function listLiftMaxes(
  userId: string,
  params?: ListLiftMaxesParams
): Promise<PaginatedResponse<LiftMax>> {
  return getRaw<PaginatedResponse<LiftMax>>(`/users/${userId}/lift-maxes`, {
    params: {
      ...buildPaginationParams(params),
      lift_id: params?.liftId,
      type: params?.type,
    },
  });
}

export async function getLiftMax(userId: string, liftMaxId: string): Promise<LiftMax> {
  return get<LiftMax>(`/users/${userId}/lift-maxes/${liftMaxId}`);
}

export async function createLiftMax(
  userId: string,
  request: CreateLiftMaxRequest
): Promise<LiftMaxWithWarnings> {
  return post<LiftMaxWithWarnings, CreateLiftMaxRequest>(`/users/${userId}/lift-maxes`, request);
}

export async function updateLiftMax(
  userId: string,
  liftMaxId: string,
  request: UpdateLiftMaxRequest
): Promise<LiftMax> {
  // Note: Backend uses /lift-maxes/{id} not /users/{userId}/lift-maxes/{id}
  return put<LiftMax, UpdateLiftMaxRequest>(`/lift-maxes/${liftMaxId}`, request);
}

export async function deleteLiftMax(userId: string, liftMaxId: string): Promise<void> {
  // Note: Backend uses /lift-maxes/{id} not /users/{userId}/lift-maxes/{id}
  return del<void>(`/lift-maxes/${liftMaxId}`);
}

export async function getLatestLiftMax(
  userId: string,
  liftId: string,
  type: MaxType
): Promise<LiftMax> {
  return get<LiftMax>(`/users/${userId}/lift-maxes/latest/${liftId}`, {
    params: { type },
  });
}

// Enrollment State Management

export async function startWeek(userId: string): Promise<EnrollmentResponse> {
  return post<EnrollmentResponse>(`/users/${userId}/enrollment/start-week`);
}

export async function completeWeek(userId: string): Promise<EnrollmentResponse> {
  return post<EnrollmentResponse>(`/users/${userId}/enrollment/complete-week`);
}

export async function startCycle(userId: string): Promise<EnrollmentResponse> {
  return post<EnrollmentResponse>(`/users/${userId}/enrollment/start-cycle`);
}

export async function completeCycle(userId: string): Promise<EnrollmentResponse> {
  return post<EnrollmentResponse>(`/users/${userId}/enrollment/complete-cycle`);
}

// Progression History

export interface ListProgressionHistoryParams extends PaginationParams {
  liftId?: string;
  progressionId?: string;
}

export async function listProgressionHistory(
  userId: string,
  params?: ListProgressionHistoryParams
): Promise<PaginatedResponse<ProgressionHistoryEntry>> {
  return get<PaginatedResponse<ProgressionHistoryEntry>>(`/users/${userId}/progression-history`, {
    params: {
      ...buildPaginationParams(params),
      lift_id: params?.liftId,
      progression_id: params?.progressionId,
    },
  });
}

// Manual Progression

export async function triggerProgression(
  userId: string,
  request: ManualTriggerRequest
): Promise<ManualTriggerResponse> {
  return post<ManualTriggerResponse, ManualTriggerRequest>(
    `/users/${userId}/progressions/trigger`,
    request
  );
}
