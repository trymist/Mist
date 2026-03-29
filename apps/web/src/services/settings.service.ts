import { apiClient } from '@/api';

export interface SystemSettings {
  wildcardDomain: string | null;
  mistAppName: string;
  allowedOrigins: string;
  productionMode: boolean;
  secureCookies: boolean;
  autoCleanupContainers: boolean;
  autoCleanupImages: boolean;
}

export interface UpdateSystemSettingsRequest {
  wildcardDomain?: string | null;
  mistAppName?: string;
  allowedOrigins?: string;
  productionMode?: boolean;
  secureCookies?: boolean;
  autoCleanupContainers?: boolean;
  autoCleanupImages?: boolean;
}

export type DockerCleanupType = 'containers' | 'images' | 'system' | 'system-all';

export interface DockerCleanupResponse {
  message: string;
  type: string;
}

export const settingsService = {
  async getSystemSettings(): Promise<SystemSettings> {
    const response = await apiClient.get<SystemSettings>('/settings/system');
    return response.data;
  },

  async updateSystemSettings(
    settings: UpdateSystemSettingsRequest
  ): Promise<SystemSettings> {
    const response = await apiClient.put<SystemSettings>('/settings/system', settings);
    return response.data;
  },

  async dockerCleanup(type: DockerCleanupType): Promise<DockerCleanupResponse> {
    const response = await apiClient.post<DockerCleanupResponse>('/settings/docker/cleanup', { type });
    return response.data;
  },
};
