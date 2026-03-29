import { useEffect, useState } from "react"
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { toast } from "sonner"
import { DeploymentMonitor } from "@/components/deployments"
import type { Deployment, App } from "@/types"
import { Loader2, Clock, CheckCircle2, XCircle, PlayCircle, AlertCircle, Square, SquareSlash } from "lucide-react"
import { deploymentsService } from "@/services"
import { cn } from "@/lib/utils"

export const DeploymentsTab = ({ appId, app }: { appId: number; app?: App }) => {
  const [deployments, setDeployments] = useState<Deployment[]>([])
  const [loading, setLoading] = useState(true)
  const [deploying, setDeploying] = useState(false)
  const [stoppingIds, setStoppingIds] = useState<Set<number>>(new Set())
  const [selectedDeployment, setSelectedDeployment] = useState<number | null>(null)

  const fetchDeployments = async () => {
    try {
      setLoading(true)
      const data = await deploymentsService.getByAppId(appId)
      setDeployments(data || [])
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to fetch deployments")
    } finally {
      setLoading(false)
    }
  }

  const handleDeploy = async () => {
    try {
      setDeploying(true)

      const deployment = await deploymentsService.create(appId)
      toast.success('Deployment started successfully')

      // Open the monitor immediately
      setSelectedDeployment(deployment.id)

      // Refresh deployments list
      await fetchDeployments()
    } catch (error) {
      console.error('Deployment error:', error)
      toast.error(error instanceof Error ? error.message : 'Failed to start deployment')
    } finally {
      setDeploying(false)
    }
  }

  const handleDeploymentComplete = () => {
    toast.success('Deployment completed successfully!')
    fetchDeployments()
  }

  const handleStop = async (deploymentId: number) => {
    if (stoppingIds.has(deploymentId)) return
    try {
      setStoppingIds(prev => new Set(prev).add(deploymentId))
      await deploymentsService.stopDeployment(deploymentId)
      toast.success('Deployment stopped')
      fetchDeployments()
    } catch (error) {
      console.error('Stop deployment error:', error)
      toast.error(error instanceof Error ? error.message : 'Failed to stop deployment')
    } finally {
      setStoppingIds(prev => {
        const next = new Set(prev)
        next.delete(deploymentId)
        return next
      })
    }
  }

  useEffect(() => {
    fetchDeployments()

    // Auto-refresh deployments every 10 seconds to catch updates
    const interval = setInterval(fetchDeployments, 10000)
    return () => clearInterval(interval)
  }, [appId])

  // Helper to get status badge with styles
  const getStatusBadge = (deployment: Deployment) => {
    const { status, stage } = deployment

    switch (status) {
      case 'success':
        return (
          <Badge className="bg-green-500 text-white flex items-center gap-1.5 border-0">
            <CheckCircle2 className="h-3 w-3" />
            Success
          </Badge>
        )
      case 'failed':
        return (
          <Badge variant="destructive" className="flex items-center gap-1.5">
            <XCircle className="h-3 w-3" />
            Failed
          </Badge>
        )
      case 'building':
      case 'deploying':
      case 'cloning':
        return (
          <Badge className="bg-blue-500 text-white flex items-center gap-1.5 animate-pulse border-0">
            <Loader2 className="h-3 w-3 animate-spin" />
            {stage.charAt(0).toUpperCase() + stage.slice(1)}
          </Badge>
        )
      case 'pending':
        return (
          <Badge variant="outline" className="flex items-center gap-1.5 border-slate-300 text-slate-600">
            <AlertCircle className="h-3 w-3" />
            Pending
          </Badge>
        )
      case 'stopped':
        return (
          <Badge className="bg-slate-400 text-white flex items-center gap-1.5 border-0">
            <SquareSlash className="h-3 w-3" />
            Stopped
          </Badge>
        )
      default:
        return <Badge variant="outline">{status}</Badge>
    }
  }

  // Check if deployment can be stopped
  const canStopDeployment = (deployment: Deployment) => {
    return ['pending', 'building', 'deploying', 'cloning'].includes(deployment.status) && deployment.status !== 'stopped'
  }

  // Get border class based on status
  const getDeploymentBorderClass = (status: string) => {
    switch (status) {
      case 'success':
        return 'border-green-200 dark:border-green-800 bg-green-50/50 dark:bg-green-950/20'
      case 'failed':
        return 'border-red-200 dark:border-red-800 bg-red-50/50 dark:bg-red-950/20'
      case 'stopped':
        return 'border-slate-200 dark:border-slate-700 bg-slate-50/50 dark:bg-slate-900/20'
      default:
        return 'border-border bg-muted/20 hover:bg-muted/30'
    }
  }

  return (
    <>
      {/* Deployment Monitor */}
      {selectedDeployment && (
        <DeploymentMonitor
          deploymentId={selectedDeployment}
          open={!!selectedDeployment}
          onClose={() => setSelectedDeployment(null)}
          onComplete={handleDeploymentComplete}
        />
      )}

      {/* Deployments Card */}
      <Card>
        <CardHeader className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <CardTitle>Deployments</CardTitle>
          <Button
            onClick={handleDeploy}
            disabled={deploying}
            className="flex items-center gap-2 w-full sm:w-auto"
          >
            {deploying ? (
              <>
                <Loader2 className="h-4 w-4 animate-spin" />
                Deploying...
              </>
            ) : (
              <>
                <PlayCircle className="h-4 w-4" />
                Deploy Now
              </>
            )}
          </Button>
        </CardHeader>

        <CardContent className="space-y-4">
          {loading && (
            <div className="flex items-center justify-center py-8 text-muted-foreground">
              <Loader2 className="h-6 w-6 animate-spin mr-2" />
              Loading deployments...
            </div>
          )}

          {!loading && deployments.length === 0 && (
            <div className="text-center py-12">
              <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-muted mb-4">
                <PlayCircle className="h-8 w-8 text-muted-foreground" />
              </div>
              <p className="text-muted-foreground mb-2">No deployments yet</p>
              <p className="text-sm text-muted-foreground">
                Click "Deploy Now" to create your first deployment
              </p>
            </div>
          )}

          {!loading && deployments.length > 0 && (
            <div className="space-y-3">
              {deployments.map((d) => (
                <div
                  key={d.id}
                  className={cn(
                    'flex flex-col sm:flex-row sm:items-start sm:justify-between p-4 rounded-lg border transition-colors gap-4',
                    getDeploymentBorderClass(d.status)
                  )}
                >
                  <div className="flex-1 space-y-2">
                    <div className="flex items-center gap-3 flex-wrap">
                      {getStatusBadge(d)}

                      <span className="text-xs text-muted-foreground font-mono">
                        #{d.id}
                      </span>

                      {/* Progress indicator for in-progress deployments */}
                      {d.status !== 'success' && d.status !== 'failed' && d.status !== 'stopped' && d.progress > 0 && (
                        <div className="flex items-center gap-2 w-full sm:w-auto">
                          <div className="w-24 bg-muted rounded-full h-1.5 overflow-hidden">
                            <div
                              className="bg-primary h-full transition-all duration-300"
                              style={{ width: `${d.progress}%` }}
                            />
                          </div>
                          <span className="text-xs text-muted-foreground">
                            {d.progress}%
                          </span>
                        </div>
                      )}
                    </div>

                    <div className="space-y-1">
                      {app?.appType === 'database' ? (
                        <p className="font-mono text-sm break-all">
                          <span className="text-primary">Version: {d.commit_hash}</span>
                          {d.commit_message && (
                            <>
                              {' – '}
                              {d.commit_message}
                            </>
                          )}
                        </p>
                      ) : (
                        <p className="font-mono text-sm break-all">
                          <span className="text-primary">{d.commit_hash.slice(0, 7)}</span>
                          {' – '}
                          {d.commit_message}
                        </p>
                      )}

                      {/* Show error message for failed deployments */}
                      {d.error_message && d.status === 'failed' && (
                        <p className="text-xs text-red-500 flex items-start gap-1">
                          <XCircle className="h-3 w-3 mt-0.5 shrink-0" />
                          <span className="break-all">{d.error_message}</span>
                        </p>
                      )}

                      {/* Show stopped message for stopped deployments */}
                      {d.status === 'stopped' && (
                        <p className="text-xs text-slate-500 flex items-start gap-1">
                          <SquareSlash className="h-3 w-3 mt-0.5 shrink-0" />
                          <span className="break-all">Deployment was stopped by user</span>
                        </p>
                      )}
                    </div>

                    <div className="flex flex-col sm:flex-row sm:items-center gap-2 sm:gap-4 text-xs text-muted-foreground">
                      <span className="flex items-center gap-1">
                        <Clock className="h-3 w-3" />
                        {new Date(d.created_at).toLocaleString()}
                      </span>

                      {d.duration && (
                        <span>
                          Duration: {d.duration}s
                        </span>
                      )}
                    </div>
                  </div>

                  <div className="flex gap-2">
                    {canStopDeployment(d) && (
                      <Button
                        variant="destructive"
                        size="sm"
                        onClick={() => handleStop(d.id)}
                        disabled={stoppingIds.has(d.id)}
                        className="flex items-center gap-1.5"
                      >
                        {stoppingIds.has(d.id) ? (
                          <Loader2 className="h-3.5 w-3.5 animate-spin" />
                        ) : (
                          <Square className="h-3.5 w-3.5" />
                        )}
                        Stop
                      </Button>
                    )}
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => setSelectedDeployment(d.id)}
                      className="flex items-center gap-1.5"
                    >
                      View Logs
                    </Button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </>
  )
}
