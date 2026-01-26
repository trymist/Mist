import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"

import { Badge } from "@/components/ui/badge"
import {
  CheckCircle,
  XCircle,
  Github,
  GitBranch,
  GitCommit,
  Rocket,
  Server,
  Activity,
  Clock,
  Loader2,
  ExternalLink, Database, Globe,
} from "lucide-react"
import { Button } from "@/components/ui/button"
import { toast } from "sonner"
import { useState, useEffect } from "react"
import type { App } from "@/types"
import { DeploymentMonitor } from "@/components/deployments"
import { deploymentsService, applicationsService } from "@/services"
import { useDomains } from "@/hooks"

interface Props {
  app: App
  latestCommit?: {
    sha: string
    html_url: string
    author?: string
    timestamp?: string
    message?: string
  } | null
}

const InfoItem = ({
  icon: Icon,
  label,
  children,
  className = ""
}: {
  icon: React.ComponentType<{ className?: string }>
  label: string
  children: React.ReactNode
  className?: string
}) => (
  <div className={`space-y-1.5 ${className}`}>
    <div className="flex items-center gap-1.5 text-xs font-medium text-muted-foreground">
      <Icon className="h-3.5 w-3.5" />
      <span>{label}</span>
    </div>
    <div className="pl-5">{children}</div>
  </div>
)

const SectionDivider = ({ title }: { title: string }) => (
  <div className="md:col-span-2 pt-3 pb-1 border-b border-border/30">
    <h3 className="text-xs font-semibold text-muted-foreground uppercase tracking-wide">
      {title}
    </h3>
  </div>
)

const getStatusConfig = (status: string) => {
  switch (status.toLowerCase()) {
    case "running":
      return {
        color: "text-green-600 dark:text-green-400",
        bgColor: "bg-green-500/10 border-green-500/20",
        icon: CheckCircle,
      }
    case "error":
    case "failed":
      return {
        color: "text-red-600 dark:text-red-400",
        bgColor: "bg-red-500/10 border-red-500/20",
        icon: XCircle,
      }
    case "deploying":
    case "building":
      return {
        color: "text-blue-600 dark:text-blue-400",
        bgColor: "bg-blue-500/10 border-blue-500/20",
        icon: Activity,
      }
    default:
      return {
        color: "text-gray-600 dark:text-gray-400",
        bgColor: "bg-gray-500/10 border-gray-500/20",
        icon: XCircle,
      }
  }
}

