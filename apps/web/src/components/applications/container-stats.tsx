import { useState, useEffect, useRef } from "react"
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Alert, AlertDescription } from "@/components/ui/alert"
import {
  Activity,
  Cpu,
  HardDrive,
  Network,
  Loader2,
  AlertCircle,
  TrendingUp,
} from "lucide-react"
import { ResponsiveContainer, AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip } from 'recharts'

interface ContainerStatsProps {
  appId: number
}

interface StatsData {
  cpu_percent: number
  memory_used_mb: number
  memory_limit_mb: number
  memory_percent: number
  network_rx_mb: number
  network_tx_mb: number
  block_read_mb: number
  block_write_mb: number
  pids: number
  timestamp?: number
}

interface StatsEvent {
  type: string
  timestamp: string
  data: StatsData | { message: string } | { container: string; state: string; status: string }
}

const formatPercentage = (value: number): string => {
  return `${value.toFixed(1)}%`
}

const formatMemory = (mb: number): string => {
  if (mb < 1024) {
    return `${mb.toFixed(2)} MB`
  }
  return `${(mb / 1024).toFixed(2)} GB`
}

const getUsageColor = (percentage: number): string => {
  if (percentage >= 80) return 'text-red-500'
  if (percentage >= 60) return 'text-yellow-500'
  return 'text-green-500'
}

