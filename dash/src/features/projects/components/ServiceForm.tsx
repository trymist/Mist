import { useState } from "react";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Checkbox } from "@/components/ui/checkbox";
import type { CreateAppRequest } from "@/types/app";

interface ServiceFormProps {
  projectId: number;
  onSubmit: (data: CreateAppRequest) => void;
  onBack: () => void;
}

export function ServiceForm({ projectId, onSubmit, onBack }: ServiceFormProps) {
  const [formData, setFormData] = useState({
    name: "",
    description: "",
    shouldExpose: false,
    exposePort: undefined as number | undefined,
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit({
      projectId,
      appType: "service",
      name: formData.name,
      description: formData.description || undefined,
      port: 3000, // Internal port
      shouldExpose: formData.shouldExpose,
      exposePort: formData.exposePort,
    });
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div>
        <h3 className="text-lg font-semibold mb-2">Create Background Service</h3>
        <p className="text-sm text-muted-foreground">
          Workers, bots, and processes that run internally without external access
        </p>
      </div>

      <div className="space-y-4">
        <div>
          <Label htmlFor="name">Service Name *</Label>
          <Input
            id="name"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            placeholder="my-discord-bot"
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
            placeholder="Brief description of your service"
            className="mt-1"
          />
        </div>

        <div className="p-4 border rounded-lg bg-muted/30 space-y-4">
          <div className="flex items-start space-x-3">
            <Checkbox
              id="shouldExpose"
              checked={formData.shouldExpose}
              onCheckedChange={(checked) => 
                setFormData({ ...formData, shouldExpose: checked as boolean })
              }
            />
            <div className="space-y-1">
              <Label htmlFor="shouldExpose" className="font-medium cursor-pointer">
                Expose External Port
              </Label>
              <p className="text-xs text-muted-foreground">
                When enabled, exposes the service on an external port accessible via IP:port.
                Useful for services that need external access.
              </p>
            </div>
          </div>

          {formData.shouldExpose && (
            <div className="pl-7 space-y-2">
              <Label htmlFor="exposePort">External Port *</Label>
              <Input
                id="exposePort"
                type="number"
                value={formData.exposePort || ""}
                onChange={(e) => {
                  const value = e.target.value ? parseInt(e.target.value) : undefined;
                  setFormData({ ...formData, exposePort: value });
                }}
                placeholder="3000"
                required
                min={1}
                max={65535}
                className="mt-1"
              />
              <p className="text-xs text-muted-foreground">
                The external port to expose the service on
              </p>
            </div>
          )}
        </div>
      </div>

      <div className="flex justify-between pt-4">
        <Button type="button" variant="outline" onClick={onBack}>
          Back
        </Button>
        <Button type="submit">Create Service</Button>
      </div>
    </form>
  );
}
