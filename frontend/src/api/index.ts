// Client configuration
export { configureClient, getClientConfig, ApiClientError } from './client';
export type { ClientConfig, RequestOptions } from './client';

// Types
export * from './types';

// Endpoints
export * as programs from './endpoints/programs';
export * as workouts from './endpoints/workouts';
export * as users from './endpoints/users';
