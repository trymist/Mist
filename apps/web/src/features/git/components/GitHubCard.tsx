import { Card, CardHeader, CardTitle, CardContent, CardFooter } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { ExternalLink, Github } from 'lucide-react';
import { getGitHubInstallUrl, getGitHubManageUrl } from '../utils';
import type { GitHubApp, User } from '@/types';

interface GitHubCardProps {
  app: GitHubApp | null;
  isInstalled: boolean;
  user: User | null;
  onCreateApp: () => void;
}

export function GitHubCard({ app, isInstalled, user, onCreateApp }: GitHubCardProps) {
  return (
    <Card className="border-border bg-card hover:border-primary transition-colors">
      <CardHeader className="pb-2 flex flex-row items-center gap-2">
        <Github className="w-5 h-5 text-muted-foreground" />
        <CardTitle className="text-lg font-semibold">GitHub</CardTitle>
      </CardHeader>

      <CardContent>
        {app ? (
          <div className="space-y-3">
            <div>
              <p className="font-medium text-foreground">{app.name}</p>
              <p className="text-sm text-muted-foreground mt-1">
                App ID: {app.appId}
              </p>
            </div>

            <div className="flex flex-wrap gap-2 pt-2">
              <Button
                onClick={() => {
                  if (user) {
                    const installUrl = getGitHubInstallUrl(
                      app.slug,
                      app.appId,
                      parseInt(user.id.toString())
                    );
                    window.open(installUrl);
                  }
                }}
                disabled={isInstalled || !user}
                className="transition-colors disabled:cursor-not-allowed"
              >
                {isInstalled ? 'Installed' : 'Install App'}
                <ExternalLink className="w-4 h-4 ml-2" />
              </Button>

              <Button
                variant="outline"
                onClick={() => window.open(getGitHubManageUrl(app.slug), '_blank')}
              >
                Manage
              </Button>
            </div>
          </div>
        ) : (
          <div className="space-y-3">
            <p className="text-sm text-muted-foreground">
              No GitHub App connected yet. Create one to enable Git deployments.
            </p>
            <Button
              onClick={onCreateApp}
              size="default"
              className="transition-colors"
            >
              Create GitHub App
            </Button>
          </div>
        )}
      </CardContent>

      {app && (
        <CardFooter className="text-xs text-muted-foreground border-t border-border pt-2">
          Created at: {new Date(app.createdAt).toLocaleString()}
        </CardFooter>
      )}
    </Card>
  );
}
