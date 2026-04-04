import type { AuditLog } from '@/types';
import { format } from 'date-fns';
import { TriggerBadge, ActionBadge, ResourceBadge } from './AuditLogBadges';
import { AuditLogDetails } from './AuditLogDetails';

interface AuditLogItemProps {
  log: AuditLog;
}

export const AuditLogItem = ({ log }: AuditLogItemProps) => {
  const getDescriptionText = () => {
    if (log.triggerType === 'user') {
      return (
        <>
          <span className="font-medium text-gray-300">
            {log.username || log.email || 'Unknown User'}
          </span>
          <span className="text-gray-400"> performed </span>
          <span className="font-medium text-gray-300">{log.action}</span>
          <span className="text-gray-400"> on </span>
          <span className="font-medium text-gray-300">{log.resourceType}</span>
        </>
      );
    }
    
    if (log.triggerType === 'webhook') {
      return (
        <>
          <span className="font-medium text-gray-300">GitHub Webhook</span>
          <span className="text-gray-400"> triggered </span>
          <span className="font-medium text-gray-300">{log.action}</span>
          <span className="text-gray-400"> on </span>
          <span className="font-medium text-gray-300">{log.resourceType}</span>
        </>
      );
    }
    
    return (
      <>
        <span className="font-medium text-gray-300">System</span>
        <span className="text-gray-400"> performed </span>
        <span className="font-medium text-gray-300">{log.action}</span>
        <span className="text-gray-400"> on </span>
        <span className="font-medium text-gray-300">{log.resourceType}</span>
      </>
    );
  };

  return (
    <div className="flex items-start gap-4 p-4 rounded-lg border border-gray-800 hover:border-gray-700 transition-colors">
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2 flex-wrap">
          <TriggerBadge triggerType={log.triggerType} />
          <ActionBadge action={log.action} />
          <ResourceBadge 
            resourceType={log.resourceType} 
            resourceId={log.resourceId} 
          />
        </div>
        
        <div className="mt-2">
          <p className="text-sm">
            {getDescriptionText()}
          </p>
          <AuditLogDetails log={log} />
        </div>
      </div>
      
      <div className="text-right text-xs text-gray-500 whitespace-nowrap">
        {format(new Date(log.createdAt), 'MMM dd, yyyy')}
        <br />
        {format(new Date(log.createdAt), 'HH:mm:ss')}
      </div>
    </div>
  );
};
