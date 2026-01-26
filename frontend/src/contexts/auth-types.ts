import type { ReactNode } from 'react';

export const AUTH_STORAGE_KEY = 'powerpro_user_id';
export const TOKEN_STORAGE_KEY = 'powerpro_token';

export interface AuthState {
  userId: string | null;
  email: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
}

export interface AuthActions {
  /** Login with user ID directly (for dev mode backward compat) */
  loginWithId: (userId: string) => void;
  /** Login with email and password credentials */
  loginWithCredentials: (email: string, password: string) => Promise<void>;
  /** Register a new user with email and password */
  registerUser: (email: string, password: string, name?: string) => Promise<void>;
  logout: () => void;
  /** Create a user with generated credentials (calls registerUser internally) */
  createUser: () => Promise<string>;
}

export type AuthContextValue = AuthState & AuthActions;

export interface AuthProviderProps {
  children: ReactNode;
}
