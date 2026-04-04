import { ResponsiveContainer, AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip } from 'recharts';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import type { SystemStats } from '../DashboardPage';

interface ChartCardProps {
  title: string;
  data: SystemStats[];
  dataKey: string;
  color: string;
  formatter?: (value: number) => string;
}

export function ChartCard({ title, data, dataKey, color, formatter }: ChartCardProps) {
  const customTooltip = ({ active, payload, label }: { 
    active?: boolean; 
    payload?: Array<{ value: number; color: string }>; 
    label?: string | number;
  }) => {
    if (active && payload && payload.length) {
      const value = payload[0].value;
      const formattedValue = formatter ? formatter(value) : value;
      const timestamp = typeof label === 'number' ? label : parseInt(String(label), 10);
      
      return (
        <div className="bg-popover p-3 border border-border rounded-md">
          <p className="text-foreground">
            {new Date(timestamp * 1000).toLocaleTimeString()}
          </p>
          <p style={{ color: payload[0].color }}>
            {title}: {formattedValue}
          </p>
        </div>
      );
    }
    return null;
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>{title}</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="h-[300px]">
          <ResponsiveContainer width="100%" height="100%">
            <AreaChart data={data}>
              <CartesianGrid strokeDasharray="3 3" stroke="#30363D" />
              <XAxis
                dataKey="timestamp"
                stroke="#C9D1D9"
                tickFormatter={(timestamp) =>
                  new Date(timestamp * 1000).toLocaleTimeString()
                }
              />
              <YAxis
                stroke="#C9D1D9"
                tickFormatter={formatter}
              />
              <Tooltip content={customTooltip} />
              <Area
                type="monotone"
                dataKey={dataKey}
                stroke={color}
                fill={color}
                fillOpacity={0.3}
              />
            </AreaChart>
          </ResponsiveContainer>
        </div>
      </CardContent>
    </Card>
  );
}
