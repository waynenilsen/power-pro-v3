import { get, post, buildPaginationParams } from '../client';
import type {
  PaginatedResponse,
  PaginationParams,
  ProgramListItem,
  ProgramDetail,
  EnrollRequest,
  EnrollmentResponse,
} from '../types';

export type ListProgramsParams = PaginationParams;

export async function listPrograms(params?: ListProgramsParams): Promise<PaginatedResponse<ProgramListItem>> {
  return get<PaginatedResponse<ProgramListItem>>('/programs', {
    params: buildPaginationParams(params),
  });
}

export async function getProgram(id: string): Promise<ProgramDetail> {
  return get<ProgramDetail>(`/programs/${id}`);
}

export async function getProgramBySlug(slug: string): Promise<ProgramDetail> {
  return get<ProgramDetail>(`/programs/by-slug/${slug}`);
}

export async function enrollInProgram(userId: string, request: EnrollRequest): Promise<EnrollmentResponse> {
  return post<EnrollmentResponse, EnrollRequest>(`/users/${userId}/program`, request);
}
