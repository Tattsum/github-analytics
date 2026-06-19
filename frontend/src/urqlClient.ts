import { Client, cacheExchange, fetchExchange } from "urql";

// The frontend always talks to a same-origin `/query` endpoint. In dev, Vite
// proxies it to the Go server; in production the Go binary serves both the
// embedded SPA and the GraphQL endpoint on the same origin. There is no auth
// (internal/local tool), so no auth exchange is configured.
export const urqlClient = new Client({
  url: "/query",
  exchanges: [cacheExchange, fetchExchange],
});
