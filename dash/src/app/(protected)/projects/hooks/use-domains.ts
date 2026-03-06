import { useState, useCallback, useEffect } from 'react';
import { toast } from 'sonner';
import { applicationsService } from '@/services';
import type { Domain } from '@/types';

interface UseDomainsOptions {
  appId: number;
  autoFetch?: boolean;
}

interface UseDomainsReturn {
  domains: Domain[];
  loading: boolean;
  error: string | null;
  fetchDomains: () => Promise<void>;
  createDomain: (domain: string) => Promise<Domain | null>;
  updateDomain: (id: number, domain: string) => Promise<Domain | null>;
  deleteDomain: (id: number) => Promise<boolean>;
  refreshDomains: () => Promise<void>;
  updateDomainInState: (updatedDomain: Domain) => void;
}

export const useDomains = (options: UseDomainsOptions): UseDomainsReturn => {
  const { appId, autoFetch = true } = options;
  
  const [domains, setDomains] = useState<Domain[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchDomains = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await applicationsService.getDomains(appId);
      setDomains(data);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to fetch domains';
      setError(errorMessage);
      toast.error(errorMessage);
    } finally {
      setLoading(false);
    }
  }, [appId]);

  const createDomain = useCallback(async (domain: string): Promise<Domain | null> => {
    try {
      const response = await applicationsService.createDomain({ appId, domain });
      // The backend now returns { domain, actionRequired, actionMessage }
      const newDomain = response?.domain || response;
      if (newDomain) {
        setDomains(prev => [...prev, newDomain]);
        toast.success('Domain added');
      }
      return newDomain;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to add domain';
      toast.error(errorMessage);
      return null;
    }
  }, [appId]);

  const updateDomain = useCallback(async (id: number, domain: string): Promise<Domain | null> => {
    try {
      const response = await applicationsService.updateDomain({ id, domain });
      // The backend now returns { domain, actionRequired, actionMessage }
      const updatedDomain = response?.domain || response;
      if (updatedDomain) {
        setDomains(prev => prev.map(d => d.id === id ? updatedDomain : d));
        toast.success('Domain updated');
      }
      return updatedDomain;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to update domain';
      toast.error(errorMessage);
      return null;
    }
  }, []);

  const deleteDomain = useCallback(async (id: number): Promise<boolean> => {
    try {
      await applicationsService.deleteDomain(id);
      setDomains(prev => prev.filter(d => d.id !== id));
      toast.success('Domain deleted');
      return true;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to delete domain';
      toast.error(errorMessage);
      return false;
    }
  }, []);

  const refreshDomains = useCallback(() => fetchDomains(), [fetchDomains]);

  const updateDomainInState = useCallback((updatedDomain: Domain) => {
    setDomains(prev => prev.map(d => d.id === updatedDomain.id ? updatedDomain : d));
  }, []);

  useEffect(() => {
    if (autoFetch) {
      fetchDomains();
    }
  }, [autoFetch, fetchDomains]);

  return {
    domains,
    loading,
    error,
    fetchDomains,
    createDomain,
    updateDomain,
    deleteDomain,
    refreshDomains,
    updateDomainInState,
  };
};
