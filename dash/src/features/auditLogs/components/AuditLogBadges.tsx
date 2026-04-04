import { Badge } from '@/components/ui/badge';
import { 
  getTriggerIcon, 
  getTriggerBadgeVariant,
  getActionBadgeColor 
} from '../utils/auditLogFormatters';

interface TriggerBadgeProps {
  triggerType: string;
}

export const TriggerBadge = ({ triggerType }: TriggerBadgeProps) => {
  const Icon = getTriggerIcon(triggerType);
  const variant = getTriggerBadgeVariant(triggerType);
  
  return (
    <Badge variant={variant} className="flex items-center gap-1">
      <Icon className="h-4 w-4" />
      {triggerType}
    </Badge>
  );
};

interface ActionBadgeProps {
  action: string;
}

export const ActionBadge = ({ action }: ActionBadgeProps) => {
  const color = getActionBadgeColor(action);
  
  return (
    <Badge className={color}>
      {action}
    </Badge>
  );
};

interface ResourceBadgeProps {
  resourceType: string;
  resourceId?: number;
}

export const ResourceBadge = ({ resourceType, resourceId }: ResourceBadgeProps) => {
  return (
    <>
      <Badge variant="outline">{resourceType}</Badge>
      {resourceId && (
        <span className="text-xs text-gray-500">
          #{resourceId}
        </span>
      )}
    </>
  );
};
