import { QueryClient } from '@tanstack/react-query';
import { ApiClientError } from '../api/client';
import { triggerUnauthorized } from '../api/auth-error-handler';

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000, // 5 minutes
      gcTime: 10 * 60 * 1000, // 10 minutes (formerly cacheTime)
      retry: (failureCount, error) => {
        // Don't retry on 4xx errors (client errors)
        if (error instanceof ApiClientError && error.status >= 400 && error.status < 500) {
          return false;
        }
        // Retry up to 3 times for other errors
        return failureCount < 3;
      },
      refetchOnWindowFocus: false,
    },
    mutations: {
      retry: false,
      onError: (error) => {
        // Trigger logout on 401 Unauthorized for mutations
        // (Queries are handled at the API client level)
        if (error instanceof ApiClientError && error.status === 401) {
          triggerUnauthorized();
        }
      },
    },
  },
});

// Query key factory for consistent key structure
export const queryKeys = {
  // Programs
  programs: {
    all: ['programs'] as const,
    lists: () => [...queryKeys.programs.all, 'list'] as const,
    list: <T extends object>(filters?: T) => [...queryKeys.programs.lists(), filters] as const,
    details: () => [...queryKeys.programs.all, 'detail'] as const,
    detail: (id: string) => [...queryKeys.programs.details(), id] as const,
    bySlug: (slug: string) => [...queryKeys.programs.all, 'slug', slug] as const,
  },

  // User-specific queries
  users: {
    all: ['users'] as const,
    detail: (userId: string) => [...queryKeys.users.all, userId] as const,
    enrollment: (userId: string) => [...queryKeys.users.detail(userId), 'enrollment'] as const,
    liftMaxes: <T extends object>(userId: string, filters?: T) =>
      [...queryKeys.users.detail(userId), 'lift-maxes', filters] as const,
    progressionHistory: <T extends object>(userId: string, filters?: T) =>
      [...queryKeys.users.detail(userId), 'progression-history', filters] as const,
  },

  // Workouts
  workouts: {
    all: ['workouts'] as const,
    current: (userId: string) => [...queryKeys.workouts.all, 'current', userId] as const,
    forDay: (userId: string, week: number, dayIndex: number) =>
      [...queryKeys.workouts.all, userId, week, dayIndex] as const,
    sessions: (userId: string) => [...queryKeys.workouts.all, 'sessions', userId] as const,
    sessionList: <T extends object>(userId: string, filters?: T) =>
      [...queryKeys.workouts.sessions(userId), 'list', filters] as const,
    sessionDetail: (userId: string, sessionId: string) =>
      [...queryKeys.workouts.sessions(userId), sessionId] as const,
    sessionCurrent: (userId: string) =>
      [...queryKeys.workouts.sessions(userId), 'current'] as const,
  },
} as const;
