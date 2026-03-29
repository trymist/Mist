import type { AuditLog } from '@/types';
import { AuditLogItem } from './AuditLogItem';

interface AuditLogListProps {
  logs: AuditLog[];
  loading?: boolean;
}

export const AuditLogList = ({ logs, loading }: AuditLogListProps) => {
  if (loading && logs.length === 0) {
    return (
      <div className="text-center py-12 text-gray-400">
        Loading audit logs...
      </div>
    );
  }

  if (logs.length === 0) {
    return (
      <div className="text-center py-12 text-gray-400">
        No audit logs found
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {logs.map((log) => (
        <AuditLogItem key={log.id} log={log} />
      ))}
    </div>
  );
};
