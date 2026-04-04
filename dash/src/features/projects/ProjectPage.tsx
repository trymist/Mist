import { FormModal } from "@/components/common/form-modal";
import { FullScreenLoading } from "@/components/common";
import { Button } from "@/components/ui/button";
import { useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { toast } from "sonner";
import type { CreateAppRequest } from "@/types/app";
import { AppCard } from "./components/AppCard";
import { CreateAppModal } from "./components/CreateAppModal";
import { ManageMembersModal } from "./components/ManageMembersModal";
import { useApplications, useProject } from "@/hooks";
import { useAuth } from "@/providers";
import { applicationsService } from "@/services";

export const ProjectPage = () => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isAddNewAppModalOpen, setIsAddNewAppModalOpen] = useState(false);
  const [isManageMembersModalOpen, setIsManageMembersModalOpen] = useState(false);

  const params = useParams();
  const navigate = useNavigate();
  const projectId = parseInt(params.projectId!);
  const { user } = useAuth();

  const { project, loading, error, updateProject, deleteProject: deleteProjectFn, refreshProject } = useProject({
    projectId,
    autoFetch: true,
  });

  const { apps, loading: fetchingApps, fetchApps
  } = useApplications({
    projectId,
    autoFetch: true,
  });

  const createNewApp = async (appData: CreateAppRequest) => {
    try {
      const response = await fetch(`/api/apps/create`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(appData),
        credentials: "include",
      });
      const data = await response.json();
      if (!data.success) throw new Error(data.error || "Failed to create app");

      toast.success("App created successfully");
      fetchApps()
      setIsAddNewAppModalOpen(false);
    }
    catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to create app");
    }
  };

  const deleteProjectHandler = async () => {
    const success = await deleteProjectFn();
    if (success) {
      navigate("/projects");
    }
  };

  const handleUpdateProject = async (projectData: {
    name: string;
    description: string;
    tags: string[];
  }) => {
    const result = await updateProject(projectData);
    if (result) {
      setIsModalOpen(false);
    }
  };

  const handleDeleteApp = async (appId: number) => {
    try {
      await applicationsService.delete(appId);
      toast.success("Application deleted successfully");
      fetchApps();
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to delete application";
      toast.error(errorMessage);
    }
  };

  if (loading) return <FullScreenLoading />;

  if (error)
    return (
      <div className="min-h-screen flex items-center justify-center bg-background">
        <div className="bg-destructive/10 border border-destructive text-destructive p-4 rounded-lg max-w-md text-center">
          {error}
        </div>
      </div>
    );

  if (!project) return null;

  // Check if current user is the project owner
  const isProjectOwner = user && Number(user.id) === Number(project.ownerId);

  return (
    <div className="flex flex-col min-h-screen bg-background">
      {/* Header */}
      <header className="border-b border-border py-6 flex flex-col sm:flex-row justify-between gap-4">
        <div>
          <h1 className="text-2xl font-semibold">{project.name}</h1>
          <p className="text-muted-foreground mt-1">{project.description}</p>
          {project.projectMembers && (
            <div className="mt-3 flex flex-wrap gap-2 items-center">
              <span className="text-sm font-medium text-foreground">Members:</span>
              {project.projectMembers.map((member: { username: string; id: number }) => (
                <span
                  key={member.id}
                  className="bg-secondary text-secondary-foreground px-2 py-1 rounded-full text-xs"
                >
                  {member.username}
                </span>
              ))}
              {isProjectOwner && (
                <Button 
                  variant="secondary" 
                  size="sm" 
                  className="ml-2"
                  onClick={() => setIsManageMembersModalOpen(true)}
                >
                  Manage Members
                </Button>
              )}
            </div>
          )}
        </div>

        <div className="flex flex-wrap gap-2 sm:flex-nowrap">
          <Button variant="outline" onClick={() => setIsModalOpen(true)}>
            Edit Project
          </Button>
          <Button variant="destructive" onClick={deleteProjectHandler}>
            Delete Project
          </Button>
          <Button variant="secondary" onClick={() => setIsAddNewAppModalOpen(true)}>
            Create App
          </Button>
        </div>
      </header>

      {/* Apps Section */}
      <main className="flex-1 overflow-y-auto py-6">

        {fetchingApps ? (
          <div className="text-muted-foreground text-center py-10">Loading apps...</div>
        ) : apps && apps.length === 0 ? (
          <div className="text-muted-foreground text-center py-10">
            No apps found for this project. Create one to get started.
          </div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            {apps && apps.map((app) => (
              <AppCard
                key={app.id}
                app={app}
                onClick={() => navigate(`/projects/${app.projectId}/apps/${app.id}`)}
                onDelete={handleDeleteApp}
              />
            ))}
          </div>
        )}
      </main>

      {/* Modals */}
      <FormModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        title="Edit Project"
        fields={[
          { label: "Name", name: "name", type: "text", defaultValue: project.name },
          { label: "Description", name: "description", type: "textarea", defaultValue: project.description },
          { name: "tags", label: "Tags", type: "tags", defaultValue: project.tags || [] },
        ]}
        onSubmit={(data) => handleUpdateProject(data as { name: string; description: string; tags: string[] })}
      />

      <CreateAppModal
        isOpen={isAddNewAppModalOpen}
        onClose={() => setIsAddNewAppModalOpen(false)}
        projectId={projectId}
        onSubmit={createNewApp}
      />

      <ManageMembersModal
        isOpen={isManageMembersModalOpen}
        onClose={() => setIsManageMembersModalOpen(false)}
        project={project}
        onSuccess={refreshProject}
      />
    </div>
  );
};
