import { useState, useCallback, useEffect } from 'react';
import { toast } from 'sonner';
import { applicationsService } from '@/services';
import type { App, UpdateAppRequest } from '@/types';

interface UseApplicationOptions {
  appId: number;
  projectId: number;
  autoFetch?: boolean;
}

interface UseApplicationReturn {
  app: App | null;
  loading: boolean;
  error: string | null;
  latestCommit: {
    sha: string;
    html_url: string;
    author?: string;
    timestamp?: string;
    message?: string;
  } | null;
  previewUrl: string;
  fetchApp: () => Promise<void>;
  fetchLatestCommit: () => Promise<void>;
  fetchPreviewUrl: () => Promise<void>;
  updateApp: (data: UpdateAppRequest) => Promise<App | null>;
  deleteApp: () => Promise<boolean>;
  refreshApp: () => Promise<void>;
}

export const useApplication = (options: UseApplicationOptions): UseApplicationReturn => {
  const { appId, projectId, autoFetch = true } = options;

  const [app, setApp] = useState<App | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [latestCommit, setLatestCommit] = useState<{
    sha: string;
    html_url: string;
    author?: string;
    timestamp?: string;
    message?: string;
  } | null>(null);
  const [previewUrl, setPreviewUrl] = useState<string>('');

  const fetchApp = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await applicationsService.getById(appId);
      setApp(data);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to fetch application';
      setError(errorMessage);
      toast.error(errorMessage);
    } finally {
      setLoading(false);
    }
  }, [appId]);

  const fetchLatestCommit = useCallback(async () => {
    try {
      const data = await applicationsService.getLatestCommit(appId, projectId);
      setLatestCommit(data);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to fetch latest commit';
      toast.error(errorMessage);
    }
  }, [appId]);

  const fetchPreviewUrl = useCallback(async () => {
    if (!app || app.status !== 'running') return;

    try {
      const data = await applicationsService.getPreviewUrl(appId);
      setPreviewUrl(data.url);
    } catch (err) {
      console.error('Failed to fetch preview URL:', err);
    }
  }, [appId, app]);

  const updateApp = useCallback(async (data: UpdateAppRequest): Promise<App | null> => {
    try {
      const updatedApp = await applicationsService.update(appId, data);
      setApp(updatedApp);
      toast.success('Application updated successfully');
      return updatedApp;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to update application';
      toast.error(errorMessage);
      return null;
    }
  }, [appId]);

  const deleteApp = useCallback(async (): Promise<boolean> => {
    try {
      // applicationsService doesn't have a delete method yet, using fetch directly
      const response = await fetch(`/api/apps/delete?id=${appId}`, {
        method: 'DELETE',
        credentials: 'include',
      });
      const data = await response.json();
      if (!data.success) throw new Error(data.error || 'Failed to delete application');

      toast.success('Application deleted successfully');
      return true;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to delete application';
      toast.error(errorMessage);
      return false;
    }
  }, [appId]);

  const refreshApp = useCallback(() => fetchApp(), [fetchApp]);

  useEffect(() => {
    if (autoFetch) {
      fetchApp();
      fetchLatestCommit();
    }
  }, [autoFetch, fetchApp, fetchLatestCommit]);

  useEffect(() => {
    fetchPreviewUrl();
  }, [app?.status, fetchPreviewUrl]);

  return {
    app,
    loading,
    error,
    latestCommit,
    previewUrl,
    fetchApp,
    fetchLatestCommit,
    fetchPreviewUrl,
    updateApp,
    deleteApp,
    refreshApp,
  };
};
