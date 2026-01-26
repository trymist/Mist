import { useState } from "react";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import type { CreateAppRequest } from "@/types/app";

interface ComposeAppFormProps {
    projectId: number;
    onSubmit: (data: CreateAppRequest) => void;
    onBack: () => void;
}

export function ComposeAppForm({ projectId, onSubmit, onBack }: ComposeAppFormProps) {
    const [formData, setFormData] = useState({
        name: "",
        description: "",
    });

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        onSubmit({
            projectId,
            appType: "compose",
            name: formData.name,
            description: formData.description || undefined,
        });
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-4">
            <div>
                <h3 className="text-lg font-semibold mb-2">Create Compose Application</h3>
                <p className="text-sm text-muted-foreground">
                    Deploy complex multi-container applications using docker-compose
                </p>
            </div>

            <div className="space-y-4">
                <div>
                    <Label htmlFor="name">Application Name *</Label>
                    <Input
                        id="name"
                        value={formData.name}
                        onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                        placeholder="my-compose-app"
                        required
                        className="mt-1"
                    />
                    <p className="text-xs text-muted-foreground mt-1">
                        Lowercase letters, numbers, and hyphens only
                    </p>
                </div>

                <div>
                    <Label htmlFor="description">Description</Label>
                    <Textarea
                        id="description"
                        value={formData.description}
                        onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                        placeholder="Brief description of your application"
                        className="mt-1"
                    />
                </div>
            </div>

            <div className="flex justify-between pt-4">
                <Button type="button" variant="outline" onClick={onBack}>
                    Back
                </Button>
                <Button type="submit">Create Application</Button>
            </div>
        </form>
    );
}
