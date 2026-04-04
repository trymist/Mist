
import React, { useEffect, useState, useMemo } from 'react';
import { toast } from 'sonner';

import { useAuth } from '@/providers';
import { useProjects } from '@/hooks';

import { Button } from '../../components/ui/button';
import { Card, CardContent } from '../../components/ui/card';
import { FormModal } from '@/components/common/form-modal';
import { Plus, Search } from 'lucide-react';
import { Input } from '../../components/ui/input';
import { FullScreenLoading } from '@/components/common';

import { canEditProject, canDeleteProject, validateProjectData, filterProjects, sortProjects } from '../projects/utils';
import type { Project, ProjectCreateInput } from '../../types';
import { ProjectCard } from './components/ProjectCard';


const EmptyState: React.FC<{ onCreateProject: () => void; canCreate: boolean }> = ({
  onCreateProject,
  canCreate
}) => (
  <div className="flex flex-col items-center justify-center py-12">
    <div className="text-center">
      <h3 className="text-lg font-semibold mb-2">No projects yet</h3>
      <p className="text-muted-foreground mb-4">
        Get started by creating your first project
      </p>
      {canCreate && (
        <Button onClick={onCreateProject}>
          <Plus className="w-4 h-4 mr-2" />
          Create first project
        </Button>
      )}
    </div>
  </div>
);

export const ProjectsPage: React.FC = () => {
  const { user } = useAuth();
  const { projects, loading, error, createProject, updateProject, deleteProject, fetchProjects } = useProjects();

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingProject, setEditingProject] = useState<Project | null>(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [sortBy] = useState<'name' | 'created' | 'updated' | 'owner'>('created');
  const [sortOrder] = useState<'asc' | 'desc'>('desc');

  const canCreateProjects = user?.isAdmin || user?.role === 'owner';

  const filteredAndSortedProjects = useMemo(() => {
    const filtered = filterProjects(projects, searchTerm);
    return sortProjects(filtered, sortBy, sortOrder);
  }, [projects, searchTerm, sortBy, sortOrder]);

  useEffect(() => {
    fetchProjects();

    document.body.style.overflow = 'hidden';
    return () => {
      document.body.style.overflow = '';
    };
  }, [fetchProjects]);

  const handleCreateOrUpdateProject = async (projectData: Record<string, unknown>) => {
    const validation = validateProjectData(projectData as unknown as ProjectCreateInput);
    if (!validation.isValid) {
      const firstError = Object.values(validation.errors)[0];
      toast.error(firstError);
      return;
    }

    const isEditing = !!editingProject;
    const result = isEditing 
      ? await updateProject(editingProject.id, projectData as unknown as ProjectCreateInput)
      : await createProject(projectData as unknown as ProjectCreateInput);

    if (result) {
      setIsModalOpen(false);
      setEditingProject(null);
    }
  };

  const handleDeleteProject = async (projectId: number) => {
    if (!confirm('Are you sure you want to delete this project?')) return;
    await deleteProject(projectId);
  };

  const handleEditProject = (project: Project) => {
    setEditingProject(project);
    setIsModalOpen(true);
  };

  const handleCreateProject = () => {
    setEditingProject(null);
    setIsModalOpen(true);
  };

  if (loading && projects.length === 0) {
    return <FullScreenLoading />;
  }

  return (
    <div className="flex flex-col h-screen bg-background overflow-hidden">
      {/* Header */}
      <div className="flex flex-col gap-4 py-6 border-b border-border shrink-0">
        <div className="flex flex-col sm:flex-row justify-between sm:items-center gap-4">
          <div>
            <h1 className="text-2xl font-bold text-foreground">Projects</h1>
            <p className="text-muted-foreground">
              Create and manage your projects
            </p>
          </div>
          {canCreateProjects && (
            <Button onClick={handleCreateProject}>
              <Plus className="w-4 h-4 mr-2" />
              New Project
            </Button>
          )}
        </div>

        {/* Search and filters */}
        {projects.length > 0 && (
          <div className="flex flex-col sm:flex-row gap-4">
            <div className="relative flex-1 max-w-md">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground w-4 h-4" />
              <Input
                placeholder="Search projects..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="pl-10"
              />
            </div>
            {/* Add more filters here if needed */}
          </div>
        )}
      </div>

      {/* Content */}
      <div className="flex-1 py-6 overflow-y-auto">
        {error && (
          <Card className="border-destructive text-destructive mb-6">
            <CardContent className="p-4 text-center">
              <p>{error}</p>
              <Button
                variant="outline"
                size="sm"
                className="mt-2"
                onClick={() => fetchProjects()}
              >
                Try Again
              </Button>
            </CardContent>
          </Card>
        )}

        {filteredAndSortedProjects.length === 0 ? (
          <EmptyState
            onCreateProject={handleCreateProject}
            canCreate={canCreateProjects}
          />
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 pb-16">
            {filteredAndSortedProjects.map((project) => (
              <ProjectCard
                key={project.id}
                project={project}
                onEdit={handleEditProject}
                onDelete={handleDeleteProject}
                canEdit={canEditProject(project, user)}
                canDelete={canDeleteProject(project, user)}
              />
            ))}
          </div>
        )}
      </div>

      <FormModal
        isOpen={isModalOpen}
        onClose={() => {
          setIsModalOpen(false);
          setEditingProject(null);
        }}
        onSubmit={handleCreateOrUpdateProject}
        title={editingProject ? 'Edit Project' : 'Create New Project'}
        fields={[
          {
            name: 'name',
            label: 'Project Name',
            type: 'text',
            required: true,
            defaultValue: editingProject?.name || '',
          },
          {
            name: 'description',
            label: 'Description',
            type: 'textarea',
            defaultValue: editingProject?.description || '',
          },
          {
            name: 'tags',
            label: 'Tags',
            type: 'tags',
            defaultValue: editingProject?.tags || [],
          },
        ]}
      />
    </div>
  );
};

export default ProjectsPage;
