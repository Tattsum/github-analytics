import type { ReactNode } from "react";
import { NavLink } from "react-router-dom";

const navLinkStyle = ({ isActive }: { isActive: boolean }) => ({
  padding: "0.5rem 0.75rem",
  borderRadius: "0.375rem",
  fontWeight: isActive ? 600 : 400,
  color: isActive ? "#2563eb" : "#374151",
  backgroundColor: isActive ? "#eff6ff" : "transparent",
});

// AppShell is the shared layout: a top navigation bar plus a content area that
// renders the routed page. Reused by every page so navigation stays consistent.
export function AppShell({ children }: { children: ReactNode }) {
  return (
    <div style={{ minHeight: "100vh" }}>
      <header
        style={{
          display: "flex",
          alignItems: "center",
          gap: "1rem",
          padding: "0.75rem 1.5rem",
          borderBottom: "1px solid #e5e7eb",
          backgroundColor: "#ffffff",
        }}
      >
        <strong style={{ fontSize: "1.05rem" }}>GitHub チーム分析</strong>
        <nav style={{ display: "flex", gap: "0.25rem" }}>
          <NavLink to="/" style={navLinkStyle} end>
            概要
          </NavLink>
          <NavLink to="/repositories" style={navLinkStyle}>
            リポジトリ
          </NavLink>
        </nav>
      </header>
      <main style={{ maxWidth: 1200, margin: "0 auto", padding: "1.5rem" }}>{children}</main>
    </div>
  );
}
