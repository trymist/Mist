

import { useNavigate } from 'react-router-dom';


import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { MoreVertical } from 'lucide-react';

import { getProjectOwnerDisplay, formatProjectDate } from '../utils';
import type { Project } from '@/types';
import { ROUTES } from '@/constants';


export const ProjectCard: React.FC<{
  project: Project;
  onEdit: (project: Project) => void;
  onDelete: (projectId: number) => void;
  canEdit: boolean;
  canDelete: boolean;
}> = ({ project, onEdit, onDelete, canEdit, canDelete }) => {
  const navigate = useNavigate();
  const ownerDisplay = getProjectOwnerDisplay(project);

  return (
    <Card
      className="relative transition-colors hover:border-primary cursor-pointer group"
      onClick={() => navigate(`${ROUTES.PROJECTS}/${project.id}`)}
    >
      <CardHeader className="flex flex-row items-start justify-between space-y-0">
        <div className="min-w-0 flex-1">
          <CardTitle className="truncate">{project.name}</CardTitle>
          <CardDescription className="truncate">
            {project.description}
          </CardDescription>
        </div>

        {(canEdit || canDelete) && (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button
                variant="ghost"
                size="icon"
                className="h-8 w-8 shrink-0 opacity-0 group-hover:opacity-100 transition-opacity"
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
                    onEdit(project);
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
                    onDelete(project.id);
                  }}
                >
                  Delete
                </DropdownMenuItem>
              )}
            </DropdownMenuContent>
          </DropdownMenu>
        )}
      </CardHeader>

      <CardContent>
        {/* Tags */}
        {project.tags && project.tags.length > 0 && (
          <div className="flex flex-wrap gap-2 mb-4">
            {project.tags.map((tag) => (
              <Badge
                key={tag}
                variant="secondary"
                className="capitalize text-xs"
              >
                {tag}
              </Badge>
            ))}
          </div>
        )}

        <Separator className="my-4" />

        {/* Owner and date info */}
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between text-xs text-muted-foreground gap-2">
          <div className="flex items-center gap-2">
            <div className="w-6 h-6 rounded-full bg-muted flex items-center justify-center">
              <span className="text-xs font-medium text-foreground">
                {ownerDisplay.initials}
              </span>
            </div>
            <span className="truncate">{ownerDisplay.name}</span>
          </div>
          <span className="shrink-0">
            {formatProjectDate(project, 'created')}
          </span>
        </div>
      </CardContent>
    </Card>
  );
};

