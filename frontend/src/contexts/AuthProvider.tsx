import { useEffect, useState, useCallback, useMemo } from 'react';
import { configureClient } from '../api/client';
import { AUTH_STORAGE_KEY, type AuthContextValue, type AuthProviderProps } from './auth-types';
import { AuthContext } from './auth-context';

function generateUserId(): string {
  return crypto.randomUUID();
}

function getInitialUserId(): string | null {
  return localStorage.getItem(AUTH_STORAGE_KEY);
}

export function AuthProvider({ children }: AuthProviderProps) {
  const [userId, setUserId] = useState<string | null>(getInitialUserId);

  // Configure API client when userId changes
  useEffect(() => {
    configureClient({ userId: userId ?? undefined });
  }, [userId]);

  const login = useCallback((newUserId: string) => {
    localStorage.setItem(AUTH_STORAGE_KEY, newUserId);
    setUserId(newUserId);
  }, []);

  const logout = useCallback(() => {
    localStorage.removeItem(AUTH_STORAGE_KEY);
    setUserId(null);
  }, []);

  const createUser = useCallback((): string => {
    const newUserId = generateUserId();
    login(newUserId);
    return newUserId;
  }, [login]);

  const value: AuthContextValue = useMemo(
    () => ({
      userId,
      isAuthenticated: userId !== null,
      isLoading: false,
      login,
      logout,
      createUser,
    }),
    [userId, login, logout, createUser]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
