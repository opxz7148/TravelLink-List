import { api } from './api';
import type { TravelPlan } from './plan_service';

export interface UserProfile {
  id: string;
  username: string;
  email: string;
  role: 'simple' | 'traveller' | 'admin';
  is_active: boolean;
  created_at: string;
  plan_count?: number;
}

export const userService = {
  /**
   * Get user profile (requires auth)
   */
  async getProfile(): Promise<UserProfile> {
    const response = await api.get('/users/me');
    return response.data;
  },

  /**
   * Get user's travel plans (requires auth)
   */
  async getMyPlans(page?: number, status?: string): Promise<{
    plans: TravelPlan[];
    total: number;
  }> {
    const response = await api.get('/users/me/plans', {
      params: { page, status },
    });
    return response.data;
  },

  /**
   * Get public user profile by username
   */
  async getUserProfile(username: string): Promise<UserProfile> {
    const response = await api.get(`/users/${username}`);
    return response.data;
  },

  /**
   * Get another user's published plans
   */
  async getUserPlans(username: string, page?: number): Promise<{
    plans: TravelPlan[];
    total: number;
  }> {
    const response = await api.get(`/users/${username}/plans`, {
      params: { page },
    });
    return response.data;
  },

  /**
   * Update own profile (requires auth)
   */
  async updateProfile(data: Partial<{
    username: string;
    email: string;
  }>): Promise<UserProfile> {
    const response = await api.patch('/users/me', data);
    return response.data;
  },

  /**
   * Change password (requires auth)
   */
  async changePassword(oldPassword: string, newPassword: string): Promise<void> {
    await api.patch('/users/me/password', {
      old_password: oldPassword,
      new_password: newPassword,
    });
  },

  /**
   * Delete own account (requires auth)
   */
  async deleteAccount(password: string): Promise<void> {
    await api.delete('/users/me', {
      data: { password },
    });
  },
};
