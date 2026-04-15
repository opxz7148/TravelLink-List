import { api } from './api';
import type { TravelPlan } from './plan_service';
import type { UserProfile } from './user_service';

export interface PromotionRequest {
  id: string;
  user_id: string;
  user?: UserProfile;
  plan_id?: string;
  plan?: TravelPlan;
  status: 'pending' | 'approved' | 'rejected';
  admin_notes: string;
  created_at: string;
  reviewed_at?: string;
}

export const promotionService = {
  /**
   * Submit a promotion request (simple user -> traveller or plan promotion)
   */
  async submitPromotionRequest(
    careerHighlight: string,
    planId?: string
  ): Promise<PromotionRequest> {
    const payload: any = planId 
      ? { plan_id: planId }
      : { career_highlight: careerHighlight };
    
    const response = await api.post('/promotions/request', payload);
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
    pagination: {
      current_page: number;
      limit: number;
      total: number;
    };
  }> {
    const response = await api.get('/admin/promotions/pending', {
      params: { page },
    });
    return response.data.data;
  },

  /**
   * Approve a promotion request (admin only)
   */
  async approvePromotion(
    requestId: string,
    adminNotes?: string
  ): Promise<{ request_id: string; status: string }> {
    const response = await api.patch(
      `/admin/promotions/${requestId}/approve`,
      { admin_notes: adminNotes || '' }
    );
    return response.data;
  },

  /**
   * Reject a promotion request (admin only)
   */
  async rejectPromotion(
    requestId: string,
    adminNotes?: string
  ): Promise<{ request_id: string; status: string }> {
    const response = await api.patch(
      `/admin/promotions/${requestId}/reject`,
      { admin_notes: adminNotes || '' }
    );
    return response.data;
  },
};
