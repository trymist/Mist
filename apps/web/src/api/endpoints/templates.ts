import { apiClient, type ApiResponse } from "..";
import type { ServiceTemplate } from "@/types/app";

export const templatesApi = {
  async list(): Promise<ApiResponse<ServiceTemplate[]>> {
    return apiClient.get('/templates/list');
  },

  async getByName(name: string): Promise<ApiResponse<ServiceTemplate>> {
    return apiClient.get(`/templates/get?name=${name}`);
  }
};
