// Client configuration
export { configureClient, getClientConfig, ApiClientError } from './client';
export type { ClientConfig, RequestOptions } from './client';

// Auth error handling
export { onUnauthorized, clearUnauthorizedHandler } from './auth-error-handler';

// Types
export * from './types';

// Endpoints
export * as programs from './endpoints/programs';
export * as workouts from './endpoints/workouts';
export * as users from './endpoints/users';
export * as lifts from './endpoints/lifts';
