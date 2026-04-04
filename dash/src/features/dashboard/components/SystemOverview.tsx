import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { 
  formatMemory, 
  formatUptime, 
  formatTemperature, 
  formatLoadAverage,
  getSystemHealthStatus,
  getHealthStatusColor 
} from '../utils';
import type { SystemStats } from '../DashboardPage';

interface SystemOverviewProps {
  stats: SystemStats | null;
}

export function SystemOverview({ stats }: SystemOverviewProps) {
  if (!stats) return null;

  const memoryPercentage = (stats.memory.used / stats.memory.total) * 100;
  const diskUsage = stats.disk.length > 0 ? 
    (stats.disk[0].usedSpace / stats.disk[0].totalSpace) * 100 : 0;
  
  const healthStatus = getSystemHealthStatus(
    stats.cpuUsage,
    memoryPercentage,
    diskUsage
  );

  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">System Health</CardTitle>
          <Badge 
            variant={healthStatus === 'healthy' ? 'default' : 'destructive'}
            className={getHealthStatusColor(healthStatus)}
          >
            {healthStatus.charAt(0).toUpperCase() + healthStatus.slice(1)}
          </Badge>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">
            {stats.uptime ? formatUptime(stats.uptime) : 'N/A'}
          </div>
          <p className="text-xs text-muted-foreground">System Uptime</p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">CPU Usage</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">
            {stats.cpuUsage.toFixed(1)}%
          </div>
          <p className="text-xs text-muted-foreground">
            {stats.cpuTemperature ? `${formatTemperature(stats.cpuTemperature)}` : 'Temp: N/A'}
          </p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Memory Usage</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">
            {memoryPercentage.toFixed(1)}%
          </div>
          <p className="text-xs text-muted-foreground">
            {formatMemory(stats.memory.used)} / {formatMemory(stats.memory.total)}
          </p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Load Average</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">
            {formatLoadAverage(stats.loadAverage.oneMinute)}
          </div>
          <p className="text-xs text-muted-foreground">
            5min: {formatLoadAverage(stats.loadAverage.fiveMinutes)} | 
            15min: {formatLoadAverage(stats.loadAverage.fifteenMinutes)}
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
