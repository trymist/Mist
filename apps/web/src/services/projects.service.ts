import type { Project, CreateProjectRequest, UpdateProjectRequest } from '@/types';

const API_BASE = '/api';

export const projectsService = {
  /**
   * Get all projects
   */
  async getAll(): Promise<Project[]> {
    const response = await fetch(`${API_BASE}/projects/getAll`, {
      credentials: 'include',
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to fetch projects');
    }

    return data.data || [];
  },

  /**
   * Get project by ID
   */
  async getById(projectId: number): Promise<Project> {
    const response = await fetch(`${API_BASE}/projects/getFromId?id=${projectId}`, {
      credentials: 'include',
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to fetch project');
    }

    return data.data;
  },

  /**
   * Create new project
   */
  async create(request: CreateProjectRequest): Promise<Project> {
    const response = await fetch(`${API_BASE}/projects/create`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(request),
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to create project');
    }

    return data.data;
  },

  /**
   * Update project
   */
  async update(projectId: number, updates: UpdateProjectRequest): Promise<Project> {
    const response = await fetch(`${API_BASE}/projects/update`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ id: projectId, ...updates }),
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to update project');
    }

    return data.data;
  },

  /**
   * Delete project
   */
  async delete(projectId: number): Promise<void> {
    const response = await fetch(`${API_BASE}/projects/delete?id=${projectId}`, {
      method: 'DELETE',
      credentials: 'include',
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to delete project');
    }
  },

  /**
   * Update project members
   */
  async updateMembers(projectId: number, userIds: number[]): Promise<Project> {
    const response = await fetch(`${API_BASE}/projects/updateMembers?id=${projectId}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ userIds }),
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to update project members');
    }

    return data.data;
  },
};
