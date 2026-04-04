import { FullScreenLoading } from "@/components/common";
import { useAuth } from "@/providers";
import { Navigate, Route, BrowserRouter as Router, Routes } from "react-router-dom"
import { Layout } from "./Layout";
import "./App.css"
import { Toaster } from "./components/ui/sonner";
import { CallbackPage } from "./pages/Callback";

import { SetupPage, LoginPage } from "./features/auth";
import { DashboardPage } from "./features/dashboard";
import { UsersPage } from "./features/users";
import { ProjectsPage } from "./features/projects";
import { GitPage } from "./features/git";

// import { DeploymentsPage } from "./pages/Deployments";
import { DatabasesPage } from "./pages/Databases";
import { LogsPage } from "./pages/Logs";
import { SettingsPage } from "./pages/Settings";
import { ProfilePage } from "./pages/Profile";
import { UpdatesPage } from "./pages/Updates";
import { AuditLogsPage } from "./features/auditLogs";
import { ProjectPage } from "./features/projects/ProjectPage";
import { AppPage } from "./features/applications/AppPage";

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
            <Route element={<Layout />} >
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

