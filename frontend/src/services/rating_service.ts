import { api } from './api';

export interface Rating {
  id: string;
  plan_id: string;
  user_id: string;
  stars: number;
  created_at: string;
}

export const ratingService = {
  /**
   * Get rating statistics for a plan
   */
  async getPlanRatingStats(planId: string): Promise<{
    average: number;
    count: number;
  }> {
    const response = await api.get(`/plans/${planId}/ratings/stats`);
    return response.data;
  },

  /**
   * Get user's rating for a plan (if exists)
   */
  async getUserRating(planId: string): Promise<Rating | null> {
    try {
      const response = await api.get(`/plans/${planId}/ratings/my-rating`);
      return response.data;
    } catch (error: any) {
      // 404 means user hasn't rated this plan
      if (error.response?.status === 404) {
        return null;
      }
      throw error;
    }
  },

  /**
   * Submit or update a rating (requires auth)
   */
  async submitRating(planId: string, stars: number): Promise<Rating> {
    if (stars < 1 || stars > 5 || !Number.isInteger(stars)) {
      throw new Error('Rating must be between 1 and 5 stars');
    }

    const response = await api.post(`/plans/${planId}/ratings`, {
      stars,
    });
    return response.data;
  },

  /**
   * Remove a rating (own ratings only)
   */
  async deleteRating(planId: string): Promise<void> {
    await api.delete(`/plans/${planId}/ratings`);
  },
};
