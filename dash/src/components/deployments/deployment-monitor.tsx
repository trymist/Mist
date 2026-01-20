import { useEffect, useRef, useState } from 'react';
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Terminal, CheckCircle2, XCircle, AlertCircle, Loader2, Square } from 'lucide-react';
import { useDeploymentMonitor } from '@/hooks';
import { deploymentsService } from '@/services';
import { LogLine } from '@/components/logs/log-line';
import { cn } from '@/lib/utils';
import { toast } from 'sonner';

interface Props {
  deploymentId: number;
  open: boolean;
  onClose: () => void;
  onComplete?: () => void;
}

export const DeploymentMonitor = ({ deploymentId, open, onClose, onComplete }: Props) => {
  const bottomRef = useRef<HTMLDivElement>(null);

  const completedRef = useRef(false);
  const [stopping, setStopping] = useState(false);

  const { logs, status, error, isConnected, isLoading, isLive, reset } = useDeploymentMonitor({
    deploymentId,
    enabled: open,
    onComplete: () => {
      if (!completedRef.current) {
        completedRef.current = true;
        onComplete?.();
      }
    },
    onError: (err) => {
      console.error('Deployment error:', err);
      toast.error(err);
    },
    onClose: () => {
      handleClose();
    },
  });

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [logs]);

  const handleClose = () => {
    completedRef.current = false;
    reset();
    onClose();
  };

  const handleStop = async () => {
    if (stopping) return;
    try {
      setStopping(true);
      await deploymentsService.stopDeployment(deploymentId);
      toast.success('Deployment stopped');
    } catch (err) {
      console.error('Stop deployment error:', err);
      toast.error(err instanceof Error ? err.message : 'Failed to stop deployment');
    } finally {
      setStopping(false);
    }
  };

  const getStatusInfo = () => {
    const statusValue = status?.status || 'pending';

    switch (statusValue) {
      case 'success':
        return {
          color: 'bg-green-500 text-white',
          icon: <CheckCircle2 className="h-4 w-4" />,
          label: 'Success',
        };
      case 'failed':
        return {
          color: 'bg-red-500 text-white',
          icon: <XCircle className="h-4 w-4" />,
          label: 'Failed',
        };
      case 'building':
      case 'deploying':
      case 'cloning':
        return {
          color: 'bg-blue-500 text-white animate-pulse',
          icon: <Loader2 className="h-4 w-4 animate-spin" />,
          label: statusValue.charAt(0).toUpperCase() + statusValue.slice(1),
        };
      case 'pending':
        return {
          color: 'bg-yellow-500 text-white',
          icon: <AlertCircle className="h-4 w-4" />,
          label: 'Pending',
        };
      default:
        return {
          color: 'bg-gray-500 text-white',
          icon: null,
          label: statusValue,
        };
    }
  };

  const canStop = () => {
    const statusValue = status?.status || 'pending';
    return ['pending', 'building', 'deploying', 'cloning'].includes(statusValue);
  };

  const statusInfo = getStatusInfo();

  return (
    <Sheet open={open} onOpenChange={handleClose}>
      <SheetContent
        side="right"
        className="w-full sm:w-[90vw] sm:max-w-[90vw] lg:w-[85vw] lg:max-w-[85vw] p-0 flex flex-col gap-0"
      >
        {/* Header */}
        <SheetHeader className="px-10 py-4 border-b bg-background/80 backdrop-blur flex flex-col sm:flex-row sm:justify-between sm:items-center gap-3 shrink-0">
          <SheetTitle className="flex items-center gap-3 text-lg">
            <Terminal className="h-5 w-5 text-primary" />
            <span>Deployment Monitor</span>
          </SheetTitle>

          <div className="flex items-center gap-3">
            {/* Connection Status */}
            <div className="flex items-center gap-2 text-xs">
              <div
                className={cn(
                  'w-2 h-2 rounded-full',
                  isLive ? (isConnected ? 'bg-green-500' : 'bg-red-500') : 'bg-blue-500'
                )}
              />
              <span className="text-muted-foreground">
                {isLive
                  ? (isConnected ? 'Live' : 'Disconnected')
                  : 'Completed'
                }
              </span>
            </div>

            {/* Deployment ID Badge */}
            <Badge variant="outline" className="font-mono text-xs px-2 py-0.5">
              #{deploymentId}
            </Badge>

            {/* Stop Button */}
            {canStop() && (
              <Button
                variant="destructive"
                size="sm"
                onClick={handleStop}
                disabled={stopping}
                className="flex items-center gap-1.5"
              >
                {stopping ? (
                  <Loader2 className="h-3.5 w-3.5 animate-spin" />
                ) : (
                  <Square className="h-3.5 w-3.5" />
                )}
                Stop
              </Button>
            )}
          </div>
        </SheetHeader>

        {/* Status Bar */}
        {status && (
          <div className="px-6 py-3 bg-muted/30 border-b flex items-center justify-between shrink-0">
            <div className="flex items-center gap-3">
              <Badge className={cn('flex items-center gap-1.5', statusInfo.color)}>
                {statusInfo.icon}
                <span>{statusInfo.label}</span>
              </Badge>
              <span className="text-sm text-muted-foreground">
                {status.message}
              </span>
            </div>

            {/* Progress Bar */}
            {status.status !== 'success' && status.status !== 'failed' && (
              <div className="flex items-center gap-3 min-w-[200px]">
                <div className="flex-1 bg-muted rounded-full h-2 overflow-hidden">
                  <div
                    className="bg-primary h-full transition-all duration-300 ease-out"
                    style={{ width: `${status.progress}%` }}
                  />
                </div>
                <span className="text-sm text-muted-foreground font-medium min-w-[3ch] text-right">
                  {status.progress}%
                </span>
              </div>
            )}
          </div>
        )}

        {/* Error Banner - Only show for live deployments */}
        {error && isLive && (
          <div className="px-6 py-3 bg-red-500/10 border-b border-red-500/30 flex items-start gap-3 shrink-0">
            <XCircle className="h-5 w-5 text-red-500 mt-0.5 shrink-0" />
            <div className="flex-1">
              <h3 className="font-semibold text-red-500">Deployment Failed</h3>
              <p className="text-sm text-red-500/90 mt-1">{error}</p>
            </div>
          </div>
        )}

        {/* Success Banner - Only show for live deployments */}
        {status?.status === 'success' && isLive && (
          <div className="px-6 py-3 bg-green-500/10 border-b border-green-500/30 flex items-center gap-3 shrink-0">
            <CheckCircle2 className="h-5 w-5 text-green-500" />
            <span className="font-semibold text-green-500">
              Deployment Successful!
            </span>
            {status.duration && (
              <span className="text-sm text-green-500/80">
                Completed in {status.duration}s
              </span>
            )}
          </div>
        )}

        {/* Logs Viewer */}
        <div className="flex-1 bg-slate-950 p-4 overflow-auto">
          {isLoading ? (
            <div className="flex flex-col items-center justify-center h-full text-gray-500">
              <Loader2 className="h-8 w-8 animate-spin mb-3" />
              <p>Loading deployment...</p>
            </div>
          ) : logs.length === 0 ? (
            <div className="flex flex-col items-center justify-center h-full text-gray-500">
              <Loader2 className="h-8 w-8 animate-spin mb-3" />
              <p>Waiting for logs...</p>
            </div>
          ) : (
            <div className="space-y-0.5">
              {logs.map((log, index) => (
                <LogLine
                  key={index}
                  line={log.line}
                  index={index}
                  showLineNumbers={false}
                  streamType={log.stream}
                />
              ))}
              <div ref={bottomRef} />
            </div>
          )}
        </div>

        {/* Footer */}
      </SheetContent>
    </Sheet>
  );
};
