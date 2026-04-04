import { useEffect, useState, useCallback } from 'react';
import { toast } from 'sonner';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { FullScreenLoading } from '@/components/common';
import { SystemOverview, ChartCard, MetricCard } from './components';
import { formatPercentage, formatMemory } from './utils';

export interface DiskInfo {
  name: string;
  totalSpace: number;
  availableSpace: number;
  usedSpace: number;
}

export interface SystemStats {
  cpuUsage: number;
  memory: {
    total: number;
    used: number;
  };
  disk: DiskInfo[];
  loadAverage: {
    oneMinute: number;
    fiveMinutes: number;
    fifteenMinutes: number;
  };
  timestamp: number;
  uptime: number;
  cpuTemperature: number;
}

export default function DashboardPage() {
  const [stats, setStats] = useState<SystemStats[]>([]);
  const [isConnected, setIsConnected] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [wsConnection, setWsConnection] = useState<WebSocket | null>(null);

  const connectWebSocket = useCallback(() => {
    try {
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const wsUrl = `${protocol}//${window.location.host}/api/ws/stats`;
      const ws = new WebSocket(wsUrl);

      ws.onopen = () => {
        setIsConnected(true);
        setIsLoading(false);
        setError(null);
      };

      ws.onmessage = (event) => {
        try {
          const newStats: SystemStats = JSON.parse(event.data);
          setStats(prevStats => {
            const updatedStats = [...prevStats, newStats];
            return updatedStats.slice(-50);
          });
        } catch (err) {
          console.error('Error parsing WebSocket message:', err);
        }
      };

      ws.onerror = (wsError) => {
        console.error('WebSocket error:', wsError);
        setIsLoading(false);
      };

      ws.onclose = (event) => {
        console.log('[Dashboard] WebSocket closed - Code:', event.code);
        setIsConnected(false);

        setTimeout(() => {
          if (!wsConnection) connectWebSocket();
        }, 5000);
      };

      setWsConnection(ws);
    } catch (error) {
      console.error('Failed to establish WebSocket connection:', error);
      setError('Failed to establish WebSocket connection');
      setIsLoading(false);
    }
  }, [wsConnection]);

  const disconnectWebSocket = useCallback(() => {
    if (wsConnection) {
      wsConnection.close();
      setWsConnection(null);
      setIsConnected(false);
    }
  }, [wsConnection]);

  const getLatestStats = (): SystemStats | null => {
    return stats.length > 0 ? stats[stats.length - 1] : null;
  };

  const getAverageCpuUsage = (): number => {
    if (stats.length === 0) return 0;
    const sum = stats.reduce((acc, stat) => acc + stat.cpuUsage, 0);
    return sum / stats.length;
  };

  const getMemoryUsagePercentage = (): number => {
    const latestStats = getLatestStats();
    if (!latestStats) return 0;
    return (latestStats.memory.used / latestStats.memory.total) * 100;
  };

  const getDiskUsagePercentage = (disk: DiskInfo): number => {
    if (!disk.totalSpace) return 0;
    return (disk.usedSpace / disk.totalSpace) * 100;
  };

  useEffect(() => {
    connectWebSocket();
    return () => disconnectWebSocket();
  }, []);

  useEffect(() => {
    if (error) toast.error(error);
  }, [error]);

  const latestStats = getLatestStats();
  const averageCpuUsage = getAverageCpuUsage();
  const memoryUsagePercentage = getMemoryUsagePercentage();

  if (isLoading) return <FullScreenLoading />;

  return (
    <div className="min-h-screen bg-background">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between py-6 border-b border-border shrink-0 gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-foreground">Dashboard</h1>
          <p className="text-muted-foreground mt-1">System monitoring and performance overview</p>
        </div>
        <div className="flex items-center gap-2">
          <div className={`w-3 h-3 rounded-full ${isConnected ? 'bg-green-500' : 'bg-red-500'}`} />
          <span className="text-sm text-muted-foreground">{isConnected ? 'Connected' : 'Disconnected'}</span>
        </div>
      </div>

      {error && (
        <div className="mt-4">
          <Alert variant="destructive">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        </div>
      )}

      <div className="py-6 space-y-6">
        <SystemOverview stats={latestStats} />

        {stats.length > 0 && (
          <div className="grid gap-6 grid-cols-1 lg:grid-cols-2">
            <ChartCard title="CPU Usage" data={stats} dataKey="cpuUsage" color="#8B5CF6" formatter={formatPercentage} />
            <ChartCard title="Memory Usage" data={stats} dataKey="memory.used" color="#A371F7" formatter={formatMemory} />
          </div>
        )}

        {stats.length > 0 && (
          <div className="grid gap-4 grid-cols-1 sm:grid-cols-2 xl:grid-cols-4">
            <MetricCard title="Current CPU" value={latestStats?.cpuUsage || 0} showUsageColor />
            <MetricCard title="Average CPU" value={averageCpuUsage} showUsageColor />
            <MetricCard title="Memory Usage" value={memoryUsagePercentage} showUsageColor />
            <MetricCard
              title="Load Average"
              value={latestStats?.loadAverage.oneMinute || 0}
              formatter={(value: number) => value.toFixed(2)}
            />
          </div>
        )}

        {/* Disk Usage */}
        {latestStats?.disk && latestStats.disk.length > 0 && (
          <div className="space-y-4">
            <h2 className="text-xl font-semibold tracking-tight text-foreground">Disk Usage</h2>

            <div className="grid gap-4 grid-cols-1 sm:grid-cols-2 xl:grid-cols-3">
              {latestStats.disk.map((d) => (
                <MetricCard
                  key={d.name}
                  title={`Disk: ${d.name}`}
                  value={getDiskUsagePercentage(d)}
                  showUsageColor
                  formatter={(v: number) => `${v.toFixed(1)}% (${formatMemory(d.usedSpace)} / ${formatMemory(d.totalSpace)})`}
                />
              ))}
            </div>
          </div>
        )}

        {stats.length === 0 && !isLoading && (
          <div className="flex flex-col items-center justify-center py-12">
            <p className="text-muted-foreground text-lg mb-4">
              {isConnected ? 'Waiting for system data...' : 'Unable to connect to system monitoring'}
            </p>
          </div>
        )}
      </div>
    </div>
  );
}
