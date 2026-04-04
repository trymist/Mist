
import type { User } from './user';

export interface Project {
  id: number;
  name: string;
  description: string;
  tags?: string[];
  ownerId: string | number;
  owner?: User;
  projectMembers: User[];
  createdAt?: string;
  updatedAt?: string;
}

export interface ProjectCreateInput {
  name: string;
  description: string;
  tags: string[];
}

export interface ProjectUpdateInput {
  name: string;
  description: string;
  tags: string[];
}

// Aliases for service layer
export type CreateProjectRequest = ProjectCreateInput;
export type UpdateProjectRequest = ProjectUpdateInput;