export const ContainerStats = ({ appId }: ContainerStatsProps) => {
  const [stats, setStats] = useState<StatsData | null>(null)
  const [statsHistory, setStatsHistory] = useState<StatsData[]>([])
  const [connected, setConnected] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [containerState, setContainerState] = useState<string | null>(null)
  const wsRef = useRef<WebSocket | null>(null)

  useEffect(() => {
    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:"
    const wsUrl = `${protocol}//${window.location.host}/api/ws/container/stats?appId=${appId}`

    const connectWebSocket = () => {
      try {
        const ws = new WebSocket(wsUrl)
        wsRef.current = ws

        ws.onopen = () => {
          console.log("Container stats WebSocket connected")
          setConnected(true)
          setError(null)
        }

        ws.onmessage = (event) => {
          try {
            const message: StatsEvent = JSON.parse(event.data)

            switch (message.type) {
              case "status":
                const statusData = message.data as { container: string; state: string; status: string }
                setContainerState(statusData.state)
                break

              case "stats":
                const statsData = message.data as StatsData
                const timestamp = Date.now() / 1000
                const statsWithTimestamp = { ...statsData, timestamp }
                setStats(statsWithTimestamp)
                setStatsHistory(prev => {
                  const updated = [...prev, statsWithTimestamp]
                  return updated.slice(-50) // Keep last 50 data points
                })
                break

              case "error":
                const errorData = message.data as { message: string }
                setError(errorData.message)
                break

              case "end":
                const endData = message.data as { message: string }
                setError(endData.message)
                ws.close()
                break
            }
          } catch (err) {
            console.error("Error parsing WebSocket message:", err)
          }
        }

        ws.onerror = (event) => {
          console.error("WebSocket error:", event)
          setError("Connection error occurred")
        }

        ws.onclose = () => {
          console.log("Container stats WebSocket disconnected")
          setConnected(false)
          wsRef.current = null

          // Attempt to reconnect after 5 seconds
          if (!error) {
            setTimeout(() => {
              if (!wsRef.current || wsRef.current.readyState === WebSocket.CLOSED) {
                connectWebSocket()
              }
            }, 5000)
          }
        }
      } catch (err) {
        console.error("Error creating WebSocket:", err)
        setError("Failed to connect to stats stream")
      }
    }

    connectWebSocket()

    return () => {
      if (wsRef.current) {
        wsRef.current.close()
        wsRef.current = null
      }
    }
  }, [appId])

  const customTooltip = ({ active, payload, label }: {
    active?: boolean
    payload?: Array<{ value: number; color: string; name: string }>
    label?: string | number
  }) => {
    if (active && payload && payload.length) {
      const value = payload[0].value
      const timestamp = typeof label === 'number' ? label : parseInt(String(label), 10)

      return (
        <div className="bg-popover p-3 border border-border rounded-md">
          <p className="text-foreground text-sm">
            {new Date(timestamp * 1000).toLocaleTimeString()}
          </p>
          <p style={{ color: payload[0].color }} className="text-sm font-semibold">
            {payload[0].name}: {formatPercentage(value)}
          </p>
        </div>
      )
    }
    return null
  }

  return (
    <div className="space-y-6">
      {/* Connection Status */}
      <div className="flex items-center gap-2">
        {connected ? (
          <Badge className="bg-green-500 text-white flex items-center gap-1.5">
            <span className="relative flex h-2 w-2">
              <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-300 opacity-75"></span>
              <span className="relative inline-flex rounded-full h-2 w-2 bg-green-100"></span>
            </span>
            Live Stats
          </Badge>
        ) : (
          <Badge variant="secondary" className="flex items-center gap-1.5">
            <Loader2 className="h-3 w-3 animate-spin" />
            Connecting...
          </Badge>
        )}
        {containerState && (
          <Badge variant="outline">
            {containerState}
          </Badge>
        )}
      </div>

      {/* Error Alert */}
      {error && (
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}
      {!error && stats && (
        <div className="grid gap-4 grid-cols-1 sm:grid-cols-2 xl:grid-cols-4">
          {/* CPU Usage */}
          <Card className="p-4">
            <p className="text-sm text-muted-foreground flex items-center gap-2">
              <Cpu className="h-4 w-4" />
              CPU Usage
            </p>
            <p className={`text-lg font-semibold ${getUsageColor(stats.cpu_percent)}`}>
              {formatPercentage(stats.cpu_percent)}
            </p>
          </Card>

          {/* Memory Usage */}
          <Card className="p-4">
            <p className="text-sm text-muted-foreground flex items-center gap-2">
              <Activity className="h-4 w-4" />
              Memory Usage
            </p>
            <p className={`text-lg font-semibold ${getUsageColor(stats.memory_percent)}`}>
              {formatPercentage(stats.memory_percent)}
            </p>
            <p className="text-xs text-muted-foreground mt-1">
              {formatMemory(stats.memory_used_mb)} / {formatMemory(stats.memory_limit_mb)}
            </p>
          </Card>

          {/* Network RX */}
          <Card className="p-4">
            <p className="text-sm text-muted-foreground flex items-center gap-2">
              <Network className="h-4 w-4" />
              Network RX
            </p>
            <p className="text-lg font-semibold">
              {formatMemory(stats.network_rx_mb)}
            </p>
          </Card>

          {/* Network TX */}
          <Card className="p-4">
            <p className="text-sm text-muted-foreground flex items-center gap-2">
              <TrendingUp className="h-4 w-4" />
              Network TX
            </p>
            <p className="text-lg font-semibold">
              {formatMemory(stats.network_tx_mb)}
            </p>
          </Card>

          {/* Block Read */}
          <Card className="p-4">
            <p className="text-sm text-muted-foreground flex items-center gap-2">
              <HardDrive className="h-4 w-4" />
              Block Read
            </p>
            <p className="text-lg font-semibold">
              {formatMemory(stats.block_read_mb)}
            </p>
          </Card>

          {/* Block Write */}
          <Card className="p-4">
            <p className="text-sm text-muted-foreground flex items-center gap-2">
              <HardDrive className="h-4 w-4" />
              Block Write
            </p>
            <p className="text-lg font-semibold">
              {formatMemory(stats.block_write_mb)}
            </p>
          </Card>

          {/* Processes */}
          {stats.pids > 0 && (
            <Card className="p-4">
              <p className="text-sm text-muted-foreground">Processes</p>
              <p className="text-lg font-semibold">
                {stats.pids}
              </p>
            </Card>
          )}
        </div>
      )}
      {/* Charts */}
      {!error && statsHistory.length > 0 && (
        <div className="grid gap-6 grid-cols-1 lg:grid-cols-2">
          {/* CPU Chart */}
          <Card>
            <CardHeader>
              <CardTitle>CPU Usage</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="h-[300px]">
                <ResponsiveContainer width="100%" height="100%">
                  <AreaChart data={statsHistory}>
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
                      tickFormatter={formatPercentage}
                    />
                    <Tooltip content={(props) => customTooltip({ ...props, payload: props.payload?.map(p => ({ ...p, name: 'CPU' })) })} />
                    <Area
                      type="monotone"
                      dataKey="cpu_percent"
                      stroke="#8B5CF6"
                      fill="#8B5CF6"
                      fillOpacity={0.3}
                    />
                  </AreaChart>
                </ResponsiveContainer>
              </div>
            </CardContent>
          </Card>

          {/* Memory Chart */}
          <Card>
            <CardHeader>
              <CardTitle>Memory Usage</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="h-[300px]">
                <ResponsiveContainer width="100%" height="100%">
                  <AreaChart data={statsHistory}>
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
                      tickFormatter={formatPercentage}
                    />
                    <Tooltip content={(props) => customTooltip({ ...props, payload: props.payload?.map(p => ({ ...p, name: 'Memory' })) })} />
                    <Area
                      type="monotone"
                      dataKey="memory_percent"
                      stroke="#A371F7"
                      fill="#A371F7"
                      fillOpacity={0.3}
                    />
                  </AreaChart>
                </ResponsiveContainer>
              </div>
            </CardContent>
          </Card>
        </div>
      )}

      {/* Metric Cards */}

      {/* Loading State */}
      {!stats && !error && connected && (
        <div className="flex items-center justify-center py-12 text-muted-foreground">
          <div className="text-center space-y-3">
            <Loader2 className="h-8 w-8 animate-spin mx-auto" />
            <p className="text-sm">Waiting for stats data...</p>
          </div>
        </div>
      )}
    </div>
  )
}
