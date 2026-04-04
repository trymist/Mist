import type { User } from '@/types';

export const generateGitHubState = (appId: number, userId: number): string => {
  const payload = { appId, userId };
  return btoa(JSON.stringify(payload));
};

export const parseGitHubState = (state: string): { appId: number; userId: number } | null => {
  try {
    return JSON.parse(atob(state));
  } catch {
    return null;
  }
};

export const getGitHubInstallUrl = (appSlug: string, appId: number, userId: number): string => {
  const state = generateGitHubState(appId, userId);
  return `https://github.com/apps/${appSlug}/installations/new?state=${state}`;
};

export const getGitHubManageUrl = (appSlug: string): string => {
  return `https://github.com/apps/${appSlug}/installations/select_target`;
};

export const canManageGitApps = (user: User | null): boolean => {
  return !!user?.isAdmin;
};

export const formatRepositoryName = (fullName: string): { owner: string; name: string } => {
  const [owner, name] = fullName.split('/');
  return { owner, name };
};

export const getProviderIcon = (provider: string): string => {
  switch (provider.toLowerCase()) {
    case 'github':
      return 'Github';
    case 'gitlab':
      return 'Gitlab';
    case 'bitbucket':
      return 'GitFork';
    case 'gitea':
      return 'GitMerge';
    default:
      return 'Git';
  }
};

export const isProviderSupported = (provider: string): boolean => {
  return ['github'].includes(provider.toLowerCase());
};

export const getRepositoryUrl = (fullName: string, provider: string = 'github'): string => {
  switch (provider.toLowerCase()) {
    case 'github':
      return `https://github.com/${fullName}`;
    case 'gitlab':
      return `https://gitlab.com/${fullName}`;
    default:
      return '#';
  }
};
