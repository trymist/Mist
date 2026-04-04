import { useEffect, useRef, useState } from "react"
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { useContainerLogs } from "@/hooks"
import { LogLine } from "@/components/logs/log-line"
import {
  Loader2,
  Wifi,
  WifiOff,
  Terminal,
  Trash2,
  Download,
  Play,
  Square,
  AlertCircle,
  CheckCircle2,
} from "lucide-react"
import { Alert, AlertDescription } from "@/components/ui/alert"

interface LiveLogsViewerProps {
  appId: number
  enabled?: boolean
}

export const LiveLogsViewer = ({ appId, enabled = true }: LiveLogsViewerProps) => {
  const [autoScroll, setAutoScroll] = useState(true)
  const logsEndRef = useRef<HTMLDivElement>(null)
  const logsContainerRef = useRef<HTMLDivElement>(null)

  const {
    logs,
    containerState,
    error,
    isConnected,
    isLoading,
    clearLogs,
  } = useContainerLogs({
    appId,
    enabled,
    onError: (error) => {
      console.error("Container logs error:", error)
    },
  })

  useEffect(() => {
    if (autoScroll && logsEndRef.current) {
      logsEndRef.current.scrollIntoView({ behavior: "smooth" })
    }
  }, [logs, autoScroll])

  useEffect(() => {
    const container = logsContainerRef.current
    if (!container) return

    const handleScroll = () => {
      const { scrollTop, scrollHeight, clientHeight } = container
      const isAtBottom = scrollHeight - scrollTop - clientHeight < 50
      setAutoScroll(isAtBottom)
    }

    container.addEventListener("scroll", handleScroll)
    return () => container.removeEventListener("scroll", handleScroll)
  }, [])

  const downloadLogs = () => {
    const logText = logs.map(log => log.line).join("\n")
    const blob = new Blob([logText], { type: "text/plain" })
    const url = URL.createObjectURL(blob)
    const a = document.createElement("a")
    a.href = url
    a.download = `container-logs-${appId}-${new Date().toISOString()}.txt`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
  }

  const getConnectionBadge = () => {
    if (isLoading) {
      return (
        <Badge variant="outline" className="flex items-center gap-1.5">
          <Loader2 className="h-3 w-3 animate-spin" />
          Connecting...
        </Badge>
      )
    }

    if (isConnected) {
      return (
        <Badge className="bg-green-500 text-white flex items-center gap-1.5">
          <Wifi className="h-3 w-3" />
          Live
        </Badge>
      )
    }

    return (
      <Badge variant="destructive" className="flex items-center gap-1.5">
        <WifiOff className="h-3 w-3" />
        Disconnected
      </Badge>
    )
  }

  return (
    <Card className="flex flex-col h-full">
      <CardHeader className="border-b border-border/50 bg-muted/30">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div className="flex items-center gap-3">
            <Terminal className="h-5 w-5 text-primary" />
            <CardTitle>Live Container Logs</CardTitle>
            {getConnectionBadge()}
          </div>

          <div className="flex items-center gap-2">
            {containerState && (
              <Badge variant="outline" className="hidden sm:flex">
                State: {containerState}
              </Badge>
            )}

            {logs.length > 0 && (
              <Button
                variant="outline"
                size="sm"
                onClick={downloadLogs}
                className="flex items-center gap-2"
              >
                <Download className="h-4 w-4" />
                Download
              </Button>
            )}

            {logs.length > 0 && (
              <Button
                variant="outline"
                size="sm"
                onClick={clearLogs}
                className="flex items-center gap-2"
              >
                <Trash2 className="h-4 w-4" />
                Clear
              </Button>
            )}
          </div>
        </div>
      </CardHeader>

      <CardContent className="flex-1 p-0 overflow-hidden flex flex-col">
        {/* Status Messages */}
        {error && (
          <Alert variant="destructive" className="m-4 mb-0">
            <AlertCircle className="h-4 w-4" />
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}

        {isLoading && (
          <div className="flex items-center justify-center p-8 text-muted-foreground">
            <Loader2 className="h-6 w-6 animate-spin mr-2" />
            Connecting to container...
          </div>
        )}

        {!isLoading && !error && logs.length === 0 && (
          <div className="flex flex-col items-center justify-center p-12 text-center">
            <Terminal className="h-12 w-12 text-muted-foreground/50 mb-4" />
            <p className="text-muted-foreground font-medium mb-1">No logs yet</p>
            <p className="text-sm text-muted-foreground">
              Waiting for container output...
            </p>
          </div>
        )}

        {logs.length > 0 && (
          <>
            <div
              ref={logsContainerRef}
              className="flex-1 overflow-y-auto bg-slate-950 text-slate-100 p-4 space-y-0.5"
              style={{ height: "calc(100vh - 400px)", minHeight: "400px" }}
            >
              {logs.map((log, index) => (
                <LogLine
                  key={index}
                  line={log.line}
                  index={index}
                  showLineNumbers={true}
                  streamType={log.stream}
                />
              ))}
              <div ref={logsEndRef} />
            </div>

            {!autoScroll && (
              <div className="border-t border-border/50 p-2 bg-muted/50 flex items-center justify-between">
                <div className="flex items-center gap-2 text-xs text-muted-foreground">
                  <AlertCircle className="h-3 w-3" />
                  <span>Auto-scroll paused (scroll to bottom to resume)</span>
                </div>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => {
                    logsEndRef.current?.scrollIntoView({ behavior: "smooth" })
                    setAutoScroll(true)
                  }}
                  className="text-xs h-7"
                >
                  Jump to bottom
                </Button>
              </div>
            )}

            <div className="border-t border-border/50 px-4 py-2 bg-muted/30 flex items-center justify-between text-xs text-muted-foreground">
              <div className="flex items-center gap-4">
                <span className="flex items-center gap-1.5">
                  <CheckCircle2 className="h-3 w-3" />
                  {logs.length} lines
                </span>
                {autoScroll && (
                  <span className="flex items-center gap-1.5 text-green-600">
                    <Play className="h-3 w-3" />
                    Auto-scrolling
                  </span>
                )}
                {!autoScroll && (
                  <span className="flex items-center gap-1.5 text-yellow-600">
                    <Square className="h-3 w-3" />
                    Paused
                  </span>
                )}
              </div>
              <div className="text-xs font-mono">
                App ID: {appId}
              </div>
            </div>
          </>
        )}
      </CardContent>
    </Card>
  )
}
