import { useState, useCallback, useEffect } from 'react';
import { toast } from 'sonner';
import { projectsService } from '@/services';
import type { Project, ProjectCreateInput } from '@/types';
import { SUCCESS_MESSAGES, ERROR_MESSAGES } from '@/constants';

interface UseProjectsOptions {
  autoFetch?: boolean;
}

interface UseProjectsReturn {
  projects: Project[];
  loading: boolean;
  error: string | null;
  fetchProjects: () => Promise<void>;
  createProject: (data: ProjectCreateInput) => Promise<Project | null>;
  updateProject: (id: number, data: ProjectCreateInput) => Promise<Project | null>;
  deleteProject: (id: number) => Promise<boolean>;
  refreshProjects: () => Promise<void>;
}

export const useProjects = (options: UseProjectsOptions = {}): UseProjectsReturn => {
  const { autoFetch = false } = options;
  
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchProjects = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await projectsService.getAll();
      setProjects(data);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : ERROR_MESSAGES.GENERIC_ERROR;
      setError(errorMessage);
      toast.error(errorMessage);
    } finally {
      setLoading(false);
    }
  }, []);

  const createProject = useCallback(async (data: ProjectCreateInput): Promise<Project | null> => {
    try {
      const project = await projectsService.create(data);
      setProjects(prev => [...prev, project]);
      toast.success(SUCCESS_MESSAGES.PROJECT_CREATED);
      return project;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : ERROR_MESSAGES.GENERIC_ERROR;
      toast.error(errorMessage);
      return null;
    }
  }, []);

  const updateProject = useCallback(async (id: number, data: ProjectCreateInput): Promise<Project | null> => {
    try {
      const updatedProject = await projectsService.update(id, data);
      setProjects(prev => prev.map(p => p.id === id ? updatedProject : p));
      toast.success(SUCCESS_MESSAGES.PROJECT_UPDATED);
      return updatedProject;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : ERROR_MESSAGES.GENERIC_ERROR;
      toast.error(errorMessage);
      return null;
    }
  }, []);

  const deleteProject = useCallback(async (id: number): Promise<boolean> => {
    try {
      await projectsService.delete(id);
      setProjects(prev => prev.filter(p => p.id !== id));
      toast.success(SUCCESS_MESSAGES.PROJECT_DELETED);
      return true;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : ERROR_MESSAGES.GENERIC_ERROR;
      toast.error(errorMessage);
      return false;
    }
  }, []);

  const refreshProjects = useCallback(() => fetchProjects(), [fetchProjects]);

  useEffect(() => {
    if (autoFetch) {
      fetchProjects();
    }
  }, [autoFetch, fetchProjects]);

  return {
    projects,
    loading,
    error,
    fetchProjects,
    createProject,
    updateProject,
    deleteProject,
    refreshProjects,
  };
};
