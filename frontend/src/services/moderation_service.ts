import { api } from './api';
import type { TravelPlan } from './plan_service';
import type { Node, AttractionNode, TransitionNode } from './node_service';

export interface FlaggedPlan extends TravelPlan {
  flag_reason?: string;
  flag_count: number;
}

export type PendingNode = (AttractionNode | TransitionNode) & {
  created_by_username: string;
};

export const moderationService = {
  /**
   * Get flagged travel plans requiring review (admin only)
   */
  async getFlaggedPlans(page?: number): Promise<{
    plans: FlaggedPlan[];
    total: number;
  }> {
    const response = await api.get('/admin/plans/flagged', {
      params: { page },
    });
    return response.data;
  },

  /**
   * Get pending user-created nodes for approval (admin only)
   */
  async getPendingNodes(page?: number): Promise<{
    nodes: PendingNode[];
    total: number;
  }> {
    const response = await api.get('/admin/nodes/pending', {
      params: { page },
    });
    return response.data;
  },

  /**
   * Delete a flagged plan (admin only)
   */
  async deletePlan(planId: string, reason?: string): Promise<void> {
    await api.delete(`/admin/plans/${planId}`, {
      data: { reason },
    });
  },

  /**
   * Suspend a plan (admin only)
   */
  async suspendPlan(planId: string): Promise<TravelPlan> {
    const response = await api.patch(`/admin/plans/${planId}/suspend`);
    return response.data;
  },

  /**
   * Approve a user-created node (admin only)
   */
  async approveNode(nodeId: string): Promise<Node> {
    const response = await api.patch(`/admin/nodes/${nodeId}/approve`);
    return response.data;
  },

  /**
   * Reject/delete a user-created node (admin only)
   */
  async rejectNode(nodeId: string, reason?: string): Promise<void> {
    await api.delete(`/admin/nodes/${nodeId}`, {
      data: { reason },
    });
  },

  /**
   * Get user management info for admins
   */
  async getUsers(page?: number, role?: string): Promise<{
    users: any[];
    total: number;
  }> {
    const response = await api.get('/admin/users', {
      params: { page, role },
    });
    return response.data;
  },

  /**
   * Update user role (admin only)
   */
  async updateUserRole(userId: string, newRole: string): Promise<any> {
    const response = await api.patch(`/admin/users/${userId}/role`, {
      role: newRole,
    });
    return response.data;
  },

  /**
   * Deactivate a user account (admin only)
   */
  async deactivateUser(userId: string): Promise<void> {
    await api.patch(`/admin/users/${userId}/deactivate`);
  },
};
