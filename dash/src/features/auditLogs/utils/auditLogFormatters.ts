import { User, Webhook, Server } from 'lucide-react';
import type { AuditLog, AuditLogDetails } from '@/types';

export const getTriggerIcon = (triggerType: string) => {
  const icons = {
    user: User,
    webhook: Webhook,
    system: Server,
  };
  return icons[triggerType as keyof typeof icons] || Server;
};

export const getTriggerBadgeVariant = (triggerType: string) => {
  const variants: Record<string, 'default' | 'secondary' | 'outline'> = {
    user: 'default',
    webhook: 'secondary',
    system: 'outline',
  };
  return variants[triggerType] || 'outline';
};

export const getActionBadgeColor = (action: string) => {
  const colors: Record<string, string> = {
    create: 'bg-green-500/20 text-green-400 hover:bg-green-500/30',
    update: 'bg-blue-500/20 text-blue-400 hover:bg-blue-500/30',
    delete: 'bg-red-500/20 text-red-400 hover:bg-red-500/30',
    login: 'bg-purple-500/20 text-purple-400 hover:bg-purple-500/30',
    signup: 'bg-indigo-500/20 text-indigo-400 hover:bg-indigo-500/30',
    logout: 'bg-gray-500/20 text-gray-400 hover:bg-gray-500/30',
    start: 'bg-green-500/20 text-green-400 hover:bg-green-500/30',
    stop: 'bg-orange-500/20 text-orange-400 hover:bg-orange-500/30',
    restart: 'bg-yellow-500/20 text-yellow-400 hover:bg-yellow-500/30',
  };
  return colors[action] || 'bg-gray-500/20 text-gray-400';
};

export const parseAuditLogDetails = (detailsString?: string): AuditLogDetails | null => {
  if (!detailsString) return null;
  
  try {
    return JSON.parse(detailsString);
  } catch {
    return null;
  }
};

export const formatWebhookDetails = (details: AuditLogDetails) => {
  const items = [];
  
  if (details.repository) {
    items.push({ label: 'Repository', value: details.repository });
  }
  if (details.branch) {
    items.push({ label: 'Branch', value: details.branch });
  }
  if (details.pusher) {
    items.push({ label: 'Pushed by', value: details.pusher });
  }
  if (details.commit_hash) {
    items.push({ label: 'Commit', value: details.commit_hash.substring(0, 7) });
  }
  
  return items;
};

export const formatChangeDetails = (details: AuditLogDetails) => {
  if (!details.changes || typeof details.changes !== 'object') {
    return null;
  }
  
  const changes = details.changes as Record<string, unknown>;
  const keys = Object.keys(changes);
  
  if (keys.length === 0) return null;
  
  return keys.join(', ');
};

export const formatGenericDetails = (details: AuditLogDetails) => {
  const relevantKeys = Object.keys(details).filter(
    key => !['trigger_type', 'data'].includes(key) && details[key]
  );
  
  if (relevantKeys.length === 0) return [];
  
  return relevantKeys.slice(0, 3).map(key => ({
    key,
    value: String(details[key]).substring(0, 50),
    truncated: String(details[key]).length > 50,
  }));
};

export const getActivityDescription = (log: AuditLog): string => {
  const actor = log.triggerType === 'user' 
    ? log.username || log.email || 'Unknown User'
    : log.triggerType === 'webhook'
    ? 'GitHub Webhook'
    : 'System';
  
  return `${actor} performed ${log.action} on ${log.resourceType}`;
};

export const getResourceFilters = () => [
  { value: 'all', label: 'All Resources' },
  { value: 'user', label: 'Users' },
  { value: 'project', label: 'Projects' },
  { value: 'application', label: 'Applications' },
  { value: 'deployment', label: 'Deployments' },
  { value: 'container', label: 'Containers' },
  { value: 'env_variable', label: 'Environment Variables' },
];
