import type {
  App,
  CreateAppRequest,
  UpdateAppRequest,
  Volume,
  CreateVolumeRequest,
  UpdateVolumeRequest
} from '@/types';

const API_BASE = '/api';

export const applicationsService = {
  async getById(appId: number): Promise<App> {
    const response = await fetch(`${API_BASE}/apps/getById`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ appId }),
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to fetch application');
    }

    return data.data;
  },

  async getByProjectId(projectId: number): Promise<App[]> {
    const response = await fetch(`${API_BASE}/apps/getByProjectId`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ projectId }),
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to fetch applications');
    }

    return data.data || [];
  },

  async create(request: CreateAppRequest): Promise<App> {
    const response = await fetch(`${API_BASE}/apps/create`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(request),
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to create application');
    }

    return data.data;
  },

  async update(appId: number, updates: UpdateAppRequest): Promise<App & {
    actionRequired: string,
    actionMessage: string
  }> {
    const response = await fetch(`${API_BASE}/apps/update`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ appId, ...updates }),
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to update application');
    }

    return data.data;
  },

  async getLatestCommit(appId: number, projectId: number): Promise<{
    sha: string;
    html_url: string;
    author?: string;
    timestamp?: string;
    message?: string;
  }> {
    const response = await fetch(`${API_BASE}/apps/getLatestCommit`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ appId, projectId }),
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to fetch latest commit');
    }

    return data.data;
  },

  async createEnvVariable(request: { appId: number; key: string; value: string }) {
    const response = await fetch(`${API_BASE}/apps/envs/create`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(request),
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to create environment variable');
    }

    return data.data;
  },

  async getEnvVariables(appId: number) {
    const response = await fetch(`${API_BASE}/apps/envs/get`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ appId }),
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to fetch environment variables');
    }

    return data.data || [];
  },

  async updateEnvVariable(request: { id: number; key: string; value: string }) {
    const response = await fetch(`${API_BASE}/apps/envs/update`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(request),
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to update environment variable');
    }

    return data.data;
  },

  async deleteEnvVariable(id: number) {
    const response = await fetch(`${API_BASE}/apps/envs/delete`, {
      method: 'DELETE',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ id }),
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to delete environment variable');
    }

    return data.data;
  },

  async createDomain(request: { appId: number; domain: string }) {
    const response = await fetch(`${API_BASE}/apps/domains/create`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(request),
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to create domain');
    }

    return data.data;
  },

  async getDomains(appId: number) {
    const response = await fetch(`${API_BASE}/apps/domains/get`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ appId }),
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to fetch domains');
    }

    return data.data || [];
  },

  async updateDomain(request: { id: number; domain: string }) {
    const response = await fetch(`${API_BASE}/apps/domains/update`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(request),
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to update domain');
    }

    return data.data;
  },

  async deleteDomain(id: number) {
    const response = await fetch(`${API_BASE}/apps/domains/delete`, {
      method: 'DELETE',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ id }),
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to delete domain');
    }

    return data.data;
  },

  async verifyDomainDNS(id: number) {
    const response = await fetch(`${API_BASE}/apps/domains/verify`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ id }),
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to verify domain DNS');
    }

    return data.data;
  },

  async getDNSInstructions(id: number) {
    const response = await fetch(`${API_BASE}/apps/domains/instructions`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ id }),
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to get DNS instructions');
    }

    return data.data;
  },

  async getPreviewUrl(appId: number): Promise<{ url: string; domain: string }> {
    const response = await fetch(`${API_BASE}/apps/getPreviewUrl`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ appId }),
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to fetch preview URL');
    }

    return data.data;
  },

  async stopContainer(appId: number): Promise<{ message: string }> {
    const response = await fetch(`${API_BASE}/apps/container/stop?appId=${appId}`, {
      method: 'POST',
      credentials: 'include',
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to stop container');
    }

    return data.data;
  },

  async startContainer(appId: number): Promise<{ message: string }> {
    const response = await fetch(`${API_BASE}/apps/container/start?appId=${appId}`, {
      method: 'POST',
      credentials: 'include',
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to start container');
    }

    return data.data;
  },

  async restartContainer(appId: number): Promise<{ message: string }> {
    const response = await fetch(`${API_BASE}/apps/container/restart?appId=${appId}`, {
      method: 'POST',
      credentials: 'include',
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to restart container');
    }

    return data.data;
  },

  async recreateContainer(appId: number): Promise<{ message: string }> {
    const response = await fetch(`${API_BASE}/apps/container/recreate?appId=${appId}`, {
      method: 'POST',
      credentials: 'include',
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to recreate container');
    }

    return data.data;
  },

  async getContainerStatus(appId: number) {
    const response = await fetch(`${API_BASE}/apps/container/status?appId=${appId}`, {
      method: 'GET',
      credentials: 'include',
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to fetch container status');
    }

    return data.data;
  },

  async getContainerLogs(appId: number, tail: number = 100): Promise<{ logs: string }> {
    const response = await fetch(`${API_BASE}/apps/container/logs?appId=${appId}&tail=${tail}`, {
      method: 'GET',
      credentials: 'include',
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to fetch container logs');
    }

    return data.data;
  },

  async getVolumes(appId: number): Promise<Volume[]> {
    const response = await fetch(`${API_BASE}/apps/volumes/get`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ appId }),
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to fetch volumes');
    }

    return data.data || [];
  },

  async createVolume(request: CreateVolumeRequest): Promise<Volume> {
    const response = await fetch(`${API_BASE}/apps/volumes/create`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(request),
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to create volume');
    }

    return data.data;
  },

  async updateVolume(request: UpdateVolumeRequest): Promise<Volume> {
    const response = await fetch(`${API_BASE}/apps/volumes/update`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(request),
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to update volume');
    }

    return data.data;
  },

  async deleteVolume(id: number): Promise<void> {
    const response = await fetch(`${API_BASE}/apps/volumes/delete`, {
      method: 'DELETE',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ id }),
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to delete volume');
    }
  },

  async delete(appId: number): Promise<void> {
    const response = await fetch(`${API_BASE}/apps/delete?id=${appId}`, {
      method: 'DELETE',
      credentials: 'include',
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to delete application');
    }
  },

  async redeploy(appId: number): Promise<{ deploymentId: number; message: string }> {
    const response = await fetch(`${API_BASE}/deployments/create`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ appId }),
    });

    const data = await response.json();
    if (!data.id) {
      throw new Error(data.error || 'Failed to trigger redeployment');
    }

    return data.data;
  },
};
