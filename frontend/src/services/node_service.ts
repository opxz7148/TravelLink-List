import { api } from './api';

/**
 * Embedded detail for an attraction node
 * Non-null for nodes with type='attraction'
 */
export interface AttractionDetail {
  node_id: string;
  name: string;
  category: string;
  location: string;
  description?: string;
  contact_info?: string;
  hours_of_operation?: string;
  estimated_visit_duration_minutes?: number;
  created_at: string;
}

/**
 * Embedded detail for a transition node
 * Non-null for nodes with type='transition'
 */
export interface TransitionDetail {
  node_id: string;
  title: string;
  mode: string;
  description?: string;
  hours_of_operation?: string;
  route_notes?: string;
  created_at: string;
}

/**
 * Complete Node with embedded type-specific details
 * Either attraction or transition detail will be populated based on type
 */
export interface Node {
  id: string;
  type: 'attraction' | 'transition';
  created_by: string;
  is_approved: boolean;
  created_at: string;
  updated_at?: string;
  
  // Embedded details - only one will be populated based on type
  attraction?: AttractionDetail | null;
  transition?: TransitionDetail | null;
}

/**
 * Helper type for working with attraction nodes with populated details
 */
export interface AttractionNode extends Node {
  type: 'attraction';
  attraction: AttractionDetail;
}

/**
 * Helper type for working with transition nodes with populated details
 */
export interface TransitionNode extends Node {
  type: 'transition';
  transition: TransitionDetail;
}

export interface ListNodesRequest {
  type?: 'attraction' | 'transition';
  search?: string;
  page?: number;
  limit?: number;
  approved_only?: boolean;
}

/**
 * Helper functions to extract detail information from nodes
 */
export function getAttractionDetail(node: Node): AttractionDetail | null {
  return node.type === 'attraction' && node.attraction ? node.attraction : null;
}

export function getTransitionDetail(node: Node): TransitionDetail | null {
  return node.type === 'transition' && node.transition ? node.transition : null;
}

/**
 * Get display name for a node (name for attractions, title for transitions)
 */
export function getNodeName(node: Node): string {
  if (node.type === 'attraction' && node.attraction) {
    return node.attraction.name;
  } else if (node.type === 'transition' && node.transition) {
    return node.transition.title;
  }
  return 'Unknown Node';
}

/**
 * Get display description for a node
 */
export function getNodeDescription(node: Node): string {
  if (node.type === 'attraction' && node.attraction) {
    return node.attraction.description || '';
  } else if (node.type === 'transition' && node.transition) {
    return node.transition.description || '';
  }
  return '';
}

export const nodeService = {
  /**
   * List approved nodes with embedded details (public)
   */
  async listApprovedNodes(params?: ListNodesRequest): Promise<{
    nodes: Node[];
    total: number;
  }> {
    const response = await api.get('/nodes', {
      params: { ...params, approved_only: true },
    });
    // Response includes nodes with embedded details
    console.log('API response for listApprovedNodes:', response.data.data);
    return response.data.data;
  },

  /**
   * List nodes by type with embedded details (approved only)
   */
  async listByType(
    type: 'attraction' | 'transition',
    params?: Omit<ListNodesRequest, 'type'>
  ): Promise<{
    nodes: Node[];
    total: number;
  }> {
    const response = await api.get('/nodes', {
      params: { ...params, type, approved_only: true },
    });
    return response.data;
  },

  /**
   * Search attractions by name
   */
  async searchAttractions(query: string, page?: number): Promise<{
    nodes: Node[];
    total: number;
  }> {
    const response = await api.get('/nodes', {
      params: { search: query, type: 'attraction', page, approved_only: true },
    });
    return response.data;
  },

  /**
   * Get node detail with embedded type-specific information
   */
  async getNodeDetail(nodeId: string): Promise<Node> {
    const response = await api.get(`/nodes/${nodeId}`);
    // Response includes complete node with embedded details (attraction or transition)
    return response.data.node || response.data;
  },

  /**
   * Create a new attraction node (traveller only)
   */
  async createAttractionNode(data: {
    name: string;
    description?: string;
    location: string;
    category?: string;
    contact_info?: string;
    hours_of_operation?: string;
  }): Promise<Node> {
    const response = await api.post('/nodes/attraction', data);
    return response.data.node || response.data;
  },

  /**
   * Create a new transition node (traveller only)
   */
  async createTransitionNode(data: {
    title: string;
    mode: string;
    description?: string;
    hours_of_operation?: string;
    route_notes?: string;
  }): Promise<Node> {
    const response = await api.post('/nodes/transition', data);
    return response.data.node || response.data;
  },

  /**
   * Get user's own draft nodes with embedded details (unapproved only)
   */
  async getUserNodes(): Promise<Node[]> {
    const response = await api.get('/nodes/my-draft');
    const nodes = response.data.data.nodes || [];
    return Array.isArray(nodes) ? nodes : [];
  },

  /**
   * Approve a user-created node (admin only)
   */
  async approveNode(nodeId: string): Promise<Node> {
    const response = await api.patch(`/admin/nodes/${nodeId}/approve`);
    return response.data.node || response.data;
  },

  /**
   * Delete a node (admin only)
   */
  async deleteNode(nodeId: string): Promise<void> {
    await api.delete(`/admin/nodes/${nodeId}`);
  },
};
