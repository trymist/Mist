import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Checkbox } from "@/components/ui/checkbox";
import { Trash2, Plus, Pencil, X, Check, FileText, Info, Eye, EyeOff, AlertTriangle, RefreshCw } from "lucide-react";
import { toast } from "sonner";
import { useEnvironmentVariables } from "@/hooks";
import { applicationsService } from "@/services";
import type { EnvVariable } from "@/types";
import { Alert, AlertDescription } from "@/components/ui/alert";

interface EnvironmentVariablesProps {
  appId: number;
}

export const EnvironmentVariables = ({ appId }: EnvironmentVariablesProps) => {
  const { envVars, loading, createEnvVar, updateEnvVar, deleteEnvVar } = useEnvironmentVariables({
    appId,
    autoFetch: true
  });

  const [newKey, setNewKey] = useState("");
  const [newValue, setNewValue] = useState("");
  const [newRuntime, setNewRuntime] = useState(true);
  const [newBuildtime, setNewBuildtime] = useState(false);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [editKey, setEditKey] = useState("");
  const [editValue, setEditValue] = useState("");
  const [editRuntime, setEditRuntime] = useState(true);
  const [editBuildtime, setEditBuildtime] = useState(false);
  const [showAddForm, setShowAddForm] = useState(false);
  const [bulkMode, setBulkMode] = useState(false);
  const [bulkText, setBulkText] = useState("");
  const [parsedVars, setParsedVars] = useState<Array<{ key: string; value: string }>>([]);
  const [isAdding, setIsAdding] = useState(false);
  const [showNewValue, setShowNewValue] = useState(false);
  const [showEditValue, setShowEditValue] = useState(false);
  const [visibleVarIds, setVisibleVarIds] = useState<Set<number>>(new Set());

  // Dialog state for redeployment
  const [actionDialogOpen, setActionDialogOpen] = useState(false);

  const parseEnvText = (text: string) => {
    const lines = text.split('\n').filter(line => line.trim());
    const parsed: Array<{ key: string; value: string }> = [];
    const errors: string[] = [];

    lines.forEach((line, index) => {
      const trimmed = line.trim();

      if (!trimmed || trimmed.startsWith('#')) {
        return;
      }

      const match = trimmed.match(/^([A-Za-z_][A-Za-z0-9_]*)\s*=\s*(.*)$/);

      if (match) {
        let [, key, value] = match;

        if ((value.startsWith('"') && value.endsWith('"')) ||
          (value.startsWith("'") && value.endsWith("'"))) {
          value = value.slice(1, -1);
        }

        parsed.push({ key: key.trim(), value: value.trim() });
      } else {
        errors.push(`Line ${index + 1}: Invalid format "${trimmed}"`);
      }
    });

    if (errors.length > 0) {
      toast.error(`Parsing errors:\n${errors.join('\n')}`);
    }

    return parsed;
  };

  const handleBulkTextChange = (text: string) => {
    setBulkText(text);
    if (text.trim()) {
      const parsed = parseEnvText(text);
      setParsedVars(parsed);
    } else {
      setParsedVars([]);
    }
  };

  const showRedeployDialog = () => {
    setActionDialogOpen(true);
  };

  const handleRedeploy = async () => {
    try {
      setActionDialogOpen(false);
      toast.info("Starting full redeployment (this may take several minutes)...");
      await applicationsService.redeploy(appId);
      toast.success("Redeployment started - environment variables will be applied");
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to trigger redeployment");
    }
  };

  const handleSkipRedeploy = () => {
    setActionDialogOpen(false);
    toast.info("Environment variable changes will take effect on next redeployment");
  };

  const handleAdd = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newKey.trim()) {
      toast.error("Key is required");
      return;
    }

    // Validate that at least one checkbox is selected
    if (!newRuntime && !newBuildtime) {
      toast.error("At least one of 'Runtime' or 'Buildtime' must be selected");
      return;
    }

    const result = await createEnvVar(newKey.trim(), newValue, newRuntime, newBuildtime);
    if (result) {
      setNewKey("");
      setNewValue("");
      setNewRuntime(true);
      setNewBuildtime(false);
      setShowNewValue(false);
      setShowAddForm(false);
      // Show redeploy dialog since env vars require full redeployment
      showRedeployDialog();
    }
  };

  const handleBulkAdd = async (e: React.FormEvent) => {
    e.preventDefault();

    if (parsedVars.length === 0) {
      toast.error("No valid environment variables to add");
      return;
    }

    setIsAdding(true);
    let successCount = 0;
    let failCount = 0;

    for (const { key, value } of parsedVars) {
      const result = await createEnvVar(key, value);
      if (result) {
        successCount++;
      } else {
        failCount++;
      }
    }

    setIsAdding(false);

    if (successCount > 0) {
      toast.success(`Added ${successCount} environment variable${successCount > 1 ? 's' : ''}`);
      // Show redeploy dialog since env vars require full redeployment
      showRedeployDialog();
    }
    if (failCount > 0) {
      toast.error(`Failed to add ${failCount} variable${failCount > 1 ? 's' : ''}`);
    }

    setBulkText("");
    setParsedVars([]);
    setShowAddForm(false);
    setBulkMode(false);
  };

  const handleUpdate = async (id: number) => {
    if (!editKey.trim()) {
      toast.error("Key is required");
      return;
    }

    // Validate that at least one checkbox is selected
    if (!editRuntime && !editBuildtime) {
      toast.error("At least one of 'Runtime' or 'Buildtime' must be selected");
      return;
    }

    const result = await updateEnvVar(id, editKey.trim(), editValue, editRuntime, editBuildtime);
    if (result) {
      setEditingId(null);
      // Show redeploy dialog since env vars require full redeployment
      showRedeployDialog();
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm("Are you sure you want to delete this environment variable?")) {
      return;
    }
    const result = await deleteEnvVar(id);
    if (result) {
      // Show redeploy dialog since env vars require full redeployment
      showRedeployDialog();
    }
  };

  const startEdit = (env: EnvVariable) => {
    setEditingId(env.id);
    setEditKey(env.key);
    setEditValue(env.value);
    // Default to runtime=true if not specified (backward compatibility)
    setEditRuntime(env.runtime !== false);
    setEditBuildtime(env.buildtime === true);
    setShowEditValue(false);
  };

  const cancelEdit = () => {
    setEditingId(null);
    setEditKey("");
    setEditValue("");
    setEditRuntime(true);
    setEditBuildtime(false);
  };

  const toggleVisibility = (id: number) => {
    setVisibleVarIds(prev => {
      const next = new Set(prev);
      if (next.has(id)) {
        next.delete(id);
      } else {
        next.add(id);
      }
      return next;
    });
  };

  if (loading) {
    return <div className="text-muted-foreground">Loading environment variables...</div>;
  }

  return (
    <Card>
      <CardHeader>
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <CardTitle>Environment Variables</CardTitle>
            <CardDescription>
              Manage environment variables for your application. Changes will be applied on next deployment.
            </CardDescription>
          </div>
          {!showAddForm && (
            <Button onClick={() => setShowAddForm(true)} size="sm" className="w-full sm:w-auto">
              <Plus className="h-4 w-4 mr-2" />
              Add Variable
            </Button>
          )}
        </div>
      </CardHeader>
      <CardContent className="space-y-4">
        <Alert>
          <Info className="h-4 w-4" />
          <AlertDescription>
            Environment variable changes require a full redeploy. Please redeploy your application after making changes for them to take effect.
          </AlertDescription>
        </Alert>

        {showAddForm && (
          <div className="space-y-4 p-4 border rounded-lg bg-muted/50">
            <div className="flex items-center gap-2 pb-2 border-b">
              <Button
                type="button"
                variant={!bulkMode ? "default" : "ghost"}
                size="sm"
                onClick={() => {
                  setBulkMode(false);
                  setBulkText("");
                  setParsedVars([]);
                }}
              >
                <Plus className="h-4 w-4 mr-2" />
                Single
              </Button>
              <Button
                type="button"
                variant={bulkMode ? "default" : "ghost"}
                size="sm"
                onClick={() => {
                  setBulkMode(true);
                  setNewKey("");
                  setNewValue("");
                }}
              >
                <FileText className="h-4 w-4 mr-2" />
                Bulk Paste
              </Button>
            </div>

            {!bulkMode ? (
              <form onSubmit={handleAdd} className="space-y-4">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="new-key">Key</Label>
                    <Input
                      id="new-key"
                      placeholder="API_KEY"
                      value={newKey}
                      onChange={(e) => setNewKey(e.target.value)}
                      autoFocus
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="new-value">Value</Label>
                    <div className="relative">
                      <Input
                        id="new-value"
                        placeholder="your-api-key-value"
                        value={newValue}
                        onChange={(e) => setNewValue(e.target.value)}
                        type={showNewValue ? "text" : "password"}
                        className="pr-10"
                      />
                      <Button
                        type="button"
                        variant="ghost"
                        size="sm"
                        className="absolute right-0 top-0 h-full px-3 hover:bg-transparent"
                        onClick={() => setShowNewValue(!showNewValue)}
                      >
                        {showNewValue ? (
                          <EyeOff className="h-4 w-4 text-muted-foreground" />
                        ) : (
                          <Eye className="h-4 w-4 text-muted-foreground" />
                        )}
                      </Button>
                    </div>
                  </div>
                </div>
                <div className="flex items-center gap-6 py-2">
                  <div className="flex items-center space-x-2">
                    <Checkbox
                      id="new-runtime"
                      checked={newRuntime}
                      onCheckedChange={(checked) => setNewRuntime(checked as boolean)}
                    />
                    <Label htmlFor="new-runtime" className="text-sm font-normal cursor-pointer">
                      Runtime
                    </Label>
                  </div>
                  <div className="flex items-center space-x-2">
                    <Checkbox
                      id="new-buildtime"
                      checked={newBuildtime}
                      onCheckedChange={(checked) => setNewBuildtime(checked as boolean)}
                    />
                    <Label htmlFor="new-buildtime" className="text-sm font-normal cursor-pointer">
                      Buildtime
                    </Label>
                  </div>
                </div>
                <div className="flex flex-col sm:flex-row gap-2">
                  <Button type="submit" size="sm" className="w-full sm:w-auto">
                    <Check className="h-4 w-4 mr-2" />
                    Add
                  </Button>
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    className="w-full sm:w-auto"
                    onClick={() => {
                      setShowAddForm(false);
                      setNewKey("");
                      setNewValue("");
                    }}
                  >
                    <X className="h-4 w-4 mr-2" />
                    Cancel
                  </Button>
                </div>
              </form>
            ) : (
              <form onSubmit={handleBulkAdd} className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="bulk-text">
                    Paste Environment Variables
                    <span className="text-xs text-muted-foreground ml-2">
                      (Format: KEY=VALUE, one per line)
                    </span>
                  </Label>
                  <Textarea
                    id="bulk-text"
                    placeholder="API_KEY=your-api-key&#10;DATABASE_URL=postgres://...&#10;PORT=3000"
                    value={bulkText}
                    onChange={(e) => handleBulkTextChange(e.target.value)}
                    rows={8}
                    className="font-mono text-sm"
                    autoFocus
                  />
                </div>

                {parsedVars.length > 0 && (
                  <div className="space-y-2">
                    <Label className="text-sm">
                      Preview ({parsedVars.length} variable{parsedVars.length > 1 ? 's' : ''} detected)
                    </Label>
                    <div className="max-h-40 overflow-y-auto space-y-1 p-2 border rounded bg-background">
                      {parsedVars.map((v, idx) => (
                        <div key={idx} className="font-mono text-xs text-muted-foreground">
                          <span className="font-semibold text-foreground">{v.key}</span>
                          <span className="mx-1">=</span>
                          <span>••••••••</span>
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                <div className="flex flex-col sm:flex-row gap-2">
                  <Button type="submit" size="sm" disabled={parsedVars.length === 0 || isAdding} className="w-full sm:w-auto">
                    <Check className="h-4 w-4 mr-2" />
                    {isAdding ? "Adding..." : `Add ${parsedVars.length} Variable${parsedVars.length > 1 ? 's' : ''}`}
                  </Button>
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    className="w-full sm:w-auto"
                    onClick={() => {
                      setShowAddForm(false);
                      setBulkMode(false);
                      setBulkText("");
                      setParsedVars([]);
                    }}
                    disabled={isAdding}
                  >
                    <X className="h-4 w-4 mr-2" />
                    Cancel
                  </Button>
                </div>
              </form>
            )}
          </div>
        )}

        {envVars.length === 0 && !showAddForm ? (
          <p className="text-muted-foreground text-center py-8">
            No environment variables added yet. Click "Add Variable" to get started.
          </p>
        ) : (
          <div className="space-y-2">
            {envVars.map((env) => env ? (
              <div key={env.id} className="p-4 border rounded-lg bg-card">
                {editingId === env.id ? (
                  <div className="space-y-4">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div className="space-y-2">
                        <Label htmlFor={`edit-key-${env.id}`}>Key</Label>
                        <Input
                          id={`edit-key-${env.id}`}
                          value={editKey}
                          onChange={(e) => setEditKey(e.target.value)}
                        />
                      </div>
                      <div className="space-y-2">
                        <Label htmlFor={`edit-value-${env.id}`}>Value</Label>
                        <div className="relative">
                          <Input
                            id={`edit-value-${env.id}`}
                            value={editValue}
                            onChange={(e) => setEditValue(e.target.value)}
                            type={showEditValue ? "text" : "password"}
                            className="pr-10"
                          />
                          <Button
                            type="button"
                            variant="ghost"
                            size="sm"
                            className="absolute right-0 top-0 h-full px-3 hover:bg-transparent"
                            onClick={() => setShowEditValue(!showEditValue)}
                          >
                            {showEditValue ? (
                              <EyeOff className="h-4 w-4 text-muted-foreground" />
                            ) : (
                              <Eye className="h-4 w-4 text-muted-foreground" />
                            )}
                          </Button>
                        </div>
                      </div>
                    </div>
                    <div className="flex items-center gap-6 py-2">
                      <div className="flex items-center space-x-2">
                        <Checkbox
                          id={`edit-runtime-${env.id}`}
                          checked={editRuntime}
                          onCheckedChange={(checked) => setEditRuntime(checked as boolean)}
                        />
                        <Label htmlFor={`edit-runtime-${env.id}`} className="text-sm font-normal cursor-pointer">
                          Runtime
                        </Label>
                      </div>
                      <div className="flex items-center space-x-2">
                        <Checkbox
                          id={`edit-buildtime-${env.id}`}
                          checked={editBuildtime}
                          onCheckedChange={(checked) => setEditBuildtime(checked as boolean)}
                        />
                        <Label htmlFor={`edit-buildtime-${env.id}`} className="text-sm font-normal cursor-pointer">
                          Buildtime
                        </Label>
                      </div>
                    </div>
                    <div className="flex flex-col sm:flex-row gap-2">
                      <Button size="sm" onClick={() => handleUpdate(env.id)} className="w-full sm:w-auto">
                        <Check className="h-4 w-4 mr-2" />
                        Save
                      </Button>
                      <Button size="sm" variant="outline" onClick={cancelEdit} className="w-full sm:w-auto">
                        <X className="h-4 w-4 mr-2" />
                        Cancel
                      </Button>
                    </div>
                  </div>
                ) : (
                  <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
                    <div className="flex-1 font-mono break-all flex items-center gap-2 flex-wrap">
                      <span className="font-semibold">{env.key}</span>
                      <span className="text-muted-foreground">=</span>
                      <span className="text-muted-foreground">
                        {visibleVarIds.has(env.id) ? env.value : "••••••••"}
                      </span>
                      <Button
                        variant="ghost"
                        size="icon"
                        className="h-6 w-6"
                        onClick={() => toggleVisibility(env.id)}
                      >
                        {visibleVarIds.has(env.id) ? (
                          <EyeOff className="h-3 w-3" />
                        ) : (
                          <Eye className="h-3 w-3" />
                        )}
                      </Button>
                      <div className="flex items-center gap-1">
                        {(env.runtime !== false) && (
                          <span className="inline-flex items-center px-1.5 py-0.5 rounded text-xs font-medium bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200">
                            R
                          </span>
                        )}
                        {env.buildtime && (
                          <span className="inline-flex items-center px-1.5 py-0.5 rounded text-xs font-medium bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200">
                            B
                          </span>
                        )}
                      </div>
                    </div>
                    <div className="flex gap-2 self-end sm:self-auto">
                      <Button
                        size="sm"
                        variant="ghost"
                        onClick={() => startEdit(env)}
                      >
                        <Pencil className="h-4 w-4" />
                      </Button>
                      <Button
                        size="sm"
                        variant="ghost"
                        onClick={() => handleDelete(env.id)}
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </div>
                  </div>
                )}
              </div>
            ) : null)}
          </div>
        )}
      </CardContent>

      {/* Redeploy Dialog */}
      <Dialog open={actionDialogOpen} onOpenChange={setActionDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <AlertTriangle className="h-5 w-5 text-yellow-500" />
              Redeployment Required
            </DialogTitle>
            <DialogDescription>
              Environment variable changes require a full redeployment to take effect.
              Would you like to redeploy now?
            </DialogDescription>
          </DialogHeader>
          <DialogFooter className="flex gap-2 sm:justify-end">
            <Button variant="outline" onClick={handleSkipRedeploy}>
              Skip for Now
            </Button>
            <Button onClick={handleRedeploy} className="flex items-center gap-2">
              <RefreshCw className="h-4 w-4" />
              Redeploy Now
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </Card>
  );
};
