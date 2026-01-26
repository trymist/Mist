import { useState, useEffect } from "react"
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { applicationsService } from "@/services"
import { toast } from "sonner"
import {
    Activity,
    RefreshCw,
    Server,
    Loader2,
    Power,
    Square,
    RotateCw,
    Box,
} from "lucide-react"

interface ComposeService {
    name: string
    status: string
    state: string
}

interface ComposeStatusData {
    name: string
    status: string
    state: string
    uptime: string
    services: ComposeService[]
    error?: string
}

interface ComposeStatusProps {
    appId: number
    onStatusChange?: () => void
}

export const ComposeStatus = ({ appId, onStatusChange }: ComposeStatusProps) => {
    const [status, setStatus] = useState<ComposeStatusData | null>(null)
    const [loading, setLoading] = useState(false)
    const [actionLoading, setActionLoading] = useState<string | null>(null)
    const [lastUpdated, setLastUpdated] = useState<Date>(new Date())

    const fetchStatus = async () => {
        try {
            setLoading(true)
            // We rely on the same endpoint, but it returns different structure for compose
            const data = await applicationsService.getContainerStatus(appId) as ComposeStatusData
            setStatus(data)
            setLastUpdated(new Date())
        } catch (error) {
            console.error("Failed to fetch compose status:", error)
        } finally {
            setLoading(false)
        }
    }

    useEffect(() => {
        fetchStatus()
        const interval = setInterval(fetchStatus, 15000)
        return () => clearInterval(interval)
    }, [appId])

    const handleStart = async () => {
        try {
            setActionLoading("start")
            await applicationsService.startContainer(appId)
            toast.success("Stack start triggered successfully")
            await fetchStatus()
            onStatusChange?.()
        } catch (error) {
            toast.error(error instanceof Error ? error.message : "Failed to start stack")
        } finally {
            setActionLoading(null)
        }
    }

    const handleStop = async () => {
        try {
            setActionLoading("stop")
            await applicationsService.stopContainer(appId)
            toast.success("Stack stop triggered successfully")
            await fetchStatus()
            onStatusChange?.()
        } catch (error) {
            toast.error(error instanceof Error ? error.message : "Failed to stop stack")
        } finally {
            setActionLoading(null)
        }
    }

    const handleRestart = async () => {
        try {
            setActionLoading("restart")
            await applicationsService.restartContainer(appId)
            toast.success("Stack restart triggered successfully")
            await fetchStatus()
            onStatusChange?.()
        } catch (error) {
            toast.error(error instanceof Error ? error.message : "Failed to restart stack")
        } finally {
            setActionLoading(null)
        }
    }

    const getOverallStateBadge = () => {
        if (!status) return null
        if (status.error) return <Badge variant="destructive">Error</Badge>

        switch (status.state) {
            case "running":
                return <Badge className="bg-green-500 text-white">Running</Badge>
            case "stopped":
                return <Badge variant="secondary">Stopped</Badge>
            case "partial":
                return <Badge className="bg-yellow-500 text-white">Partial</Badge>
            default:
                return <Badge variant="outline">{status.state}</Badge>
        }
    }

    return (
        <Card className="border-border/50">
            <CardHeader className="border-b border-border/50 bg-muted/30">
                <div className="flex items-center justify-between">
                    <CardTitle className="text-lg font-semibold flex items-center gap-2">
                        <Activity className="h-5 w-5 text-primary" />
                        Stack Status
                    </CardTitle>
                    <Button
                        variant="outline"
                        size="sm"
                        onClick={fetchStatus}
                        disabled={loading}
                    >
                        <RefreshCw className={`h-4 w-4 ${loading ? "animate-spin" : ""}`} />
                    </Button>
                </div>
            </CardHeader>

            <CardContent className="p-6">
                {!status && loading && (
                    <div className="flex items-center justify-center py-8 text-muted-foreground">
                        <Loader2 className="h-6 w-6 animate-spin mr-2" />
                        Checking status...
                    </div>
                )}

                {status && (
                    <div className="space-y-6">
                        {/* Controls */}
                        <div className="flex gap-2 pb-4 border-b border-border/50">
                            <Button
                                onClick={handleStart}
                                disabled={status.state === "running" || actionLoading !== null || loading}
                                size="sm"
                                className="flex-1"
                            >
                                {actionLoading === "start" ? <Loader2 className="h-4 w-4 mr-2 animate-spin" /> : <Power className="h-4 w-4 mr-2" />}
                                Start
                            </Button>
                            <Button
                                onClick={handleStop}
                                disabled={status.state === "stopped" || actionLoading !== null || loading}
                                size="sm"
                                variant="destructive"
                                className="flex-1"
                            >
                                {actionLoading === "stop" ? <Loader2 className="h-4 w-4 mr-2 animate-spin" /> : <Square className="h-4 w-4 mr-2" />}
                                Stop
                            </Button>
                            <Button
                                onClick={handleRestart}
                                disabled={actionLoading !== null || loading}
                                size="sm"
                                variant="outline"
                                className="flex-1"
                            >
                                {actionLoading === "restart" ? <Loader2 className="h-4 w-4 mr-2 animate-spin" /> : <RotateCw className="h-4 w-4 mr-2" />}
                                Restart
                            </Button>
                        </div>

                        {/* Overall Status */}
                        <div className="flex items-center justify-between py-2">
                            <div className="flex items-center gap-2 text-sm text-muted-foreground">
                                <Server className="h-4 w-4" />
                                <span>Overall Status</span>
                            </div>
                            <div className="flex items-center gap-2">
                                {getOverallStateBadge()}
                                <span className="text-xs text-muted-foreground">({status.status})</span>
                            </div>
                        </div>

                        {/* Services List */}
                        {status.services && status.services.length > 0 && (
                            <div className="space-y-2">
                                <h4 className="text-sm font-medium mb-2 flex items-center gap-2">
                                    <Box className="h-4 w-4" /> Services
                                </h4>
                                <div className="space-y-2">
                                    {status.services.map((svc, idx) => (
                                        <div key={idx} className="flex items-center justify-between p-3 rounded-md bg-muted/40 border border-border/40">
                                            <span className="text-sm font-mono font-medium">{svc.name}</span>
                                            <Badge variant={svc.state === 'running' ? 'default' : 'secondary'} className={svc.state === 'running' ? 'bg-green-500/10 text-green-600 hover:bg-green-500/20 border-green-500/20' : ''}>
                                                {svc.state}
                                            </Badge>
                                        </div>
                                    ))}
                                </div>
                            </div>
                        )}

                        {status.error && (
                            <div className="p-3 bg-destructive/10 text-destructive text-sm rounded-md">
                                Error: {status.error}
                            </div>
                        )}

                        <div className="pt-2 text-xs text-muted-foreground text-center">
                            Last updated: {lastUpdated.toLocaleTimeString()}
                        </div>
                    </div>
                )}
            </CardContent>
        </Card>
    )
}
