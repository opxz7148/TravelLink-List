import { api } from './api';

export interface AttractionNode {
  id: string;
  type: 'attraction';
  name: string;
  description: string;
  location: string;
  is_approved: boolean;
  created_by: string;
  created_at: string;
}

export interface TransitionNode {
  id: string;
  type: 'transition';
  duration_minutes: number;
  description: string;
  is_approved: boolean;
  created_by: string;
  created_at: string;
}

export type Node = AttractionNode | TransitionNode;

export interface ListNodesRequest {
  type?: 'attraction' | 'transition';
  search?: string;
  page?: number;
  limit?: number;
  approved_only?: boolean;
}

export const nodeService = {
  /**
   * List approved nodes (public)
   */
  async listApprovedNodes(params?: ListNodesRequest): Promise<{
    nodes: Node[];
    total: number;
  }> {
    const response = await api.get('/nodes', {
      params: { ...params, approved_only: true },
    });
    return response.data;
  },

  /**
   * List nodes by type (approved only)
   */
  async listByType(
    type: 'attraction' | 'transition',
    params?: Omit<ListNodesRequest, 'type'>
  ): Promise<{
    nodes: Node[];
    total: number;
  }> {
    const response = await api.get(`/nodes/${type}s`, {
      params: { ...params, approved_only: true },
    });
    return response.data;
  },

  /**
   * Search attractions by name or location
   */
  async searchAttractions(query: string, page?: number): Promise<{
    nodes: AttractionNode[];
    total: number;
  }> {
    const response = await api.get('/nodes/attractions/search', {
      params: { q: query, page },
    });
    return response.data;
  },

  /**
   * Get node detail
   */
  async getNodeDetail(nodeId: string): Promise<Node> {
    const response = await api.get(`/nodes/${nodeId}`);
    return response.data;
  },

  /**
   * Create a new attraction node (traveller only)
   */
  async createAttractionNode(data: {
    name: string;
    description: string;
    location: string;
  }): Promise<AttractionNode> {
    const response = await api.post('/nodes/attractions', data);
    return response.data;
  },

  /**
   * Create a new transition node (traveller only)
   */
  async createTransitionNode(data: {
    description: string;
    duration_minutes: number;
  }): Promise<TransitionNode> {
    const response = await api.post('/nodes/transitions', data);
    return response.data;
  },

  /**
   * Get user's own nodes (including unapproved)
   */
  async getUserNodes(): Promise<Node[]> {
    const response = await api.get('/nodes/my-nodes');
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
   * Delete a node (admin only)
   */
  async deleteNode(nodeId: string): Promise<void> {
    await api.delete(`/admin/nodes/${nodeId}`);
  },
};
