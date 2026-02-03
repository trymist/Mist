import { useState, useEffect } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Checkbox } from "@/components/ui/checkbox";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { Database, Server, Layers, HardDrive, MessageSquare, Package, type LucideIcon } from "lucide-react";
import { templatesApi } from "@/api/endpoints/templates";
import type { CreateAppRequest, ServiceTemplate } from "@/types/app";
import { toast } from "sonner";

interface DatabaseFormProps {
  projectId: number;
  onSubmit: (data: CreateAppRequest) => void;
  onBack: () => void;
}

const categoryIcons: Record<string, LucideIcon> = {
  database: Database,
  cache: Layers,
  queue: MessageSquare,
  storage: HardDrive,
  other: Package,
};

export function DatabaseForm({ projectId, onSubmit, onBack }: DatabaseFormProps) {
  const [templates, setTemplates] = useState<ServiceTemplate[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedTemplate, setSelectedTemplate] = useState<ServiceTemplate | null>(null);
  const [formData, setFormData] = useState({
    name: "",
    password: "",
    shouldExpose: false,
    exposePort: undefined as number | undefined,
  });

  useEffect(() => {
    loadTemplates();
  }, []);

  const loadTemplates = async () => {
    try {
      const response = await templatesApi.list();
      if (response.success && response.data) {
        setTemplates(response.data);
      } else {
        toast.error("Failed to load templates");
      }
    } catch (error) {
      toast.error("Failed to load templates");
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  const generatePassword = () => {
    const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*";
    let password = "";
    for (let i = 0; i < 16; i++) {
      password += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    setFormData({ ...formData, password });
  };

  const handleTemplateSelect = (template: ServiceTemplate) => {
    setSelectedTemplate(template);
    setFormData({
      name: template.name,
      password: "",
      shouldExpose: false,
      exposePort: undefined,
    });
    generatePassword();
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedTemplate) return;

    // Prepare environment variables based on template
    const envVars: Record<string, string> = {};
    
    // Add password if provided (most databases need this)
    if (formData.password) {
      // Common password env var names based on template
      if (selectedTemplate.name === 'postgres') {
        envVars['POSTGRES_PASSWORD'] = formData.password;
      } else if (selectedTemplate.name === 'mysql') {
        envVars['MYSQL_ROOT_PASSWORD'] = formData.password;
      } else if (selectedTemplate.name === 'mariadb') {
        envVars['MARIADB_ROOT_PASSWORD'] = formData.password;
      } else if (selectedTemplate.name === 'mongodb') {
        envVars['MONGO_INITDB_ROOT_PASSWORD'] = formData.password;
        envVars['MONGO_INITDB_ROOT_USERNAME'] = 'admin';
      } else if (selectedTemplate.name === 'redis') {
        envVars['REDIS_PASSWORD'] = formData.password;
      }
    }

    onSubmit({
      projectId,
      appType: "database",
      templateName: selectedTemplate.name,
      name: formData.name,
      port: selectedTemplate.defaultPort,
      shouldExpose: formData.shouldExpose,
      exposePort: formData.exposePort,
      envVars,
    });
  };

  if (loading) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-10 w-full" />
        <Skeleton className="h-32 w-full" />
        <Skeleton className="h-32 w-full" />
      </div>
    );
  }

  if (!selectedTemplate) {
    return (
      <div className="space-y-4">
        <div>
          <h3 className="text-lg font-semibold mb-2">Select a Database or Service</h3>
          <p className="text-sm text-muted-foreground">
            One-click deployment of popular databases and services
          </p>
        </div>

        <div className="grid gap-3 max-h-[400px] overflow-y-auto pr-2">
          {templates.map((template) => {
            const Icon = categoryIcons[template.category] || Server;
            return (
              <Card
                key={template.id}
                className="p-4 cursor-pointer hover:border-primary transition-colors"
                onClick={() => handleTemplateSelect(template)}
              >
                <div className="flex items-start gap-3">
                  <div className="p-2 rounded bg-primary/10">
                    <Icon className="h-5 w-5 text-primary" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 mb-1">
                      <h4 className="font-semibold">{template.displayName}</h4>
                      <Badge variant="secondary" className="text-xs">
                        {template.category}
                      </Badge>
                    </div>
                    {template.description && (
                      <p className="text-sm text-muted-foreground line-clamp-2">
                        {template.description}
                      </p>
                    )}
                    <div className="flex items-center gap-3 mt-2 text-xs text-muted-foreground">
                      <span>Port: {template.defaultPort}</span>
                      {template.recommendedMemory && (
                        <span>RAM: {template.recommendedMemory}MB</span>
                      )}
                      {template.dockerImageVersion && (
                        <span>v{template.dockerImageVersion}</span>
                      )}
                    </div>
                  </div>
                </div>
              </Card>
            );
          })}
        </div>

        <div className="flex justify-start pt-4">
          <Button type="button" variant="outline" onClick={onBack}>
            Back
          </Button>
        </div>
      </div>
    );
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div>
        <h3 className="text-lg font-semibold mb-2">
          Deploy {selectedTemplate.displayName}
        </h3>
        <p className="text-sm text-muted-foreground">
          {selectedTemplate.description || "Configure your database deployment"}
        </p>
      </div>

      <Card className="p-4 bg-muted/50">
        <div className="flex items-center gap-3 mb-3">
          {(() => {
            const Icon = categoryIcons[selectedTemplate.category] || Server;
            return <Icon className="h-5 w-5 text-primary" />;
          })()}
          <div>
            <p className="font-medium">{selectedTemplate.displayName}</p>
            <p className="text-xs text-muted-foreground">
              {selectedTemplate.dockerImage}
              {selectedTemplate.dockerImageVersion && `:${selectedTemplate.dockerImageVersion}`}
            </p>
          </div>
        </div>
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-2 text-xs">
          <div>
            <span className="text-muted-foreground">Port:</span>
            <span className="ml-2 font-medium">{selectedTemplate.defaultPort}</span>
          </div>
          {selectedTemplate.recommendedMemory && (
            <div>
              <span className="text-muted-foreground">Recommended RAM:</span>
              <span className="ml-2 font-medium">{selectedTemplate.recommendedMemory}MB</span>
            </div>
          )}
          {selectedTemplate.recommendedCpu && (
            <div>
              <span className="text-muted-foreground">Recommended CPU:</span>
              <span className="ml-2 font-medium">{selectedTemplate.recommendedCpu} cores</span>
            </div>
          )}
          {selectedTemplate.volumeRequired && (
            <div>
              <span className="text-muted-foreground">Storage:</span>
              <span className="ml-2 font-medium">Required</span>
            </div>
          )}
        </div>
      </Card>

      <div className="space-y-4">
        <div>
          <Label htmlFor="name">Instance Name *</Label>
          <Input
            id="name"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            placeholder={selectedTemplate.name}
            required
            className="mt-1"
          />
          <p className="text-xs text-muted-foreground mt-1">
            Unique name for this database instance
          </p>
        </div>

        <div>
          <Label htmlFor="password">Root/Admin Password</Label>
          <div className="flex gap-2 mt-1">
            <Input
              id="password"
              type="text"
              value={formData.password}
              onChange={(e) => setFormData({ ...formData, password: e.target.value })}
              placeholder="Auto-generated password"
              className="font-mono text-sm"
            />
            <Button type="button" variant="outline" onClick={generatePassword}>
              Generate
            </Button>
          </div>
          <p className="text-xs text-muted-foreground mt-1">
            Password will be stored as environment variable
          </p>
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
                When enabled, exposes the database on an external port accessible via IP:port.
                Use with caution - databases should typically not be exposed publicly.
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
                placeholder={selectedTemplate?.defaultPort.toString()}
                required
                min={1}
                max={65535}
                className="mt-1"
              />
              <p className="text-xs text-muted-foreground">
                The external port to expose the database on (default: {selectedTemplate?.defaultPort})
              </p>
            </div>
          )}
        </div>
      </div>

      <div className="flex justify-between pt-4">
        <Button
          type="button"
          variant="outline"
          onClick={() => setSelectedTemplate(null)}
        >
          Back
        </Button>
        <Button type="submit">Deploy {selectedTemplate.displayName}</Button>
      </div>
    </form>
  );
}
