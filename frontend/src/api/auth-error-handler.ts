type AuthErrorCallback = () => void;

let onUnauthorizedCallback: AuthErrorCallback | null = null;

/**
 * Registers a callback to be invoked when a 401 Unauthorized response is received.
 * This allows the AuthProvider to handle logout when auth errors occur.
 * @param callback - Function to call on 401 responses (typically triggers logout)
 */
export function onUnauthorized(callback: AuthErrorCallback): void {
  onUnauthorizedCallback = callback;
}

/**
 * Triggers the registered unauthorized callback.
 * Called by the API client when a 401 response is received.
 */
export function triggerUnauthorized(): void {
  if (onUnauthorizedCallback) {
    onUnauthorizedCallback();
  }
}

/**
 * Clears the registered unauthorized callback.
 * Should be called when the AuthProvider unmounts.
 */
export function clearUnauthorizedHandler(): void {
  onUnauthorizedCallback = null;
}
