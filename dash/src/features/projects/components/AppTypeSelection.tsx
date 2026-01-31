import { Globe, Cog, Database, Container } from "lucide-react";
import { Card } from "@/components/ui/card";
import type { AppType } from "@/types/app";

interface AppTypeSelectionProps {
  onSelect: (type: AppType) => void;
}

export function AppTypeSelection({ onSelect }: AppTypeSelectionProps) {
  const appTypes = [
    {
      type: "web" as AppType,
      icon: Globe,
      title: "Web Application",
      description: "HTTP servers, APIs, and web apps that need external access via domains",
      examples: "React apps, Next.js, Express APIs, Django, Flask",
    },
    {
      type: "service" as AppType,
      icon: Cog,
      title: "Background Service",
      description: "Workers, bots, and background processes that don't need external ports",
      examples: "Discord bots, Queue workers, Scheduled tasks, Cron jobs",
    },
    {
      type: "database" as AppType,
      icon: Database,
      title: "Database / Service",
      description: "Pre-configured databases and services deployed from official Docker images",
      examples: "PostgreSQL, Redis, MySQL, MongoDB, RabbitMQ",
    },
    {
      type: "compose" as AppType,
      icon: Container,
      title: "Docker Compose",
      description: "Deploy complex multi-container applications using docker-compose",
      examples: "Full stack apps, microservices",
    },
  ];

  return (
    <div className="space-y-4">
      <div>
        <h3 className="text-lg font-semibold mb-2">Select Application Type</h3>
        <p className="text-sm text-muted-foreground">
          Choose the type of application you want to deploy
        </p>
      </div>

      <div className="grid gap-4">
        {appTypes.map((appType) => {
          const Icon = appType.icon;
          return (
            <Card
              key={appType.type}
              className="p-4 cursor-pointer hover:border-primary transition-colors"
              onClick={() => onSelect(appType.type)}
            >
              <div className="flex items-start gap-4">
                <div className="p-3 rounded-lg bg-primary/10">
                  <Icon className="h-6 w-6 text-primary" />
                </div>
                <div className="flex-1">
                  <h4 className="font-semibold mb-1">{appType.title}</h4>
                  <p className="text-sm text-muted-foreground mb-2">
                    {appType.description}
                  </p>
                  <p className="text-xs text-muted-foreground">
                    <span className="font-medium">Examples:</span> {appType.examples}
                  </p>
                </div>
              </div>
            </Card>
          );
        })}
      </div>
    </div>
  );
}
