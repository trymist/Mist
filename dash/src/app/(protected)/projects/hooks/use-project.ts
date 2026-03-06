import { useState, useCallback, useEffect } from 'react';
import { toast } from 'sonner';
import { projectsService } from '@/services';
import type { Project, ProjectCreateInput } from '@/types';
import { SUCCESS_MESSAGES, ERROR_MESSAGES } from '@/constants';

interface UseProjectOptions {
  projectId: number;
  autoFetch?: boolean;
}

interface UseProjectReturn {
  project: Project | null;
  loading: boolean;
  error: string | null;
  fetchProject: () => Promise<void>;
  updateProject: (data: ProjectCreateInput) => Promise<Project | null>;
  deleteProject: () => Promise<boolean>;
  refreshProject: () => Promise<void>;
}

export const useProject = (options: UseProjectOptions): UseProjectReturn => {
  const { projectId, autoFetch = true } = options;
  
  const [project, setProject] = useState<Project | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchProject = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await projectsService.getById(projectId);
      setProject(data);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : ERROR_MESSAGES.GENERIC_ERROR;
      setError(errorMessage);
      toast.error(errorMessage);
    } finally {
      setLoading(false);
    }
  }, [projectId]);

  const updateProject = useCallback(async (data: ProjectCreateInput): Promise<Project | null> => {
    try {
      const updatedProject = await projectsService.update(projectId, data);
      setProject(updatedProject);
      toast.success(SUCCESS_MESSAGES.PROJECT_UPDATED);
      return updatedProject;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : ERROR_MESSAGES.GENERIC_ERROR;
      toast.error(errorMessage);
      return null;
    }
  }, [projectId]);

  const deleteProject = useCallback(async (): Promise<boolean> => {
    try {
      await projectsService.delete(projectId);
      toast.success(SUCCESS_MESSAGES.PROJECT_DELETED);
      return true;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : ERROR_MESSAGES.GENERIC_ERROR;
      toast.error(errorMessage);
      return false;
    }
  }, [projectId]);

  const refreshProject = useCallback(() => fetchProject(), [fetchProject]);

  useEffect(() => {
    if (autoFetch) {
      fetchProject();
    }
  }, [autoFetch, fetchProject]);

  return {
    project,
    loading,
    error,
    fetchProject,
    updateProject,
    deleteProject,
    refreshProject,
  };
};
