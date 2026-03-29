import { Outlet } from "react-router-dom";
import { SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar";
import { AppSidebar } from "@/components/app-sidebar";
import { AppBreadcrumbs } from "@/components/app-breadcrumbs";

export const Layout = () => {
  return (
    <SidebarProvider>
      <AppSidebar />
      <main className="w-full flex flex-col gap-y-4 pb-2 px-2 pt-1 lg:pb-4 lg:px-4 lg:pt-3 xl:pb-6 xl:px-6 xl:pt-4">
        <div className="flex items-center gap-5">
          <SidebarTrigger />
          <AppBreadcrumbs />
        </div>
        <Outlet />
      </main>
    </SidebarProvider>
  );
};
