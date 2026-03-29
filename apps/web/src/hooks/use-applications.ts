import { useState, useCallback, useEffect } from 'react';
import { toast } from 'sonner';
import { applicationsService } from '@/services';
import type { App, CreateAppRequest } from '@/types';

interface UseApplicationsOptions {
  projectId: number;
  autoFetch?: boolean;
}

interface UseApplicationsReturn {
  apps: App[];
  loading: boolean;
  error: string | null;
  fetchApps: () => Promise<void>;
  createApp: (data: CreateAppRequest) => Promise<App | null>;
  refreshApps: () => Promise<void>;
}

export const useApplications = (options: UseApplicationsOptions): UseApplicationsReturn => {
  const { projectId, autoFetch = true } = options;
  
  const [apps, setApps] = useState<App[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchApps = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await applicationsService.getByProjectId(projectId);
      setApps(data);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to fetch applications';
      setError(errorMessage);
      toast.error(errorMessage);
    } finally {
      setLoading(false);
    }
  }, [projectId]);

  const createApp = useCallback(async (data: CreateAppRequest): Promise<App | null> => {
    try {
      const app = await applicationsService.create({ ...data, projectId });
      setApps(prev => [...prev, app]);
      toast.success('Application created successfully');
      return app;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to create application';
      toast.error(errorMessage);
      return null;
    }
  }, [projectId]);

  const refreshApps = useCallback(() => fetchApps(), [fetchApps]);

  useEffect(() => {
    if (autoFetch) {
      fetchApps();
    }
  }, [autoFetch, fetchApps]);

  return {
    apps,
    loading,
    error,
    fetchApps,
    createApp,
    refreshApps,
  };
};
