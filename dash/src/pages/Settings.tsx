import { useState, useEffect } from 'react';
import { useAuth } from '@/providers';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { AlertCircle, CheckCircle2, Globe, ShieldAlert, Shield, HardDrive, Trash2 } from 'lucide-react';
import { settingsService } from '@/services';
import { toast } from 'sonner';
import { useNavigate } from 'react-router-dom';
import { Switch } from '@/components/ui/switch';
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs';

export const SettingsPage = () => {
  const { user } = useAuth();
  const navigate = useNavigate();
  const [isUpdatingSystemSettings, setIsUpdatingSystemSettings] = useState(false);

  const [wildcardDomain, setWildcardDomain] = useState('');
  const [mistAppName, setMistAppName] = useState('mist');
  const [systemSettingsError, setSystemSettingsError] = useState('');
  const [isLoadingSystemSettings, setIsLoadingSystemSettings] = useState(true);

  const [allowedOrigins, setAllowedOrigins] = useState('');
  const [productionMode, setProductionMode] = useState(false);
  const [secureCookies, setSecureCookies] = useState(false);
  const [isUpdatingSecuritySettings, setIsUpdatingSecuritySettings] = useState(false);
  const [securitySettingsError, setSecuritySettingsError] = useState('');

  const [autoCleanupContainers, setAutoCleanupContainers] = useState(false);
  const [autoCleanupImages, setAutoCleanupImages] = useState(false);
  const [isUpdatingDockerSettings, setIsUpdatingDockerSettings] = useState(false);
  const [dockerSettingsError, setDockerSettingsError] = useState('');
  const [isCleaningDocker, setIsCleaningDocker] = useState<string | null>(null);

  useEffect(() => {
    if (!user) {
      navigate('/');
      return;
    }

    if (user.role !== 'owner') {
      toast.error('Only owners can access system settings');
      navigate('/');
      return;
    }

    loadSystemSettings();
  }, [user, navigate]);

  const loadSystemSettings = async () => {
    try {
      const settings = await settingsService.getSystemSettings();
      setWildcardDomain(settings.wildcardDomain || '');
      setMistAppName(settings.mistAppName);
      setAllowedOrigins(settings.allowedOrigins || '');
      setProductionMode(settings.productionMode || false);
      setSecureCookies(settings.secureCookies || false);
      setAutoCleanupContainers(settings.autoCleanupContainers || false);
      setAutoCleanupImages(settings.autoCleanupImages || false);
    } catch (error) {
      console.error('Failed to load system settings:', error);
      toast.error('Failed to load system settings');
    } finally {
      setIsLoadingSystemSettings(false);
    }
  };

  const handleUpdateSystemSettings = async (e: React.FormEvent) => {
    e.preventDefault();
    setSystemSettingsError('');

    if (!mistAppName.trim()) {
      setSystemSettingsError('Mist app name is required');
      return;
    }

    setIsUpdatingSystemSettings(true);

    try {
      const settings = await settingsService.updateSystemSettings({
        wildcardDomain: wildcardDomain.trim() || '',
        mistAppName: mistAppName.trim(),
      });
      setWildcardDomain(settings.wildcardDomain || '');
      setMistAppName(settings.mistAppName);
      toast.success('System settings updated successfully');
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to update system settings';
      setSystemSettingsError(message);
      toast.error(message);
    } finally {
      setIsUpdatingSystemSettings(false);
    }
  };

  const handleUpdateSecuritySettings = async (e: React.FormEvent) => {
    e.preventDefault();
    setSecuritySettingsError('');
    setIsUpdatingSecuritySettings(true);

    try {
      const settings = await settingsService.updateSystemSettings({
        allowedOrigins: allowedOrigins.trim(),
        productionMode,
        secureCookies,
      });
      setAllowedOrigins(settings.allowedOrigins || '');
      setProductionMode(settings.productionMode || false);
      setSecureCookies(settings.secureCookies || false);
      toast.success('Security settings updated successfully');
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to update security settings';
      setSecuritySettingsError(message);
      toast.error(message);
    } finally {
      setIsUpdatingSecuritySettings(false);
    }
  };

  const handleUpdateDockerSettings = async (e: React.FormEvent) => {
    e.preventDefault();
    setDockerSettingsError('');
    setIsUpdatingDockerSettings(true);

    try {
      const settings = await settingsService.updateSystemSettings({
        autoCleanupContainers,
        autoCleanupImages,
      });
      setAutoCleanupContainers(settings.autoCleanupContainers || false);
      setAutoCleanupImages(settings.autoCleanupImages || false);
      toast.success('Docker settings updated successfully');
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to update Docker settings';
      setDockerSettingsError(message);
      toast.error(message);
    } finally {
      setIsUpdatingDockerSettings(false);
    }
  };

  const handleDockerCleanup = async (type: 'containers' | 'images' | 'system' | 'system-all') => {
    if (type === 'system-all' && !confirm('WARNING: This will remove ALL unused Docker images. This action cannot be undone. Continue?')) {
      return;
    }

    setIsCleaningDocker(type);
    try {
      const result = await settingsService.dockerCleanup(type);
      toast.success(result.message || 'Docker cleanup completed');
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Docker cleanup failed';
      toast.error(message);
    } finally {
      setIsCleaningDocker(null);
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
              Only owners can access system settings
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
            System Settings
          </h1>
          <p className="text-muted-foreground mt-1">
            Configure system-wide settings for your Mist instance
          </p>
        </div>
      </div>

      {/* Content with Tabs */}
      <div className="py-6 max-w-4xl">
        <Tabs defaultValue="general" className="w-full">
          <div className="w-full overflow-x-auto mb-6 pb-1">
            <TabsList className="inline-flex w-full min-w-fit">
            <TabsTrigger value="general">
              <Globe className="h-4 w-4 mr-2" />
              General
            </TabsTrigger>
            <TabsTrigger value="security">
              <Shield className="h-4 w-4 mr-2" />
              Security
            </TabsTrigger>
            <TabsTrigger value="docker">
              <HardDrive className="h-4 w-4 mr-2" />
              Docker
            </TabsTrigger>
            </TabsList>
          </div>

          {/* General Settings Tab */}
          <TabsContent value="general" className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Globe className="h-5 w-5 text-primary" />
                  Wildcard Domain Configuration
                </CardTitle>
                <CardDescription>
                  Configure wildcard domain for automatic app domain generation
                </CardDescription>
              </CardHeader>
              <CardContent>
                {isLoadingSystemSettings ? (
                  <div className="flex items-center justify-center py-8">
                    <p className="text-muted-foreground">Loading settings...</p>
                  </div>
                ) : (
                  <form onSubmit={handleUpdateSystemSettings} className="space-y-6">
                    {systemSettingsError && (
                      <Alert variant="destructive">
                        <AlertCircle className="h-4 w-4" />
                        <AlertDescription>{systemSettingsError}</AlertDescription>
                      </Alert>
                    )}

                    <div className="space-y-2">
                      <Label htmlFor="wildcardDomain">Wildcard Domain</Label>
                      <Input
                        id="wildcardDomain"
                        type="text"
                        value={wildcardDomain}
                        onChange={(e) => setWildcardDomain(e.target.value)}
                        placeholder="*.exam.ple or exam.ple"
                        disabled={isUpdatingSystemSettings}
                      />
                      <p className="text-sm text-muted-foreground">
                        When configured, apps will automatically get domains like <code className="bg-muted px-1 py-0.5 rounded">project-app.exam.ple</code>
                      </p>
                      <div className="mt-3 p-3 bg-muted rounded-md">
                        <p className="text-sm font-medium mb-2">Example:</p>
                        <ul className="text-sm text-muted-foreground space-y-1 list-disc list-inside">
                          <li>Wildcard domain: <code className="bg-background px-1 py-0.5 rounded">exam.ple</code></li>
                          <li>Project name: <code className="bg-background px-1 py-0.5 rounded">crux</code></li>
                          <li>App name: <code className="bg-background px-1 py-0.5 rounded">main</code></li>
                          <li>Generated domain: <code className="bg-background px-1 py-0.5 rounded">crux-main.exam.ple</code></li>
                        </ul>
                      </div>
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="mistAppName">Mist App Name</Label>
                      <Input
                        id="mistAppName"
                        type="text"
                        value={mistAppName}
                        onChange={(e) => setMistAppName(e.target.value)}
                        placeholder="mist"
                        disabled={isUpdatingSystemSettings}
                      />
                      <p className="text-sm text-muted-foreground">
                        Subdomain name for the Mist dashboard. With wildcard domain <code className="bg-muted px-1 py-0.5 rounded">exam.ple</code> and name <code className="bg-muted px-1 py-0.5 rounded">mist</code>,
                        Mist will be available at <code className="bg-muted px-1 py-0.5 rounded">mist.exam.ple</code>
                      </p>
                    </div>

                    <div className="flex flex-col sm:flex-row items-start sm:items-center gap-2 pt-4">
                      <Button
                        type="submit"
                        disabled={isUpdatingSystemSettings || !mistAppName.trim()}
                        className="w-full sm:w-auto"
                      >
                        {isUpdatingSystemSettings ? (
                          <>
                            <span className="animate-spin mr-2">⏳</span>
                            Updating...
                          </>
                        ) : (
                          <>
                            <CheckCircle2 className="h-4 w-4 mr-2" />
                            Update Settings
                          </>
                        )}
                      </Button>
                      <Button
                        type="button"
                        variant="outline"
                        onClick={loadSystemSettings}
                        disabled={isUpdatingSystemSettings}
                        className="w-full sm:w-auto"
                      >
                        Reset
                      </Button>
                    </div>
                  </form>
                )}
              </CardContent>
            </Card>
          </TabsContent>

          {/* Security Settings Tab */}
          <TabsContent value="security" className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Shield className="h-5 w-5 text-primary" />
                  Security Settings
                </CardTitle>
                <CardDescription>
                  Configure security settings
                </CardDescription>
              </CardHeader>
              <CardContent>
                {isLoadingSystemSettings ? (
                  <div className="flex items-center justify-center py-8">
                    <p className="text-muted-foreground">Loading settings...</p>
                  </div>
                ) : (
                  <form onSubmit={handleUpdateSecuritySettings} className="space-y-6">
                    {securitySettingsError && (
                      <Alert variant="destructive">
                        <AlertCircle className="h-4 w-4" />
                        <AlertDescription>{securitySettingsError}</AlertDescription>
                      </Alert>
                    )}

                    {/* <div className="space-y-2"> */}
                    {/*   <Label htmlFor="allowedOrigins">Allowed Origins (CORS)</Label> */}
                    {/*   <Input */}
                    {/*     id="allowedOrigins" */}
                    {/*     type="text" */}
                    {/*     value={allowedOrigins} */}
                    {/*     onChange={(e) => setAllowedOrigins(e.target.value)} */}
                    {/*     placeholder="https://mist.exam.ple,https://app.exam.ple" */}
                    {/*     disabled={isUpdatingSecuritySettings} */}
                    {/*   /> */}
                    {/*   <p className="text-sm text-muted-foreground"> */}
                    {/*     Comma-separated list of allowed origins for cross-origin WebSocket connections. Same-origin requests are always allowed. */}
                    {/*   </p> */}
                    {/*   <div className="mt-3 p-3 bg-muted rounded-md"> */}
                    {/*     <p className="text-sm font-medium mb-2">Examples:</p> */}
                    {/*     <ul className="text-sm text-muted-foreground space-y-1 list-disc list-inside"> */}
                    {/*       <li><code className="bg-background px-1 py-0.5 rounded">https://mist.exam.ple</code></li> */}
                    {/*       <li><code className="bg-background px-1 py-0.5 rounded">https://mist.exam.ple,https://app.exam.ple</code></li> */}
                    {/*     </ul> */}
                    {/*     <p className="text-sm text-muted-foreground mt-2"> */}
                    {/*       <strong>Note:</strong> Same-origin requests are always allowed automatically. Only add origins here if you need to allow cross-origin WebSocket connections. */}
                    {/*     </p> */}
                    {/*   </div> */}
                    {/* </div> */}

                    <div className="space-y-4">
                      <div className="flex items-center justify-between">
                        <div className="space-y-0.5">
                          <Label htmlFor="productionMode">Production Mode</Label>
                          <p className="text-sm text-muted-foreground">
                            Enable production mode for stricter security checks
                          </p>
                        </div>
                        <Switch
                          id="productionMode"
                          checked={productionMode}
                          onCheckedChange={setProductionMode}
                          disabled={isUpdatingSecuritySettings}
                        />
                      </div>

                      <div className="flex items-center justify-between">
                        <div className="space-y-0.5">
                          <Label htmlFor="secureCookies">Secure Cookies (HTTPS Only)</Label>
                          <p className="text-sm text-muted-foreground">
                            Enable secure flag on cookies. Only enable if using HTTPS.
                          </p>
                        </div>
                        <Switch
                          id="secureCookies"
                          checked={secureCookies}
                          onCheckedChange={setSecureCookies}
                          disabled={isUpdatingSecuritySettings}
                        />
                      </div>
                    </div>

                    <Alert>
                      <AlertCircle className="h-4 w-4" />
                      <AlertDescription>
                        <strong>Important:</strong> Only enable "Secure Cookies" if your Mist instance is running behind HTTPS.
                        Enabling this without HTTPS will prevent users from logging in.
                      </AlertDescription>
                    </Alert>

                    <div className="flex flex-col sm:flex-row items-start sm:items-center gap-2 pt-4">
                      <Button
                        type="submit"
                        disabled={isUpdatingSecuritySettings}
                        className="w-full sm:w-auto"
                      >
                        {isUpdatingSecuritySettings ? (
                          <>
                            <span className="animate-spin mr-2">⏳</span>
                            Updating...
                          </>
                        ) : (
                          <>
                            <CheckCircle2 className="h-4 w-4 mr-2" />
                            Update Settings
                          </>
                        )}
                      </Button>
                      <Button
                        type="button"
                        variant="outline"
                        onClick={loadSystemSettings}
                        disabled={isUpdatingSecuritySettings}
                        className="w-full sm:w-auto"
                      >
                        Reset
                      </Button>
                    </div>
                  </form>
                )}
              </CardContent>
            </Card>
          </TabsContent>

          {/* Docker Cleanup Settings Tab */}
          <TabsContent value="docker" className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <HardDrive className="h-5 w-5 text-primary" />
                  Docker Cleanup Settings
                </CardTitle>
                <CardDescription>
                  Configure automatic Docker cleanup after deployments and manual cleanup options
                </CardDescription>
              </CardHeader>
              <CardContent>
                {isLoadingSystemSettings ? (
                  <div className="flex items-center justify-center py-8">
                    <p className="text-muted-foreground">Loading settings...</p>
                  </div>
                ) : (
                  <div className="space-y-6">
                    <form onSubmit={handleUpdateDockerSettings} className="space-y-6">
                      {dockerSettingsError && (
                        <Alert variant="destructive">
                          <AlertCircle className="h-4 w-4" />
                          <AlertDescription>{dockerSettingsError}</AlertDescription>
                        </Alert>
                      )}

                      <div className="space-y-4">
                        <div className="flex items-center justify-between">
                          <div className="space-y-0.5">
                            <Label htmlFor="autoCleanupContainers">Auto-Cleanup Stopped Containers</Label>
                            <p className="text-sm text-muted-foreground">
                              Automatically remove stopped containers after each deployment
                            </p>
                          </div>
                          <Switch
                            id="autoCleanupContainers"
                            checked={autoCleanupContainers}
                            onCheckedChange={setAutoCleanupContainers}
                            disabled={isUpdatingDockerSettings}
                          />
                        </div>

                        <div className="flex items-center justify-between">
                          <div className="space-y-0.5">
                            <Label htmlFor="autoCleanupImages">Auto-Cleanup Dangling Images</Label>
                            <p className="text-sm text-muted-foreground">
                              Automatically remove dangling (untagged) images after each deployment
                            </p>
                          </div>
                          <Switch
                            id="autoCleanupImages"
                            checked={autoCleanupImages}
                            onCheckedChange={setAutoCleanupImages}
                            disabled={isUpdatingDockerSettings}
                          />
                        </div>
                      </div>

                      <div className="flex flex-col sm:flex-row items-start sm:items-center gap-2 pt-4 border-t">
                        <Button
                          type="submit"
                          disabled={isUpdatingDockerSettings}
                          className="w-full sm:w-auto"
                        >
                          {isUpdatingDockerSettings ? (
                            <>
                              <span className="animate-spin mr-2">⏳</span>
                              Updating...
                            </>
                          ) : (
                            <>
                              <CheckCircle2 className="h-4 w-4 mr-2" />
                              Update Settings
                            </>
                          )}
                        </Button>
                        <Button
                          type="button"
                          variant="outline"
                          onClick={loadSystemSettings}
                          disabled={isUpdatingDockerSettings}
                          className="w-full sm:w-auto"
                        >
                          Reset
                        </Button>
                      </div>
                    </form>

                    {/* Manual Cleanup Section */}
                    <div className="pt-6 border-t">
                      <h3 className="text-lg font-semibold mb-2 flex items-center gap-2">
                        <Trash2 className="h-5 w-5" />
                        Manual Docker Cleanup
                      </h3>
                      <p className="text-sm text-muted-foreground mb-4">
                        Run manual cleanup operations to free up disk space
                      </p>

                      <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                        <Button
                          variant="outline"
                          onClick={() => handleDockerCleanup('containers')}
                          disabled={isCleaningDocker !== null}
                          className="w-full"
                        >
                          {isCleaningDocker === 'containers' ? (
                            <>
                              <span className="animate-spin mr-2">⏳</span>
                              Cleaning...
                            </>
                          ) : (
                            <>
                              <Trash2 className="h-4 w-4 mr-2" />
                              Clean Stopped Containers
                            </>
                          )}
                        </Button>

                        <Button
                          variant="outline"
                          onClick={() => handleDockerCleanup('images')}
                          disabled={isCleaningDocker !== null}
                          className="w-full"
                        >
                          {isCleaningDocker === 'images' ? (
                            <>
                              <span className="animate-spin mr-2">⏳</span>
                              Cleaning...
                            </>
                          ) : (
                            <>
                              <Trash2 className="h-4 w-4 mr-2" />
                              Clean Dangling Images
                            </>
                          )}
                        </Button>

                        <Button
                          variant="outline"
                          onClick={() => handleDockerCleanup('system')}
                          disabled={isCleaningDocker !== null}
                          className="w-full"
                        >
                          {isCleaningDocker === 'system' ? (
                            <>
                              <span className="animate-spin mr-2">⏳</span>
                              Cleaning...
                            </>
                          ) : (
                            <>
                              <Trash2 className="h-4 w-4 mr-2" />
                              System Prune
                            </>
                          )}
                        </Button>

                        <Button
                          variant="destructive"
                          onClick={() => handleDockerCleanup('system-all')}
                          disabled={isCleaningDocker !== null}
                          className="w-full"
                        >
                          {isCleaningDocker === 'system-all' ? (
                            <>
                              <span className="animate-spin mr-2">⏳</span>
                              Cleaning...
                            </>
                          ) : (
                            <>
                              <Trash2 className="h-4 w-4 mr-2" />
                              Aggressive Prune (All Unused)
                            </>
                          )}
                        </Button>
                      </div>

                      <Alert className="mt-4">
                        <AlertCircle className="h-4 w-4" />
                        <AlertDescription>
                          <strong>Cleanup Types:</strong>
                          <ul className="mt-2 space-y-1 text-sm list-disc list-inside">
                            <li><strong>Stopped Containers:</strong> Removes only stopped containers</li>
                            <li><strong>Dangling Images:</strong> Removes untagged images</li>
                            <li><strong>System Prune:</strong> Removes stopped containers, unused networks, dangling images, and build cache</li>
                            <li><strong>Aggressive Prune:</strong> Removes ALL unused Docker resources including unused images (use with caution!)</li>
                          </ul>
                        </AlertDescription>
                      </Alert>
                    </div>
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
