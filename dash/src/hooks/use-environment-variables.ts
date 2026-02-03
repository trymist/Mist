import { useState, useCallback, useEffect } from 'react';
import { toast } from 'sonner';
import { applicationsService } from '@/services';
import type { EnvVariable } from '@/types';

interface UseEnvironmentVariablesOptions {
  appId: number;
  autoFetch?: boolean;
}

interface UseEnvironmentVariablesReturn {
  envVars: EnvVariable[];
  loading: boolean;
  error: string | null;
  fetchEnvVars: () => Promise<void>;
  createEnvVar: (key: string, value: string) => Promise<EnvVariable | null>;
  updateEnvVar: (id: number, key: string, value: string) => Promise<EnvVariable | null>;
  deleteEnvVar: (id: number) => Promise<boolean>;
  refreshEnvVars: () => Promise<void>;
}

export const useEnvironmentVariables = (options: UseEnvironmentVariablesOptions): UseEnvironmentVariablesReturn => {
  const { appId, autoFetch = true } = options;
  
  const [envVars, setEnvVars] = useState<EnvVariable[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchEnvVars = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const data: Array<EnvVariable | null | undefined> = await applicationsService.getEnvVariables(appId);
      // Filter out any null/undefined entries
      setEnvVars(data.filter((env): env is EnvVariable => env != null));
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to fetch environment variables';
      setError(errorMessage);
      toast.error(errorMessage);
    } finally {
      setLoading(false);
    }
  }, [appId]);

  const createEnvVar = useCallback(async (key: string, value: string): Promise<EnvVariable | null> => {
    try {
      const response = await applicationsService.createEnvVariable({ appId, key, value });
      // The backend now returns { envVariable, actionRequired, actionMessage }
      const envVar = response?.envVariable || response;
      if (envVar) {
        setEnvVars(prev => [...prev, envVar]);
        toast.success('Environment variable added');
      }
      return envVar;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to add environment variable';
      toast.error(errorMessage);
      return null;
    }
  }, [appId]);

  const updateEnvVar = useCallback(async (id: number, key: string, value: string): Promise<EnvVariable | null> => {
    try {
      const response = await applicationsService.updateEnvVariable({ id, key, value });
      // The backend now returns { envVariable, actionRequired, actionMessage }
      const updatedVar = response?.envVariable || response;
      if (updatedVar) {
        setEnvVars(prev => prev.map(v => v.id === id ? updatedVar : v));
        toast.success('Environment variable updated');
      }
      return updatedVar;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to update environment variable';
      toast.error(errorMessage);
      return null;
    }
  }, []);

  const deleteEnvVar = useCallback(async (id: number): Promise<boolean> => {
    try {
      await applicationsService.deleteEnvVariable(id);
      setEnvVars(prev => prev.filter(v => v.id !== id));
      toast.success('Environment variable deleted');
      return true;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to delete environment variable';
      toast.error(errorMessage);
      return false;
    }
  }, []);

  const refreshEnvVars = useCallback(() => fetchEnvVars(), [fetchEnvVars]);

  useEffect(() => {
    if (autoFetch) {
      fetchEnvVars();
    }
  }, [autoFetch, fetchEnvVars]);

  return {
    envVars,
    loading,
    error,
    fetchEnvVars,
    createEnvVar,
    updateEnvVar,
    deleteEnvVar,
    refreshEnvVars,
  };
};
