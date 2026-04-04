import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { GitFork, Gitlab, GitMerge } from 'lucide-react';
import { cn } from '@/lib/utils';

interface ProviderCardProps {
  name: string;
  icon: 'GitFork' | 'Gitlab' | 'GitMerge';
  isComingSoon?: boolean;
}

const iconMap = {
  GitFork,
  Gitlab,
  GitMerge,
};

export function ProviderCard({ name, icon, isComingSoon = true }: ProviderCardProps) {
  const Icon = iconMap[icon];

  return (
    <Card
      className={cn(
        'h-full flex flex-col items-center justify-between border border-dashed border-border bg-card hover:border-primary/30 transition-colors',
        isComingSoon && 'cursor-not-allowed opacity-60'
      )}
    >
      <CardHeader className="flex flex-col items-center space-y-3 pb-4">
        <Icon className="w-6 h-6 text-muted-foreground" />
        <CardTitle className="text-base font-medium">{name}</CardTitle>
      </CardHeader>
      <CardContent className="pb-6">
        {isComingSoon && (
          <Badge variant="secondary">Coming Soon</Badge>
        )}
      </CardContent>
    </Card>
  );
}
