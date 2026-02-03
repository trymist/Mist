import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { Checkbox } from "@/components/ui/checkbox";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { AlertTriangle, RefreshCw, Play } from "lucide-react";
import { toast } from "sonner";
import { applicationsService } from "@/services";
import type { App, RestartPolicy } from "@/types";

interface AppSettingsProps {
  app: App;
  onUpdate: () => void;
}

export const AppSettings = ({ app, onUpdate }: AppSettingsProps) => {
  const [port, setPort] = useState(app.port?.toString() || "");
  // Default to true for web apps, false for others (matches backend defaults)
  const defaultShouldExpose = app.appType === 'web' ? true : false;
  const [shouldExpose, setShouldExpose] = useState(app.shouldExpose ?? defaultShouldExpose);
  const [exposePort, setExposePort] = useState(app.exposePort?.toString() || "");
  const [buildCommand, setBuildCommand] = useState(app.buildCommand || "");
  const [startCommand, setStartCommand] = useState(app.startCommand || "");
  const [rootDirectory, setRootDirectory] = useState(app.rootDirectory || "");
  const [dockerfilePath, setDockerfilePath] = useState(app.dockerfilePath || "");
  const [healthcheckPath, setHealthcheckPath] = useState(app.healthcheckPath || "");
  const [cpuLimit, setCpuLimit] = useState(app.cpuLimit?.toString() || "");
  const [memoryLimit, setMemoryLimit] = useState(app.memoryLimit?.toString() || "");
  const [restartPolicy, setRestartPolicy] = useState(app.restartPolicy || "unless-stopped");
  const [deploymentStrategy, setDeploymentStrategy] = useState(app.deploymentStrategy || "auto");
  const [saving, setSaving] = useState(false);
  const [actionDialogOpen, setActionDialogOpen] = useState(false);
  const [actionRequired, setActionRequired] = useState<'restart' | 'redeploy' | null>(null);
  const [actionMessage, setActionMessage] = useState('');
  const [pendingUpdates, setPendingUpdates] = useState<any>(null);
  const [isRestarting, setIsRestarting] = useState(false);

  // Track original values to detect changes
  const [originalValues] = useState({
    port: app.port,
    shouldExpose: app.shouldExpose,
    exposePort: app.exposePort,
    rootDirectory: app.rootDirectory,
    dockerfilePath: app.dockerfilePath,
    buildCommand: app.buildCommand,
    startCommand: app.startCommand,
    healthcheckPath: app.healthcheckPath,
    restartPolicy: app.restartPolicy,
    deploymentStrategy: app.deploymentStrategy,
    cpuLimit: app.cpuLimit,
    memoryLimit: app.memoryLimit,
  });

  // Determine if port exposure settings should be shown (for all app types except compose)
  const showPortExposureSettings = app.appType !== 'compose';

  const detectChanges = () => {
    const changes: string[] = [];
    
    if (port !== originalValues.port?.toString()) changes.push('port');
    if (shouldExpose !== (originalValues.shouldExpose ?? true)) changes.push('shouldExpose');
    if (exposePort !== originalValues.exposePort?.toString()) changes.push('exposePort');
    if (rootDirectory !== originalValues.rootDirectory) changes.push('rootDirectory');
    if (dockerfilePath !== (originalValues.dockerfilePath || '')) changes.push('dockerfilePath');
    if (buildCommand !== (originalValues.buildCommand || '')) changes.push('buildCommand');
    if (startCommand !== (originalValues.startCommand || '')) changes.push('startCommand');
    if (restartPolicy !== originalValues.restartPolicy) changes.push('restartPolicy');
    if (parseFloat(cpuLimit || '0') !== (originalValues.cpuLimit || 0)) changes.push('cpuLimit');
    if (parseInt(memoryLimit || '0') !== (originalValues.memoryLimit || 0)) changes.push('memoryLimit');
    
    return changes;
  };

  const handleSave = async () => {
    try {
      setSaving(true);

      const updates: Partial<{
        rootDirectory: string;
        dockerfilePath: string | null;
        buildCommand: string | null;
        startCommand: string | null;
        healthcheckPath: string | null;
        restartPolicy: RestartPolicy;
        deploymentStrategy: string;
        port: number;
        shouldExpose: boolean;
        exposePort: number | null;
        cpuLimit: number | null;
        memoryLimit: number | null;
      }> = {};

      // Only include fields that have actually changed
      if (rootDirectory !== (originalValues.rootDirectory || '')) {
        updates.rootDirectory = rootDirectory;
      }
      if (dockerfilePath !== (originalValues.dockerfilePath || '')) {
        updates.dockerfilePath = dockerfilePath || null;
      }
      if (buildCommand !== (originalValues.buildCommand || '')) {
        updates.buildCommand = buildCommand || null;
      }
      if (startCommand !== (originalValues.startCommand || '')) {
        updates.startCommand = startCommand || null;
      }
      if (healthcheckPath !== (originalValues.healthcheckPath || '')) {
        updates.healthcheckPath = healthcheckPath || null;
      }
      if (restartPolicy !== originalValues.restartPolicy) {
        updates.restartPolicy = restartPolicy as RestartPolicy;
      }
      if (deploymentStrategy !== originalValues.deploymentStrategy) {
        updates.deploymentStrategy = deploymentStrategy;
      }
      if (shouldExpose !== (originalValues.shouldExpose ?? defaultShouldExpose)) {
        updates.shouldExpose = shouldExpose;
      }

      // Port handling
      const currentPort = port ? parseInt(port) : null;
      if (currentPort !== originalValues.port) {
        if (port) {
          const portNum = parseInt(port);
          if (isNaN(portNum) || portNum < 1 || portNum > 65535) {
            toast.error("Port must be a number between 1 and 65535");
            return;
          }
          updates.port = portNum;
        }
      }

      // Expose port handling - only include if shouldExpose changed or exposePort changed
      const currentExposePort = exposePort ? parseInt(exposePort) : null;
      if (shouldExpose !== (originalValues.shouldExpose ?? defaultShouldExpose) || 
          currentExposePort !== originalValues.exposePort) {
        if (shouldExpose) {
          if (exposePort) {
            const exposePortNum = parseInt(exposePort);
            if (isNaN(exposePortNum) || exposePortNum < 1 || exposePortNum > 65535) {
              toast.error("External port must be a number between 1 and 65535");
              return;
            }
            updates.exposePort = exposePortNum;
          } else {
            updates.exposePort = null;
          }
        } else {
          updates.exposePort = null;
        }
      }

      // CPU limit handling
      const currentCpuLimit = cpuLimit ? parseFloat(cpuLimit) : null;
      if (currentCpuLimit !== (originalValues.cpuLimit || null)) {
        if (cpuLimit) {
          const cpu = parseFloat(cpuLimit);
          if (isNaN(cpu) || cpu <= 0) {
            toast.error("CPU limit must be a positive number");
            return;
          }
          updates.cpuLimit = cpu;
        } else {
          updates.cpuLimit = null;
        }
      }

      // Memory limit handling
      const currentMemoryLimit = memoryLimit ? parseInt(memoryLimit) : null;
      if (currentMemoryLimit !== (originalValues.memoryLimit || null)) {
        if (memoryLimit) {
          const memory = parseInt(memoryLimit);
          if (isNaN(memory) || memory <= 0) {
            toast.error("Memory limit must be a positive number");
            return;
          }
          updates.memoryLimit = memory;
        } else {
          updates.memoryLimit = null;
        }
      }

      // Check if there are any changes
      if (Object.keys(updates).length === 0) {
        toast.info("No changes to save");
        return;
      }

      // Save the updates
      const response = await applicationsService.update(app.id, updates);
      toast.success("Settings saved successfully");

      // Check backend response for action requirements
      if (response.actionRequired === 'redeploy') {
        setActionRequired('redeploy');
        setActionMessage(response.actionMessage || 'These build configuration changes require a full redeployment. Would you like to redeploy now?');
        setPendingUpdates(updates);
        setActionDialogOpen(true);
      } else if (response.actionRequired === 'restart') {
        setActionRequired('restart');
        setActionMessage(response.actionMessage || 'These runtime configuration changes require restarting the container. Would you like to restart now?');
        setPendingUpdates(updates);
        setActionDialogOpen(true);
      } else {
        // No action needed, just update
        onUpdate();
      }
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to update settings");
    } finally {
      setSaving(false);
    }
  };

  const handleRestart = async () => {
    try {
      setIsRestarting(true);
      await applicationsService.recreateContainer(app.id);
      toast.success("Container recreated successfully");
      setActionDialogOpen(false);
      onUpdate();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to recreate container");
    } finally {
      setIsRestarting(false);
    }
  };

  const handleRedeploy = async () => {
    try {
      setActionDialogOpen(false);
      toast.info("Starting full redeployment (building Docker image, this may take several minutes)...");
      await applicationsService.redeploy(app.id);
      toast.success("Redeployment started - monitoring deployment logs...");
      onUpdate();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to trigger redeployment");
    }
  };

  const handleSkip = () => {
    setActionDialogOpen(false);
    setActionRequired(null);
    setPendingUpdates(null);
    onUpdate();
    toast.info("Changes saved. They will take effect on next restart/redeployment.");
  };

  return (
    <Card>
      <CardHeader className="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-4">
        <div className="flex flex-col gap-2">
          <CardTitle>Application Settings</CardTitle>
          <CardDescription>
            Configure your application settings. Changes will be applied on next deployment.
          </CardDescription>
        </div>
        <div className="flex justify-start sm:justify-end w-full sm:w-auto">
          <Button onClick={handleSave} disabled={saving}>
            {saving ? "Saving..." : "Save Settings"}
          </Button>
        </div>
      </CardHeader>
      <CardContent className="space-y-6">
        {/* App Type Badge */}
        <div className="flex items-center gap-2 pb-4 border-b">
          <Label>Application Type:</Label>
          <Badge variant="secondary" className="capitalize">
            {app.appType || 'web'}
          </Badge>
          {app.templateName && (
            <>
              <span className="text-muted-foreground">•</span>
              <Badge variant="outline">{app.templateName}</Badge>
            </>
          )}
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="space-y-2">
            <Label htmlFor="port">Port</Label>
            <Input
              id="port"
              type="number"
              placeholder="3000"
              value={port}
              onChange={(e) => setPort(e.target.value)}
              min="1"
              max="65535"
              disabled={app.appType === 'database'}
            />
            <p className="text-sm text-muted-foreground">
              {app.appType === 'database'
                ? 'Port is managed by the template'
                : 'The port your application runs on'}
            </p>
          </div>

          {showPortExposureSettings && (
            <>
              <div className="space-y-3 p-4 border rounded-lg bg-muted/30 md:col-span-2">
                <div className="flex items-start space-x-3">
                  <Checkbox
                    id="shouldExpose"
                    checked={shouldExpose}
                    onCheckedChange={(checked) => setShouldExpose(checked as boolean)}
                  />
                  <div className="space-y-1">
                    <Label htmlFor="shouldExpose" className="font-medium cursor-pointer">
                      Expose External Port
                    </Label>
                    <p className="text-sm text-muted-foreground">
                      When enabled, exposes the application on an external port accessible via IP:port. 
                      Only applies when no custom domains are configured.
                    </p>
                  </div>
                </div>

                {shouldExpose && (
                  <div className="space-y-2 pt-2 pl-7">
                    <Label htmlFor="exposePort">External Port (Optional)</Label>
                    <Input
                      id="exposePort"
                      type="number"
                      placeholder={port || "3000"}
                      value={exposePort}
                      onChange={(e) => setExposePort(e.target.value)}
                      min="1"
                      max="65535"
                      className="max-w-xs"
                    />
                    <p className="text-sm text-muted-foreground">
                      Leave empty to use the same port as the container port ({port || "3000"})
                    </p>
                  </div>
                )}
              </div>
            </>
          )}

          <div className="space-y-2">
            <Label htmlFor="rootDirectory">Root Directory</Label>
            <Input
              id="rootDirectory"
              placeholder="/"
              value={rootDirectory}
              onChange={(e) => setRootDirectory(e.target.value)}
              disabled={app.appType === 'database'}
            />
            <p className="text-sm text-muted-foreground">
              {app.appType === 'database'
                ? 'Not applicable for database apps'
                : 'The root directory of your application'}
            </p>
          </div>

          <div className="space-y-2">
            <Label htmlFor="cpuLimit">CPU Limit (cores)</Label>
            <Input
              id="cpuLimit"
              type="number"
              placeholder="e.g., 0.5 or 2"
              value={cpuLimit}
              onChange={(e) => setCpuLimit(e.target.value)}
              min="0"
              step="0.1"
            />
            <p className="text-sm text-muted-foreground">
              Maximum CPU cores (leave empty for no limit)
            </p>
          </div>

          <div className="space-y-2">
            <Label htmlFor="memoryLimit">Memory Limit (MB)</Label>
            <Input
              id="memoryLimit"
              type="number"
              placeholder="e.g., 512 or 1024"
              value={memoryLimit}
              onChange={(e) => setMemoryLimit(e.target.value)}
              min="0"
            />
            <p className="text-sm text-muted-foreground">
              Maximum memory in MB (leave empty for no limit)
            </p>
          </div>

          <div className="space-y-2">
            <Label htmlFor="restartPolicy">Restart Policy</Label>
            <select
              id="restartPolicy"
              value={restartPolicy}
              onChange={(e) => setRestartPolicy(e.target.value as RestartPolicy)}
              className="w-full bg-background border rounded-md px-3 py-2"
            >
              <option value="no">No</option>
              <option value="always">Always</option>
              <option value="on-failure">On Failure</option>
              <option value="unless-stopped">Unless Stopped</option>
            </select>
            <p className="text-sm text-muted-foreground">
              When to restart the container
            </p>
          </div>

          <div className="space-y-2">
            <Label htmlFor="deploymentStrategy">Deployment Strategy</Label>
            <select
              id="deploymentStrategy"
              value={deploymentStrategy}
              onChange={(e) => setDeploymentStrategy(e.target.value)}
              className="w-full bg-background border rounded-md px-3 py-2"
              disabled={app.appType === 'database'}
            >
              <option value="auto">Automatic</option>
              <option value="manual">Manual</option>
            </select>
            <p className="text-sm text-muted-foreground">
              {app.appType === 'database'
                ? 'Deployment strategy managed by template'
                : 'Auto: Deploy on every push. Manual: Deploy only when triggered manually'}
            </p>
          </div>

          <div className="space-y-2">
            <Label htmlFor="dockerfilePath">Dockerfile Path</Label>
            <Input
              id="dockerfilePath"
              placeholder="Dockerfile"
              value={dockerfilePath}
              onChange={(e) => setDockerfilePath(e.target.value)}
              disabled={app.appType === 'database'}
            />
            <p className="text-sm text-muted-foreground">
              {app.appType === 'database'
                ? 'Not applicable for database apps'
                : 'Path to your Dockerfile (optional)'}
            </p>
          </div>

          {/* <div className="space-y-2">
            <Label htmlFor="healthcheckPath">Health Check Path</Label>
            <Input
              id="healthcheckPath"
              placeholder="/health"
              value={healthcheckPath}
              onChange={(e) => setHealthcheckPath(e.target.value)}
              disabled={app.appType === 'database'}
            />
            <p className="text-sm text-muted-foreground">
              {app.appType === 'database'
                ? 'Health checks managed by template'
                : 'Path for health checks (optional)'}
            </p>
          </div>

          <div className="space-y-2 md:col-span-2">
            <Label htmlFor="buildCommand">Build Command</Label>
            <Input
              id="buildCommand"
              placeholder="npm run build"
              value={buildCommand}
              onChange={(e) => setBuildCommand(e.target.value)}
              disabled={app.appType === 'database'}
            />
            <p className="text-sm text-muted-foreground">
              {app.appType === 'database'
                ? 'Not applicable for database apps'
                : 'Command to build your application (optional)'}
            </p>
          </div>

          <div className="space-y-2 md:col-span-2">
            <Label htmlFor="startCommand">Start Command</Label>
            <Input
              id="startCommand"
              placeholder="npm start"
              value={startCommand}
              onChange={(e) => setStartCommand(e.target.value)}
              disabled={app.appType === 'database'}
            />
            <p className="text-sm text-muted-foreground">
              {app.appType === 'database'
                ? 'Start command managed by template'
                : 'Command to start your application (optional)'}
            </p>
          </div> */}
        </div>

      </CardContent>

      {/* Action Required Dialog */}
      <Dialog open={actionDialogOpen} onOpenChange={(open) => !isRestarting && setActionDialogOpen(open)}>
        <DialogContent onPointerDownOutside={(e) => isRestarting && e.preventDefault()}>
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <AlertTriangle className="h-5 w-5 text-yellow-500" />
              Action Required
            </DialogTitle>
            <DialogDescription>
              {isRestarting ? "Restarting container, please wait..." : actionMessage}
            </DialogDescription>
          </DialogHeader>
          <DialogFooter className="flex gap-2 sm:justify-end">
            <Button variant="outline" onClick={handleSkip} disabled={isRestarting}>
              Skip for Now
            </Button>
            {actionRequired === 'restart' && (
              <Button onClick={handleRestart} disabled={isRestarting} className="flex items-center gap-2">
                <RefreshCw className={`h-4 w-4 ${isRestarting ? 'animate-spin' : ''}`} />
                {isRestarting ? 'Restarting...' : 'Restart Now'}
              </Button>
            )}
            {actionRequired === 'redeploy' && (
              <Button onClick={handleRedeploy} className="flex items-center gap-2">
                <Play className="h-4 w-4" />
                Redeploy Now
              </Button>
            )}
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </Card>
  );
};
