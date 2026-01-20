import type { Deployment } from '@/types';

const API_BASE = '/api';

export const deploymentsService = {
  /**
   * Get deployments for an application
   */
  async getByAppId(appId: number): Promise<Deployment[]> {
    const response = await fetch(`${API_BASE}/deployments/getByAppId`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ appId }),
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to fetch deployments');
    }

    return data.data || [];
  },

  /**
   * Create new deployment
   */
  async create(appId: number): Promise<{ id: number }> {
    const response = await fetch(`${API_BASE}/deployments/create`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ appId }),
    });

    const data = await response.json();
    if (!response.ok) {
      throw new Error(data.error || 'Failed to create deployment');
    }

    return data;
  },

  /**
   * Get completed deployment logs (REST API)
   */
  async getCompletedLogs(deploymentId: number): Promise<{ deployment: Deployment; logs: string }> {
    const response = await fetch(`${API_BASE}/deployments/logs?id=${deploymentId}`, {
      credentials: 'include',
    });

    if (!response.ok) {
      if (response.status === 400) {
        throw new Error('DEPLOYMENT_IN_PROGRESS');
      }
      throw new Error('Failed to fetch deployment logs');
    }

    const data = await response.json();
    return data.data;
  },

  /**
   * Get WebSocket URL for live deployment streaming
   */
  getWebSocketUrl(deploymentId: number): string {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = window.location.host;
    return `${protocol}//${host}${API_BASE}/deployments/logs/stream?id=${deploymentId}`;
  },

  /**
   * Stop a deployment
   */
  async stopDeployment(deploymentId: number): Promise<void> {
    const response = await fetch(`${API_BASE}/deployments/stopDep`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ deploymentId }),
    });

    const data = await response.json();
    if (!response.ok) {
      throw new Error(data.message || data.error || 'Failed to stop deployment');
    }
  },
};
