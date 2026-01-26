import type { ApiError, PaginationParams } from './types';
import { triggerUnauthorized } from './auth-error-handler';

const BASE_URL = '/api';

/**
 * Extended RequestInit options with query parameter support.
 */
export interface RequestOptions extends RequestInit {
  /** Query parameters to append to the URL */
  params?: Record<string, string | number | boolean | undefined>;
}

/**
 * Custom error class for API client errors.
 * Provides structured error information including HTTP status and optional details.
 */
export class ApiClientError extends Error {
  /** HTTP status code of the failed request */
  public readonly status: number;
  /** Additional error details from the API response */
  public readonly details?: string[];

  constructor(message: string, status: number, details?: string[]) {
    super(message);
    this.name = 'ApiClientError';
    this.status = status;
    this.details = details;
  }
}

/**
 * Builds a URL with query parameters.
 * @param path - The API endpoint path
 * @param params - Optional query parameters
 * @returns The complete URL path with query string
 */
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

/**
 * Handles API response, parsing JSON or throwing appropriate errors.
 * Triggers logout on 401 Unauthorized responses.
 * @param response - The fetch Response object
 * @returns Parsed JSON response data
 * @throws ApiClientError for non-2xx responses
 */
/**
 * API response wrapper type - backend wraps all responses in { data: T }
 */
interface ApiResponse<T> {
  data: T;
}

async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    let errorMessage = `Request failed with status ${response.status}`;
    let details: string[] | undefined;

    try {
      const errorBody: ApiError = await response.json();
      errorMessage = errorBody.error.message;
      details = errorBody.error.details?.validationErrors;
    } catch {
      // Failed to parse error body, use default message
    }

    // Trigger logout on 401 Unauthorized
    if (response.status === 401) {
      triggerUnauthorized();
    }

    throw new ApiClientError(errorMessage, response.status, details);
  }

  // Handle 204 No Content
  if (response.status === 204) {
    return undefined as T;
  }

  // Backend wraps all responses in { data: T }, so unwrap it
  const json: ApiResponse<T> = await response.json();
  return json.data;
}

/**
 * Creates request headers including authentication headers.
 * @param userId - Optional user ID for X-User-ID header
 * @param isAdmin - Optional admin flag for X-Admin header
 * @param token - Optional Bearer token for Authorization header
 * @returns Headers object for fetch requests
 */
function createHeaders(userId?: string, isAdmin?: boolean, token?: string): HeadersInit {
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
  };

  // Bearer token takes precedence over X-User-ID if both are set
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  } else if (userId) {
    headers['X-User-ID'] = userId;
  }

  if (isAdmin) {
    headers['X-Admin'] = 'true';
  }

  return headers;
}

/**
 * Configuration for the API client.
 */
export interface ClientConfig {
  /** User ID to include in X-User-ID header */
  userId?: string;
  /** Whether to include X-Admin header */
  isAdmin?: boolean;
  /** Bearer token for Authorization header */
  token?: string;
}

// Storage keys for auth - must match auth-types.ts
const AUTH_STORAGE_KEY = 'powerpro_user_id';
const TOKEN_STORAGE_KEY = 'powerpro_token';

// Initialize config from localStorage synchronously on module load
function getInitialConfig(): ClientConfig {
  if (typeof window === 'undefined') return {};

  const token = localStorage.getItem(TOKEN_STORAGE_KEY);
  const userId = localStorage.getItem(AUTH_STORAGE_KEY);

  const config: ClientConfig = {};
  if (token) config.token = token;
  if (userId) config.userId = userId;

  return config;
}

let globalConfig: ClientConfig = getInitialConfig();

/**
 * Configures the global API client settings.
 * Settings are merged with existing configuration.
 * @param config - Configuration options to apply
 */
export function configureClient(config: ClientConfig): void {
  globalConfig = { ...globalConfig, ...config };
}

/**
 * Returns a copy of the current client configuration.
 * @returns Current client configuration
 */
export function getClientConfig(): ClientConfig {
  return { ...globalConfig };
}

/**
 * Handles response without unwrapping {data: T} - for paginated endpoints
 * that return {data: items[], meta: {...}} directly.
 */
