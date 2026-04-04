import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { getResourceFilters } from '../utils/auditLogFormatters';

interface AuditLogFiltersProps {
  resourceFilter: string;
  onResourceFilterChange: (value: string) => void;
  total: number;
}

export const AuditLogFilters = ({ 
  resourceFilter, 
  onResourceFilterChange,
  total 
}: AuditLogFiltersProps) => {
  const filters = getResourceFilters();

  return (
    <div className="flex flex-col sm:flex-row sm:items-center gap-4">
      <div className="flex items-center gap-2">
        <span className="text-sm text-gray-400">Filter by resource:</span>
        <Select value={resourceFilter} onValueChange={onResourceFilterChange}>
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder="Select resource" />
          </SelectTrigger>
          <SelectContent>
            {filters.map((filter) => (
              <SelectItem key={filter.value} value={filter.value}>
                {filter.label}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>
      <div className="text-sm text-gray-400">
        Total: {total} logs
      </div>
    </div>
  );
};
