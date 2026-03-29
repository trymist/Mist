import { useState, useEffect } from "react"
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { toast } from "sonner"
import { applicationsService } from "@/services"
import type { ContainerStatus } from "@/types"
import { 
  Loader2, 
  Power, 
  PowerOff, 
  RefreshCw, 
  Activity,
  CheckCircle2,
  XCircle,
  Clock
} from "lucide-react"

interface ContainerControlsProps {
  appId: number
  onStatusChange?: () => void
}

export const ContainerControls = ({ appId, onStatusChange }: ContainerControlsProps) => {
  const [status, setStatus] = useState<ContainerStatus | null>(null)
  const [loading, setLoading] = useState(false)
  const [actionLoading, setActionLoading] = useState<string | null>(null)

  const fetchStatus = async () => {
    try {
      setLoading(true)
      const data = await applicationsService.getContainerStatus(appId)
      setStatus(data)
    } catch (error) {
      console.error("Failed to fetch container status:", error)
      // Don't show error toast for status fetch
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchStatus()
    
    // Refresh status every 10 seconds
    const interval = setInterval(fetchStatus, 10000)
    return () => clearInterval(interval)
  }, [appId])

  const handleStop = async () => {
    try {
      setActionLoading("stop")
      await applicationsService.stopContainer(appId)
      toast.success("Container stopped successfully")
      await fetchStatus()
      onStatusChange?.()
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to stop container")
    } finally {
      setActionLoading(null)
    }
  }

  const handleStart = async () => {
    try {
      setActionLoading("start")
      await applicationsService.startContainer(appId)
      toast.success("Container started successfully")
      await fetchStatus()
      onStatusChange?.()
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to start container")
    } finally {
      setActionLoading(null)
    }
  }

  const handleRestart = async () => {
    try {
      setActionLoading("restart")
      await applicationsService.restartContainer(appId)
      toast.success("Container restarted successfully")
      await fetchStatus()
      onStatusChange?.()
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to restart container")
    } finally {
      setActionLoading(null)
    }
  }

  const getStatusBadge = () => {
    if (!status) return null

    switch (status.state) {
      case "running":
        return (
          <Badge className="bg-green-500 text-white flex items-center gap-1.5">
            <Activity className="h-3 w-3" />
            Running
          </Badge>
        )
      case "stopped":
        return (
          <Badge variant="secondary" className="flex items-center gap-1.5">
            <PowerOff className="h-3 w-3" />
            Stopped
          </Badge>
        )
      case "error":
        return (
          <Badge variant="destructive" className="flex items-center gap-1.5">
            <XCircle className="h-3 w-3" />
            Error
          </Badge>
        )
      default:
        return (
          <Badge variant="outline" className="flex items-center gap-1.5">
            {status.state}
          </Badge>
        )
    }
  }

  return (
    <Card>
      <CardHeader className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
        <CardTitle className="flex items-center gap-2">
          Container Controls
          {loading && <Loader2 className="h-4 w-4 animate-spin text-muted-foreground" />}
        </CardTitle>
        <Button
          variant="outline"
          size="sm"
          onClick={fetchStatus}
          disabled={loading}
        >
          <RefreshCw className={`h-4 w-4 ${loading ? "animate-spin" : ""}`} />
        </Button>
      </CardHeader>

      <CardContent className="space-y-4">
        {/* Status Information */}
        {status && (
          <div className="space-y-3">
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">Status</span>
              {getStatusBadge()}
            </div>

            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">Container</span>
              <span className="text-sm font-mono">{status.name}</span>
            </div>

            {status.state === "running" && (
              <>
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">Health</span>
                  {status.healthy ? (
                    <Badge variant="outline" className="flex items-center gap-1.5 text-green-600 border-green-600">
                      <CheckCircle2 className="h-3 w-3" />
                      Healthy
                    </Badge>
                  ) : (
                    <Badge variant="outline" className="flex items-center gap-1.5 text-yellow-600 border-yellow-600">
                      <Activity className="h-3 w-3" />
                      Unknown
                    </Badge>
                  )}
                </div>

                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground flex items-center gap-1.5">
                    <Clock className="h-3 w-3" />
                    Started At
                  </span>
                  <span className="text-sm">
                    {new Date(status.uptime).toLocaleString()}
                  </span>
                </div>
              </>
            )}
          </div>
        )}

        {!status && !loading && (
          <div className="text-center py-4 text-muted-foreground">
            <PowerOff className="h-8 w-8 mx-auto mb-2 opacity-50" />
            <p className="text-sm">No container found</p>
          </div>
        )}

        {/* Control Buttons */}
        <div className="flex gap-2 pt-2">
          <Button
            onClick={handleStart}
            disabled={status?.state === "running" || actionLoading !== null}
            className="flex-1"
            variant="default"
          >
            {actionLoading === "start" ? (
              <>
                <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                Starting...
              </>
            ) : (
              <>
                <Power className="h-4 w-4 mr-2" />
                Start
              </>
            )}
          </Button>

          <Button
            onClick={handleStop}
            disabled={status?.state !== "running" || actionLoading !== null}
            className="flex-1"
            variant="destructive"
          >
            {actionLoading === "stop" ? (
              <>
                <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                Stopping...
              </>
            ) : (
              <>
                <PowerOff className="h-4 w-4 mr-2" />
                Stop
              </>
            )}
          </Button>

          <Button
            onClick={handleRestart}
            disabled={status?.state !== "running" || actionLoading !== null}
            className="flex-1"
            variant="outline"
          >
            {actionLoading === "restart" ? (
              <>
                <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                Restarting...
              </>
            ) : (
              <>
                <RefreshCw className="h-4 w-4 mr-2" />
                Restart
              </>
            )}
          </Button>
        </div>

        {/* Info Message */}
        <p className="text-xs text-muted-foreground text-center pt-2">
          Container controls allow you to manage the running state of your application
        </p>
      </CardContent>
    </Card>
  )
}
