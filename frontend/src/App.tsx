import { Route, Routes } from "react-router-dom";
import { AppShell } from "./components/AppShell";
import { TeamOverview } from "./pages/TeamOverview";
import { MemberDetail } from "./pages/MemberDetail";
import { Repositories } from "./pages/Repositories";
import { RepositoryDetail } from "./pages/RepositoryDetail";
import { NotFoundPage } from "./pages/NotFoundPage";

// Route table for the SPA, wired to the real pages backed by the GraphQL API.
export function App() {
  return (
    <AppShell>
      <Routes>
        <Route path="/" element={<TeamOverview />} />
        <Route path="/members/:login" element={<MemberDetail />} />
        <Route path="/repositories" element={<Repositories />} />
        <Route path="/repositories/:name" element={<RepositoryDetail />} />
        <Route path="*" element={<NotFoundPage />} />
      </Routes>
    </AppShell>
  );
}
