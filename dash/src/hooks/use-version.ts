import { useState, useEffect } from 'react';
import { updatesApi } from '@/api/endpoints/updates';

interface UseVersionReturn {
  version: string;
  loading: boolean;
  error: string | null;
}

export const useVersion = (): UseVersionReturn => {
  const [version, setVersion] = useState<string>('1.0.0');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchVersion = async () => {
      try {
        setLoading(true);
        setError(null);
        const response = await updatesApi.getCurrentVersion();
        if (response.success && response.data) {
          setVersion(response.data.version);
        }
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : 'Failed to fetch version';
        setError(errorMessage);
        // Silently fail, keep default version
        console.error('Failed to fetch version:', errorMessage);
      } finally {
        setLoading(false);
      }
    };

    fetchVersion();
  }, []);

  return {
    version,
    loading,
    error,
  };
};
