import { api } from './api';

export interface Comment {
  id: string;
  plan_id: string;
  author_id: string;
  author?: {
    username: string;
  };
  text: string;
  created_at: string;
  updated_at?: string;
}

export const commentService = {
  /**
   * Get all comments for a plan
   */
  async getComments(planId: string, page?: number, limit?: number): Promise<{
    comments: Comment[];
    total: number;
  }> {
    const response = await api.get(`/plans/${planId}/comments`, {
      params: { page, limit },
    });
    return response.data;
  },

  /**
   * Post a new comment (requires auth)
   */
  async createComment(planId: string, text: string): Promise<Comment> {
    const response = await api.post(`/plans/${planId}/comments`, {
      text,
    });
    return response.data;
  },

  /**
   * Update a comment (own comments only)
   */
  async updateComment(commentId: string, text: string): Promise<Comment> {
    const response = await api.put(`/comments/${commentId}`, {
      text,
    });
    return response.data;
  },

  /**
   * Delete a comment (own comments or admin)
   */
  async deleteComment(commentId: string): Promise<void> {
    await api.delete(`/comments/${commentId}`);
  },
};
