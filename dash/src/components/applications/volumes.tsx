import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Trash2, Plus, Pencil, X, Check, HardDrive, Info, AlertTriangle, RefreshCw } from "lucide-react";
import { toast } from "sonner";
import { applicationsService } from "@/services";
import type { Volume, CreateVolumeRequest, UpdateVolumeRequest } from "@/types";
import { Switch } from "@/components/ui/switch";
import { Alert, AlertDescription } from "@/components/ui/alert";

interface VolumesProps {
  appId: number;
  appType: string;
}

export const Volumes = ({ appId, appType }: VolumesProps) => {
  const [volumes, setVolumes] = useState<Volume[]>([]);
  const [loading, setLoading] = useState(true);
  const [showAddForm, setShowAddForm] = useState(false);

  const [newName, setNewName] = useState("");
  const [newHostPath, setNewHostPath] = useState("");
  const [newContainerPath, setNewContainerPath] = useState("");
  const [newReadOnly, setNewReadOnly] = useState(false);
  const [isAdding, setIsAdding] = useState(false);

  const [editingId, setEditingId] = useState<number | null>(null);
  const [editName, setEditName] = useState("");
  const [editHostPath, setEditHostPath] = useState("");
  const [editContainerPath, setEditContainerPath] = useState("");
  const [editReadOnly, setEditReadOnly] = useState(false);
  const [isUpdating, setIsUpdating] = useState(false);
  const [actionDialogOpen, setActionDialogOpen] = useState(false);

  const dbVolumePaths: Record<string, string> = {
    postgres: "/var/lib/postgresql/data",
    mysql: "/var/lib/mysql",
    mariadb: "/var/lib/mysql",
    mongodb: "/data/db",
    redis: "/data",
  };

  useEffect(() => {
    loadVolumes();
  }, [appId]);

  const loadVolumes = async () => {
    try {
      setLoading(true);
      const data = await applicationsService.getVolumes(appId);
      setVolumes(data);
    } catch (error) {
      console.error("Failed to load volumes:", error);
      toast.error("Failed to load volumes");
    } finally {
      setLoading(false);
    }
  };

  const showRestartDialog = () => {
    setActionDialogOpen(true);
  };

  const [isRestarting, setIsRestarting] = useState(false);
  const [restartError, setRestartError] = useState<string | null>(null);

  const handleRestart = async () => {
    try {
      setIsRestarting(true);
      setRestartError(null);
      await applicationsService.recreateContainer(appId);
      toast.success("Container recreated successfully - volume changes applied");
      setActionDialogOpen(false);
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : "Failed to recreate container";
      setRestartError(errorMessage);
    } finally {
      setIsRestarting(false);
    }
  };

  const handleSkipRestart = () => {
    setActionDialogOpen(false);
    toast.info("Volume changes will take effect on next restart or redeployment");
  };

  const handleAdd = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!newName.trim() || !newHostPath.trim() || !newContainerPath.trim()) {
      toast.error("All fields are required");
      return;
    }

    setIsAdding(true);
    try {
      const request: CreateVolumeRequest = {
        appId,
        name: newName.trim(),
        hostPath: newHostPath.trim(),
        containerPath: newContainerPath.trim(),
        readOnly: newReadOnly,
      };
      await applicationsService.createVolume(request);
      toast.success("Volume created successfully");
      setNewName("");
      setNewHostPath("");
      setNewContainerPath("");
      setNewReadOnly(false);
      setShowAddForm(false);
      loadVolumes();
      // Show restart dialog since volumes require container recreation
      showRestartDialog();
    } catch (error) {
      const message = error instanceof Error ? error.message : "Failed to create volume";
      toast.error(message);
    } finally {
      setIsAdding(false);
    }
  };

  const startEdit = (volume: Volume) => {
    setEditingId(volume.id);
    setEditName(volume.name);
    setEditHostPath(volume.hostPath);
    setEditContainerPath(volume.containerPath);
    setEditReadOnly(volume.readOnly);
  };

  const cancelEdit = () => {
    setEditingId(null);
    setEditName("");
    setEditHostPath("");
    setEditContainerPath("");
    setEditReadOnly(false);
  };

  const handleUpdate = async (volumeId: number) => {
    if (!editName.trim() || !editHostPath.trim() || !editContainerPath.trim()) {
      toast.error("All fields are required");
      return;
    }

    setIsUpdating(true);
    try {
      const request: UpdateVolumeRequest = {
        id: volumeId,
        name: editName.trim(),
        hostPath: editHostPath.trim(),
        containerPath: editContainerPath.trim(),
        readOnly: editReadOnly,
      };
      await applicationsService.updateVolume(request);
      toast.success("Volume updated successfully");
      cancelEdit();
      loadVolumes();
      // Show restart dialog since volumes require container recreation
      showRestartDialog();
    } catch (error) {
      const message = error instanceof Error ? error.message : "Failed to update volume";
      toast.error(message);
    } finally {
      setIsUpdating(false);
    }
  };

  const handleDelete = async (volumeId: number) => {
    if (!confirm("Are you sure you want to delete this volume? This will not delete the actual data on disk.")) {
      return;
    }

    try {
      await applicationsService.deleteVolume(volumeId);
      toast.success("Volume deleted successfully");
      loadVolumes();
      // Show restart dialog since volumes require container recreation
      showRestartDialog();
    } catch (error) {
      const message = error instanceof Error ? error.message : "Failed to delete volume";
      toast.error(message);
    }
  };

  const getSuggestedContainerPath = () => {
    const lowerAppType = appType.toLowerCase();
    for (const [db, path] of Object.entries(dbVolumePaths)) {
      if (lowerAppType.includes(db)) {
        return path;
      }
    }
    return "";
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <HardDrive className="h-5 w-5 text-primary" />
          Volumes
        </CardTitle>
        <CardDescription>
          Manage persistent storage volumes for your application
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        {appType === 'database' && (
          <Alert>
            <Info className="h-4 w-4" />
            <AlertDescription>
              <strong>Database volumes:</strong> Ensure data persists across container restarts by mounting a volume to the database data directory.
              {getSuggestedContainerPath() && (
                <> Suggested container path: <code className="bg-muted px-1 py-0.5 rounded">{getSuggestedContainerPath()}</code></>
              )}
            </AlertDescription>
          </Alert>
        )}

        {loading ? (
          <div className="text-center py-4 text-muted-foreground">Loading volumes...</div>
        ) : (
          <>
            {volumes.length === 0 && !showAddForm ? (
              <div className="text-center py-8 text-muted-foreground">
                <HardDrive className="h-12 w-12 mx-auto mb-2 opacity-50" />
                <p>No volumes configured</p>
                <p className="text-sm">Add a volume to persist data across deployments</p>
              </div>
            ) : (
              <div className="space-y-2">
                {volumes.map((volume) => (
                  <div
                    key={volume.id}
                    className="flex items-start gap-2 p-3 border border-border rounded-md bg-card"
                  >
                    {editingId === volume.id ? (
                      <div className="flex-1 space-y-3">
                        <div>
                          <Label htmlFor={`edit-name-${volume.id}`} className="text-xs">Name</Label>
                          <Input
                            id={`edit-name-${volume.id}`}
                            value={editName}
                            onChange={(e) => setEditName(e.target.value)}
                            placeholder="data"
                            disabled={isUpdating}
                          />
                        </div>
                        <div>
                          <Label htmlFor={`edit-host-${volume.id}`} className="text-xs">Host Path</Label>
                          <Input
                            id={`edit-host-${volume.id}`}
                            value={editHostPath}
                            onChange={(e) => setEditHostPath(e.target.value)}
                            placeholder="mist-volume-123"
                            disabled={isUpdating}
                          />
                          <p className="text-xs text-muted-foreground mt-1">
                            Docker volume name or absolute host path (e.g., /var/data/myapp)
                          </p>
                        </div>
                        <div>
                          <Label htmlFor={`edit-container-${volume.id}`} className="text-xs">Container Path</Label>
                          <Input
                            id={`edit-container-${volume.id}`}
                            value={editContainerPath}
                            onChange={(e) => setEditContainerPath(e.target.value)}
                            placeholder="/data"
                            disabled={isUpdating}
                          />
                        </div>
                        <div className="flex items-center gap-2">
                          <Switch
                            id={`edit-readonly-${volume.id}`}
                            checked={editReadOnly}
                            onCheckedChange={setEditReadOnly}
                            disabled={isUpdating}
                          />
                          <Label htmlFor={`edit-readonly-${volume.id}`} className="text-xs">
                            Read Only
                          </Label>
                        </div>
                        <div className="flex gap-2">
                          <Button
                            size="sm"
                            onClick={() => handleUpdate(volume.id)}
                            disabled={isUpdating}
                          >
                            <Check className="h-3 w-3 mr-1" />
                            Save
                          </Button>
                          <Button
                            size="sm"
                            variant="outline"
                            onClick={cancelEdit}
                            disabled={isUpdating}
                          >
                            <X className="h-3 w-3 mr-1" />
                            Cancel
                          </Button>
                        </div>
                      </div>
                    ) : (
                      <>
                        <div className="flex-1">
                          <div className="font-medium text-sm">{volume.name}</div>
                          <div className="text-xs text-muted-foreground space-y-1 mt-1">
                            <div>Host: <code className="bg-muted px-1 py-0.5 rounded">{volume.hostPath}</code></div>
                            <div>Container: <code className="bg-muted px-1 py-0.5 rounded">{volume.containerPath}</code></div>
                            {volume.readOnly && <div className="text-orange-500">Read Only</div>}
                          </div>
                        </div>
                        <div className="flex gap-1">
                          <Button
                            size="sm"
                            variant="ghost"
                            onClick={() => startEdit(volume)}
                          >
                            <Pencil className="h-3 w-3" />
                          </Button>
                          <Button
                            size="sm"
                            variant="ghost"
                            onClick={() => handleDelete(volume.id)}
                          >
                            <Trash2 className="h-3 w-3 text-destructive" />
                          </Button>
                        </div>
                      </>
                    )}
                  </div>
                ))}
              </div>
            )}

            {showAddForm ? (
              <form onSubmit={handleAdd} className="space-y-3 p-3 border border-border rounded-md bg-muted/50">
                <div>
                  <Label htmlFor="new-name" className="text-xs">Volume Name</Label>
                  <Input
                    id="new-name"
                    value={newName}
                    onChange={(e) => setNewName(e.target.value)}
                    placeholder="data"
                    disabled={isAdding}
                  />
                  <p className="text-xs text-muted-foreground mt-1">
                    A descriptive name for this volume (e.g., "data", "postgres-data")
                  </p>
                </div>
                <div>
                  <Label htmlFor="new-host" className="text-xs">Host Path</Label>
                  <Input
                    id="new-host"
                    value={newHostPath}
                    onChange={(e) => setNewHostPath(e.target.value)}
                    placeholder="mist-volume-123"
                    disabled={isAdding}
                  />
                  <p className="text-xs text-muted-foreground mt-1">
                    Docker volume name (e.g., "mist-volume-123") or absolute host path (e.g., "/var/data/myapp")
                  </p>
                </div>
                <div>
                  <Label htmlFor="new-container" className="text-xs">Container Path</Label>
                  <Input
                    id="new-container"
                    value={newContainerPath}
                    onChange={(e) => setNewContainerPath(e.target.value)}
                    placeholder={getSuggestedContainerPath() || "/data"}
                    disabled={isAdding}
                  />
                  <p className="text-xs text-muted-foreground mt-1">
                    Path inside the container where the volume will be mounted
                  </p>
                </div>
                <div className="flex items-center gap-2">
                  <Switch
                    id="new-readonly"
                    checked={newReadOnly}
                    onCheckedChange={setNewReadOnly}
                    disabled={isAdding}
                  />
                  <Label htmlFor="new-readonly" className="text-xs">
                    Read Only (container cannot modify data)
                  </Label>
                </div>
                <div className="flex gap-2">
                  <Button type="submit" size="sm" disabled={isAdding}>
                    <Plus className="h-3 w-3 mr-1" />
                    {isAdding ? "Adding..." : "Add Volume"}
                  </Button>
                  <Button
                    type="button"
                    size="sm"
                    variant="outline"
                    onClick={() => {
                      setShowAddForm(false);
                      setNewName("");
                      setNewHostPath("");
                      setNewContainerPath("");
                      setNewReadOnly(false);
                    }}
                    disabled={isAdding}
                  >
                    Cancel
                  </Button>
                </div>
              </form>
            ) : (
              <Button
                onClick={() => setShowAddForm(true)}
                size="sm"
                variant="outline"
              >
                <Plus className="h-4 w-4 mr-2" />
                Add Volume
              </Button>
            )}
          </>
        )}
      </CardContent>

      {/* Restart Dialog */}
      <Dialog open={actionDialogOpen} onOpenChange={(open) => {
        if (!isRestarting) {
          setActionDialogOpen(open);
          if (!open) {
            setRestartError(null);
          }
        }
      }}>
        <DialogContent onPointerDownOutside={(e) => isRestarting && e.preventDefault()}>
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <AlertTriangle className="h-5 w-5 text-yellow-500" />
              Restart Required
            </DialogTitle>
            <DialogDescription>
              {isRestarting 
                ? "Restarting container, please wait..." 
                : "Volume changes require restarting the container to take effect. Would you like to restart now?"}
            </DialogDescription>
          </DialogHeader>
          {restartError && (
            <div className="p-3 rounded-md bg-destructive/10 text-destructive text-sm break-words max-h-32 overflow-y-auto">
              Error: {restartError}
            </div>
          )}
          <DialogFooter className="flex gap-2 sm:justify-end">
            <Button variant="outline" onClick={handleSkipRestart} disabled={isRestarting}>
              Skip for Now
            </Button>
            <Button onClick={handleRestart} disabled={isRestarting} className="flex items-center gap-2">
              <RefreshCw className={`h-4 w-4 ${isRestarting ? 'animate-spin' : ''}`} />
              {isRestarting ? 'Restarting...' : 'Restart Now'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </Card>
  );
};
