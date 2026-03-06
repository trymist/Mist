import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { MoreVertical, GitBranch, Github, Clock, ExternalLink } from "lucide-react";
import type { App } from "@/types/app";

interface AppCardProps {
  app: App;
  onClick: () => void;
  onEdit?: (app: App) => void;
  onDelete?: (appId: number) => void;
  canEdit?: boolean;
  canDelete?: boolean;
}

const getStatusColor = (status: string) => {
  switch (status.toLowerCase()) {
    case "running":
      return "bg-green-500/10 text-green-600 border-green-500/20 dark:text-green-400";
    case "error":
    case "failed":
      return "bg-red-500/10 text-red-600 border-red-500/20 dark:text-red-400";
    case "deploying":
    case "building":
      return "bg-blue-500/10 text-blue-600 border-blue-500/20 dark:text-blue-400";
    case "stopped":
      return "bg-gray-500/10 text-gray-600 border-gray-500/20 dark:text-gray-400";
    default:
      return "bg-gray-500/10 text-gray-600 border-gray-500/20 dark:text-gray-400";
  }
};

export const AppCard: React.FC<AppCardProps> = ({
  app,
  onClick,
  onEdit,
  onDelete,
  canEdit = true,
  canDelete = true,
}) => {
  return (
    <Card
      className="relative transition-all duration-200 hover:shadow-lg hover:border-foreground/20 cursor-pointer group overflow-hidden border-border/50"
      onClick={onClick}
    >
      <div className="absolute inset-0 bg-gradient-to-br from-primary/5 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300" />

      <CardHeader className="flex flex-row items-start justify-between space-y-0 pb-4">
        <div className="min-w-0 flex-1 space-y-1">
          <div className="flex items-center gap-2">
            <CardTitle className="truncate text-lg font-semibold">{app.name}</CardTitle>
            <ExternalLink className="w-4 h-4 text-muted-foreground opacity-0 group-hover:opacity-100 transition-opacity" />
          </div>
          <CardDescription className="truncate text-sm">
            {app.description || "No description provided"}
          </CardDescription>
        </div>

        {(canEdit || canDelete) && (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button
                variant="ghost"
                size="icon"
                className="h-8 w-8 shrink-0 opacity-0 group-hover:opacity-100 transition-opacity relative z-10"
                onClick={(e) => e.stopPropagation()}
              >
                <MoreVertical className="w-4 h-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              {canEdit && (
                <DropdownMenuItem
                  onClick={(e: React.MouseEvent) => {
                    e.stopPropagation();
                    onEdit?.(app);
                  }}
                >
                  Edit
                </DropdownMenuItem>
              )}
              {canDelete && (
                <DropdownMenuItem
                  className="text-destructive"
                  onClick={(e: React.MouseEvent) => {
                    e.stopPropagation();
                    onDelete?.(app.id);
                  }}
                >
                  Delete
                </DropdownMenuItem>
              )}
            </DropdownMenuContent>
          </DropdownMenu>
        )}
      </CardHeader>

      <CardContent className="space-y-4">
        {/* Status and Strategy */}
        <div className="flex flex-wrap items-center gap-2">
          {app.status && (
            <Badge
              variant="outline"
              className={`${getStatusColor(app.status)} font-medium`}
            >
              <span className="relative flex h-2 w-2 mr-1.5">
                {app.status === "running" && (
                  <>
                    <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
                    <span className="relative inline-flex rounded-full h-2 w-2 bg-green-500"></span>
                  </>
                )}
                {app.status !== "running" && (
                  <span className="relative inline-flex rounded-full h-2 w-2 bg-current opacity-75"></span>
                )}
              </span>
              {app.status.charAt(0).toUpperCase() + app.status.slice(1)}
            </Badge>
          )}
          {app.deploymentStrategy && (
            <Badge variant="secondary" className="font-mono text-xs">
              {app.deploymentStrategy}
            </Badge>
          )}
        </div>

        {/* Git Info */}
        {(app.gitRepository || app.gitBranch) && (
          <div className="space-y-2.5 rounded-lg bg-muted/50 p-3 border border-border/50">
            {app.gitRepository && (
              <div className="flex items-center gap-2 text-sm">
                <Github className="w-3.5 h-3.5 text-muted-foreground flex-shrink-0" />
                <span className="truncate font-mono text-xs text-foreground/80">
                  {app.gitRepository}
                </span>
              </div>
            )}
            {app.gitBranch && (
              <div className="flex items-center gap-2 text-sm">
                <GitBranch className="w-3.5 h-3.5 text-muted-foreground flex-shrink-0" />
                <span className="font-mono text-xs text-foreground/80">{app.gitBranch}</span>
              </div>
            )}
          </div>
        )}

        {/* Footer Metadata */}
        <div className="flex items-center justify-between text-xs text-muted-foreground pt-2 border-t border-border/50">
          <div className="flex items-center gap-4">
            {app.port ?? (
              <div className="flex items-center gap-1.5">
                <span className="text-muted-foreground">Port</span>
                <span className="font-mono text-foreground font-medium">
                  :8080
                </span>
              </div>
            )}
          </div>
          {app.createdAt && (
            <div className="flex items-center gap-1.5">
              <Clock className="w-3.5 h-3.5" />
              <span>
                {new Date(app.createdAt).toLocaleDateString(undefined, {
                  month: "short",
                  day: "numeric",
                  year: "numeric",
                })}
              </span>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
};
