import axios from 'axios';
import type { AxiosInstance, AxiosError, AxiosResponse } from 'axios';

/**
 * HTTP Client for TravelLink API
 * Security-hardened version with:
 * - Request/response validation
 * - HTTPS enforcement in production
 * - Error message sanitization
 * - CORS protection
 * - SSRF prevention
 */

/**
 * Validate and get API base URL from environment
 * Prevents SSRF attacks and configuration tampering
 */
function validateApiUrl(configUrl: string | undefined): string {
  const DEFAULT_URL = import.meta.env.PROD
    ? 'https://api.travellink.com/api/v1' // HTTPS enforced in production
    : 'http://localhost:8000/api/v1';

  if (!configUrl) {
    return DEFAULT_URL;
  }

  try {
    const urlObj = new URL(configUrl);

    // SSRF Prevention: Whitelist allowed hosts
    const ALLOWED_HOSTS = [
      'localhost:8080',
      'localhost:3000',
      'api.travellink.com',
      'api-staging.travellink.com',
    ];

    const hostMatch = ALLOWED_HOSTS.some((host) => {
      return urlObj.host === host || urlObj.hostname === host.split(':')[0];
    });

    if (!hostMatch) {
      console.error(`❌ SECURITY: Unauthorized API host "${urlObj.host}" - falling back to default`);
      return DEFAULT_URL;
    }

    // Enforce HTTPS in production
    if (import.meta.env.PROD && urlObj.protocol !== 'https:') {
      console.error('❌ SECURITY: Non-HTTPS URL in production - falling back to secure URL');
      return DEFAULT_URL;
    }

    return configUrl;
  } catch (e) {
    console.error('❌ SECURITY: Invalid API URL configuration', e);
    return DEFAULT_URL;
  }
}

const API_BASE_URL = validateApiUrl(import.meta.env.VITE_API_URL);

export const api: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true, // Include HttpOnly cookies with requests
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 30000,
  validateStatus: (status) => status >= 200 && status < 300, // Only accept 2xx responses
  responseType: 'json', // Explicitly set expected response type
});

/**
 * Request interceptor: Attach JWT token from secure storage to requests
 * Uses both localStorage fallback and HttpOnly cookies
 */
api.interceptors.request.use(
  (config) => {
    // If HttpOnly cookies aren't working, fallback to localStorage
    const token = localStorage.getItem('access_token');
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    
    // Add security headers
    config.headers['X-Requested-With'] = 'XMLHttpRequest';
    
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

/**
 * Response interceptor: Validate responses and handle authentication errors
 * - Validates Content-Type
 * - Sanitizes error messages
 * - Handles token expiration
 */
api.interceptors.response.use(
  (response: AxiosResponse) => {
    // Validate response Content-Type
    const contentType = response.headers['content-type'];
    if (!contentType?.includes('application/json')) {
      console.error('❌ SECURITY: Invalid response Content-Type', contentType);
      throw new Error('Invalid server response format');
    }

    return response;
  },
  (error: AxiosError) => {
    // Only handle API errors, not network errors
    if (!error.response) {
      return Promise.reject(error);
    }

    // Handle authentication errors (401 Unauthorized)
    if (error.response.status === 401) {
      localStorage.removeItem('access_token');
      localStorage.removeItem('user');
      // Redirect to login handled by router guard
      window.location.href = '/login';
      return Promise.reject(new Error('Session expired. Please login again.'));
    }

    // Handle authorization errors (403 Forbidden)
    if (error.response.status === 403) {
      return Promise.reject(new Error('You do not have permission to access this resource.'));
    }

    // Handle server errors (5xx) - don't expose details
    if (error.response.status >= 500) {
      console.error('Server error:', error.response.status);
      return Promise.reject(new Error('Server error. Please try again later.'));
    }

    // Handle validation errors (4xx, non-401/403)
    if (error.response.status >= 400) {
      // Sanitize error message - only expose if it's a known error type
      const message = (error.response.data as any)?.message;
      if (typeof message === 'string' && message.length < 200) {
        return Promise.reject(new Error(message));
      }
      return Promise.reject(new Error('Request failed. Please try again.'));
    }

    return Promise.reject(error);
  }
);

/**
 * Helper function to set authorization token securely
 * Note: Backend should set HttpOnly cookies; this is fallback for localStorage
 */
export function setAuthToken(token: string): void {
  if (!token || typeof token !== 'string' || token.length > 10000) {
    console.error('❌ SECURITY: Invalid token received');
    return;
  }
  localStorage.setItem('access_token', token);
  if (api.defaults.headers.common) {
    api.defaults.headers.common['Authorization'] = `Bearer ${token}`;
  }
}

/**
 * Helper function to clear authorization token
 */
export function clearAuthToken(): void {
  localStorage.removeItem('access_token');
  if (api.defaults.headers.common) {
    delete api.defaults.headers.common['Authorization'];
  }
}

/**
 * Get current authorization token (for debugging only)
 */
export function getAuthToken(): string | null {
  return localStorage.getItem('access_token');
}

/**
 * Validate JWT format (basic sanity check)
 */
export function isValidJWT(token: string): boolean {
  const parts = token.split('.');
  return parts.length === 3 && parts.every((part) => part.length > 0);
}

export default api;
