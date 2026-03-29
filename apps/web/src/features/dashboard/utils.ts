export const formatMemory = (bytes: number): string => {
  const gb = bytes / (1024 * 1024 * 1024);
  return `${gb.toFixed(2)} GB`;
};

export const formatUptime = (seconds: number): string => {
  const days = Math.floor(seconds / (24 * 60 * 60));
  const hours = Math.floor((seconds % (24 * 60 * 60)) / (60 * 60));
  const minutes = Math.floor((seconds % (60 * 60)) / 60);
  
  const parts = [];
  if (days > 0) parts.push(`${days}d`);
  if (hours > 0) parts.push(`${hours}h`);
  if (minutes > 0) parts.push(`${minutes}m`);
  
  return parts.join(' ') || '0m';
};

export const formatPercentage = (value: number): string => {
  return `${value.toFixed(1)}%`;
};

export const formatTemperature = (celsius: number): string => {
  return `${celsius.toFixed(1)}°C`;
};

export const getUsageColor = (percentage: number): string => {
  if (percentage >= 80) return 'text-red-500';
  if (percentage >= 60) return 'text-yellow-500';
  return 'text-green-500';
};

export const getDiskUsagePercentage = (used: number, total: number): number => {
  return (used / total) * 100;
};

export const formatBytes = (bytes: number): string => {
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  if (bytes === 0) return '0 B';
  
  const i = Math.floor(Math.log(bytes) / Math.log(1024));
  return `${(bytes / Math.pow(1024, i)).toFixed(2)} ${sizes[i]}`;
};

export const formatLoadAverage = (load: number): string => {
  return load.toFixed(2);
};

export const getSystemHealthStatus = (
  cpuUsage: number,
  memoryPercentage: number,
  diskPercentage: number
): 'healthy' | 'warning' | 'critical' => {
  const maxUsage = Math.max(cpuUsage, memoryPercentage, diskPercentage);
  
  if (maxUsage >= 90) return 'critical';
  if (maxUsage >= 75) return 'warning';
  return 'healthy';
};

export const getHealthStatusColor = (status: 'healthy' | 'warning' | 'critical'): string => {
  switch (status) {
    case 'healthy':
      return 'text-green-500';
    case 'warning':
      return 'text-yellow-500';
    case 'critical':
      return 'text-red-500';
    default:
      return 'text-gray-500';
  }
};
