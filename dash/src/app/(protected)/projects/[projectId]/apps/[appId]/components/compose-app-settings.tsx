import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { toast } from "sonner";
import { applicationsService } from "@/services";
import type { App } from "@/types";

interface ComposeAppSettingsProps {
    app: App;
    onUpdate: () => void;
}

export const ComposeAppSettings = ({ app, onUpdate }: ComposeAppSettingsProps) => {
    const [rootDirectory, setRootDirectory] = useState(app.rootDirectory || "");
    const [saving, setSaving] = useState(false);

    const handleSave = async () => {
        try {
            setSaving(true);

            await applicationsService.update(app.id, {
                rootDirectory,
            });

            toast.success("Settings updated successfully");
            onUpdate();
        } catch (error) {
            toast.error(error instanceof Error ? error.message : "Failed to update settings");
        } finally {
            setSaving(false);
        }
    };

    return (
        <Card>
            <CardHeader className="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-4">
                <div className="flex flex-col gap-2">
                    <CardTitle>Application Settings</CardTitle>
                    <CardDescription>
                        Configure your compose application settings.
                    </CardDescription>
                </div>
                <div className="flex justify-start sm:justify-end w-full sm:w-auto">
                    <Button onClick={handleSave} disabled={saving}>
                        {saving ? "Saving..." : "Save Settings"}
                    </Button>
                </div>
            </CardHeader>
            <CardContent className="space-y-6">
                <div className="space-y-2">
                    <Label htmlFor="rootDirectory">Root Directory</Label>
                    <Input
                        id="rootDirectory"
                        placeholder="/"
                        value={rootDirectory}
                        onChange={(e) => setRootDirectory(e.target.value)}
                    />
                    <p className="text-sm text-muted-foreground">
                        The directory containing your docker-compose.yml file
                    </p>
                </div>

                <div className="pt-4 border-t">
                    <h3 className="text-sm font-medium mb-3">Read-only Configuration</h3>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
                        <div>
                            <span className="text-muted-foreground block">Repository</span>
                            <span className="font-mono">{app.gitRepository || "Not connected"}</span>
                        </div>
                        <div>
                            <span className="text-muted-foreground block">Branch</span>
                            <span className="font-mono">{app.gitBranch}</span>
                        </div>
                    </div>
                </div>
            </CardContent>
        </Card>
    );
};
