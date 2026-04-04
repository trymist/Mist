import { useEffect, useRef, useState, useCallback } from 'react';
import type { DeploymentEvent, StatusUpdate, LogUpdate, Deployment } from '@/types/deployment';

export interface DeploymentLogEntry {
  line: string;
  stream?: 'stdout' | 'stderr';
}

interface DockerLogEntry {
  stream?: string;
  aux?: Record<string, unknown>;
  error?: string;
  errorDetail?: Record<string, unknown>;
  status?: string;       // Docker pull status
  id?: string;           // Layer ID for pull
  progress?: string;     // Progress string
  progressDetail?: Record<string, unknown>; // Detailed progress
}

/**
 * Parse Docker build JSON format and extract content with stream detection
 */
function parseDockerLogLine(line: string): DeploymentLogEntry | null {
  const trimmed = line.trim();
  
  // Check if line looks like JSON
  if (!trimmed.startsWith('{')) {
    // Not JSON, detect stream type from content
    return {
      line: trimmed,
      stream: detectStreamType(trimmed),
    };
  }

  try {
    const dockerLog: DockerLogEntry = JSON.parse(trimmed);

    // Extract content based on what's available
    if (dockerLog.error) {
      // Error field present - this is an error message
      return {
        line: dockerLog.error,
        stream: 'stderr',
      };
    }

    if (dockerLog.stream) {
      // Stream field present - this is the main content
      const content = dockerLog.stream;
      return {
        line: content,
        stream: detectStreamType(content),
      };
    }

    // Handle Docker image pull progress logs (status field)
    // Skip most progress updates to reduce noise
    if (dockerLog.status) {
      switch (dockerLog.status) {
        case 'Downloading':
        case 'Extracting':
        case 'Waiting':
        case 'Verifying Checksum':
          // Skip noisy progress updates
          return null;
        case 'Pull complete':
        case 'Download complete':
        case 'Already exists':
          // Show completion messages
          const completeMsg = dockerLog.id 
            ? `${dockerLog.status}: ${dockerLog.id}` 
            : dockerLog.status;
          return {
            line: completeMsg,
            stream: 'stdout',
          };
        default:
          // Show other status messages (like "Pulling from...")
          const statusMsg = dockerLog.id 
            ? `${dockerLog.status}: ${dockerLog.id}` 
            : dockerLog.status;
          return {
            line: statusMsg,
            stream: 'stdout',
          };
      }
    }

    // Aux field only (metadata like image IDs) - skip these
    if (dockerLog.aux) {
      return null;
    }
  } catch {
    // Failed to parse as JSON, treat as regular line
    return {
      line: trimmed,
      stream: detectStreamType(trimmed),
    };
  }

  // Unknown format, return original
  return {
    line: trimmed,
    stream: 'stdout',
  };
}

/**
 * Detect if a log line is from stderr based on common patterns
 */
function detectStreamType(line: string): 'stdout' | 'stderr' {
  const lineLower = line.toLowerCase();

  const stderrPatterns = [
    'error:',
    'err:',
    'fatal:',
    'panic:',
    'warning:',
    'warn:',
    'failed',
    'failure',
    'exception:',
    'traceback',
    'stack trace',
    ' err ',
    '[error]',
    '[err]',
    '[fatal]',
    '[panic]',
    '[warning]',
    '[warn]',
  ];

  for (const pattern of stderrPatterns) {
    if (lineLower.includes(pattern)) {
      return 'stderr';
    }
  }

  return 'stdout';
}

interface UseDeploymentMonitorOptions {
  deploymentId: number;
  enabled: boolean;
  onComplete?: () => void;
  onError?: (error: string) => void;
  onClose?: () => void;
}

