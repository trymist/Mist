import { apiClient, type ApiResponse } from "..";
import type { User } from "@/types";



export const authApi = {
  async getMe(): Promise<ApiResponse<{
    setupRequired?: boolean; user?: User
  }>> {
    return apiClient.get('/auth/me');
  },

  async login(email: string, password: string): Promise<ApiResponse<User>> {
    return apiClient.post<User>('/auth/login', { email, password });
  },

  async signup(email: string, password: string, username: string): Promise<ApiResponse<User>> {
    return apiClient.post<User>('/auth/signup', { email, password, username });
  },

  async logout(): Promise<ApiResponse<void>> {
    return apiClient.post<void>('/auth/logout');
  }


}
