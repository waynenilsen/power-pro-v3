import type { ReactNode } from 'react';

export const AUTH_STORAGE_KEY = 'powerpro_user_id';

export interface AuthState {
  userId: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
}

export interface AuthActions {
  login: (userId: string) => void;
  logout: () => void;
  createUser: () => string;
}

export type AuthContextValue = AuthState & AuthActions;

export interface AuthProviderProps {
  children: ReactNode;
}
