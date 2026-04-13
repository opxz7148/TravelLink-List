import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { authService } from '../services/auth_service';
import type { User, AuthToken } from '../services/auth_service';

/**
 * Sanitize error messages to prevent information disclosure
 */
function sanitizeErrorMessage(error: any): string {
  const status = error?.response?.status;
  const message = error?.response?.data?.message;

  // Map specific HTTP errors to safe messages
  const errorMap: Record<number, string> = {
    400: 'Invalid request. Please check your input.',
    401: 'Authentication failed. Invalid email or password.',
    403: 'Permission denied.',
    404: 'Resource not found.',
    409: 'Username or email already in use.',
    422: 'Validation error. Please check your input.',
    500: 'Server error. Please try again later.',
    502: 'Service unavailable. Please try again later.',
    503: 'Service unavailable. Please try again later.',
  };

  // If we have a mapped status code, use it
  if (status && status in errorMap) {
    return errorMap[status as keyof typeof errorMap] || 'An error occurred. Please try again later.';
  }

  // For password mismatch, provide specific message
  if (status === 400 && message?.includes('password')) {
    return 'Passwords do not match or are invalid.';
  }

  // For user not found, use generic message (prevent user enumeration)
  if (status === 401) {
    return 'Authentication failed. Please try again.';
  }

  // Default fallback
  return 'An error occurred. Please try again later.';
}

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref<User | null>(null);
  const token = ref<string>('');
  const isLoading = ref(false);
  const error = ref<string>('');

  // Computed
  const isAuthenticated = computed(() => !!token.value);
  const isTraveller = computed(() => user.value?.role === 'traveller' || user.value?.role === 'admin');
  const isAdmin = computed(() => user.value?.role === 'admin');
  const userRole = computed(() => user.value?.role || 'simple');

  // Actions
  const register = async (
    username: string,
    email: string,
    password: string,
    displayName: string = ''
  ) => {
    isLoading.value = true;
    error.value = '';
    try {
      const response = await authService.register({
        email,
        username,
        password,
        display_name: displayName,
      });

      // Extract user and token from response
      const userData = {
        id: response.id,
        username: response.username,
        email: response.email,
        display_name: response.display_name,
        role: response.role,
        created_at: response.created_at,
      };

      user.value = userData;
      token.value = response.access_token;

      // Store in localStorage for persistence
      localStorage.setItem('access_token', response.access_token);
      localStorage.setItem('user', JSON.stringify(userData));

      return userData;
    } catch (err: any) {
      error.value = sanitizeErrorMessage(err);
      throw err;
    } finally {
      isLoading.value = false;
    }
  };

  const login = async (email: string, password: string) => {
    isLoading.value = true;
    error.value = '';
    console.log('🔐 Attempting login for:', email);
    try {
      console.log('📡 Sending login request to backend...');
      const response = await authService.login({ email, password });
      console.log('✅ Login successful:', response);
      user.value = response.user;
      token.value = response.token.access_token;

      // Store in localStorage for persistence (with validation)
      if (response.token.access_token && response.token.access_token.length > 0) {
        localStorage.setItem('access_token', response.token.access_token);
        localStorage.setItem('user', JSON.stringify(response.user));
      } else {
        throw new Error('Invalid token received');
      }

      return response;
    } catch (err: any) {
      error.value = sanitizeErrorMessage(err);
      throw err;
    } finally {
      isLoading.value = false;
    }
  };

  const logout = () => {
    user.value = null;
    token.value = '';
    error.value = '';
    localStorage.removeItem('access_token');
    localStorage.removeItem('user');
  };

  const restoreSession = () => {
    const storedToken = localStorage.getItem('access_token');
    const storedUser = localStorage.getItem('user');

    if (storedToken && storedUser) {
      try {
        // Validate token format
        const parts = storedToken.split('.');
        if (parts.length !== 3) {
          throw new Error('Invalid token format');
        }

        token.value = storedToken as string;
        user.value = JSON.parse(storedUser);
      } catch (err) {
        console.error('❌ SECURITY: Invalid stored session');
        logout();
      }
    }
  };

  const updateProfile = async (data: Partial<{ username: string; email: string }>) => {
    isLoading.value = true;
    error.value = '';
    try {
      const response = await authService.updateProfile(data);
      user.value = response;
      localStorage.setItem('user', JSON.stringify(response));
      return response;
    } catch (err: any) {
      error.value = sanitizeErrorMessage(err);
      throw err;
    } finally {
      isLoading.value = false;
    }
  };

  const changePassword = async (oldPassword: string, newPassword: string) => {
    isLoading.value = true;
    error.value = '';
    try {
      await authService.changePassword(oldPassword, newPassword);
      error.value = '';
    } catch (err: any) {
      error.value = sanitizeErrorMessage(err);
      throw err;
    } finally {
      isLoading.value = false;
    }
  };

  const getCurrentUser = async () => {
    isLoading.value = true;
    error.value = '';
    try {
      const response = await authService.getCurrentUser();
      user.value = response;
      localStorage.setItem('user', JSON.stringify(response));
      return response;
    } catch (err: any) {
      error.value = sanitizeErrorMessage(err);
      throw err;
    } finally {
      isLoading.value = false;
    }
  };

  return {
    // State
    user,
    token,
    isLoading,
    error,

    // Computed
    isAuthenticated,
    isTraveller,
    isAdmin,
    userRole,

    // Actions
    register,
    login,
    logout,
    restoreSession,
    updateProfile,
    changePassword,
    getCurrentUser,
  };
});
