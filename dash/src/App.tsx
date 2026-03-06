import { FullScreenLoading } from "@/components/common";
import { useAuth } from "@/providers";
import { Navigate, Route, BrowserRouter as Router, Routes } from "react-router-dom"
import "./App.css"
import { Toaster } from "./components/ui/sonner";
import ProtectedLayout from "@/app/(protected)/layout";
import DashboardPage from "@/app/(protected)/page";
import LoginPage from "@/app/(auth)/login/page";
import SetupPage from "@/app/(auth)/setup/page";
import UsersPage from "@/app/(protected)/users/page";
import GitPage from "@/app/(protected)/git/page";
import ProjectsPage from "@/app/(protected)/projects/page";
import ProjectPage from "@/app/(protected)/projects/[projectId]/page";
import AppPage from "@/app/(protected)/projects/[projectId]/apps/[appId]/page";
import DatabasesPage from "@/app/(protected)/databases/page";
import LogsPage from "@/app/(protected)/logs/page";
import AuditLogsPage from "@/app/(protected)/audit-logs/page";
import SettingsPage from "@/app/(protected)/settings/page";
import ProfilePage from "@/app/(protected)/profile/page";
import UpdatesPage from "@/app/(protected)/updates/page";
import CallbackPage from "@/app/(protected)/callback/page";

export default function App() {
  const { setupRequired, user } = useAuth();

  if (setupRequired === null) {
    return <FullScreenLoading />;
  }
  return (
    <Router>
      <Routes>
        {setupRequired ? (
          <>
            <Route path="/setup" element={<SetupPage />} />
            <Route path="*" element={<Navigate to="/setup" replace />} />
          </>
        ) : !user ? (
          <>
            <Route path="/login" element={<LoginPage />} />
            <Route path="*" element={<Navigate to="/login" replace />} />
          </>
        ) : (
          <>
            <Route element={<ProtectedLayout />} >
              <Route path="/" element={<DashboardPage />} />
              <Route path="/users" element={<UsersPage />} />
              <Route path="/git" element={<GitPage />} />
              <Route path="/projects" element={<ProjectsPage />} />
              <Route path="/projects/:projectId" element={<ProjectPage />} />
              <Route path="/projects/:projectId/apps/:appId" element={<AppPage />} />
              {/* <Route path="/deployments" element={<DeploymentsPage />} /> */}
              <Route path="/databases" element={<DatabasesPage />} />
              <Route path="/logs" element={<LogsPage />} />
              <Route path="/audit-logs" element={<AuditLogsPage />} />
              <Route path="/settings" element={<SettingsPage />} />
              <Route path="/profile" element={<ProfilePage />} />
              <Route path="/updates" element={<UpdatesPage />} />
              <Route path="/callback" element={<CallbackPage />} />
            </Route>
            <Route path="*" element={<Navigate to="/" replace />} />
          </>
        )}

      </Routes>
      <Toaster />
    </Router>
  )
}
