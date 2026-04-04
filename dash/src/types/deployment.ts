// Enhanced Deployment type with new fields
export interface Deployment {
  id: number;
  app_id: number;
  commit_hash: string;
  commit_message: string;
  triggered_by?: number;
  logs?: string;
  status: string;
  stage: string;
  progress: number;
  error_message?: string;
  created_at: string;
  started_at?: string;
  finished_at?: string;
  duration?: number;
}

// Request types
export interface CreateDeploymentRequest {
  appId: number;
}

// WebSocket event types
export interface DeploymentEvent {
  type: 'log' | 'status' | 'progress' | 'error';
  timestamp: string;
  data: LogUpdate | StatusUpdate;
}

export interface LogUpdate {
  line: string;
  stream?: 'stdout' | 'stderr';
  timestamp: string;
}

export interface StatusUpdate {
  deployment_id: number;
  status: string;
  stage: string;
  progress: number;
  message: string;
  error_message?: string;
  duration?: number;
}
