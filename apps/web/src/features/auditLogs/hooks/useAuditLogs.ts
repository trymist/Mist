import { useState, useEffect, useCallback } from 'react';
import { auditLogsService } from '@/services';
import type { AuditLog } from '@/types';

interface UseAuditLogsOptions {
  limit?: number;
  resourceFilter?: string;
  autoFetch?: boolean;
}

interface UseAuditLogsReturn {
  logs: AuditLog[];
  loading: boolean;
  error: string | null;
  total: number;
  page: number;
  totalPages: number;
  setPage: (page: number) => void;
  setResourceFilter: (filter: string) => void;
  refetch: () => Promise<void>;
}

export const useAuditLogs = (
  options: UseAuditLogsOptions = {}
): UseAuditLogsReturn => {
  const { limit = 50, resourceFilter: initialFilter = 'all', autoFetch = true } = options;
  
  const [logs, setLogs] = useState<AuditLog[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(0);
  const [resourceFilter, setResourceFilter] = useState(initialFilter);

  const fetchLogs = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      
      const params: {
        limit: number;
        offset: number;
        resourceType?: string;
      } = {
        limit,
        offset: page * limit,
      };
      
      if (resourceFilter !== 'all') {
        params.resourceType = resourceFilter;
      }
      
      const response = await auditLogsService.getAll(params);
      setLogs(response.logs || []);
      setTotal(response.total || 0);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch audit logs');
    } finally {
      setLoading(false);
    }
  }, [page, resourceFilter, limit]);

  useEffect(() => {
    if (autoFetch) {
      fetchLogs();
    }
  }, [fetchLogs, autoFetch]);

  const totalPages = Math.ceil(total / limit);

  return {
    logs,
    loading,
    error,
    total,
    page,
    totalPages,
    setPage,
    setResourceFilter,
    refetch: fetchLogs,
  };
};
