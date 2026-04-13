import { api } from './api';

export interface PromotionRequest {
  id: string;
  user_id: string;
  status: 'pending' | 'approved' | 'rejected';
  created_at: string;
  approved_at?: string;
  rejection_reason?: string;
}

export const promotionService = {
  /**
   * Submit a promotion request (simple user -> traveller)
   */
  async submitPromotionRequest(
    careerHighlight: string
  ): Promise<PromotionRequest> {
    const response = await api.post('/promotions/request', {
      career_highlight: careerHighlight,
    });
    return response.data;
  },

  /**
   * Get current user's promotion status
   */
  async getUserPromotionStatus(): Promise<PromotionRequest | null> {
    try {
      const response = await api.get('/promotions/my-status');
      return response.data;
    } catch (error: any) {
      if (error.response?.status === 404) {
        return null;
      }
      throw error;
    }
  },

  /**
   * Get all pending promotion requests (admin only)
   */
  async getPendingRequests(page?: number): Promise<{
    requests: PromotionRequest[];
    total: number;
  }> {
    const response = await api.get('/admin/promotions/pending', {
      params: { page },
    });
    return response.data;
  },

  /**
   * Approve a promotion request (admin only)
   */
  async approvePromotion(requestId: string): Promise<PromotionRequest> {
    const response = await api.patch(
      `/admin/promotions/${requestId}/approve`
    );
    return response.data;
  },

  /**
   * Reject a promotion request (admin only)
   */
  async rejectPromotion(
    requestId: string,
    reason?: string
  ): Promise<PromotionRequest> {
    const response = await api.patch(
      `/admin/promotions/${requestId}/reject`,
      { reason }
    );
    return response.data;
  },
};
