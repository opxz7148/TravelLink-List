import { api } from './api';

export interface TravelPlan {
  id: string;
  title: string;
  destination: string;
  author_id: string;
  author?: {
    username: string;
  };
  status: 'draft' | 'published' | 'archived';
  rating_count: number;
  rating_sum: number;
  comment_count: number;
  created_at: string;
  updated_at: string;
}

export interface PlanDetail extends TravelPlan {
  nodes?: PlanNode[];
}

export interface PlanNode {
  id: string;
  type: 'attraction' | 'transition';
  sequence_position: number;
  details?: {
    name?: string;
    description?: string;
    location?: string;
    duration_minutes?: number;
  };
}

export interface ListPlansRequest {
  page?: number;
  limit?: number;
  status?: string;
  sort?: 'recent' | 'popular' | 'rated';
}

export interface SearchPlansRequest {
  q: string;
  destination?: string;
  page?: number;
  limit?: number;
}

export const planService = {
  /**
   * List published travel plans with pagination
   */
  async listPlans(params?: ListPlansRequest): Promise<{
    plans: TravelPlan[];
    total: number;
    page: number;
  }> {
    const response = await api.get('/plans', { params });
    return response.data;
  },

  /**
   * Search travel plans by keyword or destination
   */
  async searchPlans(params: SearchPlansRequest): Promise<{
    plans: TravelPlan[];
    total: number;
  }> {
    const response = await api.get('/plans/search', { params });
    return response.data;
  },

  /**
   * Get plan details with linked list of nodes
   */
  async getPlanDetail(planId: string): Promise<PlanDetail> {
    const response = await api.get(`/plans/${planId}`);
    return response.data;
  },

  /**
   * Get user's own travel plans (private - requires auth)
   */
  async getUserPlans(): Promise<TravelPlan[]> {
    const response = await api.get('/users/me/plans');
    return response.data;
  },

  /**
   * Create a new draft travel plan (traveller only)
   */
  async createDraftPlan(title: string, destination: string): Promise<TravelPlan> {
    const response = await api.post('/plans', {
      title,
      destination,
    });
    return response.data;
  },

  /**
   * Add nodes to plan in sequence (traveller only)
   */
  async addNodesToPlan(
    planId: string,
    nodeIds: string[]
  ): Promise<PlanDetail> {
    const response = await api.patch(`/plans/${planId}/nodes`, {
      node_ids: nodeIds,
    });
    return response.data;
  },

  /**
   * Reorder nodes in plan (traveller only)
   */
  async reorderNodes(
    planId: string,
    orderedNodeIds: string[]
  ): Promise<PlanDetail> {
    const response = await api.patch(`/plans/${planId}/nodes/reorder`, {
      node_ids: orderedNodeIds,
    });
    return response.data;
  },

  /**
   * Remove node from plan (traveller only)
   */
  async removeNodeFromPlan(planId: string, nodeId: string): Promise<void> {
    await api.delete(`/plans/${planId}/nodes/${nodeId}`);
  },

  /**
   * Publish a draft plan (traveller only)
   */
  async publishPlan(planId: string): Promise<TravelPlan> {
    const response = await api.patch(`/plans/${planId}/publish`);
    return response.data;
  },

  /**
   * Update plan metadata (traveller only)
   */
  async updatePlan(
    planId: string,
    data: Partial<{ title: string; destination: string }>
  ): Promise<TravelPlan> {
    const response = await api.patch(`/plans/${planId}`, data);
    return response.data;
  },

  /**
   * Delete a plan (traveller/admin)
   */
  async deletePlan(planId: string): Promise<void> {
    await api.delete(`/plans/${planId}`);
  },

  /**
   * Suspend a plan (admin only)
   */
  async suspendPlan(planId: string): Promise<TravelPlan> {
    const response = await api.patch(`/admin/plans/${planId}/suspend`);
    return response.data;
  },
};
