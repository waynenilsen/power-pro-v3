import type { ApiError, PaginationParams } from './types';

const BASE_URL = '/api';

export interface RequestOptions extends RequestInit {
  params?: Record<string, string | number | boolean | undefined>;
}

export class ApiClientError extends Error {
  public readonly status: number;
  public readonly details?: string[];

  constructor(message: string, status: number, details?: string[]) {
    super(message);
    this.name = 'ApiClientError';
    this.status = status;
    this.details = details;
  }
}

function buildUrl(path: string, params?: Record<string, string | number | boolean | undefined>): string {
  const url = new URL(`${BASE_URL}${path}`, window.location.origin);

  if (params) {
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined) {
        url.searchParams.append(key, String(value));
      }
    });
  }

  return url.pathname + url.search;
}

async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    let errorMessage = `Request failed with status ${response.status}`;
    let details: string[] | undefined;

    try {
      const errorBody: ApiError = await response.json();
      errorMessage = errorBody.error;
      details = errorBody.details;
    } catch {
      // Failed to parse error body, use default message
    }

    throw new ApiClientError(errorMessage, response.status, details);
  }

  // Handle 204 No Content
  if (response.status === 204) {
    return undefined as T;
  }

  return response.json();
}

function createHeaders(userId?: string, isAdmin?: boolean): HeadersInit {
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
  };

  if (userId) {
    headers['X-User-ID'] = userId;
  }

  if (isAdmin) {
    headers['X-Admin'] = 'true';
  }

  return headers;
}

export interface ClientConfig {
  userId?: string;
  isAdmin?: boolean;
}

let globalConfig: ClientConfig = {};

export function configureClient(config: ClientConfig): void {
  globalConfig = { ...globalConfig, ...config };
}

export function getClientConfig(): ClientConfig {
  return { ...globalConfig };
}

export async function get<T>(path: string, options?: RequestOptions): Promise<T> {
  const url = buildUrl(path, options?.params);

  const response = await fetch(url, {
    method: 'GET',
    headers: createHeaders(globalConfig.userId, globalConfig.isAdmin),
    ...options,
  });

  return handleResponse<T>(response);
}

export async function post<T, B = unknown>(path: string, body?: B, options?: RequestOptions): Promise<T> {
  const url = buildUrl(path, options?.params);

  const response = await fetch(url, {
    method: 'POST',
    headers: createHeaders(globalConfig.userId, globalConfig.isAdmin),
    body: body ? JSON.stringify(body) : undefined,
    ...options,
  });

  return handleResponse<T>(response);
}

export async function put<T, B = unknown>(path: string, body?: B, options?: RequestOptions): Promise<T> {
  const url = buildUrl(path, options?.params);

  const response = await fetch(url, {
    method: 'PUT',
    headers: createHeaders(globalConfig.userId, globalConfig.isAdmin),
    body: body ? JSON.stringify(body) : undefined,
    ...options,
  });

  return handleResponse<T>(response);
}

export async function del<T = void>(path: string, options?: RequestOptions): Promise<T> {
  const url = buildUrl(path, options?.params);

  const response = await fetch(url, {
    method: 'DELETE',
    headers: createHeaders(globalConfig.userId, globalConfig.isAdmin),
    ...options,
  });

  return handleResponse<T>(response);
}

export function buildPaginationParams(pagination?: PaginationParams): Record<string, string | number | undefined> {
  if (!pagination) return {};

  return {
    page: pagination.page,
    pageSize: pagination.pageSize,
    sortBy: pagination.sortBy,
    sortOrder: pagination.sortOrder,
  };
}
