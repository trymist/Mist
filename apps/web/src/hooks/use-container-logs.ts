import { useEffect, useRef, useState, useCallback } from 'react';

export interface ContainerLogEntry {
  line: string;
  stream?: 'stdout' | 'stderr';
}

interface ContainerLogEvent {
  type: 'log' | 'status' | 'error' | 'end';
  timestamp: string;
  data: {
    line?: string;
    stream?: 'stdout' | 'stderr';
    message?: string;
    container?: string;
    state?: string;
    status?: string;
  };
}

interface UseContainerLogsOptions {
  appId: number;
  enabled: boolean;
  onError?: (error: string) => void;
}

export const useContainerLogs = ({
  appId,
  enabled,
  onError,
}: UseContainerLogsOptions) => {
  const [logs, setLogs] = useState<ContainerLogEntry[]>([]);
  const [containerState, setContainerState] = useState<string>('');
  const [error, setError] = useState<string | null>(null);
  const [isConnected, setIsConnected] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<number>(0);
  const reconnectAttemptsRef = useRef(0);
  const onErrorRef = useRef(onError);
  const intentionalCloseRef = useRef(false);

  useEffect(() => {
    onErrorRef.current = onError;
  }, [onError]);

  const connectWebSocket = useCallback(() => {
    if (!enabled) return;

    if (wsRef.current && (wsRef.current.readyState === WebSocket.OPEN || wsRef.current.readyState === WebSocket.CONNECTING)) {
      return;
    }

    try {
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const host = window.location.host;
      const ws = new WebSocket(`${protocol}//${host}/api/ws/container/logs?appId=${appId}`);
      wsRef.current = ws;

      ws.onopen = () => {
        setIsConnected(true);
        setError(null);
        setIsLoading(false);
        reconnectAttemptsRef.current = 0;
      };

      ws.onmessage = (event) => {
        if (event.data instanceof Blob) {
          return;
        }

        try {
          const logEvent: ContainerLogEvent = JSON.parse(event.data);

          switch (logEvent.type) {
            case 'log': {
              if (logEvent.data.line && logEvent.data.line.trim()) {
                setLogs((prev) => [
                  ...prev,
                  {
                    line: logEvent.data.line!,
                    stream: logEvent.data.stream,
                  },
                ]);
              }
              break;
            }

            case 'status': {
              if (logEvent.data.state) {
                setContainerState(logEvent.data.state);
              }
              break;
            }

            case 'error': {
              const errorMsg = logEvent.data.message || 'Unknown error';
              console.error('[ContainerLogs] Error:', errorMsg);
              setError(errorMsg);
              onErrorRef.current?.(errorMsg);
              break;
            }

            case 'end': {
              setError('Log stream ended');
              break;
            }
          }
        } catch (err) {
          console.error('[ContainerLogs] Error parsing message:', err);
        }
      };

      ws.onerror = (event) => {
        console.error('[ContainerLogs] WebSocket error:', event);
        setIsConnected(false);
      };

      ws.onclose = (event) => {
        console.log('[ContainerLogs] WebSocket closed - Code:', event.code);
        setIsConnected(false);
        wsRef.current = null;

        if (!intentionalCloseRef.current && enabled && reconnectAttemptsRef.current < 5) {
          const delay = Math.min(1000 * Math.pow(2, reconnectAttemptsRef.current), 10000);
          reconnectAttemptsRef.current++;

          reconnectTimeoutRef.current = window.setTimeout(() => {
            connectWebSocket();
          }, delay);
        } else if (reconnectAttemptsRef.current >= 5) {
          setError('Failed to connect after multiple attempts');
          setIsLoading(false);
        }

        // Reset the flag
        intentionalCloseRef.current = false;
      };
    } catch (err) {
      console.error('[ContainerLogs] Error creating WebSocket:', err);
      setError('Failed to establish connection');
      setIsLoading(false);
    }
  }, [appId, enabled]);

  useEffect(() => {
    if (!enabled) {
      // Close connection when disabled
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }
      if (wsRef.current) {
        intentionalCloseRef.current = true;
        wsRef.current.close();
        wsRef.current = null;
      }
      setIsConnected(false);
      setIsLoading(false);
      return;
    }

    // Only connect if not already connected or connecting
    const shouldConnect = !wsRef.current ||
      wsRef.current.readyState === WebSocket.CLOSED ||
      wsRef.current.readyState === WebSocket.CLOSING;

    if (shouldConnect) {
      connectWebSocket();
    }

    return () => {
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }
      if (wsRef.current) {
        intentionalCloseRef.current = true;
        wsRef.current.close();
        wsRef.current = null;
      }
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [enabled, appId]);

  const clearLogs = () => {
    setLogs([]);
  };

  const disconnect = () => {
    if (wsRef.current) {
      intentionalCloseRef.current = true;
      wsRef.current.close();
    }
    setIsConnected(false);
  };

  return {
    logs,
    containerState,
    error,
    isConnected,
    isLoading,
    clearLogs,
    disconnect,
  };
};
