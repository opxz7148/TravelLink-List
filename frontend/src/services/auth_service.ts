import { api } from './api';

/**
 * API Response envelope structure that wraps all backend responses
 */
interface ApiResponse<T> {
  success: boolean;
  api_version: string;
  data: T | null;
  error: {
    code: string;
    message: string;
    details?: unknown;
  } | null;
  timestamp: string;
}

export interface User {
  id: string;
  username: string;
  email: string;
  display_name: string;
  role: 'simple' | 'traveller' | 'admin';
  created_at: string;
  last_login_at?: string;
}

export interface AuthToken {
  access_token: string;
  token_type: 'Bearer';
  expires_in: number;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  username: string;
  password: string;
  display_name?: string;
}

/**
 * Login response data structure
 */
interface LoginResponse {
  user: User;
  access_token: string;
  token_type: 'Bearer';
  expires_in: number;
}

/**
 * Register response data structure
 */
interface RegisterResponse {
  user: User;
  access_token: string;
  token_type: 'Bearer';
  expires_in: number;
}

export const authService = {
  /**
   * Register a new user account
   * POST /api/v1/auth/register
   * @param data - Registration data (email, username, password, optional display_name)
   * @returns User object
   */
  async register(data: RegisterRequest): Promise<User & AuthToken> {
    const response = await api.post<ApiResponse<RegisterResponse>>('/auth/register', {
      email: data.email,
      username: data.username,
      password: data.password,
      display_name: data.display_name || '',
    });

    // Extract data from ResponseEnvelope
    const responseData = response.data.data;
    if (!responseData) {
      throw new Error('Invalid response from server');
    }

    return {
      // User data
      id: responseData.user.id,
      username: responseData.user.username,
      email: responseData.user.email,
      display_name: responseData.user.display_name,
      role: responseData.user.role,
      created_at: responseData.user.created_at,
      // Auth token
      access_token: responseData.access_token,
      token_type: responseData.token_type,
      expires_in: responseData.expires_in,
    };
  },

  /**
   * Login with email and password
   * POST /api/v1/auth/login
   * @param data - Login credentials (email, password)
   * @returns User object and auth token
   */
  async login(data: LoginRequest): Promise<{ user: User; token: AuthToken }> {
    const response = await api.post<ApiResponse<LoginResponse>>('/auth/login', data);

    // Extract data from ResponseEnvelope
    const responseData = response.data.data;
    if (!responseData) {
      throw new Error('Invalid response from server');
    }

    return {
      user: responseData.user,
      token: {
        access_token: responseData.access_token,
        token_type: responseData.token_type,
        expires_in: responseData.expires_in,
      },
    };
  },

  /**
   * Logout user
   * POST /api/v1/auth/logout
   * Note: Backend returns 204 No Content, so no data extraction needed
   */
  async logout(): Promise<void> {
    await api.post('/auth/logout');
    // Client-side logout handled by store
  },

  /**
   * Get current user profile (requires auth)
   * GET /api/v1/users/me or /api/v1/users/:id
   */
  async getCurrentUser(): Promise<User> {
    const response = await api.get<ApiResponse<{ user: User }>>('/users/me');
    const responseData = response.data.data;
    if (!responseData) {
      throw new Error('Invalid response from server');
    }
    return responseData.user;
  },

  /**
   * Update user profile (requires auth)
   */
  async updateProfile(
    data: Partial<{ username: string; email: string; display_name: string }>
  ): Promise<User> {
    const response = await api.put<ApiResponse<{ user: User }>>('/users/me', data);
    const responseData = response.data.data;
    if (!responseData) {
      throw new Error('Invalid response from server');
    }
    return responseData.user;
  },

  /**
   * Change user password (requires auth)
   */
  async changePassword(oldPassword: string, newPassword: string): Promise<void> {
    await api.post('/users/me/change-password', {
      old_password: oldPassword,
      new_password: newPassword,
    });
  },

  /**
   * Validate if token is still valid
   */
  async validateToken(token: string): Promise<boolean> {
    try {
      await api.get('/auth/validate', {
        headers: { Authorization: `Bearer ${token}` },
      });
      return true;
    } catch {
      return false;
    }
  },
};
