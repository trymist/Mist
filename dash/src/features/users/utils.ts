import type { User } from '@/types';

export const getRoleStyles = (role: string): string => {
  switch (role) {
    case 'owner':
      return 'bg-purple-500/20 text-purple-400';
    case 'admin':
      return 'bg-blue-500/20 text-blue-400';
    default:
      return 'bg-muted text-muted-foreground';
  }
};

export const getUserInitials = (username: string): string => {
  return username.charAt(0).toUpperCase();
};

export const formatUserId = (id: string | number): string => {
  return `User ID: ${id}`;
};

export const canManageUsers = (user: User | null): boolean => {
  return !!user?.isAdmin;
};

export const getRoleHierarchy = (): Array<{ label: string; value: string }> => {
  return [
    { label: 'User', value: 'user' },
    { label: 'Admin', value: 'admin' },
    { label: 'Owner', value: 'owner' },
  ];
};

export const getAvailableRoles = (currentUser: User | null): Array<{ label: string; value: string }> => {
  const baseRoles = [{ label: 'User', value: 'user' }];

  if (currentUser?.isAdmin) {
    baseRoles.push({ label: 'Admin', value: 'admin' });
  }

  return baseRoles;
};
