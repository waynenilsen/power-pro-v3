import { get, buildPaginationParams } from '../client';
import type { Lift, PaginatedResponse, PaginationParams } from '../types';

export type ListLiftsParams = PaginationParams;

export async function listLifts(params?: ListLiftsParams): Promise<PaginatedResponse<Lift>> {
  return get<PaginatedResponse<Lift>>('/lifts', {
    params: buildPaginationParams(params),
  });
}

export async function getLift(id: string): Promise<Lift> {
  return get<Lift>(`/lifts/${id}`);
}

export async function getLiftBySlug(slug: string): Promise<Lift> {
  return get<Lift>(`/lifts/by-slug/${slug}`);
}
