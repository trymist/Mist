import { FormModal } from "@/components/common/form-modal";
import { FullScreenLoading } from "@/components/common";
import { Button } from "@/components/ui/button";
import { useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { useApplication } from "@/hooks";
import { TabsList, Tabs, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import { AppInfo, GitProviderTab, EnvironmentVariables, Domains, AppSettings, LiveLogsViewer, AppStats, Volumes, ContainerStats } from "@/components/applications";
import { DeploymentsTab } from "@/components/deployments";
import { ComposeAppPage } from "./ComposeAppPage";


export const AppPage = () => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [activeTab, setActiveTab] = useState("info");

  const params = useParams();
  const navigate = useNavigate();

  const appId = useMemo(() => Number(params.appId), [params.appId]);
  const projectId = useMemo(() => Number(params.projectId), [params.projectId])

  const {
    app,
    loading,
    error,
    latestCommit,
    previewUrl,
    updateApp,
    deleteApp,
    refreshApp,
  } = useApplication({
    appId,
    autoFetch: true,
    projectId
  });

  const deleteAppHandler = async () => {
    const success = await deleteApp();
    if (success) {
      navigate(-1);
    }
  };

  const handleUpdateApp = async (appData: {
    name: string;
    description: string;
  }) => {
    const result = await updateApp(appData);
    if (result) {
      setIsModalOpen(false);
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

  if (!app) return null;

  if (app.appType === 'compose') {
    return <ComposeAppPage />;
  }

  return (
    <div className="flex flex-col min-h-screen bg-background">
      {/* Header */}
      <header className="border-b border-border py-6 flex flex-col sm:flex-row justify-between gap-4">
        <div>
          <h1 className="text-2xl font-semibold">{app.name}</h1>
          <p className="text-muted-foreground mt-1">{app.description}</p>
        </div>

        <div className="flex flex-wrap gap-2 sm:flex-nowrap">
          <Button variant="outline" onClick={() => setIsModalOpen(true)}>
            Edit App
          </Button>
          <Button variant="destructive" onClick={deleteAppHandler}>
            Delete App
          </Button>
        </div>
      </header>

      {/* App Info */}
      <main className="flex-1 overflow-y-auto py-6">
        <Tabs defaultValue="info" className="w-full" onValueChange={setActiveTab}>
          <div className="w-full overflow-x-auto mb-6 pb-1">
            <TabsList className="inline-flex w-full min-w-fit">
              <TabsTrigger value="info">Info</TabsTrigger>
              {app.appType !== 'database' && <TabsTrigger value="git">Git</TabsTrigger>}
              <TabsTrigger value="environment">Environment</TabsTrigger>
              {app.appType === 'web' && <TabsTrigger value="domains">Domains</TabsTrigger>}
              <TabsTrigger value="deployments">Deployments</TabsTrigger>
              <TabsTrigger value="stats">Stats</TabsTrigger>
              <TabsTrigger value="logs">Logs</TabsTrigger>
              <TabsTrigger value="settings">Settings</TabsTrigger>
            </TabsList>
          </div>

          {/* ✅ INFO TAB */}
          <TabsContent value="info" className="space-y-6">
            <div className="grid grid-cols-1 xl:grid-cols-3 gap-6">
              <div className="xl:col-span-2">
                <AppInfo app={app} latestCommit={latestCommit} />
              </div>
              <div>
                <AppStats appId={app.id} appStatus={app.status} app={app} previewUrl={previewUrl} onStatusChange={refreshApp} />
              </div>
            </div>
          </TabsContent>

          {app.appType !== 'database' && (
            <TabsContent value="git" className="space-y-6">
              <GitProviderTab app={app} />
            </TabsContent>
          )}

          {/* ✅ ENVIRONMENT TAB */}
          <TabsContent value="environment" className="space-y-6">
            <EnvironmentVariables appId={app.id} />
          </TabsContent>

          {/* ✅ DOMAINS TAB */}
          {app.appType === 'web' && (
            <TabsContent value="domains" className="space-y-6">
              <Domains appId={app.id} />
            </TabsContent>
          )}

          {/* ✅ DEPLOYMENTS TAB */}
          <TabsContent value="deployments">
            <DeploymentsTab appId={app.id} app={app} />
          </TabsContent>

          {/* ✅ STATS TAB */}
          <TabsContent value="stats" className="space-y-6">
            <ContainerStats appId={app.id} />
          </TabsContent>

          <TabsContent value="logs" className="h-full">
            <LiveLogsViewer appId={app.id} enabled={activeTab === "logs"} />
          </TabsContent>

          <TabsContent value="settings" className="space-y-6">
            <AppSettings app={app} onUpdate={refreshApp} />
            <Volumes appId={app.id} appType={app.appType} />
          </TabsContent>
        </Tabs>
      </main>

      {/* Edit Modal */}
      <FormModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        title="Edit App"
        fields={[
          { label: "App Name", name: "name", type: "text", defaultValue: app.name },
          { label: "Description", name: "description", type: "textarea", defaultValue: app.description || "" },
        ]}
        onSubmit={(data) => handleUpdateApp(data as { name: string; description: string })}
      />
    </div>
  );
};