async function handleRawResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    let errorMessage = `Request failed with status ${response.status}`;
    let details: string[] | undefined;

    try {
      const errorBody: ApiError = await response.json();
      errorMessage = errorBody.error.message;
      details = errorBody.error.details?.validationErrors;
    } catch {
      // Failed to parse error body, use default message
    }

    if (response.status === 401) {
      triggerUnauthorized();
    }

    throw new ApiClientError(errorMessage, response.status, details);
  }

  if (response.status === 204) {
    return undefined as T;
  }

  return response.json();
}

/**
 * Performs a GET request to the API.
 * @param path - The API endpoint path
 * @param options - Optional request options including query parameters
 * @returns Promise resolving to the response data
 * @throws ApiClientError for non-2xx responses
 */
export async function get<T>(path: string, options?: RequestOptions): Promise<T> {
  const url = buildUrl(path, options?.params);

  const response = await fetch(url, {
    method: 'GET',
    headers: createHeaders(globalConfig.userId, globalConfig.isAdmin, globalConfig.token),
    ...options,
  });

  return handleResponse<T>(response);
}

/**
 * Performs a GET request without unwrapping {data: T}.
 * Use for paginated endpoints that return {data: items[], meta: {...}}.
 */
export async function getRaw<T>(path: string, options?: RequestOptions): Promise<T> {
  const url = buildUrl(path, options?.params);

  const response = await fetch(url, {
    method: 'GET',
    headers: createHeaders(globalConfig.userId, globalConfig.isAdmin, globalConfig.token),
    ...options,
  });

  return handleRawResponse<T>(response);
}

/**
 * Performs a POST request to the API.
 * @param path - The API endpoint path
 * @param body - Optional request body (will be JSON stringified)
 * @param options - Optional request options including query parameters
 * @returns Promise resolving to the response data
 * @throws ApiClientError for non-2xx responses
 */
export async function post<T, B = unknown>(path: string, body?: B, options?: RequestOptions): Promise<T> {
  const url = buildUrl(path, options?.params);

  const response = await fetch(url, {
    method: 'POST',
    headers: createHeaders(globalConfig.userId, globalConfig.isAdmin, globalConfig.token),
    body: body ? JSON.stringify(body) : undefined,
    ...options,
  });

  return handleResponse<T>(response);
}

/**
 * Performs a PUT request to the API.
 * @param path - The API endpoint path
 * @param body - Optional request body (will be JSON stringified)
 * @param options - Optional request options including query parameters
 * @returns Promise resolving to the response data
 * @throws ApiClientError for non-2xx responses
 */
export async function put<T, B = unknown>(path: string, body?: B, options?: RequestOptions): Promise<T> {
  const url = buildUrl(path, options?.params);

  const response = await fetch(url, {
    method: 'PUT',
    headers: createHeaders(globalConfig.userId, globalConfig.isAdmin, globalConfig.token),
    body: body ? JSON.stringify(body) : undefined,
    ...options,
  });

  return handleResponse<T>(response);
}

/**
 * Performs a DELETE request to the API.
 * @param path - The API endpoint path
 * @param options - Optional request options including query parameters
 * @returns Promise resolving to the response data (or void for 204 responses)
 * @throws ApiClientError for non-2xx responses
 */
export async function del<T = void>(path: string, options?: RequestOptions): Promise<T> {
  const url = buildUrl(path, options?.params);

  const response = await fetch(url, {
    method: 'DELETE',
    headers: createHeaders(globalConfig.userId, globalConfig.isAdmin, globalConfig.token),
    ...options,
  });

  return handleResponse<T>(response);
}

/**
 * Builds query parameters object from pagination options.
 * Filters out undefined values for clean query strings.
 * @param pagination - Optional pagination parameters
 * @returns Query parameters object for use with request options
 */
export function buildPaginationParams(pagination?: PaginationParams): Record<string, string | number | undefined> {
  if (!pagination) return {};

  return {
    page: pagination.page,
    pageSize: pagination.pageSize,
    sortBy: pagination.sortBy,
    sortOrder: pagination.sortOrder,
  };
}
