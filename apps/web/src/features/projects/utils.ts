
import type { Project, User } from '../../types';
import { formatDate } from '../../utils/formatters/date';
import { getInitials } from '../../utils/formatters/text';

export function getProjectOwnerDisplay(project: Project): {
  name: string;
  initials: string;
  avatar?: string;
} {
  const owner = project.owner;
  
  if (!owner) {
    return {
      name: 'Unknown User',
      initials: '?',
    };
  }

  return {
    name: owner.username,
    initials: getInitials(owner.username, 2),
    avatar: undefined, // Could be extended to support avatars
  };
}

export function formatProjectDate(project: Project, type: 'created' | 'updated' = 'created'): string {
  const dateString = type === 'created' ? project.createdAt : project.updatedAt;
  return formatDate(dateString || '');
}

export function canEditProject(project: Project, user: User | null): boolean {
  if (!user) return false;
  if (user.role === 'owner') return true;
  if (user.role === 'admin') return true;
  return project.ownerId === user.id;
}

export function canDeleteProject(project: Project, user: User | null): boolean {
  if (!user) return false;
  if (user.role === 'owner') return true;
  if (user.role === 'admin') return true;
  return project.ownerId === user.id;
}

export function getProjectTagStyle(tag: string): string {
  const colors = [
    'bg-blue-100 text-blue-800',
    'bg-green-100 text-green-800',
    'bg-yellow-100 text-yellow-800',
    'bg-red-100 text-red-800',
    'bg-purple-100 text-purple-800',
    'bg-pink-100 text-pink-800',
    'bg-indigo-100 text-indigo-800',
  ];
  
  const hash = tag.split('').reduce((a, b) => {
    a = ((a << 5) - a) + b.charCodeAt(0);
    return a & a;
  }, 0);
  
  return colors[Math.abs(hash) % colors.length];
}

export function validateProjectData(data: {
  name: string;
  description: string;
  tags: string[];
}): { isValid: boolean; errors: Record<string, string> } {
  const errors: Record<string, string> = {};

  if (!data.name.trim()) {
    errors.name = 'Project name is required';
  } else if (data.name.length > 100) {
    errors.name = 'Project name must be less than 100 characters';
  }

  if (data.description.length > 500) {
    errors.description = 'Description must be less than 500 characters';
  }

  if (data.tags.length > 10) {
    errors.tags = 'Maximum 10 tags allowed';
  }

  for (const tag of data.tags) {
    if (tag.length > 20) {
      errors.tags = 'Each tag must be less than 20 characters';
      break;
    }
    if (!/^[a-zA-Z0-9_-]+$/.test(tag)) {
      errors.tags = 'Tags can only contain letters, numbers, underscores, and hyphens';
      break;
    }
  }

  return {
    isValid: Object.keys(errors).length === 0,
    errors,
  };
}

export function filterProjects(
  projects: Project[],
  searchTerm: string,
  filters: {
    tags?: string[];
    owner?: string;
    dateRange?: { start: Date; end: Date };
  } = {}
): Project[] {
  return projects.filter(project => {
    if (searchTerm) {
      const searchLower = searchTerm.toLowerCase();
      const matchesName = project.name.toLowerCase().includes(searchLower);
      const matchesDescription = project.description.toLowerCase().includes(searchLower);
      const matchesTags = project.tags?.some(tag => 
        tag.toLowerCase().includes(searchLower)
      );
      const matchesOwner = project.owner?.username.toLowerCase().includes(searchLower);
      
      if (!matchesName && !matchesDescription && !matchesTags && !matchesOwner) {
        return false;
      }
    }

    if (filters.tags && filters.tags.length > 0) {
      const hasMatchingTag = filters.tags.some(filterTag =>
        project.tags?.some(projectTag => 
          projectTag.toLowerCase() === filterTag.toLowerCase()
        )
      );
      if (!hasMatchingTag) return false;
    }

    if (filters.owner) {
      if (project.owner?.username !== filters.owner) return false;
    }

    if (filters.dateRange) {
      const projectDate = new Date(project.createdAt || '');
      if (projectDate < filters.dateRange.start || projectDate > filters.dateRange.end) {
        return false;
      }
    }

    return true;
  });
}

export function sortProjects(
  projects: Project[],
  sortBy: 'name' | 'created' | 'updated' | 'owner',
  sortOrder: 'asc' | 'desc' = 'asc'
): Project[] {
  return [...projects].sort((a, b) => {
    let comparison = 0;

    switch (sortBy) {
      case 'name':
        comparison = a.name.localeCompare(b.name);
        break;
      case 'created':
        comparison = new Date(a.createdAt || '').getTime() - new Date(b.createdAt || '').getTime();
        break;
      case 'updated':
        comparison = new Date(a.updatedAt || '').getTime() - new Date(b.updatedAt || '').getTime();
        break;
      case 'owner':
        comparison = (a.owner?.username || '').localeCompare(b.owner?.username || '');
        break;
    }

    return sortOrder === 'asc' ? comparison : -comparison;
  });
}
