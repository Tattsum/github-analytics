# frontend

React + Vite SPA for the GitHub team analytics web app.

- Package manager: pnpm (see repo `mise.toml` for toolchain versions)
- GraphQL client: urql, posting to a same-origin `/query`
- Typed operations: graphql-codegen (client preset) reads `../graph/*.graphqls`
  and generates `src/gql/`. Use documents via `graphql(\`...\`)` + `useQuery`.
- Charts: Recharts (wrapped by `src/components/BarChart.tsx`)
- Build output: `frontend/dist` (embedded by the Go binary for same-origin serving)
- Dev: `pnpm dev` runs Vite, which proxies `/query` to the Go server on :8080

## Scripts

- `pnpm dev` — start the Vite dev server
- `pnpm codegen` — regenerate typed GraphQL hooks into `src/gql/`
- `pnpm build` — type-check and produce `dist/`
- `pnpm test` — run unit tests (Vitest)

## Layout

- `src/components/` — shared building blocks (`AppShell`, `BarChart`)
- `src/lib/ranking.ts` — client-side `sortBy` / `rank` helpers (sorting,
  ranking and comparison are computed on the frontend, not in GraphQL)
- `src/pages/` — routed pages (placeholders; Phase 1 fills them in)
- `src/gql/` — generated, committed typed GraphQL output (run `pnpm codegen`)
