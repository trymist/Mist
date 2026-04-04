import { apiClient, type ApiResponse } from "..";

export interface VersionInfo {
  version: string;
}

export const updatesApi = {
  async getCurrentVersion(): Promise<ApiResponse<VersionInfo>> {
    return apiClient.get('/updates/version');
  },
};
