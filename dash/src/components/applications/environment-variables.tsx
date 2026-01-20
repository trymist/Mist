import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Trash2, Plus, Pencil, X, Check, FileText, Info } from "lucide-react";
import { toast } from "sonner";
import { useEnvironmentVariables } from "@/hooks";
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
  const [editingId, setEditingId] = useState<number | null>(null);
  const [editKey, setEditKey] = useState("");
  const [editValue, setEditValue] = useState("");
  const [showAddForm, setShowAddForm] = useState(false);
  const [bulkMode, setBulkMode] = useState(false);
  const [bulkText, setBulkText] = useState("");
  const [parsedVars, setParsedVars] = useState<Array<{ key: string; value: string }>>([]);
  const [isAdding, setIsAdding] = useState(false);

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

  const handleAdd = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newKey.trim()) {
      toast.error("Key is required");
      return;
    }

    const result = await createEnvVar(newKey.trim(), newValue);
    if (result) {
      setNewKey("");
      setNewValue("");
      setShowAddForm(false);
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

    const result = await updateEnvVar(id, editKey.trim(), editValue);
    if (result) {
      setEditingId(null);
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm("Are you sure you want to delete this environment variable?")) {
      return;
    }
    await deleteEnvVar(id);
  };

  const startEdit = (env: EnvVariable) => {
    setEditingId(env.id);
    setEditKey(env.key);
    setEditValue(env.value);
  };

  const cancelEdit = () => {
    setEditingId(null);
    setEditKey("");
    setEditValue("");
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
                    <Input
                      id="new-value"
                      placeholder="your-api-key-value"
                      value={newValue}
                      onChange={(e) => setNewValue(e.target.value)}
                      type="password"
                    />
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
                        <Input
                          id={`edit-value-${env.id}`}
                          value={editValue}
                          onChange={(e) => setEditValue(e.target.value)}
                          type="password"
                        />
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
                    <div className="flex-1 font-mono break-all">
                      <span className="font-semibold">{env.key}</span>
                      <span className="text-muted-foreground ml-2">=</span>
                      <span className="ml-2 text-muted-foreground">••••••••</span>
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
    </Card>
  );
};