export const useDeploymentMonitor = ({
  deploymentId,
  enabled,
  onComplete,
  onError,
  onClose,
}: UseDeploymentMonitorOptions) => {
  const [logs, setLogs] = useState<DeploymentLogEntry[]>([]);
  const [status, setStatus] = useState<StatusUpdate | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isConnected, setIsConnected] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [isLive, setIsLive] = useState(false);

  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<number>(0);
  const reconnectAttemptsRef = useRef(0);
  const hasFetchedRef = useRef(false);
  const hasCompletedRef = useRef(false);

  const fetchCompletedDeployment = useCallback(async () => {
    if (hasFetchedRef.current) return;

    try {
      setIsLoading(true);
      const response = await fetch(`/api/deployments/logs?id=${deploymentId}`, {
        credentials: 'include',
      });

      if (!response.ok) {
        if (response.status === 400) {
          setIsLive(true);
          setIsLoading(false);
          hasFetchedRef.current = true;
          return;
        }
        if (response.status === 404) {
          const result = await response.json();
          const errorMsg = result.message || 'Deployment log file not found';
          setError(errorMsg);
          onError?.(errorMsg);
          onClose?.();
          setIsLoading(false);
          hasFetchedRef.current = true;
          return;
        }
        throw new Error('Failed to fetch deployment logs');
      }

      const result = await response.json();
      const deployment: Deployment = result.data.deployment;
      const logsContent: string = result.data.logs;

      if (logsContent) {
        setLogs(
          logsContent
            .split('\n')
            .filter(line => line.length > 0)
            .map(line => parseDockerLogLine(line))
            .filter(entry => entry !== null) as DeploymentLogEntry[]
        );
      }

      setStatus({
        deployment_id: deployment.id,
        status: deployment.status,
        stage: deployment.stage,
        progress: deployment.progress,
        message: deployment.status === 'success'
          ? 'Deployment completed successfully'
          : deployment.status === 'stopped'
          ? 'Deployment was stopped by user'
          : 'Deployment failed',
        error_message: deployment.error_message,
        duration: deployment.duration,
      });

      if (deployment.status === 'failed' && deployment.error_message) {
        setError(deployment.error_message);
      } else if (deployment.status === 'stopped') {
        setError('Deployment was stopped by user');
      }

      setIsLoading(false);
      hasFetchedRef.current = true;
    } catch (err) {
      console.error('[DeploymentMonitor] Error fetching completed deployment:', err);
      setError('Failed to load deployment logs');
      setIsLoading(false);
      hasFetchedRef.current = true;
    }
  }, [deploymentId]);

  const connectWebSocket = useCallback(() => {
    if (!enabled || !isLive) return;

    // Prevent duplicate connections
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      console.log('[DeploymentMonitor] WebSocket already connected, skipping');
      return;
    }

    try {
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const host = window.location.host;
      const ws = new WebSocket(`${protocol}//${host}/api/deployments/logs/stream?id=${deploymentId}`)
      wsRef.current = ws;

      ws.onopen = () => {
        console.log('[DeploymentMonitor] WebSocket connected');
        setIsConnected(true);
        setError(null);
        setIsLoading(false);
        reconnectAttemptsRef.current = 0;
      };

      ws.onmessage = (event) => {
        // Handle ping messages
        if (event.data instanceof Blob) {
          return;
        }

        try {
          const deploymentEvent: DeploymentEvent = JSON.parse(event.data);

          switch (deploymentEvent.type) {
            case 'log': {
              const logData = deploymentEvent.data as LogUpdate;
              if (logData.line && logData.line.trim()) {
                setLogs((prev) => [
                  ...prev,
                  {
                    line: logData.line,
                    stream: logData.stream,
                  },
                ]);
              }
              break;
            }

            case 'status': {
              const statusData = deploymentEvent.data as StatusUpdate;
              setStatus(statusData);

              // Handle completion - only call once
              if (statusData.status === 'success' && !hasCompletedRef.current) {
                console.log('[DeploymentMonitor] Deployment completed successfully');
                hasCompletedRef.current = true;
                onComplete?.();
              }

              // Handle errors
              if (statusData.status === 'failed' && statusData.error_message) {
                console.error('[DeploymentMonitor] Deployment failed:', statusData.error_message);
                setError(statusData.error_message);
                onError?.(statusData.error_message);
              }

              // Handle stopped status
              if (statusData.status === 'stopped' && !hasCompletedRef.current) {
                console.log('[DeploymentMonitor] Deployment stopped by user');
                hasCompletedRef.current = true;
                setError('Deployment was stopped by user');
              }
              break;
            }

            case 'error': {
              const errorMsg = (deploymentEvent.data as { message?: string }).message || 'Unknown error';
              console.error('[DeploymentMonitor] Error event:', errorMsg);
              setError(errorMsg);
              onError?.(errorMsg);

              // If it's a log file not found error, close the viewer
              if (errorMsg.includes('log file not found') || errorMsg.includes('Failed to read deployment logs')) {
                onClose?.();
              }
              break;
            }
          }
        } catch (err) {
          console.error('[DeploymentMonitor] Error parsing message:', err);
        }
      };

      ws.onerror = (event) => {
        console.error('[DeploymentMonitor] WebSocket error:', event);
        setIsConnected(false);
      };

      ws.onclose = (event) => {
        console.log('[DeploymentMonitor] WebSocket closed - Code:', event.code, 'Reason:', event.reason);
        setIsConnected(false);
        wsRef.current = null;

        // Check the current status from state to decide on reconnection
        setStatus((currentStatus) => {
          // Don't reconnect if deployment is complete (success, failed, or stopped)
          if (currentStatus?.status === 'success' || currentStatus?.status === 'failed' || currentStatus?.status === 'stopped') {
            console.log('[DeploymentMonitor] Deployment complete, not reconnecting');
            return currentStatus;
          }

          // Reconnect with exponential backoff only if still enabled
          if (enabled && reconnectAttemptsRef.current < 10) {
            const delay = Math.min(1000 * Math.pow(2, reconnectAttemptsRef.current), 30000);
            reconnectAttemptsRef.current++;

            console.log(`[DeploymentMonitor] Reconnecting in ${delay}ms (attempt ${reconnectAttemptsRef.current})`);
            reconnectTimeoutRef.current = window.setTimeout(() => {
              connectWebSocket();
            }, delay);
          } else if (reconnectAttemptsRef.current >= 10) {
            setError('Failed to connect after multiple attempts');
            setIsLoading(false);
          }

          return currentStatus;
        });
      };
    } catch (err) {
      console.error('[DeploymentMonitor] Error creating WebSocket:', err);
      setError('Failed to establish connection');
      setIsLoading(false);
    }
  }, [deploymentId, enabled, isLive, onComplete, onError, onClose]);

  useEffect(() => {
    if (enabled) {
      fetchCompletedDeployment();
    }

    return () => {
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, [enabled, fetchCompletedDeployment]);

  // Separate effect for WebSocket management that only runs when isLive changes
  useEffect(() => {
    if (!isLive || !enabled) {
      return;
    }

    // Only connect if we don't have an active connection
    if (!wsRef.current || wsRef.current.readyState === WebSocket.CLOSED) {
      connectWebSocket();
    }

    return () => {
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }
      if (wsRef.current) {
        wsRef.current.close();
        wsRef.current = null;
      }
    };
  }, [isLive, enabled]);

  const reset = () => {
    setLogs([]);
    setStatus(null);
    setError(null);
    setIsConnected(false);
    setIsLoading(true);
    setIsLive(false);
    hasFetchedRef.current = false;
    hasCompletedRef.current = false;
  };

  return {
    logs,
    status,
    error,
    isConnected,
    isLoading,
    isLive,
    reset,
  };
};
