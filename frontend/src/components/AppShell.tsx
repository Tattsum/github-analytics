import { useEffect, useRef, useState, type ReactNode } from "react";
import { NavLink } from "react-router-dom";
import { mq } from "../styles/breakpoints";

const NAV_ITEMS: ReadonlyArray<{ to: string; label: string; end?: boolean }> = [
  { to: "/", label: "概要", end: true },
  { to: "/repositories", label: "リポジトリ" },
];

const navLinkStyle = ({ isActive }: { isActive: boolean }) => ({
  padding: "0.5rem 0.75rem",
  borderRadius: "0.375rem",
  fontWeight: isActive ? 600 : 400,
  color: isActive ? "#2563eb" : "#374151",
  backgroundColor: isActive ? "#eff6ff" : "transparent",
});

// AppShell is the shared layout: a top navigation bar plus a content area that
// renders the routed page. Above the mobile breakpoint the nav is a horizontal
// row; at <=768px it collapses into a hamburger that opens a dropdown below the
// header (closed on link tap, outside click, or Escape).
export function AppShell({ children }: { children: ReactNode }) {
  const [menuOpen, setMenuOpen] = useState(false);
  const headerRef = useRef<HTMLElement>(null);

  useEffect(() => {
    if (!menuOpen) {
      return;
    }
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        setMenuOpen(false);
      }
    };
    const onPointerDown = (event: MouseEvent) => {
      if (headerRef.current && !headerRef.current.contains(event.target as Node)) {
        setMenuOpen(false);
      }
    };
    document.addEventListener("keydown", onKeyDown);
    document.addEventListener("mousedown", onPointerDown);
    return () => {
      document.removeEventListener("keydown", onKeyDown);
      document.removeEventListener("mousedown", onPointerDown);
    };
  }, [menuOpen]);

  return (
    <div css={{ minHeight: "100vh" }}>
      <header
        ref={headerRef}
        css={{
          position: "relative",
          display: "flex",
          alignItems: "center",
          gap: "1rem",
          padding: "0.75rem 1.5rem",
          borderBottom: "1px solid #e5e7eb",
          backgroundColor: "#ffffff",
        }}
      >
        <strong css={{ fontSize: "1.05rem" }}>GitHub チーム分析</strong>

        <nav
          css={{
            display: "flex",
            gap: "0.25rem",
            [mq.mobile]: { display: "none" },
          }}
        >
          {NAV_ITEMS.map((item) => (
            <NavLink key={item.to} to={item.to} style={navLinkStyle} end={item.end}>
              {item.label}
            </NavLink>
          ))}
        </nav>

        <button
          type="button"
          aria-label="メニュー"
          aria-expanded={menuOpen}
          onClick={() => setMenuOpen((open) => !open)}
          css={{
            display: "none",
            marginLeft: "auto",
            padding: "0.4rem 0.6rem",
            fontSize: "1.25rem",
            lineHeight: 1,
            color: "#374151",
            backgroundColor: "transparent",
            border: "1px solid #e5e7eb",
            borderRadius: "0.375rem",
            cursor: "pointer",
            [mq.mobile]: { display: "inline-flex" },
          }}
        >
          ☰
        </button>

        {menuOpen && (
          <nav
            css={{
              position: "absolute",
              top: "100%",
              left: 0,
              right: 0,
              display: "flex",
              flexDirection: "column",
              gap: "0.25rem",
              padding: "0.5rem",
              backgroundColor: "#ffffff",
              borderBottom: "1px solid #e5e7eb",
              boxShadow: "0 8px 16px rgba(0, 0, 0, 0.08)",
              zIndex: 10,
            }}
          >
            {NAV_ITEMS.map((item) => (
              <NavLink
                key={item.to}
                to={item.to}
                style={navLinkStyle}
                end={item.end}
                onClick={() => setMenuOpen(false)}
              >
                {item.label}
              </NavLink>
            ))}
          </nav>
        )}
      </header>
      <main
        css={{
          maxWidth: 1200,
          margin: "0 auto",
          padding: "1.5rem",
          [mq.mobile]: { padding: "1rem" },
        }}
      >
        {children}
      </main>
    </div>
  );
}
