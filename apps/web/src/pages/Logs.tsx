import { useEffect, useRef, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Terminal, Trash2, Download, Pause, Play, AlertCircle, RefreshCw } from "lucide-react";
import { LogLine } from "@/components/logs/log-line";
import { toast } from "sonner";

interface LogEvent {
  type: string;
  timestamp: string;
  data: {
    line?: string;
    message?: string;
  };
}

export const LogsPage = () => {
  const [logs, setLogs] = useState<string[]>([]);
  const [connected, setConnected] = useState(false);
  const [paused, setPaused] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const wsRef = useRef<WebSocket | null>(null);
  const logsEndRef = useRef<HTMLDivElement>(null);
  const logsContainerRef = useRef<HTMLDivElement>(null);
  const pausedLogsRef = useRef<string[]>([]);

  const scrollToBottom = () => {
    if (!paused && logsEndRef.current) {
      logsEndRef.current.scrollIntoView({ behavior: "smooth" });
    }
  };

  useEffect(() => {
    scrollToBottom();
  }, [logs, paused]);

  const connectWebSocket = () => {
    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
    const wsUrl = `${protocol}//${window.location.host}/api/ws/system/logs`;

    const ws = new WebSocket(wsUrl);
    wsRef.current = ws;

    ws.onopen = () => {
      setConnected(true);
      setError(null);
      console.log("Connected to system logs");
    };

    ws.onmessage = (event) => {
      try {
        const logEvent: LogEvent = JSON.parse(event.data);

        if (logEvent.type === "log" && logEvent.data.line) {
          const line = logEvent.data.line;
          if (paused) {
            pausedLogsRef.current.push(line);
          } else {
            setLogs((prev) => [...prev, line]);
          }
        } else if (logEvent.type === "error") {
          setError(logEvent.data.message || "Unknown error");
          toast.error(logEvent.data.message || "Error streaming logs");
        } else if (logEvent.type === "connected") {
          toast.success("Connected to system logs");
        } else if (logEvent.type === "end") {
          toast.info("Log stream ended");
        }
      } catch (err) {
        console.error("Failed to parse log event:", err);
      }
    };

    ws.onerror = (error) => {
      console.error("WebSocket error:", error);
      setError("WebSocket connection error");
      toast.error("Failed to connect to system logs");
    };

    ws.onclose = () => {
      setConnected(false);
      console.log("Disconnected from system logs");
    };
  };

  const disconnectWebSocket = () => {
    if (wsRef.current) {
      wsRef.current.close();
      wsRef.current = null;
    }
  };

  useEffect(() => {
    connectWebSocket();

    return () => {
      disconnectWebSocket();
    };
  }, []);

  const handleClearLogs = () => {
    setLogs([]);
    pausedLogsRef.current = [];
    toast.success("Logs cleared");
  };

  const handleDownloadLogs = () => {
    const allLogs = [...logs, ...pausedLogsRef.current];
    const logText = allLogs.join("\n");
    const blob = new Blob([logText], { type: "text/plain" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `mist-system-logs-${new Date().toISOString()}.txt`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
    toast.success("Logs downloaded");
  };

  const handleTogglePause = () => {
    if (paused) {
      // Resume: add all paused logs to the main logs
      setLogs((prev) => [...prev, ...pausedLogsRef.current]);
      pausedLogsRef.current = [];
      setPaused(false);
      toast.info("Logs resumed");
    } else {
      setPaused(true);
      toast.info("Logs paused");
    }
  };

  const handleReconnect = () => {
    disconnectWebSocket();
    setLogs([]);
    pausedLogsRef.current = [];
    setPaused(false);
    setError(null);
    connectWebSocket();
    toast.info("Reconnecting...");
  };

  return (
    <div className="min-h-screen bg-background">
      {/* Page Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between py-6 border-b border-border gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-foreground">
            System Logs
          </h1>
          <p className="text-muted-foreground mt-1">
            Real-time logs from the Mist systemd service
          </p>
        </div>
        <div className="flex items-center gap-2">
          <div className={`w-3 h-3 rounded-full ${connected ? "bg-green-500" : "bg-red-500"}`} />
          <span className="text-sm text-muted-foreground">
            {connected ? "Connected" : "Disconnected"}
          </span>
        </div>
      </div>

      {/* Content */}
      <div className="py-6 space-y-6">
        {/* Error Alert */}
        {error && (
          <Alert variant="destructive">
            <AlertCircle className="h-4 w-4" />
            <AlertDescription className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
              <span>{error}</span>
              <Button
                variant="outline"
                size="sm"
                onClick={handleReconnect}
                className="w-full sm:w-auto"
              >
                <RefreshCw className="h-4 w-4 mr-2" />
                Reconnect
              </Button>
            </AlertDescription>
          </Alert>
        )}

        {/* Logs Card */}
        <Card>
          <CardHeader className="border-b border-border">
            <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4">
              <div>
                <CardTitle className="flex items-center gap-2">
                  <Terminal className="h-5 w-5" />
                  Live Logs
                </CardTitle>
                <CardDescription>
                  Streaming logs from journalctl (last 100 lines + live updates)
                </CardDescription>
              </div>
              <div className="flex flex-wrap items-center gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={handleTogglePause}
                  className="flex items-center gap-2 flex-1 sm:flex-initial"
                >
                  {paused ? (
                    <>
                      <Play className="h-4 w-4" />
                      Resume
                    </>
                  ) : (
                    <>
                      <Pause className="h-4 w-4" />
                      Pause
                    </>
                  )}
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={handleDownloadLogs}
                  disabled={logs.length === 0}
                  className="flex items-center gap-2 flex-1 sm:flex-initial"
                >
                  <Download className="h-4 w-4" />
                  <span className="hidden sm:inline">Download</span>
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={handleClearLogs}
                  disabled={logs.length === 0}
                  className="flex items-center gap-2 flex-1 sm:flex-initial"
                >
                  <Trash2 className="h-4 w-4" />
                  <span className="hidden sm:inline">Clear</span>
                </Button>
              </div>
            </div>
          </CardHeader>
          <CardContent className="p-0">
            <div
              ref={logsContainerRef}
              className="h-[calc(100vh-320px)] sm:h-[calc(100vh-300px)] overflow-y-auto overflow-x-auto bg-slate-950 p-4"
            >
              {logs.length === 0 && !connected && (
                <div className="flex items-center justify-center h-full text-muted-foreground">
                  <p>Connecting to system logs...</p>
                </div>
              )}
              {logs.length === 0 && connected && (
                <div className="flex items-center justify-center h-full text-muted-foreground">
                  <p>Waiting for logs...</p>
                </div>
              )}
              <div className="space-y-0.5">
                {logs.map((log, index) => (
                  <LogLine
                    key={index}
                    line={log}
                    index={index}
                    showLineNumbers={false}
                  />
                ))}
              </div>
              {paused && pausedLogsRef.current.length > 0 && (
                <div className="mt-4 p-2 bg-yellow-500/10 border border-yellow-500/20 rounded text-yellow-500">
                  {pausedLogsRef.current.length} new log entries (paused)
                </div>
              )}
              <div ref={logsEndRef} />
            </div>
          </CardContent>
        </Card>

        {/* Stats */}
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 text-sm text-muted-foreground">
          <div className="flex items-center gap-2">
            <Terminal className="h-4 w-4" />
            <span>Total log entries: {logs.length + pausedLogsRef.current.length}</span>
          </div>
          {paused && (
            <Badge variant="secondary" className="flex items-center gap-1.5 w-fit">
              <Pause className="h-3 w-3" />
              Paused
            </Badge>
          )}
        </div>
      </div>
    </div>
  );
};
