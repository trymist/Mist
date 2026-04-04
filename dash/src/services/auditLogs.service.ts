import type { AuditLogsResponse } from '@/types';

const API_BASE = '/api';

export interface GetAuditLogsParams {
  limit?: number;
  offset?: number;
  resourceType?: string;
  resourceId?: number;
}

export const auditLogsService = {
  /**
   * Get all audit logs (admin only)
   */
  async getAll(params?: GetAuditLogsParams): Promise<AuditLogsResponse> {
    const queryParams = new URLSearchParams();
    
    if (params?.limit) queryParams.append('limit', params.limit.toString());
    if (params?.offset) queryParams.append('offset', params.offset.toString());
    if (params?.resourceType) queryParams.append('resourceType', params.resourceType);
    
    const url = `${API_BASE}/audit-logs${queryParams.toString() ? `?${queryParams.toString()}` : ''}`;
    
    const response = await fetch(url, {
      method: 'GET',
      credentials: 'include',
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to fetch audit logs');
    }

    return data.data;
  },

  /**
   * Get audit logs by resource
   */
  async getByResource(
    resourceType: string,
    resourceId: number,
    params?: { limit?: number; offset?: number }
  ): Promise<AuditLogsResponse> {
    const queryParams = new URLSearchParams({
      resourceType,
      resourceId: resourceId.toString(),
    });
    
    if (params?.limit) queryParams.append('limit', params.limit.toString());
    if (params?.offset) queryParams.append('offset', params.offset.toString());
    
    const response = await fetch(
      `${API_BASE}/audit-logs/resource?${queryParams.toString()}`,
      {
        method: 'GET',
        credentials: 'include',
      }
    );

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to fetch audit logs');
    }

    return data.data;
  },
};
