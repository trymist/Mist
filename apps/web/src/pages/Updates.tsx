import { useState, useEffect } from 'react';
import { useAuth } from '@/providers';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { AlertCircle, CheckCircle2, Download, RefreshCw, ShieldAlert, History, Clock } from 'lucide-react';
import { toast } from 'sonner';
import { useNavigate } from 'react-router-dom';
import { Badge } from '@/components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';

interface UpdateInfo {
  currentVersion: string;
  latestVersion: string;
  updateAvailable: boolean;
  releaseNotes: string;
  releaseName: string;
}

interface UpdateLog {
  id: number;
  versionFrom: string;
  versionTo: string;
  status: 'in_progress' | 'success' | 'failed';
  logs: string;
  errorMessage?: string;
  startedBy: number;
  startedAt: string;
  completedAt?: string;
  username: string;
}

export const UpdatesPage = () => {
  const { user } = useAuth();
  const navigate = useNavigate();
  const [isLoading, setIsLoading] = useState(true);
  const [isCheckingUpdate, setIsCheckingUpdate] = useState(false);
  const [isUpdating, setIsUpdating] = useState(false);
  const [updateInfo, setUpdateInfo] = useState<UpdateInfo | null>(null);
  const [error, setError] = useState('');
  const [updateLogs, setUpdateLogs] = useState<string[]>([]);
  const [updateHistory, setUpdateHistory] = useState<UpdateLog[]>([]);
  const [isLoadingHistory, setIsLoadingHistory] = useState(true);

  useEffect(() => {
    if (!user) {
      navigate('/');
      return;
    }

    if (user.role !== 'owner') {
      toast.error('Only owners can access updates');
      navigate('/');
      return;
    }

    checkForUpdates();
    loadUpdateHistory();
  }, [user, navigate]);

  const checkForUpdates = async () => {
    setIsCheckingUpdate(true);
    setError('');
    try {
      const response = await fetch('/api/updates/check', {
        credentials: 'include',
      });

      if (!response.ok) {
        throw new Error('Failed to check for updates');
      }

      const data = await response.json();
      if (data.success) {
        setUpdateInfo(data.data);
      } else {
        throw new Error(data.message || 'Failed to check for updates');
      }
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to check for updates';
      setError(message);
      toast.error(message);
    } finally {
      setIsLoading(false);
      setIsCheckingUpdate(false);
    }
  };

  const loadUpdateHistory = async () => {
    setIsLoadingHistory(true);
    try {
      const response = await fetch('/api/updates/history?limit=20', {
        credentials: 'include',
      });

      if (!response.ok) {
        throw new Error('Failed to load update history');
      }

      const data = await response.json();
      if (data.success) {
        setUpdateHistory(data.data || []);
      }
    } catch (err) {
      console.error('Failed to load update history:', err);
    } finally {
      setIsLoadingHistory(false);
    }
  };

  const triggerUpdate = async () => {
    // Enhanced confirmation dialog
    const confirmMessage = `
⚠️  IMPORTANT: System Update Confirmation

This will update Mist to version ${updateInfo?.latestVersion}

What will happen:
✓ Database will be backed up automatically
✓ Git repository will be tagged for rollback
✓ Code will be updated from GitHub
✓ Backend and CLI will be rebuilt
✓ Mist service will restart (brief downtime)
✓ Your applications will continue running

Safety measures in place:
• Automatic rollback on failure
• Health checks after update
• Update lock to prevent conflicts
• Full update logs saved

The dashboard will be unavailable for 1-2 minutes during the update.

Do you want to proceed with the update?
    `.trim();

    if (!confirm(confirmMessage)) {
      return;
    }

    setIsUpdating(true);
    setUpdateLogs([]);
    setError('');

    try {
      const response = await fetch('/api/updates/trigger', {
        method: 'POST',
        credentials: 'include',
      });

      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.message || 'Failed to trigger update');
      }

      // Stream the logs
      const reader = response.body?.getReader();
      const decoder = new TextDecoder();

      if (reader) {
        while (true) {
          const { done, value } = await reader.read();
          if (done) break;

          const text = decoder.decode(value);
          const lines = text.split('\n').filter(line => line.trim());
          setUpdateLogs(prev => [...prev, ...lines]);
        }
      }

      toast.success('Update completed successfully');

      // Reload history
      loadUpdateHistory();

      // Wait a bit for the service to restart, then reload
      setTimeout(() => {
        window.location.reload();
      }, 5000);
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Update failed';
      setError(message);
      toast.error(message);

      // Reload history even on failure to show the failed attempt
      loadUpdateHistory();
    } finally {
      setIsUpdating(false);
    }
  };

  if (!user) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <p className="text-muted-foreground">Loading...</p>
      </div>
    );
  }

  if (user.role !== 'owner') {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <Card className="max-w-md">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <ShieldAlert className="h-5 w-5 text-destructive" />
              Access Denied
            </CardTitle>
            <CardDescription>
              Only owners can access system updates
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Button onClick={() => navigate('/')}>Go to Dashboard</Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background">
      <div className="py-6 border-b border-border">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-foreground">
            System Updates
          </h1>
          <p className="text-muted-foreground mt-1">
            Check and install updates for your Mist instance
          </p>
        </div>
      </div>

      <div className="py-6 max-w-4xl">
        <Tabs defaultValue="updates" className="w-full">
          <TabsList className="mb-6">
            <TabsTrigger value="updates">
              <Download className="h-4 w-4 mr-2" />
              Updates
            </TabsTrigger>
            <TabsTrigger value="history">
              <History className="h-4 w-4 mr-2" />
              History
            </TabsTrigger>
          </TabsList>

          <TabsContent value="updates">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Download className="h-5 w-5 text-primary" />
                  Update Status
                </CardTitle>
                <CardDescription>
                  Current version and available updates
                </CardDescription>
              </CardHeader>
              <CardContent>
                {isLoading ? (
                  <div className="flex items-center justify-center py-8">
                    <p className="text-muted-foreground">Checking for updates...</p>
                  </div>
                ) : (
                  <div className="space-y-6">
                    {error && (
                      <Alert variant="destructive">
                        <AlertCircle className="h-4 w-4" />
                        <AlertDescription>{error}</AlertDescription>
                      </Alert>
                    )}

                    {updateInfo && (
                      <>
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                          <div className="p-4 border rounded-lg">
                            <p className="text-sm text-muted-foreground mb-1">Current Version</p>
                            <p className="text-2xl font-bold">{updateInfo.currentVersion}</p>
                          </div>
                          <div className="p-4 border rounded-lg">
                            <p className="text-sm text-muted-foreground mb-1">Latest Version</p>
                            <div className="flex items-center gap-2">
                              <p className="text-2xl font-bold">{updateInfo.latestVersion}</p>
                              {updateInfo.updateAvailable && (
                                <Badge variant="default">New</Badge>
                              )}
                            </div>
                          </div>
                        </div>

                        {updateInfo.updateAvailable ? (
                          <Alert>
                            <CheckCircle2 className="h-4 w-4" />
                            <AlertDescription>
                              <strong>Update Available!</strong> A new version of Mist is ready to install.
                            </AlertDescription>
                          </Alert>
                        ) : (
                          <Alert>
                            <CheckCircle2 className="h-4 w-4" />
                            <AlertDescription>
                              You are running the latest version of Mist.
                            </AlertDescription>
                          </Alert>
                        )}

                        {updateInfo.releaseNotes && (
                          <div className="space-y-2">
                            <h3 className="text-lg font-semibold">Release Notes</h3>
                            {updateInfo.releaseName && (
                              <p className="text-sm font-medium text-muted-foreground">
                                {updateInfo.releaseName}
                              </p>
                            )}
                            <div className="p-4 bg-muted rounded-lg">
                              <pre className="text-sm whitespace-pre-wrap font-mono overflow-auto max-h-64">
                                {updateInfo.releaseNotes}
                              </pre>
                            </div>
                          </div>
                        )}

                        {updateInfo.updateAvailable && (
                          <Alert>
                            <AlertCircle className="h-4 w-4" />
                            <AlertDescription>
                              <strong>Before updating:</strong>
                              <ul className="mt-2 space-y-1 text-sm list-disc list-inside">
                                <li>Database backup will be created automatically</li>
                                <li>The dashboard will be unavailable for 1-2 minutes</li>
                                <li>Your applications will continue running</li>
                                <li>Automatic rollback will occur if update fails</li>
                                <li>Review the release notes above before proceeding</li>
                              </ul>
                            </AlertDescription>
                          </Alert>
                        )}

                        <div className="flex flex-col sm:flex-row items-start sm:items-center gap-2 pt-4">
                          <Button
                            onClick={triggerUpdate}
                            disabled={!updateInfo.updateAvailable || isUpdating || isCheckingUpdate}
                            className="w-full sm:w-auto"
                          >
                            {isUpdating ? (
                              <>
                                <span className="animate-spin mr-2">⏳</span>
                                Updating...
                              </>
                            ) : (
                              <>
                                <Download className="h-4 w-4 mr-2" />
                                Install Update
                              </>
                            )}
                          </Button>
                          <Button
                            variant="outline"
                            onClick={checkForUpdates}
                            disabled={isUpdating || isCheckingUpdate}
                            className="w-full sm:w-auto"
                          >
                            {isCheckingUpdate ? (
                              <>
                                <span className="animate-spin mr-2">⏳</span>
                                Checking...
                              </>
                            ) : (
                              <>
                                <RefreshCw className="h-4 w-4 mr-2" />
                                Check Again
                              </>
                            )}
                          </Button>
                        </div>

                        {isUpdating && updateLogs.length > 0 && (
                          <div className="space-y-2">
                            <h3 className="text-lg font-semibold flex items-center gap-2">
                              <span className="animate-spin">⏳</span>
                              Update in Progress
                            </h3>
                            <Alert>
                              <AlertCircle className="h-4 w-4" />
                              <AlertDescription>
                                Please do not close this window or refresh the page until the update completes.
                              </AlertDescription>
                            </Alert>
                            <div className="p-4 bg-black text-green-400 rounded-lg font-mono text-sm overflow-auto max-h-96">
                              {updateLogs.map((log, index) => (
                                <div key={index}>{log}</div>
                              ))}
                            </div>
                          </div>
                        )}

                        {!isUpdating && (
                          <Alert>
                            <AlertCircle className="h-4 w-4" />
                            <AlertDescription>
                              <strong>Safety Features:</strong> Automatic database backup, git rollback tags,
                              health checks, and automatic rollback on failure are all enabled.
                              If something goes wrong, the system will automatically revert to the previous version.
                            </AlertDescription>
                          </Alert>
                        )}
                      </>
                    )}
                  </div>
                )}
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="history">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <History className="h-5 w-5 text-primary" />
                  Update History
                </CardTitle>
                <CardDescription>
                  View past system updates and their status
                </CardDescription>
              </CardHeader>
              <CardContent>
                {isLoadingHistory ? (
                  <div className="flex items-center justify-center py-8">
                    <p className="text-muted-foreground">Loading history...</p>
                  </div>
                ) : updateHistory.length === 0 ? (
                  <div className="flex flex-col items-center justify-center py-12">
                    <History className="h-12 w-12 text-muted-foreground mb-4" />
                    <p className="text-muted-foreground">No update history available</p>
                  </div>
                ) : (
                  <div className="space-y-4">
                    {updateHistory.map((log) => (
                      <div
                        key={log.id}
                        className="border rounded-lg p-4 hover:bg-muted/50 transition-colors"
                      >
                        <div className="flex items-start justify-between mb-2">
                          <div className="flex items-center gap-2">
                            <h4 className="font-semibold">
                              {log.versionFrom} → {log.versionTo}
                            </h4>
                            <Badge
                              variant={
                                log.status === 'success'
                                  ? 'default'
                                  : log.status === 'failed'
                                    ? 'destructive'
                                    : 'secondary'
                              }
                            >
                              {log.status}
                            </Badge>
                          </div>
                          <div className="flex items-center gap-2 text-sm text-muted-foreground">
                            <Clock className="h-4 w-4" />
                            {new Date(log.startedAt).toLocaleString()}
                          </div>
                        </div>

                        <p className="text-sm text-muted-foreground mb-2">
                          Started by {log.username}
                        </p>

                        {log.completedAt && (
                          <p className="text-sm text-muted-foreground mb-2">
                            Completed: {new Date(log.completedAt).toLocaleString()}
                          </p>
                        )}

                        {log.errorMessage && (
                          <Alert variant="destructive" className="mb-2">
                            <AlertCircle className="h-4 w-4" />
                            <AlertDescription>{log.errorMessage}</AlertDescription>
                          </Alert>
                        )}

                        {log.logs && (
                          <details className="mt-2">
                            <summary className="cursor-pointer text-sm text-primary hover:underline">
                              View Logs
                            </summary>
                            <div className="mt-2 p-3 bg-black text-green-400 rounded-lg font-mono text-xs overflow-auto max-h-64">
                              <pre className="whitespace-pre-wrap">{log.logs}</pre>
                            </div>
                          </details>
                        )}
                      </div>
                    ))}
                  </div>
                )}
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
};
