import { useEffect, useState } from "react"
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card"
import { Label } from "@/components/ui/label"
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs"
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from "@/components/ui/select"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { toast } from "sonner"
import type { App } from "@/types/app"
import { Github, GitBranch } from "lucide-react"
import { Skeleton } from "@/components/ui/skeleton"

interface GitProviderTabProps {
  app: App
}

export const GitProviderTab = ({ app }: GitProviderTabProps) => {
  const [provider, setProvider] = useState("github")

  const [, setGithubApp] = useState<{ name: string; slug: string } | null>(null)
  const [isInstalled, setIsInstalled] = useState<boolean>(false)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const [gitProviderId, setGitProviderId] = useState<number | null>(app.gitProviderId || null)

  const [repos, setRepos] = useState<Array<{ id: number; full_name: string; clone_url?: string }>>([])
  const [branches, setBranches] = useState<Array<{ name: string }>>([])
  const [selectedRepo, setSelectedRepo] = useState(app.gitRepository || "")
  const [selectedRepoCloneUrl, setSelectedRepoCloneUrl] = useState(app.gitCloneUrl || "")
  const [selectedBranch, setSelectedBranch] = useState(app.gitBranch || "")
  const [isRepoLoading, setIsRepoLoading] = useState(true)
  const [isBranchLoading, setIsBranchLoading] = useState(false)

  // Public Git state
  const [publicGitUrl, setPublicGitUrl] = useState(app.gitCloneUrl || "")
  const [publicGitBranch, setPublicGitBranch] = useState(app.gitBranch || "main")
  const [isSavingPublicGit, setIsSavingPublicGit] = useState(false)

  // this is for github app fetching
  // FIX: name of this function should be changed
  const fetchApp = async () => {
    try {
      setIsLoading(true)
      setError(null)

      const response = await fetch("/api/github/app", { credentials: "include" })
      const data = await response.json()

      if (data.success) {
        setGithubApp(data.data.app)
        setIsInstalled(data.data.isInstalled)
      } else {
        setError(data.error || "Failed to load GitHub App details")
      }
    } catch (error) {
      console.error('Failed to load GitHub App details:', error);
      setError("Failed to load GitHub App details")
    } finally {
      setIsLoading(false)
    }
  }

  // fetch user's git providers to get the git_provider_id
  const fetchGitProviders = async () => {
    try {
      const response = await fetch("/api/users/git-providers", { credentials: "include" })
      const data = await response.json()

      if (data.success && data.data && data.data.length > 0) {
        // find the GitHub provider
        // it is only the time being when there's only one git provider, and should be changed
        const githubProvider = data.data.find((p: any) => p.provider === "github")
        if (githubProvider) {
          setGitProviderId(githubProvider.id)
        }
      }
    } catch (error) {
      console.error('Failed to load git providers:', error)
    }
  }

  const fetchRepos = async () => {
    try {
      setIsRepoLoading(true)
      const res = await fetch("/api/github/repositories", { credentials: "include" })
      const data = await res.json()
      setRepos(data || [])
    } catch {
      toast.error("Failed to load repositories")
    } finally {
      setIsRepoLoading(false)
    }
  }

  const fetchBranchList = async (repoFullName: string) => {
    try {
      setIsBranchLoading(true)
      const res = await fetch(`/api/github/branches`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ repo: repoFullName }),
        credentials: "include",
      })
      const data = await res.json()

      if (data.success) setBranches(data.data)
    } catch {
      toast.error("Failed to load branches")
    } finally {
      setIsBranchLoading(false)
    }
  }

  useEffect(() => {
    fetchApp()
    fetchRepos()
    fetchGitProviders()
  }, [])

  useEffect(() => {
    if (selectedRepo) fetchBranchList(selectedRepo)
  }, [selectedRepo])

  useEffect(() => {
    if (app.gitRepository) {
      setRepos([{ id: 0, full_name: app.gitRepository, clone_url: app.gitCloneUrl || undefined }])
    }
    if (app.gitBranch) {
      setBranches([{ name: app.gitBranch }])
    }
  }, [app])

  const saveGitConfig = async () => {
    try {
      const res = await fetch("/api/apps/update", {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({
          appId: app.id,
          gitProviderId: gitProviderId,
          gitRepository: selectedRepo,
          gitBranch: selectedBranch,
          gitCloneUrl: selectedRepoCloneUrl || `https://github.com/${selectedRepo}.git`,
        }),
      })

      const data = await res.json()
      if (!data.success) throw new Error(data.error)

      toast.success("Git provider configuration saved")
    } catch {
      toast.error("Failed to save configuration")
    }
  }

  const savePublicGitConfig = async () => {
    if (!publicGitUrl.trim()) {
      toast.error("Please enter a Git URL")
      return
    }

    try {
      setIsSavingPublicGit(true)
      const res = await fetch("/api/apps/update", {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({
          appId: app.id,
          gitProviderId: null,
          gitRepository: null,
          gitBranch: publicGitBranch || "main",
          gitCloneUrl: publicGitUrl.trim(),
        }),
      })

      const data = await res.json()
      if (!data.success) throw new Error(data.error)

      toast.success("Public Git configuration saved")
    } catch {
      toast.error("Failed to save configuration")
    } finally {
      setIsSavingPublicGit(false)
    }
  }

  return (
    <Tabs defaultValue="github" value={provider} onValueChange={setProvider} className="w-full space-y-8">

      {/* ✅ PROVIDER LIST */}
      <div className="w-full overflow-x-auto pb-1">
        <TabsList className="inline-flex w-full min-w-fit">
          <TabsTrigger value="github" className="flex items-center gap-2">
            <Github className="h-4 w-4" />
            GitHub
          </TabsTrigger>

          <TabsTrigger value="public-git" className="flex items-center gap-2">
            <GitBranch className="h-4 w-4" />
            Public Git
          </TabsTrigger>

          {/* <TabsTrigger value="gitlab" disabled className="flex items-center gap-2 opacity-70"> */}
          {/*   <Gitlab className="h-4 w-4" /> */}
          {/*   GitLab */}
          {/* </TabsTrigger> */}
          {/**/}
          {/* <TabsTrigger value="gitea" disabled className="flex items-center gap-2 opacity-70"> */}
          {/*   <SiGitea className="h-4 w-4" /> */}
          {/*   Gitea */}
          {/* </TabsTrigger> */}
          {/**/}
          {/* <TabsTrigger value="bitbucket" disabled className="flex items-center gap-2 opacity-70"> */}
          {/*   <SiBitbucket className="h-4 w-4" /> */}
          {/*   Bitbucket */}
          {/* </TabsTrigger> */}
        </TabsList>
      </div>

      {/* ✅ GITHUB TAB CONTENT */}
      <TabsContent value="github">
        {/* ✅ Loading state */}
        {isLoading && (
          <Card>
            <CardHeader>
              <CardTitle>Loading…</CardTitle>
            </CardHeader>
          </Card>
        )}

        {/* ✅ Error state */}
        {!isLoading && error && (
          <Card className="border-red-500">
            <CardHeader>
              <CardTitle className="text-red-500">Error Loading GitHub App</CardTitle>
              <CardDescription>{error}</CardDescription>
            </CardHeader>
          </Card>
        )}

        {/* ✅ NOT INSTALLED → Show connection card */}
        {!isLoading && !isInstalled && (
          <Card>
            <CardHeader>
              <CardTitle>GitHub App Not Connected</CardTitle>
              <CardDescription>
                You need to connect your GitHub App to enable repository syncing.
              </CardDescription>
            </CardHeader>

            <CardContent>
              <Button asChild>
                <a href="/git">Connect GitHub App</a>
              </Button>
            </CardContent>
          </Card>
        )}

        {!isLoading && isInstalled && (
          <Card>
            <CardHeader>
              <CardTitle>GitHub Repository</CardTitle>
              <CardDescription>
                Select the repository and branch to link with your app.
              </CardDescription>
            </CardHeader>

            <CardContent className="space-y-6">

              <div className="flex flex-col md:flex-row gap-6">

                {/* ✅ Repo select or skeleton */}
                <div className="flex-1">
                  <Label className="text-muted-foreground">Repository</Label>

                  {isRepoLoading ? (
                    <Skeleton className="w-full h-10 mt-2" />
                  ) : (
                    <Select
                      value={selectedRepo}
                      onValueChange={(value) => {
                        setSelectedRepo(value)
                        // Find the repo object to get the clone_url
                        const repo = repos.find(r => r.full_name === value)
                        if (repo?.clone_url) {
                          setSelectedRepoCloneUrl(repo.clone_url)
                        } else {
                          // Fallback: construct GitHub URL
                          setSelectedRepoCloneUrl(`https://github.com/${value}.git`)
                        }
                      }}
                    >
                      <SelectTrigger className="mt-2 w-full">
                        <SelectValue placeholder="Select a repository" />
                      </SelectTrigger>

                      <SelectContent>
                        {repos.map((repo) => (
                          <SelectItem key={repo.id} value={repo.full_name}>
                            {repo.full_name}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  )}
                </div>

                {/* ✅ Branch select or skeleton */}
                <div className="flex-1">
                  <Label className="text-muted-foreground">Branch</Label>

                  {isBranchLoading ? (
                    <Skeleton className="w-full h-10 mt-2" />
                  ) : (
                    <Select value={selectedBranch} onValueChange={setSelectedBranch}>
                      <SelectTrigger className="mt-2 w-full">
                        <SelectValue placeholder="Select a branch" />
                      </SelectTrigger>

                      <SelectContent>
                        {branches.map((branch) => (
                          <SelectItem key={branch.name} value={branch.name}>
                            {branch.name}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  )}
                </div>

              </div>

              <Button onClick={saveGitConfig} className="w-fit">
                Save Configuration
              </Button>
            </CardContent>
          </Card>
        )}
      </TabsContent>

      {/* ✅ PUBLIC GIT TAB CONTENT */}
      <TabsContent value="public-git">
        <Card>
          <CardHeader>
            <CardTitle>Public Git Repository</CardTitle>
            <CardDescription>
              Deploy from any public Git repository by providing the URL and branch.
            </CardDescription>
          </CardHeader>

          <CardContent className="space-y-6">
            <div className="flex flex-col md:flex-row gap-6">
              <div className="flex-1">
                <Label className="text-muted-foreground">Git URL</Label>
                <Input
                  className="mt-2"
                  placeholder="https://github.com/user/repo.git"
                  value={publicGitUrl}
                  onChange={(e) => setPublicGitUrl(e.target.value)}
                />
                <p className="text-xs text-muted-foreground mt-1">
                  Enter the full clone URL of your public repository
                </p>
              </div>

              <div className="flex-1">
                <Label className="text-muted-foreground">Branch</Label>
                <Input
                  className="mt-2"
                  placeholder="main"
                  value={publicGitBranch}
                  onChange={(e) => setPublicGitBranch(e.target.value)}
                />
                <p className="text-xs text-muted-foreground mt-1">
                  The branch to deploy from (default: main)
                </p>
              </div>
            </div>

            <Button onClick={savePublicGitConfig} disabled={isSavingPublicGit} className="w-fit">
              {isSavingPublicGit ? "Saving..." : "Save Configuration"}
            </Button>
          </CardContent>
        </Card>
      </TabsContent>

      {/* ✅ OTHER PROVIDERS (Disabled) */}
      <TabsContent value="gitlab">
        <Card>
          <CardHeader>
            <CardTitle>GitLab Support</CardTitle>
            <CardDescription>GitLab integration is coming soon.</CardDescription>
          </CardHeader>
        </Card>
      </TabsContent>

      <TabsContent value="gitea">
        <Card>
          <CardHeader>
            <CardTitle>Gitea Support</CardTitle>
            <CardDescription>Gitea integration is coming soon.</CardDescription>
          </CardHeader>
        </Card>
      </TabsContent>

      <TabsContent value="bitbucket">
        <Card>
          <CardHeader>
            <CardTitle>Bitbucket Support</CardTitle>
            <CardDescription>Bitbucket integration is coming soon.</CardDescription>
          </CardHeader>
        </Card>
      </TabsContent>

    </Tabs>
  )
}
