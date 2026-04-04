export interface AuditLog {
  id: number;
  userId?: number;
  username?: string;
  email?: string;
  action: string;
  resourceType: string;
  resourceId?: number;
  resourceName?: string;
  details?: string;
  ipAddress?: string;
  userAgent?: string;
  triggerType: 'user' | 'webhook' | 'system';
  createdAt: string;
}

export interface AuditLogDetails {
  before?: Record<string, unknown>;
  after?: Record<string, unknown>;
  reason?: string;
  extra?: Record<string, unknown>;
  trigger_type?: string;
  repository?: string;
  branch?: string;
  pusher?: string;
  commit_hash?: string;
  commit_message?: string;
  app_id?: number;
  changes?: Record<string, unknown>;
  [key: string]: unknown;
}

export interface AuditLogsResponse {
  logs: AuditLog[];
  total: number;
  limit: number;
  offset: number;
}
