import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { getRoleStyles, getUserInitials, formatUserId } from '../utils';
import type { User } from '@/types';

interface UserCardProps {
  user: User;
  onClick?: (user: User) => void;
}

export function UserCard({ user, onClick }: UserCardProps) {
  return (
    <Card
      className="cursor-pointer border-border bg-card hover:border-primary transition-colors"
      onClick={() => onClick?.(user)}
    >
      <CardHeader className="pb-2">
        <div className="flex items-start justify-between">
          <div className="flex items-center gap-3">
            <div className="flex h-9 w-9 items-center justify-center rounded-full bg-muted text-foreground overflow-hidden border-2 border-border">
              {user.avatarUrl ? (
                <img src={user.avatarUrl} alt={user.username} className="h-full w-full object-cover" />
              ) : (
                getUserInitials(user.username)
              )}
            </div>
            <div>
              <CardTitle className="text-lg font-semibold text-foreground">
                {user.username}
              </CardTitle>
              <CardDescription className="text-sm text-muted-foreground">
                {user.email}
              </CardDescription>
            </div>
          </div>
          <Badge
            variant="secondary"
            className={`capitalize ${getRoleStyles(user.role)}`}
          >
            {user.role}
          </Badge>
        </div>
      </CardHeader>
      <CardContent className="border-t border-border pt-3">
        <p className="text-sm text-muted-foreground font-mono break-all">
          {formatUserId(user.id)}
        </p>
      </CardContent>
    </Card>
  );
}
