import { SidebarMenu, SidebarMenuButton, SidebarMenuItem } from "@/components/ui/sidebar"
import { Loader2 } from "lucide-react"

export function VersionSwitcher({
  version,
  loading = false,
}: {
  version: string
  loading?: boolean
}) {
  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <SidebarMenuButton
          size="lg"
          className="cursor-default hover:bg-transparent focus-visible:ring-0 focus-visible:ring-offset-0"
        >
          <img src="/mist.png" alt="Mist Icon" className="size-10" />
          <div className="flex flex-col gap-0.5 leading-none">
            <span className="font-medium">Mist</span>
            <span className="text-sm text-muted-foreground flex items-center gap-1">
              {loading ? (
                <>
                  <Loader2 className="h-3 w-3 animate-spin" />
                  <span>Loading...</span>
                </>
              ) : (
                `v${version}`
              )}
            </span>
          </div>
        </SidebarMenuButton>
      </SidebarMenuItem>
    </SidebarMenu>
  )
}
