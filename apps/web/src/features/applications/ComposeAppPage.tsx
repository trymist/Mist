import { FormModal } from "@/components/common/form-modal";
import { FullScreenLoading } from "@/components/common";
import { Button } from "@/components/ui/button";
import { useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { useApplication } from "@/hooks";
import { TabsList, Tabs, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import { AppInfo, GitProviderTab, EnvironmentVariables, LiveLogsViewer, Volumes } from "@/components/applications";
import { ComposeStatus } from "@/components/applications/compose-status";
import { ComposeAppSettings } from "@/components/applications/compose-app-settings";
import { DeploymentsTab } from "@/components/deployments";


export const ComposeAppPage = () => {
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

    return (
        <div className="flex flex-col min-h-screen bg-background">
            {/* Header */}
            <header className="border-b border-border py-6 flex flex-col sm:flex-row justify-between gap-4">
                <div>
                    <div className="flex items-center gap-2">
                        <h1 className="text-2xl font-semibold">{app.name}</h1>
                        <span className="px-2 py-0.5 rounded-full bg-blue-500/10 text-blue-500 text-xs font-medium border border-blue-500/20">Compose</span>
                    </div>
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
                            <TabsTrigger value="git">Git</TabsTrigger>
                            <TabsTrigger value="environment">Environment</TabsTrigger>
                            <TabsTrigger value="deployments">Deployments</TabsTrigger>
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
                                <ComposeStatus appId={app.id} onStatusChange={refreshApp} />
                            </div>
                        </div>
                    </TabsContent>

                    <TabsContent value="git" className="space-y-6">
                        <GitProviderTab app={app} />
                    </TabsContent>

                    {/* ✅ ENVIRONMENT TAB */}
                    <TabsContent value="environment" className="space-y-6">
                        <EnvironmentVariables appId={app.id} />
                    </TabsContent>

                    {/* ✅ DEPLOYMENTS TAB */}
                    <TabsContent value="deployments">
                        <DeploymentsTab appId={app.id} app={app} />
                    </TabsContent>

                    {/* Stats tab removed as requested ("keep the stats if its possible to get the compose stats" -> we put compose stats in info/status card for now as overall stats. Detailed container stats are hard for multiple containers yet. User said 'keep only the necessary things') */}
                    {/* Actually user said "keep the stats if its possible to get the compose stats". But our backend currently doesn't provide resource stats for compose. */}
                    {/* I will hide the stats tab for now. */}

                    <TabsContent value="logs" className="h-full">
                        <LiveLogsViewer appId={app.id} enabled={activeTab === "logs"} />
                    </TabsContent>

                    <TabsContent value="settings" className="space-y-6">
                        <ComposeAppSettings app={app} onUpdate={refreshApp} />
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
