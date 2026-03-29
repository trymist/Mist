import type { AuditLog } from '@/types';
import { 
  parseAuditLogDetails,
  formatWebhookDetails,
  formatChangeDetails,
  formatGenericDetails
} from '../utils/auditLogFormatters';

interface AuditLogDetailsProps {
  log: AuditLog;
}

export const AuditLogDetails = ({ log }: AuditLogDetailsProps) => {
  const details = parseAuditLogDetails(log.details);
  
  if (!details) return null;

  // For webhook events
  if (log.triggerType === 'webhook') {
    const webhookItems = formatWebhookDetails(details);
    
    if (webhookItems.length === 0) return null;
    
    return (
      <div className="text-xs text-gray-400 mt-1 space-y-0.5">
        {webhookItems.map((item) => (
          <div key={item.label}>
            {item.label}: {item.value}
          </div>
        ))}
      </div>
    );
  }

  // For user actions with changes
  const changes = formatChangeDetails(details);
  if (changes) {
    return (
      <div className="text-xs text-gray-400 mt-1">
        Changed: {changes}
      </div>
    );
  }

  // Generic details
  const genericItems = formatGenericDetails(details);
  if (genericItems.length > 0) {
    return (
      <div className="text-xs text-gray-400 mt-1">
        {genericItems.map((item) => (
          <div key={item.key}>
            {item.key}: {item.value}
            {item.truncated ? '...' : ''}
          </div>
        ))}
      </div>
    );
  }

  return null;
};
