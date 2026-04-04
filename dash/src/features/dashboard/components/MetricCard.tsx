import { Card } from '@/components/ui/card';
import { formatPercentage, getUsageColor } from '../utils';

interface MetricCardProps {
  title: string;
  value: number;
  unit?: string;
  formatter?: (value: number) => string;
  showUsageColor?: boolean;
}

export function MetricCard({ 
  title, 
  value, 
  unit = '', 
  formatter = formatPercentage,
  showUsageColor = false 
}: MetricCardProps) {
  const displayValue = formatter(value);
  const colorClass = showUsageColor ? getUsageColor(value) : '';

  return (
    <Card className="p-4">
      <p className="text-sm text-muted-foreground">{title}</p>
      <p className={`text-lg font-semibold ${colorClass}`}>
        {displayValue}{unit}
      </p>
    </Card>
  );
}
