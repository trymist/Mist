import { VersionSwitcher } from "@/components/version-switcher"
import { useSidebar } from "@/components/ui/sidebar"
import { useVersion } from "@/hooks"


import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarRail,
  SidebarFooter,
} from "@/components/ui/sidebar"
import {
  Home,
  Settings,
  Users,
  GitBranch,
  Book,
  LifeBuoy,
  FolderGit2,
  User,
  // Server,
  LogOut,
  FileText,
  ScrollText,
  RefreshCw,
} from "lucide-react"
import { Link, useLocation } from "react-router-dom"
import { useAuth } from "@/providers"
import { Button } from "@/components/ui/button"

const useNavData = () => {
  const location = useLocation()

  const data = {
    navMain: [
      {
        title: "Home",
        items: [
          {
            title: "Dashboard",
            url: "/",
            icon: Home,
            isActive: location.pathname === "/",
            newTab: false
          },
          {
            title: "Projects",
            url: "/projects",
            icon: FolderGit2,
            isActive: location.pathname.startsWith("/projects"),
            newTab: false
          },
          // {
          //   title: "Deployments",
          //   url: "/deployments",
          //   icon: Server,
          //   isActive: location.pathname === "/deployments",
          // },
          {
            title: "Audit Logs",
            url: "/audit-logs",
            icon: FileText,
            isActive: location.pathname === "/audit-logs",
            newTab: false
          },
          {
            title: "System Logs",
            url: "/logs",
            icon: ScrollText,
            isActive: location.pathname === "/logs",
            newTab: false
          },

        ],
      },
      {
        title: "Settings",
        items: [
          {
            title: "Users",
            url: "/users",
            icon: Users,
            isActive: location.pathname === "/users",
            newTab: false
          },
          {
            title: "System",
            url: "/settings",
            icon: Settings,
            isActive: location.pathname === "/settings",
            newTab: false
          },
          {
            title: "Profile",
            url: "/profile",
            icon: User,
            isActive: location.pathname === "/profile",
            newTab: false
          },
          {
            title: "Git",
            url: "/git",
            icon: GitBranch,
            isActive: location.pathname === "/git",
            newTab: false
          },
        ],
      },
      {
        title: "Extras",
        items: [
          {
            title: "Updates",
            url: "/updates",
            icon: RefreshCw,
            isActive: location.pathname === "/updates",
            newTab: false
          },
          {
            title: "Documentation",
            url: "https://trymist.cloud/guide/what-is-mist.html",
            icon: Book,
            isActive: location.pathname === "/docs",
            newTab: true
          },
          {
            title: "Contribute",
            url: "https://github.com/corecollectives/mist",
            icon: LifeBuoy,
            isActive: location.pathname === "/support",
            newTab: true
          },
        ],
      },
    ],
  }

  return { data }
}
export function AppSidebar(props: React.ComponentProps<typeof Sidebar>) {
  const { data } = useNavData()
  const { user, logout } = useAuth()
  const { state } = useSidebar()
  const { version, loading } = useVersion()

  const isCollapsed = state === "collapsed"

  return (
    <Sidebar {...props} collapsible="icon" variant="floating">
      <SidebarHeader>
        <VersionSwitcher version={version} loading={loading} />
      </SidebarHeader>

      <SidebarContent>
        {data.navMain.map((group) => (
          <SidebarGroup key={group.title}>
            {!isCollapsed && <SidebarGroupLabel>{group.title}</SidebarGroupLabel>}
            <SidebarGroupContent>
              <SidebarMenu>
                {group.items.map((item) => {
                  const Icon = item.icon
                  return (
                    <SidebarMenuItem key={item.title}>
                      <SidebarMenuButton asChild isActive={item.isActive}>
                        <Link to={item.url} referrerPolicy="no-referrer" target={item.newTab ? "_blank" : "_self"} className="flex items-center gap-2">
                          <Icon className="h-4 w-4" />
                          {!isCollapsed && <span>{item.title}</span>}
                        </Link>
                      </SidebarMenuButton>
                    </SidebarMenuItem>
                  )
                })}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        ))}
      </SidebarContent>

      <SidebarFooter className="border-t border-border/40 p-3 flex flex-row items-center justify-between">
        {isCollapsed ? (
          <div className="h-8 w-8 rounded-full overflow-hidden border-2 border-border bg-muted flex items-center justify-center mx-auto">
            {user?.avatarUrl ? (
              <img src={user.avatarUrl} alt="Profile" className="h-full w-full object-cover" />
            ) : (
              <User className="h-4 w-4 text-muted-foreground" />
            )}
          </div>
        ) : (
          <>
            <div className="flex items-center w-full gap-2">
              <div className="h-8 w-8 rounded-full overflow-hidden border-2 border-border bg-muted flex items-center justify-center shrink-0">
                {user?.avatarUrl ? (
                  <img src={user.avatarUrl} alt="Profile" className="h-full w-full object-cover" />
                ) : (
                  <User className="h-4 w-4 text-muted-foreground" />
                )}
              </div>
              <div className="flex flex-col leading-none">
                <span className="text-sm font-medium">{user?.username || "Guest"}</span>
                <span className="text-xs text-muted-foreground truncate max-w-[120px]">
                  {user?.email || "Not signed in"}
                </span>
              </div>
            </div>

            <Button
              variant="ghost"
              size="icon"
              onClick={logout}
              title="Logout"
            >
              <LogOut className="h-4 w-4" />
            </Button>
          </>
        )}
      </SidebarFooter>

      <SidebarRail />
    </Sidebar>
  )
}
