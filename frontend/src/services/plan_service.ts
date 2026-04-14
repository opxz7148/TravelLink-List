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
  description?: string | null;                 // Plan-specific notes
  estimated_price_cents?: number | null;       // Plan-specific cost
  duration_minutes?: number | null;            // Plan-specific duration
  details?: {
    // For attractions
    name?: string;
    description?: string;
    location?: string;
    category?: string;
    // For transitions
    title?: string;
  };
}

/**
 * Plan-specific details for a node when creating/editing a plan
 * Allows customization of node properties per plan
 */
export interface NodeDetailForPlan {
  node_id: string;
  description?: string; // Plan-specific notes (max 500 chars)
  estimated_price_cents?: number; // Cost in cents (e.g., 1500 = $15.00)
  duration_minutes?: number; // Duration in minutes for this node in this plan
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
    return response.data.data;
  },

  /**
   * Search travel plans by keyword or destination
   */
  async searchPlans(params: SearchPlansRequest): Promise<{
    plans: TravelPlan[];
    total: number;
  }> {
    const response = await api.get('/plans/search', { params });
    // Backend response envelope: { success, api_version, data: { plans, pagination }, timestamp }
    return response.data.data;
  },

  /**
   * Get plan details with linked list of nodes
   */
  async getPlanDetail(planId: string): Promise<PlanDetail> {
    const response = await api.get(`/plans/${planId}`);
    // Backend response envelope: { success, api_version, data: { plan, nodes }, timestamp }
    return {
      ...response.data.data.plan,
      nodes: response.data.data.nodes,
    };
  },

  /**
   * Get user's own travel plans (private - requires auth)
   */
  async getUserPlans(): Promise<TravelPlan[]> {
    const response = await api.get('/users/me/plans');
    // Backend response envelope: { success, api_version, data: [...], timestamp }
    return response.data.data;
  },

  /**
   * Create a new travel plan with nodes and plan-specific details in a single request (traveller only)
   * @param title Plan title
   * @param destination Plan destination
   * @param nodeDetails Array of node details with plan-specific information (description, price, duration)
   * @param status Plan initial status: 'draft' or 'published' (default: 'draft')
   * @returns Complete plan with added nodes
   */
  async createPlanWithNodes(
    title: string,
    destination: string,
    nodeDetails: NodeDetailForPlan[],
    status: 'draft' | 'published' = 'draft'
  ): Promise<PlanDetail> {
    const response = await api.post(`/plans?status=${status}`, {
      title,
      destination,
      nodes: nodeDetails,
    });
    // Backend response envelope: { success, api_version, data: { plan, nodes }, timestamp }
    return {
      ...response.data.data.plan,
      nodes: response.data.data.nodes,
    };
  },

  /**
   * Create a new draft travel plan (traveller only)
   * @deprecated Use createPlanWithNodes instead
   */
  async createDraftPlan(title: string, destination: string): Promise<TravelPlan> {
    const response = await api.post('/plans?status=draft', {
      title,
      destination,
      nodes: [],
    });
    // Backend response envelope: { success, api_version, data: { plan, ... }, timestamp }
    return response.data.data.plan;
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
    // Backend response envelope: { success, api_version, data: { plan, nodes }, timestamp }
    return {
      ...response.data.data.plan,
      nodes: response.data.data.nodes,
    };
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
    // Backend response envelope: { success, api_version, data: { plan, nodes }, timestamp }
    return {
      ...response.data.data.plan,
      nodes: response.data.data.nodes,
    };
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
    // Backend response envelope: { success, api_version, data: { plan, ... }, timestamp }
    return response.data.data.plan;
  },

  /**
   * Update plan metadata (traveller only)
   */
  async updatePlan(
    planId: string,
    data: Partial<{ title: string; destination: string }>
  ): Promise<TravelPlan> {
    const response = await api.patch(`/plans/${planId}`, data);
    // Backend response envelope: { success, api_version, data: { plan, ... }, timestamp }
    return response.data.data.plan;
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
    // Backend response envelope: { success, api_version, data: { plan, ... }, timestamp }
    return response.data.data.plan;
  },
};
