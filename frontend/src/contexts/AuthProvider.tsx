import { useEffect, useState, useCallback, useMemo } from 'react';
import { configureClient } from '../api/client';
import { onUnauthorized, clearUnauthorizedHandler } from '../api/auth-error-handler';
import { queryClient } from '../lib/query';
import * as authApi from '../api/endpoints/auth';
import { AUTH_STORAGE_KEY, TOKEN_STORAGE_KEY, type AuthContextValue, type AuthProviderProps } from './auth-types';
import { AuthContext } from './auth-context';

function generateUserId(): string {
  return crypto.randomUUID();
}

function getStoredUserId(): string | null {
  return localStorage.getItem(AUTH_STORAGE_KEY);
}

function getStoredToken(): string | null {
  return localStorage.getItem(TOKEN_STORAGE_KEY);
}

// Get initial values for state (client.ts handles its own initialization)
const initialToken = getStoredToken();
const initialUserId = getStoredUserId();

export function AuthProvider({ children }: AuthProviderProps) {
  const [userId, setUserId] = useState<string | null>(initialUserId);
  const [email, setEmail] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(!!initialToken); // Only loading if we need to validate token

  // Validate stored token on mount (client already configured synchronously above)
  useEffect(() => {
    async function validateToken() {
      if (!initialToken) {
        setIsLoading(false);
        return;
      }

      try {
        // Validate token by fetching user info
        const user = await authApi.me();
        setUserId(user.id);
        setEmail(user.email);
        localStorage.setItem(AUTH_STORAGE_KEY, user.id);
      } catch {
        // Token is invalid, clear it
        localStorage.removeItem(TOKEN_STORAGE_KEY);
        localStorage.removeItem(AUTH_STORAGE_KEY);
        configureClient({ token: undefined, userId: undefined });
        setUserId(null);
      }

      setIsLoading(false);
    }

    validateToken();
  }, []);

  const loginWithId = useCallback((newUserId: string) => {
    localStorage.setItem(AUTH_STORAGE_KEY, newUserId);
    localStorage.removeItem(TOKEN_STORAGE_KEY); // Clear any existing token
    setUserId(newUserId);
    setEmail(null);
    configureClient({ userId: newUserId, token: undefined });
  }, []);

  const loginWithCredentials = useCallback(async (emailInput: string, password: string) => {
    const response = await authApi.login({ email: emailInput, password });

    // Store token and userId
    localStorage.setItem(TOKEN_STORAGE_KEY, response.token);
    localStorage.setItem(AUTH_STORAGE_KEY, response.user.id);

    // Update state
    setUserId(response.user.id);
    setEmail(response.user.email);

    // Configure client with both userId and token
    configureClient({ userId: response.user.id, token: response.token });
  }, []);

  const registerUser = useCallback(async (emailInput: string, password: string, name?: string) => {
    // First register the user
    await authApi.register({ email: emailInput, password, name });

    // Then login to get the token
    await loginWithCredentials(emailInput, password);
  }, [loginWithCredentials]);

  const logout = useCallback(async () => {
    // Try to call logout API (ignore errors if token is already invalid)
    try {
      await authApi.logout();
    } catch {
      // Ignore logout API errors
    }

    // Clear storage
    localStorage.removeItem(AUTH_STORAGE_KEY);
    localStorage.removeItem(TOKEN_STORAGE_KEY);

    // Clear client config
    configureClient({ userId: undefined, token: undefined });

    // Clear query cache
    queryClient.clear();

    // Reset state
    setUserId(null);
    setEmail(null);
  }, []);

  // Register the unauthorized handler to trigger logout on 401 responses
  useEffect(() => {
    onUnauthorized(logout);
    return () => {
      clearUnauthorizedHandler();
    };
  }, [logout]);

  const createUser = useCallback(async (): Promise<string> => {
    // Generate random credentials for the new user
    const newUserId = generateUserId();
    const generatedEmail = `${newUserId}@powerpro.local`;
    const generatedPassword = crypto.randomUUID();

    try {
      // Try to register with backend
      await registerUser(generatedEmail, generatedPassword);
      return userId!; // userId is set by registerUser
    } catch {
      // If backend registration fails, fall back to dev mode (local-only)
      loginWithId(newUserId);
      return newUserId;
    }
  }, [registerUser, loginWithId, userId]);

  const value: AuthContextValue = useMemo(
    () => ({
      userId,
      email,
      isAuthenticated: userId !== null,
      isLoading,
      loginWithId,
      loginWithCredentials,
      registerUser,
      logout,
      createUser,
    }),
    [userId, email, isLoading, loginWithId, loginWithCredentials, registerUser, logout, createUser]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