export const AppInfo = ({ app, latestCommit }: Props) => {
  const [deploying, setDeploying] = useState(false)
  const [logsOpen, setLogsOpen] = useState(false)
  const [deploymentId, setDeploymentId] = useState<number | null>(null)
  const [previewUrl, setPreviewUrl] = useState<string>("")

  // Fetch domains for web apps
  const { domains } = useDomains({
    appId: app.id,
    autoFetch: app.appType === 'web' && app.status === 'running'
  })

  const statusConfig = getStatusConfig(app.status)

  useEffect(() => {
    const fetchPreviewUrl = async () => {
      try {
        const data = await applicationsService.getPreviewUrl(app.id)
        setPreviewUrl(data.url)
      } catch (error) {
        console.error("Failed to fetch preview URL:", error)
      }
    }

    if (app.status === "running" && app.appType === 'web' && domains.length === 0) {
      fetchPreviewUrl()
    }
  }, [app.id, app.status, app.appType, domains.length])

  const handleDeploy = async () => {
    try {
      setDeploying(true)

      const deployment = await deploymentsService.create(app.id)

      toast.success("Deployment started!")

      setDeploymentId(deployment.id)
      setLogsOpen(true)

    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Deployment failed")
    } finally {
      setDeploying(false)
    }
  }

  return (
    <Card className="border-border/40 shadow-sm pt-0">
      <CardHeader className="border-b border-border/40 pt-5 flex flex-col sm:flex-row justify-between gap-3 bg-gradient-to-br from-muted/40 to-muted/20">
        <div>
          <CardTitle className="text-lg font-semibold">Application Overview</CardTitle>
          <p className="text-xs text-muted-foreground mt-0.5">View and manage your application details</p>
        </div>
        <Button
          onClick={handleDeploy}
          disabled={deploying}
          className="flex items-center gap-2 w-full sm:w-auto"
          size="sm"
        >
          {deploying ? (
            <>
              <Loader2 className="h-4 w-4 animate-spin" />
              Deploying...
            </>
          ) : (
            <>
              <Rocket className="h-4 w-4" />
              Deploy
            </>
          )}
        </Button>
      </CardHeader>

      <CardContent className="px-6 py-2">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-4">
          {/* Status & Access Section */}
          <div className="md:col-span-2 space-y-4 ">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-x-8 gap-y-4">
              {/* Status */}
              <InfoItem icon={Activity} label="Status">
                <div className="flex items-center gap-3 flex-wrap">
                  <Badge
                    variant="outline"
                    className={`${statusConfig.bgColor} ${statusConfig.color} font-medium border px-3 py-1`}
                  >
                    <span className="relative flex h-2 w-2 mr-2">
                      {app.status === "running" && (
                        <>
                          <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
                          <span className="relative inline-flex rounded-full h-2 w-2 bg-green-500"></span>
                        </>
                      )}
                      {app.status !== "running" && (
                        <span className="relative inline-flex rounded-full h-2 w-2 bg-current opacity-75"></span>
                      )}
                    </span>
                    {app.status.charAt(0).toUpperCase() + app.status.slice(1)}
                  </Badge>
                </div>
              </InfoItem>

              {/* Deployment Strategy */}
              <InfoItem icon={Server} label="Deployment Strategy">
                <Badge variant="secondary" className="font-mono text-xs px-3 py-1">
                  {app.deploymentStrategy || "Not specified"}
                </Badge>
              </InfoItem>
            </div>

            {/* Domains or Preview URL - only for web apps */}
            {app.appType === 'web' && app.status === 'running' && (
              <InfoItem icon={Globe} label="Access URLs">
                <div className="flex flex-col gap-2">
                  {domains.length > 0 ? (
                    domains.map((domain) => (
                      <a
                        key={domain.id}
                        href={`https://${domain.domain}`}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="inline-flex items-center gap-2 text-sm text-primary hover:underline font-medium w-fit group p-1.5 -ml-1.5 rounded-md hover:bg-muted/50 transition-colors"
                      >
                        <ExternalLink className="h-3.5 w-3.5 flex-shrink-0 group-hover:translate-x-0.5 group-hover:-translate-y-0.5 transition-transform" />
                        <span className="font-mono break-all">{domain.domain}</span>
                        {domain.sslStatus === 'active' && (
                          <Badge variant="outline" className="bg-green-500/10 border-green-500/20 text-green-600 dark:text-green-400 text-xs">
                            SSL
                          </Badge>
                        )}
                      </a>
                    ))
                  ) : previewUrl ? (
                    <a
                      href={previewUrl}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="inline-flex items-center gap-2 text-sm text-primary hover:underline font-medium w-fit group p-1.5 -ml-1.5 rounded-md hover:bg-muted/50 transition-colors"
                    >
                      <ExternalLink className="h-3.5 w-3.5 group-hover:translate-x-0.5 group-hover:-translate-y-0.5 transition-transform" />
                      <span className="font-mono break-all">{previewUrl.replace(/^https?:\/\//, '')}</span>
                    </a>
                  ) : (
                    <p className="text-muted-foreground text-sm">No domains configured</p>
                  )}
                </div>
              </InfoItem>
            )}
          </div>

          {app.appType === 'database' ? (
            <>
              <SectionDivider title="Database Configuration" />

              <InfoItem icon={Database} label="Service Template">
                <Badge variant="outline" className="font-mono text-xs px-3 py-1">
                  {app.templateName || "Unknown"}
                </Badge>
              </InfoItem>

              {app.port && app.port !== 0 && (
                <InfoItem icon={Server} label="Port">
                  <div className="font-mono text-sm font-semibold">
                    {app.port}
                  </div>
                </InfoItem>
              )}
            </>
          ) : app.appType === 'compose' ? (
            <>
              <SectionDivider title="Repository & Configuration" />

              {/* Git Repository */}
              <InfoItem icon={Github} label="Repository">
                {app.gitRepository ? (
                  <a
                    href={`https://github.com/${app.gitRepository}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="font-mono text-xs text-primary hover:underline flex items-center gap-2 group w-fit"
                  >
                    <span className="truncate">{app.gitRepository}</span>
                    <ExternalLink className="w-3 h-3 opacity-0 group-hover:opacity-100 transition-all group-hover:translate-x-0.5 group-hover:-translate-y-0.5" />
                  </a>
                ) : (
                  <p className="text-muted-foreground text-xs">Not connected</p>
                )}
              </InfoItem>

              {/* Branch */}
              <InfoItem icon={GitBranch} label="Branch">
                <Badge variant="outline" className="font-mono text-xs px-3 py-1">
                  {app.gitBranch || "Not specified"}
                </Badge>
              </InfoItem>

              {/* Latest Commit */}
              {latestCommit && (
                <>
                  <SectionDivider title="Latest Commit" />

                  <div className="md:col-span-2">
                    <InfoItem icon={GitCommit} label="Commit Details">
                      <div className="space-y-2 p-3 rounded-lg bg-muted/30 border border-border/50">
                        <div className="flex items-center gap-3 flex-wrap">
                          <a
                            href={latestCommit.html_url}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="font-mono text-xs text-primary hover:underline inline-flex items-center gap-1.5 group"
                          >
                            <span className="font-semibold">{latestCommit.sha.slice(0, 7)}</span>
                            <ExternalLink className="w-3 h-3 opacity-0 group-hover:opacity-100 transition-all group-hover:translate-x-0.5 group-hover:-translate-y-0.5" />
                          </a>
                          {latestCommit.author && (
                            <span className="text-xs text-muted-foreground">
                              by {latestCommit.author}
                            </span>
                          )}
                          {latestCommit.timestamp && (
                            <span className="flex items-center gap-1 text-xs text-muted-foreground">
                              <Clock className="w-3 h-3" />
                              {new Date(latestCommit.timestamp).toLocaleString()}
                            </span>
                          )}
                        </div>
                        {latestCommit.message && (
                          <p className="text-xs text-foreground/80 leading-relaxed">{latestCommit.message}</p>
                        )}
                      </div>
                    </InfoItem>
                  </div>
                </>
              )}
            </>
          ) : (
            <>
              {/* Git Configuration Section */}
              <SectionDivider title="Repository & Configuration" />

              {/* Git Repository */}
              <InfoItem icon={Github} label="Repository">
                {app.gitRepository ? (
                  <a
                    href={`https://github.com/${app.gitRepository}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="font-mono text-xs text-primary hover:underline flex items-center gap-2 group w-fit"
                  >
                    <span className="truncate">{app.gitRepository}</span>
                    <ExternalLink className="w-3 h-3 opacity-0 group-hover:opacity-100 transition-all group-hover:translate-x-0.5 group-hover:-translate-y-0.5" />
                  </a>
                ) : (
                  <p className="text-muted-foreground text-xs">Not connected</p>
                )}
              </InfoItem>

              {/* Branch */}
              <InfoItem icon={GitBranch} label="Branch">
                <Badge variant="outline" className="font-mono text-xs px-3 py-1">
                  {app.gitBranch || "Not specified"}
                </Badge>
              </InfoItem>

              {/* Port */}
              {app.port !== 0 && (
                <InfoItem icon={Server} label="Port">
                  <div className="font-mono text-sm font-semibold">
                    {app.port || <s className="text-red-500">"Not configured"</s>}
                  </div>
                </InfoItem>
              )}

              {/* Root Directory */}
              {app.rootDirectory && (
                <InfoItem icon={Server} label="Root Directory">
                  <div className="font-mono text-sm text-muted-foreground">
                    {app.rootDirectory}
                  </div>
                </InfoItem>
              )}

              {/* Build & Start Commands */}
              {(app.buildCommand || app.startCommand) && (
                <SectionDivider title="Commands" />
              )}

              {/* Build Command */}
              {app.buildCommand && (
                <InfoItem icon={Server} label="Build Command" className="md:col-span-2">
                  <div className="font-mono text-xs bg-muted/50 border border-border/50 rounded-md px-3 py-2 overflow-x-auto">
                    {app.buildCommand}
                  </div>
                </InfoItem>
              )}

              {/* Start Command */}
              {app.startCommand && (
                <InfoItem icon={Rocket} label="Start Command" className="md:col-span-2">
                  <div className="font-mono text-xs bg-muted/50 border border-border/50 rounded-md px-3 py-2 overflow-x-auto">
                    {app.startCommand}
                  </div>
                </InfoItem>
              )}

              {/* Latest Commit */}
              {latestCommit && (
                <>
                  <SectionDivider title="Latest Commit" />

                  <div className="md:col-span-2">
                    <InfoItem icon={GitCommit} label="Commit Details">
                      <div className="space-y-2 p-3 rounded-lg bg-muted/30 border border-border/50">
                        <div className="flex items-center gap-3 flex-wrap">
                          <a
                            href={latestCommit.html_url}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="font-mono text-xs text-primary hover:underline inline-flex items-center gap-1.5 group"
                          >
                            <span className="font-semibold">{latestCommit.sha.slice(0, 7)}</span>
                            <ExternalLink className="w-3 h-3 opacity-0 group-hover:opacity-100 transition-all group-hover:translate-x-0.5 group-hover:-translate-y-0.5" />
                          </a>
                          {latestCommit.author && (
                            <span className="text-xs text-muted-foreground">
                              by {latestCommit.author}
                            </span>
                          )}
                          {latestCommit.timestamp && (
                            <span className="flex items-center gap-1 text-xs text-muted-foreground">
                              <Clock className="w-3 h-3" />
                              {new Date(latestCommit.timestamp).toLocaleString()}
                            </span>
                          )}
                        </div>
                        {latestCommit.message && (
                          <p className="text-xs text-foreground/80 leading-relaxed">{latestCommit.message}</p>
                        )}
                      </div>
                    </InfoItem>
                  </div>
                </>
              )}
            </>
          )}

          {/* Healthcheck */}
          {app.healthcheckPath && (
            <>
              <SectionDivider title="Health Check" />
              <InfoItem icon={Activity} label="Configuration" className="md:col-span-2">
                <div className="flex items-center gap-3 flex-wrap">
                  <code className="font-mono text-xs bg-muted/50 border border-border/50 rounded-md px-3 py-1.5">
                    {app.healthcheckPath}
                  </code>
                  <Badge variant="secondary" className="text-xs">
                    Interval: {app.healthcheckInterval}s
                  </Badge>
                </div>
              </InfoItem>
            </>
          )}

          {/* Timestamps */}
          <SectionDivider title="Metadata" />

          <InfoItem icon={Clock} label="Timestamps" className="md:col-span-2">
            <div className="flex flex-col sm:flex-row sm:flex-wrap gap-3 sm:gap-6 text-xs">
              <div className="flex items-center gap-2">
                <span className="text-muted-foreground">Created:</span>
                <span className="font-medium text-foreground break-all">
                  {new Date(app.createdAt).toLocaleString()}
                </span>
              </div>
              <div className="flex items-center gap-2">
                <span className="text-muted-foreground">Updated:</span>
                <span className="font-medium text-foreground break-all">
                  {new Date(app.updatedAt).toLocaleString()}
                </span>
              </div>
            </div>
          </InfoItem>
        </div>
      </CardContent>

      {deploymentId !== null && (
        <DeploymentMonitor
          deploymentId={deploymentId}
          open={logsOpen}
          onClose={() => {
            setLogsOpen(false)
            setDeploymentId(null)
          }}
          onComplete={() => {
            toast.success("Deployment completed successfully!")
          }}
        />
      )}
    </Card>
  )
}
